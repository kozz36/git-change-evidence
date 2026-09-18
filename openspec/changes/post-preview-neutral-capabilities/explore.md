# Explore: Post-preview neutral capabilities

## Purpose and status

This is a planning-only exploration for the new active change `post-preview-neutral-capabilities`. It uses `bootstrap-neutral-core-extraction/reconciliation-matrix.md` as its evidence source. The bootstrap change remains active and immutable at its canonical **3/12** task ledger; this change neither reconciles nor edits its tasks, artifacts, lifecycle state, or historical RED→GREEN record.

The work is explicitly post-preview and nonblocking for the planned source/dev identity `v0.1.0-preview.1`. `git-change-evidence` remains evidence-only: none of the proposed capabilities may approve, reject, gate, block, merge, deploy, release, or assign delivery authority.

## Proposal-level scope

Only the following residual capabilities are in scope; no additional bootstrap residue or multilingual census work is included.

1. **2.2 — Real-temporary-Git traversal-race proof.** Add a future proof using real temporary Git repositories for the traversal/replacement race conditions currently tested only against plain temporary directories, while retaining committed-snapshot isolation.
2. **2.3 — No-alternate-default / production-equivalent profile bridge.** Establish a future neutral/profile-adapter boundary that demonstrates production-equivalent fallback when an alternate default is not supplied, without changing neutral accounting semantics or allowing profile policy to alter neutral contract meaning.
3. **3.2 — Explicit redaction contract.** Define and test redaction across report, projection, and diagnostic surfaces for names, secrets, environment values, paths, and command output; outputs must not leak those values or imply authority.
4. **3.3 — Moving-symbolic-ref proof or API-boundary decision.** Resolve in later specification/design whether publication must accept and revalidate a symbolic Git reference, or whether the public/shell API boundary must explicitly accept only already-resolved immutable identity. This exploration makes no API choice.
5. **3.4 — Supplied forecast comparison.** Introduce a future supplied-input comparison only. It must cover valid, missing, divergent, unavailable, and threshold observations while never authoring, mutating, inferring, or gating a forecast.
6. **4.1 — One-way CNSIC profile and parity.** Add a separately bounded CNSIC profile adapter that translates CNSIC vocabulary one way into neutral inputs, rejects authority injection, preserves distinct-profile neutrality, and proves frozen-vector plus real-Git parity against the immutable predecessor oracle `a0fc7b26ff8a0e0a61baa586b32c46841611c806`.
7. **4.2 — Separately governed reversible CNSIC handoff/observation/deletion.** Define a distinct, governed handoff record covering dependency integration, reversible cutover, owner observation, and local-code deletion. It must not edit CNSIC in this change and cannot assert authority over CNSIC decisions.
8. **4.3 — Post-profile pinned-Git/distribution decision.** Record a future decision after core and first-profile stability that evaluates pinned-Git consumption and distribution options without treating preview identity as a release, package, publication, or distribution commitment.

## Constraints and non-goals

- Preserve the immutable predecessor corpus and bootstrap/OpenSpec history as evidence; do not backfill chronology or change bootstrap checkbox state.
- Preserve Functional Core / Imperative Shell and Hexagonal boundaries: typed immutable inputs in core; Git/filesystem/CLI effects in shell; configuration decoding and project vocabulary only in adapters.
- Keep neutral public APIs at the module root or the narrowly approved `census` subpackage, CLI adaptation in `cmd/git-change-evidence/`, and non-public adapters/shells in `internal/`.
- Do not select a concrete configuration-file syntax, broad distribution, packaging, signing/provenance, or release automation through this change.
- Do not implement multilingual census, CNSIC edits, external publication, tags, releases, packages, uploads, or delivery authority.
- The reconciliation matrix classifies 2.2, 2.3, 3.2, and 3.3 as partial; 3.4, 4.1, and 4.2 as absent; and 4.3 as a superseded bootstrap-era decision that remains deferred under current governance.

