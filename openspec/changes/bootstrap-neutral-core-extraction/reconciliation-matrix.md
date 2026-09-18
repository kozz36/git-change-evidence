# Bootstrap Neutral Core Extraction — Reconciliation Evidence Matrix

## Scope and methodology

- **Matrix timestamp:** 2026-09-18T09:21:08Z.
- **Workspace observed:** `ef425bd53594670da2ff57173ddcc78013e6ccb8`.
- **Native-status starting fact:** the supplied `gentle-ai.sdd-status` v2 fact says this bootstrap change is `ready/apply`; proposal, specs, design, tasks, and apply-progress are done; no bootstrap verify report exists; and the canonical task ledger is **3/12 complete**. This matrix records that supplied status as input and does not reconstruct or alter native lifecycle state.
- **Artifact method:** complete reads of the active proposal, design, tasks, apply-progress, exploration, all seven active change specs, `openspec/config.yaml`, and `AGENTS.md`; targeted reads of current source/tests; parent-reported initial CodeGraph absence followed by delegated explorer initialization and CodeGraph exploration before filesystem mapping; verified Git history for cited commits; and complete reads of relevant later Public Preview and Go AST Census artifacts.
- **Classification meaning:** `proven-present` means current implementation/documentary evidence covers the bounded task intent; `partial` means a material portion exists but a stated semantic proof is missing; `absent` means no implementation/handoff surface was found; `superseded-decision` means current governance supersedes the bootstrap decision framing while explicitly retaining the relevant decision as deferred. These are evidence classifications only.
- **Chronology rule:** current source/tests and later commits establish present bytes and executable behavior, not a retroactive RED→GREEN sequence for an unchecked bootstrap task. “Checkbox reconciliation supportable” therefore means whether the canonical checkbox could be retained or changed without inventing chronology; it does not authorize any checkbox edit.
- **Preservation:** `tasks.md` remains unchanged. Its original ledger is reported unchanged at **3/12** (`1.1`, `1.2`, and `2.1` checked; nine later items unchecked).

## Native status, current-governance, and residue census

| Surface | Observed fact | Provenance |
|---|---|---|
| Bootstrap lifecycle input | Supplied native STATUS v2 fact: `ready/apply`; active artifacts done; bootstrap verify report missing; ledger 3/12. | Task instruction; not regenerated from OpenSpec artifacts. |
| Current preview decision | `v0.1.0-preview.1` is a planned source/dev identity, not a tag, release, package, asset, distribution channel, or visibility change. Broad distribution, packaging, signing/provenance, release automation, and configuration syntax remain deferred. | `AGENTS.md` “Current public-preview governance” and “Deferred Decisions”; `openspec/config.yaml` `context`; `docs/product-charter.md` current-preview section. |
| CNSIC state | The immutable predecessor is `a0fc7b26ff8a0e0a61baa586b32c46841611c806`; it remains evidence, not a current profile implementation. | `AGENTS.md`; `compatibility_test.go`; `testdata/compatibility/v1/manifest.json`. |
| Baseline host residue | `.atl/skill-registry.md` modified, `.gitignore` modified, and `.pi/` untracked before matrix creation. They are separate host/tooling residue and are not candidate changes. | `git status --short` at matrix start. |
| CodeGraph residue | The parent’s initial CodeGraph call reported no `.codegraph/` index. The delegated explorer then initialized `.codegraph/`; it is explorer-generated non-candidate tooling residue, not baseline residue. | Parent initial CodeGraph result; subsequent delegated explorer initialization and post-exploration directory check. |
| CNSIC profile / handoff discovery | No current-tree profile/handoff source path was found; `find openspec -iname '*handoff*'` produced no path. CNSIC strings in compatibility fixtures and generated `.atl` metadata are not profile or handoff implementation. | `git ls-files` profile/handoff path census; `find` result; `compatibility_test.go`; `.atl/skill-registry.md` classified separately. |

