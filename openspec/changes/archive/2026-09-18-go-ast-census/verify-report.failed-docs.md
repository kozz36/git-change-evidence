# Canonical verification: go-ast-census

## Current verdict

**FAIL — current documentary requirement GAC-008-S01 is not satisfied.** Current executable behavior and the residual C3 assertion correction PASS. This is a fresh canonical optional SDD verification, not either historical FAIL replayed and not a substitute manual verification workflow. No production defect was reproduced.

Census: **10 requirements / 30 scenarios inspected; 9 requirements / 29 scenarios supported** (27 runtime scenarios PASS, two manual scenarios supported, one manual scenario FAIL). All nine task checkboxes are checked; no exact unchecked implementation lines matching `^\s*- \[ \]` remain. Completion markers do not cure the documentation finding.

## Authority and provenance

- Selected change: `go-ast-census`; store: OpenSpec; phase: verify only (Public Preview Task 4).
- Branch: `docs/go-ast-census-archive`; HEAD and observed origin/main: `16aaa4498f91c6d1ff91b330e1df6a47db26c922`. No fetch performed.
- Consumed the exact parent session preflight (auto/OpenSpec/ask-on-risk) without persisting choices. The explicit current Task-4-only request limits execution to this single phase; no back-to-back phase follows.
- Fresh native command: `gentle-ai sdd-status go-ast-census --contract gentle-ai.sdd-status/v2`, exit 0. Native v2 selects this change, apply all_done, verify ready, archive ready, tasks 9/9, nextRecommended archive, blockedReasons [], notes [].
- Native actionContext is repo-local; workspaceRoot and sole allowedEditRoot are `/data/Projects/git-change-evidence-worktrees/go-ast-census-archive`. Git root and resolved implementation paths were checked inside that root. Report is the only authored tracked artifact.
- Native status output SHA-256: `600c04e97320347eb6d5878c042c28eac3ef1793125081acdb89c6d42bd85072` (exact combined stdout/stderr, including native Markdown wrapper). Readiness is not recomputed from this report.
- Provider uses classical optional verification. No retired attestation envelope, attempt lifecycle, or legacy report-validation command was required or synthesized.
- Read current proposal, spec, design, tasks, apply-progress, config, guidance, consumer docs, both test files, production API/extractor and root inventory accessors. Historical failed-initial report was read; canonical prior bytes were read and confirmed identical to failed-c3 by `cmp -s` (exit 0).
- Skill resolution: fallback-registry; no parent skill paths were injected. Loaded `/home/kozz36/.agents/skills/go-testing/SKILL.md` through the registry fallback and global `/home/kozz36/.pi/agent/gentle-ai/support/strict-tdd-verify.md`. No project strict-TDD override exists. No separate sdd-verify SKILL.md was found; the assigned native SDD verify executor contract governs this run. Parent should inject indexed paths next time.
- CodeGraph was absent; mandatory `gentle-ai codegraph init --cwd "$PWD"` succeeded, followed by read-only `codegraph_codegraph_explore` calls through functions.mcp for GoASTV1 and validation/span/order helpers. Local ignored `.codegraph/` was created, not a source change.
- Initial worktree was NOT clean: `.atl/skill-registry.md` and `.gitignore` were already modified and `.pi/` untracked before verification. These ambient files were not authored or cleaned by this executor.

## Commands and reproducible evidence

Commands below were executed through `functions.bash` using Python subprocess capture with `/bin/bash`, in the authoritative workspace. Logs contain exact combined stdout/stderr, not reconstructed results. Evidence directory: `/tmp/gac-verify-evidence-868pkifa/`; `commands.json` records commands, exits and output hashes. Logs 0–7 correspond to the rows below. Temporary evidence is local, not uploaded or committed.

