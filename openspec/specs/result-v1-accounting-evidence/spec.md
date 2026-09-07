# Result V1 Accounting Evidence Specification

## Purpose

Define the committed-snapshot scope, exclusive classification, recomputable totals, and factual observation model carried by Result V1.

## Requirements

### Requirement: Valid immutable committed-snapshot boundary

Result V1 construction MUST consume only a `CommittedSnapshot` of committed changes. It MUST NOT derive, augment, or replace its entries from Git, refs, the index, the working tree, the filesystem, untracked Inventory V1 records, or any other acquisition source.

The constructor MUST reject a snapshot unless its `Base` and `Head` are immutable native Git object IDs: each MUST be exactly 40 or exactly 64 lowercase hexadecimal characters, and both MUST have the same native width. Symbolic names, abbreviated IDs, uppercase or mixed-case IDs, malformed IDs, and cross-width revision pairs are invalid. A base equal to head is permitted; it represents an immutable empty comparison only when the supplied entry sequence is empty.

For Result V1 accounting, each committed entry's primary `Path` MUST be interpreted as the exact byte sequence of that Go string, without UTF-8 validation, Unicode normalization, case folding, locale collation, separator replacement, path cleaning, or conversion. It MUST be a nonempty relative slash-separated path with no NUL, leading slash, empty component, `.` component, or `..` component. Entries with duplicate primary path bytes are invalid. The current `Path`, not `PreviousPath`, MUST determine classification. Status, object ID, mode, kind, and binary metadata MUST NOT alter a classification, a total, or an observation; `Lines.Countable` exclusively determines whether its additions/deletions are countable.

A Result V1 constructor MUST retain only defensively owned evidence derived from the supplied immutable snapshot. It MUST NOT claim to prove snapshot acquisition, object existence, repository membership, moving-reference stability, or working-tree isolation; those are acquisition-boundary concerns outside Result V1.

#### Scenario: Valid native SHA-256 snapshot with hostile path bytes
- **GIVEN** a committed snapshot with distinct 64-character lowercase-hex base and head IDs and an entry whose primary path string contains raw bytes `0xff 0x0a 0x61`
- **WHEN** Result V1 is constructed
- **THEN** the entry is eligible for accounting, its exact bytes are preserved through Base64 evidence, and neither text conversion nor acquisition is performed

#### Scenario: Invalid revision or path input
- **GIVEN** a snapshot with an abbreviated revision, a symbolic ref, mixed 40- and 64-character IDs, a nonempty entry sequence with equal base and head IDs, an empty path, `/absolute`, `a//b`, `a/../b`, a NUL path, or duplicate byte-identical primary paths
- **WHEN** Result V1 construction is requested
- **THEN** construction fails fatally and produces no document

#### Scenario: Rename classification uses the committed primary path
- **GIVEN** a valid snapshot entry has `PreviousPath` `docs/old.txt`, primary `Path` `src/new.go`, and Policy V1 assigns those paths to different categories
- **WHEN** Result V1 is constructed
- **THEN** the entry appears once and is classified from the exact bytes of `src/new.go`; its previous path does not create another classification or contribute another total

### Requirement: Inventory is provenance-only and excluded from accounting

A valid exact Inventory V1 antecedent is required for every Result V1 document, but its entries are provenance evidence only. Inventory V1 content bytes, content digests, byte lengths, paths, entry count, empty/nonempty state, and any untracked-file meaning MUST NOT contribute a committed entry classification, category total, non-countable count, threshold actual, ratio numerator, ratio denominator, availability decision, or exceeded value.

Result V1 MUST account every and only the committed entries in the supplied snapshot. A valid empty Inventory V1 remains a valid required antecedent; it does not mean Result V1 has no committed entries. Conversely, a nonempty Inventory V1 does not produce a Result V1 entry when the committed snapshot is empty.

#### Scenario: Inventory identity changes without changing accounting
- **GIVEN** the same valid policy and committed snapshot, plus two valid inventories linked to that policy that have different untracked entries and therefore different Inventory V1 digests
- **WHEN** a Result V1 document is constructed with each inventory
- **THEN** the committed entries, categories, totals, and observations are identical, while each Result V1 retains its respective inventory link and therefore has a distinct canonical identity

#### Scenario: Nonempty inventory with empty committed snapshot
- **GIVEN** a valid nonempty Inventory V1 and a valid committed snapshot with zero entries
- **WHEN** Result V1 is constructed
- **THEN** it contains no committed classifications and no inventory entry is converted into a line total or non-countable count

