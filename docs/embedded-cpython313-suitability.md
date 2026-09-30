# Embedded CPython 3.13 suitability — S1 (#151)

**Original S1 fixtures PASS; independent byte-validation mismatches discovered and corrected; parser suitability NOT PROVEN.** This private raw-byte probe is not a production census API, envelope, dependency choice, encoding policy, or delivery authority. S2/S3 are not implemented. Resource/cancellation, final host imports/isolation, and supply-chain provenance remain **UNKNOWN**. Earlier sample startup/heap measurements do not settle those axes.

## Experiment and results

The existing count smoke is unchanged. New `smoke/spans.go` accepts an explicitly supplied batch of paths/bytes, retains zero-match files, and returns file SHA-256 plus call/function half-open byte offsets, copied physical slices, and SHA-256 for each slice. Input order is retained; matches within a file are sorted by AST positions. There is no filesystem discovery or inventory validation policy. Paths are labels, not opened files.

A trusted helper parses and compiles each source before walking the AST. Neither the resulting code object nor scanned source is executed. The exact predicate is `Call(func=Attribute(value=Name(id='os'), attr='open'))`; shadowed `os` counts, while bare attributes, nested receivers, and indexed functions do not. Parentheses surrounding a function are absent from its AST Attribute span but remain in the enclosing call span. Nested direct calls each produce a match. Any file error returns a nil batch, discarding earlier matches.

UTF-8 AST columns are mapped through original physical line starts rather than treated as absolute raw offsets. The corrected helper receives original Python bytes, including BOM and line endings, for native parsing and compilation. Physical line mapping accounts for the runtime's BOM removal and CRLF/CR normalization without altering file hashes or returned slices. Tests compare literal expected fragments and original-byte offsets, including multiline parenthesized/nested calls with multibyte characters before and inside calls. Generic annotations/type aliases compile in a `.pyi` zero-match fixture. Malformed syntax and parse-valid top-level `return`/`nonlocal` fail an entire mixed valid/invalid batch. A scanned `raise RuntimeError` remains unexecuted.

| Encoding/line-ending fixture | Observed result |
| --- | --- |
| UTF-8 without cookie, LF / CRLF / CR | PASS: exact raw-byte spans and hashes |
| UTF-8 BOM, each of LF / CRLF / CR | PASS: same, including first-line calls |
| `coding: utf-8`; second-line `coding=utf_8` after shebang | PASS |
| Latin-1 cookie, including ASCII body | Explicit unsupported-encoding error; nil result |
| BOM plus Latin-1 cookie | Explicit unsupported-encoding error; nil result |
| Active unknown encoding cookie | Native syntax error; nil result |
| Invalid UTF-8 byte `ff` | Explicit unsupported-byte error; nil result |

This is deliberately conservative experimental rejection, **not** a selected production encoding contract. Other codec aliases and exhaustive grammar coverage are not demonstrated. After the correction below, original bytes cross the bridge as a copied Python bytes value. Native parse/compile validates byte semantics before matches; `tokenize.detect_encoding` then limits span mapping to observed UTF-8 spellings (`utf-8`, `utf-8-sig`, `utf8`). Non-UTF8 raw buffers remain explicitly unsupported by the Go mapping precheck. An active Latin-1 cookie remains unsupported even with ASCII contents; ignored second-line cookie comments do not select an encoding. No adversarial size/depth or hard cancellation guarantees follow. Bounded execution below confines this experiment; it is not proof of parser resource controls.

## Reproduction and TDD evidence

The observed run used a fresh external scratch directory and retained #149 artifacts outside the checkout. Inputs were copied read-only: `go-python` (matching generated bridge/stdlib) and `python-wasm/build/wasm2go/internal/internal/wasm2go` as `bundle`. No generated runtime is checked in and no root/nested committed module dependency changed. See [#149 input identities](embedded-cpython313-build-port.md).

