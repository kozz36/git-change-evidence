# Tasks: Go AST selector-call census

> **Planning forecast only.** This ordinary documentation correction neither authorizes implementation nor selects a delivery shape. A future implementation preserves all 10 requirements and all 30 scenarios; it must not reduce acceptance, relocate production into tests, or claim census completeness from an API-only stub.

## Review Workload Forecast

The selected product goal is **one complete GoASTV1 vertical**: value API, validated supplied scope, AST extraction, physical evidence, ownership, external usability, documentation, and narrow layout alignment. It is one acceptance unit, even if a later human selects a review presentation with several dependent partitions.

| Measure | Complete vertical estimate |
|---|---:|
| Novel production in `census/query.go` plus `census/go_ast.go` | **325–395** lines; **HARD ceiling: 400** |
| Tests | 450–700 lines |
| Consumer documentation | 60–100 lines |
| Layout alignment | 6–12 lines |
| Gross authored implementation | **841–1,207** lines, excluding planning artifacts (which count if included in review) |
| Delivery strategy | `ask-on-risk`; size decision pending before implementation |

A single complete implementation review at gross 400 is impossible under this forecast. A later human may either explicitly accept a `size:exception` for one oversized review or explicitly select a coherent reviewed stack. The following stack is only a forecast: it is not adopted, does not grant an exception, and its partitions are **not independently executable increments or independently GREEN** when an earlier partition lacks the complete entry-point dependencies.

| Optional review partition of the one vertical | Production | Tests | Docs | Layout | Gross | Dependency and review meaning |
|---|---:|---:|---:|---:|---:|---|
| A. Public aperture, first external entry-point RED, and layout alignment | 80–100 | 95–130 | 0 | 6–12 | 181–242 | `AGENTS.md` and `layout.public_api` change before or in the same diff that first creates `census`; the early external consumer test constrains the API, but this partition alone cannot claim a working `GoASTV1`. |
| B. Inventory/supplied-set admission and atomic failure boundary | 80–100 | 115–170 | 0 | 0 | 195–270 | Depends on A’s contract shape; tests and implementation remain part of the same complete vertical, not a separately releasable census. |
| C. Complete-file AST candidate collection and match-limit atomicity | 90–105 | 120–200 | 0 | 0 | 210–305 | Depends on A–B; only the full A–D stack has the complete public entry point and acceptance surface. |
| D. Physical evidence, ownership/order, external hardening, and consumer docs | 75–90 | 120–200 | 60–100 | 0 | 255–390 | Depends on A–C and closes the vertical; it is the first point at which final GREEN and all acceptance evidence are meaningful. |
| **Aggregate** | **325–395** | **450–700** | **60–100** | **6–12** | **841–1,207** | Arithmetic equals the complete-vertical forecast. |

## Decision and stop condition

A human must make the size/delivery decision **before any future implementation**. No chain, review partition, or `size:exception` is selected or authorized here. Stop before apply, delivery, commit, or publication regardless of any future planning status.

## Exact future edit boundary

Only a separately authorized complete-vertical implementation may edit these real targets:

| Future path | Future responsibility |
|---|---|
| `census/query.go` | Public value contracts, defaults, private result bindings, and defensive accessors. |
| `census/go_ast.go` | Validated-scope-to-owned-result pipeline, including parsing, spans, hashes, ordering, and atomic errors. |
| `census/go_ast_test.go` | Co-located behavior, boundary, arithmetic, physical-evidence, and private span-safety tests. |
| `census/external_api_test.go` | `census_test` external-consumer and ownership tests. |
| `docs/census.md` | Consumer scope, predicate, bindings, limits, and non-guarantees. |
| `AGENTS.md` | Only the narrow public `census` subpackage layout exception. |
| `openspec/config.yaml` | Only `layout.public_api` for that same exception; its testing rubric remains immutable. |

`openspec/config.yaml` is also a **read-only verification reference** outside its narrowly authorized future `layout.public_api` edit. In particular, no future task may alter or treat these as edit targets: `testing.commands.primary`, `testing.commands.concurrency`, `testing.commands.format_check`, or the testing rubric.

## One contingent complete-vertical work unit

