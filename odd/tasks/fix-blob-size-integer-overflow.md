# Fix Blob Size Integer Overflow

- **Issue:** #125
- **Parent program:** #104
- **Repository status:** public repository; no tag, release, package, asset, or publication is created by this work.
- **Root cause:** On 32-bit and 64-bit targets, `strconv.ParseUint` accepts the platform-width unsigned maximum; converting that value to `int` overflows before the existing byte-cap comparison, allowing the oversized declaration to reach the blob-content read.

## Tasks

- [x] Regression/TDD: add the platform-width unsigned maximum declaration case to prove `BlobOversized` and no additional blob-content read. **Evidence:** implementation commit `b374aacdd0915d18763f12500f2feb1ac061b5be`, tree `93c944caf517c26f3e45c053869f41f5772d72b1`; on the observed 64-bit host, that maximum equaled `math.MaxUint64`, and focused RED was `blob acquisition: racing_object; want oversized`, then focused GREEN passed after the production correction.
- [x] Minimal production fix: reject declared unsigned sizes above `MaxAcquiredBlobBytes` or the request byte cap before conversion to `int`. **Evidence:** implementation commit `b374aacdd0915d18763f12500f2feb1ac061b5be`, tree `93c944caf517c26f3e45c053869f41f5772d72b1`; the guard runs before conversion and content read.
- [x] Independent verification and delivery: independently verify the committed candidate, review and merge it through an issue-linked PR, and require CodeQL alert #1 to resolve from merged code without dismissal before repeating the publication gates. **Evidence:** portable correction commit `247460912aa96e3ffb63e598db16a25ba2df804e`, tree `1eea89ae63cc8fee3dda063d1a15bb0e220c269b`; [PR #126](https://github.com/kozz36/git-change-evidence/pull/126) passed Go, policy, and CodeQL checks and merged as `3a0d058dca093947fed95a8065175d70a5e45d5b`, with that same tree and parents `362bbb39cf0668fd5ddb106e906ae0b8be3659c2` + `247460912aa96e3ffb63e598db16a25ba2df804e`; main CodeQL run `35418425010` succeeded and alert #1 became `fixed` at `2026-09-19T03:24:35Z` with no dismissal.

## TDD evidence

- **RED:** `go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1` failed before the production edit: on the observed 64-bit host, the platform-width unsigned maximum equaled `math.MaxUint64` and produced `blob acquisition: racing_object; want oversized`. The regression also records the no-additional-content-read requirement.
- **GREEN:** `go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1` passed after the production correction.

## Verification evidence

- `go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1`: passed after the production correction.
- `go test ./internal/git`: passed.
- `go test ./...`: passed.
- `go test -race ./...`: passed.
- `GOARCH=386 go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1`: compiled, executed, and passed.
- `GOARCH=386 go test ./internal/git -count=1`: compiled, executed, and passed.
- `go vet ./...`: passed.
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"`: passed.
- `git diff --check`: passed.
- Independent committed-candidate verification passed after correcting the initial non-portable test fixture; the final fixture uses the platform-width unsigned maximum and retains kill power on 32-bit and 64-bit targets.
- PR #126 checks passed: `Go checks` run `35417941106`, `Validate PR policy` run `35417941125`, pull-request CodeQL analysis run `35417939562` with zero results, and CodeQL check run `105830200916`.
- Main CodeQL run `35418425010` passed after merge; alert #1 is fixed without dismissal. The repository remains public with no tag, release, package, asset, or upload.
