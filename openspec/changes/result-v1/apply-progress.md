# Apply Progress: Result V1

## PR6 — Model slice

- **Work unit:** PR6 Result V1 model only, in the selected `stacked-to-main` chain.
- **Completed:** public constants, discriminated entry/observation views, private `ResultDocumentV1` state, and defensive inspection accessors.
- **Excluded:** wire encoding, Base64/LF, hashing, constructors, decoding, observation derivation, antecedent/snapshot validation, report binding, CLI, shell, policy, and configuration work remain deferred.
- **Prior evidence:** merged PR57 (`31bcc06d9147adae1c830040d494ab5bdba5a659`) completed Units 0–1; its package-boundary/baseline evidence was independently verified as `14046` (parent-provided evidence).

## Completed tasks

- [x] Reconciled Units 0–1 in `tasks.md` to the merged PR57 and evidence `14046`.
- [x] Added the PR6 model subtask while keeping Unit 2 unchecked and explicitly reserving its wire/identity remainder for PR7.
- [x] Added Result V1 contract constants, view types, private document state, and defensive accessors.
- [x] Added focused model tests for zero values, populated values, copied bytes/slices/nested paths, repeated accessor reads, empty populated views, and discriminators.

## TDD Cycle Evidence

Toolchain for every Go command: `GOTOOLCHAIN=go1.25.10` with that toolchain's `GOROOT/bin` first on `PATH`.

| Stage | Command | Actual result |
|---|---|---|
| RED | `go test . -run '^TestResultV1Model'` | Failed as expected (exit 1): `ResultDocumentV1`, Result V1 views, and discriminators were undefined. |
| GREEN (first compile) | `go test . -run '^TestResultV1Model'` | Failed (exit 1): two clone helpers shadowed Go's `copy` builtin; corrected before accepting GREEN. |
| GREEN | `go test . -run '^TestResultV1Model'` | Passed: `ok github.com/kozz36/git-change-evidence 0.018s`. |
| TRIANGULATE | `go test . -run '^TestResultV1Model'` | Passed after adding the populated-document empty-view cases: `ok ... 0.001s`. |
| REFACTOR | `go test . -run '^TestResultV1Model'` | Passed (`cached`) after review; no behavior-preserving refactor was needed because distinct clone helpers preserve readable typed ownership boundaries. |

## Verification

