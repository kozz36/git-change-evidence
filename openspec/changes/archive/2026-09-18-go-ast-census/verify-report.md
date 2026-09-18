# Canonical re-verification: go-ast-census

## Verdict and scope

**PASS — 10/10 requirements and 30/30 scenarios supported.** All 27 behavioral scenarios have passing executable evidence; the three explicitly configured manual scenarios are supported by current documentary inspection. Residual C3 and corrected GAC-008-S01 pass. All configured full/race/format/build gates passed. No requirement, scenario, gate, preservation check, or canonical execution step failed. This is one fresh canonical SDD verification for Public Preview Task 4, not a manual substitute verdict or a delivery approval.

All nine implementation task checkboxes are checked. Direct scan of `tasks.md` for `^\s*- \[ \]` found **no unchecked implementation lines**, including outside any partial slice. Historical proposal acceptance checklists are not substituted for tasks.

## Authority, status and ownership

- Consumed the exact parent SDD Session Preflight: auto / openspec / ask-on-risk, 400-line review budget, chain deferred. The explicit Task-4-only single-phase request restricts this run; no preferences were confirmed or persisted and no subsequent phase was launched.
- Change `go-ast-census`; OpenSpec; branch `docs/go-ast-census-archive`; HEAD `8318462d556be5e09de8ebb69c6f9664e28d68bf` (bounded documentary correction).
- Native `gentle-ai.sdd-status` v2 is authoritative: apply all_done, verify ready, archive ready, tasks 9/9, nextRecommended archive, blockedReasons [], notes []. Explicit optional verification is admitted without changing that recommendation.
- Native actionContext: repo-local; workspaceRoot and sole allowedEditRoot `/data/Projects/git-change-evidence-worktrees/go-ast-census-archive`. Git root and resolved source/test/document paths are within this root; artifacts were read directly from the supplied OpenSpec locators.
- Exact native command-output SHA-256 (including Markdown wrapper): `600c04e97320347eb6d5878c042c28eac3ef1793125081acdb89c6d42bd85072`. Status was not reconstructed from artifacts. Post-write status/hash is returned separately.
- Classical optional verification requires no retired attestation envelope, validator, attempt settlement, or manual replacement workflow. None was synthesized.
- Read proposal, specification, design, tasks, apply-progress, configuration, both complete production/test files, current consumer docs, root inventory constructor/accessors and guidance. CodeGraph index existed and reported up to date; read-only `codegraph status` and `codegraph query 'census'` preceded source exploration.
- Skill resolution: **fallback-path**, degraded self-healing because parent injected no skill paths. Loaded `/home/kozz36/.agents/skills/sdd-verify/SKILL.md`, its shared phase protocol, `/home/kozz36/.agents/skills/go-testing/SKILL.md`, and global `/home/kozz36/.pi/agent/gentle-ai/support/strict-tdd-verify.md`. No project override exists. Parent should inject indexed paths next time.
- Initial Git status contained existing modifications to `.atl/skill-registry.md` and `.gitignore`. They were neither authored nor cleaned by this executor. Only the canonical report is an authorized repository write; local command logs reside outside the repository.

## Commands and actual results

Exact commands, exit codes and combined stdout/stderr hashes are retained in `/tmp/gac-reverify-0ftfsq76/commands.json`; logs are numbered below. No evidence was uploaded. All listed commands exited 0.

