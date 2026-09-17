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

- [x] 6. **Create the accurate public command entrypoint and Apache license.**
- [ ] 7. **Reconcile current governance and author the separate consumer guide.**
- [ ] 8. **Complete the bounded publication audit and narrow active-OpenSpec reconciliation.**
- [ ] 9. **Perform the implementation-slice review for the selected stacked-PR sequence.**

## Surgical correction: duplicate-option parity fixture

- An independent verifier found that the prior `duplicate` row appended `--receiver os` after `--`, exercising path validation rather than duplicate-option rejection.
- The row now places the repeated `--receiver os` before the separator and keeps `match.go` after `--`; an isolation assertion proves parsing fails with the duplicate and succeeds after removing only that pair.
- `cmd/git-change-evidence/census_test.go` measured 99 changed lines before and 117 after this correction, below the 160-line limit.
- Passed: focused census/discovery command, `go test ./...`, and the exact configured check-only formatting command. Race is N/A because no concurrency behavior changed.
- Revised PR 1 attributable review surface: 379 changed lines (286 Go, 83 progress, 10 task-checkbox additions/deletions), within the 400-line budget. Tasks 1–5 remain complete; tasks 6–9 remain untouched.


## PR 2 public entrypoint and license (task 6)

- Native status consumed: `gentle-ai.sdd-status` v2, `applyState: ready`, `nextRecommended: apply`; `actionContext` is repo-local at the canonical workspace with its root authorized and no warnings.
- Completed: task 6; its persisted checkbox is now `[x]`. PR boundary: PR 2 → PR 1, task 6 only, with delivery `auto-chain` / `stacked-to-main` and a 400-line budget.
- Files changed: `LICENSE`, `README.md`, `docs/census.md`, this task record, and this cumulative progress record.
- License evidence: `LICENSE` is byte-identical to `https://www.apache.org/licenses/LICENSE-2.0.txt`, SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`, 202 lines, appendix present.
- Documentation evidence: manual/pattern review confirmed English portable prose, valid README links to `AGENTS.md` and `docs/census.md`, and a truthful future-facing (not broken) `docs/agent-consumers.md` reference because task 7 remains pending. Census guidance states canonical/legacy equivalence, supplied syntax-only scope, Linux descriptor confinement, and non-Linux unavailability.
- Commands passed: exact two `go build -o` commands into a removed temporary directory outside the checkout; both names ran help/version; canonical/legacy valid output was byte-identical and parsed as `git-change-evidence.census-go-ast-cli/v1`; canonical/legacy `--invalid` failures both returned exit 2 with empty stdout and `invalid input` stderr.
- TDD mode: skipped under the `default` rubric row for this docs/license-only slice; no production Go behavior changed and no strict-TDD cycle was invented.
- Workload: measured after all task-6 updates and progress-truth correction as 323 additions plus deletions (202 LICENSE, 41 README, 56 census guide, 2 task-checkbox, and 22 progress-record lines); below 400, no `size:exception`.
- Deviations: none. Runtime `.atl/skill-registry.md` and `.gitignore` were outside scope and left unchanged.
- Residual legal/audit boundaries: Apache-2.0 presence does not clear third-party attribution, IP ownership, confidentiality, export, NOTICE, or other legal questions; no copyright owner, NOTICE, legal-clearance, or third-party-rights claim was invented. External settings, history, assets, and legal review remain owner-controlled and uninspected by this task.

## Remaining tasks

- [ ] 7. **Reconcile current governance and author the separate consumer guide.** Update `AGENTS.md`, `openspec/config.yaml`, and `docs/product-charter.md` only to replace contradicted current-facing bootstrap statements with current facts: module/source/tests/CI exist; Apache-2.0 and the preview identity are selected; CLI-first positioning retains supported Go APIs; Go 1.25.10, Functional Core/Imperative Shell, layout, test commands, thresholds, evidence-only authority, and immutable CNSIC predecessor `a0fc7b26ff8a0e0a61baa586b32c46841611c806` remain. Preserve concrete configuration syntax, broad distribution, signing, and release automation as deferred; add a clear current-versus-historical boundary rather than rewriting bootstrap provenance. Create `docs/agent-consumers.md`, distinct from `AGENTS.md`, organized as consume → validate → reproduce → interpret. Cover canonical bytes/final LF/digest preservation, Base64 raw paths, revision-bound antecedents, census supplied-byte scope and non-atomic-tree limitation, stream/exit handling (including census 0/2/3/4 and documented existing command-specific exits), stale/incomplete/unavailable handling, independent regeneration, and consumer-owned delivery decisions. Verify all newly authored material is English and portable, every claim maps to code/specification, no material grants GCE delivery authority, and historical/local-path/literal fixtures remain unedited. **Depends on:** task 6. **Finish:** current governance and consumer guidance are accurate without erasing history or selecting excluded publication/distribution work. **Rollback:** revert only the four listed current-facing files, leaving historical OpenSpec and CNSIC provenance intact.
- [ ] 8. **Complete the bounded publication audit and narrow active-OpenSpec reconciliation.** Create `openspec/changes/public-preview-readiness/publication-audit.md` with per-surface rows for inspection method/scope or revision, observed fact, classification, coverage limitation, and owner follow-up. Classify current docs/examples, `.atl/skill-registry.md` (generated local-path metadata, not deletion authority), historical OpenSpec/CNSIC/local-path provenance, license/dependency/third-party questions, and any potential sensitive finding without reproducing secrets. Explicitly label Git history, unpublished refs, external settings/secrets/metadata, outside-tree generated or release assets, and legal/IP clearance as uninspected unless genuinely inspected under separate authority. Record owner-controlled decision points before any later visibility change, tag/release/asset publication, destructive remediation, credential action, or legal/confidentiality determination; these are informational consent boundaries, not GCE approvals or automated gates. Inspect active `openspec/specs/**/spec.md`; make only the justified prose-only update in `openspec/specs/inventory-v1-contract/spec.md` to distinguish its existing provenance sequence and WU-4C reference as historical while retaining the normative acquisition exclusion. Leave `openspec/specs/evidence-v1-provenance/spec.md` unchanged unless a concrete current-facing contradiction is found, and record the intentionally retained result. Do not modify `openspec/changes/archive/**`, prior validation/progress records, predecessor hashes, `.atl/skill-registry.md`, or literal historical evidence. **Depends on:** tasks 6–7. **Finish:** audit facts, benign context, actionable questions, and coverage gaps are distinguishable; only narrow active-spec wording changes occur. **Rollback:** preserve `publication-audit.md` as audit history even when implementation or findings change; append a dated correction, supersession, or withdrawal record instead of reverting or removing the audit, and revert only this task’s targeted active-spec prose. Preserve all history and audit trail.
- [ ] 9. **Perform the implementation-slice review for the selected stacked-PR sequence.** Recalculate additions plus deletions by unit and category, including the full Apache license, every documentation change, tests, mechanical reconciliation, and audit artifact; compare each selected PR with the 400-line budget without code-golfing or hiding formatting churn. Confirm each candidate unit still has its required tests/docs, focused command result, full `go test ./...` result, exact format-check result, applicable race-test result or recorded non-concurrency rationale, binary smoke result, documentation/license checklist, and independent rollback boundary. Recheck that no excluded work occurred: multilingual census, capabilities schema, evidence-schema/digest/API changes, internal/core policy changes, visibility change, tag/release/asset publication, history rewrite, deletion of `.atl` or archives, CNSIC cutover, MEH integration, packaging/signing/release automation, or delivery authority. **Depends on:** tasks 1–8. **Finish:** implementation evidence and actual review surface are ready for the selected stacked-PR reviews; no apply/verify/archive lifecycle claim, publication action, branch/commit/PR creation, or authority verdict is made by GCE. **Rollback:** use each unit’s independent rollback boundary rather than modifying historical artifacts or external state.

## PR 3 governance and consumer guidance (task 7)

- Native status consumed: `gentle-ai.sdd-status` v2; `applyState: ready`, `nextRecommended: apply`, canonical repo-local workspace authorized, and no `actionContext` warning.
- Delivery boundary: PR 3 → PR 2, task 7 only; `auto-chain`, `stacked-to-main`, review budget 400. Pre-write forecast was 249 additions plus 3 deletions (252 changed lines), below budget; final LOC tally is recorded after the completed diff check.
- Completed task 7 and persisted its `[x]` checkbox after observed documentation checks. Tasks 1–6 remain `[x]`; tasks 8–9 remain `[ ]`.
- Files changed: `AGENTS.md`, `docs/product-charter.md`, `docs/consumer-integration-guide.md`, `tasks.md`, and this cumulative progress record. No other authorized or runtime metadata file was changed.
- Governance evidence: current material records existing module/source/tests/CI, Apache-2.0, planned source/dev `v0.1.0-preview.1`, CLI/machine-contract primacy, both pre-v1 census names, supported public root/census Go APIs, Go 1.25.10 architecture/testing/thresholds, deferred configuration syntax/distribution/signing/release automation, and preserved CNSIC predecessor `a0fc7b26ff8a0e0a61baa586b32c46841611c806`.
- Consumer-guide evidence: added the authorized separate guide at `docs/consumer-integration-guide.md`, organized consume → validate → reproduce → interpret. It documents canonical bytes/final LF/digests, Base64 raw paths, revision-bound Policy/Inventory/Result/Report antecedents, census supplied bytes and non-atomic scope, streams, census exits 0/2/3/4, documented command-specific exits, stale/incomplete/unavailable handling, independent regeneration, and consumer-owned delivery decisions without assigning GCE authority.
- Checks passed: `git diff --check`; portable-path scan; new-guide Markdown-link target check; canonical and legacy command consistency across `AGENTS.md`, `README.md`, `docs/product-charter.md`, `docs/census.md`, and the new guide; and exact multilingual roadmap sequence checks. TDD was skipped under the configured default rubric because this slice changes documentation only; no Go behavior changed.
- Deviation: task text names `docs/agent-consumers.md` and `openspec/config.yaml`, but the explicit PR3 authorization required `docs/consumer-integration-guide.md` and forbade config edits. The authorized path was created; config was not modified.
- Residual: `README.md` remains outside authorized scope and still says `docs/agent-consumers.md` is planned (line 32); it does not link to `docs/consumer-integration-guide.md`. Tasks 8 and 9 remain intentionally unchecked. Runtime `.atl/skill-registry.md` and `.gitignore` were outside scope and remain unchanged.
- Final PR3 LOC: 176 additions plus 3 deletions = 179 changed lines, including the new guide and task/progress updates; within the 400-line budget.

## Remaining tasks

- [ ] 8. **Complete the bounded publication audit and narrow active-OpenSpec reconciliation.**
- [ ] 9. **Perform the implementation-slice review for the selected stacked-PR sequence.**

## PR 3 surgical verifier correction: AGENTS module/source statement

- Finding: the current `AGENTS.md` Go Layout section incorrectly said this planning slice had no module identity, `go.mod`, or Go source files.
- Correction: it now names `github.com/kozz36/git-change-evidence` and states that its `go.mod`, Go source files, and co-located Go tests are present, while retaining Go 1.25.10, layout, testing, architecture, and historical-artifact boundaries.
- Checks passed before retaining task 7 as `[x]`: `git diff --check`; `go.mod` and Go-source presence checks; exact module/baseline text checks; stale-assertion absence check; and persisted-task checkbox check.
- Revised PR3 LOC: 186 additions plus 4 deletions = 190 changed lines, matching the correction forecast and within the 400-line budget.
- Out-of-scope residuals: `openspec/config.yaml` still has stale CI/license deferrals, and `README.md` still names planned `docs/agent-consumers.md`; neither file was edited.
- Scope: only `AGENTS.md` and this cumulative progress record were changed for this correction; no code, phase, commit, branch, push, PR, publication, or visibility action occurred.

## PR 3 surgical verifier correction: README consumer-guide link

- Finding: `README.md` still described `docs/agent-consumers.md` as a future guide even though the authorized consumer guide exists at `docs/consumer-integration-guide.md`.
- Correction: replaced only that placeholder with a current Markdown link to [the Consumer Integration Guide](../../../docs/consumer-integration-guide.md).
- Checks passed: `git diff --check`; all local Markdown-link targets in `README.md` and `docs/consumer-integration-guide.md`; and canonical/legacy census command consistency across current guidance.
- Revised PR3 LOC: 196 additions plus 7 deletions = 203 changed lines, within the 400-line budget.
- Residual status: the README guide-path inconsistency is resolved. `openspec/config.yaml` retains stale CI/license deferrals and remains outside this correction's authorization.
- Scope: only `README.md` and this cumulative progress record were changed for this correction; no task, code, phase, commit, branch, push, PR, publication, tag, release, or visibility action occurred.
