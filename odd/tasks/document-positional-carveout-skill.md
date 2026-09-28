# Document positional carveout evidence in the GCE consumer skill

Issue: https://github.com/kozz36/git-change-evidence/issues/137
Branch: `docs/gce-carveout-positional-reference` (worktree `gce-carveout-reference`)

## Objective and boundary
The existing three-positional CLI returns technical counts for one source/new-file pair, but the consumer skill currently only routes census and Go APIs. Add accurate consumer guidance and a negative inference scenario, then observe one actual agent run. This does not change the CLI, publish a release, certify the SD-7 table, authorize a CNSIC pilot, or grant delivery authority. Preserve argv/stdout/stderr/exit status and immutable revision IDs as separate observations.

## Work and acceptance
- [x] T1 — Confirm command behavior, skill surfaces, duplicate issue search, and isolate a clean worktree. Evidence: `main.go` positional handling; issue #22 is command implementation, not skill guidance; #137 created with readback; worktree starts at `fbdaa78`.
- [x] T2 — Update skill, consumer reference, and one synthetic pair-only scenario. The delegated writer changed exactly three files, and `git diff --check` and introduced Markdown links passed. Strict TDD disabled for this docs-only task by explicit user choice. One pre-existing broken scenario-spec link remains outside this scope.
- [x] T3 — Focused CLI test passed (exit 0), introduced links and diff check passed; a fresh Pi headless evaluation actually loaded the skill, read one reference, and refused pair-to-table/approval inference. Exact prompt, output, trace hash and **observed out-of-scope reference read** are in `skills/gce-evidence/assets/positional-agent-evaluation.md`. One independent verifier found no factual guidance defect but did not inspect the private trace; an existing scenario-spec link remains broken.
- [ ] T4 — Staged five-file diff passed `git diff --cached --check` (57 authored changed lines); native assessment: medium, RDD off, independent verifier required and completed. Work-unit commit created with a conventional message. PR delivery remains blocked: issue #137 has no `status:approved`, which GCE PR policy requires; no approval label, push, or PR was attempted. Ask-on-risk remains the delivery strategy.

## Verification and recovery
Exact focused check: `go test ./cmd/git-change-evidence -run 'TestRunCLIPreserves(ThreeArgumentSpecialBaseReferences|ImmutableBaseForGoASTPath)$' -count=1`. Check relative links and `git diff --check`; read back full changed files. Real-agent run must distinguish skill loading from a fixture prompt. Runtime harness: actual agent trace for the pair-level interpretation (not a CLI change). Rollback boundary: this task's skill/reference/scenario/evaluation changes and this task file only; no CLI or archive rewrite.

Memory mirror pending: this Pi session is bound to CNSIC; Engram refused a GCE-project write. Local task document is the surviving progress copy. Current next action: obtain target-host-confirmed maintainer approval of issue #137 through the owner's normal workflow, re-read its label, then decide whether to push/open a PR under GCE policy. Independent read-only verifier found no factual guidance defect. Preserve the unrelated pre-existing broken scenario link as a disclosed limitation, not a silent fix.