| Log | Exact command | Exit | Output SHA-256 |
|---|---|---:|---|
| 0 | `gentle-ai sdd-status go-ast-census --contract gentle-ai.sdd-status/v2` | 0 | 600c04e97320347eb6d5878c042c28eac3ef1793125081acdb89c6d42bd85072 |
| 1 | `go version` | 0 | 98a041ad3c95c51454b9bfebc9934a61d9573cfdf1548c40c3309879f5f1feec |
| 2 | `go test ./census -run '^TestGAC005OwnershipAndOrder$/^GAC-005-S02_raw_path_order_is_deterministic$'` | 0 | 81af61bc5a63d4898489083cd13b19977c3e562a39cdb65836cee890519442c1 |
| 3 | `go test ./census` | 0 | 81af61bc5a63d4898489083cd13b19977c3e562a39cdb65836cee890519442c1 |
| 4 | `go test ./...` | 0 | 53fef60ca0fb35a69af23f75d2e36220fd551e1eb83f653ccf5dff173d4ea78c |
| 5 | `go test -race ./...` | 0 | 53fef60ca0fb35a69af23f75d2e36220fd551e1eb83f653ccf5dff173d4ea78c |
| 6 | `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 &#124; xargs -0 -r gofmt -l)"` | 0 | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| 7 | `go test ./...` | 0 | 53fef60ca0fb35a69af23f75d2e36220fd551e1eb83f653ccf5dff173d4ea78c |
| 8 | `go test -count=1 -v -coverprofile=/tmp/gac-reverify-0ftfsq76/coverage.out ./census` | 0 | 96dfb117eac96f30496ab28eb31f866f2adaf7af3d49e79e0aca4acc73eeb503 |
| 9 | `git diff --check` | 0 | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| 10 | `go test -count=1 ./...` | 0 | 1147675957e9d13ae9da40bbdfed6de6add6c6b62c7abd60489593d2ab29100b |
| 11 | `go test -race -count=1 ./...` | 0 | c785924e2e9d80791dd1b5c0a6e5bf7000c5df5fce21380dfd3a3fa8faa1cbff |
| 12 | `go test -count=1 ./census -run '^TestGAC005OwnershipAndOrder$/^GAC-005-S02_raw_path_order_is_deterministic$'` | 0 | f0244630eedff8c8dba471f48599c037a33d0eb44f6ba88f430bad2ecea8fef6 |

Logs 4/5/6/7 are the exact configured primary/race/format/build commands; build is explicitly `go test ./...`, not an invented `go build` replacement. Exact configured test commands returned cached successes. Logs 8/10/11/12 additionally force fresh census, full, race and isolated-C3 execution with `-count=1`; all passed. Six packages passed both uncached full and race runs. Nine census top-level tests and all subtests passed without skips. Check-only gofmt produced no output and made no edits. `git diff --check` passed.

Runtime: `go version go1.27.1-X:nodwarf5 linux/amd64`. WARNING: Go 1.25.10 is the module baseline, but that exact toolchain was not exercised. No separate linter is configured or was run. Go compilation/type checking passed through test/build gates.

## Complete requirements/scenarios census

Aliases from `census/go_ast_test.go`: V=TestGAC001AndGAC002Validation, C=TestGAC003Candidates, P=TestGAC004PhysicalEvidence, O=TestGAC005OwnershipAndOrder, L=TestGAC006LimitsAndGAC007Atomicity, S=TestGAC007SpanSafetyAndArithmetic, T=TestTriangulatedAtomicityAndOrder, H=TestDefaultLimitsAndIndependentHashes. E=`census/external_api_test.go::TestExternalValueAPIAndOwnership`. V/C/P/O/L/S/T/H provide unit/component evidence; E provides external-contract integration evidence. Only the three configured manual exceptions use documentary evidence.

| Scenario | Result and actual support |
|---|---|
| GAC-001-S01 | PASS — E: exact inventory/query/extractor bindings, literal spans and independent hashes. |
| GAC-001-S02 | PASS — V: unsupported version, empty receiver, empty selector and keyword fail with typed error and zero result. |
| GAC-002-S01 | PASS — V: duplicate supplied path fails before malformed duplicate content is parsed. |
| GAC-002-S02 | PASS — V: missing supplied path returns typed error and zero result; admission-before-parse confirmed in GoASTV1. |
| GAC-002-S03 | PASS — V: supplied extra path fails atomically before parsing its malformed content. |
| GAC-002-S04 | PASS — V: same-length digest substitution and unequal-length drift fail separately. |
| GAC-002-S05 | PASS — V: actual zero inventory value rejected without private-state fabrication. |
| GAC-002-S06 | PASS — V: root-constructed empty scope succeeds with digest binding and non-nil empty matches. |
| GAC-003-S01 | PASS — C: exact count and literal direct, parenthesized and single-index call fragments. |
| GAC-003-S02 | PASS — C: nested parentheses and multi-index wrappers retain outer-call fragment. |
| GAC-003-S03 | PASS — C: shadowed spelling and ambiguous indexed value included; alias, comment, string, non-selector and receiver-wrapper decoys excluded. |
| GAC-003-S04 | PASS — C: three repeated/legacy usages yield candidates; source/result contract contains no invariant classification. |
| GAC-004-S01 | PASS — P/E: literal UTF-8, CRLF, multiline and line-directive physical positions; independent call and selector coordinates. |
| GAC-004-S02 | PASS — P/H/E: changed locations/arguments differ in spans and file/fragment hashes; independent literal digest oracles. |
| GAC-005-S01 | PASS — O/E: mutations to original bytes, exported paths and match collection do not alter later retained observations. |
| GAC-005-S02 | PASS — O/T: raw-path cardinality guard and byte-order oracle; both input permutations compared to independently built full MatchV1 records, spans and hashes. |
| GAC-005-S03 | PASS — E: real external-package value consumer constructs inventory, calls exported API and verifies ownership. |
| GAC-006-S01 | PASS — L: exact file-count limit succeeds; one-unit excess fails with zero result. |
| GAC-006-S02 | PASS — L: per-file equality and one-unit excess, typed error and atomic failure. |
| GAC-006-S03 | PASS — L/S: aggregate equality/excess and near-uint64 overflow-safe arithmetic. |
| GAC-006-S04 | PASS — L: match equality succeeds; excess discards collected matches and all bindings. |
| GAC-006-S05 | PASS — L/H: zero file count rejected; complete default tuple independently asserted. |
| GAC-006-S06 | PASS — L: zero per-file byte limit independently rejected. |
| GAC-006-S07 | PASS — L: zero aggregate byte limit independently rejected. |
| GAC-006-S08 | PASS — L: zero match limit rejected even with empty scope. |
| GAC-007-S01 | PASS — L/T: malformed source after a candidate and after an earlier valid file publishes no partial result. |
| GAC-007-S02 | PASS — V/L/S: public error classes clear all result accessors; NoPos, foreign-file, past-EOF and reversed private spans fail; error propagation inspected. |
| GAC-008-S01 | PASS — manual: current docs explicitly explain shadow/index/repeated/legacy candidates and semantic, invariant, resource and authority non-guarantees. |
| GAC-009-S01 | PASS — manual: census-only layout exception and unchanged testing section verified against original introduction. |
| GAC-010-S01 | PASS — manual: explicit ask-on-risk pause in tasks and recorded accepted size:exception in progress; production 336/400 lines. |

