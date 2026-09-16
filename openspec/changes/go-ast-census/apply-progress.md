# Apply progress: go-ast-census

## Status consumed

- Native `gentle-ai.sdd-status` v2 selected `go-ast-census` with `applyState: ready`, `nextRecommended: apply`, and all proposal, specification, design, and task artifacts complete.
- The authoritative workspace and only allowed edit root are `/data/Projects/git-change-evidence-worktrees/go-ast-census`.
- The parent supplied the active launch authorization and retains attempt settlement; this work did not acquire, settle, reset, or expose an attempt token.
- Delivery decision: `exception-ok (size:exception)` applies to the approved complete vertical. Chain strategy is not applicable. The hard 400 physical-line novel-production cap remains in force.

## TDD Cycle Evidence

| Cycle | RED evidence | GREEN evidence | TRIANGULATE / REFACTOR evidence |
|---|---|---|---|
| Complete GoASTV1 vertical | Before production files existed, `go test ./census` exited 1 with undefined `File`, `QueryV1`, `ResultV1`, and `SelectorCallQueryV1` symbols from the complete RED suite. | After both production files were added and gofmt was applied to owned Go files, `go test ./census` exited 0. | Independent earlier-valid-file parser-atomicity and input-permutation order cases were added, then `go test ./census` exited 0. Refactor review retained the clear API/extractor split, added consumer documentation, and a final focused rerun exited 0. |

## Verification evidence

