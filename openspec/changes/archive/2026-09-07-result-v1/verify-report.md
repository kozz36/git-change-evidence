```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:15c8576cc6cd497419f8a5f4d1c7dd7572c727799025f621cfb3a85bf66f2870
verdict: pass_with_warnings
blockers: 0
critical_findings: 0
requirements: 13/13
scenarios: 38/38
test_command: GOTOOLCHAIN=go1.25.10 go test ./... -count=1 -json
test_exit_code: 0
test_output_hash: sha256:b0d9c249486e909037250e8245932ce9541c2930630a585041a92e0509dfe0b9
build_command: GOTOOLCHAIN=go1.25.10 go test ./... -count=1
build_exit_code: 0
build_output_hash: sha256:aac83722cc6729ca1cd9cda1415bb61d9faa060cd4f807cf890d06bad7c2e3b5
```

# Verification Report — Result V1 authorized re-verification

## Current verdict and historical disposition

**PASS WITH WARNINGS.** All 13 requirements / 38 scenarios have passing runtime coverage, with the direct/private/composed qualifications below. All 12/12 authoritative implementation checkboxes are checked. Current blockers: **none**; current CRITICAL findings: **0**. This is independent formal verification, not delivery authority or phase advancement.

- Workspace `/data/Projects/git-change-evidence-worktrees/result-v1-unit9-mutation-census`; branch `test/result-v1-unit9-mutation-census`; HEAD `17854e7539ca571107a01b35907695a979f1295c` unchanged.
- The prior formal report, SHA-256 `4d1fca4f02823705d8e023206f42899cf06a0ba016a9d2b93b1c40e89276fa2a`, failed with F1/F2 scenario gaps and F3 unavailable history (11/13 requirements, 36/38 scenarios). It is retained **byte-for-byte** in the historical appendix; its FAIL, hashes, counts and pending instructions describe the prior run only, not the corrected candidate.
- User-authorized remediation added two supplemental tests and recovered historical evidence, without production or existing-test changes. This run independently read the full proposal, four specs, design, tasks, progress, config, changed Go source/tests and relevant dependencies. The original failure was an evidence gap, not a reproduced source defect.

| Finding | Current disposition / independently confirmed evidence |
|---|---|
| F1 | **CLOSED.** `result_v1_construction_evidence_test.go:9-44`, `TestNewAccountingResultV1ClassifiesRenameByPrimaryPath`: the policy maps `docs/old.txt` to category 0 and `src/new.go` to category 1; public construction asserts one exact `src/new.go`/`source` countable entry (4 additions, 2 deletions), zero docs/other totals, and source-only totals in policy order. The private premise guard is not the sole oracle. |
| F2 | **CLOSED.** Same file `:46-86`, `TestNewAccountingResultV1ReproducesIdentityFromIndependentAntecedents`: two separate strict Policy decodes and two separate strict Inventory decodes use copied canonical antecedent bytes; each Inventory links its respective Policy. Two public Result constructions use independently built equivalent snapshots and assert both derived links and equal canonical bytes/digests. This is independent reconstruction, not alias assignment or the different-chain decoder fixture. |
| F3 | **CLOSED.** Supplied raw session SHA-256 and nine exact command/result links verified below; pre-change safety net and historical RED/GREEN/TRIANGULATE/REFACTOR are now available. No current RED was invented or replayed. |
| W1 | **RETAINED WARNING:** `.pi/gentle-ai/support/strict-tdd-verify.md` is absent. Config/prompt and assigned strict module were applied directly; no strict-TDD check was skipped. |
| W2 | **RETAINED WARNING:** native SHA-256 revisions plus hostile bytes have composed/private boundary coverage, not one public exact-fixture golden. The matrix preserves this distinction. |
| W3 | **RETAINED WARNING:** direct valid-report/invalid-binding validator, invalid-subject bound-builder, deletion-only accumulation overflow, and legacy required-line-total overflow vectors remain absent. Shared helpers and unchanged source support these nonblocking residuals; no new requirement or cached-view hardening mandate is imposed. |

## Current complete scenario matrix

Actual headings independently counted: contract **4/12**, accounting **5/14**, Result provenance **3/8**, Evidence provenance **1/4** = **13/38**. Registry keys CON through REPORT refer to the exact source/test names preserved in the appendix registry; all were cross-checked against current files and observed passing in FOCUS/FULL. `REN` and `REP` are the two tests named above. `P` = public, `H` = private helper, `C` = composed runtime boundaries plus inspected wiring; C is not an unexecuted end-to-end claim. All current matrix rows were re-evaluated, not merely inherited from the old verdict.

### result-v1-contract — 4/4 requirements, 12/12 scenarios covered

| Requirement | Exact scenario | Trace / result |
|---|---|---|
| Neutral Result V1 construction and inspection | Construct a neutral result from exact typed evidence | CON + OWN + DEC: PASS (P); literal public canonical golden, digest, links, revisions and views. |
| Same | Caller offers only digest assertions | ANT + OWN + DEC: PASS (C); private digest-only forgeries rejected, public zero antecedents return zero, public signatures require typed documents rather than legacy policy/digest inputs. No negative-compilation test claimed. |
| Same | Failed construction has no partial result | OWN + DEC: PASS (P); antecedent, malformed snapshot, late path, stored/required-sum overflow return exact zero. |
| Closed Result V1 wire contract and content identity | Canonical hostile-path fixture | WIRE + CON: PASS (C); exact `/wpmaWxl` golden is private encoder coverage with synthetic links; public golden separately verifies real antecedents and hostile-byte encoding. |
| Same | Empty committed change set | INV + WIRE + CAN: PASS (C/P); distinct-revision empty public snapshot with nonempty inventory, policy-order zero totals, arrays and derived identity. |
| Same | Alternate wire representation | CAN + SHAPE + VALUES: PASS (P/H); order, whitespace, LF, escaped category, null array, Base64, fractional/exponent/string numbers rejected. |
| Defensive ownership and deterministic construction | Caller mutates a source or returned buffer | OWN + DEC + VALUES: PASS (P/H); effective full-buffer alternates and independently copied oracles. |
| Same | Equivalent snapshot orders | CON: PASS (P); raw prefix ordering and identical canonical bytes/digest. |
| Strict antecedent-aware canonical decoding | Wrong exact antecedent despite valid syntax | CAN + ANT + DEC: PASS (P/H); valid alternate and cross-policy chains, no digest-syntax shortcut. |
| Same | Recomputed evidence differs | DEC + SHAPE: PASS (P/H); classification/measurement/counts/totals/threshold/ratio/order mutations; forbidden unavailable actual is private shape coverage. |
| Same | Unknown, duplicate, and authority-bearing nested fields | SHAPE + ROOT + VALUES + DEC: PASS (C); nested duplicate/escaped/unknown/authority paths run through private validators/root/materializer; public zero-result materializer funnel separately passed. Not all nested rows execute the public decoder. |
| Same | Policy-derived authority-like string remains evidence | CAN: PASS (P), valid `release` policy construction/decode; `merge` also has H shape coverage. |

### result-v1-accounting-evidence — 5/5 requirements, 14/14 scenarios covered

| Requirement | Exact scenario | Trace / result |
|---|---|---|
| Valid immutable committed-snapshot boundary | Valid native SHA-256 snapshot with hostile path bytes | SNAP + VALUES + WIRE + CON: PASS (C), qualified by W2; private native-ID/raw-byte boundary and encoding/public wiring, not a public exact-fixture run. |
| Same | Invalid revision or path input | SNAP + VALUES + OWN + DEC: PASS (H/P); both widths, malformed/mixed/equal-with-entry IDs, invalid/duplicate raw paths; public fatal funnels. |
| Same | Rename classification uses the committed primary path | REN: PASS (P); distinct old/current categories, exactly one literal current-path entry and complete policy-order totals (F1 closed). |
| Inventory is provenance-only and excluded from accounting | Inventory identity changes without changing accounting | INV: PASS (P); empty versus nonempty linked inventory, identical entries/totals/observations and distinct links/digests. |
| Same | Nonempty inventory with empty committed snapshot | INV: PASS (P); no entries, explicit policy-category zero totals. |
| Exclusive classification and recomputable policy-category totals | Overlap, default, and non-countable evidence | DOMAIN + CON + COMP: PASS (C/P); private exact four-path vector and public equivalent overlap/default/non-countable vector; literal exclusive classifications/totals. |
| Same | Constructor input order does not establish result order | CON: PASS (P), exact `a`, `a/ff`, `a0`, `ff` order. |
| Same | Realizable arithmetic overflow is fatal rather than invented evidence | DOMAIN + OBS + OWN + DEC: PASS (H/P); MaxUint64+1 stored additions, required threshold/numerator/reference sums, late overflow and exact-zero public results. |
| Closed factual threshold, ratio, and unavailable-observation model | Available threshold and exact ratio | OBS + COMP + DEC: PASS (H/P); 2+1 > 2, 3/4 > 0.5, exact fields/limits. |
| Same | Exact decimal comparison intentionally differs at a legacy float64 precision edge | PREC: PASS (H + legacy P); exact 9007199254740993/18014398509481984 > 0.5 versus legacy false. The fixture's symbolic snapshot IDs are used only by legacy Account/private domain, not a valid public Result construction. |
| Same | Non-countable and zero-denominator observations remain embedded | OBS + CON + COMP: PASS (H/P); unavailable false variants and available-zero distinction, no invented line values. |
| Same | Exceeded evidence is not authority | CON + OBS + DEC + WIRE: PASS (C/P); literal closed evidence and strict decoding, no authority result field. |
| Legacy accounting compatibility | Legacy caller remains unaffected | DOMAIN + unchanged-source comparison: PASS (P); no Policy/Inventory requirement added to Account, legacy validation/order/float arithmetic retained. |
| Same | Result and legacy accounting agree on shared-domain vectors | COMP + PREC + DOMAIN: PASS (P/H); literal expected entries plus differential totals/availability/ordinary ratios and accumulation failure. |