## Classification census

| Classification | Count | Tasks |
|---|---:|---|
| `proven-present` | 4 | 1.1, 1.2, 2.1, 3.1 |
| `partial` | 4 | 2.2, 2.3, 3.2, 3.3 |
| `absent` | 3 | 3.4, 4.1, 4.2 |
| `superseded-decision` | 1 | 4.3 |
| **Total** | **12** | **Original ledger remains 3/12; this census does not change it.** |

## Evidence matrix

### 1.1 — Bootstrap boundary decision

| Field | Evidence |
|---|---|
| Original acceptance intent | “Record the minimal root `AGENTS.md`, Go 1.25.10 target, public root API/`cmd/git-change-evidence/`/`internal/` layout, co-located `*_test.go` tests, primary/race/format gates, immutable typed Go configuration boundary; defer module/publication identity, packaging, CI, release, license, distribution, and concrete syntax.” — `tasks.md` Phase 1, item 1.1. |
| Classification | `proven-present` |
| Provenance | `AGENTS.md` records architecture, Go 1.25.10, root/census public API, CLI/internal layout, co-located tests, and typed adapter-only configuration; `openspec/config.yaml` records the test commands and layout; `proposal.md`/`design.md` record the original boundary. Bootstrap decision commit `73e2e6f4c610c6c6d44acdb00925d4f52a536366` (`chore(project): establish Go-first bootstrap decisions (#2)`). |
| Commands and actual results | The configured full/race/format commands in the command table below passed in the current tree. Historical `apply-progress.md` records planning-only checks and N/A production TDD for this item. |
| Semantic gaps | The originally deferred module/license/CI identity decisions are now superseded by later current governance; this does not erase the historical bootstrap record. Concrete configuration syntax, broad distribution, packaging, signing/provenance, and release automation remain deferred. |
| Checkbox reconciliation supportable? | **Yes, retain checked.** `tasks.md` and `apply-progress.md` both record completion, and commit `73e2e6f` is verified. |

### 1.2 — Strict-TDD/testing rubric and zero baseline

| Field | Evidence |
|---|---|
| Original acceptance intent | “Refresh `openspec/config.yaml` as the authoritative Go strict-TDD/testing rubric; record exact commands, layout thresholds, and zero-file initial crossing baseline without fabricating Go source or tests.” — `tasks.md` Phase 1, item 1.2. |
| Classification | `proven-present` |
| Provenance | `openspec/config.yaml`: `strict_tdd: true`; exact primary, race, and check-only format commands; rubric rows; layout thresholds; `initial_file_crossing_baseline`. `apply-progress.md` records the planning-only validation. Commit `73e2e6f4c610c6c6d44acdb00925d4f52a536366` is verified. |
| Commands and actual results | Exact current configured primary, race, and format commands passed; results are recorded below. |
| Semantic gaps | The zero-file baseline is historical: the current repository now contains a module, implementation, and tests. The rubric itself remains current. |
| Checkbox reconciliation supportable? | **Yes, retain checked.** The checked item, historical progress evidence, configuration content, and bootstrap commit align. |

### 2.1 — Closed contracts, canonical identity, and predecessor corpus