### Requirement: Exclusive classification and recomputable policy-category totals

Result V1 MUST classify every valid committed snapshot entry exactly once. It MUST apply the exact raw-byte glob semantics, category precedence, and configured default category of the validated Policy V1 document: scan categories in their Policy V1 order, assign the first category with a matching glob, and assign the policy default category when none matches. A path MUST NOT appear in more than one classification, be omitted, be classified from a normalized path, or be assigned by a caller-supplied category.

The `entries` array MUST be sorted by raw primary path bytes in ascending lexicographic byte order. For each Policy V1 category, `totals` MUST contain exactly one record in Policy V1 category order, including categories with all zero values. Its `additions` is the exact sum of countable classified-entry additions, its `deletions` is the exact sum of countable classified-entry deletions, and its `non_countable` is the exact count of classified entries with `Lines.Countable == false`. Non-countable entries MUST use the `non_countable` measurement variant and MUST contribute neither their additions nor deletions to any line total, even if their input line fields are nonzero.

Classifications, totals, and observations emitted by construction or accepted by decoding MUST equal the values determined by the exact antecedents and entry evidence. Construction and decoding MUST reject a candidate whose classifications, totals, or observations are inconsistent with those values. Any overflow while accumulating a stored category addition or deletion, or while computing a required line-threshold actual or ratio numerator or denominator, is fatal: no valid Result V1 document may be emitted or accepted.

#### Scenario: Overlap, default, and non-countable evidence
- **GIVEN** Policy V1 orders `raw` before `source`, has matching raw-byte globs for both categories, and has `other` as default; the snapshot contains a raw-byte `.go` path matching both, `src/main.go`, `assets/blob` with `Countable == false`, and unmatched `note.txt`
- **WHEN** Result V1 is constructed
- **THEN** the overlapping raw path is classified only as `raw`, `src/main.go` is classified only as `source`, `assets/blob` has a `non_countable` measurement in its one assigned category, `note.txt` is assigned only to `other`, and every total is exactly recomputable from those four classifications

#### Scenario: Constructor input order does not establish result order
- **GIVEN** equivalent snapshots list valid unique paths in reverse raw-byte order and in sorted raw-byte order, including `a`, `a/` followed by byte `0xff`, `a0`, and byte `0xff`
- **WHEN** Result V1 is constructed from each
- **THEN** both entry arrays are ordered `a`, `a/` followed by byte `0xff`, `a0`, byte `0xff`, and their canonical bytes and digests are identical

#### Scenario: Realizable arithmetic overflow is fatal rather than invented evidence
- **GIVEN** two countable entries classified to one category have additions of `MaxUint64` and `1`, respectively, or a required line-threshold or ratio line total includes one countable entry with additions of `MaxUint64` and deletions of `1`
- **WHEN** Result V1 construction or strict decoding recomputes the evidence
- **THEN** it fails fatally with no result document; it does not wrap, saturate, emit an approximate value, or mark the overflow as a successful total

### Requirement: Closed factual threshold, ratio, and unavailable-observation model

Result V1 MUST derive observations only from the exact validated Policy V1 document and recomputed category totals. It MUST emit line-threshold observations first, one for each Policy V1 category with a nonzero line threshold in policy category order. It MUST then emit ratio observations, one for each Policy V1 ratio in that policy's declared ratio order. No other observation kind, reordered observation, duplicate observation, omitted required observation, or observation for a zero line threshold is valid.

For a line threshold, the `limit` MUST equal that category's Policy V1 threshold. It is available only if that category has no non-countable entries; when available, `actual` MUST equal its additions plus deletions and `exceeded` MUST be true exactly when `actual > limit`. When unavailable, the unavailable variant MUST be used and `exceeded` MUST be false.

For a ratio, `category`, `reference`, and decimal-string `limit` MUST equal the corresponding Policy V1 ratio exactly. It is available only when both category totals have no non-countable entries, both line totals can be represented, and the reference line total is nonzero. When available, `numerator` and `denominator` MUST be the exact respective line totals, and `exceeded` MUST be true exactly when the rational `numerator / denominator` is greater than the Policy V1 decimal limit, with no floating-point rounding, textual reformatting, or invented approximation. When unavailable, the unavailable variant MUST be used and `exceeded` MUST be false.