### result-v1-provenance — 3/3 requirements, 8/8 scenarios covered

| Requirement | Exact scenario | Trace / result |
|---|---|---|
| Exact validated antecedent chain | Byte-identical antecedents reproduce provenance | REP: PASS (P); two separately strict-decoded equal Policy/Inventory chains, two equivalent snapshots, both Result links and equal bytes/digests (F2 closed). |
| Same | Inventory policy link is wrong but all digests are well formed | ANT + OWN + CAN: PASS (H/P); isolated supplied cached-link mismatch, valid cross-policy inventory rejection, public zero results. |
| Same | Same accounting with a different antecedent is not interchangeable | INV + CAN: PASS (P); different valid inventory identity leaves accounting unchanged and cannot decode original Result with alternate inventory. |
| Result identity is system-derived and complete | One raw path byte changes the identity | WIRE + CON: PASS (C); private ff/fe, LF/CRLF and padding vectors distinguish bytes/digests; public golden binds complete wire. |
| Same | Reordered candidate is not repaired | DEC: PASS (P); entry/total/observation swaps reject with exact zero Result. |
| Result V1 remains neutral evidence-only provenance | Authority-bearing candidate provenance | SHAPE + ROOT + VALUES: PASS (C); closed normalized metadata/discriminator rejection at nested boundaries; unknown CNSIC/project keys remain forbidden by closure, not by a new vocabulary policy. |
| Same | Policy-derived authority-like value remains evidence | CAN + SHAPE: PASS (P/H); exact valid `release` values accepted, structural strings not misrepresented as canonical Policy limits. |
| Same | Out-of-scope integration request | Public signatures + CON/DEC/BIND and inspected dependencies: PASS boundary coverage; pure typed calls expose no acquisition/delivery request API. No external integration or Git-authenticity proof is claimed. |

### evidence-v1-provenance — 1/1 requirement, 4/4 scenarios covered

| Requirement | Exact scenario | Trace / result |
|---|---|---|
| Provenance and reproducibility | Bound report derives all provenance | BIND + REPORT: PASS (P); fixed P/I/R digest literals, five derived values, identical legacy report bytes/digest, no Result body. |
| Same | Well-formed but mismatched accounting digest | BIND: PASS (P); independent valid alternate chain/report supplies the different accounting digest; isolated validator error. All five provenance fields separately tested. |
| Same | Standalone report bytes are not provenance proof | REPORT: PASS (P); standalone syntactic decode succeeds, bound validation fails against unrelated Result. |
| Same | No projection or authority expansion | BIND + REPORT + CLI focus: PASS (C); binding fixture has exceeded evidence and produces identical unchanged legacy bytes; old JSON/human projections independently pass. No direct bound-report projection test claimed. |

## Fresh command and runtime evidence

Commands ran from this worktree. Go reported `go1.25.10 linux/amd64`; `GOTOOLCHAIN=go1.25.10 go env GOROOT` selected that toolchain's `bin` first on PATH for check-only `gofmt` (both exit 0). Hashes are SHA-256: command strings without a terminal LF; outputs are exact combined stdout/stderr captured in memory. No test log/profile/document artifact was created.

NEW:
```bash
GOTOOLCHAIN=go1.25.10 go test . -count=1 -run '^(TestNewAccountingResultV1ClassifiesRenameByPrimaryPath|TestNewAccountingResultV1ReproducesIdentityFromIndependentAntecedents)$' -json
```
FOCUS (all selected names checked against source and the JSON top-level PASS set):
```bash
GOTOOLCHAIN=go1.25.10 go test . -count=1 -run '^(TestAccountUsesFirstMatchingRawByteGlobAndDefault|TestAccountRejectsInvalidPolicyBeforeAccounting|TestAccountReturnsEmptyResultOnOverflow|TestPathGlobGrammar|TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable|TestAccountingDomainRejectsCheckedAccumulationOverflow|TestAccountSharedDomainTriangulatesLegacyBehavior|TestResultAccountingDomainProjectsPolicyViewLosslessly|TestNewReportV1CanonicalBytesAndDigest|TestDecodeCanonicalRejectsClosedAndAuthorityBearingDocuments|TestNewReportV1RejectsNonSHA256ProvenanceDigest|TestProjectEvidenceCanonicalJSON|TestProjectEvidenceHumanText|TestProjectEvidenceRejectsInvalidEvidenceAndUnsupportedFormat|TestProjectEvidenceOwnsCanonicalJSONOutput|TestValidateResultAntecedentsRequiresExactPolicyInventoryChain|TestNewAccountingResultV1InventoryOnlyChangesIdentity|TestResultV1AndAccountAgreeOnSharedDomainVectors|TestNewAccountingResultV1ConstructsExclusiveCanonicalEvidence|TestDecodeAccountingResultV1RejectsNoncanonicalRepresentations|TestDecodeAccountingResultV1RequiresExactAntecedentsAndAcceptsIndependentChain|TestDecodeAccountingResultV1AcceptsEmptyEqualRevisionEvidence|TestResultV1ObservationShapesRejectClosedShapesDuplicatesAndAuthorityKeys|TestResultV1RevisionShapesRejectClosedShapesDuplicatesAndAuthorityKeys|TestResultV1RootEnvelopeShapeDelegatesClosedSubtrees|TestResultV1RootEnvelopeShapeAcceptsStructuralSubtreeVariants|TestResultV1RootEnvelopeShapeRejectsClosedRoot|TestDecodeAccountingResultV1RoundTripsOwnedSelfConsistentEvidence|TestDecodeAccountingResultV1RejectsSemanticAndOrderMutations|TestDecodeAccountingResultV1ReturnsZeroOnFatalBoundaries|TestResultV1EntryShapesRejectClosedShapesDuplicatesAndAuthorityKeys|TestResultV1TotalShapesRejectClosedShapesDuplicatesAndAuthorityKeys|TestResultV1DecodeValuesMaterializesOwnedCandidate|TestResultV1DecodeValuesRejectsLexicalAndStructuralCandidates|TestResultV1DecodeValuesAcceptsLexicallyCanonicalEmptyPath|TestResultV1DecodeValuesDefersSnapshotSemanticsToConstructor|TestResultV1ExactRatioPrecisionEdge|TestResultV1PrecisionEdgeIsTheOnlyAvailableRatioComparatorDifference|TestResultV1ObservationsUseClosedOrderAndUnavailableVariants|TestResultV1StrictGreaterThanExcludesEqualityBoundary|TestResultV1RejectsRequiredLineTotalOverflow|TestResultV1NonCountableTotalsRemainUnavailableBeforeAddition|TestNewAccountingResultV1OwnsIngressAndEgress|TestNewAccountingResultV1ReturnsZeroOnFatalFailure|TestResultV1ReportBindingDerivesAllProvenance|TestValidateReportV1ResultProvenanceRejectsEveryMismatchAndPreservesLegacySurfaces|TestNewReportV1FromResultReturnsZeroEvidenceForInvalidChains|TestResultV1ReportBindingPreservesLegacyStandaloneReports|TestValidateResultSnapshotAcceptsNativeIDsAndRejectsInvalidPaths|TestResultV1ModelZeroValue|TestResultV1ModelPopulatedAccessors|TestResultV1ModelDefensiveCopiesPreserveState|TestResultV1ModelReturnsEmptyViewsForPopulatedDocument|TestResultV1ModelDiscriminators|TestResultV1CanonicalWireUsesClosedVariants|TestResultV1CanonicalWirePreservesPathByteIdentity|TestResultV1CanonicalWireHostilePathAndIdentity|TestResultV1CanonicalWireUsesEmptyArrays|TestNewAccountingResultV1ClassifiesRenameByPrimaryPath|TestNewAccountingResultV1ReproducesIdentityFromIndependentAntecedents)$' -json
```
CLI:
```bash
GOTOOLCHAIN=go1.25.10 go test ./cmd/git-change-evidence -count=1 -run '^(TestRunProjectCanonicalJSON|TestRunProjectHumanTextAndRepeatedInput)$' -json
```

