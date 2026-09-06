# Tasks: Result/Accounting V1 Canonical Contract

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 4,495–5,125 additions + deletions: production 1,450–1,710; tests 1,640–1,995; existing planning artifacts 1,380; future apply evidence 25–40 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | Selected `stacked-to-main` chain below; every PR targets `main` only after its predecessor lands |
| Delivery strategy | chained PRs |
| Chain strategy | `stacked-to-main` |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: `stacked-to-main`
Size exception: not selected
400-line budget risk: High

This forecast includes every existing Result V1 planning artifact, all focused/adversarial/differential/mutation-kill coverage, and the future apply-progress evidence. The owner selected chained PRs with `stacked-to-main`; no size exception is selected. Each PR remains at or below 400 additions + deletions without code-golf or artificial splits.

### Selected stacked-to-main chain

- Each PR targets `main` only after its predecessor lands; PR1 is current.
- PR1 exploration+proposal (379); PR2 three Result specs (321); PR3 design+evidence-v1 spec (366); PR4 tasks (314).
- PR5 accounting foundation (235–290); PR6 Result model (255–305); PR7 wire+identity (285–330); PR8 exact observations (280–335); PR9 antecedents+snapshot (275–330).
- PR10 constructor (285–335); PR11 ownership (135–160); PR12 legacy adapter+differential evidence (255–320); PR13 decoder shape (350–400 target); PR14 authority hardening (75–95).
- PR15 semantic decoder (280–335); PR16 canonical rejection (145–175); PR17 report binding (235–295); PR18 apply evidence/final census (25–40).
- PR13/PR15 must stop rather than hide overage if actual scope exceeds 400; no code-golf or artificial split.

### Production forecast (changed lines, additions + deletions)

| Planned path | Estimate | Purpose |
|---|---:|---|
| `accounting_domain.go` | 140–160 | private shared first-match/default classification and adapters |
| `accounting_domain_totals.go` | 5–5 | checked accumulation helper separated from classification |
| `accounting.go` | 110–145 | behavior-preserving legacy adapter/refactor |
| `result_v1.go` | 135–155 | exported contract types, immutable document, cloning accessors |
| `result_v1_wire.go` | 140–160 | ordered root/entry/total wire structs, Base64, LF, SHA-256 |
| `result_v1_wire_observations.go` | 5–5 | concrete observation-wire variants |
| `result_v1_observations.go` | 135–160 | threshold/ratio variants, checked totals, exact rational comparator |
| `result_v1_antecedents.go` | 130–155 | exact Policy→Inventory revalidation and snapshot boundary validation |
| `result_v1_build.go` | 140–160 | Result construction, sorting, totals/observation assembly |
| `result_v1_build_ownership.go` | 5–5 | defensive ownership ingress/egress helper |
| `result_v1_decode.go` | 135–155 | decoder entry point, exact rebuild/equality, zero-on-error boundary |
| `result_v1_decode_shape.go` | 140–160 | duplicate-aware closed object/array/variant parser |
| `result_v1_decode_shape_keys.go` | 5–5 | normalized structural-key/authority-key helper |
| `result_v1_decode_values.go` | 135–160 | primitive, Base64, path, revision, ordering, and authority-key validation |
| `contract.go` | 90–120 | additive Result-bound Report V1 builder and validator only |
| **Production subtotal** | **1,450–1,710** | **sum of the 15 production paths** |

### Test forecast (changed lines, additions + deletions)

| Planned path | Estimate | Purpose |
|---|---:|---|
| `accounting_test.go` | 90–125 | legacy `Account` regression characterization after refactor |
| `result_v1_test.go` | 120–150 | shared exact antecedent/snapshot fixtures and zero-value assertions |
| `result_v1_wire_test.go` | 140–160 | hostile-byte canonical wire/digest and empty-array vectors |
| `result_v1_wire_identity_test.go` | 0–5 | conditional identity-vector split if the wire test reaches its 160-line target |
| `result_v1_antecedents_test.go` | 145–160 | exact antecedent and inventory-only vectors |
| `result_v1_snapshot_test.go` | 0–15 | conditional revision/path boundary-vector split |
| `result_v1_construct_test.go` | 145–160 | exclusive classifications and policy-order totals |
| `result_v1_order_test.go` | 0–15 | conditional deterministic-order vector split |
| `result_v1_ownership_test.go` | 130–155 | source/egress defensive-copy and failed-result vectors |
| `result_v1_observations_test.go` | 145–160 | threshold, unavailable, exact-rational, and overflow vectors |
| `result_v1_observations_precision_test.go` | 0–15 | conditional precision/equality-boundary vector split |
| `result_v1_compatibility_test.go` | 145–160 | `Account` shared-domain differential vectors authored in unit 5 |
| `result_v1_compatibility_edge_test.go` | 0–15 | conditional precision-edge differential split authored in unit 3 |
| `result_v1_decode_shape_test.go` | 145–160 | closed shapes, duplicates (including escaped), and types |
| `result_v1_decode_authority_test.go` | 0–15 | conditional authority-key/discriminator vector split |
| `result_v1_decode_semantics_test.go` | 145–160 | links, recomputation, raw-path/category/total/observation/order mutations |
| `result_v1_decode_recompute_test.go` | 0–15 | conditional recomputation-mutation vector split |
| `result_v1_decode_canonical_test.go` | 145–160 | whitespace, escaping, LF, Base64, numeric, and trailing-value rejection |
| `result_v1_decode_alternates_test.go` | 0–15 | conditional alternate-wire vector split |
| `result_v1_report_binding_test.go` | 145–160 | derived report provenance and mismatch rejection |
| `result_v1_report_legacy_test.go` | 0–15 | conditional legacy Report V1 preservation-vector split |
| **Test subtotal** | **1,640–1,995** | **sum of the 21 test paths** |

