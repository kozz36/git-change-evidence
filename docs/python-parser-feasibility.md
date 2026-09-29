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

## F2: bounded resources, supply chain, standalone runtime

Reproduce from `internal/pythonparseprobe`:

```sh
CGO_ENABLED=0 GOTOOLCHAIN=go1.26.0 go test -c -o /tmp/gce-pythonparseprobe.test .
file /tmp/gce-pythonparseprobe.test
ldd /tmp/gce-pythonparseprobe.test  # expected: not a dynamic executable
for mode in depth breadth token files matches; do
  GCE_PROBE_WORKLOAD="$mode" GOMEMLIMIT=128MiB timeout 8s \
    /tmp/gce-pythonparseprobe.test -test.run '^TestResourceWorkload$' -test.v
done
env -i PATH=/nonexistent HOME=/nonexistent /tmp/gce-pythonparseprobe.test \
  -test.run '^TestExactCallShape$' -test.v
```

The test generates depth 200 parenthesis pairs, breadth 2,000 assignment lines, a 1 MiB string token, 64 `.pyi` files, and 10,000 call lines respectively; each mode is a separate process. The limits above are **external** 8-second timeout and soft Go GC memory limit, not parser-enforced file/match count, hard heap, cancellation or stack bounds. Observed on linux/amd64, Go 1.26.0: depth, breadth, token, files exit 0 (0.002, 0.006, 0.005, 0.001 seconds shell wall respectively); the combined test generator, probe reflection walk, and parser in the matches case exit **124 after 8.002 seconds**, without result; this measurement does not isolate parser time. At the passing test checkpoints, `HeapSys` was 7,536,640; 11,862,016; 11,993,088; and 3,866,624 bytes respectively, `StackSys` was 851,968; 655,360; 557,056; 327,680 bytes. These are checkpoint runtime counters, **not** peak process RSS or guaranteed worst-case maxima. `ulimit -v 262144` could not be used: Go runtime failed to reserve page-summary virtual address space before test execution (exit 2), so a hard process memory limit remains **UNKNOWN**. The combined 10,000-match probe is an observed time-bound **FAIL** at this limit (not a parser-only benchmark); bounded time/heap/stack/cancellation for arbitrary legal input remain **UNKNOWN**.

The dependency archive contains `parser2/{ast.go,lex.go,parser.go,stmt.go}` under the pinned module zip; no separate downloaded grammar artifact appears in that package. The README claims Python 3.14 grammar, but independent whole-grammar certification is **UNKNOWN**. `go list -m all` in the nested module listed only the probe module and `github.com/tamnd/gopapy/v2 v2.0.0`; `go list -deps` showed only its `parser2` package beyond the standard library and the probe. The candidate's bundled `LICENSE` says MIT, copyright 2026 Duc-Tam Nguyen (SHA-256 `17aabccf56cbcdbe1ceae5416effc3cf9e44c7b0e5339281751a53d44fbae41f`); module zip SHA-256 `5724ca80af87a3fb7c5f58607660ad8cfda67fd77a9b038c070382f0c834a2f7`. This is source attribution and dependency inspection, not legal clearance.

`CGO_ENABLED=0` produced a statically linked ELF (`ldd`: `not a dynamic executable`); the exact-call tests passed under `env -i PATH=/nonexistent HOME=/nonexistent`, with no Python executable available on PATH. This demonstrates that isolated binary on this host, not Go 1.25 compatibility or a reproducible cross-platform release build. **Fail-closed conclusion:** the pinned candidate does not currently satisfy the required span and Go-baseline axes; resource worst-case controls are unproved and one bounded match run timed out. Stop before adoption; any alternative or changed baseline needs a separate owner decision.