| Field | Evidence |
|---|---|
| Original acceptance intent | “RED then GREEN closed/versioned contracts, canonical bytes/digest/provenance, authority rejection, and frozen vectors/manifest with origin identities.” — `tasks.md` Phase 2, item 2.1. |
| Classification | `proven-present` |
| Provenance | Root symbols `NewReportV1`, `NewReportV1FromResult`, `ValidateReportV1ResultProvenance`, `Evidence.CanonicalBytes`, and `Evidence.Digest` in `contract.go`; strict decoding in `decode.go`; tests `TestNewReportV1CanonicalBytesAndDigest`, `TestDecodeCanonicalRejectsClosedAndAuthorityBearingDocuments`, `TestW2FrozenFixtureUsesNeutralPublicChain`, and `TestW2SyntheticCopySeparatesLegacyThresholdFromNeutralAvailability`. Corpus provenance is in `testdata/compatibility/v1/{contracts-v1.json,accounting-v1.json,manifest.json}` and `compatibility_test.go`. Verified commit `9e8a4e7a1bb0bf9ae06de36d02b23eb2dffd18b9` (`feat: add neutral W1/W2 evidence contracts (#4)`). |
| Commands and actual results | `go test ./...` passed, covering the root contract tests. `apply-progress.md` separately records the original 2.1 RED/GREEN/correction evidence. |
| Semantic gaps | The corpus is successor-authored metadata binding predecessor identities; `apply-progress.md` explicitly records correction of the earlier overbroad frozen-vector wording. No gap in the bounded corrected task statement was found. |
| Checkbox reconciliation supportable? | **Yes, retain checked.** The checked item, historical TDD record, corpus provenance, and verified commit agree. |

### 2.2 — Immutable Git snapshot and separate inventory

| Field | Evidence |
|---|---|
| Original acceptance intent | “RED then GREEN immutable Git/inventory: SHA-1/256, hostile paths, rename/binary/symlink/gitlink, missing/racing objects, dirty/index isolation, argv/closed-env/replacement isolation, and no-follow races in real Git.” — `tasks.md` Phase 2, item 2.2. |
| Classification | `partial` |
| Provenance | `internal/git.Acquire`/`acquire`/`resolve`/`gitOutputBounded` in `internal/git/snapshot.go`; `AcquireBlob` in `internal/git/blob.go`; inventory `Acquire` in `internal/inventory`. `gitOutputBounded` supplies explicit argv and controlled Git environment; snapshot tests use real temporary Git repositories. Tests: `TestAcquirePreservesCommittedEvidence`, `TestAcquireReportsTypedFailuresWithoutSnapshot`, `TestAcquirePreservesSHA256Identity`, `TestAcquireBlobUsesCommittedObjectIdentity`, `TestAcquireBlobIgnoresIndexAndWorktree`, `TestAcquireBlobPreservesLiteralPathsAndSymlinks`, `TestAcquireBlobReportsTypedFailuresWithoutContent`, `TestAcquireBlobPreservesSHA256Identity`, `TestAcquireRejectsUnsafeOrUnavailablePathWithoutRecords`, `TestAcquireRejectsReplacementRaceWithoutRecords`, and `TestAcquireDoesNotChangeCommittedSnapshot`. Verified implementation commits: `8c77cd4e222298f34e0c18ac6edf9c82e1ebffec` (snapshot), `e38374421789657c261cfadf4c5484cfe0aecbdc` (inventory), and `d52583c38055cbc693bce90bfd6fe4bef17ab5db` (blob). |
| Commands and actual results | `go test -count=1 ./internal/git ./internal/inventory` passed: `internal/git` 0.160s; `internal/inventory` 0.019s. These tests create and use real temporary Git repositories where their test bodies require them. Full/race gates also passed. |
| Semantic gaps | The original acceptance requires traversal-race verification in real temporary Git repositories. `internal/inventory/inventory_unix_test.go` runs `TestAcquireRejectsUnsafeOrUnavailablePathWithoutRecords` and `TestAcquireRejectsReplacementRaceWithoutRecords` against plain temporary directories; its separate Git-backed `TestAcquireDoesNotChangeCommittedSnapshot` proves committed-snapshot isolation, not those traversal races. The original RED→GREEN chronology is also not recorded in the bootstrap `apply-progress.md` for this unchecked task. |
| Checkbox reconciliation supportable? | **No.** Current behavior and later commits support present code, but `tasks.md` remains unchecked and the bootstrap progress ledger contains no task-2.2 RED/GREEN chronology; changing it would manufacture chronology. |

### 2.3 — Injected exclusive policy accounting