- `go test ./census` — RED: exit 1, intentional absent end-to-end API; GREEN, triangulation, and refactor reruns: exit 0.
- `go test ./...` — exit 0 for root, `census`, CLI, and all internal packages.
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` — exit 0 with no output.
- `go test -race ./...` — N/A: the added API is pure and synchronous and introduces no concurrency-bearing behavior.

## Scenario evidence

- GAC-001-S01 and GAC-005-S03: `census/external_api_test.go` constructs root inventory values, asserts bindings and independent SHA-256 values, and mutates inputs and exported paths.
- GAC-001-S02; GAC-002-S01 through S06; GAC-006-S01 through S08; and GAC-007-S01: named `census/go_ast_test.go` validation, bounds, and atomicity rows assert typed field/code errors and zero results.
- GAC-003-S01 through S04: exact literal-oracle tests cover direct, parenthesized, indexed, multi-indexed, nested, shadowed, ambiguous indexed-value, decoy, repeated, and legacy calls.
- GAC-004-S01 and S02: independent UTF-8, CRLF, multiline, line-directive, shifted-position, and fragment/file hash assertions cover physical evidence.
- GAC-005-S01 and S02: mutation and raw byte order/permutation tests assert defensive ownership and deterministic output.
- GAC-007-S02: private physical-span and aggregate-arithmetic boundary tests cover no-position, foreign-file, past-EOF, reversed, and near-`uint64` conditions.
- GAC-008-S01: `docs/census.md` documents syntax-only scope, shadow/index ambiguity, semantic non-guarantees, evidence-only boundary, and lack of parser resource-isolation guarantees.
- GAC-009-S01: the only `AGENTS.md` change is the `census` exception and the only configuration change is `layout.public_api`; the testing rubric diff is empty.
- GAC-010-S01: measured novel production is 336 physical lines across the two approved files, under the hard 400 ceiling; the approved `exception-ok (size:exception)` boundary applies to the complete review unit.

## Threshold and package-boundary review

- `census/go_ast_test.go` first crossed the configured 160 changed-lines-per-Go-file threshold at 343 added physical lines during the complete RED scenario matrix. Its rationale is the required 30-scenario, independently-oracled behavior and private-boundary coverage; no acceptance logic is relocated into tests.
- Before adding the two production files after that crossing, the package boundary was reviewed: the approved narrow public `census` subpackage depends one way on the root inventory contracts, introduces no new package or source root, and keeps the API/extractor split confined to `census/query.go` and `census/go_ast.go`.
- `census/external_api_test.go` is at 51 added physical lines and has not crossed its per-file threshold.
- `census/go_ast.go` first crossed the configured 160 changed-lines-per-Go-file threshold at 256 added physical lines. Its rationale is the inseparable complete-file parse, wrapper normalization, validation, physical-span, atomic collection, and ordering pipeline required by the approved vertical; it remains the only extractor file and stays one-way from `census` to root inventory APIs.
- `census/query.go` is at 80 added physical lines and has not crossed its per-file threshold. No additional Go files or packages are planned after the reviewed two-file API/extractor split.

## Completed work and task evidence

- Created the complete RED suite before production, then added `census/query.go` and `census/go_ast.go` as the only production implementation files.
- Added `docs/census.md` and the narrow, synchronized public-layout exception in `AGENTS.md` and `openspec/config.yaml` before or with the public package introduction.
- Completed and persisted all five implementation-owned task rows: boundary alignment, RED, GREEN, TRIANGULATE, and REFACTOR are visibly marked `- [x]` in `tasks.md`.
- No unchecked implementation task rows remain.

## Files and measured workload

| Surface | Added lines | Deleted lines | Notes |
|---|---:|---:|---|
| `census/query.go` | 80 | 0 | Public contracts, defaults, result bindings, defensive accessors. |
| `census/go_ast.go` | 256 | 0 | Validation, whole-file AST traversal, physical evidence, ordering, atomic failures. |
| Production total | 336 | 0 | Physical lines after gofmt; within the hard 400-line cap. |
| `census/go_ast_test.go` | 368 | 0 | Co-located scenario and private-boundary coverage. |
| `census/external_api_test.go` | 51 | 0 | External value-API and ownership coverage. |
| `docs/census.md` | 39 | 0 | Consumer contract and non-guarantees. |
| `AGENTS.md`, `openspec/config.yaml` | 2 | 2 | Only the approved narrow layout wording. |
| `tasks.md` checkbox state | 5 | 5 | Five persisted completion markers only. |

The measured authored gross is 870 additions plus 7 deletions (877 total), including this 69-line cumulative progress artifact; it remains below the parent-authorized 1,800 gross cap. This single complete vertical is the approved `exception-ok (size:exception)` work-unit/PR boundary; no commit, PR, or delivery action was performed.

## Deviations, remaining work, and action context

- Design deviations: none. No semantic resolution, decoder, filesystem discovery, policy input, hunk API, concurrency, or additional package was introduced.
- Remaining tasks: none; there are no exact unchecked `- [ ]` lines in the persisted task artifact.
- Action-context warning: editing remained inside the authoritative repository root and the user-supplied allowed surfaces. The parent retains native lifecycle settlement and next-phase authority.

## Evidence correction — C4/C5/FAIL preservation

This cumulative section corrects conflicting earlier progress statements. In particular, the earlier claims that every task was complete, no unchecked tasks remained, and this actor neither acquired nor continued an attempt are superseded by the observed task artifact and this continuation. The completed original vertical remains recorded below as historical work; this section does not claim a current passing test as evidence of an earlier phase.

### Historical sources and chronology

- The retained original writer task record at `/home/kozz36/.pi/agent/gentle-agents/tasks/mu2l2jvk-4-cuc2.json` is available. Its ordered thread shows the complete RED suite written before production, then `go test ./census` failing with undefined `File`, `QueryV1`, `ResultV1`, and `SelectorCallQueryV1` (thread item 28, exit 1).
- That same record shows `census/query.go` and `census/go_ast.go` written after RED; a `gofmt -w` plus focused run passed (item 45); added triangulation cases plus focused run passed (item 49); then focused, primary, and check-only formatting commands passed (items 53, 55, and 57). Its final refactor-focused check plus check-only formatting also passed (item 77). The record is direct historical command evidence, not reconstructed chronology.
- The interrupted correction transcript at `/home/kozz36/.pi/agent/gentle-agents/sessions/2026-09-15T13-22-13-206Z_01a0a53b-7855-74b3-a86e-0659a4f97fce.jsonl` records `go test ./census` passing at `2026-09-15T13:26:08.952Z`. The supplied interruption context records two retained checkbox updates and an abort before the C5 configuration write; the pre-write read in this continuation confirmed that `rules.verify.manual_evidence` was absent.
- No standalone raw test output establishes a separate RED for C1/C2 or the C3 oracle correction. Those tests characterize already-correct production, and no temporary production mutation or invented failure was used. Their historical write ordering is therefore unavailable and recorded as N/A rather than inferred.

### TDD Cycle Evidence

| Work and applicability | RED (written / passed) | GREEN (written / passed) | TRIANGULATE (written / passed) | REFACTOR (written / passed) | Safety net and status |
|---|---|---|---|---|---|
| Original complete GoASTV1 vertical; Go source/test matched the authoritative strict-TDD rubric. | Written before production; passed as an intentional failing focused run (task-record item 28, exit 1). | Production written after RED; focused pass recorded at item 45. | Alternate atomicity/order cases written after GREEN; focused pass recorded at item 49. | API/extractor split and docs retained; focused plus formatting pass recorded at item 77. | Primary pass recorded at item 55; check-only formatting pass at items 57 and 77. Concurrency/race: N/A because the surface is synchronous and creates no concurrent behavior. |
| Retained C1/C2 characterization and C3 TRIANGULATE correction; Go tests were already written before this continuation and no production file is authorized. | N/A — behavior was already correct, and no synthetic regression was permitted. | Passed: interrupted focused run at `2026-09-15T13:26:08.952Z`; this continuation's focused `go test ./census` also passed. | Readback passed: `TestTriangulatedAtomicityAndOrder` independently builds literal `MatchV1` records with call/selector spans and SHA-256 values, guards exact count in both input orders, and compares full ordered records with `reflect.DeepEqual`; it does not fabricate lower-key ties. | N/A — no production or test refactor was performed by this continuation. | This continuation's `go test ./...` and the configured check-only formatting command both passed. Concurrency/race: N/A; no Go behavior changed. |
| C4 progress, C5 manual evidence, and historical FAIL copy; documentary/process artifacts only. | N/A — no Go source or test was written. | N/A — no new executable behavior. | N/A. | N/A. | Manual evidence is limited by `rules.verify.manual_evidence` to GAC-008-S01, GAC-009-S01, and GAC-010-S01; all behavioral scenarios remain under the unchanged strict-TDD rubric. |

### Retained test readback and current verification

- C1 is present in `TestGAC001AndGAC002Validation` as the explicit empty-selector `invalid_identifier` case. C2 is present in `TestGAC003Candidates` as a source comment containing matching `os.Open` text that yields zero matches.
- C3 is present in `TestTriangulatedAtomicityAndOrder`: each forward and backward result first has an exact-count guard, then is compared to independently constructed complete `MatchV1` records. The expected records include raw paths, whole-file and fragment digests, and call/selector positions; counts alone cannot satisfy the oracle.
- Current commands, run after native continuation acquired `proceed` and before documentary-only writes, all exited 0: `go test ./census`; `go test ./...`; and `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"`. No Go file changed after those commands, and no formatting write was performed.