### Existing planning-artifact forecast (changed lines, additions + deletions)

| Existing path | Measured lines | Inclusion |
|---|---:|---|
| `openspec/changes/result-v1/exploration.md` | 234 | existing planning artifact |
| `openspec/changes/result-v1/proposal.md` | 145 | existing planning artifact |
| `openspec/changes/result-v1/specs/result-v1-contract/spec.md` | 138 | existing planning artifact |
| `openspec/changes/result-v1/specs/result-v1-accounting-evidence/spec.md` | 116 | existing planning artifact |
| `openspec/changes/result-v1/specs/result-v1-provenance/spec.md` | 67 | existing planning artifact |
| `openspec/changes/result-v1/specs/evidence-v1-provenance/spec.md` | 38 | existing planning artifact |
| `openspec/changes/result-v1/design.md` | 328 | existing planning artifact |
| `openspec/changes/result-v1/tasks.md` | 314 | existing planning artifact |
| **Existing-artifact subtotal** | **1,380** | **234 + 145 + 138 + 116 + 67 + 38 + 328 + 314** |

### Future evidence-artifact forecast (changed lines, additions + deletions)

| Future path | Estimate | Purpose |
|---|---:|---|
| `openspec/changes/result-v1/apply-progress.md` | 25–40 | required first applicable layout-threshold crossing, focused RED/GREEN/triangulation, mutation-kill, verification, and revert evidence; do not create in this phase |
| **Future-evidence subtotal** | **25–40** | **apply-only artifact** |

**Arithmetic:** production **1,450–1,710** + tests **1,640–1,995** = Go implementation **3,090–3,705**; + existing artifacts **1,380** = **4,470–5,085**; + future evidence **25–40** = **4,495–5,125** changed lines.

No rename is assumed or omitted: listed `result_v1*` paths are additions, `accounting.go` and `contract.go` are edits, and the eight planning artifacts are included as existing delivery content. The explicitly named 0-minimum test paths are conditional responsibility splits only when their primary co-located test would cross 160 changed lines; their maximums are included above. No production/test file is planned above the `openspec/config.yaml` 160 changed-lines-per-Go-file design target, and a split is not permission to evade the package-boundary review.

**Confidence:** medium: approved behavior and test matrix drive the range, while implementation discovery determines whether a conditional test split is needed. The forecast is well above 400 lines; the owner-selected chain constrains delivery before apply.

> Layout control: `openspec/config.yaml` records a 160 changed-lines-per-Go-file threshold and an already-exceeded public-root-file baseline. Keep files below that threshold only where it follows a real responsibility boundary; do not compress code or split mechanically to game the review budget. Before the first applicable crossing, obtain the required package-boundary review and record it during apply in `apply-progress.md`; do not create that later-phase artifact now.

## Shared execution rules

- Work in the public root package only. Do not change `cmd/`, `internal/`, Policy V1, Inventory V1, Report V1 wire decoding/projections, publication, CLI behavior, or the historical bootstrap artifacts.
- For every Go-producing work unit: add the named table-driven RED test first, run its focused command and observe failure, implement only the named GREEN behavior, add the named TRIANGULATE/adversarial vectors, then refactor and rerun the focused command. Keep tests co-located and use `t.Run` case names by scenario.
- A failed constructor or decoder must return the exact zero `ResultDocumentV1`; each new fatal-path test must assert nil canonical bytes/views and empty digest/links/revisions.
- Use only raw path bytes (`[]byte(change.Path)`), `bytes.Compare`, `base64.StdEncoding`, checked `uint64` addition, and exact decimal/rational arithmetic. Do not normalize paths, use maps for canonical JSON, or introduce authority/delivery vocabulary.
- No concurrency-bearing behavior is designed. Race testing is **N/A** unless implementation introduces goroutines, channels, mutexes, atomics, or shared mutable execution; if it does, run `go test -race ./...` before closing that work unit and report the result.

## Work units

### Implementation/progress checklist (authoritative)