| Field | Evidence |
|---|---|
| Original acceptance intent | “RED then GREEN injected policy: overlap/invalid/default classification, non-countable data, text totals, non-gating thresholds.” — `tasks.md` Phase 2, item 2.3. |
| Classification | `partial` |
| Provenance | `Account` and policy validation in `accounting.go`; shared classification in `accounting_domain.go`; versioned policy contract in `policy_contract.go`. Tests: `TestAccountUsesFirstMatchingRawByteGlobAndDefault`, `TestAccountRejectsInvalidPolicyBeforeAccounting`, `TestAccountReturnsEmptyResultOnOverflow`, `TestAccountingDomainPreservesFirstMatchDefaultAndNonCountable`, `TestAccountSharedDomainTriangulatesLegacyBehavior`, and `TestAccountingPolicyV1Contract`. Observations expose thresholds as data (`Available`/`Exceeded`) without a delivery action. Verified commits: `2363ecd388d7c9fb11daf02fff7ddf5740b1f19d` (accounting) and `31bcc06d9147adae1c830040d494ab5bdba5a659` (shared-domain refactor). |
| Commands and actual results | `go test -count=1 . -run 'Test(Account|Accounting|ResultAccounting)'` passed in 0.001s. Full/race gates passed. |
| Semantic gaps | The original specification requires production-equivalent fallback when no alternate default is supplied. Current accounting requires an explicit named default (`accounting.go:100`; `accounting_test.go:56`; `policy_contract.go:94`), and no profile bridge proves the original fallback behavior. The bootstrap record also does not preserve this unchecked task’s original RED/GREEN chronology. |
| Checkbox reconciliation supportable? | **No.** Present behavior is not enough to convert the unchecked historical task into a completed chronological record. |

### 3.1 — Bounded immutable carveout and positional edge

| Field | Evidence |
|---|---|
| Original acceptance intent | “RED then GREEN W3 bounded immutable carveout: invalid bound, unavailable, timeout, deterministic match, dirty-tree isolation; positional legacy edge only.” — `tasks.md` Phase 3, item 3.1. |
| Classification | `proven-present` |
| Provenance | `PrepareCarveout`/`invalidCarveoutBound` in `carveout.go`; `MatchPreparedCarveout` plus deadline and pair-cap handling in `matcher.go`; positional CLI `run`, `measure`, `prepared`, and `matched` in `cmd/git-change-evidence/main.go`; immutable blob acquisition in `internal/git/blob.go`. Tests: `TestPrepareCarveoutReportsInvalidBounds`, `TestPrepareCarveoutDistinguishesTypedStates`, `TestPrepareCarveoutHonorsLimitBoundaries`, `TestPrepareCarveoutReportsHostileInputResourceBoundsBeforeAllocation`, `TestMatchPreparedCarveoutMatchesDeterministically`, `TestMatchPreparedCarveoutDeadlineBoundaries`, `TestMatchPreparedCarveoutPreservesFuzzyAndCapBounds`, `TestRunUsesCommittedBlobsAndFormatsExactly`, and `TestMatchedSharesAbsoluteDeadline`. Verified commits: `5f8091792c3c562b40b22b9d29740bbbd60bd891` (prepared inputs), `41d144c81c6ade24ba6a274d2c4a157f3067121e` (bounded/cancelled blobs), and `3104231d8c8acb86f77f16d74acecab41184f2d1` (positional edge). |
| Commands and actual results | `go test -count=1 . -run 'Test(PrepareCarveout|MatchPreparedCarveout)' && go test -count=1 ./cmd/git-change-evidence -run 'Test(RunUsesCommittedBlobsAndFormatsExactly|MatchedSharesAbsoluteDeadline)'` passed: root 0.001s; CLI 0.018s. The CLI test constructs a real temporary Git repository. |
| Semantic gaps | No missing bounded acceptance behavior was found: unavailable and invalid inputs are typed, deadline expiration maps to timeout evidence, and the positional path uses committed blobs despite dirty/index content. The original task’s RED→GREEN chronology remains absent from the bootstrap ledger. |
| Checkbox reconciliation supportable? | **No.** The current behavior is present, but the unchecked task has no preserved bootstrap-task chronology sufficient for a canonical checkbox change. |

