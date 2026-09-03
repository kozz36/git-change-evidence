# Design: Inventory V1 Canonical Contract

## Decision

Add a pure, additive Inventory V1 document contract in the root `changeevidence` package. It accepts only a validated `PolicyDocumentV1` and raw path/content pairs, derives all hashes and lengths itself, and emits or accepts exactly one canonical JSON representation. It does not alter `inventory.go`, legacy inventory acquisition, or any Policy V1 API.

The constructor boundary is also the acquisition seam: an empty input means the caller has already verified its requested scope is empty. There is deliberately no status, availability, filesystem, or discovery API in this contract.

## Public API

Add the following root-package declarations; all fields not shown are private implementation state.

```go
const InventoryContractV1 ContractVersion = "git-change-evidence.inventory/v1"

type InventoryEntryInput struct {
    Path    []byte
    Content []byte
}

type InventoryEntryV1 struct {
    Path          []byte
    ContentSHA256 Digest
    ByteLength    uint64
}

type InventoryDocumentV1 struct {
    // private: canonical, digest, accountingPolicyDigest, entries
}

func NewInventoryV1(policy PolicyDocumentV1, entries []InventoryEntryInput) (InventoryDocumentV1, error)
func DecodeInventoryV1(raw []byte, policy PolicyDocumentV1) (InventoryDocumentV1, error)

func (i InventoryDocumentV1) CanonicalBytes() []byte
func (i InventoryDocumentV1) Digest() DocumentDigest
func (i InventoryDocumentV1) AccountingPolicyDigest() DocumentDigest
func (i InventoryDocumentV1) Entries() []InventoryEntryV1
```

`InventoryEntryInput` intentionally has no digest, length, text, or project-compatibility fields. `InventoryEntryV1` intentionally has no content field: the output exposes the path and derived evidence only.

### Zero-value and failure behavior

| Value or call | Required behavior |
|---|---|
| Zero `InventoryDocumentV1` | `CanonicalBytes()` is `nil`; `Digest()` and `AccountingPolicyDigest()` are `""`; `Entries()` is `nil`. It represents no document. |
| Valid empty document | Canonical bytes contain `"entries":[]`; both digests are populated; `Entries()` returns a non-`nil`, zero-length slice. |
| Failed constructor or decoder | Returns `InventoryDocumentV1{}` and a non-`nil` typed error. No partial entries, link, canonical bytes, or digest escape. |
| `CanonicalBytes()` | Returns a fresh copy on every call. |
| `Entries()` | Returns a fresh slice and fresh copy of every `Path` on every call. `ContentSHA256` and `ByteLength` are scalar evidence values. |

`DocumentDigest` remains the existing document-digest type. The per-entry hash uses the existing generic `Digest` type and is named `ContentSHA256` so it is not confused with the Inventory document digest.

## Canonical document and provenance binding

The sole V1 wire shape is encoded with an ordered private struct, then `json.Marshal`, followed by exactly one `\n` byte:

```json
{"schema":"git-change-evidence.inventory/v1","accounting_policy_sha256":"<64 lowercase hex>","entries":[{"path_b64":"<padded RFC 4648 standard base64>","content_sha256":"<64 lowercase hex>","byte_length":<uint64>}]}
```

The struct declaration fixes root order as `schema`, `accounting_policy_sha256`, `entries`, and entry order as `path_b64`, `content_sha256`, `byte_length`. Its entries slice is always allocated, so a valid empty inventory encodes as `[]`, never `null`.

Before either construction or decoding accepts a policy, an internal `validatedPolicyDigest(policy)` routine will:

1. obtain a defensive copy with `policy.CanonicalBytes()`;
2. reject an empty value and pass the bytes to shipped `DecodeAccountingPolicyV1`;
3. require the decoded policy's canonical bytes to equal the supplied bytes;
4. independently SHA-256 those exact canonical bytes, compare the result with both `policy.Digest()` and the decoded policy's `Digest()`; and
5. return that recomputed `DocumentDigest` as the only allowed `accounting_policy_sha256` value.

Thus the link is based on the exact antecedent bytes, not on a caller assertion. Construction serializes that returned digest; decoding recomputes it from the supplied expected policy and compares it with the parsed link before a document can be returned. No Inventory V1 constructor or decoder accepts a policy digest in lieu of the document.

