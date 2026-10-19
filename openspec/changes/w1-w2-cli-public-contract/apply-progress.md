# W1/W2 implementation boundaries

## Current local source status

The local source now wires `w1w2 inventory`, `result`, `report` and two-argument
family help. The original contract-only plan below is retained as history, not
as the current execution checklist. No distribution, release or consumer adoption is implied.

| Surface | Observed local status |
|---|---|
| Input, selected acquisition, bounded Git, composition | Focused tests and scoped independent verification completed; canonical bytes, retained binding, no-fetch/process cancellation and unavailable observations covered. |
| Root dispatch and help | RED then GREEN; 119 CLI-package test events passed without skips. Both isolated builds and 33 binary differential cases passed, allowing only one root-help entry; three new stage flows also passed. Existing seven Go files retained. |
| Remaining qualification | Root-dispatch independent checks passed. The added qualification tests have passed focused checks; their independent check and final updated full/race/format run remain pending. No formal OpenSpec-validator PASS is claimed. |

The revised implementation forecast is at most 3,900 authored changed lines including
tests/docs, allowing the bounded final qualification supplement. It supersedes the
original estimates retained below; it is not workload calibration. Do not omit
tests, compress code or waive proof to fit a forecast.

## Local scenario evidence map

This maps representative executed checks, not exhaustive OS/filesystem or workload
certification. Final supplement verification remains pending as noted above.

| Scenario | Scoped execution evidence |
|---|---|
| R1 Raw paths and safe content | Raw-name guards, descriptor reads, leaf payload tests and canonical selected inventory. |
| R1 Empty requested selection | Accessible-root empty inventory succeeds; inaccessible-root request fails. |
| R1 Unsafe or malformed selection | Closed JSON/Base64/path guards, intermediate links and nonregular leaves rejected. |
| R2 Exact success streams | Versioned constructors compared byte-for-byte, including LF, for all three stages. |
| R2 Ambiguous stdin or noncanonical input | Multiple stdin selectors rejected; unused input guards and exact canonical decoding. |
| R2 Output failure | Short/error writers return 9; diagnostic failure cannot turn an error into success. |
| R3 Change outside inventory | Real committed ranges retain changes outside the selection; worktree edits do not replace them. |
| R3 Literal refs and unavailable Git evidence | Literal metacharacter ref, immutable resolution, qualified local no-fetch trap and typed failures. |
| R3 Binding mismatch | Different valid retained policy/inventory antecedents fail with empty stdout. |
| R4 Observations are data | Zero, exceeded and unavailable ratio observations preserve code 0 and canonical output. |
| R4 No implicit discovery or effects | Explicit selection and known fixture source/index/ref checks; no completeness or global no-write certificate. |
| R5 Unsafe, unavailable or bounded acquisition | Guard/error/cancellation mappings and actual unsupported-runtime refusals before input acquisition. |
| R5 Version separation | Unsupported selection/core versions fail explicitly without decoder fallback. |
| R5 Detectable acquisition race | Before-open, post-read and captured-link replacement cases fail without partial evidence. |
| R5 Initial Linux capability boundary | Linux acquisition plus actual js/wasm execution of all three unsupported-platform refusals; Windows is compile-only. |
| R5 Individual and aggregate budget boundaries | Input/path/content/output cases plus fixed Git 16 MiB/64 MiB/10,000 below/exact/above runner responses. |
| R5 Cooperative deadline and external termination | Owned Git FIFO rendezvous: actual 30.10 s expiry returns 6; externally killed owned CLI records -9, not a promised family status. |
| R6 Legacy triple collisions | Ref/stage-help triples retain legacy measurement; discovery/census exceptions stay first. |
| R6 Differential preservation | 33 baseline binary cases allow only the single root-help entry; independent capture comparison completed. |

Decoded-path total overflow cannot be independently eligible: 128 × 4095 = 524160
bytes is below 512 KiB. Ordinary Linux symlink payloads cannot reach the 1 MiB
content cap; regular-content and mixed per-path charging are tested without inventing
such links. No privileged device fixtures were created; generic nonregular rejection
is exercised. js/wasm execution is not an all-OS or Linux syscall-filter certificate.
The fixed Git matrix uses the existing runner seam, not 64 MiB of physical Git traffic.
Calibration, hard RSS, every-write detection, atomic snapshots, all-alias provenance,
blocked-I/O wall time and universal forced-termination diagnostics are not certified.

## Original contract-only boundary review

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
