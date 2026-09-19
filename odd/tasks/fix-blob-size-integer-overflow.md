# Fix Blob Size Integer Overflow

- **Issue:** #125
- **Parent program:** #104
- **Repository status:** public repository; no tag, release, package, asset, or publication is created by this work.
- **Root cause:** On 64-bit targets, `strconv.ParseUint` accepts a declared blob size of `math.MaxUint64`; converting that value to `int` overflows before the existing byte-cap comparison, allowing the oversized declaration to reach the blob-content read.

## Tasks

- [ ] Regression/TDD: add the `math.MaxUint64` declaration case to prove `BlobOversized` and no additional blob-content read.
- [ ] Minimal production fix: reject declared unsigned sizes above `MaxAcquiredBlobBytes` or the request byte cap before conversion to `int`.
- [ ] Verification/evidence: run the required focused and repository checks and record the observed results.

## TDD evidence

- **RED:** `go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1` failed before the production edit: `blob acquisition: racing_object; want oversized` for the `math.MaxUint64` declaration. The regression also records the no-additional-content-read requirement.
- **GREEN:** `go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1` passed after the production correction.

## Verification evidence

- `go test ./internal/git -run '^TestAcquireBlobRejectsObjectOverByteCapBeforeContentRead$' -count=1`: passed after the production correction.
- `go test ./internal/git`: passed.
- `go test ./...`: passed.
- `go test -race ./...`: passed.
- `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"`: passed.
- `git diff --check`: passed.