The Inventory document digest is SHA-256 of the complete canonical Inventory bytes, including the final LF. It is computed only after canonical encoding and is never caller supplied.

## Entry data flow, validation, and ownership

### Construction

1. Validate the exact Policy V1 antecedent as above.
2. Allocate an owned entry sequence. For each input, copy `Path` and copy `Content` at ingress.
3. Validate the copied raw path bytes; compute SHA-256 and `uint64(len(content))` from the copied content bytes; then discard the temporary content copy. The retained document state contains only copied path bytes, `ContentSHA256`, and `ByteLength`.
4. Sort retained entries with `bytes.Compare` on raw paths. Reject byte-identical adjacent paths after sorting; do not merge them, even if their content is identical.
5. Canonically encode the ordered, validated entries, derive the document digest, and retain private copies only.

This deliberately makes NUL, CR/LF differences, invalid UTF-8, and high-bit content bytes part of content identity. It also prevents post-construction mutation of the caller's entry slice, path buffers, or content buffers from changing the document. Content is not kept in `InventoryDocumentV1`, its canonical form, or its view.

### Raw-path rule

Path validation operates only on `[]byte`, never through `string`, `path.Clean`, Unicode processing, or a platform path API. A path is valid precisely when it is nonempty and relative; does not start with `/`; contains no NUL; and has no empty, `.` or `..` slash-delimited component. `/` is the only separator. Backslashes, controls other than NUL, invalid UTF-8, and high-bit bytes are opaque component bytes and remain valid and byte-distinct.

Sorting is raw-byte lexicographic `bytes.Compare`; a shorter prefix sorts first. Therefore `a`, `a/\xff`, `a0`, `\xff` order exactly in that sequence, independent of locale, UTF-8, or caller input order.

## Strict validating decode

`DecodeInventoryV1` is validating, not reparative. It first validates the supplied expected policy, then executes this sequence:

1. Reject non-UTF-8 input before JSON parsing.
2. Parse the root and each entry with an `encoding/json.Decoder` token scanner that requires an object, records decoded member names, rejects duplicate keys (including escaped spellings), rejects unknown keys, and requires decoder EOF after the object. A map-only `json.Unmarshal` is not sufficient because it loses duplicate keys.
3. Require all root and entry fields, reject incorrect JSON types, and parse `entries` as an array only. `null` is never an array.
4. Decode each `path_b64` with `base64.StdEncoding`; require `StdEncoding.EncodeToString(decoded) == supplied` before validating the resulting raw path. This rejects URL encoding, missing padding, whitespace, and alternate spellings.
5. Require both digest strings to satisfy the existing `isDigest` lowercase-hex rule. Parse `byte_length` from its `json.RawMessage` only if every byte is an ASCII decimal digit, it has no nonzero leading-zero form, and `strconv.ParseUint(..., 10, 64)` succeeds. This rejects signs, fractions, exponents, strings, coercions, and overflow.
6. Preserve parsed path/digest/length values in internal evidence entries. Validate duplicate paths and strict raw-byte ascending order without sorting. Compare the parsed policy link with the independently derived expected-policy digest.
7. Rebuild canonical Inventory V1 bytes from those validated values and require `bytes.Equal(reencoded, raw)`. Any different field order, whitespace, escape spelling, omitted LF, or other alternate representation fails rather than being normalized.

The decoder can validate a transmitted entry digest's form, ordering, and canonical document membership, but cannot rederive that digest because content is intentionally not carried on the wire. Exact-content derivation is therefore enforced at the constructor boundary; strict decoder acceptance is of its closed canonical evidence representation and exact policy antecedent link.

Authority-bearing unknown fields are rejected at both levels before canonical equality. The inventory scanner will use the existing authority-word convention plus the Policy V1 extensions, recognizing approval/approve, rejection/reject, denial/deny, blocking/block, gating/gate, merge, deploy, release, delivery/delivery-authority after case and `_`/`-`/`.` normalization.

## Typed contract errors

Use the existing `*ContractError` type exclusively, with dot-qualified `Field` values and lower-snake-case `Code` values. No Inventory-specific error type is introduced.

