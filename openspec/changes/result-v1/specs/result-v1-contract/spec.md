# Result V1 Contract Specification

## Purpose

Define a neutral, immutable, canonical Result V1 document in the public root package. Result V1 is the third provenance document in the sequence **Accounting Policy V1 → Inventory V1 → Accounting Result V1**.

## ADDED Requirements

### Requirement: Neutral Result V1 construction and inspection

The system MUST provide a public, versioned Result V1 contract that constructs a document only from all of the following typed inputs:

1. an exact, validated `PolicyDocumentV1`;
2. an exact, validated `InventoryDocumentV1`; and
3. a valid immutable `CommittedSnapshot`.

The constructor MUST expose neither a caller-supplied policy, inventory, or result digest as an authority boundary nor the legacy `AccountingPolicy` as its public policy input. For a valid document, it MUST expose defensively owned canonical bytes, a system-derived Result V1 document digest, the derived policy and inventory links, immutable revisions, ordered entry classifications, category totals, and observations.

A zero value, a value returned with an error, or any otherwise invalid Result V1 value MUST expose no canonical bytes, digest, antecedent link, revisions, entries, totals, or observations as a valid document.

The Result V1 functional-core contract MUST NOT invoke Git, resolve a ref, inspect the index or working tree, acquire a snapshot or inventory, access the filesystem, parse external configuration, publish, project, operate a CLI, or introduce CNSIC or delivery-authority behavior.

#### Scenario: Construct a neutral result from exact typed evidence
- **GIVEN** an exact validated Policy V1 document, an exact validated Inventory V1 document linked to that policy, and a valid immutable committed snapshot
- **WHEN** the caller constructs and inspects Result V1 through the public root-package API
- **THEN** it receives an immutable canonical document, its derived digest, exact antecedent links, snapshot revisions, classifications, totals, and observations without any Git, filesystem, CLI, profile, or caller-asserted digest operation

#### Scenario: Caller offers only digest assertions
- **GIVEN** a caller supplies a policy digest, inventory digest, result digest, or separately shaped legacy `AccountingPolicy` but omits the required exact antecedent documents
- **WHEN** it attempts Result V1 construction or validating decode
- **THEN** the operation fails and no Result V1 document is emitted

#### Scenario: Failed construction has no partial result
- **GIVEN** a valid policy and inventory but a snapshot with a malformed revision or a category accumulation that overflows
- **WHEN** Result V1 construction is requested
- **THEN** construction returns a fatal error and the returned result exposes no partial classifications, totals, observations, canonical bytes, or digest

### Requirement: Closed Result V1 wire contract and content identity

Result V1 canonical bytes MUST be exactly one UTF-8 JSON object with this closed root shape and root-field order:

```json
{"schema":"git-change-evidence.accounting-result/v1","accounting_policy_sha256":"<64 lowercase hex>","inventory_sha256":"<64 lowercase hex>","revisions":{"base":"<40 or 64 lowercase hex>","head":"<same native width lowercase hex>"},"entries":[<entry>],"totals":[<total>],"observations":[<observation>]}
```

The schema value MUST be exactly `git-change-evidence.accounting-result/v1`. The root MUST contain exactly the seven fields shown, in the shown order. It MUST contain no insignificant whitespace and MUST end in exactly one LF (`\n`) byte. `entries`, `totals`, and `observations` MUST always be arrays, including `[]`; they MUST NOT be `null`.

Each `entries` element MUST be exactly one of these discriminated shapes, with fields in the stated order:

```json
{"path_b64":"<padded RFC 4648 standard Base64>","category":"<exact Policy V1 category name>","measurement":{"kind":"countable","additions":<unsigned integer>,"deletions":<unsigned integer>}}
{"path_b64":"<padded RFC 4648 standard Base64>","category":"<exact Policy V1 category name>","measurement":{"kind":"non_countable"}}
```

Each `totals` element MUST be exactly:

```json
{"category":"<exact Policy V1 category name>","additions":<unsigned integer>,"deletions":<unsigned integer>,"non_countable":<unsigned integer>}
```

Each `revisions` object MUST contain exactly `base` and `head`, in that order. A line-threshold observation MUST be exactly one of:

