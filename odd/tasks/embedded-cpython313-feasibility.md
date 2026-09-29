# Embedded CPython 3.13 feasibility for the Python census (#147)

Status: issue approved; evidence-only feasibility, not parser adoption, census implementation, or a CNSIC pilot.
Branch: `test/embedded-cpython313-feasibility`; base: `5db09940d8e4a7863eff784d66f532d79c09f9d8` (`origin/main`).
Recovery mirror: Engram `odd/embedded-cpython313-feasibility/tasks` in the CNSIC-bound orchestrator session; GCE-project memory writes are unavailable from this session.

## Work units

- [x] T1 — Pin a real CPython 3.13 revision and examine whether `goccy/python-wasm` can reproducibly generate the standalone Go 1.25-compatible runtime, or capture the first concrete blocker. **Observed**: CPython v3.13.13 at `01104ce1`, pinned wasmify image `sha256:fe54e565`, build exit 2 at missing `Tools/wasm/wasi`, before compilation. Report: `docs/embedded-cpython313-feasibility.md`; work-unit commit and root tests recorded below.
- [x] T2 — **Not applicable (early blocker)**: T1 generated no 3.13 candidate, so AST/compile, spans, atomicity, binary footprint, concurrency and capabilities were not tested. Each axis is UNKNOWN in the report, not a pass. No Go code or nested module was created.
- [x] T3 — Independent read-only audit found two report defects; correction `f444648` passed focused re-audit and root checks. Owner separately authorized publication; branch and evidence-only PR [#148](https://github.com/kozz36/git-change-evidence/pull/148) were published with approved issue #147 and `type:chore`. PR policy, Go checks and CodeQL passed on the published candidate. **Merge and parser adoption remain unauthorized.**

## Scope and checks

Allowed experimental surfaces: `internal/embeddedcpythonprobe/**` (isolated nested module, only if candidate generated), `docs/embedded-cpython313-feasibility.md`, and this task file. External toolchain scratch lives outside this repository and is never committed. No root `go.mod`/`go.sum`, production census API/CLI, copied CPython runtime, new release, or downstream integration. A 3.13 oracle and Docker toolchain may be used to generate evidence, not as runtime dependencies of GCE.

Go changes use RED → GREEN → TRIANGULATE → REFACTOR. Run focused nested-module tests explicitly, plus root `GOTOOLCHAIN=go1.25.10 go test ./...` and check-only formatting; root tests alone exclude nested modules. Record full commands, status, toolchain/OS, size/resource numbers and any skipped checks. An early build failure is a valid negative result, not permission to silently switch to 3.14 or relax the contract.

## Observed checks and commits

- `GOTOOLCHAIN=go1.25.10 go test ./...`: PASS (six root packages). `gofmt -l` on tracked Go files: empty; `git diff --check`: PASS. No Go files changed, and no nested probe module exists, so no nested tests were run.
- First pinned build attempt: exit 2, blocked at CPython 3.13's missing `Tools/wasm/wasi`; exact command/log digest and all unknown axes in the report. No static binary was built. The upstream reported 3.14 footprint is not a GCE measurement.
- T1 evidence commit: `83bf51f792fe27d668bfa5e9c7ff5d55f6f934fd` (`test(python): record embedded 3.13 build blocker`).
- T3 independent audit of frozen report found a false claim about the WASI triple and incomplete repro checkout commands; corrected both in `f444648`. Focused re-audit: PASS; full root suite rerun with `-count=1`: PASS. PR #148 was published after explicit owner decision; CI PR policy, Go checks, CodeQL analyses: PASS. No merge or adoption decision has been granted.

Issue: https://github.com/kozz36/git-change-evidence/issues/147 (`status:approved`, verified on target). The evidence report's `feasible | infeasible | undetermined` assessment never authorizes adoption or delivery.
