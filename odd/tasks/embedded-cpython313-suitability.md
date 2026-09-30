# Embedded CPython 3.13 strict suitability (#151)

Status: approved issue #151 (GitHub label read back); isolated evidence only, not parser adoption.
Base: c080d5b0930066c8f71c22e4944c6a3ff0236e46. Branch: test/embedded-cpython313-suitability.
Mirror: Engram odd/embedded-cpython313-suitability/tasks in CNSIC-bound parent session.

## Tasks

- [x] S1 — Delegated writer: raw-byte span and parse+compile probe fixtures, exact predicate and mixed-file all-or-error demonstration. Preparation/multi-file trigger; no production API. Detect supported encodings and fail unsupported ones explicitly; do not select production encoding policy. RED/GREEN/TRIANGULATE/REFACTOR from repo Go instructions; Go 1.25.10 root and explicit nested tests. Original eight tests and independent checks PASS, but differential native compile(bytes) reproduced false acceptance of BOM+coding:utf8 and false rejection of second-line cookies after code/prior cookie. S1 remains incomplete; owner authorized narrow byte-validation correction. Independent evidence hashes: results 66fbba0136b33cfb097cc2b1c04ff5ec24f460ccfb20f8d0ff354c3f9c498b9b, log 26edd7d361db94841534411f8ab3fc99db681880a8be09990b9319bbdba77cd1. Corrected native-byte parse/compile boundary and three regression cases; writer nine tests and fresh independent nine-plus-differential tests PASS. Source hashes and historical/negative/correction evidence retained in report. Independent mutations not rerun; broader suitability remains NOT PROVEN. Local work-unit commit follows verification; publication pending separate decision.
- [ ] S2 — Delegated bounded resource/cancellation and final host import evidence; report PASS/FAIL/UNKNOWN independently. Offline runs <=4 CPUs, 8 GiB, 20 minutes; no scanned source execution. Stop a failed axis and preserve evidence. No production limits or signed provenance claims.
- [ ] S3 — Independent verification of frozen evidence and limitation wording; exact parser suitability NOT PROVEN if any required axis remains FAIL/UNKNOWN. Offer publication only after separate owner decision.

## Boundaries

Allowed: internal/embeddedcpythonprobe/**, docs/embedded-cpython313-suitability.md, this task file. External disposable scratch may reuse generated #149 assets; do not mutate retained prior evidence. No root dependencies, generated runtime in Git, production census, parser selection, release or project-specific pilot. No push/PR/merge without separate decision.

Delivery strategy: ask-on-risk → owner selected stacked-to-main (small PRs to main) when cumulative forecast exceeds ~400. Slice 1: corrected S1 spans/validity/probe batch, ~425 authored lines including historical negative/correction/audit evidence and task document; kept as one coherent proof slice rather than trimming tests. Independent verification PASS; local commit and publication decision next. Slice 2: resource/capability evidence on a fresh branch after the first evidence slice is separately authorized and merged. Publication/merge remain separate decisions; no automatic push. Per-task heuristic is advisory; preserve tests and readable code.

## Evidence

Approval: https://github.com/kozz36/git-change-evidence/issues/151 OPEN, status:approved.
Mapping: existing smoke takes a decoded string and returns counts only; AST UTF-8 columns cannot directly index arbitrary original encoding. Missing physical spans/batch/resource guarantees are proof gaps, not authorized contract changes.
