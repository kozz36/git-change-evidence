# Proposal: Preserve bootstrap evidence and separate post-preview capabilities

## Intent and resolved product decision

**Preservar y separar**: preserve `bootstrap-neutral-core-extraction` active and immutable at exactly **3/12**, and plan its eight residual capabilities separately under `post-preview-neutral-capabilities`.

This is **proposal-level planning only**. All specifications, design, tasks, implementation, and verification are future post-preview work, nonblocking for the planned source/dev identity `v0.1.0-preview.1`. That identity is not a published tag, release, package, asset, distribution channel, or visibility change. There is no unresolved product decision and no further product interview is required; later technical choices remain explicitly bounded below.

`git-change-evidence` remains evidence-only. It may report technical states, measurements, and observations, but must never approve, reject, gate, block, merge, deploy, release, or assign delivery authority.

## Motivation and evidence

Maintainers and future neutral API/profile consumers need a clear distinction between present behavior, historical completion records, and future capabilities. Conflating those would manufacture bootstrap chronology or make optional post-preview work appear necessary for the source/dev preview.

The direct inputs are [the canonical exploration](explore.md), [the verified bootstrap reconciliation matrix](../bootstrap-neutral-core-extraction/reconciliation-matrix.md), root `AGENTS.md`, and `openspec/config.yaml`.

The matrix classifies 2.2, 2.3, 3.2, and 3.3 as **partial**; 3.4, 4.1, and 4.2 as **absent**; and 4.3 as **superseded-decision**, with the distribution decision still deferred. It reports current tests and later implementation evidence, not a recovered historical RED→GREEN sequence. In particular, even proven-present 3.1 remains outside this proposal: its present behavior does not authorize editing its historical checkbox.

This proposal consumes those recorded observations without claiming a new source audit, test run, implementation, or native lifecycle transition.

## Goals and scope boundaries

- Make exactly the eight residual capabilities below independently understandable and suitable for later specification.
- Preserve neutral contract meaning, committed evidence identity, and one-way adapter boundaries.
- Describe future success evidence without claiming any residual capability is complete.
- Separate CNSIC handoff governance and distribution evaluation from neutral implementation and preview identity.

The only artifact created now is `openspec/changes/post-preview-neutral-capabilities/proposal.md`. The exploration and every bootstrap artifact, including the reconciliation matrix, tasks, apply-progress, and historical chronology, remain unchanged. Bootstrap remains active at exactly 3/12; no reconciliation, archival, or closure is proposed.

## Capabilities and proposal-level requirements

The identifiers below are traceability references to bootstrap residue, not new task checkboxes or completion claims.

### 2.2 — Real-Git traversal-race proof

Plan traversal/replacement-race proofs in real temporary Git repositories. The matrix distinguishes existing plain-directory race tests from the separate Git-backed committed-snapshot isolation proof; neither should be represented as the missing combined proof.

Future evidence should exercise unsafe traversal and replacement races through production-relevant paths, return no records for unsafe acquisition, and retain committed-snapshot isolation despite worktree/index changes. Deterministic synchronization should avoid timing-only assertions.

### 2.3 — No-alternate-default / production-equivalent profile bridge

Plan an adapter bridge that supplies production-equivalent fallback when no alternate default is supplied. Current neutral accounting requires an explicit named default; this proposal does not claim fallback is already implemented or relax neutral accounting semantics.

Future evidence should demonstrate omitted-alternate fallback and supplied-alternate behavior, invalid/overlapping classification handling, and unchanged neutral totals and threshold observations. Project defaults and configuration decoding belong in adapters, never as CNSIC-specific core rules.

### 3.2 — Explicit redaction contract and tests

Plan an explicit contract for redaction across report, projection, and diagnostic surfaces covering names, secrets, environment values, paths, and command output, without authority implications. Generic errors or omission from human output alone do not establish comprehensive redaction; the matrix notes canonical report bytes retain caller `Subject`.

