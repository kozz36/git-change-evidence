```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:05a270a6b2c8c2867ca36e8ecaa43acc589be2edadcdf83d8221c38042e76bbf
verdict: fail
blockers: 5
critical_findings: 5
requirements: 4/10
scenarios: 24/30
test_command: go test ./...
test_exit_code: 0
test_output_hash: sha256:53fef60ca0fb35a69af23f75d2e36220fd551e1eb83f653ccf5dff173d4ea78c
build_command: go test ./...
build_exit_code: 0
build_output_hash: sha256:53fef60ca0fb35a69af23f75d2e36220fd551e1eb83f653ccf5dff173d4ea78c
```

## Verification Report

**Change:** go-ast-census. **Mode:** strict TDD. **Verdict:** FAIL, evidence incomplete; no failing production behavior was reproduced.

All required commands passed. This does not establish all 30 scenarios: three behavioral scenarios have partial test oracles, and three documentation/process scenarios have manual evidence but no covering runtime test. The installed phase skill requires passing covering tests and does not permit counting manual-only scenarios without a project exception. The latter is an evidence-method mismatch, not an observed product defect.

### Completeness and authoritative context

- Read proposal, specification, design, tasks, apply-progress, all four `census/*.go` files, `docs/census.md`, `openspec/config.yaml`, `go.mod`, root inventory constructor/decoder/accessor contracts, `contract.go`, and root inventory contract/ownership tests directly.
- Actual specification headings: **10 requirements, 30 scenarios**. Fully runtime-covered requirements: GAC-002, GAC-004, GAC-006, GAC-007 (**4/10**). GAC-001, GAC-003, GAC-005 are partial; GAC-008 through GAC-010 are manually supported but runtime-untested.
- **5/5 task checkboxes checked; 0 unchecked.** No lines matching `^\s*- \[ \]` remain in `tasks.md`. Checked state does not prove the full-record triangulation obligation: see C3.
- Consumed native `gentle-ai.sdd-status` v2: unique selected change, OpenSpec store, verify ready, apply all_done, all five tasks complete, no blocked reasons. No readiness was inferred from apply-progress.
- `actionContext.mode` is `repo-local`; workspace and allowed edit root are `/data/Projects/git-change-evidence-worktrees/go-ast-census`. All implementation targets are beneath that real Git root. This verifier's narrower write boundary is this report only.
- CodeGraph was not reprobed or initialized; used the explicitly approved direct-read fallback.
- No review-authority prerequisite was imposed. Findings concern technical evidence and SDD completeness, not delivery authority.

### Commands and results

Commands ran synchronously in the authoritative workspace. Hashes cover exact combined stdout/stderr bytes, not rendered tool output.

| Purpose | Exact command | Exit | Observation |
|---|---|---:|---|
| Focused | `go test ./census` | 0 | Census package passed from Go cache. |
| Primary test | `go test ./...` | 0 | Root, census, CLI, internal/git, internal/inventory, internal/publication passed from cache. |
| Independent build invocation | `go test ./...` | 0 | Separately invoked as configured; same six cached successes and same exact output hash. |
| Check-only formatting | Command below | 0 | Empty output; no format writes. |
| Fresh existing-test execution and coverage | `go test -count=1 -v -coverprofile=/dev/stdout ./census` | 0 | All 9 top-level tests and their subtests passed; 93.8% statement coverage. No test, fixture, harness, or profile file was created in the workspace. |

Exact formatting command:

```sh
test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"
```

Additional exact output SHA-256 values:

- Focused: `81af61bc5a63d4898489083cd13b19977c3e562a39cdb65836cee890519442c1`.
- Formatting: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- Fresh verbose/coverage run: `05897cba01cc0e6695b080cd0d5e35c442170995a645513afe7943c00c31839b`.
- `go test -race ./...`: not run, N/A; the added surface is pure and synchronous, with no concurrency-bearing execution.

### Scenario compliance matrix

