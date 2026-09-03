# Inventory V1 Implementation Tasks

Inventory V1 is a pure, additive root-package contract. It does not authorize or include acquisition/discovery/WU-4C, Result V1, W5, CNSIC, publication, sync/archive, Git, CLI, or changes to legacy inventory behavior.

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 735 additions + deletions for contract work; 747 only if the conditional 12-line config-context correction is approved before apply |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | 2A constructor/canonical evidence (365) → 2B strict decoder (370); conditional config-context correction is artifact-only and remains within the applicable slice |
| Delivery strategy | ask-on-risk → owner selected chained delivery |
| Chain strategy | stacked-to-main |

Decision needed before apply: No — owner selected the 2A → 2B chain
Chained PRs recommended: Yes — accepted by owner
Chain strategy: stacked-to-main
400-line budget risk: High — contained by two sub-400 slices

**Apply gate:** the owner resolved the 735-line budget risk by selecting the 2A → 2B stacked-to-main chain. Work units `kozz36/git-change-evidence#38` (2A) and `#39` (2B) are open with `status:approved`. These tasks still do not authorize implementation: interactive approval, fresh native status, and runtime acquire remain required before each apply slice.

### Exact forecast arithmetic

All estimates are additions unless explicitly stated; no production or test deletions are planned. The two implementation slices are independently reviewable and each remains below 400 changed lines.

| Slice | Class | Additions + deletions addends | Subtotal |
|-------|-------|-------------------------------|---------:|
| 2A | production | `inventory_v1.go` 52 + `inventory_v1_build.go` 103 | 155 |
| 2A | tests | `inventory_v1_contract_test.go` 145 + `inventory_v1_ownership_test.go` 65 | 210 |
| 2A | mechanical | 0 | 0 |
| 2A | artifacts | 0 | 0 |
| **2A total** | — | 155 + 210 + 0 + 0 | **365** |
| 2B | production | `inventory_v1_decode.go` 155 | 155 |
| 2B | tests | `inventory_v1_decode_test.go` 150 + `inventory_v1_decode_strict_test.go` 65 | 215 |
| 2B | mechanical | 0 | 0 |
| 2B | artifacts | 0 | 0 |
| **2B total** | — | 155 + 215 + 0 + 0 | **370** |
| Conditional pre-apply context correction | production | 0 | 0 |
| Conditional pre-apply context correction | tests | 0 | 0 |
| Conditional pre-apply context correction | mechanical | 0 | 0 |
| Conditional pre-apply context correction | artifacts | `openspec/config.yaml`: 6 additions + 6 deletions | 12 |
| **Total without / with conditional correction** | — | 365 + 370 / 365 + 370 + 12 | **735 / 747** |

### Slice limits and stop rule

- 2A changes four root-package Go files; 2B changes three. Each is within the eight-changed-Go-files-per-package limit.
- Every planned Go file is at most 155 changed lines, below the 160-line per-Go-file limit. The design's former 170/190-line file targets are refined downward here by separating ownership and strict-vector tests, not by compressing production code or tests.
- If necessary evidence makes either slice exceed 400 changes or a Go file exceed 160, stop after this honest 2A/2B split, report the revised arithmetic, and request a delivery-size decision. Do not code-golf tests, docs, or error assertions.

## Preconditions and Evidence Rules

- **Strict TDD:** for each behavior below, record the named failing RED command and failure reason before production code; then record the GREEN, TRIANGULATE, and REFACTOR focused-test results. Each negative assertion uses `errors.As` for `*ContractError` and records the expected `Field` and `Code`.
- **Runtime harness:** N/A for both slices. This is a pure root-package value/decoder contract with no executable, filesystem, network, Git, CLI, goroutine, or adapter boundary. The focused Go tests are the runtime evidence.
- **Race command:** N/A. Neither slice creates, synchronizes, or depends on concurrent execution; do not run `go test -race ./...` unless implementation changes that fact.
- **Legacy fence:** do not edit `inventory.go`, `internal/inventory/**`, `policy_*.go`, or their tests. A new compatibility smoke test must exercise `NewUntrackedInventory`, `Records`, and `InventoryUnavailableError` unchanged; the full suite is the adapter compatibility evidence.
- **Stale config context:** `openspec/config.yaml` historically says that no Go source, Go tests, or `go.mod` exist. That prose is stale; its strict-TDD commands and layout/threshold rules remain operative. Before apply, decide whether a bounded context/baseline correction is needed for truthful apply evidence. If needed, request it as the separately estimated artifact-only `openspec/config.yaml` change above; it may correct only stale context/baseline text, not behavior, test commands, layout rules, or scope. Do not edit it in this phase or fold it into either feature behavior.