- [x] Unit 0 — Apply gate, baseline, and layout review (completed by merged PR57 at `31bcc06d9147adae1c830040d494ab5bdba5a659`; its baseline/package-boundary evidence was independently verified as evidence `14046`).
- [x] Unit 1 — Shared accounting classification and totals foundation (completed by merged PR57; see independent evidence `14046`).
- [x] Unit 2 — Canonical Result V1 types, wire encoding, and identity primitive (PR6 model slice and PR7 wire/identity remainder complete).
  - [x] PR6 — Public Result V1 model types, private document state, and defensive inspection accessors.
  - [x] PR7 — Canonical wire structs, encoding, Base64/LF, and system-derived identity.
- [x] Unit 3 — Exact-rational Result observation derivation.
- [ ] Unit 4 — Exact antecedent and immutable snapshot validation boundary.
- [ ] Unit 5 — Result constructor, deterministic ordering, and defensive ownership.
- [ ] Unit 6 — Strict decoder: duplicate-aware closed structural shape.
- [ ] Unit 7 — Strict decoder: value validation, recomputation, and noncanonical rejection.
- [ ] Unit 8 — Antecedent-aware Report V1 provenance binding without report-surface expansion.
- [ ] Unit 9 — Mutation-kill evidence and final verification census.

### 0. Apply gate, baseline, and layout review (no implementation)

- **Dependencies:** none; this gate must complete before any Go-producing work unit.

- **Start / discovery targets:** confirm the current contracts and regression seams in `accounting.go:Account`, `policy_contract.go:PolicyDocumentV1`, `inventory_v1_build.go:validatedPolicyDigest`, `inventory_v1_decode.go:DecodeInventoryV1`, `snapshot.go:CommittedSnapshot`, `contract.go:NewReportV1`, and `decode.go:DecodeCanonical`. Confirm the governed commands in `openspec/config.yaml`.
- **Delivery decision:** use the owner-selected `stacked-to-main` chain with no size exception; each PR targets `main` only after its predecessor lands. Unit 0 was completed by PR57 at `31bcc06d9147adae1c830040d494ab5bdba5a659`; independent evidence `14046` verified its baseline and package-boundary review. PR6 records only its bounded root-package review in `apply-progress.md`.
- **Baseline verification:**
  ```bash
  go test . -run '^(TestAccountUsesFirstMatchingRawByteGlobAndDefault|TestAccountRejectsInvalidPolicyBeforeAccounting|TestAccountReturnsEmptyResultOnOverflow|TestNewReportV1CanonicalBytesAndDigest|TestDecodeCanonicalRejectsClosedAndAuthorityBearingDocuments|TestProjectEvidenceCanonicalJSON|TestProjectEvidenceHumanText)$'
  go test ./...
  test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"
  git diff --check
  git diff --stat
  ```
- **Finish / rollback:** capture the baseline results and the package-boundary review decision in apply-time evidence only; no production behavior or later-phase artifact is changed by this unit, so rollback is N/A.

### 1. Shared accounting classification and totals foundation

- **Dependencies:** unit 0 gate and baseline evidence.

- **Files / symbols:** add `accounting_domain.go` with private `accountingDomain` and legacy/Result adapters plus `accounting_domain_totals.go` with checked classify-and-accumulate helpers; narrowly refactor `accounting.go:Account`, `validateAccountingPolicy`, `matchingCategory`, `addLines`, and `lineTotal`; extend `accounting_test.go`.
- **Requirement coverage:** `result-v1-accounting-evidence` **Exclusive classification and recomputable policy-category totals** / *Overlap, default, and non-countable evidence* and *Realizable arithmetic overflow is fatal rather than invented evidence*; **Legacy accounting compatibility** / *Legacy caller remains unaffected*. It establishes the common mechanics consumed by the later Result-only scenarios.
- **RED:** add `TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable` and `TestAccountingDomainRejectsCheckedAccumulationOverflow`; first run:
  ```bash
  go test . -run '^(TestAccountUsesFirstMatchingRawByteGlobAndDefault|TestAccountReturnsEmptyResultOnOverflow|TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable|TestAccountingDomainRejectsCheckedAccumulationOverflow)$'
  ```
- **GREEN:** route legacy `Account` through one private raw-byte domain that retains category order, first-match glob precedence, default selection, non-countable handling, checked totals, legacy errors, and legacy observation conversion/`float64` behavior.
- **TRIANGULATE:** table-test overlapping `0xff` glob paths, unmatched defaults, non-countable entries with nonzero input lines, all-zero categories, and `MaxUint64 + 1`; assert exact legacy totals, observation order, and empty result on failure.
- **REFACTOR / verify:** preserve `AccountingPolicy`, `AccountingResult`, and all exported `Account` behavior; rerun the focused command and `go test ./...`. **Rollback:** revert only `accounting_domain.go`, the `accounting.go` delegation, and matching regression tests; legacy accounting remains self-contained.

### 2. Canonical Result V1 types, wire encoding, and identity primitive

- **Dependencies:** unit 0 gate and baseline evidence; it is otherwise additive and independent of the accounting refactor.