Requirement census: GAC-001 **2/2**, GAC-002 **6/6**, GAC-003 **4/4**, GAC-004 **2/2**, GAC-005 **3/3**, GAC-006 **8/8**, GAC-007 **2/2**, GAC-008 **1/1**, GAC-009 **1/1**, GAC-010 **1/1**. Total **10 requirements / 30 scenarios**, none removed, rewritten or left unsupported. Support combines actual execution and inspected implementation; it does not claim every defensive branch was dynamically reached.

## Required manual evidence and historical dispositions

**GAC-008-S01:** Read the complete current `docs/census.md`, especially `Candidate meaning and limits`, against C's shadowed/indexed/repeated/legacy fixtures. It now states syntax-only textual observation; includes shadowed exact spellings and ambiguous indexed-value shapes; rejects invariant classification for repeated/legacy occurrences; disclaims package, method, import, type, alias, shadow, dependency and semantic identity; restricts completeness to validated supplied inventory and exact syntax query; disclaims parser heap, stack, time and process isolation; grants no approval, gating, merge, deploy or release authority. Current docs therefore cure historical F1. The CLI contract was not altered by this verifier.

**GAC-009-S01:** `git diff 3767f9b^ 3767f9b -- AGENTS.md openspec/config.yaml` shows only the narrow census public-layout exception and separately authorized three-scenario manual allowance. Current guidance retains the root-only default apart from census. Python byte comparison of the complete `testing:` section against `3767f9b^:openspec/config.yaml` returned equality; SHA-256 `873913843bc084eade9f979dd910d6e102412363f0b7e5ec853211c9148fad33`. Existing testing commands and rubric remain unchanged.

**GAC-010-S01:** `wc -l census/*.go` reports production 80 + 256 = **336** lines, tests 419 + 51 = **470** lines. Tasks explicitly paused for delivery selection; persisted apply-progress explicitly records approved `exception-ok (size:exception)` for the complete vertical, plus threshold/package-boundary review. The optional stack was not adopted; no assigned partial slice is claimed. The raw original launch record is unavailable, so historical authorization chronology is supported by persisted progress, not independently reopened raw transport. No new exception is inferred. This report-only task adds no implementation or delivery scope.

Historical C1/C2 are resolved by the explicit empty-selector and matching-comment-decoy cases. C3 is resolved: O captures matches, requires exactly four, then iterates that captured slice; T guards both lengths and checks independent complete ordered records. The isolated O subtest ran uncached. Lower-key exact ties remain defensive ordering cases with no demonstrated public reachability; no fabricated AST or public testing aperture is required. C4 has truthful TDD applicability/chronology tables; C5 has precisely the three configured manual exceptions. Historical reports remain historical FAILs, not rewritten as successes.

## Strict TDD and assertion quality