## Work Unit 2A — Construct Immutable Canonical Evidence

**Start boundary:** approved work unit `kozz36/git-change-evidence#38`, approved stacked-to-main delivery, and no `inventory_v1*.go` files exist. <br>
**File targets:** `inventory_v1.go`, `inventory_v1_build.go`, `inventory_v1_contract_test.go`, and `inventory_v1_ownership_test.go` only. <br>
**Finish boundary:** `NewInventoryV1` produces the sole canonical V1 document from a revalidated `PolicyDocumentV1`, and its accessors expose defensively owned derived evidence. Strict wire decoding remains absent and is deferred to 2B. <br>
**Spec coverage:** Contract—Neutral root-package API, Closed wire contract (canonical fixture and empty inventory), Policy antecedent binding (all construction scenarios), Verified-empty scope, Additive compatibility; Content Identity—all construction-side scenarios under content identity, raw-path fidelity, ordering, and defensive ownership.

### 2A TDD checklist

- [x] **2A.1 RED — public canonical API and antecedent link.** In `inventory_v1_contract_test.go`, add `TestNewInventoryV1CanonicalContract` using `canonicalHighBitPathFixture`: assert the exact schema, padded Base64 `/wph`, final LF, independently calculated Inventory digest, and a typed `accounting_policy`/`invalid_document` failure for zero policy. Run `go test . -run '^TestNewInventoryV1CanonicalContract$'`; record the non-zero failure before production changes. *(Specs: Contract—Construct and inspect a neutral document; Canonical byte fixture; Invalid zero policy document; Same policy yields a reproducible antecedent link.)*
- [x] **2A.2 GREEN — implement only the public document surface and construction/link path.** Add `InventoryContractV1`, `InventoryEntryInput`, `InventoryEntryV1`, `InventoryDocumentV1`, `NewInventoryV1`, and copy-returning accessors in `inventory_v1.go`; add canonical ordered-wire encoding, final-LF document SHA-256, and `validatedPolicyDigest` revalidation in `inventory_v1_build.go`. Run `go test . -run '^TestNewInventoryV1CanonicalContract$'`; record exit 0 and the exact canonical fixture/digest observation. *(Specs: same as 2A.1; no caller-supplied inventory/content identity API.)*
- [x] **2A.3 TRIANGULATE — prove byte-identity, not a trusted policy digest.** Extend `TestNewInventoryV1CanonicalContract` with `policyLinkFixture`: independently built byte-identical policies must yield identical links/bytes/digests, while a byte-different valid policy produces a distinct link and document. Assert `InventoryEntryInput` has no digest or length input fields. Re-run `go test . -run '^TestNewInventoryV1CanonicalContract$'`; record exit 0. *(Specs: Contract—Caller supplies only an opaque inventory or content hash; Same policy yields a reproducible antecedent link; Policy byte identity changes.)*
- [x] **2A.4 REFACTOR — isolate policy validation without changing the surface.** Refactor only private helpers in `inventory_v1.go` and `inventory_v1_build.go` to keep antecedent validation, canonical encoding, and copying explicit; retain the fixture assertions unchanged. Re-run `go test . -run '^TestNewInventoryV1CanonicalContract$'`; record exit 0. *(Specs: Contract—Neutral root-package API and Policy V1 antecedent binding.)*

