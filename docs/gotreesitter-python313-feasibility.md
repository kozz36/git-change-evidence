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

## T2: deferred axes

Resource worst-case controls, hard memory/time bounds, cancellation behavior, independent supply-chain inspection, and standalone build/run: **UNKNOWN pending T2**. The early grammar failure makes a full resource campaign uneconomic. No root `go.mod`/`go.sum`, production API/CLI, release, or pilot was changed.
