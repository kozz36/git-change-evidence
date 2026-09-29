# Port embedded CPython 3.13 build recipe (#149)

Status: issue approved; evidence-only build-port, not parser adoption or Python census implementation.
Branch: `test/embedded-cpython313-build-port`; base `0385c84d5832b3dbf832e1d24102943d4df12165` (`origin/main`).
Recovery mirror: Engram `odd/embedded-cpython313-build-port/tasks` in the CNSIC-bound orchestrator session; GCE-project memory writes unavailable here.

## Work units

- [x] P1 — Pinned CPython 3.13.13, upstream recipe and image; a small patch fixes the helper/config paths, aligns the host build and turns off host sockets/subprocess. Full bounded WASM → Go generation exited **0** in 5m58s; optimized WASM/Go bundle, hashes, warnings and the **not-yet-running** boundary appear in `docs/embedded-cpython313-build-port.md`. Clean pre-run patch applies to a fresh checkout. Root tests and commit recorded below.
- [x] P2 — Assembled generated bridge/bundle with matching 3.13 stdlib in scratch, built a static `CGO_ENABLED=0` Go 1.25.10 smoke executable and ran it in networkless/read-only BusyBox without Python. Four focused tests PASS for interpreter identity, exact shape, 3.13 stub syntax, AST-only invalid `return` rejected by compile, malformed source and nonexecution of scanned source. A widened receiver predicate mutation died. Binary/startup/heap sample recorded; physical spans, multi-file atomicity, adversarial bounds and supply-chain certification remain UNKNOWN. Commit and checks below.
- [ ] P3 — Independently verify immutable pins, adaptations, log provenance, current report and exact PASS/FAIL/UNKNOWN scope; fix concrete report defects and rerun checks. Parent alone may offer push/PR/merge after a separate owner decision. No adoption or CNSIC pilot from a successful candidate.

## Allowed surfaces and budget

`internal/embeddedcpythonprobe/**`, `docs/embedded-cpython313-build-port.md`, this task file. Scratch upstream clone/build outputs remain outside the GCE repo; no generated CPython runtime or root module edits. Stop to re-scope if the port exceeds a small reviewable patch or requires broad upstream patches. GCE architecture, Go 1.25.10 and strict Python validity/byte-span/atomicity requirements remain invariant.

Go TDD: RED → GREEN → TRIANGULATE → REFACTOR. Root `GOTOOLCHAIN=go1.25.10 go test ./...`, check-only gofmt and `git diff --check`; nested module tests explicitly if one is created. Record exact build commands, pinned hashes, exit codes, limits, and every skipped check. One conventional commit for each completed implementation work unit. No push/PR/merge without owner decision.

## Observed checks / commits

- P1: pinned Docker full generation exit 0, four CPUs/8 GiB, 5m58s, log SHA-256 `ab6b6ddcc1c9ba1f00765e9dc416b4106016f8db7054563eae8dc4124c2f4a32`. `git apply --unidiff-zero --check` and application on fresh pinned recipe, pre-run config hash, `bash -n`, `git diff --check` observed PASS. This is generator evidence only, not a Go runtime test.
- P1 work-unit commit: `90b7c36f21f303543874c63cd7989eccc9167d17` (`test(python): port pinned 3.13 wasm build recipe`).
- P2: RED undefined `runtimeVersion`/`inspect` → GREEN all four tests on generated 3.13 assets with Go1.25.10, `CGO_ENABLED=0`. Binary 38,237,673 bytes; BusyBox read-only/networkless runtime version 3.13.13, direct calls 1. One widened-attribute mutation → targeted test FAIL (`2` vs `1`), restored test PASS. Exact commands/digests and residual UNKNOWN in report. P2 work-unit commit: pending readback.
- P3 independent audit and delivery: pending.

Issue: https://github.com/kozz36/git-change-evidence/issues/149 (`status:approved`, read back). Prior indeterminate report: `docs/embedded-cpython313-feasibility.md` (#147/#148).