Strict TDD is active. The original Go work and test corrections match the authoritative Go-source-or-test rubric. This verification and the bounded GAC-008 correction add no executable behavior; their RED cycle is N/A, not manufactured. Historical progress contains original, characterization and residual-C3 TDD Cycle Evidence tables and documentary applicability notes.

| Check | Finding |
|---|---|
| Evidence table | Present; original RED/GREEN/TRIANGULATE/REFACTOR and correction safety nets documented. |
| Test files | Both referenced files exist and were read completely; 2/2 exercised successfully. |
| Historical RED | Progress records undefined API symbols before production; current GREEN is not substituted for that historical RED. |
| Current GREEN | Confirmed by uncached focused/census/full/race runs; check-only formatting also passes. |
| Triangulation | Actual independent drift, parser atomicity, physical coordinates, hashes and complete-order alternatives exist and pass. |
| Safety net/refactor | Recorded original and correction focused/full checks; already-correct characterization and docs-only N/A explanations retained. |

No missing or incomplete applicable TDD cycle found. WARNING: original raw writer record `/home/kozz36/.pi/agent/gentle-agents/tasks/mu2l2jvk-4-cuc2.json` is absent; historical order cannot be independently re-opened. The persisted chronology is disclosed, not replayed or invented.

Assertion audit of both files found **0 CRITICAL** banned patterns: no tautologies, ghost loops, type-only assertions alone, smoke-only tests or implementation-detail CSS assertions. Loops over returned candidates guard cardinality; arithmetic tests invoke production; meaningful empty outcomes have positive companions. Two noncritical quality limitations remain: ownership tests recheck selected fields rather than the entire record, and private span negatives check error presence rather than typed details. Neither recreates C3 or removes scenario support.

Test layers: **8 unit/component top-level tests in 1 file**, **1 external-contract integration test in 1 file**, **0 E2E** tests for this pure API; total **9**, all passing. No tests or goldens were changed.

Statement coverage from the fresh profile: `go_ast.go` 124/133 = **93.23%**, `query.go` 11/11 = **100%**; weighted total 135/144 = **93.75%** (Go: 93.8%), unweighted file mean 96.62%. Uncovered go_ast.go ranges: 48–49, 65–67, 70–72, 150–151, 160–161, 230–231, 247 (nine statements). Query has none. Go does not supply line/branch percentages here; no fabricated line/branch metric is reported. Configured coverage threshold is 0.

## Design coherence and remaining warnings

The pure public census depends one-way on root inventory contracts. Root constructors own inventory uniqueness. Census validates canonical identity, exact supplied raw-path set, lengths and hashes before whole-file parsing; SkipObjectResolution avoids semantic resolution. Only the approved three callee wrappers normalize. Physical unadjusted spans and literal full-call hashes are defensively owned; all failures return zero results. Sorting follows raw path, call start/end, selector start/end, documented in design. There is no process, filesystem, policy decoding or delivery-authority effect in the core.

Design deviation (nonblocking): validateSources copies content before byte-budget validation, whereas design planned deferred copying. Limits still precede parsing and are not resource-isolation promises; current docs explicitly explain that boundary. Unreachable defensive span branches are statically inspected, not claimed covered by runtime tests. Runtime-version and raw-history limitations remain warnings above, not suppressed findings.

## Historical preservation and next action

Before replacement, canonical bytes exactly equaled `verify-report.failed-docs.md`. All three historical files were compared byte-for-byte with their HEAD blobs and hashed; equality was true for each:

| Historical report | SHA-256 |
|---|---|
| verify-report.failed-initial.md | f088b195421dc4ba682c6fa48b2d016de86715c64cee9d0a8db9493d80dba7e9 |
| verify-report.failed-c3.md | 67973625eb60816fa3489713054e359914e3391dcb083a39e2b072914d7a8a97 |
| verify-report.failed-docs.md | e9938fc780edb3cd06adf7c7299c045aa5abccd0be313bddc5691ea76918551f |

Final report hash, post-write historical equality, tracked-file preservation, Git status and fresh native status are returned separately to avoid a self-referential hash. Exact verification blockers: **none**. Task 4's requested passing verification premise is met. **Task 5 is ready in native status**, with nextRecommended **archive** unchanged; actual archive still requires fresh native status and real permission/safety checks. No archive or delivery action was executed or authorized by this report.

Only this canonical report was authored in the repository. No source/doc/config/spec/task/progress/history edits, commit, archive, sync, push, PR, tag, visibility change, release, upload or child subagent occurred.
