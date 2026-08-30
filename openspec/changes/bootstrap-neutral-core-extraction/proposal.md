# Proposal: Bootstrap Neutral Core Extraction

## Intent

Extract a neutral, evidence-only core without mechanically copying CNSIC. The CNSIC W1/W2 merge `a0fc7b26ff8a0e0a61baa586b32c46841611c806` and frozen vectors are immutable semantic/parity inputs; the charter governs product boundaries.

## Scope

### In Scope

- Bootstrap design for boundaries, contract ownership, vector execution, and test strategy.
- Neutral W1/W2 semantics: immutable Git evidence, separate inventory, canonical contracts, policy injection, and exclusive accounting.
- Re-derived historical W3 bounded carveout, W4 CLI/publication, and W5 forecast comparison.
- Later CNSIC profile/cutover after parity; distribution evaluation after stabilization.

### Out of Scope

- Selecting packaging, concrete source layout, license, runtime matrix, CI, release automation, or configuration syntax.
- Giving core outputs approval, rejection, blocking, merge, deployment, or release authority.

## Capabilities

### New Capabilities

- `git-snapshot-inventory`: Derive immutable changes and separate untracked inventory.
- `neutral-evidence-contracts`: Validate versioned documents and canonical bytes.
- `exclusive-policy-accounting`: Apply injected categories and evidence-only threshold observations.
- `bounded-carveout-analysis`: Compare immutable content through a pure bounded matcher.
- `evidence-cli-publication`: Provide deterministic projections and conflict-safe publication.
- `forecast-comparison`: Compare supplied forecasts without authoring forecasts or verdicts.
- `project-profile-adaptation`: Keep project vocabulary, thresholds, and compatibility outside the core.

### Modified Capabilities

None.

## Approach

Use **Branch by Abstraction** with a package-first neutral boundary. Frozen W1/W2 inputs remain the predecessor oracle until the successor passes vector and real-Git parity. Neutral specs/contracts then become the core source of truth; CNSIC policy stays profile-owned. Advance W3/W4/W5 serially, then parity-gate CNSIC cutover. Distribution remains an owner decision. Extraction gates are external, never core verdicts.

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| Neutral core | New | Contracts, accounting, carveout, reports, comparisons |
| Imperative shell | New | Git/filesystem effects, CLI, retries, publication |
| Profile edge | New | CNSIC translation and compatibility |

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Semantic or profile contamination | Medium | Require vector parity and one-way profile-to-core dependency |
| Mutable input or publication races | Medium | Use immutable objects, bounded revalidation, atomic exclusive writes |
| False delivery authority | Medium | Use neutral evidence vocabulary and negative contract tests |

## Rollback Plan

Keep both paths behind the abstraction. Parity failure stops source-of-truth transition or cutover; route consumers to the predecessor and retain evidence for diagnosis.

## Dependencies

- Immutable CNSIC merge, frozen W1/W2 vectors/manifest, and owner-approved compatibility decisions.

## Success Criteria

- [ ] Frozen vectors and real-Git cases produce equivalent canonical evidence.
- [ ] Core, shell, profile, and compatibility responsibilities remain separate.
- [ ] W3/W4/W5 behavior is re-derived and evidence-only.
- [ ] CNSIC cutover is parity-gated and reversible.
- [ ] All charter-open decisions remain unresolved.