### 1. Complete GoASTV1 vertical: constrained externally, implemented end to end

- [x] **Boundary alignment before/with the first public package introduction** — In the same future work unit, update only the narrow `census` exception in `AGENTS.md` and `openspec/config.yaml` `layout.public_api` before or with the first `census` source file. Do not defer this to a later work unit, widen another layout rule, or change any testing configuration.
- [x] **RED** — First create `census/external_api_test.go` as `package census_test` so a genuine external consumer calls the eventual value entry point with a root-created inventory and can observe bindings and defensive output ownership. In the same bounded RED cycle, add only behavior tests in `census/go_ast_test.go` that the complete work unit will implement: query/default/limit validation; inventory and supplied-set negatives; complete-file candidate and parser atomicity; limit equality/excess; physical spans/hashes; mutation/order; and private invalid-span/arithmetic safety. Record the known focused `go test ./census` failure, including the absent or incomplete end-to-end `GoASTV1` entry point. This intentional RED is not a claim that a contract-only stub is usable or complete.
- [x] **GREEN** — Create `census/query.go` **and** `census/go_ast.go` together as the complete minimal vertical needed by the RED suite: exact public version/value contract and defaults; root accessor-based inventory identity; byte-for-byte supplied paths, length, and digest admission before parsing; complete-file `go/parser` traversal with only approved callee wrappers; match-limit atomicity; owned physical spans and literal fragment hashes; raw-byte deterministic order; private result bindings and defensive accessors; and typed `ce.ContractError` zero-result failures. Keep all production behavior in those two production files, with no semantic resolution, decoder, filesystem, policy, nil-inventory alternative, or production helper hidden in tests. Make the early external test GREEN only as part of this complete pipeline, never through API-only stubs.
- [x] **TRIANGULATE** — Within this same work unit, add independent literal-oracle and alternate cases for same-length versus unequal-length drift, malformed source after a candidate or earlier valid file, valid empty inventory, equality before one-unit excess, near-`uint64` aggregate arithmetic, UTF-8/CRLF/multiline/Go line directives physical positions, shifted locations/arguments, EOF/nested wrappers, raw newline/high-bit/backslash/case-distinct paths, and output/input mutation. Each added test either exercises behavior implemented in this work unit or remains an explicitly recorded RED until its paired complete-pipeline behavior is added; compare full ordered records, spans, and hashes rather than counts alone.
- [x] **REFACTOR** — Refactor only the two production files and their co-located fixtures/helpers after the complete focused suite is GREEN. Preserve the approved wrapper set (`ParenExpr`, `IndexExpr`, `IndexListExpr`), literal independent coordinate/hash oracles, zero-result atomicity, external value usability, and all scenario mappings. Add `docs/census.md` in this complete vertical with exact validated-scope completeness, syntax-only textual candidates, shadow/index ambiguity, lack of semantic/invariant classification, evidence-only authority boundary, and limits’ lack of parser resource-isolation guarantees.

## Authorized evidence correction — 2026-09-15

The TRIANGULATE checkbox is reopened because verification finding C3 identifies a path-only permutation oracle where this task requires full ordered records, spans, and hashes. Complete exact-count guards and realistic independent ordered-record expectations; distinguish defensive lower-key ties without fabricating ASTs or expanding the API. The other four historical checkboxes remain unchanged; reopening does not deny their recorded execution.

This approved follow-up permits only tests, truthful apply-progress evidence, task completion, the exact historical FAIL copy, and `openspec/config.yaml` `rules.verify` for the three named manual-evidence scenarios below. It supersedes the earlier layout-only configuration edit boundary solely for that allowance. All testing rubric modes, rows, commands, ten requirements, thirty scenarios, production files, and consumer documentation remain unchanged.