Test aliases identify exact top-level names; scenario-labelled subtests are quoted below with their source spelling (Go's verbose output replaces spaces with underscores).

- V = `census/go_ast_test.go::TestGAC001AndGAC002Validation`.
- C = `census/go_ast_test.go::TestGAC003Candidates`.
- P = `census/go_ast_test.go::TestGAC004PhysicalEvidence`.
- O = `census/go_ast_test.go::TestGAC005OwnershipAndOrder`.
- L = `census/go_ast_test.go::TestGAC006LimitsAndGAC007Atomicity`.
- S = `census/go_ast_test.go::TestGAC007SpanSafetyAndArithmetic`.
- T = `census/go_ast_test.go::TestTriangulatedAtomicityAndOrder`.
- H = `census/go_ast_test.go::TestDefaultLimitsAndIndependentHashes`.
- E = `census/external_api_test.go::TestExternalValueAPIAndOwnership`.

COMPLIANT means a covering existing test passed in the fresh run. PARTIAL means passing tests omit a required part; UNTESTED means no covering runtime test. Manual findings are retained, not counted as runtime completion.

| Scenario | Test/evidence | Result and limitations |
|---|---|---|
| GAC-001-S01 | E | COMPLIANT: exact inventory/query/extractor bindings, full literal match coordinates and hashes. |
| GAC-001-S02 | V / `GAC-001-S02 invalid query` | PARTIAL: unsupported version, empty receiver, keyword selector tested; empty selector absent (C1). |
| GAC-002-S01 | V / `GAC-002-S01 duplicate supplied path` | COMPLIANT: different-content duplicate produces typed path error and zero result. |
| GAC-002-S02 | V / `GAC-002-S02 missing supplied path` | COMPLIANT: missing entry produces typed error and zero result; malformed entry is missing, not a supplied parse-precedence oracle. |
| GAC-002-S03 | V / `GAC-002-S03 extra supplied path` | COMPLIANT: malformed extra produces path error before parsing and zero result. |
| GAC-002-S04 | V / `GAC-002-S04 length and digest drift` | COMPLIANT: independent same-length digest and unequal-length errors, zero results. |
| GAC-002-S05 | V / `GAC-002-S05 zero inventory` | COMPLIANT: literal value API, no forged private state. Root constructor/decoder own duplicate inventory rejection. |
| GAC-002-S06 | V / `GAC-002-S06 valid empty scope` | COMPLIANT: root-created empty inventory, retained inventory binding, non-nil empty matches. |
| GAC-003-S01 | C / `GAC-003-S01 exact and single-index selector calls` | COMPLIANT: count checked before loop; direct/parenthesized/indexed literal fragments in order. |
| GAC-003-S02 | C / `GAC-003-S02 nested parentheses and index-list normalize` | COMPLIANT: nested parentheses and multi-index fragment; normalized-selector coordinates not separately asserted for this fixture. |
| GAC-003-S03 | C / `GAC-003-S03 shadowed and ambiguous candidates remain syntax-only`, plus S01 fixture | PARTIAL: shadow, alias, string, non-selector, receiver-wrapper and indexed ambiguity covered; comment decoy absent (C2). |
| GAC-003-S04 | C / `GAC-003-S04 repeated legacy sites have no classification` | COMPLIANT: three usages, count and middle literal fragment; no classification field exists. |
| GAC-004-S01 | P / `GAC-004-S01 UTF8 CRLF multiline and line directives`, E | COMPLIANT: independent literal physical call/selector spans; Unicode call byte positions. |
| GAC-004-S02 | P / `GAC-004-S02 equal counts retain distinct locations and fragments`, H, E | COMPLIANT: shifted/changed argument evidence differs; H/E independently hash literal files/fragments. Position and argument changes are combined, not isolated. |
| GAC-005-S01 | O / `GAC-005-S01 input and output mutation cannot alter evidence`, E | COMPLIANT: input path/content, output path and record/collection mutations; only selected retained fields rechecked, not entire evidence. |
| GAC-005-S02 | O / `GAC-005-S02 raw path order is deterministic`, T / `GAC-005-S02 input permutation retains the complete ordered evidence` | PARTIAL: path-only assertions, no full-record/final-tie oracle; isolated ordering subtest has unguarded result loop (C3). |
| GAC-005-S03 | E | COMPLIANT: actual external consumer of value API, literal match evidence and defensive observation checks. |
| GAC-006-S01 | L / `GAC-006-S01 exact file count then excess` | COMPLIANT: fixed two-file input, limit two succeeds, one fails atomically. |
| GAC-006-S02 | L / `GAC-006-S02 and S03 exact byte budgets then excess` | COMPLIANT: equality succeeds, decrementing per-file limit by one fails atomically. |
| GAC-006-S03 | Same L subtest; S / `GAC-006-S03 near uint64 aggregate arithmetic` | COMPLIANT: aggregate equality/excess and subtraction-safe near-uint64 cases. |
| GAC-006-S04 | L / `GAC-006-S04 exact match count then atomic excess` | COMPLIANT: two matches accepted, limit one publishes zero result; within one file, not cross-file. |
| GAC-006-S05 | L / `GAC-006-S05 zero files`, H | COMPLIANT: zero invalid independently; exact complete default tuple. |
| GAC-006-S06 | L / `GAC-006-S06 zero file bytes` | COMPLIANT: independent zero limit, typed error and zero result. |
| GAC-006-S07 | L / `GAC-006-S07 zero total bytes` | COMPLIANT: independent zero limit, typed error and zero result. |
| GAC-006-S08 | L / `GAC-006-S08 zero matches` | COMPLIANT: independent zero limit, including empty supplied scope. |
| GAC-007-S01 | L / `GAC-007-S01 parser error after candidate is atomic`; T / `GAC-007-S01 earlier valid file cannot escape later parser error` | COMPLIANT: partial-file and later-file parse failures publish no accumulated result. |
| GAC-007-S02 | V/L shared `requireContract`/`requireZero`; S / `GAC-007-S02 invalid positions and spans are rejected` | COMPLIANT: public failures assert all zero accessors; private NoPos/foreign/past-EOF/reversed spans fail. Unreachable production span-error propagation supported by inspection, not forged AST/inventory injection. |
| GAC-008-S01 | `docs/census.md`, C fixtures | UNTESTED: manual documentation review confirms semantic/invariant/isolation disclaimers; no test checks consumer documentation (C5). |
| GAC-009-S01 | `git diff -- AGENTS.md openspec/config.yaml` | UNTESTED: manual diff proves only narrow census layout changes and unchanged testing rubric; no covering runtime test (C5). |
| GAC-010-S01 | Tasks forecast, apply-progress, explicit current human confirmation, physical line measurement | UNTESTED: explicit size exception is recorded and confirmed, production cap met; no covering runtime test of delivery decision history (C5). |

### Implementation correctness and design coherence

| Surface | Independent inspection finding |
|---|---|
| Neutral API / inventory identity | Private root inventory value and defensive accessors; census checks canonical SHA-256, policy link, non-nil entries. No policy argument, decoder round trip, nil alternative, or fabricated private state. |
| Scope admission | Raw-byte sorted paths, duplicates, both membership directions, length before digest; all files validated before any parsing. Root inventory constructor sorts and rejects duplicates; decoder preserves canonical order. |
| AST predicate | Full `parser.ParseFile` with typed byte content and `SkipObjectResolution`; only ParenExpr/IndexExpr/IndexListExpr unwrapped; identifier receiver exact. Traversal continues into call children. No filesystem, import/type/alias/shadow resolution or policy effects. |
| Physical evidence | File ownership and raw token bounds checked before potentially clamping APIs; `PositionFor(..., false)` ignores directives. Half-open strict spans, selector containment, whole-file and exact call-fragment SHA-256. |
| Ownership/order | Inputs copied; private result stores value bindings; `Matches` clones collection and paths. Sort key is raw path, call start/end, selector start/end, as documented. Implementation appears correct; oracle insufficiency remains C3. |
| Atomicity/limits | Failures return `ResultV1{}`; parse errors discard partial AST and prior matches. Positive uint64 limits use inclusive boundaries and subtraction-before-addition aggregate arithmetic. Match cap checked before append. No parser isolation claim. |
| Architecture/layout | Two production files in the sole approved public census exception; one-way dependency on root contracts; no shell effects or extra package. |
| Design deviation | `validateSources` copies all content at lines 137-140 before byte-budget checks, whereas design deferred content copying until after budgets. Bounds still precede parsing and no allocation/isolation guarantee is made: WARNING, not a reproduced spec violation. |

### Strict TDD compliance

Resolved the authoritative `testing.rubric`: Go-source-or-test matches; strictest mode and evidence union select strict-tdd. The default row is not selected. Concurrency row is N/A.

| Check | Result | Evidence |
|---|---|---|
| TDD Cycle Evidence table present | Yes, incomplete schema (C4) | One complete-vertical cycle at apply-progress lines 10-13. |
| Reported tests exist | 2/2 | Both new co-located files read fully. |
| Historical RED | Reported, not independently replayed | Table records `go test ./census` exit 1 and undefined File/QueryV1/ResultV1/SelectorCallQueryV1 before production. No separate historical output log or chronological source snapshot was supplied/retrieved. |
| Current GREEN | Confirmed | Fresh existing suite: 9/9 top-level tests pass. |
| Historical TRIANGULATE | Reported; cases exist | Earlier-valid-file atomicity and input permutation are present. Full-record permutation claim exceeds actual assertions (C3). |
| Historical REFACTOR / focused rerun | Reported only for chronology | Table records final focused exit 0 and retained API/extractor split; fresh GREEN independently confirms current state, not past ordering. |
| Safety net / required table fields | Incomplete (C4) | No SAFETY NET column or explicit per-task applicability; no required Written/Passed status markers or separate triangulation/refactor fields. Go files are new; layout/config were modified narrowly. |

Both test files exist and current GREEN is proven. The artifact narrates one cycle for the five checked steps; it does not provide five independently complete evidence rows. Missing table fields are a reporting-completeness finding, not proof that RED or refactoring never happened. Frozen tests were not rewritten to manufacture historical RED.

### Assertion quality and limitations

| File / lines | Finding | Severity |
|---|---|---|
| `census/go_ast_test.go:217-224` | Ordering subtest loops over `result.Matches()` without its own exact-count guard. Selected alone, an empty result makes its assertions vacuous. The sibling ownership subtest guards four matches in the full run, but is not a prerequisite for isolated selection. | CRITICAL, part of C3 |
| `census/go_ast_test.go:335-346` | Subtest named complete ordered evidence compares only two paths, not file/fragment digests or spans. It cannot establish full-record equality or final position tie-breakers. | CRITICAL, part of C3 |
| `census/go_ast_test.go:209-224`; `census/external_api_test.go:45-50` | Mutation checks protect selected fields but do not compare the complete retained record and bindings after every mutation. External test otherwise has strong independent literal coordinate/hash oracles. | WARNING |
| `census/go_ast_test.go:292-319` | Invalid-span tests check an error exists, not typed field/code or zero Span; no positive EOF endpoint, equal-key tie-breaker, or selector-containment failure fixture. These omissions reduce design triangulation, without a reproduced production failure. | WARNING |

No tautologies, CSS/implementation-detail assertions, mock-heavy assertions, production behavior hidden in test helpers, or smoke-only test files were found. Candidate iteration checks cardinality before looping. Empty-scope and nonmatching-query checks have meaningful non-empty companions. Indexing assumptions can panic on regression (therefore fail); they do not silently pass.

### Test layers, coverage, and quality metrics

- Unit/component behavior: 8 top-level tests in `census/go_ast_test.go`, using real root-created inputs and pure production APIs/private safety helpers.
- External public-contract integration: 1 top-level test in `census/external_api_test.go`; real root inventory plus public census, no mocks or filesystem.
- E2E: 0, not needed for the pure synchronous API. Total: 9 top-level tests across 2 files; all pass.
- Go statement coverage: **93.8% aggregate**; configuration threshold is 0. Go's profile is statement/block coverage, not line or branch coverage.
- `census/query.go`: all instrumented statements covered (100%).
- `census/go_ast.go`: 9 uncovered statements in profile ranges 48-49, 65-67, 70-72, 150-151, 160-161, 230-231, 247. Per-file percentage was not separately computed. These include defensive span errors, alternate membership exits, and exact equal-key return; coverage does not replace scenario/oracle review.
- Formatting passed; Go compilation/type checking passed through configured commands. No separate linter was selected or run. No coverage profile was persisted.

### Review workload and frozen boundary

- Measured physical production: `census/query.go` **80**, `census/go_ast.go` **256**, total **336/400**, including comments and blanks. No cap waiver or test relocation.
- Tests: 368 + 51 = **419** physical lines (below the 450-700 forecast, which is an estimate, not a substitute for acceptance). Consumer docs: **39** lines.
- `git diff -- AGENTS.md openspec/config.yaml` shows exactly the two narrow layout replacements. Testing rubric/commands unchanged. Existing ambient `.atl/skill-registry.md`, `.gitignore`, and `.pi/` changes were present before verification and were not edited or attributed to implementation.
- Tasks retain the old planning-only pending-decision wording; apply-progress explicitly records `exception-ok (size:exception)` for the complete vertical, and the current human prompt confirms approval. This later approval supplies the boundary; the forecast's optional stack is not an adopted chain. No PR exists as verification evidence and none was created.
- Threshold crossings and package-boundary rationale are recorded in apply-progress. Complete vertical, not an API-only slice, is present.
- Before and after runtime execution, all **167** Git-listed tracked/untracked nonignored files had the same manifest digest as the envelope. Digest algorithm: SHA-256 over byte-sorted raw Git paths, each followed by NUL and the binary SHA-256 of its file bytes. The report does not exist in that baseline and is excluded from the evidence revision.
- This report is the only permitted authored change, remains below 400 gross lines, and is admitted as exact in-memory bytes before its sole write. No corrections, new harness, child actor, archive, commit, push, delivery action, or model-setting changes were performed.

### Exact blockers

1. **C1 — CRITICAL / partial GAC-001-S02:** no empty-selector test; the three invalid-query rows at `census/go_ast_test.go:65-67` do not cover that explicitly required case.
2. **C2 — CRITICAL / partial GAC-003-S03:** candidate fixture at `census/go_ast_test.go:121` has string/alias/non-selector decoys but no comment containing matching call text. Required comment-decoy behavior has no covering test.
3. **C3 — CRITICAL / partial GAC-005-S02 and assertion quality:** path-only permutation oracle, unguarded isolated ordering loop, and absent full-record/final-position tie-breaker evidence contradict the checked TRIANGULATE task's full ordered-record obligation.
4. **C4 — CRITICAL / incomplete strict-TDD evidence:** apply-progress's combined table omits required safety-net/status/separate-phase fields from installed strict-TDD verification guidance. Existing narrative and present GREEN do not make those fields complete; historical command order remains reported rather than independently evidenced.
5. **C5 — CRITICAL / runtime-untested GAC-008-S01, GAC-009-S01, GAC-010-S01:** manual docs/layout/size-decision checks are favorable, but no covering runtime tests exist and config provides no manual-verification exception to the installed phase skill. This is an evidence-method mismatch for inherently documentary/process scenarios, not a finding of incorrect docs, layout, size approval, or delivery behavior.

**Conclusion:** FAIL; not ready for archive. Report the five bounded findings for an explicit human/orchestrator decision. No automatic remediation, further attempt, or acceptance reduction is proposed.