- [x] **2A.5 RED — exact content identity.** In `inventory_v1_contract_test.go`, add `TestNewInventoryV1ExactContentAndPaths` with `exactContentBytesFixture`: exact SHA-256/length for `0x00,0xff,0x0d,0x0a`, valid zero-length content, and distinct `line\n` versus `line\r\n` evidence. Run `go test . -run '^TestNewInventoryV1ExactContentAndPaths$'`; record the non-zero failure. *(Specs: Content—Exact content digest and length; Similar text with different bytes; Proposed digest or length conflicts with content.)*
- [x] **2A.6 GREEN — derive evidence from copied content only.** In `inventory_v1_build.go`, copy each input path/content at ingress, derive SHA-256 and `uint64(len(content))` from the copied content, retain no content, and expose only `Path`, `ContentSHA256`, and `ByteLength`. Run `go test . -run '^TestNewInventoryV1ExactContentAndPaths$'`; record exit 0. *(Specs: Content—System-derived exact-content identity and Content is not retained as mutable public state.)*
- [x] **2A.7 TRIANGULATE — preserve raw paths and deterministic exact sets.** Extend `TestNewInventoryV1ExactContentAndPaths` with `rawByteOrderingFixture`: reverse versus sorted `a`, `a/\xff`, `a0`, and `\xff` inputs must have identical canonical bytes/digests and required raw-byte order; reject duplicate identical paths with same and different content and every invalid form (empty, absolute, empty component, dot, dot-dot, NUL). Include valid byte-distinct invalid-UTF-8, case, backslash, and control-byte paths. Re-run the focused command; record exit 0 and typed `entries.path` codes. *(Specs: Content—Raw path-byte fidelity and validity; Byte-distinct paths are not normalized together; Reordered equivalent caller inputs; Prefix and non-UTF-8 sort case; Duplicate path with conflicting contents.)*
- [x] **2A.8 REFACTOR — make byte-only path handling reviewable.** Refactor only private raw-path validation/sort/duplicate helpers in `inventory_v1_build.go`; retain `[]byte` handling and `bytes.Compare`, with no string/path normalization. Re-run `go test . -run '^TestNewInventoryV1ExactContentAndPaths$'`; record exit 0. *(Specs: Content—Raw path-byte fidelity and Deterministic exact-path set and ordering.)*

- [x] **2A.9 RED — empty, zero, ownership, and additive coexistence.** In `inventory_v1_ownership_test.go`, add `TestInventoryV1OwnershipEmptyAndLegacyCoexistence`: valid empty input must encode `"entries":[]` and a non-nil empty view; zero/failed documents must expose nil/empty evidence; ingress and egress mutation attempts must not change observations; legacy `NewUntrackedInventory`, `Records`, and `InventoryUnavailableError` remain unchanged. Run `go test . -run '^TestInventoryV1OwnershipEmptyAndLegacyCoexistence$'`; record the non-zero failure. *(Specs: Contract—Empty canonical inventory, Verified empty requested scope, Acquisition failure is not empty inventory, and Legacy inventory coexistence; Content—Caller mutates construction buffers, Caller mutates returned inventory data, Content is not retained as mutable public state.)*
- [x] **2A.10 GREEN — enforce defensive ownership and verified-empty semantics.** Complete copy-on-egress behavior for `CanonicalBytes()` and `Entries()` in `inventory_v1.go`, ensure valid empty uses an allocated entries slice and has no status/availability variant, and do not touch legacy files. Run `go test . -run '^TestInventoryV1OwnershipEmptyAndLegacyCoexistence$'`; record exit 0. *(Specs: same as 2A.9; Contract—Out-of-scope consumer integration.)*
- [x] **2A.11 TRIANGULATE — mutate every mutable boundary.** Extend the ownership test to mutate/reslice the caller entry slice, path/content buffers, returned canonical bytes, returned entry slice, and returned nested paths; assert later canonical bytes, order, policy link, and both digests remain original. Re-run the focused command; record exit 0. *(Specs: Content—Caller mutates construction buffers and Caller mutates returned inventory data.)*
- [x] **2A.12 REFACTOR and close 2A.** Refactor only the new 2A files for named fixtures and private copy helpers; then record these successful closure commands:
  ```bash
  go test . -run '^(TestNewInventoryV1CanonicalContract|TestNewInventoryV1ExactContentAndPaths|TestInventoryV1OwnershipEmptyAndLegacyCoexistence)$'
  go test ./...
  test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"
  ```
  Record runtime harness **N/A** and race test **N/A** with the precondition rationale above. *(Specs: all 2A-mapped requirements; Contract—Additive compatibility and historical-authority isolation.)*

**2A rollback boundary:** delete only the four new `inventory_v1*` files listed above. This removes the new construction contract and its tests without changing `PolicyDocumentV1`, `inventory.go`, legacy acquisition, or any consumer integration.