Future specification must distinguish permitted stable evidence identifiers from prohibited caller-controlled content and address canonical-byte/provenance compatibility explicitly. Future tests should inject prohibited values into each exposed route and demonstrate no leakage, including failures, without silently changing existing versioned contract meaning.

### 3.3 — Moving-ref proof or immutable-identity API-boundary decision

Plan a technical decision before extending publication: either accept and revalidate symbolic Git references with a moving-ref race proof, or explicitly define an immutable-identity-only public/shell boundary that excludes symbolic-ref handling. Both alternatives satisfy the resolved product scope; this proposal selects neither API design.

Future evidence must match the chosen boundary: deterministic moving-ref behavior and revalidation for symbolic-ref support, or an explicit contract and tests demonstrating immutable-identity-only behavior. Preserve atomic, exclusive, content-addressed local publication and existing interruption/conflict protections. Do not describe current publication as accepting or re-resolving symbolic refs.

### 3.4 — Supplied forecast comparison

Plan comparison of supplied forecast inputs only, covering valid, missing, divergent, unavailable, and threshold observations. No current neutral forecast API is claimed; predecessor fixtures are not such an implementation.

Future evidence should establish deterministic observation-only results, caller input ownership, and negative cases proving forecasts are never authored, mutated, inferred, or used to gate delivery. Comparison must remain independent of delivery outcomes.

### 4.1 — One-way CNSIC profile and parity

Plan a bounded profile adapter translating CNSIC vocabulary and defaults one way into typed neutral inputs. It must reject authority injection and preserve distinct-profile neutrality rather than embedding CNSIC policy in the core.

Future evidence should include adapter-boundary rejection tests, distinct-profile proofs, and frozen-vector plus real-Git parity against immutable predecessor `a0fc7b26ff8a0e0a61baa586b32c46841611c806`. Preserve that corpus as evidence. Existing W1/W2 compatibility fixtures do not establish a standalone profile or its parity suite.

### 4.2 — Separately governed reversible CNSIC handoff/observation/deletion

Plan a separately governed CNSIC handoff record covering dependency integration, reversible cutover, an owner observation period, and local-code deletion criteria. Its preparation follows profile parity; actual CNSIC integration, observation, and deletion remain outside this change and under CNSIC owner control.

Future documentary evidence should identify ownership, integration prerequisites, rollback steps, observation evidence, and deletion conditions. The record cannot authorize CNSIC changes, claim observations not performed, or turn neutral measurements into delivery decisions. No CNSIC repository is edited here.

### 4.3 — Post-profile pinned-Git/distribution decision

Plan a decision record evaluating pinned-Git consumption and distribution alternatives only after core and first-profile stability is evidenced. Record tradeoffs and retained deferrals without assuming a publication commitment.

Future success is a traceable evaluation under then-current governance, not a tag, release, package, upload, or visibility change. The planned preview identity does not satisfy this capability. Broad distribution, packaging, signing/provenance, and release automation remain deferred rather than selected by this proposal.

## Non-goals and preserved behavior

No additional bootstrap residue, multilingual census, production changes, tests, configuration changes, general documentation changes, or external actions belong to this proposal phase. Do not choose configuration-file syntax, package formats, signing, release automation, or distribution channels. Do not rewrite literal fixtures, local-path records, historical OpenSpec/CNSIC evidence, or bootstrap completion chronology.

Future work must preserve Functional Core / Imperative Shell with Hexagonal adapters: immutable typed neutral core inputs; Git/filesystem/CLI effects in shell; configuration decoding and project vocabulary in adapters. Public APIs remain at the module root except the approved `census` subpackage; CLI adaptation remains in `cmd/git-change-evidence/`, with non-public shell/adapters in `internal/` and co-located Go tests.

The CLI and machine-readable contracts remain the primary consumer surface. Preserve canonical `gce census go-ast`, the pre-v1 `git-change-evidence census-go-ast` compatibility command, the root public Go API, and the public `census` API. No change to their established meaning or behavior is proposed here.

## Dependencies and order for future phases

