# Neutral core — report privacy 3.2

## Authority
Approved bootstrap-neutral-core-extraction task 3.2; overseer 01a11fe8; base 8c14cc288ff031c4840b374acbc23290ae6623ad.
Jose's explicit choice relayed by overseer (ce6ea4d8-520a-4ad7-a20d-2ad6de999c42, seq13): preserve_canonical_subject_fidelity.
Preserve canonical JSON Subject and bytes/digest/schema. Canonical evidence may contain sensitive caller content; emitter must avoid secrets and handle evidence as sensitive.
Human projection and diagnostics remain privacy-preserving and evidence-only. No silent Subject sanitization/rejection, export/transport or authority expansion.
RDD off; local verified commits allowed; no push/PR/merge, destructive changes, Pi/Herdr or CNSIC pilot. Preserve foreign roots.

## Tasks
- [x] U32-1 — Document the approved canonical/human/diagnostic privacy boundary in authoritative specs. **done**. Verified local commit f09adab8db88bcab5f4aeedfa6bcc756a71f2b01; two specs +4/-2, parent tracking included; readback/diff PASS; no meaningful behavioral RED.
- [x] U32-2 — Verify existing report/diagnostic behavior and identify only meaningful missing assertions. **done**. Commit ca2db5ed29095c69e684d5ecb0144ce4e993268d. Go1.25.10 root selected 13 top/40 subtest PASS, CLI 4/4 PASS, zero FAIL/SKIP; six-package suite uncached, format/diff PASS. Root selection includes two incidental carveout tests, not remapped.
- [x] U32-3 — Add justified sensitive Subject contract coverage; verify bounded candidate and close canonical 3.2 locally. **done**. Evidence commit 55b528e03c8f90237cd5e8a51a11ede894220278; five final checks PASS; only 3.2 checked, 7/12.

## Allowed surfaces and verification
Initial writer: only openspec/changes/bootstrap-neutral-core-extraction/specs/neutral-evidence-contracts/spec.md and specs/evidence-cli-publication/spec.md (same change root).
No production/test changes initially. Passive clarification, no meaningful behavioral RED; structural readback/diff checks.
Verifier source edits: none; per-run runtime writes only <private-evidence> and children, 0700, pinned Go1.25.10 with owned HOME/cache/temp.
Authorized next test-only writer: report_test.go and cmd/git-change-evidence/main_test.go, two boundary cases for sensitive-looking Subject fidelity/human omission; target 70–95 additions, hard stop before exceeding 120 additions. No production edits. Private writer verification <private-evidence> (0700), owned caches/temp.
No helpers/frameworks/dependencies/new schemas or parsers; no actual secrets/environment/production data. Synthetic sentinels only if required.
Closure may update canonical tasks.md/apply-progress.md and parent tracking after proof; only 3.2 expected 7/12, future tasks untouched.

## Acceptance and proof limits
Canonical report bytes and digest preserve caller Subject exactly, including synthetic sensitive-looking text; neither canonical output nor arbitrary Subject claimed sanitized.
Human projection omits sensitive Subject and raw names/paths/environment/command output; fixed diagnostics do not echo error/input details or imply approval/gating/delivery.
Report/projection revisions/provenance remain deterministic and immutable; distinguish caller-injected data from automatic host-data collection (none demonstrated).
Observe actual executed test assertions, counts/failures/skips and commands; identify existing tests versus missing significant negatives before writes.
Coverage-only tests must not fabricate semantic RED or justify production changes absent a proven defect. Historical RED not rederived.
Candidate unit boundary is base 8c14cc2, never accumulated feature branch; respect bounded review workload.

## State
Current canonical progress 7/12, only 3.2 newly checked; later tasks unchanged. Existing human/diagnostic guards and canonical fidelity retained; two synthetic coverage cases added without production changes.
Subject boundary resolved by explicit human decision and committed narrow SoT clarification f09adab. Existing checks all PASS under that contract; logs <private-evidence>. No implementation defect; production validation/hash/schema unchanged.
Significant missing domain: generic fixtures cannot exclude selective sensitive-Subject sanitization/rejection or sensitive-only human echo. Add only synthetic path/credential/environment/process-text Subject cases at core projection/decode and CLI canonical/human boundaries; assert decoded Subject equality, unchanged canonical bytes/digest, human omission/exact expected projection and empty stderr.
Coverage-only exception: tests characterize correct existing behavior; no meaningful semantic RED or historical RED claim, no production change justified. Do not mutate production to fabricate RED.
Writer completed exactly TestProjectEvidenceSyntheticSubjectBoundary (+46) and TestRunProjectSyntheticSubjectBoundary (+44): 90 additions, no deletions or production changes. Core exact decoded Subject/canonical bytes/digest/provenance and human fixed-label omission; CLI exact canonical/core-human outputs, all five synthetic sentinels absent in human, success/empty stderr.
Writer Go1.25.10 focused checks PASS (two root tests, one CLI test/two subtests), zero failures/skips; allowed-file gofmt/diff PASS. Logs <private-evidence>. Parent root-test spot check and path/numstat confirm scope.
ASSESS medium/runtime-large/RDD-off: writer self-verification true, independent verifier false; no consumed authority. Final functional command runner PASS: root two top-level tests, CLI one top/two subtests, focused FAIL/SKIP 0; six-package suite all uncached; exact format/diff checks PASS. Logs <private-evidence> (0700).
Final initial/final HEAD/branch/status and all three candidate files/diff byte-identical; metadata preserved; no index/ref byte-identity claim. Application stderr assertions empty; standard Go module-download stderr not application output. No race or historical/semantic RED/architecture-review claim.
Closure: docs tasks +1/-1, progress +8/-0; readback/count/diff and parent spot check PASS. Evidence commit 55b528e includes two tests plus tracking/task/progress (105 additions/3 deletions), postcommit clean; final ASSESS medium/runtime-large/RDD-off, self-verification/no separate verifier/no consumed authority.
Next: read-only mapping of approved 3.3 CLI/publication semantics and race proofs, no delivery action or publication.