### 3.2 — Report/projection and diagnostic redaction

| Field | Evidence |
|---|---|
| Original acceptance intent | “RED then GREEN W4 report/projection and diagnostics redaction: names, secrets, environment, paths, command output never leak or imply authority.” — `tasks.md` Phase 3, item 3.2. |
| Classification | `partial` |
| Provenance | Report/provenance constructors: `NewReportV1`, `NewReportV1FromResult`, and `ValidateReportV1ResultProvenance` in `contract.go`. Projections: `ProjectEvidence` and `ProjectionError` in `report.go`. Tests: `TestProjectEvidenceCanonicalJSON`, `TestProjectEvidenceHumanText`, `TestProjectEvidenceRejectsInvalidEvidenceAndUnsupportedFormat`, `TestProjectEvidenceOwnsCanonicalJSONOutput`, `TestResultV1ReportBindingDerivesAllProvenance`, and `TestValidateReportV1ResultProvenanceRejectsEveryMismatchAndPreservesLegacySurfaces`. Verified report/projection commits: `3ff69ff72edb72b9452ec4f7541febe88980058e` and `4b8f53f3de67e0ab4cfeaaa31d7d2fe57df44833`. |
| Commands and actual results | `go test -count=1 . -run 'Test(ProjectEvidence|ResultV1ReportBinding|ValidateReportV1ResultProvenance)'` passed in 0.002s. Full/race gates passed. |
| Semantic gaps | Current human projection omits `Subject` (`TestRunProjectHumanTextAndRepeatedInput`), and errors are generic. However canonical report bytes retain caller `Subject` (`canonicalReport.Subject` in `contract.go`), and no dedicated test injects or proves redaction of names, secrets, environment values, repository paths, or command output across report/diagnostic surfaces. The stated redaction acceptance is therefore not fully evidenced. |
| Checkbox reconciliation supportable? | **No.** It is unchecked; the material semantic redaction gap and absent historical RED/GREEN record preclude a chronology-safe change. |

### 3.3 — CLI technical outcomes and exclusive publication

| Field | Evidence |
|---|---|
| Original acceptance intent | “RED then GREEN W4 CLI/publication: technical outcomes; destination/interruption/content identity/moving-ref race proofs; atomic exclusive writes.” — `tasks.md` Phase 3, item 3.3. |
| Classification | `partial` |
| Provenance | CLI `runPublish` and exit mapping in `cmd/git-change-evidence/main.go`; CLI tests `TestRunPublishSuccess`, `TestRunPublishMapsFailuresAndInput`, and `TestRunCommandDispatchesPublishBeforeCWD`. Linux publication adapter `publish`, `rootDirectory`, `shaDirectory`, `stage`, and `finalize` in `internal/publication/publication_linux.go`; tests `TestPublishWritesCanonicalEvidenceAtDigestIdentity`, `TestPublishRejectsExistingEntriesAndSymlinkRoots`, `TestPublishCancellationAndExclusivity`, and `TestSyscallSeams`. Verified commits: `3ff69ff72edb72b9452ec4f7541febe88980058e` (projection), `db8c5465b458f2d5e9c9270d237ea5f8b302c93a` (publication adapter), and `c3a279668da5b36c192402c772894717230285b0` (CLI wiring). |
| Commands and actual results | `go test -count=1 ./internal/publication ./cmd/git-change-evidence -run 'Test(Publish|RunPublish)'` passed: publication 0.001s; CLI 0.004s. Full/race gates passed. |
| Semantic gaps | Canonical content-addressed identity, interruption handling, conflict/non-overwrite behavior, descriptor-safe paths, and exclusive/atomic publication are evidenced. No current publication API accepts or re-resolves a symbolic Git ref, and no moving-ref revalidation/race test was found in the adapter or CLI. The original moving-ref proof is not currently represented. |
| Checkbox reconciliation supportable? | **No.** The task remains unchecked; the moving-ref semantic gap and absent task chronology prevent a supportable checkbox change. |