- **Files / symbols:** add `result_v1.go` (`AccountingResultContractV1`, entry/measurement/total/observation views, private `ResultDocumentV1` state and cloning accessors), `result_v1_wire.go` (ordered root/entry/total wire structs and canonical encoder), and `result_v1_wire_observations.go` (concrete observation variants); add `result_v1_test.go`, `result_v1_wire_test.go`, and, only if needed to retain the 160-line target, `result_v1_wire_identity_test.go`.
- **Requirement coverage:** `result-v1-contract` **Closed Result V1 wire contract and content identity** / *Canonical hostile-path fixture*, *Empty committed change set*, and *Alternate wire representation*; `result-v1-provenance` **Result identity is system-derived and complete** / *One raw path byte changes the identity*.
- **RED:** add `TestResultV1CanonicalWireHostilePathAndIdentity` and `TestResultV1CanonicalWireUsesEmptyArrays`; run:
  ```bash
  go test . -run '^(TestResultV1CanonicalWireHostilePathAndIdentity|TestResultV1CanonicalWireUsesEmptyArrays)$'
  ```
- **GREEN:** encode the exact seven root fields in contractual order, concrete per-variant observation objects only, padded standard Base64 paths, non-nil `[]` arrays, one terminal LF, and SHA-256 over those final bytes. Keep digest out of the body and return defensive copies.
- **TRIANGULATE:** assert the specified `/wpmaWxl` hostile-path fixture; differentiate `0xff` from `0xfe` and LF-adjacent byte paths; assert zero observations and all-zero category totals serialize without `null`.
- **REFACTOR / verify:** no map-based encoding, path text conversion, or caller identity input. Rerun the focused command and `go test ./...`. **Rollback:** remove only the Result type/wire primitive and its tests; no legacy or report surface changes.

### 3. Exact-rational Result observation derivation

- **Dependencies:** units 1 (shared domain) and 2 (Result types/wire primitive).

- **Files / symbols:** add `result_v1_observations.go` with checked required line-total helpers and `exactRatioExceeded` using `math/big.Int` plus the exact Policy V1 decimal; add `result_v1_observations_test.go` and, only if needed to retain the 160-line target, `result_v1_observations_precision_test.go` or `result_v1_compatibility_edge_test.go` for the focused legacy precision differential.
- **Requirement coverage:** `result-v1-accounting-evidence` **Closed factual threshold, ratio, and unavailable-observation model** / *Available threshold and exact ratio*, *Exact decimal comparison intentionally differs at a legacy float64 precision edge*, *Non-countable and zero-denominator observations remain embedded*, and *Exceeded evidence is not authority*.
- **RED:** add `TestResultV1ObservationsUseClosedOrderAndUnavailableVariants`, `TestResultV1ExactRatioPrecisionEdge`, the legacy differential `TestResultV1PrecisionEdgeIsTheOnlyAvailableRatioComparatorDifference`, `TestResultV1StrictGreaterThanExcludesEqualityBoundary`, and `TestResultV1RejectsRequiredLineTotalOverflow`; run before exact-comparator GREEN behavior:
  ```bash
  go test . -run '^(TestResultV1ObservationsUseClosedOrderAndUnavailableVariants|TestResultV1ExactRatioPrecisionEdge|TestResultV1PrecisionEdgeIsTheOnlyAvailableRatioComparatorDifference|TestResultV1StrictGreaterThanExcludesEqualityBoundary|TestResultV1RejectsRequiredLineTotalOverflow)$'
  ```
- **GREEN:** emit nonzero line thresholds in Policy V1 category order followed by declared ratios; distinguish unavailable from available zero; use strict `>`; preserve the exact decimal string; fail on every stored/required sum overflow. Do not alter `Account`'s `float64` comparison.
- **TRIANGULATE:** cover `2+1 > 2`, equality boundaries for both a line threshold and exact ratio (`actual == limit` and `numerator / denominator == limit` must not exceed), zero denominator, non-countable numerator/reference, exact `9007199254740993 / 18014398509481984 > 0.5`, ordinary shared-domain ratios, and additions-plus-deletions overflow. Compare that >2^53 fixture with legacy `Account` and prove its historical `float64` result is the intentional lone available-ratio comparator divergence.
- **REFACTOR / verify:** retain only factual `exceeded` evidence and the two closed kinds. Rerun the focused command, the legacy accounting focus from unit 1, and `go test ./...`. **Rollback:** remove the Result-only observation helper/tests without changing the legacy observation loop.

### 4. Exact antecedent and immutable snapshot validation boundary

- **Dependencies:** unit 2 for shared Result fixtures/types and unit 0 baseline; it consumes only shipped Policy/Inventory/snapshot contracts.

