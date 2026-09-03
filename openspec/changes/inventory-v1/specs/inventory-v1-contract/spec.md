# Inventory V1 Contract Specification

## Purpose

Define a neutral, immutable, canonical Inventory V1 document in the public root Go package. It is the second document in the provenance sequence after Accounting Policy V1 and before any future Accounting Result V1.

## Requirements

### Requirement: Neutral root-package Inventory V1 API

The system MUST provide a public, versioned Inventory V1 API in the module's root Go package. The API MUST let a caller construct an Inventory V1 document from (1) an exact validated Accounting Policy V1 document and (2) the caller's exact sequence of raw-path/content entry values. The API MUST also expose a validating decoder that accepts canonical Inventory V1 bytes and the exact expected Accounting Policy V1 document.

The API MUST expose, for a valid Inventory V1 document, its canonical bytes, system-derived document digest, linked policy digest, and ordered entry view. The entry construction boundary MUST accept content bytes and raw path bytes; an API that accepts only UTF-8 strings or a caller-provided content digest or byte length does not satisfy this requirement. A zero or failed document value MUST NOT present canonical bytes, a digest, a policy link, or entries as though it were valid.

The API MUST be usable as a standalone generic integrity-library contract. It MUST NOT expose CNSIC-specific vocabulary, a project-specific compatibility hash, filesystem discovery, descriptor handling, safe acquisition, publication, synchronization, archive, Git, CLI, Result V1, W5, or delivery-authority behavior.

#### Scenario: Construct and inspect a neutral document
- **GIVEN** a validated canonical Accounting Policy V1 document and two caller-provided raw-path/content pairs
- **WHEN** the caller constructs and then inspects Inventory V1 through the public root-package API
- **THEN** the caller receives an immutable Inventory V1 document with its canonical bytes, document digest, derived policy link, and ordered entry view, without using a filesystem, Git, CNSIC, or a caller-supplied identity assertion

#### Scenario: Caller supplies only an opaque inventory or content hash
- **GIVEN** a caller has an opaque inventory hash or a proposed entry content hash but does not supply the corresponding entry content bytes
- **WHEN** the caller attempts Inventory V1 construction
- **THEN** construction fails and no Inventory V1 provenance document is emitted

#### Scenario: Invalid zero policy document
- **GIVEN** a caller provides a zero, invalid, or otherwise noncanonical Accounting Policy V1 value
- **WHEN** the caller attempts Inventory V1 construction or validating decode
- **THEN** the operation fails and does not treat the policy value, its digest, or any Inventory V1 output as authoritative

### Requirement: Closed Inventory V1 wire contract

Inventory V1 canonical bytes MUST encode exactly one UTF-8 JSON object with this closed V1 shape and field order:

```json
{"schema":"git-change-evidence.inventory/v1","accounting_policy_sha256":"<64 lowercase hex>","entries":[{"path_b64":"<canonical standard base64>","content_sha256":"<64 lowercase hex>","byte_length":<unsigned decimal integer>}]}
```

The root object MUST contain exactly `schema`, `accounting_policy_sha256`, and `entries`, in that order. Every entry object MUST contain exactly `path_b64`, `content_sha256`, and `byte_length`, in that order. The `schema` value MUST be exactly `git-change-evidence.inventory/v1`. The `entries` value MUST be an array, including `[]` for an empty inventory, never `null`. The canonical JSON object MUST contain no insignificant whitespace and MUST end with one LF (`\n`) byte. All digest fields MUST contain exactly 64 lowercase hexadecimal characters.

`path_b64` MUST use padded RFC 4648 standard Base64 for the original raw path bytes. It MUST NOT use URL Base64, omitted padding, embedded whitespace, or a noncanonical encoding. `byte_length` MUST be a JSON unsigned integer with no sign, fraction, exponent, leading zero, overflow, coercion, or string representation.

#### Scenario: Canonical byte fixture
- **GIVEN** a valid canonical policy `P`, its independently computed digest `D = SHA-256(P.CanonicalBytes())`, and one entry with raw path bytes `0xff 0x0a 0x61`, content bytes `x`, SHA-256 content digest `2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881`, and byte length `1`
- **WHEN** the Inventory V1 document is constructed from `P`
- **THEN** its canonical bytes are exactly the UTF-8 byte concatenation `{"schema":"git-change-evidence.inventory/v1","accounting_policy_sha256":"` + `D` + `","entries":[{"path_b64":"/wph","content_sha256":"2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881","byte_length":1}]}\n`, with `D` replaced by its 64 lowercase hexadecimal characters, and its document digest is SHA-256 over exactly those bytes

#### Scenario: Empty canonical inventory
- **GIVEN** a validated Accounting Policy V1 document and a verified requested scope containing no entries
- **WHEN** the caller constructs Inventory V1 with an empty entry sequence
- **THEN** the canonical document contains `"entries":[]`, contains neither `null` nor an availability/status field, and has a system-derived digest

#### Scenario: Alternate wire representation
- **GIVEN** otherwise valid Inventory V1 semantics encoded with a different root or entry field order, leading or trailing whitespace, no final LF, an escaped equivalent JSON string, URL/unpadded Base64, `null` entries, or a numeric string for `byte_length`
- **WHEN** those bytes are submitted to the validating decoder
- **THEN** decoding fails without rewriting or normalizing the input

### Requirement: Policy V1 antecedent binding

Inventory V1 construction and validating decode MUST require the exact canonical Accounting Policy V1 antecedent document, not merely a purported digest. The system MUST compute the antecedent link as SHA-256 over that policy's exact canonical bytes and MUST encode that derived value as `accounting_policy_sha256`.

