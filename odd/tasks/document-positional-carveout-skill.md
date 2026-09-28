# Document positional carveout evidence in the GCE consumer skill

Issue: https://github.com/kozz36/git-change-evidence/issues/137
Branch: `docs/gce-carveout-positional-reference` (worktree `gce-carveout-reference`)

## Objective and boundary
The existing three-positional CLI returns technical counts for one source/new-file pair, but the consumer skill currently only routes census and Go APIs. Add accurate consumer guidance and a negative inference scenario, then observe one actual agent run. This does not change the CLI, publish a release, certify the SD-7 table, authorize a CNSIC pilot, or grant delivery authority. Preserve argv/stdout/stderr/exit status and immutable revision IDs as separate observations.

## Work and acceptance
- [x] T1 — Confirm command behavior, skill surfaces, duplicate issue search, and isolate a clean worktree. Evidence: `main.go` positional handling; issue #22 is command implementation, not skill guidance; #137 created with readback; worktree starts at `fbdaa78`.
- [x] T2 — Update skill, consumer reference, and one synthetic pair-only scenario. The delegated writer changed exactly three files, and `git diff --check` and introduced Markdown links passed. Strict TDD disabled for this docs-only task by explicit user choice. One pre-existing broken scenario-spec link remains outside this scope.
- [x] T3 — Focused CLI test passed (exit 0), introduced links and diff check passed; a fresh Pi headless evaluation actually loaded the skill, read one reference, and refused pair-to-table/approval inference. Exact prompt, output, trace hash and **observed out-of-scope reference read** are in `skills/gce-evidence/assets/positional-agent-evaluation.md`. One independent verifier found no factual guidance defect but did not inspect the private trace; an existing scenario-spec link remains broken.
- [x] T4 — The five-file work-unit commit `da24f1563a38adf701a3aca869fc109b45ecea81` passed staged diff check (57 authored changed lines); native assessment: medium, RDD off, independent verifier completed. Issue #137 gained `status:approved` on GitHub. PR #138 (`type:docs`) opened from the isolated branch; Analyze (actions), Analyze (go), CodeQL, Go checks, and Validate PR policy passed on `da24f15`. An earlier policy run was cancelled and superseded by a successful run. No merge or release attempted.

## Verification and recovery
Exact focused check: `go test ./cmd/git-change-evidence -run 'TestRunCLIPreserves(ThreeArgumentSpecialBaseReferences|ImmutableBaseForGoASTPath)$' -count=1`. Check relative links and `git diff --check`; read back full changed files. Real-agent run must distinguish skill loading from a fixture prompt. Runtime harness: actual agent trace for the pair-level interpretation (not a CLI change). Rollback boundary: this task's skill/reference/scenario/evaluation changes and this task file only; no CLI or archive rewrite.

Memory mirror pending: this Pi session is bound to CNSIC; Engram refused a GCE-project write. Local task document is the surviving progress copy. Delivery: https://github.com/kozz36/git-change-evidence/pull/138 — review/merge remains an owner decision, not an inference from evidence or passing CI. This task-document closure is a separate follow-up commit, so re-check its own PR checks before calling the branch green. Independent read-only verifier found no factual guidance defect. Preserve the unrelated pre-existing broken scenario link as a disclosed limitation, not a silent fix.
