# Apply Progress: Public Preview readiness

## Scope and authority

- Active change: `public-preview-readiness`; native `gentle-ai.sdd-status` v2 consumed before edits.
- Native state: `applyState: ready`, `nextRecommended: apply`, dependencies unblocked.
- `actionContext`: `repo-local` workspace `/data/Projects/git-change-evidence-worktrees/public-preview`; allowed root is the workspace root. No action-context warnings.
- Session slice: PR 1 CLI → PR 0 planning, tasks 1–5 only. Tasks 6–9 are explicitly out of scope.
- Delivery: `auto-chain`, `stacked-to-main`, review budget 400 changed lines; the corrected local chain is PR 0 planning → `main`, PR 1 CLI → PR 0, PR 2 public entrypoint/license → PR 1, PR 3 governance/consumer guide → PR 2, and PR 4 audit/spec reconciliation → PR 3. No `size:exception` or push is authorized.

## Task progress

- [x] Task 1 — measured before implementation edits; persisted checkbox updated in `tasks.md`.
- [x] Task 2 — RED tests added in both allowed co-located test files; persisted checkbox updated after the focused command failed only for missing production helpers.
- [x] Task 3 — added the CLI-local version constant, exact discovery helper/writer, and allocated canonical normalization at both dispatch boundaries; persisted checkbox updated after GREEN.
- [x] Task 4 — added canonical/legacy success and failure parity coverage, retained unrelated positional behavior, and completed the named-binary smoke matrix; persisted checkbox updated after TRIANGULATE.
- [x] Task 5 — no structural refactor was needed beyond check-only formatting; focused, primary, and configured format checks passed after gofmt. Persisted checkbox updated after REFACTOR.
- [ ] Tasks 6–9 — deliberately not executed in this PR 1 slice.

## Pre-edit threshold census (Task 1)

Measurement time: before any implementation or test edit.

| Measurement | Additions | Deletions | Changed lines | Result |
|---|---:|---:|---:|---|
| PR 1 attributable working-tree implementation delta | 0 | 0 | 0 | Within 400 |
| `cmd/git-change-evidence/main.go` | 0 | 0 | 0 | Within 160 |
| `cmd/git-change-evidence/main_test.go` | 0 | 0 | 0 | Within 160 |
| `cmd/git-change-evidence/census_test.go` | 0 | 0 | 0 | Within 160 |

`git merge-base main HEAD...HEAD` includes inherited branch history and is not a PR 1 work-unit measurement. The actual tracked working tree before this slice had only pre-existing `.atl/skill-registry.md` (95 additions, 42 deletions) and `.gitignore` (2 additions); both are outside authorized edit surfaces and remain untouched. The planning artifacts predate this implementation slice. Later complete measurement found 765 OpenSpec planning lines that the implementation forecast had omitted; combining those artifacts with the code surface would produce an approximately 1,049-line review. The owner therefore assigned the planning artifacts to local PR 0 rather than folding them into PR 1. This delivery correction is not a `size:exception` and grants no push authority. No license, documentation, mechanical, audit, or Go implementation delta has yet been created for PR 1.

Configured thresholds inspected in `openspec/config.yaml`: CLI package contains seven Go files against `cli_go_files: 6`; this pre-existing count is recorded as an existing threshold condition. This unit adds no Go files, so it does not cross the file-count rule or request a package-boundary review. Per-file changes will be remeasured before every write and after each task; work must stop before any existing allowed Go file exceeds 160 additions plus deletions.

## TDD Cycle Evidence

| Task | Phase | Test / command | Result |
|---|---|---|---|
| 1 | Prerequisite | Git working-tree and Go-file numstat census | Passed: zero PR 1 implementation/test delta before edits |
| 2 | RED | `go test ./cmd/git-change-evidence -run 'Test.*(Census|Discovery|Version|Normalize|RunCommand|RunCLI)'` | Expected failure: only undefined `normalizeCensusArgs` and `discoveryText` production helpers; no production code was edited first |
| 3 | GREEN | `go test ./cmd/git-change-evidence -run 'Test.*(Census|Discovery|Version|Normalize|RunCommand|RunCLI)'` | Passed after the minimal shell-only implementation; the test-only canonical-prefix predicate was corrected to exclude unrelated `census other` input |
| 4 | TRIANGULATE | Focused command above, rerun after census success/failure parity and positional-boundary tests | Passed twice: canonical/legacy matched successful output bytes and failure exit/stdout/stderr for valid, malformed, unavailable, parser, writer, and resource-bound cases |
| 4 | TRIANGULATE smoke | `tmpdir="$(mktemp -d)" && go build -o "$tmpdir/gce" ./cmd/git-change-evidence && go build -o "$tmpdir/git-change-evidence" ./cmd/git-change-evidence` plus both-name help/version/canonical/legacy success/failure matrix | Passed; temporary directory removed by shell trap |
| 5 | REFACTOR | Focused command; `go test ./...`; configured check-only `gofmt` command | Passed after formatting the three allowed changed Go files; no behavioral refactor was needed |