The commands below preserve the observed command structure with host-specific locations expressed as variables. These portable forms were **not rerun**; log hashes and results describe the original observed run. From the target checkout, supply `PRIOR_ARTIFACTS` as the external directory containing the retained #149 assets and `SCRATCH_PARENT` as an external scratch parent:

```sh
GCE_ROOT=$(git rev-parse --show-toplevel)
prior_artifacts=${PRIOR_ARTIFACTS:?set to the retained #149 artifact directory}
scratch_parent=${SCRATCH_PARENT:?set to an external scratch parent}
scratch=$(mktemp -d "$scratch_parent/gce-embedded-313-151-s1-XXXXXX")
cd "$GCE_ROOT"
```

Preparation (performed once before RED in the observed run):

```sh
mkdir -p "$scratch"
cp -R "$prior_artifacts/go-python" "$scratch/go-python"
cp -R "$prior_artifacts/python-wasm/build/wasm2go/internal/internal/wasm2go" "$scratch/bundle"
cp -R internal/embeddedcpythonprobe/smoke "$scratch/smoke"
```

RED container command (observed output redirected to external scratch `red.log`):

```sh
timeout 1200 docker run --rm --network none --cpuset-cpus=0-3 --memory=8g -v "$scratch":/work -w /work/smoke -e GOPROXY=off -e GOSUMDB=off -e CGO_ENABLED=0 golang@sha256:154bd7001b6eb339e88c964442c0ad6ed5e53f09844cc818a41ce4ecb3ce3b43 bash -c 'go mod edit -replace=github.com/goccy/go-python=/work/go-python -replace=github.com/goccy/pythonwasm2go=/work/bundle && go mod tidy && go test -count=1 ./...'
```

Before every subsequent nested run, copy only current authored Go sources:

```sh
cp internal/embeddedcpythonprobe/smoke/*.go "$scratch/smoke/"
timeout 1200 docker run --rm --network none --cpuset-cpus=0-3 --memory=8g -v "$scratch":/work -w /work/smoke -e GOPROXY=off -e GOSUMDB=off -e CGO_ENABLED=0 golang@sha256:154bd7001b6eb339e88c964442c0ad6ed5e53f09844cc818a41ce4ecb3ce3b43 go test -count=1 ./...
```

Historical initial S1 evidence (before independent negative testing): RED exit 1: new tests could not compile (`undefined: probeBatch`, `probeInput`, `probeSpan`). GREEN exit 0: all eight nested tests passed. TRIANGULATE strengthened the span test to multiline calls on the BOM-bearing first line, preserving multibyte prefixes. The restored final run passed all eight tests. No concurrency-bearing surface added; no race claim. Formatting is check-only, no refactor needed.

| Log (external scratch) | SHA-256 |
| --- | --- |
| `red.log` | `960909285767d291f0e53cf8d355a6e3ca11cf62a1017d592ccabd3a5ad0df13` |
| `green.log` | `31cca5e5a7c4d5e2202788c150297b1582a40c8146f7aa40e21bea3df25f0c26` |
| `mutation.log` | `a238dc44e3542a47988ba8a5a350e83eb020776b0f14cc8bb47f5079b11b9491` |
| `restored.log` | `4d65c444c37acf4f4365ed9e7696777bb205f6aa898b144ed3666c9ad47abd63` |

### One declared mutation

Declared load-bearing site: raw physical offset anchoring, `at := begin + column`, in complete changed implementation file `spans.go`. Census command: `grep -c 'at := begin + column' internal/embeddedcpythonprobe/smoke/spans.go` → **1** site. Physically changed that site to `at := begin + column + 1`, copied it into scratch, and ran the same bounded container with exact test arguments:

```sh
timeout 1200 docker run --rm --network none --cpuset-cpus=0-3 --memory=8g -v "$scratch":/work -w /work/smoke -e GOPROXY=off -e GOSUMDB=off -e CGO_ENABLED=0 golang@sha256:154bd7001b6eb339e88c964442c0ad6ed5e53f09844cc818a41ce4ecb3ce3b43 go test -count=1 -run '^TestPhysicalSpansPinnedToRawBytesNeverNormalizedColumns$' ./...
```

| Mutation applied | Test that died | Sites / mutations |
| --- | --- | --- |
| Physical offsets shifted +1 byte | `TestPhysicalSpansPinnedToRawBytesNeverNormalizedColumns`: wrong literal physical fragments, exit 1 | 1 / 1 |

Mutation restored in repository and scratch before final nested run. This one-site demonstration does not claim exhaustive mutation coverage of other predicate/atomicity sites.

## Narrow byte-validation correction

Independent differential execution **PASS meant mismatches confirmed, not parity**. Original immutable evidence is retained externally: `independent-negative/results-readable.json` SHA-256 `66fbba0136b33cfb097cc2b1c04ff5ec24f460ccfb20f8d0ff354c3f9c498b9b`; `independent-negative/readable-run.log` SHA-256 `26edd7d361db94841534411f8ab3fc99db681880a8be09990b9319bbdba77cd1`. Both hashes were read back before correction. No independent files or historical S1 logs were modified.

| Exact independent fixture | Original probe / pinned `compile(bytes)` | Corrected probe |
| --- | --- | --- |
| BOM + `# coding: utf8\nos.open('x')\n` | Accepted 1 / SyntaxError `encoding problem: utf8 with BOM` | SyntaxError; nil batch |
| `os.open('x')\n# coding: made-up\n` | Rejected / accepted | Accepted 1 |
| `# coding: utf-8\n# coding: latin-1\nos.open('x')\n` | Rejected / accepted | Accepted 1 |
| BOM + `# coding: utf-8\nos.open('x')\n` | Both accepted | Accepted 1 |
| `# coding: utf8\nos.open('x')\n` without BOM | Both accepted | Accepted 1 |

The new regression test checks native `compile` on original bytes against each premise, then checks probe parity, exact fragments, and nil whole-batch failure after an earlier matching file. Code objects are never executed. It does not claim the independent verifier tested a Latin-1 comment after first-line code; its exact second fixture was `made-up`.

Correction RED killed all three mismatched subcases with the original code. The first full-suite correction attempt still rejected the no-BOM `utf8` control because `detect_encoding` returned `utf8`; that observed failure is preserved as `correction/green.log`. Explicitly supporting that spelling after native byte validation fixed the control; the full suite passed nine top-level tests. REFACTOR removed unused normalized-text construction; the restored full suite remained green.

Correction logs are separate from historical S1 logs:

| External log | SHA-256 |
| --- | --- |
| `correction/red.log` | `bf8e91bc54f1a21444f64046a5c8647694a739e937f7cb0ea4fcebd89acee910` |
| `correction/green.log` (failed intermediate attempt) | `566737c99fc8c44b9af4a6dd5d071f437d32415729bcb2ead93647da072b67fe` |
| `correction/green-final.log` | `4695b144145593d395ba6f81dc99f34c9221de0d92d12207c5f69baad230f32e` |
| `correction/mutation.log` | `2484754529b703f16c0e515b6585d36c6280b6c639b678602a4bd00570e1c4da` |
| `correction/mutation-bytes.log` | `6cf0e278ba5176593eb801abdf86834861486b5356d33ba64608a9440e7f80d5` |
| `correction/restored.log` | `d7c1fccb2ec3ca3f24ff118122c5af2c9c81622e548a420aab8d361bc7d3b5bf` |
| `correction/restored-final.log` | `75e9eed4334d1cae2247bd28dd5ea5ee568712f04453bdc8b74e128939236e36` |