## Work Unit 2B — Reject Every Noncanonical Wire Form

**Start boundary:** approved work unit `kozz36/git-change-evidence#39`; 2A is present, green, and merged to main; no Inventory V1 decoder exists. <br>
**File targets:** `inventory_v1_decode.go`, `inventory_v1_decode_test.go`, and `inventory_v1_decode_strict_test.go` only. <br>
**Finish boundary:** `DecodeInventoryV1` validates the exact expected policy, accepts only constructor-equivalent canonical bytes, returns no partial document on rejection, and leaves 2A behavior unchanged. <br>
**Spec coverage:** Contract—Closed wire contract (alternate representation), Policy binding (wrong link), Strict validating decode (all rejection classes); Content Identity—non-UTF-8 round trip, invalid encoded paths, duplicate paths, and decoder noncanonical order.

### 2B TDD checklist

- [ ] **2B.1 RED — canonical decode and closed object shape.** In `inventory_v1_decode_test.go`, add `TestDecodeInventoryV1CanonicalAndClosedShape`: decode the 2A high-bit canonical fixture with its exact policy, assert a byte-identical path/view/document, then table-drive truncated input, concatenated documents, non-UTF-8 transport, non-object root, `entries` type drift, missing fields, unknown fields, and duplicate root/entry fields including escaped duplicate names. Run `go test . -run '^TestDecodeInventoryV1CanonicalAndClosedShape$'`; record the non-zero failure. *(Specs: Contract—Alternate wire representation; Duplicate or unknown field; Realistic malformed transport input; Content—Non-UTF-8 raw path round trip.)*
- [ ] **2B.2 GREEN — implement duplicate-aware closed decode.** Add `DecodeInventoryV1` and local streaming/token-scanner helpers in `inventory_v1_decode.go`; reject non-UTF-8, non-object/trailing JSON, duplicate/unknown/missing fields, wrong types, unsupported schema, and all authority-bearing normalized field names at root and entry levels. Revalidate the supplied policy before parsing. Run `go test . -run '^TestDecodeInventoryV1CanonicalAndClosedShape$'`; record exit 0 and each asserted `ContractError` field/code. *(Specs: Contract—Strict validating decode and Invalid zero policy document.)*
- [ ] **2B.3 TRIANGULATE — exercise authority claims and canonical layout.** Extend the closed-shape table with `approval`, `rejection`, `blocking`, `gating`, `merge`, `deploy`, `release`, and `delivery` at both levels, plus root/entry field reorder, whitespace, absent final LF, and escaped-equivalent string cases. Re-run the focused command; record exit 0 and `authority_field`/`noncanonical` evidence. *(Specs: Contract—Duplicate or unknown field; Alternate wire representation; Strict validating decode.)*
- [ ] **2B.4 REFACTOR — retain closed parsing locally.** Refactor only private scanner/type helpers in `inventory_v1_decode.go`; do not alter the Policy V1 decoder or generic legacy decoder. Re-run `go test . -run '^TestDecodeInventoryV1CanonicalAndClosedShape$'`; record exit 0. *(Specs: Contract—Strict validating decode; Additive compatibility and historical-authority isolation.)*

- [ ] **2B.5 RED — strict entry values and decoded path validity.** In `inventory_v1_decode_test.go`, add `TestDecodeInventoryV1ValueValidation` with `base64VariantFixture`, `rawPathViolationFixture`, `numericCoercionOverflowFixture`, and `digestShapeFixture`: URL/unpadded/whitespace/escaped Base64; every invalid raw path form; signed/fraction/exponent/string/leading-zero/overflow lengths; and uppercase/40-character content or policy digests. Run `go test . -run '^TestDecodeInventoryV1ValueValidation$'`; record the non-zero failure. *(Specs: Contract—Closed wire contract and Realistic malformed transport input; Content—Invalid path form.)*
- [ ] **2B.6 GREEN — validate wire values without coercion or normalization.** In `inventory_v1_decode.go`, require padded standard Base64 that round-trips byte-identically, validate decoded raw paths as bytes, require lowercase 64-hex digests, and parse only canonical ASCII-decimal `uint64` lengths. Run `go test . -run '^TestDecodeInventoryV1ValueValidation$'`; record exit 0 and the asserted field/code for every vector. *(Specs: same as 2B.5; Content—Raw path-byte fidelity and validity.)*
- [ ] **2B.7 TRIANGULATE — retain opaque valid bytes.** Extend the value test with a valid invalid-UTF-8/control/high-bit path and a zero-length evidence entry, and prove the decoded view is byte-identical while no text normalization occurs. Re-run the focused command; record exit 0. *(Specs: Content—Non-UTF-8 raw path round trip; Byte-distinct paths are not normalized together; System-derived exact-content identity.)*
- [ ] **2B.8 REFACTOR — centralize only decoder-local value checks.** Refactor private Base64, path, digest, and uint validation helpers without broadening any accepted spelling. Re-run `go test . -run '^TestDecodeInventoryV1ValueValidation$'`; record exit 0. *(Specs: Contract—Closed Inventory V1 wire contract and Strict validating decode.)*

