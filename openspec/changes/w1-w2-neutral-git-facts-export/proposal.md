# Approved contract — Export neutral committed Git facts independently

**Contract and isolated implementation are approved under S5–S6 in the ODD tracker.** Runtime capability is not yet available; implementation evidence is recorded in [apply progress](apply-progress.md). The delivered `inventory`, `result` and `report` stages and their V1 documents remain unchanged.

## Review path

Read [design](design.md), [proposed requirements](specs/w1-w2-neutral-git-facts-export/spec.md), then the [conditional implementation plan](tasks.md). The [ODD tracker](../../../odd/tasks/neutral-w1-w2-git-facts-export.md) records the documentary authorization and direction choice.

| Status | Scope |
|---|---|
| Source baseline | Main `edeffba9dcada6c5512e3f7d623d72bbe39adfc5`, tree `cda14216b9c37e8caca338a046f210cc05b0a19f`; attribution, not a consumer pin. |
| Owner-approved | The separate neutral contract, its fields/measurement/linkage/bounds and isolated implementation/tests/docs/private verification with local work-unit commits. |
| Approved, not yet delivered | `w1w2 snapshot` and `git-change-evidence.committed-snapshot/v1`; no availability or consumer-parity claim. |
| Not authorized | Dependency/version changes, push, release, distribution, pin or consumer adoption. |

## Problem and proposed outcome

The public typed snapshot contains Git metadata and numeric numstat, but Accounting Result V1 intentionally projects entries to path/category/measurement. It cannot itself export both sides, modes, object identities or side lengths. Current acquisition also does **not** obtain per-side LOC or retain original numstat tokens.

Propose a public neutral document carrying immutable revision identity, complete committed-entry/side facts and explicitly defined byte measurements. This is **not merely serializing the existing snapshot**: bounded blob measurement, token retention and complete facts validation are additional future work.

The export is independent of accounting categories, project profiles, mechanical annotations, inventory enumeration and delivery decisions. It provides technical facts, never delivery authority.

## Compatibility and boundaries

- Preserve existing CLI routing, legacy triples, root SDK and canonical V1 bytes; no silent result/report extension.
- Keep acquisition in internal shell adapters and immutable validation/assembly in the public root core.
- Preserve raw path bytes and native Git IDs; a document SHA-256 is a separate identity.
- No full blob payloads in output, worktree-content substitution, implicit fetch or inventory filtering.
- No claim that existing Report V1 provenance covers the new document, or that richer facts establish consumer parity.

## Approval and evidence still required

The owner approved the concrete contract and implementation under S5–S6. Before each Go unit, freeze its vectors, scoped layout review and calibrated forecast; implementation still requires observed tests and applicable risk-gated checks. Inventory/profile/annotation adaptation, consumer byte/schema/error parity and adoption remain separate work.

Current work is Markdown only. Functional tests, acquisition, builds and parity are **NOT RUN**. Repository rubric default is `skip` for passive documentation; future Go work requires its strict TDD gates. Historical validator unavailability is retained; no installation or retry is implied. No API availability, release, compatibility certificate or implementation approval is claimed.