| Log | Exact command / purpose | Exit | Output SHA-256 |
|---|---|---:|---|
| 0 | `go test ./census -run '^TestGAC005OwnershipAndOrder$/^GAC-005-S02_raw_path_order_is_deterministic$'` | 0 | f0244630eedff8c8dba471f48599c037a33d0eb44f6ba88f430bad2ecea8fef6 |
| 1 | `go test ./census` | 0 | f0244630eedff8c8dba471f48599c037a33d0eb44f6ba88f430bad2ecea8fef6 |
| 2 | `go test ./...` — configured primary | 0 | 4ed47e797b2ce3b2359ab469613313ea1f90e7004cf88cdb203d67145f178dd5 |
| 3 | `go test -race ./...` — configured concurrency command, expressly requested | 0 | 1f01ab5ca30d6faaa548f710942e000a7f550e7b4d58dba04cc2c6d291ff3761 |
| 4 | Exact configured format command below | 0 | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| 5 | `go test ./...` — separate configured build invocation | 0 | 53fef60ca0fb35a69af23f75d2e36220fd551e1eb83f653ccf5dff173d4ea78c |
| 6 | `go test -count=1 -v -coverprofile=/tmp/gac-verify-evidence-868pkifa/coverage.out ./census` | 0 | 96dfb117eac96f30496ab28eb31f866f2adaf7af3d49e79e0aca4acc73eeb503 |
| 7 | `gentle-ai sdd-status go-ast-census --contract gentle-ai.sdd-status/v2` | 0 | 600c04e97320347eb6d5878c042c28eac3ef1793125081acdb89c6d42bd85072 |

```sh
test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"
```

Full and race commands passed all six packages. Fresh census execution passed nine top-level tests and all subtests without skips. Format emitted no output and performed no writes. Runtime: `go version go1.27.1-X:nodwarf5 linux/amd64`; go.mod baseline is 1.25.10. WARNING: these results do not independently establish execution under Go 1.25.10. No separate linter was configured or run; Go compilation/type checking passed through the test/build commands.

## Requirements and scenarios

Test aliases in `census/go_ast_test.go`: V=TestGAC001AndGAC002Validation, C=TestGAC003Candidates, P=TestGAC004PhysicalEvidence, O=TestGAC005OwnershipAndOrder, L=TestGAC006LimitsAndGAC007Atomicity, S=TestGAC007SpanSafetyAndArithmetic, T=TestTriangulatedAtomicityAndOrder, H=TestDefaultLimitsAndIndependentHashes. E=`census/external_api_test.go::TestExternalValueAPIAndOwnership`. Runtime rows are unit/component evidence except E, which is external-contract integration evidence. Manual rows are only the three configured exceptions.

| Scenario | Current evidence and result |
|---|---|
| GAC-001-S01 | PASS, E: exact inventory/query/extractor bindings, literal coordinates and independent hashes. |
| GAC-001-S02 | PASS, V: unsupported version, empty receiver, empty selector, keyword rejection; typed errors and zero result. |
| GAC-002-S01 | PASS, V: duplicate supplied path rejected before malformed-source parsing. |
| GAC-002-S02 | PASS, V: missing supplied path, typed error and zero result; validation-before-parser confirmed in source. |
| GAC-002-S03 | PASS, V: malformed extra path rejected before parsing. |
| GAC-002-S04 | PASS, V: same-length digest drift and unequal-length drift fail independently. |
| GAC-002-S05 | PASS, V: actual zero inventory value; no private-state fabrication. |
| GAC-002-S06 | PASS, V: root-created valid empty inventory retains binding and non-nil empty matches. |
| GAC-003-S01 | PASS, C: guarded cardinality and literal direct/parenthesized/indexed calls. |
| GAC-003-S02 | PASS, C: nested parentheses and IndexList wrappers retain full-call fragment. |
| GAC-003-S03 | PASS, C: shadowed and ambiguous indexed candidates included; alias/comment/string/non-selector decoys excluded. |
| GAC-003-S04 | PASS, C: repeated legacy usages yield distinct candidates without classification. |
| GAC-004-S01 | PASS, P/E: literal UTF-8/CRLF/multiline/directive physical call/selector positions. |
| GAC-004-S02 | PASS, P/H/E: changed location/arguments produce different ranges and independent file/fragment hashes. |
| GAC-005-S01 | PASS, O/E: caller input/output mutation does not change retained observations. |
| GAC-005-S02 | PASS, O/T: corrected isolated exact-count guard; raw-byte paths; independent complete ordered records in both input orders, including multiple positions per path. |
| GAC-005-S03 | PASS, E: genuine external value consumer, bindings and defensive exported observations. |
| GAC-006-S01 | PASS, L: file-count equality/excess and zero failure result. |
| GAC-006-S02 | PASS, L: per-file-byte equality and one-unit excess. |
| GAC-006-S03 | PASS, L/S: aggregate equality/excess and near-uint64 arithmetic. |
| GAC-006-S04 | PASS, L: exact match count and atomic one-unit excess. |
| GAC-006-S05 | PASS, L/H: zero file-count rejected; exact complete default tuple. |
| GAC-006-S06 | PASS, L: independent zero per-file-byte rejection. |
| GAC-006-S07 | PASS, L: independent zero aggregate-byte rejection. |
| GAC-006-S08 | PASS, L: zero match limit rejected even for empty scope. |
| GAC-007-S01 | PASS, L/T: malformed source after candidate and after earlier valid file publishes no partial evidence. |
| GAC-007-S02 | PASS, V/L/S: public errors clear all result bindings/matches; private NoPos/foreign/past-EOF/reversed spans fail; defensive propagation inspected. |
| GAC-008-S01 | FAIL, manual: current consumer docs omit explicit invariant-classification and parser-isolation non-guarantees, and shadow/index ambiguity explanation. |
| GAC-009-S01 | Supported, manual: narrow census exception in AGENTS/config; original merge diff and unchanged testing section verified. |
| GAC-010-S01 | Supported, manual: tasks record ask-on-risk pause; apply-progress explicitly records approved size:exception; 336/400 production lines. Original launch record no longer available for independent chronology. |

