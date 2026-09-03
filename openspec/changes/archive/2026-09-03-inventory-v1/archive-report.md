# Archive Report: Inventory V1

## Status
**PASS — archive readiness under ordinary `disabled/unmanaged` policy.**
Archive readiness is resolved from merged Git/OpenSpec evidence under ordinary policy, not native structured state.

## Artifacts Read
- `openspec/changes/inventory-v1/{proposal,tasks,apply-progress,verify-report,sync-report}.md`
- Active specs: `openspec/changes/inventory-v1/specs/{inventory-v1-contract,inventory-v1-content-identity}/spec.md`
- Canonical specs: `openspec/specs/{inventory-v1-contract,inventory-v1-content-identity}/spec.md`

## Canonical Synchronization
- Both canonical domains are synced from their active sources.
- `inventory-v1-contract` ADDED: `Neutral root-package Inventory V1 API`; `Closed Inventory V1 wire contract`; `Policy V1 antecedent binding`; `Strict validating decode`; `Verified-empty scope and acquisition separation`; `Additive compatibility and historical-authority isolation`.
- `inventory-v1-content-identity` ADDED: `System-derived exact-content identity`; `Raw path-byte fidelity and validity`; `Deterministic exact-path set and ordering`; `Defensive ownership of Inventory V1 entry material`.
- Byte identity: PASS — `cmp -s` matched both active-source/canonical-spec pairs.
- Same-domain active census: one source for each domain; no conflict warning.

## Closure Evidence
- Task census: PASS — zero unchecked tasks; 2A.1–2A.12 and 2B.1–2B.12 are checked.
- `verify-report.md`: PASS; `sync-report.md`: PASS for ordinary `disabled/unmanaged` synchronization.
- Stale-checkbox reconciliation: the historical 2A-close unchecked-2B statement is time-scoped chronology; proof is the current 2B-complete section and zero-unchecked-task census.
- Pre-report closure budget: 344/400; projected final total: 373/400 changed lines.
- Issue #46 is approved, and the user approved the canonical dated archive move to `openspec/changes/archive/2026-09-03-inventory-v1/`.
- Native #4040 attempt remains unresolved; this report claims no native settle, verification, or approval.
- Rollback boundary: before or after the parent-controlled move, reverse only the archive-directory move; do not alter product code, canonical specs, or native state.

## Key Learnings
- Ordinary archive readiness can rely on merged Git/OpenSpec evidence while native state is unresolved.
