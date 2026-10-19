# Neutral snapshot — implementation evidence

Contract and isolated implementation/local commits are approved under S5–S6. No push, release, pin or adoption.
This ledger distinguishes approved scope, frozen expectations and actually executed checks.

## Layout and unit review

Baseline: edeffba9dcada6c5512e3f7d623d72bbe39adfc5 / tree cda14216b9c37e8caca338a046f210cc05b0a19f.
Parent package-boundary review: 31 root production Go files already exceed the configured 12; the first new API increases that to 32.
A neutral public API belongs at root per AGENTS; do not hide it in an internal/profile package or add a new public subpackage.
One new production file is below the 8 changed-production-files/package limit and is planned below 160 lines. Existing production files remain unchanged.
The inventory builder supplies only embedded len/hash/path work, not the required standalone line/NUL measurement; no neutral reusable primitive/name collision was found.
Calibration: its small pure derivation and existing external-API test style support a first-unit forecast of production40 + tests110 + mechanical0 + artifacts16 =166 authored diff lines. Counts are estimates, not a floor or reason to compress.
First Go surfaces only: snapshot_document_v1_content.go and snapshot_document_v1_content_test.go. Later units need their own review/calibration.

## Frozen first-unit vectors — NGE-4

Proposed Go identifiers are implementation-owned: SnapshotContentV1 and MeasureSnapshotContentV1.
Each public call must return byte length, NUL flag and LF-logical-line count together and leave supplied bytes unchanged.

| Input bytes | Expected byte_length / contains_nul / logical_lines |
|---|---|
| nil; empty | 0 / false / 0 for each |
| LF; LF LF | 1 / false / 1; 2 / false / 2 |
| a LF; a; a LF b | 2 / false / 1; 1 / false / 1; 3 / false / 2 |
| CR LF; a CR LF b; a CR b | 2 / false / 1; 4 / false / 2; 3 / false / 1 |
| NUL; a NUL LF; a LF NUL | 1 / true / 1; 3 / true / 1; 3 / true / 2 |
| FF LF FE (hex bytes) | 3 / false / 2 |
| UTF-8 é followed by LF | 3 / false / 1 |

## Actual execution and boundaries

- Go1.25.10 executable version verified; default host Go1.27.1 is not used. No tool installation or dependency/version change.
- Private per-run tmp/cache/logs are owned and ignored without modifying common Git excludes/config. No production data or credentials.
- RED/GREEN/focused/full suite: NOT RUN yet. Byte metrics are not accounting LOC or a binary-classification verdict.
- Initial unit Risk: none; additive pure API in new files only, no existing shared behavior changed. Differential CLI checks are N/A for this unit.
- Rollback boundary for this first unit: remove only the new content API and its tests; no existing command/document/core behavior changes.
- Future constructors/decoders/acquisition/CLI/parity remain pending. Existing synthetic JSON and historical document checks are not executable evidence.