### C5, FAIL preservation, task state, and lifecycle

- `openspec/config.yaml` now adds only `rules.verify.manual_evidence.go-ast-census.scenarios` for GAC-008-S01, GAC-009-S01, and GAC-010-S01. Testing commands, testing rubric rows, modes, and evidence methods were not edited.
- Before this correction changed any canonical report, the exact 175-line canonical FAIL report was copied to `verify-report.failed-initial.md`. `cmp -s` passed and both files SHA-256 to `f088b195421dc4ba682c6fa48b2d016de86715c64cee9d0a8db9493d80dba7e9`; the canonical report remains unchanged. Its FAIL verdict and five findings are retained verbatim in the copy.
- C4, C5, and Preserve historical FAIL are the three implementation-owned rows completed by this continuation. Their persisted task checkboxes are marked immediately after their corresponding artifact completion; no unchecked implementation row remains after the final task readback.
- The native continuation acquired `proceed` for work unit `go-ast-census-evidence-correction` with the supplied active attempt. Settlement remains parent-owned. The native passing-settle obligation requires remediation binding to `sha256:f088b195421dc4ba682c6fa48b2d016de86715c64cee9d0a8db9493d80dba7e9` and distinct current verification evidence; this continuation did not settle.
- Fresh untracked admission retained the original owned `census/`, `docs/census.md`, and `openspec/changes/go-ast-census/` artifacts, admitted only the new historical FAIL copy, and excluded ambient `.pi/gentle-ai/sdd-preflight.json`. No production, documentation, specification, rubric, command, delivery, or ambient `.pi` write was made by this continuation.
- Workload/PR boundary: this is the already-authorized evidence-correction work unit under its native 600-gross objective; it neither changes the 336-line frozen production total nor selects a delivery action. Independent final verification remains a separate phase and is not claimed complete here.