Requirement totals: GAC-001 2/2; GAC-002 6/6; GAC-003 4/4; GAC-004 2/2; GAC-005 3/3; GAC-006 8/8; GAC-007 2/2; GAC-008 0/1; GAC-009 1/1; GAC-010 1/1. No scenario or requirement was removed or rewritten.

## Current finding and historical dispositions

**F1 — CRITICAL / documentary specification failure:** `docs/census.md` now documents the CLI, not the earlier complete public API contract. Its syntax-only/no-semantic-resolution statement is useful but does not state that invariant classification and parser resource isolation are not guaranteed. It does not explain shadowed exact spellings or ambiguous indexed-value candidates. GAC-008 explicitly requires those consumer-facing non-guarantees. A targeted `git grep -n -i -E 'shadow|ambiguous.*index|parser.*(isolation|heap|stack)|invariant classification|tie.break' -- README.md docs census/query.go` returned no matches. Full docs/census.md was read; historical planning text is not a replacement for current consumer documentation. No fix was made.

Historical C1 and C2 are resolved: explicit empty-selector and matching-text comment-decoy cases exist and pass. Historical C3 is resolved: `got := result.Matches()` is followed by exact expected-length validation before iteration in the isolated raw-path subtest; T checks full independently constructed records and both lengths. The isolated subtest was executed, not inferred from its sibling. Lower-key exact ties are defensive implementation cases without established public reachability; no fabricated AST is required.

Historical C4 has corrected TDD tables with chronology, applicability and safety-net fields. Historical C5 has the explicit unchanged three-scenario manual allowance. That allowance permits documentary inspection; it does not make missing current documentation compliant. Historical favorable documentation findings describe older bytes and cannot be carried forward as current PASS.

## Strict TDD and assertion audit

Strict TDD is active in config and parent request. All-matching-row/strictest-wins selects the Go-source-or-test disciplines for the original work and corrections. No Go files were changed during this report-only phase. The synchronous API adds no concurrency, but the full race command was nevertheless run as requested.

| Check | Current assessment |
|---|---|
| TDD Cycle Evidence tables | Present: original complete vertical, characterization/documentary correction, residual C3 correction, with explicit N/A cases. |
| Referenced tests exist | 2/2 test files exist in the authoritative workspace and were read completely. |
| Historical RED | Recorded intentional undefined-symbol failure before production in apply-progress; not manufactured or replayed now. |
| Current GREEN | Confirmed by fresh nine-test census execution, focused residual C3, full, race, and format commands. |
| Triangulation | Actual alternate drift, parser atomicity, literal coordinates/hashes, raw-path and complete-record order cases exist and pass. |
| Safety net / refactor | Original and correction evidence recorded; documentary and already-correct characterization N/A explanations retained. |