- **Files / symbols:** add `result_v1_antecedents.go` with `validateResultAntecedents` and `validateResultSnapshot`, reusing strict Policy/Inventory decoders and `validInventoryPath`-equivalent raw-path rules; extend `result_v1_test.go`, add `result_v1_antecedents_test.go`, and add `result_v1_snapshot_test.go` only if its focused boundary cases are needed to retain the 160-line target.
- **Requirement coverage:** `result-v1-provenance` **Exact validated antecedent chain** / *Byte-identical antecedents reproduce provenance*, *Inventory policy link is wrong but all digests are well formed*, *Same accounting with a different antecedent is not interchangeable*; `result-v1-accounting-evidence` **Valid immutable committed-snapshot boundary** / *Valid native SHA-256 snapshot with hostile path bytes*, *Invalid revision or path input*, and *Rename classification uses the committed primary path*.
- **RED:** add `TestValidateResultAntecedentsRequiresExactPolicyInventoryChain` and `TestValidateResultSnapshotAcceptsNativeIDsAndRejectsInvalidPaths`; run:
  ```bash
  go test . -run '^(TestValidateResultAntecedentsRequiresExactPolicyInventoryChain|TestValidateResultSnapshotAcceptsNativeIDsAndRejectsInvalidPaths)$'
  ```
- **GREEN:** re-decode antecedent canonical bytes, recompute and compare digests, require the inventory policy link, copy immutable views, validate same-width 40/64 lowercase revisions, equal revisions only with no entries, byte-identical unique primary paths, and use only `Path`/`Lines`.
- **TRIANGULATE:** test zero/stale/substituted documents, digest-shaped assertions, wrong linked inventory, SHA-1 and SHA-256 pairs, cross-width/uppercase/symbolic IDs, equal revisions with entries, NUL/absolute/empty/dot/dotdot paths, duplicate raw paths, and rename `PreviousPath` exclusion.
- **REFACTOR / verify:** inventory remains a provenance link only—never an accounting source. Rerun the focused command and `go test ./...`. **Rollback:** remove only the validation seam/tests; no snapshot acquisition, Git, or filesystem behavior is added.

### 5. Result constructor, deterministic ordering, and defensive ownership

- **Dependencies:** units 1–4; construction joins the shared domain, Result types/wire, observations, and validated inputs.

- **Files / symbols:** add `result_v1_build.go`, `result_v1_build_ownership.go`, and finish `result_v1.go:NewAccountingResultV1`, `CanonicalBytes`, `Digest`, `AccountingPolicyDigest`, `InventoryDigest`, `Revisions`, `Entries`, `Totals`, and `Observations`; add `result_v1_construct_test.go`, `result_v1_ownership_test.go`, and the required early compatibility file `result_v1_compatibility_test.go`; add `result_v1_order_test.go` only where its focused scope is needed to retain the 160-line target.
- **Requirement coverage:** `result-v1-contract` **Neutral Result V1 construction and inspection** / *Construct a neutral result from exact typed evidence*, *Caller offers only digest assertions*, *Failed construction has no partial result*; **Defensive ownership and deterministic construction** / *Caller mutates a source or returned buffer* and *Equivalent snapshot orders*; `result-v1-accounting-evidence` **Inventory is provenance-only and excluded from accounting** / *Inventory identity changes without changing accounting* and *Nonempty inventory with empty committed snapshot*; **Exclusive classification and recomputable policy-category totals** / *Overlap, default, and non-countable evidence* and *Constructor input order does not establish result order*; **Legacy accounting compatibility** / *Result and legacy accounting agree on shared-domain vectors*. These public differentials are deliberately authored here, before constructor GREEN behavior, rather than deferred to unit 9.
- **RED:** add `TestNewAccountingResultV1ConstructsExclusiveCanonicalEvidence`, `TestNewAccountingResultV1OwnsIngressAndEgress`, `TestNewAccountingResultV1InventoryOnlyChangesIdentity`, `TestNewAccountingResultV1ReturnsZeroOnFatalFailure` and `TestResultV1AndAccountAgreeOnSharedDomainVectors`; run before constructor GREEN behavior:
  ```bash
  go test . -run '^(TestNewAccountingResultV1ConstructsExclusiveCanonicalEvidence|TestNewAccountingResultV1OwnsIngressAndEgress|TestNewAccountingResultV1InventoryOnlyChangesIdentity|TestNewAccountingResultV1ReturnsZeroOnFatalFailure|TestResultV1AndAccountAgreeOnSharedDomainVectors|TestAccountUsesFirstMatchingRawByteGlobAndDefault|TestAccountReturnsEmptyResultOnOverflow)$'
  ```
- **GREEN:** construct only from exact typed Policy/Inventory/Snapshot inputs; classify every committed primary path once through the shared domain, sort entries with `bytes.Compare`, retain Policy V1 total order, derive observations, canonicalize/digest once, and own all ingress/egress slices. Do not expose legacy `AccountingPolicy` as a Result input.
- **TRIANGULATE:** test reverse input order (`a`, `a/\xff`, `a0`, `\xff`), overlapping/default categories, non-countable entries, empty committed snapshots with nonempty inventory, changing inventory identity without changing accounting fields, and mutation of input entries/returned bytes/nested paths/views. In the same early differential suite, compare legacy `Account` and Result classifications/totals/availability/ordinary-ratio outcomes on shared-domain vectors; the precision-edge comparator differential is already RED/GREEN/TRIANGULATE-covered in unit 3.
- **REFACTOR / verify:** assert exact zero Result state on every antecedent, snapshot, or overflow failure. Rerun the focused command, units 1–4 focused commands, and `go test ./...`. **Rollback:** remove Result construction/ownership files and tests while leaving Policy, Inventory, snapshot, and legacy Account APIs usable.