- `go test ./...` passed for the root package and `cmd/git-change-evidence`, `internal/git`, `internal/inventory`, and `internal/publication`.
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` passed.
- `git diff --check` passed.
- `go test -race ./...` is N/A: this slice adds no concurrency-bearing behavior.

## Layout and workload boundary

The root had 32 Go files before this slice and 34 after it, already above the configured public-root threshold of 12. The bounded package-boundary review approves this additive neutral public API at the root under `AGENTS.md` and the approved design; its co-located test stays at the root. No `cmd/`, `internal/`, new shell/domain, or layout restructuring was added. Two root Go files changed, both under the 160-line threshold: `result_v1.go` (111) and `result_v1_test.go` (150).

This is the PR6 model boundary only. Its Go additions are 261 lines; documentation reconciliation/progress remains within the 400-line PR limit. No size exception is used.

## Files changed

- `result_v1.go`
- `result_v1_test.go`
- `openspec/changes/result-v1/tasks.md`
- `openspec/changes/result-v1/apply-progress.md`

## Remaining work and risks

- Unit 2 remains partial until PR7 supplies canonical wire structs, encoding, Base64/LF, and identity; all later Result V1 units remain pending.
- The strict-TDD support file `.pi/gentle-ai/support/strict-tdd.md` is absent; the configured RED → GREEN → TRIANGULATE → REFACTOR contract was followed directly.
- No implementation or test was added for deferred APIs, so future slices must not treat this model state as a valid constructible document without their validation/wire boundaries.

## PR7 — Wire and identity remainder

- **Work unit:** PR7 canonical private wire encoder and final-byte identity primitive, in the selected `stacked-to-main` chain; depends on merged PR6 model state.
- **Completed:** all seven ordered root fields, concrete countable/non-countable measurements, four observation wire variants, padded standard Base64 paths, compact JSON plus one LF, SHA-256 identity over final bytes, and privately owned Result model assembly.
- **Excluded:** public Result construction/decode, sorting/classification, observation derivation, antecedent/snapshot validation, Report binding, CLI/shell/configuration, and every later Result V1 unit.

### Completed tasks

- [x] Completed Unit 2 and its PR7 wire/identity subtask in `tasks.md`.
- [x] Added the private ordered wire structs and encoder without a map or a serialized self-digest.
- [x] Added dedicated hostile-byte, empty-array, closed-variant, Base64-padding, and identity tests; PR6's fake `modelResultDocument` was not used as provenance evidence.

### TDD Cycle Evidence

Toolchain for every Go command: `GOTOOLCHAIN=go1.25.10` with that toolchain's `GOROOT/bin` first on `PATH`.

| Stage | Command | Actual result |
|---|---|---|
| RED | `go test . -run '^(TestResultV1CanonicalWireHostilePathAndIdentity|TestResultV1CanonicalWireUsesEmptyArrays)$'` | Failed as expected (exit 1): `resultDocument` was undefined. |
| GREEN | `go test . -run '^(TestResultV1CanonicalWireHostilePathAndIdentity|TestResultV1CanonicalWireUsesEmptyArrays)$'` | Passed: `ok github.com/kozz36/git-change-evidence 0.001s`. |
| TRIANGULATE RED | `go test . -run '^TestResultV1CanonicalWire'` | Failed as expected (exit 1): the new closed-variant fixture received `invalid contract observations: unsupported`. |
| TRIANGULATE GREEN | `go test . -run '^TestResultV1CanonicalWire'` | Passed: `ok github.com/kozz36/git-change-evidence 0.001s`. |
| REFACTOR | `gofmt -w result_v1_wire.go result_v1_wire_observations.go result_v1_wire_test.go result_v1_wire_identity_test.go && go test . -run '^TestResultV1CanonicalWire'` | Passed after format-only refactor (`cached`). |

### Verification

- `go test ./...` passed for root, CLI, Git, inventory, and publication packages.
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` passed.
- `git diff --check` passed.
- `go test -race ./...` is N/A: this pure encoder adds no concurrency-bearing behavior.

### Layout, workload, and PR boundary

PR7 adds four root-package Go files: `result_v1_wire.go` (108), `result_v1_wire_observations.go` (65), `result_v1_wire_test.go` (65), and `result_v1_wire_identity_test.go` (59), for 297 Go additions and no Go deletions. The root now has 38 Go files. The user separately approved PR7's bounded four-file public-root addition; prior Unit 0/PR6 evidence remains historical context. Neutral private wire code and co-located root tests belong with the existing public Result V1 model, with no new public API, `cmd/`, `internal/`, or layout restructuring. All four files are under 160 changed lines and four changed Go files remain below the eight-file package threshold.

The planned five-line observation helper forecast was inaccurate: concrete closed wire types require 65 lines. This is a forecast correction, not a design deviation. The PR7 Go implementation remains within its 285–330 forecast at 297 additions and below the owner-approved 400-total-changed-line PR budget, including documentation. The current PR7 diff is 352 additions and 3 deletions (355 changed lines): 297 Go additions plus 55 documentation additions and 3 documentation deletions. Recorded delivery through PR6 is 1,994 changed lines (1,380 planning + 290 PR5 + 324 PR6); including this uncommitted 355-line PR7 candidate, projected cumulative delivery is 2,349. The historical 4,495–5,125 total forecast is retained; its arithmetic remaining forecast after PR7 is 2,146–2,776 changed lines.

### Files changed

- `result_v1_wire.go`
- `result_v1_wire_observations.go`
- `result_v1_wire_test.go`
- `result_v1_wire_identity_test.go`
- `openspec/changes/result-v1/tasks.md`
- `openspec/changes/result-v1/apply-progress.md`

### Remaining work and risks

- Units 3–9 remain pending; no later-unit code or public API was added.
- The strict-TDD support file `.pi/gentle-ai/support/strict-tdd.md` remains absent; the configured RED → GREEN → TRIANGULATE → REFACTOR contract was followed directly.
- Candidate remains uncommitted for the parent-owned independent check and lifecycle settlement.

## PR8 — Exact-rational observations