## Residual C3 correction — 2026-09-15

### Status, scope, and ownership

- Consumed native `gentle-ai.sdd-status` v2 for `go-ast-census`: `applyState: ready`, `nextRecommended: apply`, eight of nine tasks complete before this correction, no blockers, and only the repository root as `actionContext.allowedEditRoots`.
- Continued the active native attempt for `go-ast-census-residual-c3`; `sdd-attempt acquire` returned `state: proceed` with the inherited one-attempt, 250-gross-line objective. No reset, new harness, production mutation, delivery action, or external-test/configuration/specification/documentation change occurred.
- The only behavioral edit is the raw-path ordering subtest in `census/go_ast_test.go`: it captures `got := result.Matches()`, asserts `len(got) == len(want)` before iteration, and ranges over `got`.
- The current 126-line canonical FAIL report was copied byte-for-byte to `openspec/changes/go-ast-census/verify-report.failed-c3.md`; `cmp -s` passed and both files SHA-256 to `67973625eb60816fa3489713054e359914e3391dcb083a39e2b072914d7a8a97`. The earlier historical copy remains SHA-256 `f088b195421dc4ba682c6fa48b2d016de86715c64cee9d0a8db9493d80dba7e9`.

### TDD Cycle Evidence

| Work | RED | GREEN | TRIANGULATE | REFACTOR | Safety net |
|---|---|---|---|---|---|
| Residual C3 raw-path count guard | N/A: this strengthens an already-correct assertion, and no artificial failure or production mutation was permitted. | The exact isolated subtest passed before and after the guard was added. | The guard preserves the existing literal raw-path order oracle while preventing an empty result from satisfying the isolated loop vacuously. | N/A: no refactor was needed or performed. | `go test ./census`, `go test ./...`, and the configured check-only formatting command passed after the test edit. |

The authoritative rubric selected the `go-source-or-test` strict-TDD row by all-matching-row, strictest-wins, evidence-union resolution; the default row was not selected. The concurrency-bearing row and `go test -race ./...` are N/A because this test-only correction adds no concurrent behavior. All 10 requirements and 30 scenarios remain unchanged.

### Commands and task evidence

- `go test ./census -run '^TestGAC005OwnershipAndOrder$/^GAC-005-S02_raw_path_order_is_deterministic$'` — exit 0 before and after the correction.
- `go test ./census` — exit 0 after the correction.
- `go test ./...` — exit 0 after the correction.
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` — exit 0 after the correction with no output.
- Completed implementation-owned task: **TRIANGULATE**. Its sole reopened checkbox is updated in `tasks.md` immediately after this evidence write and must be reread before return.
- Remaining implementation tasks after that checkbox update: none; no exact unchecked `- [ ]` task line remains.

### Workload and next phase

- This is the one native `go-ast-census-residual-c3` correction work unit, bounded at 250 gross changed lines. Its review boundary excludes frozen product, configuration, external tests, specifications, documentation, and installed skills; independent verification remains a separately authorized future phase and is not claimed here.
- Design deviations: none. Native settlement must bind a distinct correction evidence revision to failed evidence `sha256:67973625eb60816fa3489713054e359914e3391dcb083a39e2b072914d7a8a97`.