### 6. Strict decoder: duplicate-aware closed structural shape

- **Dependencies:** units 2, 4, and 5; shape parsing needs the Result model, validated antecedents, and builder-backed recomputation boundary.

- **Files / symbols:** add `result_v1_decode.go`, `result_v1_decode_shape.go`, and `result_v1_decode_shape_keys.go` with Result-specific wrappers/key normalization around `decodePolicyObject`; add `result_v1_decode_shape_test.go` and `result_v1_decode_authority_test.go` only if needed to retain the 160-line target.
- **Requirement coverage:** `result-v1-contract` **Strict antecedent-aware canonical decoding** / *Unknown, duplicate, and authority-bearing nested fields* and *Policy-derived authority-like string remains evidence*; `result-v1-provenance` **Result V1 remains neutral evidence-only provenance** / *Authority-bearing candidate provenance* and *Policy-derived authority-like value remains evidence*.
- **RED:** add `TestDecodeAccountingResultV1RejectsClosedShapesDuplicatesAndAuthorityKeys`; run:
  ```bash
  go test . -run '^TestDecodeAccountingResultV1RejectsClosedShapesDuplicatesAndAuthorityKeys$'
  ```
- **GREEN:** require every field at root/revisions/entry/measurement/total/observation depth; reject unknown and duplicate decoded keys including escaped duplicates; recognize normalized authority-bearing keys/discriminators only as structural metadata; accept permitted exact Policy-derived values such as `merge` or `release`.
- **TRIANGULATE:** table-test missing fields, `null`, arrays/booleans/strings where objects/arrays are required, unknown fields, `approval`, `merge-gate`, `delivery_authority`, duplicate and escaped-duplicate `schema`/`path_b64`, unsupported schema/kind, and policy-derived `release` category values.
- **REFACTOR / verify:** preserve correct field-path `ContractError` values and never use `json.Unmarshal` alone where it could collapse duplicates. Rerun the focused command and `go test ./...`. **Rollback:** remove only Result decoder shape code/tests; construction remains available without a decoder.

### 7. Strict decoder: value validation, recomputation, and noncanonical rejection

- **Dependencies:** unit 6 plus units 3–5; semantic decode completes the shape parser by rebuilding through the exact observation and constructor paths.

- **Files / symbols:** add `result_v1_decode_values.go`; finish `DecodeAccountingResultV1` in `result_v1_decode.go`; add `result_v1_decode_semantics_test.go`, `result_v1_decode_canonical_test.go`, and, only if needed to retain the 160-line target, `result_v1_decode_recompute_test.go` and `result_v1_decode_alternates_test.go`.
- **Requirement coverage:** `result-v1-contract` **Strict antecedent-aware canonical decoding** / *Wrong exact antecedent despite valid syntax* and *Recomputed evidence differs*; **Closed Result V1 wire contract and content identity** / *Alternate wire representation*; `result-v1-provenance` **Result identity is system-derived and complete** / *Reordered candidate is not repaired*.
- **RED:** add `TestDecodeAccountingResultV1RoundTripsOnlyExactCanonicalEvidence`, `TestDecodeAccountingResultV1RejectsValueAndSemanticMutations`, and `TestDecodeAccountingResultV1RejectsNoncanonicalRepresentations`; run:
  ```bash
  go test . -run '^(TestDecodeAccountingResultV1RoundTripsOnlyExactCanonicalEvidence|TestDecodeAccountingResultV1RejectsValueAndSemanticMutations|TestDecodeAccountingResultV1RejectsNoncanonicalRepresentations)$'
  ```
- **GREEN:** reject non-UTF-8/trailing JSON, malformed lowercase digests/revisions, noncanonical unsigned integers/booleans/Base64/paths, duplicate paths, wrong variants and orders; rebuild a copied minimal snapshot from candidate entries through the constructor; accept only byte-for-byte canonical equality and return zero Result on every failure.
- **TRIANGULATE:** mutate one link, revision, raw path byte, classification, countable/non-countable marker, total, availability/actual/numerator/denominator/exceeded flag, entry/total/observation order, whitespace/string escaping/final LF, URL or unpadded Base64, numeric strings/fractions/exponents, and a valid-looking different antecedent. Each must fail without repair.
- **REFACTOR / verify:** prove valid canonical decode returns owned views; rerun the focused command and `go test ./...`. **Rollback:** remove Result decoder/value tests only; no changes to `DecodeCanonical` or Inventory decoder behavior.

### 8. Antecedent-aware Report V1 provenance binding without report-surface expansion