- **Work unit:** Unit 3 only, the approved `stacked-to-main` PR8 exact-observation slice. Baseline is `origin/main` at `3f2ef5c0b69b4e034aa590012be3018cba4b20a9` (merged PR61). The historical PR7 entry records its then-uncommitted state; PR7 is now delivered by that merge. Recorded delivery through PR7 is 2,349 changed lines; the current uncommitted PR8 candidate is 290 changed lines, so the projection is 2,349 + 290 = 2,639 changed lines.
- **Completed:** private derivation from a validated `PolicyView` and recomputed `CategoryTotal` values; closed threshold-then-ratio order; exact Policy V1 decimal limits; unavailable variants; checked required line totals; and `math/big.Int` strict rational comparison.
- **Excluded:** antecedent/snapshot validation (Unit 4), construction/ordering/ownership, decoder work, legacy API changes, report binding, and all other deferred Result V1 work remain untouched.

### Completed tasks

- [x] Completed Unit 3 in `tasks.md`; Units 4–9 remain unchanged and pending.
- [x] Added `result_v1_observations.go` with non-countable-before-sum availability, fatal no-partial-observation overflow behavior, and both-side ratio total evaluation.
- [x] Added table-driven closed-order/unavailable, equality, overflow, ordinary-compatibility, and precision-edge tests. Supplemental coverage proves fail-atomicity when an available threshold precedes a late threshold or ratio overflow, and compares equality, non-countable numerator/reference, zero-denominator, and available-zero ratio vectors with legacy `Account`. The optional precision test split is a cohesive exact-ratio/legacy-float boundary and keeps each Go file below 160 lines.

### TDD Cycle Evidence

Toolchain for every Go command: `GOTOOLCHAIN=go1.25.10` with that toolchain's `GOROOT/bin` first on `PATH`.

| Stage | Command | Actual result |
|---|---|---|
| RED | `go test . -run '^(TestResultV1ObservationsUseClosedOrderAndUnavailableVariants|TestResultV1ExactRatioPrecisionEdge|TestResultV1PrecisionEdgeIsTheOnlyAvailableRatioComparatorDifference|TestResultV1StrictGreaterThanExcludesEqualityBoundary|TestResultV1RejectsRequiredLineTotalOverflow)$'` | Failed as expected (exit 1): `resultObservations` was undefined in all five named tests. |
| GREEN | Same focused command | Passed: `ok github.com/kozz36/git-change-evidence 0.001s`. |
| TRIANGULATE | `go test . -run '^TestResultV1NonCountableTotalsRemainUnavailableBeforeAddition$'` | Passed: the additional unavailable vector proves non-countability is checked before an otherwise-overflowing addition/deletion sum. |
| REFACTOR | `gofmt -w result_v1_observations.go result_v1_observations_test.go result_v1_observations_precision_test.go`, then the required focused command | Passed: format-only refactor; no behavior change. |
| Legacy regression | `go test . -run '^(TestAccountUsesFirstMatchingRawByteGlobAndDefault|TestAccountReturnsEmptyResultOnOverflow|TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable|TestAccountingDomainRejectsCheckedAccumulationOverflow)$'` | Passed: `ok github.com/kozz36/git-change-evidence 0.001s`. |
| Supplemental-test baseline (not RED) | Required PR8 focused command before the supplemental test edit | Passed (`cached`): production was already correct and no production file was changed. |
| Additional negative verification | Required PR8 focused command after adding the late-overflow and differential rows | Passed: those new rows are supplemental verification of existing behavior, not a new production TDD cycle. |
| Supplemental test refactor | `gofmt -w result_v1_observations_test.go result_v1_observations_precision_test.go`, then required PR8 and legacy focused commands | Passed (`cached`): format-only review made no production change. |

### Verification

