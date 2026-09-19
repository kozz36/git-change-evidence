# Fix Blob Size Integer Overflow

- **Issue:** #125
- **Parent program:** #104
- **Repository status:** public repository; no tag, release, package, asset, or publication is created by this work.
- **Root cause:** On 32-bit and 64-bit targets, `strconv.ParseUint` accepts the platform-width unsigned maximum; converting that value to `int` overflows before the existing byte-cap comparison, allowing the oversized declaration to reach the blob-content read.

## Tasks

- [x] Regression/TDD: add the platform-width unsigned maximum declaration case to prove `BlobOversized` and no additional blob-content read. **Evidence:** implementation commit `b374aacdd0915d18763f12500f2feb1ac061b5be`, tree `93c944caf517c26f3e45c053869f41f5772d72b1`; on the observed 64-bit host, that maximum equaled `math.MaxUint64`, and focused RED was `blob acquisition: racing_object; want oversized`, then focused GREEN passed after the production correction.
- [x] Minimal production fix: reject declared unsigned sizes above `MaxAcquiredBlobBytes` or the request byte cap before conversion to `int`. **Evidence:** implementation commit `b374aacdd0915d18763f12500f2feb1ac061b5be`, tree `93c944caf517c26f3e45c053869f41f5772d72b1`; the guard runs before conversion and content read.
- [ ] Independent verification and delivery: independently verify the committed candidate, review and merge it through an issue-linked PR, and require CodeQL alert #1 to resolve from merged code without dismissal before repeating the publication gates.

## TDD evidence

- **RED:** `go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1` failed before the production edit: on the observed 64-bit host, the platform-width unsigned maximum equaled `math.MaxUint64` and produced `blob acquisition: racing_object; want oversized`. The regression also records the no-additional-content-read requirement.
- **GREEN:** `go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1` passed after the production correction.

## Verification evidence

- `go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1`: passed after the production correction.
- `go test ./internal/git`: passed.
- `go test ./...`: passed.
- `go test -race ./...`: passed.
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"`: passed.
- `git diff --check`: passed.