## Files changed

- `openspec/changes/public-preview-readiness/tasks.md` — persisted task 1–5 checkbox updates only.
- `openspec/changes/public-preview-readiness/apply-progress.md` — cumulative progress record.
- `cmd/git-change-evidence/main_test.go` — RED coverage for canonical malformed dispatch, exact discovery, writer failures, negative discovery, and retained unrelated positional behavior.
- `cmd/git-change-evidence/census_test.go` — RED/GREEN/TRIANGULATE coverage for normalization, caller-slice nonmutation, and canonical/legacy valid/failing census parity.
- `cmd/git-change-evidence/main.go` — CLI-local preview version, exact discovery, discovery writer handling, and idempotent canonical normalization before census dispatch.

## Verification

The focused RED command failed as expected for absent production helpers; focused GREEN, TRIANGULATE, and REFACTOR commands passed. The configured primary `go test ./...` and exact check-only formatting command passed. The final Go threshold census is `main.go` 58, `main_test.go` 111 (110 additions + 1 deletion), and `census_test.go` 99 changed lines; all remain below 160. Race testing is N/A: this slice added no goroutines, synchronization, or concurrency-dependent behavior.

## Workload and PR boundary

- PR 1 CLI → PR 0 planning is complete for tasks 1–5 only; no task 6–9 work was performed. After the surgical verifier correction, its attributable review surface is 379 changed lines (286 Go, 83 progress lines, 10 task-checkbox additions/deletions), within 400.
- No file was added, so the observed existing `cli_go_files: 7` versus configured 6 threshold did not cause a new crossing or require a package-boundary review.
- No per-file 160-line threshold was crossed; no `size:exception` is needed or claimed.

## Deviations and risks

- Delivery correction: 765 OpenSpec planning lines were absent from the implementation forecast; reviewing code plus those artifacts together would be approximately 1,049 lines. The owner inserted local PR 0 planning → `main` and retargeted PR 1 CLI → PR 0, followed by PR 2 → PR 1, PR 3 → PR 2, and PR 4 → PR 3.
- This correction preserves review scope; it is not a `size:exception` and does not authorize branches, commits, pushes, or PR creation.
- Runtime `.atl`/`.gitignore` metadata is outside authorization and is restored when agent tooling regenerates it.
- Stop and return partial for package-boundary review before any existing allowed Go file reaches more than 160 changed lines.

## Remaining tasks

- [ ] 6. **Create the accurate public command entrypoint and Apache license.**
- [ ] 7. **Reconcile current governance and author the separate consumer guide.**
- [ ] 8. **Complete the bounded publication audit and narrow active-OpenSpec reconciliation.**
- [ ] 9. **Perform the implementation-slice review for the selected stacked-PR sequence.**

## Surgical correction: duplicate-option parity fixture

- An independent verifier found that the prior `duplicate` row appended `--receiver os` after `--`, exercising path validation rather than duplicate-option rejection.
- The row now places the repeated `--receiver os` before the separator and keeps `match.go` after `--`; an isolation assertion proves parsing fails with the duplicate and succeeds after removing only that pair.
- `cmd/git-change-evidence/census_test.go` measured 99 changed lines before and 117 after this correction, below the 160-line limit.
- Passed: focused census/discovery command, `go test ./...`, and the exact configured check-only formatting command. Race is N/A because no concurrency behavior changed.
- Revised PR 1 attributable review surface: 379 changed lines (286 Go, 83 progress, 10 task-checkbox additions/deletions), within the 400-line budget. Tasks 1–5 remain complete; tasks 6–9 remain untouched.