- `go test ./...` passed for root, CLI, Git, inventory, and publication packages.
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` passed.
- `git diff --check` passed.
- `go test -race ./...` is N/A: this pure derivation adds no concurrency-bearing behavior.

### Layout, workload, and PR boundary

The user separately approved this bounded PR8 root-package addition. The root had 38 Go files before this slice and has 41 after it; the three added co-located files are `result_v1_observations.go` (84 lines), `result_v1_observations_test.go` (92 lines), and `result_v1_observations_precision_test.go` (58 lines). They add 234 Go lines with no Go deletions, every file is below the 160-line threshold, and three changed Go files remain below the eight-file package threshold. The optional precision test split follows the exact-comparison responsibility boundary rather than mechanically spreading coverage.

This is an uncommitted PR8/Unit 3 candidate only. Its actual delta against `origin/main` is 289 additions plus 1 deletion (290 total): 234 Go additions, 55 documentation additions, and 1 documentation deletion. It remains within the approved 400-total-changed-line slice budget. No delivery, commit, push, merge, or lifecycle settlement occurred; the parent owns independent checking and settlement.

### Files changed

- `result_v1_observations.go`
- `result_v1_observations_test.go`
- `result_v1_observations_precision_test.go`
- `openspec/changes/result-v1/tasks.md`
- `openspec/changes/result-v1/apply-progress.md`

### Remaining work and risks

- Units 4–9 remain pending; PR8 deliberately does not add `NewAccountingResultV1`, antecedent/snapshot validation, decoder, or report-binding behavior.
- The strict-TDD support file `.pi/gentle-ai/support/strict-tdd.md` remains absent, so the configured RED → GREEN → TRIANGULATE → REFACTOR contract was applied directly.
- Parent verification must retain the uncommitted-candidate distinction and independently confirm final diff accounting before any lifecycle action.

## PR9 — Exact antecedents and immutable snapshot boundary

- **Work unit:** Unit 4 only, the approved `stacked-to-main` PR9 slice. PR8 was delivered at 2,639 changed lines by merged PR63 (`e89596b4df82b6943b56853096ceef367adcdc58`); this candidate remains uncommitted for parent-owned independent checking and settlement.
- **Completed:** exact Policy V1 → Inventory V1 canonical revalidation with recomputed SHA-256 links; copied decoded policy material; native same-width snapshot revision validation; raw-byte primary path validation and duplicate rejection; and fail-atomic private validation outputs.
- **Excluded:** Result construction, ordering/classification integration, public ownership integration, decoder, report binding, legacy/Policy/Inventory/snapshot edits, shell work, and Units 5–9 remain deferred.

### Completed tasks

- [x] Completed Unit 4 in `tasks.md`; Units 5–9 remain untouched and pending.
- [x] Added private `validateResultAntecedents` and `validateResultSnapshot` without public constructors or decoder integration.
- [x] Added table-driven exact-chain, hostile raw-path, metadata-exclusion, atomic-failure, and defensive private-output tests.

### TDD Cycle Evidence

The strict-TDD support file `.pi/gentle-ai/support/strict-tdd.md` is absent; the configured RED → GREEN → TRIANGULATE → REFACTOR contract was applied directly. Every Go command used `GOTOOLCHAIN=go1.25.10` with that toolchain's `GOROOT/bin` first on `PATH`.

| Stage | Command | Actual result |
|---|---|---|
| RED | `go test . -run '^(TestValidateResultAntecedentsRequiresExactPolicyInventoryChain|TestValidateResultSnapshotAcceptsNativeIDsAndRejectsInvalidPaths)$'` | Failed as expected (exit 1): both private validators were undefined. |
| GREEN | Same focused command | Passed: `ok github.com/kozz36/git-change-evidence 0.001s`. |
| TRIANGULATE | `go test . -run '^TestValidateResult(Antecedents|Snapshot)'` | Passed after zero/stale/substituted/noncanonical antecedent, cross-width/invalid-path, late-failure, raw-byte, and ownership vectors. |
| REFACTOR | `gofmt -w result_v1_antecedents.go result_v1_antecedents_test.go`, then the required focused command | Passed: format-only refactor; no behavior change. |

### Verification

- `go test ./...` passed for root, CLI, Git, inventory, and publication packages.
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` passed.
- `git diff --check` passed.
- `go test -race ./...` is N/A: this synchronous pure validation slice introduces no concurrency-bearing behavior.

### Layout, workload, and PR boundary

PR9 adds two root-package Go files: `result_v1_antecedents.go` (64 lines) and `result_v1_antecedents_test.go` (142 lines). Both are below the 160-line threshold; the root increases from 41 to 43 Go files, and two changed Go files remain below the eight-file package threshold. This uncommitted candidate is 255 changed lines (254 additions, 1 deletion): 206 Go additions plus 48 documentation additions and 1 documentation deletion. This is the bounded antecedent/snapshot validation seam only; no size exception is used. The historical 4,495–5,125 forecast remains context rather than a guarantee; its parent-provided arithmetic remaining range after delivered PR8 is 1,856–2,486 changed lines.