| Condition | `Field` | `Code` |
|---|---|---|
| Zero, invalid, noncanonical, or internally inconsistent supplied policy | `accounting_policy` | `invalid_document` |
| Invalid constructor raw path | `entries.path` | `invalid_path` |
| Duplicate constructor raw path | `entries.path` | `duplicate_path` |
| Non-UTF-8 transport | `document` | `invalid` |
| Malformed, non-object, or trailing JSON | `document` or `entries` | `invalid_object` |
| Missing root or entry field | `document.<field>` or `entries.<field>` | `required` |
| Duplicate root or entry field | `document.<field>` or `entries.<field>` | `duplicate_field` |
| Unrecognized field | `document.<field>` or `entries.<field>` | `unknown_field` |
| Authority-bearing field | `document.<field>` or `entries.<field>` | `authority_field` |
| Wrong schema | `schema` | `unsupported_version` |
| Non-array `entries` | `entries` | `invalid_array` |
| Invalid `path_b64` type or Base64 | `entries.path_b64` | `invalid` or `invalid_base64` |
| Decoded invalid path | `entries.path_b64` | `invalid_path` |
| Invalid content or policy digest spelling | corresponding wire field | `invalid_digest` |
| Invalid or overflowing length | `entries.byte_length` | `invalid_uint` |
| Duplicate decoded raw path | `entries.path_b64` | `duplicate_path` |
| Decoded entries not strictly raw-byte ordered | `entries` | `noncanonical_order` |
| Parsed policy link differs from supplied Policy V1 | `accounting_policy_sha256` | `mismatch` |
| Any remaining semantic-equivalent but byte-different representation | `document` | `noncanonical` |

The constructor only emits construction-relevant rows; the decoder emits the wire rows. Tests will assert errors through `errors.As(err, *ContractError)` and assert both field and code, following current root-package contract tests.

## File decomposition

Only new root-package files are needed. Existing `inventory.go`, all `policy_*.go` files, legacy adapter files, and their tests remain unchanged.

| File | Responsibility | Target size |
|---|---|---:|
| `inventory_v1.go` | Public constant, input/view/document types, constructor and accessor signatures | <= 90 lines |
| `inventory_v1_build.go` | Policy validation/link derivation, defensive copying, raw-path validation, entry derivation/order checks, canonical encoding | <= 170 lines |
| `inventory_v1_decode.go` | Streaming duplicate-aware closed-object decode, wire type parsing, Base64/uint checks, policy-link and exact-reencode checks | <= 190 lines |
| `inventory_v1_contract_test.go` | Canonical construction, identity, ordering, empty, zero-value, ownership, and legacy-coexistence tests | <= 220 lines |
| `inventory_v1_decode_test.go` | Strict decoder positive/negative transport and authority-field vectors | <= 260 lines |

The decoder keeps its duplicate-aware scanner local rather than changing the existing Policy V1 decoder or generic legacy decoder. It may reuse package-private digest and authority helpers where their current behavior already matches, but must not weaken or alter any existing contract.

## Strict-TDD verification strategy

### Canonical fixtures

Use deterministic, named fixture classes rather than ad hoc hashes:

- **`canonicalHighBitPathFixture`**: Policy fixture plus path `[]byte{0xff, '\n', 'a'}`, content `x`, expected exact JSON and independently calculated Inventory digest.
- **`verifiedEmptyFixture`**: Valid policy and allocated empty entry input; expects `"entries":[]`, a nonempty document digest, and non-`nil` empty entry view.
- **`exactContentBytesFixture`**: `bin/data` with `0x00, 0xff, 0x0d, 0x0a`, plus zero-length content and `line\n`/`line\r\n` variants to establish byte-exact hash and length behavior.
- **`rawByteOrderingFixture`**: Unsorted `a`, `a/\xff`, `a0`, and `\xff`; expects the required raw-byte view and identical bytes for pre-sorted versus reverse input.
- **`policyLinkFixture`**: Two independently created byte-identical policies, then a byte-different valid policy; establishes reproducible and non-interchangeable links.
- **`ownershipFixture`**: Mutable outer entry slice, paths, content, returned paths, and canonical bytes; mutation after construction or egress must not change later observations.

### Realistic negative fixtures

`TestDecodeInventoryV1FailsClosed` will table-drive at least these named classes, with the exact expected `ContractError` field and code:

