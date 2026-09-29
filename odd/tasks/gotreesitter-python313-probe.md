# Probe pure-Go Tree-sitter Python 3.13 feasibility (#145)

Status: approved test-only feasibility work; not parser selection or census implementation.
Branch: `test/gotreesitter-python313-probe`; base: `342134b058a771ef0a0ab226ba4ce12be0ffd25a`.
Recovery mirror: Engram `odd/gotreesitter-python313-probe/tasks` in the CNSIC-bound orchestrator session (#19798); GCE-project memory writes are unavailable from this session.

## Work units

- [x] T1 — Pin `odvcencio/gotreesitter` v0.55.1 and its Python grammar lock; characterize representative Python 3.13 script and `.pyi` fixtures, especially accepted-invalid examples from tree-sitter-python #178, exact `Name.Attribute`-equivalent call shape, half-open physical byte spans and atomic failure. Route: delegated bounded writer, strict RED → GREEN → TRIANGULATE → REFACTOR. A grammar divergence is a FAIL, not an implicit scope change.
- [x] T2 — If useful after the syntax gate, characterize bounded workload, early stops, `CGO_ENABLED=0` binary, license/dependencies; document commands and PASS/FAIL/UNKNOWN. If early disqualification makes deeper axes uneconomic, mark them UNKNOWN rather than manufacturing passes. Route: same single writer, independent verification after code/test changes.
- [ ] T3 — Independent review of the frozen evidence and candidate PR preparation; merge only on the owner's candidate-specific delivery decision. Route: read-only verifier and parent orchestration.

Forecast: approximately 250–400 authored diff lines for an early disqualifier; report when a broader probe risks exceeding this. One conventional work-unit commit per completed implementation task with observed tests and documentation. Do not check off unverified outcomes.

## Verification authority and scope

GCE `AGENTS.md` and `openspec/config.yaml`: strict TDD for Go files, primary `go test ./...`, check-only gofmt; a nested module must ALSO run `GOTOOLCHAIN=go1.25.10 go test -count=1 ./...` explicitly because root tests exclude it. `go test -race ./...` only for concurrency-bearing behavior. The fixed `.pi/gentle-ai/project-policy.yaml` is absent; the native verified mode applies.

Allowed experimental surfaces: `internal/gotreesitterprobe/**`, `docs/gotreesitter-python313-feasibility.md`, and this task file. No root `go.mod`/`go.sum`, public census API/CLI, Go baseline change, runtime CPython/cgo requirement, release, CNSIC pilot, or delivery authority. A CPython 3.13 oracle may be used for development comparison, never required by the GCE executable.