```json
{"kind":"line_threshold","category":"<exact Policy V1 category name>","limit":<unsigned integer>,"available":true,"actual":<unsigned integer>,"exceeded":<boolean>}
{"kind":"line_threshold","category":"<exact Policy V1 category name>","limit":<unsigned integer>,"available":false,"exceeded":false}
```

A ratio observation MUST be exactly one of:

```json
{"kind":"ratio","category":"<exact Policy V1 category name>","reference":"<exact Policy V1 category name>","limit":"<exact canonical Policy V1 decimal>","available":true,"numerator":<unsigned integer>,"denominator":<nonzero unsigned integer>,"exceeded":<boolean>}
{"kind":"ratio","category":"<exact Policy V1 category name>","reference":"<exact Policy V1 category name>","limit":"<exact canonical Policy V1 decimal>","available":false,"exceeded":false}
```

The `kind` value is the complete observation vocabulary. Result V1 MUST NOT encode a generic diagnostic/status object, a verdict, a synthetic aggregate total, a caller-supplied identity, or any field beyond these closed shapes.

All SHA-256 fields MUST contain exactly 64 lowercase hexadecimal characters. Every unsigned integer MUST be a base-10 JSON integer with no sign, fraction, exponent, leading zero, overflow, coercion, or string representation. `path_b64` MUST use padded RFC 4648 standard Base64 of the original path bytes; URL Base64, omitted padding, embedded whitespace, and alternate encodings are invalid. All category and ratio-limit strings MUST be byte-for-byte values from the validated Policy V1 antecedent. Canonical JSON string escaping MUST be reproduced exactly by canonical re-encoding; an escaped but semantically equivalent alternative is noncanonical.

The Result V1 document digest MUST be SHA-256 over exactly these canonical bytes, including the final LF. The system MUST derive it itself; it MUST NOT serialize the digest into the Result V1 body or accept it from a caller.

#### Scenario: Canonical hostile-path fixture
- **GIVEN** an exact policy with policy digest `P`, a linked exact inventory with digest `I`, base revision `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`, head revision `bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb`, and one countable committed entry whose exact path bytes are `0xff 0x0a 0x66 0x69 0x6c 0x65`, category is `source`, additions are `2`, and deletions are `1`
- **WHEN** Result V1 is constructed and that policy has no enabled threshold or ratio observations
- **THEN** its canonical bytes are exactly `{"schema":"git-change-evidence.accounting-result/v1","accounting_policy_sha256":"` + `P` + `","inventory_sha256":"` + `I` + `","revisions":{"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},"entries":[{"path_b64":"/wpmaWxl","category":"source","measurement":{"kind":"countable","additions":2,"deletions":1}}],"totals":[{"category":"source","additions":2,"deletions":1,"non_countable":0}],"observations":[]}\n`, with `P` and `I` replaced by their 64 lowercase hexadecimal values, and its digest is SHA-256 over exactly those bytes

#### Scenario: Empty committed change set
- **GIVEN** valid exact antecedents and a valid immutable snapshot with distinct valid base and head IDs but zero committed entries
- **WHEN** Result V1 is constructed
- **THEN** it contains `"entries":[]`, contains one all-zero total for each policy category in policy order, includes only observations required by the policy, and has a derived digest

#### Scenario: Alternate wire representation
- **GIVEN** otherwise valid Result V1 semantics encoded with a different field order, whitespace, absent final LF, an escaped-equivalent string, `null` for an array, a URL/unpadded Base64 path, a numeric string, or a floating-point observation value
- **WHEN** the bytes are submitted to the validating decoder
- **THEN** decoding fails without normalizing, rewriting, or accepting the alternative

### Requirement: Defensive ownership and deterministic construction

Result V1 MUST ensure that mutation of caller-owned snapshot data after construction, mutation of bytes or views returned from a Result V1 document, or mutation of independently returned antecedent views does not alter the valid Result V1 document, any subsequently returned view or bytes, its classifications, totals, observations, links, or digest.

For equivalent valid inputs, construction MUST sort classifications by ascending lexicographic comparison of their original raw path bytes, with a shorter prefix before its extension. Construction MUST derive category material in the Policy V1 category order and observations in the closed order specified by the accounting-evidence requirements. Caller entry order MUST NOT affect canonical bytes or digest.