WARNING: `/home/kozz36/.pi/agent/gentle-agents/tasks/mu2l2jvk-4-cuc2.json` is now absent. Historical original RED/refactor order and size authorization are supported by persisted apply-progress and historical verifier observations, not independently reopened raw logs in this run. Current GREEN is not historical RED. No missing TDD table or new incomplete cycle was found.

Assertion audit of both complete files: no tautologies, ghost loops, type-only assertions alone, smoke-only tests, CSS assertions, or production behavior hidden in fixtures. Candidate and ordering loops guard cardinality; private arithmetic cases invoke production; zero-result cases have meaningful positive companions. Remaining noncritical quality limitations: mutation tests recheck selected fields rather than complete records; private invalid-span tests assert error presence instead of typed details. These are warnings, not the former C3 ghost-loop defect.

Layers: eight unit/component top-level tests in one file; one external public-contract integration test in one file; zero E2E tests for this pure API. Total nine top-level tests, all PASS.

Statement coverage: go_ast.go 124/133 = 93.23%; query.go 11/11 = 100%; weighted aggregate 135/144 = 93.75% (Go displays 93.8%); unweighted file mean 96.62%. Config threshold is 0. Go profile does not provide line/branch percentages. Uncovered go_ast.go ranges: 48–49, 65–67, 70–72, 150–151, 160–161, 230–231, 247 (nine statements); query.go has none.

## Architecture, manual evidence and workload

Pure census depends one-way on root inventory values. It validates canonical identity, supplied raw-path bijection, byte lengths and hashes before parsing; consumes complete supplied bytes with SkipObjectResolution; normalizes only ParenExpr/IndexExpr/IndexListExpr; reports physical unadjusted owned spans and literal call hashes; returns zero results on errors. Result accessors defensively clone matches/paths. Sorting implements the design's raw path, call start/end, selector start/end key. No semantic or delivery-authority machinery was introduced in this API.

Design warning retained: validateSources copies content before byte-budget validation rather than deferring copies as designed. Limits still precede parsing; this is not a demonstrated behavioral failure, and it makes the missing consumer resource-isolation caveat particularly important.

Manual GAC-009 evidence: `git diff 3767f9b^ 3767f9b -- AGENTS.md openspec/config.yaml` shows the census-only layout exception and separately authorized manual allowance. The complete current `testing:` section is byte-identical to `3767f9b^:openspec/config.yaml`; section SHA-256 `873913843bc084eade9f979dd910d6e102412363f0b7e5ec853211c9148fad33`. Current guidance retains the narrow exception; no testing command or rubric was edited.

Manual GAC-010 evidence: production 80 + 256 = 336 physical lines, under hard 400; tests 419 + 51 = 470 lines. Apply-progress explicitly records size:exception for the complete vertical and threshold/package-boundary review. Tasks' optional stack was never adopted; no assigned partial slice or chain is claimed. Current report-only work does not grant a fresh size exception, implement additional tasks, or expand PR scope. Historical CLI additions on main are not attributed to this report-only phase.

## Preservation, exact blockers and next action

Both historical failures remain immutable and retain their historical verdicts:

- `verify-report.failed-initial.md`: SHA-256 `f088b195421dc4ba682c6fa48b2d016de86715c64cee9d0a8db9493d80dba7e9`.
- `verify-report.failed-c3.md`: SHA-256 `67973625eb60816fa3489713054e359914e3391dcb083a39e2b072914d7a8a97` (also the canonical report's prior hash).

Exact current verification finding: F1, GAC-008-S01 documentary non-guarantees missing. No unfinished task checkbox, failing executable command, or native blockedReason is asserted. This report grants no delivery authority.

Task 5 distinction: native archive dependency remains ready and nextRecommended remains archive; optional failed verification does not override native archive admission. The requested clean-PASS premise for Public Preview Task 4 is not met. The parent must receive this new documentation finding and decide the requested next work; no archive was executed. Archive would still require fresh native status and actual safety/permission checks.

Only the canonical report was authored. No source/test/config fix, archive, spec sync, commit, push, PR, tag, visibility change, release, upload or child subagent was performed. Exact final report hash and post-write Git/native-status evidence are returned separately to avoid a self-referential report hash.
