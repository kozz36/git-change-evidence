# Public forecast ingestion PR5

Issue: #162. Work unit: W1. Status: in progress.

- [ ] W1 — Transport immutable forecast ingestion and verify exact bytes, privacy and Go behavior.

## Scope

- `forecast_v1.go`
- `forecast_v1_test.go`
- `openspec/changes/bootstrap-neutral-core-extraction/design.md`
- `openspec/changes/bootstrap-neutral-core-extraction/specs/forecast-comparison/spec.md`
- `odd/tasks/public-forecast-ingestion-pr5.md` — public tracking only, at most 60 lines.

Preserve the four transported technical files exactly at the admitted source boundary and preserve inherited public prefixes. Do not export the historical execution record. No additional redactions, redesign, comparison behavior or unrelated source changes.

## Checks and evidence

Frozen transport has no meaningful behavioral RED: retain the original tests and run fresh verification, without manufacturing test-first evidence.
Exact-byte transport, documentation patch application and public-base footprint checks: passed. Official Go 1.25.10 directed, focused, full and race checks: passed. Check-only formatting, whitespace, current-slice privacy and bounded preservation: passed.
Evidence is scoped to the current slice and observed baselines; preparation capture is qualified, and historical unknowns and unaudited host/cache state remain unchanged.
Work-unit commit: `637df1e5b4572bb641efd539cb8208f4d1ea666e` (`feat: add immutable forecast ingestion`). Recorded before the native review freeze; W1 remains open for the later gates.

## Later gates

Native RDD is enabled; candidate consent and approval are pending. A checkbox is not review authority.
Current matching-HEAD CI and parent-owned publication and merge remain pending. No source changes are planned after native freeze.
