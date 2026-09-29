# Gotreesitter Python 3.13 feasibility (#145)

Status: **test-only, grammar/atomicity FAIL; do not adopt**. This isolated probe is not a census implementation, complete Python 3.13 certification, or delivery decision. The source of requirements is [the design](python-census-design.md). Root `go test ./...` does **not** run the nested module.

## Reproduction and provenance

Candidate `github.com/odvcencio/gotreesitter` tag `v0.55.1`, peeled Git commit `92db945f28de67be51de8235c9cd4e25a900f648`; nested `internal/gotreesitterprobe/go.mod` pins it with module sum `h1:/UabDnEwB4usEs/bd6UqRfxh5tf+ByWr3iGt+JjwS2I=` and go.mod sum `h1:hBVkghd0paaYAVwd2087vfwdeU984bQbMo9LvpE0moo=`. Its `grammars/languages.lock` line 15 pins tree-sitter-python `26855eabccb19c6abf499fbc5b8dc7cc9ab8bc64`; the embedded grammar blob is consumed via `grammars.PythonLanguage()`. This records upstream artifact identity, not independent grammar equivalence. The module declares Go 1.22.0; the probe targets GCE's Go 1.25.10.

Run `cd internal/gotreesitterprobe && GOTOOLCHAIN=go1.25.10 go test -count=1 -v ./...`. The test expectations characterize actual candidate behavior, so a green characterization test **does not** mean the feasibility gate passes. CPython oracle here is locally identified `python3.13` 3.13.13, development-only, never a required GCE runtime: compile each exact validity-table source with `compile(source, '<probe>', 'exec')`, catching `SyntaxError` (do not use only `ast.parse`, which admits top-level `return`). No CPython process is invoked by the Go tests. The invalid newline case is from [tree-sitter-python #178](https://github.com/tree-sitter/tree-sitter-python/issues/178).

## T1: syntax gate first

| Fixture | CPython 3.13 `compile` | Gotreesitter `ParseStrict` + `!HasErrorOrMissing` | Result |
| --- | --- | --- | --- |
| `os.open('x')` | valid | valid | agrees |
| `def foo(x)\n:\n    return x + 2\n` (#178) | invalid | **valid** | **FAIL** |
| `a = (\n  1 +\n  2\n)\n` (#178 valid counterpart) | valid | valid | agrees |
| `os.open(\n` | invalid | invalid | agrees |
| `1 = os.open('x')` | invalid | invalid | agrees |
| `return os.open('x')` at module top level | invalid | **valid** | **FAIL** |
| `def f[T](x: T) -> T:\n    return x\n` | valid | valid | agrees |

An invalid file with the #178 newline and `os.open('x')` in its body yields one candidate and no grammar error. `ParseStrict` only rejects *early stops*, not all CPython-invalid syntax. This is a direct atomic-validity failure: a clean tree cannot certify the design's rejection requirement. **Stop unsupported adoption**; no fallback runtime or silently broadened validity claim. Representative annotated/overload-decorated generic `.pyi` and async-call stub parses; broad `.pyi` coverage remains UNKNOWN. A syntactically malformed file appended to a valid call yields nil matches and an error for that fixture, but that local pass cannot offset the accepted-invalid files.

The exact-shape fixture returns direct, nested inner, and shadowed `os.open` calls, excluding attribute-only, `os.path.open` and indexed-function forms. A source with `π`, `é`, `猫` and mixed CRLF/LF yields half-open **physical byte** call spans `[11,34)`, `[19,33)`, `[89,101)` and exact slices, including nested calls. These are representative tests, not a versioned census envelope or proof of arbitrary encoding support. The experimental walker has no resource guarantee and must not be reused as production behavior.

## T2: bounded supply-chain and standalone observations

On linux/amd64, run from `internal/gotreesitterprobe`:

```sh
GOTOOLCHAIN=go1.25.10 go list -m all
GOTOOLCHAIN=go1.25.10 go list -deps .
CGO_ENABLED=0 GOTOOLCHAIN=go1.25.10 go test -c -o /tmp/gce-gotreesitterprobe.test .
file /tmp/gce-gotreesitterprobe.test
ldd /tmp/gce-gotreesitterprobe.test
env -i PATH=/nonexistent HOME=/nonexistent /tmp/gce-gotreesitterprobe.test -test.run '^TestPython313Validity$' -test.v
```

The module graph contains gotreesitter v0.55.1, `golang.org/x/sync v0.11.0`, `gopkg.in/yaml.v3 v3.0.1`, `github.com/kr/pretty v0.1.0`, and `gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127`; the compiled probe dependency list includes gotreesitter core and grammars, their internal packages, and standard library, not the other graph modules. `file` reported statically linked ELF; `ldd` reported `not a dynamic executable`. The isolated test passed with PATH/HOME pointing nowhere (no Python executable available via PATH). This is a **PASS for this host/test binary only**, not a cross-platform reproducible release.

The pinned archive SHA-256 is `8f4612523343296acf146aaf6d4bee1f75d7de27ab866e2370edf07f6a4fe599`; its embedded `grammars/grammar_blobs/python.bin` SHA-256 is `cde4a67dc6af6e1232dbbd1eab8618478d1d73727020e8a8002542390a452d37`, and `grammars/languages.lock` SHA-256 is `34fde4ff60bd77539d8b2afcbeaad41c78f496bfac409df1724b68d6508bb2f6`. The bundled `LICENSE` says MIT, copyright 2026 Oscar Villavicencio (SHA-256 `b174fbe1e1cffafb096528de3e8361c90a0d97e7e3f1a32061aae51e983a7cfc`); `THIRD_PARTY_NOTICES` covers separately vendored grammar material. These are provenance/attribution observations, **not legal clearance** or independent proof of blob-source equivalence.

**UNKNOWN:** parser-enforced worst-case time, hard heap/stack bounds, cancellation/timeout behavior under adversarial legal inputs, breadth/depth/huge-token/match-count workload, exhaustive Python 3.13 or `.pyi` grammar coverage, non-UTF-8 encodings, multi-platform standalone builds. APIs such as `SetParseWorkLimits`, `SetTimeoutMicros`, and `SetCancellationFlag` exist but are **not tested here**. The decisive #178 grammar false success makes a deeper resource campaign uneconomic; no axis is promoted to PASS on the mere existence of an API. No root `go.mod`/`go.sum`, production API/CLI, release, or pilot was changed.