- [ ] **2B.9 RED — order, duplicate, policy replay, and exact re-encoding.** In `inventory_v1_decode_strict_test.go`, add `TestDecodeInventoryV1OrderLinkAndCanonicality` with `duplicatePathAndOrderFixture`, `wrongPolicyReplayFixture`, and `noncanonicalPresentationFixture`: duplicate paths with same/different contents, reverse raw-byte order, a well-formed link from another valid policy, and semantically equivalent noncanonical presentation. Run `go test . -run '^TestDecodeInventoryV1OrderLinkAndCanonicality$'`; record the non-zero failure. *(Specs: Contract—Wrong but well-formed policy link and Alternate wire representation; Content—Decoder receives noncanonical entry order and Duplicate path with conflicting contents.)*
- [ ] **2B.10 GREEN — reject rather than repair.** In `inventory_v1_decode.go`, require strictly ascending raw-byte paths without sorting, reject duplicates, recompute/compare the expected policy link from the supplied canonical policy, rebuild canonical bytes, and require byte equality before returning a document. Run `go test . -run '^TestDecodeInventoryV1OrderLinkAndCanonicality$'`; record exit 0 and `duplicate_path`, `noncanonical_order`, `mismatch`, and `noncanonical` evidence. *(Specs: Contract—Policy V1 antecedent binding and Strict validating decode; Content—Deterministic exact-path set and ordering.)*
- [ ] **2B.11 TRIANGULATE — prove round-trip plus no partial escape.** Extend the strict table so each rejection asserts the returned `InventoryDocumentV1{}` has nil canonical bytes/entries and empty digests; add the positive high-bit constructor → decode → view round trip. Re-run the focused command; record exit 0. *(Specs: Contract—Neutral root-package Inventory V1 API and Strict validating decode; Content—Non-UTF-8 raw path round trip.)*
- [ ] **2B.12 REFACTOR and close 2B.** Refactor only the three new 2B files for named negative fixture classes and shared test assertions; then record these successful closure commands:
  ```bash
  go test . -run '^(TestDecodeInventoryV1CanonicalAndClosedShape|TestDecodeInventoryV1ValueValidation|TestDecodeInventoryV1OrderLinkAndCanonicality)$'
  go test ./...
  test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"
  ```
  Record runtime harness **N/A** and race test **N/A** with the precondition rationale above. *(Specs: all 2B-mapped requirements; Contract—Additive compatibility and historical-authority isolation.)*

**2B rollback boundary:** delete only `inventory_v1_decode.go`, `inventory_v1_decode_test.go`, and `inventory_v1_decode_strict_test.go`. 2A remains a valid construction-only contract; no legacy, policy, acquisition, publication, sync/archive, Result V1, W5, CNSIC, CLI, or Git behavior is reverted or changed.

## Completion Evidence Ledger

Before marking either work unit complete, record in apply evidence:

- the RED command, non-zero result, and named missing/rejected behavior for every numbered RED task;
- the GREEN, TRIANGULATE, and REFACTOR focused command result for the same named test;
- the 2A and 2B closure outputs for `go test ./...` and the check-only gofmt command;
- runtime harness **N/A** and `go test -race ./...` **N/A**, each with the pure/no-concurrency rationale; and
- the per-slice `git diff --stat` arithmetic, changed root Go-file count, per-file changed-line maximum, rollback boundary, and confirmation that legacy targets were unmodified.