- **Dependencies:** units 5 and 7; a binding may derive Report V1 provenance only from a constructible and strictly decodable Result V1.

- **Files / symbols:** narrowly extend `contract.go` with `ResultV1ReportBinding`, `NewReportV1FromResult`, and `ValidateReportV1ResultProvenance`; add `result_v1_report_binding_test.go` and, only if needed to retain the 160-line target, `result_v1_report_legacy_test.go`; retain `contract_test.go`, `decode_test.go`, `report_test.go`, and `cmd/git-change-evidence/main_test.go` unchanged.
- **Requirement coverage:** `evidence-v1-provenance` **Provenance and reproducibility** / *Bound report derives all provenance*, *Well-formed but mismatched accounting digest*, *Standalone report bytes are not provenance proof*, and *No projection or authority expansion*.
- **RED:** add `TestResultV1ReportBindingDerivesAllProvenance` and the Report compatibility regression `TestValidateReportV1ResultProvenanceRejectsEveryMismatchAndPreservesLegacySurfaces` here, before binding GREEN behavior (not in unit 9); run:
  ```bash
  go test . -run '^(TestResultV1ReportBindingDerivesAllProvenance|TestValidateReportV1ResultProvenanceRejectsEveryMismatchAndPreservesLegacySurfaces|TestNewReportV1CanonicalBytesAndDigest|TestDecodeCanonicalRejectsClosedAndAuthorityBearingDocuments|TestProjectEvidenceCanonicalJSON|TestProjectEvidenceHumanText)$'
  ```
- **GREEN:** strict-validate the exact Policy→Inventory→Result chain, derive all five report provenance values from the Result, then call unchanged `NewReportV1`; validate a supplied Evidence canonically and compare every derived provenance value. Do not change `NewReportV1`, `DecodeCanonical`, report wire ordering, projections, CLI, or publication.
- **TRIANGULATE:** reject one well-formed-but-different policy, inventory, accounting digest, base, or head; prove standalone `DecodeCanonical` remains syntactic rather than bound proof; prove exceeded evidence does not alter reports/projections or create authority vocabulary.
- **REFACTOR / verify:** rerun the focused command, `go test ./...`, and the unchanged CLI focus:
  ```bash
  go test ./cmd/git-change-evidence -run '^(TestRunProjectCanonicalJSON|TestRunProjectHumanTextAndRepeatedInput)$'
  ```
  **Rollback:** remove only the additive binding type/APIs/tests; legacy report construction, decode, CLI, projection, and publication remain unchanged.

### 9. Mutation-kill evidence and final verification census

- **Dependencies:** units 1–8. This is verification/census only: it creates no compatibility test and performs no RED/GREEN behavior cycle.

- **Files / symbols:** no new source or compatibility test is added here. Census the tests authored in units 1, 3, 5, and 8: `accounting_test.go`, `result_v1_observations_test.go`, `result_v1_compatibility_test.go`, optional `result_v1_compatibility_edge_test.go`, and `result_v1_report_binding_test.go`/optional `result_v1_report_legacy_test.go`.
- **Requirement coverage:** verification census only. The legacy `Account` shared-domain scenario is created in unit 5, while its precision-edge differential and exact-ratio/equality boundaries are created in unit 3; legacy caller preservation is created in unit 1; the legacy Report V1 preservation scenario is created in unit 8. This unit adds no requirement behavior.
- **Verification / census:** rerun the already-authored differential and legacy regressions; they must remain green without behavior changes:
  ```bash
  go test . -run '^(TestAccountUsesFirstMatchingRawByteGlobAndDefault|TestAccountReturnsEmptyResultOnOverflow|TestResultV1ExactRatioPrecisionEdge|TestResultV1StrictGreaterThanExcludesEqualityBoundary|TestResultV1AndAccountAgreeOnSharedDomainVectors|TestResultV1PrecisionEdgeIsTheOnlyAvailableRatioComparatorDifference|TestResultV1ReportBindingDerivesAllProvenance|TestValidateReportV1ResultProvenanceRejectsEveryMismatchAndPreservesLegacySurfaces)$'
  ```