- [x] **C1/C2 characterization** — Add empty-selector rejection and matching-text comment-decoy cases in `census/go_ast_test.go`; use `census/external_api_test.go` only if needed. Characterize existing behavior without temporary production mutations or invented failing behavior.
- [x] **C4 truthful evidence** — Complete the evidence in `openspec/changes/go-ast-census/apply-progress.md` using available historical sources; explicitly mark unavailable chronology and N/A fields. Never describe current GREEN as historical RED.
- [x] **C5 documentary evidence allowance** — Add an explicit `rules.verify` manual-evidence allowance only for change `go-ast-census`, scenarios GAC-008-S01, GAC-009-S01, and GAC-010-S01. Keep strict TDD for behavioral scenarios and preserve all existing testing commands and rubric rows.
- [x] **Preserve historical FAIL** — Copy the exact unchanged `openspec/changes/go-ast-census/verify-report.md` bytes to `openspec/changes/go-ast-census/verify-report.failed-initial.md` before any later canonical report replacement. Retain its verdict and findings; record actual correction evidence in apply-progress.

Planning maintenance is separately authorized up to 40 gross added/deleted lines before acquisition. The correction is one bounded attempt up to 600 gross lines, including historical report preservation and later checkbox completion; independent verification is separately bounded to one attempt up to 400 gross lines, only if native routing permits. No extra reset, production change, acceptance reduction, unmanaged execution, or delivery is authorized.

## Residual C3 correction — 2026-09-15

TRIANGULATE is reopened solely for the remaining unguarded raw-path ordering subtest: capture matches, assert exact expected length, then iterate that captured slice. This follow-up permits only that test correction, truthful apply-progress evidence, checkbox completion, and exact preservation of the current FAIL as `openspec/changes/go-ast-census/verify-report.failed-c3.md`. Planning maintenance is limited to 10 gross lines; the separately admitted correction to 250 gross lines including the 126-line FAIL copy, then one independent verification up to 400 gross lines if native routing permits. No configuration, external-test, production, specification, or installed-skill edits; no manufactured RED, further automatic reset, or delivery.

## Acceptance traceability retained in the complete work unit

| Requirement | Scenarios retained | Primary future evidence |
|---|---|---|
| GAC-001 | GAC-001-S01, GAC-001-S02 | External value entry point plus co-located version/query/zero-result tests. |
| GAC-002 | GAC-002-S01, GAC-002-S02, GAC-002-S03, GAC-002-S04, GAC-002-S05, GAC-002-S06 | Accessor identity, raw-path set, length/digest precedence, zero inventory, and constructed empty-scope tests. |
| GAC-003 | GAC-003-S01, GAC-003-S02, GAC-003-S03, GAC-003-S04 | Complete-file parser fixtures for direct/wrapped/indexed, shadowed/ambiguous/decoy, repeated, and legacy candidates. |
| GAC-004 | GAC-004-S01, GAC-004-S02 | Literal physical range and fragment-hash oracles for UTF-8, CRLF, multiline, directives, and changed locations/arguments. |
| GAC-005 | GAC-005-S01, GAC-005-S02, GAC-005-S03 | Input/output ownership mutation, raw-byte ordering/final tie, and external-consumer tests. |
| GAC-006 | GAC-006-S01, GAC-006-S02, GAC-006-S03, GAC-006-S04, GAC-006-S05, GAC-006-S06, GAC-006-S07, GAC-006-S08 | Exact default tuple plus count, per-file, aggregate, and match equality/excess/zero-limit tests. |
| GAC-007 | GAC-007-S01, GAC-007-S02 | Malformed-source atomicity and every failure class returning the zero result. |
| GAC-008 | GAC-008-S01 | `docs/census.md` reviewed against shadowed, indexed, repeated, and legacy fixtures. |
| GAC-009 | GAC-009-S01 | Narrow `census` layout review; testing rubric is byte-for-byte unchanged. |
| GAC-010 | GAC-010-S01 | Measured 325–395 production forecast against the hard 400 ceiling and explicit pending size decision. |

## Future validation references (read-only configuration)

After separately authorized implementation reaches the complete-vertical GREEN state, retain the unchanged validation obligations:

```text
go test ./census
go test ./...
```

The second command is the unchanged `openspec/config.yaml` `testing.commands.primary` value. Use the check-only formatting gate by reading and running the unchanged `openspec/config.yaml` `testing.commands.format_check` value; it is a read-only configuration reference, not an edit target or an embedded replacement command. Run the exact `openspec/config.yaml` `testing.commands.concurrency` value only if the authorized Go behavior introduces concurrency-bearing execution. Do not run product tests, builds, or configuration commands for this planning correction.