| ID / exact command where not expanded above | Exit; observed result | Command SHA-256 | Output SHA-256 |
|---|---|---|---|
| NEW | 0; 2 top-level tests, 0 subtests | `dae84143b80402f87b29fea681126361724b7dcf25dc179da055ba635b8715a4` | `c5aac1cadb9f5ab63ed7f416fe768d71e9e70c39bfb3d3b2a7d91ff026e65ee0` |
| FOCUS | 0; 60 top-level, 684 subtest PASS events | `8b20214dbc6c89259ed657f44e1dfb642f1a9fd51a1acfa9a7bf1b58c3c2d74d` | `2d8ab2eb776194cf3dd935be6e495fc48678d5f10e1517c2426dd9a6673b751e` |
| CLI | 0; exactly 2 top-level tests | `f096f5483eeaa1ffa073c72b5f013fde4db58891cf7896f1a7780cea41c17990` | `b836ebeb4b055c8f54900b205e3a66b1e4f46c032a9e8b3b93b0c8e3f208e324` |
| BUILD: `GOTOOLCHAIN=go1.25.10 go test ./... -count=1` | 0; independent configured build/primary execution, five packages | `666e35311c141e11bf1aba0b54414e0af31315f173ceb18eb452e33cc1f9bb9f` | `aac83722cc6729ca1cd9cda1415bb61d9faa060cd4f807cf890d06bad7c2e3b5` |
| FULL: `GOTOOLCHAIN=go1.25.10 go test ./... -count=1 -json` | 0; 124 top-level, 928 subtest PASS events | `65c2f745f89c384dc1e241aa645867b2fd7076e1bd0eefc3bfafbcc1f539a04b` | `b0d9c249486e909037250e8245932ce9541c2930630a585041a92e0509dfe0b9` |
| COVER: `GOTOOLCHAIN=go1.25.10 go test -cover ./...` | 0; root fresh, other four packages cached; informational threshold 0 | `97b4f91e348427e5eda576e56710324e2d86c6901cd215f0c5305905f6664634` | `759bac0a9715e4d2ba931a6e697cd9b0eea51a7239eb35f05d2a5fb4a663f972` |
| FMT: exact configured command below | 0; empty output | `638bb9ec6e1345ec8613b87474e47b7359ca694c77eb1d40764fbde975d92663` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `git diff --check` | 0; empty output | `466c2f308b48c7661d646fdd068fbecea974c665fe65dbf8ed508f224180ce0b` | same empty-output hash |
| `git diff --no-index --check /dev/null result_v1_construction_evidence_test.go` | 1; empty diagnostics, normal no-index difference status, not a whitespace failure | `54bbde4b629d77a072e3671f67bb8f42ff1d1b3359e84f512a9c28facd6631ed` | same empty-output hash |
| `git diff --no-index --check /dev/null openspec/changes/result-v1/verify-report.md` | 1; empty diagnostics, same no-index semantics; prior report and persisted replacement checked | `3c3d179a86841c32f275ec5dab730e263eb287f90f67fe1dbf64b376ba83afd3` | same empty-output hash |
| `git diff --stat` | 0; tasks/progress 103 additions, 1 deletion | `bcf996cc2c0b7eb7e8ddcd15c1c22e596e11fd6556591a38ebc741509209e1a5` | `44f7e6c082607669b41e62dd3cd1df590f7da2031421f6516baf8a9fbe0a9f0a` |
| `git diff --numstat` | 0; progress 102/0, tasks 1/1 | `b80e81bb535d5373b0bcffab313126d80bf123960a1d602c4150ff7dc262c07a` | `ae0326ced695758e6adee0597224f2f569e820645e25c72b07c3a17876635b13` |

```bash
test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"
```

| Package | Top-level PASS | Subtest PASS events | Statements covered |
|---|---:|---:|---:|
| root | 89 | 888 | 94.9% |
| cmd/git-change-evidence | 13 | 13 | 83.5% |
| internal/git | 13 | 23 | 85.2% |
| internal/inventory | 4 | 4 | 75.8% |
| internal/publication | 5 | 0 | 79.3% |
| **Total** | **124** | **928** | No aggregate invented |

**0 failing / 0 skipped runtime events; 0 missing selected tests.** `GOTOOLCHAIN=go1.25.10 go list -f '{{.ImportPath}}|{{.Dir}}|{{join .TestGoFiles ","}}|{{join .XTestGoFiles ","}}' ./...` (exit 0) identified current-platform files; every top-level name in them executed. `TestAcquireReportsUnsupportedPlatform` is excluded by its existing non-Linux build constraint, not a skipped/passing Linux test. Subtest events include intermediate nodes, not assertion counts. No standalone linter was configured; compilation/test checks passed. Coverage is package statements only, not changed-file/branch coverage; no profile was written on the report-only surface. Race remains N/A: no new concurrency.

Tooling failures are distinguished from product failures: MCP returned `MCP not initialized`; read-only `codegraph status` and `codegraph explore AccountingResultV1` succeeded. Index reported one pending added file (the new test), so its bytes were read directly; no init/index/sync was performed. The initial exact `cat .pi/gentle-ai/support/strict-tdd-verify.md` exited 1 (missing file). Three inline `python3 - <<'PY'` census wrappers exited 1: first wrongly included the non-Linux test in expected names, second had a `subprocess.run` argument TypeError while assembling the go-list command, third incorrectly required no-index diff exit 0. These were verifier harness errors, not code/test failures; no repository files were changed. The first two full runs were not used as final evidence; FULL above was rerun and platform-aware reconciliation passed. The no-index commands were rerun with their exact status/empty diagnostics preserved above.

## Strict TDD, recovered history, and assertion quality

**Strict-TDD compliant with disclosed history/tooling qualifications.** Apply-progress contains actual stage/command/result `TDD Cycle Evidence` tables, including corrective records. Units **1–8: 8/8** Go-producing units now have available historical cycle evidence; Unit 0 is baseline/layout and Unit 9 is mutation/census, not a new behavior cycle. Units 2–8 remain recorded history rather than independently witnessed past executions. Referenced files exist; renamed historical planned test names map to the current actual registry. Current GREEN is independently established for all 24 changed test files (53 top-level functions, including existing functions in modified `accounting_test.go`), compared with 23/51 before correction.

F1/F2 are accurately qualified **supplemental test-only coverage of already-correct production**. The earlier baseline selected six existing names while the two new names were absent: that was not RED and was never accepted as a failing test. No production fix/refactor was needed; current focused/full/check-only gates passed. No blanket TDD exception, physical mutant, manufactured undefined-symbol RED, or historical replay was authorized or performed.

F3 source: `/home/kozz36/.local/share/pi/subagents/sessions/2026-09-04T12-56-31-002Z_01a06c7d-fc1a-7d4b-8d78-f72bce0f75f7.jsonl`, verified SHA-256 `c35a6dfafdc84bb1d2b9dce49da1826ff509a0c223326e32361ea3685446bfd5`. Only named records were inspected for command/result linkage; no full-conversation dump is included. Each pair below has matching tool-call ID and `bash` tool name. These are **2026-09-04 UTC historical commands/results**, not current execution. Parent observation `14046` is corroborative context, not the proof substituted for these raw records. Original omission reflected the previous supplied-artifact scope; it did not prove TDD was not performed.

| Stage; exact record pair (lines) | Historical command/result facts | Command SHA-256 / recorded-result-text SHA-256 |
|---|---|---|
| Baseline `faf02a8e→23c62a0f` (28→35) | 12:57:30; seven-name pre-change focus, cached PASS | `27ea420ccff5a42eb355db92ddf58027209063215efb50ed1217887275776253` / `ae1008b614f21f37333dac7cb41387e8540ae5df6237bb2e7785e0c29d74cbf9` |
| RED `a92abee9→16abeb2c` (43→44) | 12:58:54; four-name foundation focus; exit 1, undefined `legacyAccountingDomain` / `accountingEntry` | `5c8cdbb981188369a0f4368e9c82e41bf9a1d3e736d82a23e1a25b1dc5f5dc00` / `00d751312084e9acedae04c75984e936f2fb4fcddca81dd039bf71adea6c669c` |
| GREEN `88d33d02→0aa7bb0e` (51→52) | 12:59:58; same foundation focus, PASS 0.001s | `5c8cdbb981188369a0f4368e9c82e41bf9a1d3e736d82a23e1a25b1dc5f5dc00` / `0424937ff475aee5b629800a27f2aa0ce7aaa1373fb78158e025fee8fc6bef06` |
| TRIANGULATE `b7c3e8dc→1399f152` (55→56) | 13:00:57; adds `TestAccountSharedDomainTriangulatesLegacyBehavior`, PASS 0.001s | `c66b37e4e0d0bc30eaaf49f896f0418879269c3591f11706ae3878fbdd15f244` / `0424937ff475aee5b629800a27f2aa0ce7aaa1373fb78158e025fee8fc6bef06` |
| REFACTOR `c0966aa0→73d21fb5` (57→58) | 13:01:17; historical gofmt then expanded focus cached PASS and diff checks | `c99ccdb1fbf448c5b70b2f94b43841f3861a17eda66b2dd22c42d17f2b65b774` / `e95c06fca8e61939bf50fcc3edc7bb6cea78b61f071c25a1f4e2c9098317b4c7` |
| Adapter `0d379941→29adafdc` (61→62) | 13:01:50; adds `TestResultAccountingDomainProjectsPolicyViewLosslessly`, PASS 0.010s | `eadc2ca9019320aac2df91c690bf91c41e68508c4cbb848331ad79ba02fe6785` / `8fd79a58456accca1e24d80355f38c6b6270c495f9b843980473058a1867eb51` |
| Final focus `210d5aba→518269eb` (69→70) | 13:02:28; foundation focus PASS 0.006s | `5c8cdbb981188369a0f4368e9c82e41bf9a1d3e736d82a23e1a25b1dc5f5dc00` / `daefb769121540d9088e1ff8660f4bfb2ee72b2df3559a6e7d36ac59d0740637` |
| Final full `210d5aba→e3a40873` (69→71) | `GOTOOLCHAIN=go1.25.10 go test ./...`; five packages PASS | `b902a8f7ffb819678445298c322f07e87eae42276819777f720380e9bf400992` / `a60a589883e037e80f15f36402d512f0fef26a7447778294941a4aa856140697` |
| Final checks `210d5aba→20e1405f` (69→72) | configured check-only gofmt and `git diff --check`; no output, isError false | `172c394120ed451df1767da0d1d9f573fe4062824f47841e73922f9f11b7829a` / `025e763c5d283f5563e843dc9016daa9427b98edf3d660539db90b067cfebffb` |