For validating decode, the system MUST independently derive the expected policy digest from the supplied Policy V1 document and MUST reject the Inventory V1 bytes when their `accounting_policy_sha256` differs. A syntactically valid digest, a digest from a different valid policy document, or a caller-provided digest without the exact policy document MUST NOT establish the link.

#### Scenario: Same policy yields a reproducible antecedent link
- **GIVEN** two independently constructed but byte-identical valid Accounting Policy V1 documents and the same Inventory V1 entry inputs
- **WHEN** each is used to construct Inventory V1
- **THEN** both inventories contain the same derived `accounting_policy_sha256`, canonical bytes, and Inventory V1 document digest

#### Scenario: Wrong but well-formed policy link
- **GIVEN** canonical Inventory V1 bytes whose `accounting_policy_sha256` has been replaced with the valid 64-lowercase-hex digest of a different Accounting Policy V1 document
- **WHEN** the bytes are decoded with the actual expected policy document
- **THEN** decoding fails and returns no Inventory V1 document

#### Scenario: Policy byte identity changes
- **GIVEN** two valid Policy V1 documents with different canonical bytes but otherwise similar semantic intent
- **WHEN** identical entry inputs are used to construct Inventory V1 from each policy
- **THEN** the resulting inventories retain their respective derived policy digests and MUST NOT be treated as interchangeable provenance documents

### Requirement: Strict validating decode

The Inventory V1 decoder MUST be closed and validating. It MUST reject malformed JSON; non-UTF-8 input; a non-object root; trailing JSON values; missing, unknown, or duplicate fields at every object level; invalid JSON types; unsupported schema; invalid Base64; invalid path values; duplicate paths; malformed digest values; invalid byte lengths; noncanonical entry order; an incorrectly linked policy; and any input whose bytes differ from the canonical re-encoding of the validated document.

The decoder MUST NOT silently add, drop, coerce, reorder, repair, deduplicate, normalize, or otherwise reinterpret input. It MUST reject consumer-specific and authority-bearing fields, including approval, rejection, blocking, gating, merge, deploy, release, and delivery claims, wherever such fields occur.

#### Scenario: Duplicate or unknown field
- **GIVEN** canonical Inventory V1 bytes modified to duplicate `path_b64` in one entry or to add `approval`, `unknown`, or `delivery` at the document or entry level
- **WHEN** the bytes are submitted to the validating decoder with the correct policy document
- **THEN** decoding fails with typed contract evidence and returns no partial document

#### Scenario: Realistic malformed transport input
- **GIVEN** a producer transmits a truncated JSON document, two concatenated JSON documents, a non-UTF-8 byte sequence, a boolean in place of `entries`, a signed or overflowing `byte_length`, or a 40-character/uppercase content digest
- **WHEN** the receiver validates the input as Inventory V1
- **THEN** validation fails before the document is accepted or exposed

### Requirement: Verified-empty scope and acquisition separation

An Inventory V1 document with zero entries MUST mean that the caller has successfully verified the exact requested scope and found it empty. It MUST NOT mean unavailable, unsupported, missing, unsafe, racing, failed, or unrequested acquisition. The pure Inventory V1 document and API MUST have no availability/status variant and MUST NOT translate an acquisition failure into an empty Inventory V1 document.

Inventory V1 construction MUST consume only the exact entry set supplied by its caller. It MUST NOT enumerate directories, open paths, inspect the filesystem, follow links, infer omitted paths, augment the set, or define discovery semantics. Safe filesystem acquisition and adapter integration remain outside this contract and are deferred to WU-4C.

#### Scenario: Verified empty requested scope
- **GIVEN** an adapter or caller has completed acquisition of an explicitly requested scope and verified that it contains no inventory entries
- **WHEN** it invokes the pure Inventory V1 constructor with that empty exact set and a valid policy
- **THEN** the result is a valid empty Inventory V1 document and denotes only that verified empty scope

#### Scenario: Acquisition failure is not empty inventory
- **GIVEN** acquisition of an explicitly requested path is unavailable because it is missing, unsafe, racing, unsupported, or otherwise unverifiable
- **WHEN** the acquisition boundary reports that outcome
- **THEN** it remains an unavailable outcome outside Inventory V1, and neither the adapter nor the pure contract converts it into an empty Inventory V1 document

### Requirement: Additive compatibility and historical-authority isolation

Inventory V1 MUST be additive. It MUST NOT change the public behavior, data model, ordering, defensive-copy semantics, or unavailable-outcome behavior of `UntrackedRecord`, `UntrackedInventory`, `NewUntrackedInventory`, `Records`, `InventoryUnavailableError`, or existing internal inventory acquisition adapters.

This change MUST NOT alter, supersede, reinterpret, or claim completion under the historical `bootstrap-neutral-core-extraction` authority. Accounting Result V1, CNSIC integration, publication, sync/archive, W5, and Result V1 are outside this Inventory V1 contract.

#### Scenario: Legacy inventory coexistence
- **GIVEN** existing code constructs a legacy `UntrackedInventory` with raw Go-string paths and observes its records or an `InventoryUnavailableError`
- **WHEN** Inventory V1 is added
- **THEN** the legacy code retains its existing behavior unchanged and does not need to construct, decode, or consume Inventory V1

#### Scenario: Out-of-scope consumer integration
- **GIVEN** a caller requests CNSIC inventory hashing, Result V1 assembly, publication, archive synchronization, or filesystem acquisition through Inventory V1
- **WHEN** the request is evaluated
- **THEN** it is outside this contract and no such integration or authority behavior is supplied