- **`transportTruncationFixture`**, **`concatenatedDocumentsFixture`**, and **`nonUTF8TransportFixture`**;
- **`rootOrEntriesTypeDriftFixture`** for array/null/boolean/string substitution;
- **`duplicateAndMissingFieldFixture`** at root and entry levels, including escaped duplicate field names;
- **`authorityInjectionFixture`** for `approval`, `rejection`, `blocking`, `gating`, `merge`, `deploy`, `release`, and `delivery` at both levels;
- **`base64VariantFixture`** for URL alphabet, unpadded, whitespace, and escaped-but-equivalent forms;
- **`rawPathViolationFixture`** for empty, absolute, empty-component, dot, dot-dot, and NUL paths;
- **`duplicatePathAndOrderFixture`** for same-path/same-content, same-path/different-content, and out-of-order raw bytes;
- **`numericCoercionOverflowFixture`** for signed, fractional, exponent, string, leading-zero, and overflowing lengths;
- **`digestShapeFixture`** for uppercase and 40-character content/policy digests;
- **`wrongPolicyReplayFixture`** for a well-formed link from a different valid Policy V1 document; and
- **`noncanonicalPresentationFixture`** for field reordering, whitespace, no final LF, and JSON escape spelling changes.

The constructor test suite will also prove that no caller-provided digest or length can influence the output: the public input type has no such field, and calculated expected values are compared to output from the supplied content alone.

Follow strict TDD in two loops: write one focused failing fixture; make the smallest production change to pass it; triangulate with a byte-distinct or malformed companion; then refactor only while the same focused tests remain green. Run focused root tests while developing each loop and `go test ./...` at the end. Existing legacy inventory and adapter tests are run unchanged as compatibility evidence; no new acquisition test is introduced here.

## Preliminary 2A/2B work units (planning only)

These are cohesive, test-with-code slices sized below the 400-line review budget. They are not tasks, commits, or authorization to deliver.

| Unit | Cohesive behavior and files | Preliminary review budget | Completion and rollback boundary |
|---|---|---:|---|
| **2A — construct immutable canonical evidence** | Add `inventory_v1.go`, `inventory_v1_build.go`, and `inventory_v1_contract_test.go`. Cover policy validation/link derivation, exact content hash/length, raw-path rules and sorting, canonical bytes/digest, empty and zero behavior, defensive ownership, and an unchanged legacy API smoke test. | about 330–370 lines | Focused construction tests plus `go test ./...`; remove only these new Inventory V1 construction files to roll back. |
| **2B — reject every noncanonical wire form** | Add `inventory_v1_decode.go` and `inventory_v1_decode_test.go`. Cover duplicate-aware closed decode, all realistic-negative fixture classes, policy replay rejection, and canonical decode round trip. | about 320–380 lines | Focused decode tests plus `go test ./...`; remove only the new decoder and decoder tests, leaving the 2A constructor contract intact. |

If actual test evidence makes either slice exceed 400 changed lines, stop after this one honest split and ask for a delivery-size decision rather than compressing tests or merging the work units.

## Out-of-scope acquisition seam and rollback

`NewInventoryV1(policy, entries)` consumes exactly the entries handed to it. A future WU-4C-style adapter may obtain those raw paths and content only after its own safe-acquisition checks, and may call the constructor with an empty slice only after it has verified an empty requested scope. Missing, symlink, special, racing, unsupported, or unverifiable acquisition remains an external unavailable outcome (including existing `InventoryUnavailableError` behavior); it is never converted to Inventory V1.

Rollback is confined to deleting the new Inventory V1 source and focused tests. It does not modify or roll back Accounting Policy V1, `UntrackedRecord`, `UntrackedInventory`, `InventoryUnavailableError`, internal adapters, publication, Result V1, CLI, Git behavior, or the historical bootstrap authority.

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| A permissive JSON path silently accepts alternate bytes | Token-based duplicate-aware closed parsing plus exact canonical re-encode equality. |
| Raw paths are normalized accidentally | Keep paths as bytes from ingress through comparison, validation, encoding, and egress; test high-bit/control cases. |
| A policy digest is trusted without its document | Require and revalidate `PolicyDocumentV1`, recompute SHA-256 from its canonical bytes, and compare the link. |
| Mutable content or views alter evidence | Copy mutable ingress, retain no content, and copy every bytes-based egress. |
| Empty inventory is mistaken for unavailable acquisition | Carry no status variant and document the verified-empty caller precondition at the pure-function seam. |
| Contract implementation disrupts legacy APIs | Add isolated root files only; leave `inventory.go` and adapters untouched and run their existing tests unchanged. |