Declared correction mutation site: the single bridge call supplying original bytes to `_probe_spans` in complete `spans.go`. Census: `grep -c 'module.CallMethod(ctx, "_probe_spans"' internal/embeddedcpythonprobe/smoke/spans.go` → **1**. An initial string-valued mutation killed the test with a bridge TypeError (`correction/mutation.log`), not a byte-validity parity failure. The behavior-level mutation at that same site retained the bytes type but stripped the BOM: `python.ValueOf([]byte(strings.TrimPrefix(string(input.Bytes), "\xef\xbb\xbf")))`. The targeted `bom_utf8_alias` subcase then failed because invalid BOM+utf8 was accepted (`correction/mutation-bytes.log`). Restored the original byte value before the final full-suite run. Site census: **1**; two mutation attempts at that same site, including **1 behavior-level BOM-stripping kill**.

Correction runs reused matching copied #149 dependencies in external scratch, refreshing only current authored Go files. Portable command structure below substitutes that observed external scratch location; these portable forms were not rerun. RED and mutation used the targeted test argument; GREEN and restored omitted `-run` for the full suite. Output was redirected to the corresponding separate `correction/*.log` file.

```sh
mkdir -p "$scratch/correction"
cp internal/embeddedcpythonprobe/smoke/*.go "$scratch/smoke/"
timeout 1200 docker run --rm --network none --cpuset-cpus=0-3 --memory=8g -v "$scratch":/work -w /work/smoke -e GOPROXY=off -e GOSUMDB=off -e CGO_ENABLED=0 golang@sha256:154bd7001b6eb339e88c964442c0ad6ed5e53f09844cc818a41ce4ecb3ce3b43 go test -count=1 -run '^TestByteValidationPinnedToOriginalBytesNeverDecodedCookieGuess$' ./...
timeout 1200 docker run --rm --network none --cpuset-cpus=0-3 --memory=8g -v "$scratch":/work -w /work/smoke -e GOPROXY=off -e GOSUMDB=off -e CGO_ENABLED=0 golang@sha256:154bd7001b6eb339e88c964442c0ad6ed5e53f09844cc818a41ce4ecb3ce3b43 go test -count=1 ./...
```

This correction settles only the named byte-validation fixtures. Resources/cancellation, final imports/isolation, and provenance remain **UNKNOWN**; overall suitability is **NOT PROVEN**.

## Independent corrected-candidate verification

A read-only verifier reran nine candidate tests plus a separate differential test against pinned CPython 3.13.13 `compile(bytes)`: exit 0, reported test duration 0.950s. All three regressions and two controls agreed; accepted inputs retained original-byte spans/hashes. Before/after comparisons bound tested copies to current `spans.go` SHA-256 `3e7e06a13c594ffa26f69b9686e417815deb6182bccd8d8df33a49bfdc70544d` and `spans_test.go` SHA-256 `7561c434c699fba2644c9335be2b57919f485eb7af2c673b545349882a9a38e0`. Copied dependency assets matched #149. External independent results SHA-256 `6129baf78ed190c826616f669d2fa95f95a2fc9825ebc45ab4d1a092d0c8738d`; run log `5c8bd440b1914811a94ae27021fdafdfa2ce3b5fd2daf656010838dc40d235dc`.

The verifier used the same offline image/CPU/memory/timeout parameters above, with a fresh `independent-corrected` working directory, host UID and writable scratch Go cache. Historical mutation logs/hashes were checked but mutations were not independently rerun. This verifies representative S1 fixtures, not exhaustive language/encoding coverage, production atomicity, resource limits or suitability.

## Separate root checks

From target: `GOTOOLCHAIN=go1.25.10 go test -count=1 ./...` passed six root packages. Root tests exclude this nested smoke module; **only the explicit nested container runs above prove its tests**. `git diff --check` and check-only `gofmt -l internal/embeddedcpythonprobe/smoke/spans.go internal/embeddedcpythonprobe/smoke/spans_test.go` produce no output. No push, PR, merge, scanned-source execution, or production adoption follows from these checks; the parent records S1 in a local work-unit commit separately.
