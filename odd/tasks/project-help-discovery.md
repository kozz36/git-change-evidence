# Add exact project-help discovery

Issue: https://github.com/kozz36/git-change-evidence/issues/139
Branch/worktree: `fix/project-help` at `/data/Projects/git-change-evidence-worktrees/gce-project-help`; base `c29c6505802750b05197c52693e558bffa0006c8`.

## Objective and boundary
`project --help` currently enters the format parser and exits 2 with `invalid input`; root and census help use a side-effect-free discovery table. Extend that existing pattern only for exact `project --help`: stdout synopsis and `canonical-json`/`human-text`, empty stderr, exit 0, no stdin/CWD/Git/publication access. Do not add `-h`, change invalid formats or projections, release, or start any CNSIC integration.

## Work and acceptance
- [x] P1 — Reproduce and identify root cause on current origin/main. Observed `go run ./cmd/git-change-evidence project --help`: stdout empty, stderr `invalid input\nexit status 2\n`, wrapper exit 1/CLI exit 2; root `--help` prints usage, stderr empty, wrapper exit 0. Issue #139 created with readback; new worktree isolated from dirty root.
- [x] P2 — Delegated writer added `TestDiscoveryProjectHelpDispatchIsSideEffectFree` for runCommand/runCLI, exact stdout/empty stderr/exit0 and zero stdin/CWD/Git/publication use. RED observed before production edit: `go test ./cmd/git-change-evidence -run '^TestDiscoveryProjectHelpDispatchIsSideEffectFree$' -count=1` exited 1 (both dispatchers returned exit2/empty stdout/invalid input). Existing invalid-format and 3-positional tests preserved.
- [x] P3 — Minimal `projectHelp` constant and exact discoveryText case implemented; focused GREEN and full `go test ./...` exited 0, gofmt/diff checks empty, runtime `go run ./cmd/git-change-evidence project --help` returned stdout help/empty stderr/exit0. TRIANGULATE checked `project -h` and `project --help extra` remain undiscovered; no refactor needed. Independent verifier reran focused/full tests and runtime positive/negative (`project other`: inner exit2/empty stdout), found no defect; full Go run was cached, no physical mutant applied.
- [x] P4 — Work-unit commit `04206b16672350d4c3bce15199eb0ffcd80e57f5` includes code, regression and task doc; staged diff check passed (52 authored lines), native risk high/RDD off and independent verifier completed. Issue #139 gained `status:approved`. Branch pushed; PR #140 (`type:bug`, `Closes #139`) opened at head `83235af2d777742655abdc31f82f9e2cd9a9d8e1` and five checks passed. No merge or release. This task-document closure adds a later commit, so confirm the PR's current-head checks again.

## Verification
Exact focused runner: `go test ./cmd/git-change-evidence -run 'TestDiscovery|TestRunProjectRejectsInvalidInputWithoutWritingStdout' -count=1`; full `go test ./...`; `gofmt -l cmd/git-change-evidence/main.go cmd/git-change-evidence/main_test.go`; `git diff --check`; runtime `go run ./cmd/git-change-evidence project --help` with separate streams/status. Record RED/GREEN outputs, any skips/failures, rollback boundary (discovery case/constant and test only), and actual commit hash. No authority comes from CLI evidence.

Memory mirror pending: this Pi session is bound to CNSIC and cannot save the GCE-project Engram mirror. Preserve this local task file and synchronize from a GCE-bound Pi session later.
Delivery: https://github.com/kozz36/git-change-evidence/pull/140. Review/merge remains an explicit owner decision, not an inference from passing checks. Preserve the clean branch and worktree.
