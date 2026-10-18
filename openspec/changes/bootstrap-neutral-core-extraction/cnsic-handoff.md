# Draft CNSIC handoff record

## Status and scope

This GCE-side record prepares a bounded future consumer-integration handoff.
It does not create or adopt a consumer integration SDD. A separately planned
and human-approved consumer SDD is a prerequisite to future adoption, not to
preparing this document. No consumer integration, pilot, cutover, local-code
deletion, or delivery decision is performed or authorized here.

Historical completion of core/profile work and documentary handoff/evaluation
work does not complete any future consumer stage below. Technical evidence is
not human transition approval. GCE reports evidence only and assigns no delivery
authority.

## Available surfaces and bounded evidence

The neutral versioned contracts, public root Go API, and CLI are technical
surfaces available for later integration planning. The existing
`internal/cnsicprofile.DecodePolicyV1` is an internal adapter, not a public
consumer import API or a consumer-ready integration boundary. This handoff
selects no public profile exposure, bridge, implementation, or configuration syntax.

The existing standalone profile work covers policy admission, authority-field
rejection, immutability, distinct-profile neutrality, frozen W2 checks, and
SHA-1/SHA-256 Git commit/blob bindings with dirty/index isolation.
Real-Git parity is **W2-bounded change-only with empty inventory**.
It is not full-tree inventory, general CNSIC integration parity, a fresh
predecessor execution, historical cryptographic proof, or cutover approval.
These are bounded descriptions of prior technical work, not new verification
results from preparing this handoff.

## Future consumer-owned stages — none completed here

### 1. Dependency integration

- [ ] Obtain human approval for a separate CNSIC integration SDD and its scope.
- [ ] Plan the consumer-owned adapter/integration boundary around neutral inputs;
  record API/CLI contract mappings and unsupported or unavailable outcomes.
- [ ] Record the approved dependency mechanism and version pin after a separate
  human decision; package/distribution choices remain UNDECIDED here.
- [ ] Produce consumer-specific frozen-vector and applicable real-Git parity
  artifacts, including inventory scope, object identities, and evidence limits.

### 2. Reversible cutover

- [ ] Document and verify a retained predecessor route and reversal procedure;
  select the predecessor on failed consumer parity or unavailable evidence.
- [ ] Record consumer-owner approval before any proposed source-of-truth change;
  parity evidence alone does not approve transition, shipping, or delivery.
- [ ] Retain a reversible adapter until owner-approved deprecation; capture
  reversal verification artifacts and unresolved failures in the consumer SDD.

No implemented fallback, concurrent dual running, or production gate is asserted.

### 3. Owner observation

- [ ] Obtain explicit human decisions for every UNDECIDED item below before use.
- [ ] Capture observations against the approved consumer acceptance criteria,
  including unavailable evidence, parity failures, and any reversal observations.
- [ ] Preserve an owner decision record accepting or declining the proposed
  deprecation and rollback-window closure. Time passing is not acceptance.

### 4. Conditional local-code deletion

- [ ] In the future consumer SDD, inventory proposed local-code removals and
  retained compatibility/reversal responsibilities; no filenames are inferred here.
- [ ] Require accepted consumer parity, completed owner observation, explicit
  deprecation approval, and approved rollback-window closure before deletion.
- [ ] Record the removal scope, prerequisite evidence, owner decisions, and
  post-removal verification artifacts before carrying out approved removals.

## Decisions requiring future consumer human approval

| Decision | Status |
|---|---|
| Observation owner | UNDECIDED |
| Observation duration | UNDECIDED |
| Consumer acceptance criteria | UNDECIDED |
| Deprecation approval | UNDECIDED |
| Rollback-window closure criteria | UNDECIDED |
| Access and ownership responsibilities | UNDECIDED |
| Consumption boundary and integration mechanism | UNDECIDED |
| Concrete revision pin and change handling | UNDECIDED |
| Actual distribution | UNDECIDED |

No names, deadlines, numeric defaults, or automatic approvals are assigned.

## Next owner action

The future consumer owner should propose the separate integration SDD, its
boundary, evidence plan, and decision owners for human approval. This record is
input to that proposal, not its creation or adoption.
The [Git-SHA provenance evaluation](distribution-evaluation.md#why-evaluate-git-sha-provenance)
records the selected evaluation direction, not a concrete consumption pin,
adoption mechanism, or distribution implementation. Consumer changes and
local-code deletion remain separately approved future work.