- **Mutation-kill evidence:** temporarily apply each one-line mutant locally, run the named focused test expecting failure, then revert the mutant before proceeding. Record the mutant, command, observed failing test name, and clean-revert confirmation in apply-time `openspec/changes/result-v1/apply-progress.md`; do not create that later-phase artifact in this tasks phase.

  | Load-bearing predicate | Temporary mutant | Required killing test |
  |---|---|---|
  | raw ordering | remove/invert `bytes.Compare` entry sort | `TestNewAccountingResultV1ConstructsExclusiveCanonicalEvidence`, `TestDecodeAccountingResultV1RejectsValueAndSemanticMutations` |
  | provenance | bypass inventory-policy/result-link equality | `TestValidateResultAntecedentsRequiresExactPolicyInventoryChain`, `TestDecodeAccountingResultV1RejectsValueAndSemanticMutations` |
  | Base64/raw bytes | use URL Base64 or normalize `[]byte(change.Path)` | `TestResultV1CanonicalWireHostilePathAndIdentity`, `TestDecodeAccountingResultV1RejectsNoncanonicalRepresentations` |
  | exclusivity | select last match or count an inventory entry | `TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable`, `TestNewAccountingResultV1InventoryOnlyChangesIdentity` |
  | overflow | replace checked addition with wrapping `+` | `TestAccountingDomainRejectsCheckedAccumulationOverflow`, `TestResultV1RejectsRequiredLineTotalOverflow` |
  | exact rational | replace `math/big.Int` cross multiplication in `exactRatioExceeded` with `float64` comparison | `TestResultV1ExactRatioPrecisionEdge` (the `9007199254740993 / 18014398509481984 > 0.5` fixture) |
  | strict comparison | replace strict `>` with `>=` in line-threshold or `exactRatioExceeded` comparison | `TestResultV1StrictGreaterThanExcludesEqualityBoundary` (equal threshold and equal exact-ratio fixtures) |
  | digest completeness | omit final LF or hash a partial wire body | `TestResultV1CanonicalWireHostilePathAndIdentity`, `TestDecodeAccountingResultV1RoundTripsOnlyExactCanonicalEvidence`, `TestResultV1ReportBindingDerivesAllProvenance` |

- **Full verification:** after every implementation-unit refactor has rerun its focused test, run the exact policy and available repository checks:
  ```bash
  go test ./...
  test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"
  go test -cover ./...
  git diff --check
  git diff --stat
  git diff --numstat
  ```
  `go test -cover ./...` is informational because `coverage_threshold: 0`; report its actual output rather than inventing a target. Race gate remains N/A for this pure immutable contract; if concurrency was introduced contrary to this plan, also run `go test -race ./...` and treat failure as blocking verification.
- **Finish / rollback:** confirm no changes under `cmd/`, `internal/`, policy/inventory contracts, report wire/projection, publication, or bootstrap artifacts. Roll back the Result V1 and bound-report units independently: revert their files and the accounting-domain refactor together only if legacy accounting must return to its prior private implementation.

## Requirement and scenario traceability

| Specification requirement | Exact scenarios | Work units |
|---|---|---|
| `result-v1-contract` — Neutral Result V1 construction and inspection | Construct a neutral result from exact typed evidence; Caller offers only digest assertions; Failed construction has no partial result | 4, 5 |
| `result-v1-contract` — Closed Result V1 wire contract and content identity | Canonical hostile-path fixture; Empty committed change set; Alternate wire representation | 2, 5, 7 |
| `result-v1-contract` — Defensive ownership and deterministic construction | Caller mutates a source or returned buffer; Equivalent snapshot orders | 5 |
| `result-v1-contract` — Strict antecedent-aware canonical decoding | Wrong exact antecedent despite valid syntax; Recomputed evidence differs; Unknown, duplicate, and authority-bearing nested fields; Policy-derived authority-like string remains evidence | 6, 7 |
| `result-v1-accounting-evidence` — Valid immutable committed-snapshot boundary | Valid native SHA-256 snapshot with hostile path bytes; Invalid revision or path input; Rename classification uses the committed primary path | 4, 5 |
| `result-v1-accounting-evidence` — Inventory is provenance-only and excluded from accounting | Inventory identity changes without changing accounting; Nonempty inventory with empty committed snapshot | 4, 5 |
| `result-v1-accounting-evidence` — Exclusive classification and recomputable policy-category totals | Overlap, default, and non-countable evidence; Constructor input order does not establish result order; Realizable arithmetic overflow is fatal rather than invented evidence | 1, 5 |
| `result-v1-accounting-evidence` — Closed factual threshold, ratio, and unavailable-observation model | Available threshold and exact ratio; Exact decimal comparison intentionally differs at a legacy float64 precision edge; Non-countable and zero-denominator observations remain embedded; Exceeded evidence is not authority | 3, 5, 7 |
| `result-v1-accounting-evidence` — Legacy accounting compatibility | Legacy caller remains unaffected; Result and legacy accounting agree on shared-domain vectors | 1, 5 |
| `result-v1-provenance` — Exact validated antecedent chain | Byte-identical antecedents reproduce provenance; Inventory policy link is wrong but all digests are well formed; Same accounting with a different antecedent is not interchangeable | 4, 5, 7 |
| `result-v1-provenance` — Result identity is system-derived and complete | One raw path byte changes the identity; Reordered candidate is not repaired | 2, 5, 7 |
| `result-v1-provenance` — Result V1 remains neutral evidence-only provenance | Authority-bearing candidate provenance; Policy-derived authority-like value remains evidence; Out-of-scope integration request | 0, 6, 7, 8 |
| `evidence-v1-provenance` — Provenance and reproducibility | Bound report derives all provenance; Well-formed but mismatched accounting digest; Standalone report bytes are not provenance proof; No projection or authority expansion | 8 |