### 3.4 — Supplied forecast comparison

| Field | Evidence |
|---|---|
| Original acceptance intent | “RED then GREEN W5 supplied forecast: valid, missing, divergent, unavailable, threshold; never author, mutate, infer, or gate.” — `tasks.md` Phase 3, item 3.4. |
| Classification | `absent` |
| Provenance | Repository-wide Go symbol/test search for `Forecast` returned no current source or test match. The only forecast material located is immutable predecessor compatibility fixture content (`testdata/compatibility/v1/contracts-v1.json`) and historical planning artifacts; those are not a neutral forecast API. |
| Commands and actual results | `go test ./...` and `go test -race ./...` passed for existing packages; neither establishes an absent W5 surface. |
| Semantic gaps | No supplied-forecast input, comparison result, missing/divergent/unavailable handling, threshold observation, ownership protection, or non-authoring test/API exists in the current tree. |
| Checkbox reconciliation supportable? | **No.** The unchecked state is consistent with the observed absence; no completion chronology or behavior exists to reconcile. |

### 4.1 — Standalone CNSIC profile

| Field | Evidence |
|---|---|
| Original acceptance intent | “RED then GREEN standalone CNSIC profile: one-way vocabulary/authority rejection, distinct-profile neutrality, frozen-vector and real-Git parity.” — `tasks.md` Phase 4, item 4.1. |
| Classification | `absent` |
| Provenance | `git ls-files` profile/handoff discovery found no profile package or source. `AGENTS.md` and `docs/product-charter.md` identify CNSIC as a future first consumer/profile and retain the predecessor oracle. Existing `compatibility_test.go` and `testdata/compatibility/v1/*` bind a W1/W2 predecessor corpus, not a standalone CNSIC profile or real-Git profile-parity suite. |
| Commands and actual results | Full and race suites passed for current packages; no profile package/test was listed or executed because none was found. |
| Semantic gaps | No one-way CNSIC adapter/profile, authority-injection rejection at that edge, distinct-profile proof, or frozen-vector plus real-Git parity operation exists. |
| Checkbox reconciliation supportable? | **No.** The unchecked state matches the observed absence. |

### 4.2 — Separate CNSIC SDD handoff

| Field | Evidence |
|---|---|
| Original acceptance intent | “Create separate CNSIC SDD handoff: dependency integration, reversible cutover, owner observation period, local-code deletion; no CNSIC edit here.” — `tasks.md` Phase 4, item 4.2. |
| Classification | `absent` |
| Provenance | `find openspec -iname '*handoff*' -print` returned no handoff path. The active bootstrap proposal/design describe a later cutover concept only; `docs/product-charter.md` retains CNSIC profile/cutover as a later milestone. No current-tree `cnsic`/profile/handoff implementation path was found. |
| Commands and actual results | The discovery command returned no paths. Full and race suites passed for existing code but do not supply handoff evidence. |
| Semantic gaps | No separately scoped CNSIC handoff artifact, dependency integration record, reversible cutover plan, owner observation period, or local-code deletion plan is present. |
| Checkbox reconciliation supportable? | **No.** The unchecked state matches absence, and no completion chronology exists. |

### 4.3 — Distribution/pinned-Git decision