Historical hashes above include the original `cd /data/Projects/git-change-evidence-worktrees/result-v1-pr5-accounting && ` command prefix; result hashes cover stored text (including tool indentation/status), not reconstructed process streams. The exact focus strings are retained in progress's F3 command block. The verified refactor command, **not run now**, was:
```bash
cd /data/Projects/git-change-evidence-worktrees/result-v1-pr5-accounting && GOTOOLCHAIN=go1.25.10 gofmt -w accounting_domain.go accounting_domain_totals.go accounting.go accounting_test.go && GOTOOLCHAIN=go1.25.10 go test . -run '^(TestAccountUsesFirstMatchingRawByteGlobAndDefault|TestAccountReturnsEmptyResultOnOverflow|TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable|TestAccountingDomainRejectsCheckedAccumulationOverflow|TestAccountSharedDomainTriangulatesLegacyBehavior)$' && git diff --check && git diff --numstat && git diff --stat
```

Assertion audit independently reread all 24 changed test files: **0 confirmed tautologies, ghost loops, type-only-alone or smoke-only behavioral substitutes, CSS assertions, or mock-heavy substitutes**. Literal nonempty fixtures, guarded cardinalities, exact error field/code plus zero results, actual ownership alternates, independent expected wire hashes, precise rational boundaries, and the two new public scenario premises remain meaningful. Constant/model/synthetic-link tests are supplemental contracts, not proof of valid public provenance. Pure core unit/component tests include private seams and composed public APIs; CLI focus is two in-process shell/core integration tests; no browser/network E2E or consumer integration is claimed.

## Current design, mutation, workload and integrity conclusions

- Re-inspected pure constructor → antecedents/snapshot → shared first-match/check-add domain → exact observations → ordered wire/identity flow, materialize/rebuild/link/equality decoder and canonical-byte-authoritative report binding. Source and runtime remain coherent with the approved design/error-boundary clarification. No new semantic defect or authority/effect expansion found; the historical design qualifications remain applicable, not new blockers.
- **Unit 9: recorded 11/11 physical assertion kills; independently rerun now: 0 mutants.** All listed killers passed on restored source; all six restoration hashes in the appendix matched current bytes again. No exhaustive mutation adequacy is claimed, and no production mutation was performed.
- Review Workload Forecast: selected `stacked-to-main` respected. This is only the assigned report-only formal re-verification of the authorized F1/F2/F3 correction; no chained implementation slice added. No `size:exception` used. User's current budget is 400 production changed LOC, tests/docs separate; correction production delta **0**, new test **86** lines below 160, docs **104** changed lines (102 progress + 2 tasks). Root census **68**, repository **83** Go files; the approved supplemental file is the only new root Go path. Historical multi-slice source census against `a2b9654`: 17 production files / 1,076 changed LOC, now 24 test files / 2,246 changed LOC; not one oversized PR or a new exception.
- This executor changed **only this existing report**; no code/tests/tasks/progress or other artifact edited. Final report is **443 lines**; untracked-inclusive current work candidate = 86 test + 104 docs + 443 report = **633 changed lines**, below 650. Pre-remediation report was 242 lines and its correction candidate was 432 total; those historical counts are not silently rewritten.
- Integrity uses sorted relative-path manifest lines `SHA256  path` plus LF, with no `./` prefix. Envelope evidence revision covers all tracked candidate file bytes plus the new test, excluding this report; it is not merely HEAD and is not a caller-supplied settlement digest. Before/after manifest equality is checked around report persistence.

| Before = after protected surface | SHA-256 |
|---|---|
| All 82 original tracked Go files; individually equal to HEAD and original formal verification | `9e95f953e6ebe929d35ab2ff1fca0f91f8af794d6d937eacce8bb3916cdf7794` |
| All 83 Go files including supplemental test | `f6d6ccb3ac59c548d95f2e665a6e7303f504a341c12c97545afb5c961964e0e8` |
| Repository production-only manifest under the definition above | `1d945dff74d22eac6450e8c27b5d2b80e29fba4e68ad5bb1df0be8425484ce54` |
| Supplemental test | `41d5e0e5374cb18b20ae4e94dc2eb7244dc3c3bd55dde9d25bfde833ff4544cc` |
| `tasks.md` | `475b01e9d8523ed21ced86680be28de55a4ee7200c26c543fc852c8306e9778a` |
| `apply-progress.md` | `6d438276a6c309ae61e62ddfcc9df2d3d2dc99cae309384c2d65612c7e85146d` |
| Full candidate manifest excluding this report | `15c8576cc6cd497419f8a5f4d1c7dd7572c727799025f621cfb3a85bf66f2870` |

Progress's opaque production hash `a363aa6bc71170177655170ca1458d52f2d71f04a18cf737980d806ebd189a81` is historical, with no manifest serialization supplied; it is not represented as reproduced under the definition above. Independent complete original-Go manifest equality and individual HEAD byte comparison establish unchanged production/existing tests without relying on that opaque value. F3 raw history requires the original local session store; it is not repository-portable proof.

Admission: `gentle-ai sdd-verify-validate --input - --requirements 13 --scenarios 38` receives these exact candidate bytes before any report write. Persistence occurs only after exit-0 native admission; readback checks exact bytes, historical appendix hash and protected manifests. No fresh child attempt, correction settlement, delivery, checkout/reset/cleanup, memory search, sync or archive occurs. MCP/memory unavailable; this OpenSpec report is the only persistence claimed. `skill_resolution: injected` (compact parent standards plus explicitly requested Go-testing skill and assigned verify/strict phase modules).

**Next recommended:** parent readback and existing-attempt handling; **sync may be recommended on this passing formal report, but is not executed here**. Parent owns active correction settlement. W1–W3, local-history portability and coverage granularity remain disclosed; external snapshot authenticity/acquisition, consumer policy and lifecycle authority remain out of scope.

# Historical appendix — original formal FAIL, preserved verbatim

Everything below is the original 242-line report at SHA-256 `4d1fca4f02823705d8e023206f42899cf06a0ba016a9d2b93b1c40e89276fa2a`. Its envelope, findings and next-step instructions are historical. The current admitted envelope and F1–F3 dispositions above supersede them only for the corrected candidate.

```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:c8dfa42cf3f46e3f21b3165923ba1c299584740a18a943d7a0f2ce271b10f6bc
verdict: fail
blockers: 3
critical_findings: 3
requirements: 11/13
scenarios: 36/38
test_command: GOTOOLCHAIN=go1.25.10 go test ./... -count=1 -json
test_exit_code: 0
test_output_hash: sha256:52ca8788b18030f27bed1527443ef10f006f8f5944b97478dfb5140ca6525937
build_command: GOTOOLCHAIN=go1.25.10 go test ./... -count=1
build_exit_code: 0
build_output_hash: sha256:1fc855715234092ded952f79a7e1f417b3f9cda1f412c80a2919f8c81fbb994c
```

# Verification Report — Result V1

## Verdict and scope

**FAIL: complete formal verification is not established.** All executed runtime gates passed; the blockers are two required-scenario evidence gaps and incomplete locally available historical TDD evidence, not confirmed implementation failures. No source defect was confirmed and no failing runtime behavior was reproduced.

- Workspace: `/data/Projects/git-change-evidence-worktrees/result-v1-unit9-mutation-census`.
- Branch: `test/result-v1-unit9-mutation-census`; HEAD: `17854e7539ca571107a01b35907695a979f1295c`.
- Scope: entire approved Result construction/decode/accounting/provenance/report-binding contract; this phase writes this report only. No corrections, physical mutations, lifecycle, delivery, sync, archive, or consumer integration were performed.
- Read: root AGENTS, product charter, `openspec/config.yaml`, complete proposal, four specs, design, tasks, cumulative apply-progress, all 40 changed Go files relative to design baseline `a2b9654`, and relevant unchanged snapshot/policy/inventory/report/CLI dependencies and tests.
- Counts from actual headings: contract 4 requirements/12 scenarios; accounting 5/14; Result provenance 3/8; Evidence provenance 1/4; total **13/38**. PASS below means covering runtime tests plus inspected source; it is not a claim of exhaustive inputs or exclusively public coverage.
- `evidence_revision` hashes the sorted tracked-file manifest excluding this report: each line is lowercase SHA-256, two spaces, relative path, LF. It includes the pre-existing dirty tasks/progress bytes, not merely HEAD.

## Findings and exact blockers