## Dependencies and proposed ordering

1. Specify the neutral boundaries and negative authority/redaction invariants before implementation details.
2. Address 2.2 independently as a shell-level real-Git proof; it has no required profile or forecast dependency.
3. Resolve 2.3 before 4.1 because the profile bridge/default semantics determine how a CNSIC adapter can translate its inputs without altering neutral accounting meaning.
4. Define 3.2 before any new CLI/profile diagnostics so redaction applies consistently to all later exposed surfaces.
5. Resolve the 3.3 symbolic-ref/API-boundary question in spec/design before adding any moving-ref test or publication behavior.
6. Define 3.4 as a pure supplied-input comparison after the neutral result/projection contract is identified; it must remain independent of delivery outcomes.
7. Implement and parity-test 4.1 only after the profile bridge is specified; then create 4.2 as a separately governed, reversible handoff artifact with an owner observation period and deletion criteria.
8. Consider 4.3 only after the core and first profile are stable; it remains a decision/evaluation, not a publication action.

## Alternatives and deferred API-boundary question

For 3.3, later spec/design must choose between: (a) a symbolic-ref-aware shell interface that resolves and revalidates a moving reference immediately before publication, with an explicit race proof; or (b) an immutable-identity-only publication boundary, documented as intentionally outside symbolic-ref handling. The choice must be driven by the defined consumer boundary, failure semantics, and redaction requirements; this exploration does not prejudge it.

For 2.3/4.1, the alternative of embedding CNSIC fallback rules in the neutral core is excluded because it violates the one-way profile-adapter boundary. For 4.2, a combined profile-and-handoff change is excluded so integration governance, owner observation, reversibility, and deletion remain separately reviewable. For 4.3, treating the preview identity as distribution is excluded by current governance.

## Risks and planning responses

- **Race tests may be flaky or fail to reproduce filesystem timing.** Future design should use deterministic synchronization seams around real temporary Git repositories while preserving production-path relevance.
- **Redaction may conflict with canonical evidence/provenance.** Future contracts must distinguish permitted stable identifiers from prohibited caller-controlled names, secrets, environment values, paths, and command output, and test every projection and diagnostic route.
- **A profile bridge may leak project semantics into core.** Keep fallback translation and configuration decoding in adapters; core receives already-typed neutral inputs.
- **Symbolic-ref support could widen publication responsibility.** Make the API-boundary decision first and preserve atomic/exclusive content-addressed publication behavior regardless of the choice.
- **Forecast comparison could be misread as a recommendation or gate.** Use observation-only results and explicit negative tests for authoring, mutation, inference, and delivery control.
- **CNSIC parity/cutover could be mistaken for authority or external change.** Bind only to the immutable predecessor evidence, retain one-way translation, and keep handoff/deletion owner-governed and reversible.
- **Distribution evaluation could be mistaken for release work.** Keep it post-profile, decision-only, and explicitly non-publishing.

## Testing and strict TDD

`openspec/config.yaml` sets `strict_tdd: true` for future Go source or test changes, including RED → GREEN → TRIANGULATE → REFACTOR, focused tests before production changes, focused reruns after refactoring, `go test ./...`, check-only `gofmt`, and `go test -race ./...` for concurrency-bearing behavior. This exploration changes only a planning artifact, so implementation TDD and test execution are **N/A** now; no production or test behavior is proposed as completed.

## Why this is post-preview and nonblocking

Each scoped item closes a residual semantic proof, introduces a new optional capability, or records a later governance decision rather than repairing the current preview consumer surface. The verified matrix shows existing preview-era implementation and tests for the surrounding contracts while identifying these items as partial, absent, or still deferred. Therefore no item is required to claim, publish, distribute, gate, or otherwise advance `v0.1.0-preview.1`; all must remain independently planned, evidence-only, and nonblocking until separately specified, designed, tasked, implemented, and verified.