| Field | Evidence |
|---|---|
| Original acceptance intent | “Decide distribution after core/profile stability; evaluate pinned-Git consumption without commitment.” — `tasks.md` Phase 4, item 4.3. |
| Classification | `superseded-decision` |
| Provenance | Current governance was introduced by later Public Preview work: `AGENTS.md`, `openspec/config.yaml`, and `docs/product-charter.md` select module identity, Apache-2.0, and planned source/dev preview identity `v0.1.0-preview.1`, while explicitly stating that it is **not** distribution and that broad distribution, packaging, signing/provenance, and release automation are deferred. `docs/product-charter.md` retains Milestone 5: “Evaluate pinned Git consumption and choose packaging, versioning, and release automation only after the core and first profile are stable.” Archived Public Preview artifacts, especially `openspec/changes/archive/2026-09-17-public-preview-readiness/{proposal.md,design.md,publication-audit.md}`, explicitly distinguish source-build preparation from publication/distribution. |
| Commands and actual results | Current governance and later archived artifacts were read; no distribution command, tag, release, package, upload, or external action was run. |
| Semantic gaps | The preview identity is a current source/dev decision, not evidence of a distribution decision. The first CNSIC profile is absent, so the originally stated “after core/profile stability” condition has not become a completed distribution evaluation. No pinned-Git evaluation record was found. |
| Checkbox reconciliation supportable? | **No.** The later decision supersedes portions of bootstrap-era identity deferral but expressly leaves the requested distribution/pinned-Git decision deferred; changing this unchecked item would conflate preview identity with distribution and fabricate chronology. |

## Re-run command results

The following commands were run in the foreground against the current worktree. Durations are the Go tool’s reported package elapsed times where it printed them; the check-only formatter produced no output.

| Exact command | Actual result |
|---|---|
| `go test -count=1 ./internal/git ./internal/inventory` | PASS — `internal/git` 0.160s; `internal/inventory` 0.019s. |
| `go test -count=1 . -run 'Test(Account|Accounting|ResultAccounting)'` | PASS — root package 0.001s. |
| `go test -count=1 . -run 'Test(PrepareCarveout|MatchPreparedCarveout)' && go test -count=1 ./cmd/git-change-evidence -run 'Test(RunUsesCommittedBlobsAndFormatsExactly|MatchedSharesAbsoluteDeadline)'` | PASS — root 0.001s; CLI 0.018s. |
| `go test -count=1 . -run 'Test(ProjectEvidence|ResultV1ReportBinding|ValidateReportV1ResultProvenance)'` | PASS — root package 0.002s. |
| `go test -count=1 ./internal/publication ./cmd/git-change-evidence -run 'Test(Publish|RunPublish)'` | PASS — publication 0.001s; CLI 0.004s. |
| `go test ./...` | PASS — root 0.020s; census 0.001s; CLI 0.027s; internal/git 0.154s; internal/inventory 0.017s; internal/publication 0.001s. |
| `go test -race ./...` | PASS — root 1.204s; census 1.013s; CLI 1.180s; internal/git 1.218s; internal/inventory 1.027s; internal/publication 1.007s. |
| `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"` | PASS — exit 0; no output; check-only. |

## Chronology caveat

`git log --all -- openspec/changes/bootstrap-neutral-core-extraction/{tasks.md,apply-progress.md}` shows the bootstrap task record was last changed by the verified bootstrap commits `73e2e6f` and `9e8a4e7`. Later task-surface implementation commits exist, including snapshot (`8c77cd4`), inventory (`e383744`), accounting (`2363ecd`), carveout (`5f80917`, `3104231`), report/projection (`3ff69ff`, `4b8f53f`), publication (`db8c546`, `c3a2796`), and result composition (`ede5d73`), but they did not update this bootstrap ledger. Their current source/tests are useful behavioral evidence; they are not evidence that the original bootstrap task sequence was performed or recorded in the original order.

Accordingly, this matrix preserves the canonical **3/12** ledger. It records technically present later behavior without backfilling task checkboxes, historical RED/GREEN evidence, apply-progress entries, or lifecycle conclusions.