| ID | Severity / causality | Evidence and disposition |
|---|---|---|
| F1 | CRITICAL — required scenario PARTIAL/UNTESTED, not a runtime failure | Accounting / “Rename classification uses the committed primary path”: the only Result test supplying `PreviousPath` is `result_v1_snapshot_test.go:12`, using `/old` and calling private `validateResultSnapshot`. It verifies current-path extraction but supplies no policy with distinct old/new categories and never asserts the resulting classification, entry cardinality, or totals. Public constructor/compatibility fixtures supply no previous path. Source `result_v1_antecedents.go:45-55` correctly drops metadata, and the shared classifier is tested separately; neither establishes the complete approved rename premise at runtime. Complete scenario proof is missing. |
| F2 | CRITICAL — required scenario PARTIAL/UNTESTED, not a runtime failure | Result provenance / “Byte-identical antecedents reproduce provenance”: `TestValidateResultAntecedentsRequiresExactPolicyInventoryChain` reconstructs a chain by strict decoding cached bytes and compares private validation outputs; it never constructs and compares two Results from independently reconstructed equal chains. `TestNewAccountingResultV1ConstructsExclusiveCanonicalEvidence` compares reordered snapshots using the same antecedent values; the independent public decoder chain test deliberately uses a different policy. None supplies the specified two equal independent antecedent pairs and compares Result links/bytes/digest. Deterministic source is supportive, not replacement runtime evidence. |
| F3 | CRITICAL — incomplete available strict-TDD evidence | Unit 1's historical RED/GREEN/TRIANGULATE/REFACTOR and pre-change focused safety-net evidence is not included in the supplied artifact set. Apply-progress cites merged PR57 and external observation `14046`, but provides no Unit 1 cycle table or its commands/results. Read-only `git show --stat 31bcc06d9147adae1c830040d494ab5bdba5a659` confirms four Go changes; that tree has no Result apply-progress artifact. Current foundation tests pass. This is an evidence-availability blocker, **not** a conclusion that TDD was not performed. Parent must provide the referenced historical evidence for complete TDD verification; no memory retrieval or historical reconstruction was fabricated here. |
| W1 | WARNING — support/tooling qualification | `.pi/gentle-ai/support/strict-tdd-verify.md` is absent. Required TDD and assertion checks were performed using config/prompt and the assigned skill's strict module; strict mode was not relaxed. |
| W2 | WARNING — coverage granularity | Native 64-character IDs plus hostile bytes are exercised by the private snapshot validator, not by an end-to-end public Result fixture. The concrete path there is `ff 0a 0d 61`; `ff 0a 61` is separately checked in materialization. Wire/Base64 and public ownership tests cover the composing boundaries. This is composed coverage, not a directly executed public SHA-256/hostile-path golden. |
| W3 | WARNING — scoped residuals, not added requirements | Unit 8 still lacks direct valid-report/invalid-binding validator cases and invalid-subject bound-constructor coverage. Its builder invalid-chain tests and shared binding helper support source correctness; the validator's invalid-report case exits before binding validation. Preserve the disclosed nonblocking distinction. Deletion-only accumulation overflow and legacy required-line-total-overflow preservation also lack direct dedicated vectors; checked addition and unchanged legacy observation source support them. |

**Exact blockers:** F1, F2, F3. No correction is authorized in this phase. A parent/user decision is required before any follow-up evidence or tests; this report makes no delivery/archive-readiness claim.

## Scenario traceability

All test names in the registry below were observed passing in the 58-test anchored focus and again in the full suite. `P` = direct public API; `H` = private helper; `C` = composed boundary coverage (runtime tests of constituent behavior plus inspected wiring, not an unexecuted end-to-end claim). Tests are in-process Go tests, not external consumer integration.

### Source and runtime-test registry

| Key | Source | Passing tests (exact names; source file indicated) |
|---|---|---|
| CON | `result_v1_build.go`, `accounting_domain.go`, `result_v1_wire.go` | P: `result_v1_construct_test.go` — `TestNewAccountingResultV1ConstructsExclusiveCanonicalEvidence` |
| ANT | `result_v1_antecedents.go`, strict Policy/Inventory decoders | H: `result_v1_antecedents_test.go` — `TestValidateResultAntecedentsRequiresExactPolicyInventoryChain` |
| INV | `result_v1_build.go`, `result_v1_antecedents.go` | P: `result_v1_antecedents_test.go` — `TestNewAccountingResultV1InventoryOnlyChangesIdentity` |
| SNAP | `result_v1_antecedents.go:32-57`, `snapshot.go` | H: `result_v1_snapshot_test.go` — `TestValidateResultSnapshotAcceptsNativeIDsAndRejectsInvalidPaths` |
| OWN | `result_v1.go`, `result_v1_build.go` | P: `result_v1_ownership_test.go` — `TestNewAccountingResultV1OwnsIngressAndEgress`, `TestNewAccountingResultV1ReturnsZeroOnFatalFailure` |
| WIRE | `result_v1_wire.go`, `result_v1_wire_observations.go` | H: `result_v1_wire_test.go` — `TestResultV1CanonicalWireHostilePathAndIdentity`, `TestResultV1CanonicalWireUsesEmptyArrays`; `result_v1_wire_identity_test.go` — `TestResultV1CanonicalWireUsesClosedVariants`, `TestResultV1CanonicalWirePreservesPathByteIdentity` |
| OBS | `result_v1_observations.go`, `policy_primitives.go` | H: `result_v1_observations_test.go` — `TestResultV1ObservationsUseClosedOrderAndUnavailableVariants`, `TestResultV1StrictGreaterThanExcludesEqualityBoundary`, `TestResultV1RejectsRequiredLineTotalOverflow`, `TestResultV1NonCountableTotalsRemainUnavailableBeforeAddition` |
| PREC | `result_v1_observations.go:78-84`, unchanged legacy comparator | H + legacy P: `result_v1_observations_precision_test.go` — `TestResultV1ExactRatioPrecisionEdge`, `TestResultV1PrecisionEdgeIsTheOnlyAvailableRatioComparatorDifference` |
| DOMAIN | `accounting_domain.go`, `accounting_domain_totals.go`, `accounting.go` | H/P: `accounting_test.go` — `TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable`, `TestAccountingDomainRejectsCheckedAccumulationOverflow`, `TestResultAccountingDomainProjectsPolicyViewLosslessly`, `TestAccountSharedDomainTriangulatesLegacyBehavior`, `TestAccountUsesFirstMatchingRawByteGlobAndDefault`, `TestAccountRejectsInvalidPolicyBeforeAccounting`, `TestAccountReturnsEmptyResultOnOverflow`, `TestPathGlobGrammar` |
| COMP | shared domain, public builder, legacy Account | P: `result_v1_compatibility_test.go` — `TestResultV1AndAccountAgreeOnSharedDomainVectors` |
| SHAPE | all `result_v1_decode_*shape*.go`, `policy_primitives.go:decodePolicyObject` | H: corresponding entry, observation, revisions, totals, root shape test files — `TestResultV1EntryShapesRejectClosedShapesDuplicatesAndAuthorityKeys`, `TestResultV1ObservationShapesRejectClosedShapesDuplicatesAndAuthorityKeys`, `TestResultV1RevisionShapesRejectClosedShapesDuplicatesAndAuthorityKeys`, `TestResultV1TotalShapesRejectClosedShapesDuplicatesAndAuthorityKeys`, `TestResultV1RootEnvelopeShapeRejectsClosedRoot` |
| ROOT | `result_v1_decode_root_shape.go` and delegated subtrees | H composition: `result_v1_decode_root_integration_test.go` — `TestResultV1RootEnvelopeShapeDelegatesClosedSubtrees`, `TestResultV1RootEnvelopeShapeAcceptsStructuralSubtreeVariants` |
| VALUES | `result_v1_decode_values.go` | H/composed P: `result_v1_decode_values_test.go` — `TestResultV1DecodeValuesMaterializesOwnedCandidate`, `TestResultV1DecodeValuesRejectsLexicalAndStructuralCandidates`, `TestResultV1DecodeValuesAcceptsLexicallyCanonicalEmptyPath`, `TestResultV1DecodeValuesDefersSnapshotSemanticsToConstructor` |
| DEC | `result_v1_decode.go` and builder | P: `result_v1_decode_semantics_test.go` — `TestDecodeAccountingResultV1RoundTripsOwnedSelfConsistentEvidence`, `TestDecodeAccountingResultV1RejectsSemanticAndOrderMutations`, `TestDecodeAccountingResultV1ReturnsZeroOnFatalBoundaries` |
| CAN | `result_v1_decode_values.go`, `result_v1_decode.go:16-24` | P: `result_v1_decode_canonical_test.go` — `TestDecodeAccountingResultV1RejectsNoncanonicalRepresentations`, `TestDecodeAccountingResultV1RequiresExactAntecedentsAndAcceptsIndependentChain`, `TestDecodeAccountingResultV1AcceptsEmptyEqualRevisionEvidence` |
| BIND | `contract.go:71-120` | P: `result_v1_report_binding_test.go` — `TestResultV1ReportBindingDerivesAllProvenance`, `TestValidateReportV1ResultProvenanceRejectsEveryMismatchAndPreservesLegacySurfaces`, `TestNewReportV1FromResultReturnsZeroEvidenceForInvalidChains` |
| REPORT | unchanged `NewReportV1`, `DecodeCanonical`, `report.go` | P: `result_v1_report_legacy_test.go` — `TestResultV1ReportBindingPreservesLegacyStandaloneReports`; `contract_test.go` — `TestNewReportV1CanonicalBytesAndDigest`; `decode_test.go` — `TestDecodeCanonicalRejectsClosedAndAuthorityBearingDocuments`, `TestNewReportV1RejectsNonSHA256ProvenanceDigest`; `report_test.go` — `TestProjectEvidenceCanonicalJSON`, `TestProjectEvidenceHumanText`, `TestProjectEvidenceRejectsInvalidEvidenceAndUnsupportedFormat`, `TestProjectEvidenceOwnsCanonicalJSONOutput` |

### result-v1-contract — 4/4 requirements, 12/12 scenarios covered