1. Specify neutral invariants, authority exclusions, and redaction requirements before new exposed diagnostics or profile surfaces.
2. Develop 2.2 independently of profile and forecast work.
3. Resolve 2.3 bridge/default semantics before implementing 4.1; keep translation one-way.
4. Resolve the bounded 3.3 API choice before moving-ref tests or publication changes, consistent with 3.2 redaction and provenance requirements.
5. Specify 3.4 against the identified neutral result/projection contract, independently of delivery decisions.
6. Implement and establish parity for 4.1 before preparing the separately governed 4.2 handoff; do not combine parity with external cutover/deletion.
7. Evaluate 4.3 only after core and first-profile stability. It does not authorize distribution.

These are planning dependencies, not a task ledger, implementation schedule, or permission to proceed into later phases. No delivery chain or review-budget exception is selected.

## Affected areas

| Area | Potential future impact, subject to later design |
|---|---|
| Git/inventory shell | Real-Git traversal-race evidence for 2.2. |
| Neutral accounting and profile boundary | Typed fallback bridge and unchanged accounting meaning for 2.3. |
| Report/projection and CLI diagnostics | Explicit redaction and compatibility handling for 3.2. |
| Publication shell/public boundary | Bounded identity decision and matching proof for 3.3. |
| Neutral comparison contracts | Supplied-input forecast observations for 3.4. |
| Non-public profile adapters and parity tests | One-way CNSIC translation for 4.1. |
| Separate handoff and decision records | Owner-governed 4.2 and post-stability 4.3; no external execution. |

These are impact hypotheses from the supplied artifacts, not commitments to particular files, signatures, or new packages.

## Risks and responses

| Risk | Planning response |
|---|---|
| Race proofs become flaky or detached from production paths | Specify deterministic synchronization and real-Git fixtures while retaining snapshot isolation. |
| Redaction undermines canonical identity or compatibility | Define permitted identifiers and explicit compatibility/versioning treatment before changing serialization. |
| CNSIC fallback contaminates neutral accounting | Require adapter-only translation and distinct-profile neutrality evidence. |
| Symbolic-ref support broadens publication responsibility unnecessarily | Decide the consumer API boundary first; preserve local atomic/exclusive publication invariants. |
| Forecast or parity observations imply authority | Require negative authority tests and observation-only contracts; keep owner governance separate. |
| Cutover/deletion is premature or irreversible | Separate the handoff, observation period, rollback, and owner-controlled deletion conditions. |
| Preview preparation is mistaken for distribution or historical completion | Repeat source/dev-only status; retain the immutable 3/12 ledger and explicit deferred distribution decision. |

## Rollback

This proposal changes no runtime state. It can be revised or withdrawn through a later explicit planning decision without editing bootstrap history or the exploration. No runtime rollback, external rollback, or deletion is executed now.

Future designs must identify reversibility before implementation. In particular, the separate CNSIC handoff must retain a rollback path through cutover and observation and must not assume local-code deletion is automatically reversible or authorized. Distribution evaluation remains non-publishing.

## Success criteria and future evidence

Proposal success means the eight and only eight residual capabilities are traceable to the exploration/matrix, bootstrap remains active and immutable at 3/12, and no new capability or chronology is claimed complete. The confirmed product decision is preserved without reopening discovery.

Future phase success requires the capability-specific evidence above, plus explicit negative authority and compatibility coverage. Specifications, design decisions, task completion, implementation results, verification reports, and owner observations must be recorded prospectively in their proper phase; none is supplied by this proposal. Existing matrix test results are historical input evidence, not tests rerun here or proof of absent capabilities.

## Testing and strict TDD

**Planning-only strict TDD: N/A.** No Go source or test changes are made, and no tests are run in this phase.

Future Go work follows `openspec/config.yaml`: RED → GREEN → TRIANGULATE → REFACTOR, focused tests before production changes, focused reruns after refactoring, `go test ./...`, the configured check-only formatting command, and `go test -race ./...` for concurrency-bearing behavior. Future evidence must preserve actual execution chronology rather than backfill bootstrap records.
