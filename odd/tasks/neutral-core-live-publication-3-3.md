# Neutral core — live publication 3.3

## Authorization
Jose chose authorize_gce_livepublication_api through overseer 01a11fe8; relay 4a40d4f1-dcf8-420a-b0f0-c74410f7112b seq16.
Starting base 7b91a5afccb93f59b84f38b7f9a9cab1dedf86f8; initial canonical 7/12 with 3.3 open. Product evidence only, no delivery/gating authority.
Authorized opt-in Go shell API, not a CLI/parser/assembler: Publish(ctx, repo, baseRef, headRef, destination, build).
Builder captures policy/profile/Subject and receives immutable IDs; caller builds complete report using existing neutral contracts.
Acquire → build/validate report provenance for same pair → reacquire/revalidate → existing frozen publisher. Changed pair discards and rebuilds all evidence.
Three complete attempts; exhausted movement returns existing typed SnapshotRacing, with no publisher invocation or artifact.
Snapshot semantics: last observed matching pair before Publish; no atomic read of both refs, no Git/filesystem transaction, no post-check ref-stability guarantee.
Existing frozen publication and CLI remain independent of Git. Canonical Subject/bytes/hash/schema/validation and human/diagnostic privacy remain unchanged.
No publication/push/PR/merge, destructive operations, Pi/Herdr, CNSIC pilot/DB/seed/rebuild/login. Coordinated CNSIC data/runtime hold remains until explicit release.

## Tasks
- [x] U33-1 — Implement authorized shell API and narrow contract clarification using observed behavioral RED → GREEN. **done**; local commit 7ef8055e9060d6de4dcfc4b2ffe51291c6d7b66f.
- [x] U33-2 — Run focused/full/race checks and follow actual candidate assessment; verify complete bounded 3.3 matrix. **done**; verification work unit 468d5647a999283602ed1edfee215726db77216e, five-group code PASS and bounded doc PASS.
- [x] U33-3 — Reconcile only verified canonical 3.3/progress and commit locally with identity evidence. **done**; closure commit 3a2507e1b0fab50098d930714967297bc8cc59f5; canonical 8/12, next 3.4 W5 supplied forecast.

## Allowed edit surfaces and budget
shell/livepublication/publish.go
shell/livepublication/publish_test.go
openspec/changes/bootstrap-neutral-core-extraction/design.md
openspec/changes/bootstrap-neutral-core-extraction/specs/evidence-cli-publication/spec.md
This parent tracking file; closure tasks.md/apply-progress.md only after verified acceptance.
No other production/test files authorized; stop with minimal proposal if exports/helpers outside this scope are required.
Estimate 195–305 implementation/test/doc lines; stop before 400 diff lines including tracking/doc changes, split or ask overseer before expansion.
Writer verification writes only <private-evidence> and children (0700); pinned Go1.25.10, owned HOME/cache/module/temp and synthetic real-Git/filesystem fixtures.
Next verifier source edits none; private <private-evidence> (0700), no shared resource permissions.

## Acceptance and checks
Reachable public shell API with public caller-nameable builder/input/result/error surface; Git effects stay outside neutral root and existing publisher unchanged.
Real-Git stabilization forces full rebuild and publishes only evidence for final observed pair, canonical bytes/digest/provenance bound to it.
Persistent movement exhausts three attempts with typed technical failure and no artifact; builder/provenance mismatch rejects before publication.
Post-check ref movement documents observed-snapshot limit, not impossible atomicity; no raw names/secrets/environment/process output inserted into human diagnostics.
Behavioral RED must execute compiled tests against a minimal non-publishing implementation scaffold; missing function/package compilation failure alone is wiring, not semantic RED.
GREEN and focused checks, full suite and applicable race/format/diff verification must be observed. No historical RED or exhaustive/cross-platform/power-loss claims.
Existing exclusive publication/content identity/technical exit evidence must be mapped distinctly (real filesystem vs fake/seam/inspection), with any necessary additional scope proposed rather than silently expanded.
Canonical progress is 8/12 after complete 3.3 verification and local closure; no tests/runtime operations against CNSIC or production data.

## Observed evidence
Worker mv0vdeeh-w-nxd5: compiled non-publishing scaffold RED seven behavioral leaves (exit 1), GREEN seven leaves (exit 0); historical scaffold not independently replayed.
Independent verifier mv0vmq0q-x-o0nb: five check groups exit 0; live 3 top/6 subtests (7 leaves), publication 5 top, CLI 34 top/74 subtests, zero FAIL/SKIP; all seven packages and race three packages PASS; format/diff PASS.
Native ASSESS was unassessable (untracked declaration), RDD off/runtime-large: conservative independent-verifier plan followed, no native review authority.
Doc correction mv0vtn5j-y-ckrv and independent revalidation mv0vvr8i-z-akw0 PASS: caller root unchanged, relative identity derived; required 186 >160 boundary record in apply-progress; exact two-Go before/after/current hashes equal, no global metadata/exclusivity claim.
Limits: external caller compile and new nil/cancel guards inspected only; Git acquisition uses background context; fresh NewReportV1 fixtures do not prove caller inventory/accounting/FromResult recomputation; post-check real-ref/real-publisher seam demonstrates observation limit, not atomicity/power-loss safety.
Private evidence: <private-evidence>,final,docfix,docverify}-20261009; no runtime artifacts committed. Implementation commit 7ef8055e9060d6de4dcfc4b2ffe51291c6d7b66f; verification commit 468d5647a999283602ed1edfee215726db77216e; canonical closure 3a2507e1b0fab50098d930714967297bc8cc59f5.
Bounded 3.3 matrix accepted: new public real-Git/real-filesystem stable/rebuild, real-ref movement with zero-call spies for rejection, post-check real-ref/real-publisher observation seam; existing frozen real-filesystem exclusive/content/preservation/cancellation and syscall staging/rollback seams, fake-CLI technical exits/root forwarding, canonical/human privacy. No required code/test gap remains; executed/inspected limits above retained.