| Requirement | Exact scenario | Trace / result |
|---|---|---|
| Neutral Result V1 construction and inspection | Construct a neutral result from exact typed evidence | CON + OWN + DEC: PASS (P); literal public canonical golden, digest, links, revisions and views. |
| Same | Caller offers only digest assertions | ANT + OWN + DEC: PASS (C); private digest-only forgeries rejected, public zero antecedents return zero, public signatures require typed documents rather than legacy policy/digest inputs. No negative-compilation test claimed. |
| Same | Failed construction has no partial result | OWN + DEC: PASS (P); antecedent, malformed snapshot, late path, stored/required-sum overflow return exact zero. |
| Closed Result V1 wire contract and content identity | Canonical hostile-path fixture | WIRE + CON: PASS (C); exact `/wpmaWxl` golden is private encoder coverage with synthetic links; public golden separately verifies real antecedents and hostile-byte encoding. |
| Same | Empty committed change set | INV + WIRE + CAN: PASS (C/P); distinct-revision empty public snapshot with nonempty inventory, policy-order zero totals, arrays and derived identity. |
| Same | Alternate wire representation | CAN + SHAPE + VALUES: PASS (P/H); order, whitespace, LF, escaped category, null array, Base64, fractional/exponent/string numbers rejected. |
| Defensive ownership and deterministic construction | Caller mutates a source or returned buffer | OWN + DEC + VALUES: PASS (P/H); effective full-buffer alternates and independently copied oracles. |
| Same | Equivalent snapshot orders | CON: PASS (P); raw prefix ordering and identical canonical bytes/digest. |
| Strict antecedent-aware canonical decoding | Wrong exact antecedent despite valid syntax | CAN + ANT + DEC: PASS (P/H); valid alternate and cross-policy chains, no digest-syntax shortcut. |
| Same | Recomputed evidence differs | DEC + SHAPE: PASS (P/H); classification/measurement/counts/totals/threshold/ratio/order mutations; forbidden unavailable actual is private shape coverage. |
| Same | Unknown, duplicate, and authority-bearing nested fields | SHAPE + ROOT + VALUES + DEC: PASS (C); nested duplicate/escaped/unknown/authority paths run through private validators/root/materializer; public zero-result materializer funnel separately passed. Not all nested rows execute the public decoder. |
| Same | Policy-derived authority-like string remains evidence | CAN: PASS (P), valid `release` policy construction/decode; `merge` also has H shape coverage. |

### result-v1-accounting-evidence — 4/5 requirements, 13/14 scenarios covered

| Requirement | Exact scenario | Trace / result |
|---|---|---|
| Valid immutable committed-snapshot boundary | Valid native SHA-256 snapshot with hostile path bytes | SNAP + VALUES + WIRE + CON: PASS (C), qualified by W2; private native-ID/raw-byte boundary and encoding/public wiring, not a public exact-fixture run. |
| Same | Invalid revision or path input | SNAP + VALUES + OWN + DEC: PASS (H/P); both widths, malformed/mixed/equal-with-entry IDs, invalid/duplicate raw paths; public fatal funnels. |
| Same | Rename classification uses the committed primary path | SNAP only extracts metadata-independent current path; no distinct-category rename accounting oracle: **GAP F1**. |
| Inventory is provenance-only and excluded from accounting | Inventory identity changes without changing accounting | INV: PASS (P); empty versus nonempty linked inventory, identical entries/totals/observations and distinct links/digests. |
| Same | Nonempty inventory with empty committed snapshot | INV: PASS (P); no entries, explicit policy-category zero totals. |
| Exclusive classification and recomputable policy-category totals | Overlap, default, and non-countable evidence | DOMAIN + CON + COMP: PASS (C/P); private exact four-path vector and public equivalent overlap/default/non-countable vector; literal exclusive classifications/totals. |
| Same | Constructor input order does not establish result order | CON: PASS (P), exact `a`, `a/ff`, `a0`, `ff` order. |
| Same | Realizable arithmetic overflow is fatal rather than invented evidence | DOMAIN + OBS + OWN + DEC: PASS (H/P); MaxUint64+1 stored additions, required threshold/numerator/reference sums, late overflow and exact-zero public results. |
| Closed factual threshold, ratio, and unavailable-observation model | Available threshold and exact ratio | OBS + COMP + DEC: PASS (H/P); 2+1 > 2, 3/4 > 0.5, exact fields/limits. |
| Same | Exact decimal comparison intentionally differs at a legacy float64 precision edge | PREC: PASS (H + legacy P); exact 9007199254740993/18014398509481984 > 0.5 versus legacy false. The fixture's symbolic snapshot IDs are used only by legacy Account/private domain, not a valid public Result construction. |
| Same | Non-countable and zero-denominator observations remain embedded | OBS + CON + COMP: PASS (H/P); unavailable false variants and available-zero distinction, no invented line values. |
| Same | Exceeded evidence is not authority | CON + OBS + DEC + WIRE: PASS (C/P); literal closed evidence and strict decoding, no authority result field. |
| Legacy accounting compatibility | Legacy caller remains unaffected | DOMAIN + unchanged-source comparison: PASS (P); no Policy/Inventory requirement added to Account, legacy validation/order/float arithmetic retained. |
| Same | Result and legacy accounting agree on shared-domain vectors | COMP + PREC + DOMAIN: PASS (P/H); literal expected entries plus differential totals/availability/ordinary ratios and accumulation failure. |

### result-v1-provenance — 2/3 requirements, 7/8 scenarios covered

| Requirement | Exact scenario | Trace / result |
|---|---|---|
| Exact validated antecedent chain | Byte-identical antecedents reproduce provenance | ANT proves fresh decoded validation outputs; CON proves same-chain order convergence, but neither compares Results from two equal independent chains: **GAP F2**. |
| Same | Inventory policy link is wrong but all digests are well formed | ANT + OWN + CAN: PASS (H/P); isolated supplied cached-link mismatch, valid cross-policy inventory rejection, public zero results. |
| Same | Same accounting with a different antecedent is not interchangeable | INV + CAN: PASS (P); different valid inventory identity leaves accounting unchanged and cannot decode original Result with alternate inventory. |
| Result identity is system-derived and complete | One raw path byte changes the identity | WIRE + CON: PASS (C); private ff/fe, LF/CRLF and padding vectors distinguish bytes/digests; public golden binds complete wire. |
| Same | Reordered candidate is not repaired | DEC: PASS (P); entry/total/observation swaps reject with exact zero Result. |
| Result V1 remains neutral evidence-only provenance | Authority-bearing candidate provenance | SHAPE + ROOT + VALUES: PASS (C); closed normalized metadata/discriminator rejection at nested boundaries; unknown CNSIC/project keys remain forbidden by closure, not by a new vocabulary policy. |
| Same | Policy-derived authority-like value remains evidence | CAN + SHAPE: PASS (P/H); exact valid `release` values accepted, structural strings not misrepresented as canonical Policy limits. |
| Same | Out-of-scope integration request | Public signatures + CON/DEC/BIND and inspected dependencies: PASS boundary coverage; pure typed calls expose no acquisition/delivery request API. No external integration or Git-authenticity proof is claimed. |

### evidence-v1-provenance — 1/1 requirement, 4/4 scenarios covered

| Requirement | Exact scenario | Trace / result |
|---|---|---|
| Provenance and reproducibility | Bound report derives all provenance | BIND + REPORT: PASS (P); fixed P/I/R digest literals, five derived values, identical legacy report bytes/digest, no Result body. |
| Same | Well-formed but mismatched accounting digest | BIND: PASS (P); independent valid alternate chain/report supplies the different accounting digest; isolated validator error. All five provenance fields separately tested. |
| Same | Standalone report bytes are not provenance proof | REPORT: PASS (P); standalone syntactic decode succeeds, bound validation fails against unrelated Result. |
| Same | No projection or authority expansion | BIND + REPORT + CLI focus: PASS (C); binding fixture has exceeded evidence and produces identical unchanged legacy bytes; old JSON/human projections independently pass. No direct bound-report projection test claimed. |

## Correctness, compatibility and design coherence

- **Antecedents:** strict Policy decode plus SHA-256 of canonical bytes, strict linked Inventory decode plus supplied digest/link equality; construction uses the freshly decoded policy view. Inventory entries never enter accounting.
- **Accounting:** one shared first-match raw-byte classifier and checked uint64 accumulator; non-countable entries ignore supplied lines; zero categories retain policy order. Result required line sums fail atomically; non-countable availability is checked before its otherwise unnecessary sum. Both ratio sides are evaluated, including reference overflow when numerator is unavailable.
- **Exact ratios:** big.Int cross multiplication preserves canonical Policy decimal digits/scale and strict `>`; no float64 conversion on Result's path. Precision/equality vectors prove the load-bearing differences, not a theorem from the test name “Only…”. Legacy float64 observation code is unchanged.
- **Canonical decode:** duplicate-aware raw-child parsing precedes typed extraction; exact schema/digests/Base64, candidate snapshot validation, fresh builder recomputation, separate links and whole-byte equality. No repair or partial return. Private structural acceptance of arbitrary strings/zero denominators/unavailable exceeded=true is intentional; canonical rebuild rejects invalid complete evidence.
- **Report binding:** approved minimal canonical-byte authority is respected. Both APIs freshly decode Result bytes; report validator decodes report bytes before comparing all five fields. Deliberately corrupted private cached digest/link/revision views do not displace canonical authority. Mandatory reflection/all-cached-view integrity is **not** imposed; the broad earlier identity wording is read under the supplied approved binding clarification.
- **Design deviation:** design's antecedent-first decoder pseudocode differs from shipped materialize-first order. The approved later error boundary establishes no global multi-error precedence. Inspected source and passing isolated error funnels show no behavioral contradiction under that clarification.
- **Legacy preservation:** `git diff a2b9654 HEAD -- accounting.go contract.go` shows shared-domain factoring and additive binding APIs only. No Go changes under `cmd/` or `internal/`, Policy/Inventory/snapshot, standalone report decoder or projections. Full regressions cover those unchanged packages. The pre-existing report accepts mixed-width native revision pairs; Result intentionally requires equal widths. No compatibility requirement is extended beyond its approved domain.
- **No-authority/no-effects:** Result production imports are pure standard-library helpers; no Git, filesystem, external configuration, publication, or CNSIC effects. Existing shell tests use local fixtures; they are not proof of an externally asserted snapshot's authenticity.

