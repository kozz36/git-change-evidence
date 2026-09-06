# Apply Progress: Result V1

## Current slice

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