### Files changed

- `result_v1_antecedents.go`
- `result_v1_antecedents_test.go`
- `openspec/changes/result-v1/tasks.md`
- `openspec/changes/result-v1/apply-progress.md`

### Remaining work and risks

- Units 5–9 remain pending; the private outputs are deliberately not wired into a constructor or public Result V1 state.
- The strict-TDD support file remains absent; no concurrency behavior was introduced.
- Parent must independently check this uncommitted PR9 projection before lifecycle settlement; no commit, push, merge, or delivery action occurred.

## PR9 corrective evidence — antecedent test closure

- **Scope:** Corrective verification only for the existing Unit 4 private antecedent/snapshot seam; no production, configuration, deferred API, or lifecycle surface changed.
- **Findings closed:** strengthened returned-value ownership, failure atomicity, supplied-inventory link binding, and PR9 approval/accounting evidence.
- **Delivery state:** this PR9 candidate is uncommitted and not delivered; parent-owned independent checking remains required.

### Corrective finding coverage

1. Antecedent ownership now independently checks the original policy view, same-chain later view, fresh decode from caller-cached canonical material, and canonical-accessor buffer mutations against preserved bodies and identities.
2. Every negative antecedent row captures policy view and both links, then asserts the typed zero view and empty digests; this includes valid policy progress followed by a late stale supplied-inventory digest.
3. An isolated valid-bytes/correct-digest inventory fixture changes only its cached supplied policy link to another valid digest. Strict inventory decoding otherwise succeeds, so this proves the validator's additional supplied-link predicate.
4. The cohesive snapshot boundary test and its assertion helper moved to `result_v1_snapshot_test.go` because the antecedent coverage otherwise exceeded the 160-line Go-file limit; this is a responsibility split, not compression.

### TDD Cycle Evidence

The historical PR9 RED/GREEN/TRIANGULATE/REFACTOR record above is retained. The strict-TDD support file remains absent. These are supplemental tests against independently inspected, already-correct production, so no fabricated RED or production change was made.

| Stage | Command | Actual result |
|---|---|---|
| Supplemental verification | `go test . -run '^TestValidateResult(Antecedents|Snapshot)'` | Passed: added negative and ownership vectors exercise all four corrective findings. |
| REFACTOR | `gofmt -w result_v1_antecedents_test.go result_v1_snapshot_test.go`, then the focused command | Passed (`cached`); cohesive test-file move only. |

### Verification

- Focused validator tests passed after the supplemental coverage and after the format-only refactor.
- `go test ./...` passed.
- The configured check-only `gofmt` gate passed.
- `git diff --check` passed.
- `go test -race ./...` is N/A: no concurrency-bearing behavior was introduced.

### Layout, workload, approval, and PR boundary

The user separately approved PR9's bounded root-package addition under the chosen 400-total-changed-line budget and approved the root/worktree issue #64 context. That approval is specific to this PR9 work unit; it neither inherits nor grants blanket approval from historical PR6–PR8 review context.

The exact current PR9 candidate is 349 changed lines: 252 Go additions (`result_v1_antecedents.go` 64, `result_v1_antecedents_test.go` 108, `result_v1_snapshot_test.go` 80) plus documentation 96 additions and 1 deletion. Every changed Go file is at or below 160 lines; three root-package Go files are below the eight-file package limit. Recorded delivery remains 2,639 changed lines through merged PR8; projected delivery is 2,639 + 349 = 2,988 changed lines if this candidate is later accepted. PR9 is not delivered.

### Files changed

- `result_v1_antecedents.go` (pre-existing PR9 candidate production)
- `result_v1_antecedents_test.go` (corrected antecedent coverage)
- `result_v1_snapshot_test.go` (cohesive snapshot coverage split)
- `openspec/changes/result-v1/apply-progress.md`

### Remaining risks

- Unit 5 public constructor/ownership tests remain out of scope; this private validator coverage does not claim their public behavior.
- Parent must independently verify the uncommitted candidate and preserve the separately scoped PR9 approval before any lifecycle action.