A non-countable entry, an unavailable threshold, an unavailable ratio, and an exceeded threshold or ratio are embedded factual evidence, not construction failure, approval, denial, warning verdict, gate, block, or delivery decision. Result V1 MUST NOT replace them with a generic error/diagnostic array, fictional zero, estimated measurement, or authority-bearing status. Malformed or mismatched antecedents, invalid snapshot/revision/path input, invalid classifications or recomputations, and arithmetic overflow are fatal and MUST produce no Result V1 document.

#### Scenario: Available threshold and exact ratio
- **GIVEN** a category total of 2 additions and 1 deletion, a line threshold of 2, a reference category total of 4 additions and 0 deletions, and a Policy V1 ratio limit `0.5`
- **WHEN** Result V1 observations are derived
- **THEN** the threshold observation is available with `actual:3` and `exceeded:true`, and the ratio observation is available with `numerator:3`, `denominator:4`, and `exceeded:true`

#### Scenario: Exact decimal comparison intentionally differs at a legacy float64 precision edge
- **GIVEN** countable category and reference line totals of numerator `9007199254740993` and denominator `18014398509481984`, respectively, no non-countable entries in either category, and a Policy V1 ratio limit `0.5`
- **WHEN** Result V1 derives the available ratio observation
- **THEN** it reports `numerator:9007199254740993`, `denominator:18014398509481984`, and `exceeded:true` from exact rational comparison, because the numerator is exactly one greater than half the denominator
- **AND** legacy `Account` may report `exceeded:false` after historical `float64` rounding; that intentional precision-edge difference does not weaken Result V1 exact Policy V1 decimal semantics or require any legacy API change

#### Scenario: Non-countable and zero-denominator observations remain embedded
- **GIVEN** a threshold category has one non-countable classified entry, and a configured ratio has a countable zero-line reference category
- **WHEN** Result V1 is constructed
- **THEN** the corresponding threshold and ratio use their unavailable variants with `exceeded:false`, the non-countable entry remains explicit, construction succeeds, and no fictional actual or ratio is emitted

#### Scenario: Exceeded evidence is not authority
- **GIVEN** a valid policy threshold or ratio is exceeded by committed accounting evidence
- **WHEN** Result V1 is constructed or decoded
- **THEN** it emits only the factual `exceeded:true` observation and does not emit approval, rejection, pass/fail, gate, blocking, merge, deployment, release, or delivery semantics

### Requirement: Legacy accounting compatibility

For equivalent valid committed snapshots and equivalent policy categories, Result V1 accounting outcomes MUST be observably equivalent to legacy `Account` outcomes for precedence, default assignment, countable/non-countable treatment, totals, threshold availability and strict greater-than behavior, ratio availability, accumulation failure behavior, and ratio comparisons in the shared domain where legacy `Account` historical `float64` comparison has the same mathematical result as exact Policy V1 decimal rational comparison. Result V1 MUST use exact Policy V1 decimal rational comparison for every available ratio, even where that intentionally differs at a `float64` precision edge. `Account`, `AccountingPolicy`, `AccountingResult`, and their public observable behavior MUST remain compatible; callers of `Account` MUST NOT be required to construct, decode, or consume Result V1.

#### Scenario: Legacy caller remains unaffected
- **GIVEN** existing code calls `Account` with a valid legacy policy and committed snapshot
- **WHEN** Result V1 support is present
- **THEN** the caller receives the same legacy totals, observations, errors, ordering, and public data behavior as before and need not provide Policy V1 or Inventory V1

#### Scenario: Result and legacy accounting agree on shared-domain vectors
- **GIVEN** a Policy V1 document representing a valid legacy accounting policy and a committed snapshot containing first-match, default, countable, non-countable, totals, threshold, ordinary ratio, and overflow vectors whose legacy `float64` ratio comparisons have the same mathematical result as exact Policy V1 decimal rational comparison
- **WHEN** Result V1 and legacy `Account` account their equivalent inputs
- **THEN** their accounting outcomes agree, while Result V1 additionally supplies exact antecedent links, per-entry evidence, canonical bytes, and a digest

## Validation

> **Validated against code:** 2026-09-07 against branch `test/result-v1-unit9-mutation-census` at `17854e7539ca571107a01b35907695a979f1295c`. Formal verification [`verify-report.md`](../../changes/archive/2026-09-07-result-v1/verify-report.md) SHA-256 `b2d41ca63bef0d6fa9e378a7d64d12ebe021fa708b895b35a0c0df1067f0d2ae` reports **PASS WITH WARNINGS**: 13/13 requirements, 38/38 scenarios, and successful focused, full, build, coverage, format, and diff checks.