## Execution evidence

Every Go invocation used `GOTOOLCHAIN=go1.25.10`; `go version` reported `go1.25.10 linux/amd64` (exit 0). `GOTOOLCHAIN=go1.25.10 go env GOROOT` selected that toolchain's `bin` first on PATH for the separate `gofmt` executable (exit 0). Output hashes below are SHA-256 of exact combined stdout/stderr bytes captured in memory.

| ID | Exact command | Exit / result / output SHA-256 |
|---|---|---|
| T1 | Anchored expanded command immediately below | 0; all 58 selected real top-level names observed, none missing; 684 passing subtest events; `38cc9d292d01e08b68695a4d4c0b6f1bc9e3698728a7d3ccb340387a9901f5c6` |
| B1 | `GOTOOLCHAIN=go1.25.10 go test ./... -count=1` | 0; independent configured build/primary invocation, all five packages; `1fc855715234092ded952f79a7e1f417b3f9cda1f412c80a2919f8c81fbb994c` |
| T2 | `GOTOOLCHAIN=go1.25.10 go test ./cmd/git-change-evidence -count=1 -run '^(TestRunProjectCanonicalJSON|TestRunProjectHumanTextAndRepeatedInput)$' -json` | 0; exactly 2 top-level tests; `d331999913d812ec1409e45dec827f94a4794c3c4326bb8607b92df64edddea8` |
| T3 | `GOTOOLCHAIN=go1.25.10 go test ./... -count=1 -json` | 0; full fresh census below; `52ca8788b18030f27bed1527443ef10f006f8f5944b97478dfb5140ca6525937` |
| C1 | `GOTOOLCHAIN=go1.25.10 go test -cover ./...` | 0; informational (root/publication cached); `2527c2c152fbb4f3923df1e44ed0fbc7dcdc26b6688c52e39a6a2c70e467b122` |
| FMT | `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` | 0, empty output; `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| DIFF | `git diff --check` | 0, empty output; same empty-output hash as FMT. |
| STAT | `git diff --stat` | 0; pre-existing 46 additions/1 deletion in tasks/progress; `1cecc93335fef44a039532790b10e017f57997392af66a7ebef83d10534a98d8` |
| NUM | `git diff --numstat` | 0; progress 45/0, tasks 1/1; `ced5513d6e797a2d6cb60f366a56eea76a7f0fc5d087834f4b5906f7cb635a0a` |

T1 (test names extracted from the actual selected co-located files before execution; the JSON pass-name set was compared with all 58 selected names):

```bash
GOTOOLCHAIN=go1.25.10 go test . -count=1 -run '^(TestAccountUsesFirstMatchingRawByteGlobAndDefault|TestAccountRejectsInvalidPolicyBeforeAccounting|TestAccountReturnsEmptyResultOnOverflow|TestPathGlobGrammar|TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable|TestAccountingDomainRejectsCheckedAccumulationOverflow|TestAccountSharedDomainTriangulatesLegacyBehavior|TestResultAccountingDomainProjectsPolicyViewLosslessly|TestNewReportV1CanonicalBytesAndDigest|TestDecodeCanonicalRejectsClosedAndAuthorityBearingDocuments|TestNewReportV1RejectsNonSHA256ProvenanceDigest|TestProjectEvidenceCanonicalJSON|TestProjectEvidenceHumanText|TestProjectEvidenceRejectsInvalidEvidenceAndUnsupportedFormat|TestProjectEvidenceOwnsCanonicalJSONOutput|TestValidateResultAntecedentsRequiresExactPolicyInventoryChain|TestNewAccountingResultV1InventoryOnlyChangesIdentity|TestResultV1AndAccountAgreeOnSharedDomainVectors|TestNewAccountingResultV1ConstructsExclusiveCanonicalEvidence|TestDecodeAccountingResultV1RejectsNoncanonicalRepresentations|TestDecodeAccountingResultV1RequiresExactAntecedentsAndAcceptsIndependentChain|TestDecodeAccountingResultV1AcceptsEmptyEqualRevisionEvidence|TestResultV1ObservationShapesRejectClosedShapesDuplicatesAndAuthorityKeys|TestResultV1RevisionShapesRejectClosedShapesDuplicatesAndAuthorityKeys|TestResultV1RootEnvelopeShapeDelegatesClosedSubtrees|TestResultV1RootEnvelopeShapeAcceptsStructuralSubtreeVariants|TestResultV1RootEnvelopeShapeRejectsClosedRoot|TestDecodeAccountingResultV1RoundTripsOwnedSelfConsistentEvidence|TestDecodeAccountingResultV1RejectsSemanticAndOrderMutations|TestDecodeAccountingResultV1ReturnsZeroOnFatalBoundaries|TestResultV1EntryShapesRejectClosedShapesDuplicatesAndAuthorityKeys|TestResultV1TotalShapesRejectClosedShapesDuplicatesAndAuthorityKeys|TestResultV1DecodeValuesMaterializesOwnedCandidate|TestResultV1DecodeValuesRejectsLexicalAndStructuralCandidates|TestResultV1DecodeValuesAcceptsLexicallyCanonicalEmptyPath|TestResultV1DecodeValuesDefersSnapshotSemanticsToConstructor|TestResultV1ExactRatioPrecisionEdge|TestResultV1PrecisionEdgeIsTheOnlyAvailableRatioComparatorDifference|TestResultV1ObservationsUseClosedOrderAndUnavailableVariants|TestResultV1StrictGreaterThanExcludesEqualityBoundary|TestResultV1RejectsRequiredLineTotalOverflow|TestResultV1NonCountableTotalsRemainUnavailableBeforeAddition|TestNewAccountingResultV1OwnsIngressAndEgress|TestNewAccountingResultV1ReturnsZeroOnFatalFailure|TestResultV1ReportBindingDerivesAllProvenance|TestValidateReportV1ResultProvenanceRejectsEveryMismatchAndPreservesLegacySurfaces|TestNewReportV1FromResultReturnsZeroEvidenceForInvalidChains|TestResultV1ReportBindingPreservesLegacyStandaloneReports|TestValidateResultSnapshotAcceptsNativeIDsAndRejectsInvalidPaths|TestResultV1ModelZeroValue|TestResultV1ModelPopulatedAccessors|TestResultV1ModelDefensiveCopiesPreserveState|TestResultV1ModelReturnsEmptyViewsForPopulatedDocument|TestResultV1ModelDiscriminators|TestResultV1CanonicalWireUsesClosedVariants|TestResultV1CanonicalWirePreservesPathByteIdentity|TestResultV1CanonicalWireHostilePathAndIdentity|TestResultV1CanonicalWireUsesEmptyArrays)$' -json
```

| Package | Passing top-level tests | Passing subtest events | Statement coverage |
|---|---:|---:|---:|
| root | 87 | 888 | 94.9% |
| cmd/git-change-evidence | 13 | 13 | 83.5% |
| internal/git | 13 | 23 | 85.2% |
| internal/inventory | 4 | 4 | 75.8% |
| internal/publication | 5 | 0 | 79.3% |
| **Total** | **122** | **928** | No unweighted aggregate asserted |

No failing or skipped test events occurred. Subtest events can include intermediate `t.Run` nodes and are **not** assertion counts or additional top-level tests. Coverage threshold is 0. Package statement coverage is informational; no per-file/branch profile was written because only this report is an allowed artifact. No new race-bearing behavior exists, so `go test -race ./...` is N/A, not an executed PASS. No additional packages, installs, Docker, external consumer tests or CNSIC access were used.

Read-only tooling: MCP CodeGraph lookup returned “MCP not initialized” (no process exit code); upstream `codegraph status` and `codegraph explore ResultV1` both succeeded (exit 0), index up to date. No init/index/sync was run. A directed source-read loop exited 1 because two guessed dependency filenames did not exist; the actual `policy_primitives.go`/`policy_decode.go` were subsequently located with `git ls-files '*policy*go' '*inventory*go'` and read successfully (exit 0). The exact failed loop was:

```bash
for f in result_v1.go result_v1_antecedents.go accounting_domain.go accounting_domain_totals.go accounting.go result_v1_observations.go result_v1_wire_observations.go result_v1_decode_values.go result_v1_decode*shape*.go contract.go snapshot.go policy_object.go policy_decimal.go; do case "$f" in *_test.go) continue;; esac; echo "===== $f"; cat "$f"; done
```

Other inspected read-only commands succeeded: `git rev-parse --show-toplevel`, `git rev-parse HEAD`, `git status --short`, `git branch --show-current`, `git diff HEAD -- '*.go'`, `git diff --numstat a2b9654 HEAD -- '*.go'`, `git diff a2b9654 HEAD -- accounting.go contract.go`, `git show --stat --oneline 31bcc06d9147adae1c830040d494ab5bdba5a659`, `git ls-tree -r --name-only 31bcc06d9147adae1c830040d494ab5bdba5a659 openspec/changes/result-v1`, and `gentle-ai sdd-verify-validate --help` (all exit 0). Missing support-file existence checks are disclosed above, not represented as test failures.

## Strict TDD and assertion-quality assessment

Config's active authoritative `go-source-or-test` rubric requires RED → GREEN → TRIANGULATE → REFACTOR, focused tests before production and after refactor, full Go tests and check-only formatting. This verify-only phase introduces no implementation cycle.

| Check | Result |
|---|---|
| TDD Cycle Evidence tables exist | Yes: local stage/command/result tables for Units 2–8, including bounded structural/materialization/public slices and corrective records. Semantic stages are present despite not using the skill template's emoji-column layout. |
| Historical RED | Recorded undefined-symbol failures and subsequent test-oracle/compile corrections are distinguished. Existence of tests is confirmed; past command execution was not independently witnessed. Unit 1's external reference is insufficient local evidence: F3. |
| GREEN now | Confirmed for every existing referenced test file through T1/T3; all 23 Go test files changed since the design baseline were read and executed. |
| TRIANGULATE | Actual nonempty tables cover numeric limits, availability/equality, exact precision, structural closure, independent mismatched chains and ownership controls. The mandatory scenario gaps F1/F2 remain; no table-size claim is substituted for those missing premises. |
| REFACTOR and safety net | Recorded focused reruns after format/readability work for Units 2–8; some historical runs explicitly cached. Current reruns are uncached. Modified legacy foundation safety-net chronology cannot be certified without the referenced Unit 1 record. |
| Unit 0 / Unit 9 | No production TDD cycle expected: Unit 0 is a baseline/layout gate whose approval is cited; Unit 9 records temporary-mutation baseline/kill/restore, not newly authored behavior. |

Task status: **12/12 implementation checkbox entries checked**, consisting of Units 0–9 plus two Unit 2 subcheckboxes; no unchecked implementation task. Proposal success checkboxes and historical “remaining work” paragraphs are planning/history, not the authoritative current task census. Implementation checkbox completion is not equivalent to all scenarios being runtime-proven. For the eight Go-producing units (1–8), seven have locally available stage evidence; Unit 1 is reference-only. Formal strict-TDD verification is therefore incomplete, even though current GREEN is true.

Assertion audit: all 23 changed test files (51 top-level test functions, including pre-existing functions in modified `accounting_test.go`) were read. No confirmed tautology, vacuous assertion loop, type-only behavioral oracle, smoke-only behavioral suite, CSS assertion, or mock-heavy substitute was found. Iteration fixtures are literal/nonempty or guarded by expected counts; equality comparisons use behavioral values. The observation closure helper file has no top-level test of its own but is called by the observation-shape top-level test and executed. The model discriminator test checks production constants against literal contract strings; it is supplemental constant-contract coverage, not independent construction/provenance proof. Fake model documents and synthetic wire links were not counted as validated public antecedents.

The Unit 7A ownership correction overwrites the same raw backing buffer with a distinct equal-length valid alternate, verifies both alternate links/revisions/path/counts, then verifies original retained values. Unit 7B separately requires a successful alternate decode with changed source additions (index two), actual CR/LF in decoded JSON plus equivalent StdEncoding output, and isolated snapshot/overflow `ContractError` values with zero Result. These corrected premises executed now; the historical premise-kill runs did not.

Layer census: the 51 functions in 23 changed files are Go in-process core unit/component tests (private helpers plus composed public contracts), not network/browser E2E tests. CLI's 2 focused projection tests are in-process shell/core integration. The full suite also executes pre-existing local Git/filesystem adapter tests; its 122 total is not a Result-only unit count. No extra linter/standalone static analyzer was configured or run; Go compilation/test validation passed. Coverage reports only package statements, not changed-file lines or branches.

## Unit 9 mutation record versus independent observation

**Recorded: 11/11 assertion kills. Independently rerun in this phase: 0 mutants.** Independently observed now: restored-source GREEN killers, relevant assertions/premises, six matching restored file hashes, all-Go equality to HEAD. Apply-progress contains descriptions/commands/failure names but not raw per-mutant output transcripts; reported kills are retained as historical evidence, not newly witnessed execution or exhaustive mutation adequacy.

| Recorded mutant | Independently inspected killer premise | Scope limitation |
|---|---|---|
| ORD-1 sort reversal | CON literal four-entry raw order and canonical golden | Public constructor killer, not a new decode-order mutation run. |
| PATH-1 rune normalization | CON expects raw ff paths/Base64 | Distinguishes invalid-UTF-8 replacement bytes. |
| B64-1 URL encoding | WIRE exact `/wpmaWxl` | Private encoder; no claim this kills every decoder Base64 weakness. |
| PROV-1 supplied Inventory cached policy-link bypass | ANT valid bytes/correct digest with only supplied cached link changed | Kills that specific additional predicate, not all policy/inventory/result-link bypasses. |
| CLASS-1 last-match | DOMAIN category 0 versus 1 on overlapping ff glob | First-match predicate, not a separate inventory-contribution mutant. |
| OVF-1 checked-add bypass | DOMAIN MaxUint64+1; OBS required-sum overflows | Does not separately mutation-prove every accumulator call site. |
| RATIO-1 float comparator | PREC exact 2^53 edge with expected true/legacy false | Private Result derivation and legacy differential. |
| GT-LINE-1 strict threshold to >= | OBS available equal actual/limit expects not exceeded | Equality premise inspected. |
| GT-RATIO-1 strict ratio to >= | OBS available equal rational/limit expects not exceeded | Equality premise inspected. |
| DIG-LF-1 omit terminal LF | WIRE exact golden includes LF | Byte oracle, not a self-rebuild-only comparison. |
| DIG-HASH-1 hash body without LF | WIRE hashes independently shaped expected bytes | Independent expected SHA-256, not hashing actual bytes as the sole oracle. |

All six recorded restoration SHA-256 values matched current bytes:

```text
4db67187203056b083aa3dd35b02bc2d6060b334e894d711f448ccf6d48d0f0c  result_v1_build.go
bb4df4247338f40fe92384c7dacbaf51bbbf8f0fab3b2ef720041c89dbbb9ac5  result_v1_antecedents.go
d35b274ae19fd32603b822c9e0886a1a45ff94f0caced45ede443a6aa67fbc0d  result_v1_wire.go
7a6ef6e613dfaba388e9f7d08031c78704852752dc110f6db1d50cf78d2a921d  accounting_domain.go
15e7ac7dc54437cd692cc5fa23e254ffd76a9ea0a6132e5599b50a8d01f9be48  accounting_domain_totals.go
75c8a4b816778d545b75eff4bb419671944dc7958fb27f9131117bc5f10995cf  result_v1_observations.go
```

## Review workload, integrity and handoff

- Selected chain is `stacked-to-main`; report-only verification respects that assigned boundary. No chained implementation slice was added or expanded. No `size:exception` selected or inferred; the owner's current limit is 400 **production** LOC with documentation/tests separate and a 650-total runtime ceiling.
- This phase: **0 production changes, 0 test changes, one new report**. Pre-existing Unit 9 tasks/progress remain 46 additions + 1 deletion = **47 total**, not new verification edits.
- Historical source census against design baseline `a2b9654`: 17 production files / 1,076 changed lines; 23 test files / 2,160 changed lines. This is a multi-slice net Go diff, **not** a PR-sized exception or the entire authored historical delivery total. Existing planning/evidence documents are excluded from that Go-only number.
- Root Go census is **67**, repository Go census **82**. Unit 9's “all 67 Go/test files” wording describes the root census, not all internal/CLI files. Current `git diff HEAD -- '*.go'` is empty across the whole repository. Historical layout threshold crossings/owner approvals are recorded per slice; this phase adds no Go/layout crossing.
- Before/after verification source/test manifest (all 82 Go files, sorted relative paths, SHA-256/two spaces/path/LF): `9e95f953e6ebe929d35ab2ff1fca0f91f8af794d6d937eacce8bb3916cdf7794`.
- `tasks.md` before/after SHA-256: `475b01e9d8523ed21ced86680be28de55a4ee7200c26c543fc852c8306e9778a`.
- `apply-progress.md` before/after SHA-256: `e48c0563c1a43f3bc46ede93449e61eb807aecea5719a8ca8cefe7ebb2dcc4d4`.
- Combined sorted tasks/progress manifest before/after: `8f1c1495eaf332d1593c15670f2486cc0c7ed6dc77bf9e30fc29eb6e4b022d5a`. All tracked evidence excluding this report matched the envelope evidence revision before persistence; final readback rechecks it.
- Admission command: `gentle-ai sdd-verify-validate --input - --requirements 13 --scenarios 38`; the exact candidate bytes are submitted before this file is written, and only an exit-0 admitted candidate is persisted. Failure evidence is intentional and does not authorize settlement/delivery.
- Skills: `skill_resolution: fallback-path`. The parent supplied the exact Go-testing skill path, but no compact Project Standards block; the assigned sdd-verify phase skill and strict module were loaded from their known skill path. Inject compact standards next time. No memory search or external memory persistence was performed; this OpenSpec report is the persisted artifact.

**Next recommended:** parent readback and evaluation of F1–F3; await an explicit user decision. No automatic apply, mutation replay, sync, archive, delivery, or phase-4 action. External Git object existence/authenticity, repository membership, moving-ref stability and consumer policy remain explicitly out of scope.