#### Scenario: Caller mutates a source or returned buffer
- **GIVEN** a caller constructs a snapshot from a mutable entry slice, constructs Result V1 from that snapshot, and receives canonical bytes and a classification view
- **WHEN** it mutates the original entry slice after construction, or mutates the returned canonical bytes, classification slice, or nested path bytes
- **THEN** subsequent Result V1 views and bytes from the already constructed document remain byte-identical to their original values

#### Scenario: Equivalent snapshot orders
- **GIVEN** two otherwise equivalent valid snapshots contain the same committed entries in opposite input orders
- **WHEN** each is constructed with the same exact antecedents
- **THEN** both Result V1 documents contain entries in the same raw-byte path order and have byte-identical canonical bytes and the same document digest

### Requirement: Strict antecedent-aware canonical decoding

The validating Result V1 decoder MUST accept candidate Result V1 bytes only together with the exact expected `PolicyDocumentV1` and exact expected `InventoryDocumentV1`. Acceptance requires exact canonical Policy V1 and Inventory V1 antecedents, verification of the Policy V1→Inventory V1→Result V1 links, and exact equality of each link with the digest determined by its linked antecedent. Any invalid or mismatched antecedent or link MUST cause decoding to fail.

The decoder MUST reject malformed JSON; non-UTF-8 input; a non-object root; trailing JSON values; unsupported schema; wrong JSON type; a missing, unknown, or duplicate field at every object level; malformed digest, revision, Base64, decimal, boolean, or unsigned-integer value; invalid path; duplicate path; unsupported discriminator; noncanonical field or array ordering; a policy/inventory link mismatch; an invalid antecedent; a classification that does not follow the policy; totals or observations that do not exactly recompute; a noncanonical re-encoding; and an authority-bearing field name or structural metadata.

The decoder MUST NOT add, drop, coerce, reorder, sort, deduplicate, repair, normalize, infer, or otherwise reinterpret candidate bytes. It MUST reject an unknown field name or structural metadata, including an object key or discriminator value, that introduces approval, rejection, denial, gating, blocking, merge, deployment, release, delivery, or authority semantics, including case- and separator-variants. It MUST evaluate values of recognized fields only under their closed-shape and antecedent constraints; the occurrence of one of those words in a valid Policy V1-derived category or another permitted string value MUST NOT by itself cause rejection. On every fatal decoder failure it MUST return no partial valid Result V1 value.

#### Scenario: Wrong exact antecedent despite valid syntax
- **GIVEN** canonical Result V1 bytes linked to a valid policy `P1` and inventory `I1`, plus different valid policy `P2` or inventory `I2` with a well-formed 64-character digest
- **WHEN** the bytes are decoded with `P2`, `I2`, or an inventory not linked to the supplied policy
- **THEN** decoding fails and returns no Result V1 document even if a supplied digest has valid syntax

#### Scenario: Recomputed evidence differs
- **GIVEN** otherwise structurally valid Result V1 bytes whose entry category is not the first matching Policy V1 category, whose total omits one entry, whose unavailable observation contains an actual value, or whose ratio numerator, denominator, availability, or exceeded value differs from recomputation
- **WHEN** the bytes are decoded with the exact antecedents
- **THEN** decoding fails without correcting the evidence for the caller

#### Scenario: Unknown, duplicate, and authority-bearing nested fields
- **GIVEN** canonical Result V1 bytes modified to duplicate `path_b64`, add `unknown` to a measurement, add an `approval` field to a total, add a `delivery_authority` field to an observation, or add a `merge-gate` field to revisions
- **WHEN** the bytes are submitted to the validating decoder with the exact antecedents
- **THEN** decoding fails and returns no partial document; the unknown or authority-bearing field names are rejected rather than accepted as neutral metadata

#### Scenario: Policy-derived authority-like string remains evidence
- **GIVEN** exact antecedents whose validated Policy V1 document defines a category named `release` or `merge`, and otherwise valid Result V1 bytes use that exact category value in their permitted `category` fields
- **WHEN** the bytes are submitted to the validating decoder with the exact antecedents
- **THEN** the category value is evaluated as policy-derived evidence and is not rejected solely because it contains authority vocabulary
