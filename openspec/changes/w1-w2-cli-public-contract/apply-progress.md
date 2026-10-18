# W1/W2 implementation boundaries

Local implementation is authorized; this contract work unit introduces no executable CLI behavior. The selected design remains the target, not a claim of implemented support or calibration. See [design](design.md), [requirements](specs/w1-w2-cli/spec.md) and [remaining work](tasks.md).

## Package-boundary review

The implementation keeps Functional Core / Imperative Shell and leaves public root Go contracts and census unchanged. Existing `internal/git.Acquire` behavior is preserved; contextual bounded acquisition uses a new internal seam.

| Surface | Reviewed boundary |
|---|---|
| Existing CLI | Seven Go files: four production and three tests. Conservatively review the configured six-file threshold counting tests; add no CLI files. Keep main/help dispatch edits thin in existing files. |
| New shell leaves | `internal/w1w2/contract`, `input`, `acquire`, `cli`; share typed request/profile/path validation before independent input/filesystem units. No public package or configuration syntax. |
| Internal Git | New bounded/contextual acquisition and co-located tests; no change to legacy unbounded acquisition semantics. |
| Growth checks | Plan at most eight changed Go files per package, fifteen Go files per internal leaf and 160 authored changed lines per file; record and review actual crossings before expansion. Do not compress code, comments or tests to fit. |
| Dependencies | Existing Go 1.25.10 and `golang.org/x/sys v0.38.0`; no new dependency, installer or distribution change. |

## Planned local work units

| Unit | Behavior and verification |
|---|---|
| Contract | Versioned CLI specification and explicit package/verification boundaries; no executable behavior. |
| Shared protocol | `internal/w1w2/contract`: fixed profile, typed `Request`/`Failure`, `ValidatePath`/`ValidatePaths`; raw-name admission without authorization or filesystem-confinement claims. Tests pin selected constants, raw identities, exact metadata names and path/count boundaries; this unit does not wire CLI operations. |
| Input | Closed selection/flags/one-stdin/bounded documents; focused public-boundary tests. |
| Regular acquisition | Descriptor-relative roots/files/types/content/deadline guards; synthetic fixtures and independent security verification. |
| Link acquisition | Leaf/dangling-link payloads, hardlink accounting and detectable races; separate security verification. |
| Git | Bounded contextual full immutable range; real synthetic Git and runner fault tests, legacy acquisition preserved. |
| Composition | Exact inventory/result/report constructors, antecedents, diagnostics and buffered output; success/failure/no-write tests. |
| CLI | Dispatch/help, all 19 scenarios, clean baseline differential, full/race/format checks and accurate support documentation. |

Initial forecast: approximately 1,200–1,900 authored changed lines including tests/docs, not calibrated. Prepare a feature-branch review chain with cohesive local units; no push, PR/tracker publication or merge is authorized. About 400 lines per unit is a review heuristic, not permission to omit tests or force size-only restructuring.

## Verification boundaries

Use Go 1.25.10, co-located tests and RED → GREEN → TRIANGULATE → REFACTOR. Keep fixtures, caches and evidence in private excluded per-run state, never production data or credentials. Source edits and verification writes have distinct bounded scopes; real-repository Git belongs to the coordinator. Run the configured check-only formatting command in a clean candidate snapshot so it does not inspect runtime caches, dependencies or generated fixtures. Committed tests do not read Git history or a fixed commit; baseline differential state stays private. A failed/unavailable check remains failed/unavailable; evidence never decides delivery.
