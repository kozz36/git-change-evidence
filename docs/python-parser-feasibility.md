# Gopapy parser feasibility observations (#143)

Status: **test-only, fail-closed; not a parser selection or Python census implementation**. This report is the experimental evidence for the design in [python-census-design.md](python-census-design.md). No production API/CLI or root module dependency is changed. The nested `internal/pythonparseprobe` Go module is deliberately outside root `go test ./...`; always run its tests explicitly.

## Pinned candidate and reproduction

- Candidate: `github.com/tamnd/gopapy/v2` tag `v2.0.0`, source revision `8432e27c5ae6705a4f91c67d6132a9ca46701b34`, module sum `h1:JN1e4dAHEPbsOZBC6RVKfbWYcWjCW5tvg17TfhZWkbA=`, go.mod sum `h1:3Q3MQ/O9WhNsOI+FSjW/RhgzUyUs5ZQSCvW6g0jz1tE=`. Its `parser2` declares a Python 3.14 target upstream; this probe has not certified whole-grammar coverage.
- `internal/pythonparseprobe/go.mod` pins the dependency inside an isolated module. It declares Go **1.26**, versus the GCE baseline **1.25.10**. This is a baseline compatibility **FAIL**, not a request to upgrade GCE.
- From repository root: `cd internal/pythonparseprobe && GOTOOLCHAIN=go1.26.0 go test -v ./...`; baseline check: `cd internal/pythonparseprobe && GOTOOLCHAIN=go1.25.10 go test ./...`. Run root `go test ./...` separately; it does **not** test the probe.

## F1: syntax and byte-position observations

| Axis | Observation | Status |
| --- | --- | --- |
| Exact syntax predicate | Direct, shadowed, nested inner call recognized; attribute-only, chained receiver and indexed function excluded in table tests. No name resolution. | PASS for listed fixtures only |
| `.py` and `.pyi` | Async call, overload-decorated annotated stub, PEP 695 alias and `match` fixture parsed in the isolated tests. | PASS for listed fixtures; broad grammar UNKNOWN |
| Malformed atomic input | A valid candidate followed by malformed stub returns `nil` matches and an error. A future-looking `defer` form is rejected in one fixture. | PASS for listed fixtures; unsupported-grammar rejection generally UNKNOWN |
| Physical byte spans | On `π = 'é'\r\nos.open(os.open('猫'))\r\n`, AST call-position offsets observed **18, 26**, while actual physical call starts are **11, 19**. The AST `Call` exposes only `P` (line/column), no end offset. The probe does not invent a half-open function/call span or a fragment hash. | **FAIL**: required byte-span evidence not demonstrated |

These are characterization tests, not a production extractor. The reflection walk in the probe has no resource-depth guarantee. A green probe test means its observations match this pinned revision, **not** that every feasibility axis passed. Do not adopt this candidate on this evidence.

F2 resource, supply-chain and standalone measurements: pending.
