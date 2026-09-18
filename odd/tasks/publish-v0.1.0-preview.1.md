# Publish v0.1.0-preview.1

## Outcome

Publish the first truthful Public Preview only after lifecycle evidence, owner-gate
audits, final human confirmation, and a final-commit rebuild all agree. This plan
coordinates future work; it grants no approval, release, visibility, delivery,
or destructive authority by itself.

## Authorized program decisions

| Topic | Decision |
| --- | --- |
| Audit first | Run an automated legal, confidentiality, dependency, secret, Git-history, unpublished-ref, and hosting-controls audit before any lifecycle or publication action. Automated results are evidence, not legal clearance. |
| Visibility stop | Keep the repository private until a human reviews the final audit findings and explicitly confirms the exact visibility/publication target. |
| Lifecycle debt | Before release, canonically verify/archive `go-ast-census` and reconcile `bootstrap-neutral-core-extraction`. Map implemented behavior instead of duplicating it; separate genuinely pending W5, profile, handoff, and later distribution work. |
| Preview assets | Upload exactly one deterministic source archive, Linux `amd64` and `arm64` archives containing both executable names, one CycloneDX SBOM per Linux target, and `SHA256SUMS`. Assets are unsigned. |
| Formal agent skill | Defer skill design and implementation until after the preview is publicly verified. |
| Publication authority | Issue, PR, merge, visibility, tag, release, upload, and any remediation or history action require the authority named by their task. A passing check never supplies that authority. |

## Current starting evidence

- Release-preparation branch: `docs/prepare-v0.1.0-preview.1`.
- Preparation work unit: `79de3adfcfd240925537835453b1f6c51d3b3916`.
- Evidence ledger: `ba59a284ebdfa20bb100962aa0ffd8ca45dceac0`.
- `origin/main` preparation predecessor:
  `16aaa4498f91c6d1ff91b330e1df6a47db26c922`.
- Existing preparation plan remains authoritative historical evidence at
  `odd/tasks/prepare-v0.1.0-preview.1.md`; its Task 2 blocker must not be
  erased or represented as complete until canonical archive evidence exists.
- No tag, release, uploaded asset, public-visibility change, package
  publication, or formal agent skill exists under this program.

## Execution rules

1. Use a dedicated branch and isolated Git worktree for each lifecycle change
   and for release publication preparation. Never reuse a dirty worktree.
2. Keep exactly one program task in progress. Record exact commits, trees,
   commands, outputs, PRs, and external URLs only after observing them.
3. Canonical SDD tooling owns verify, sync, archive, and lifecycle mutation.
   Never manually move an OpenSpec change or rewrite historical reports.
4. Every repository change is a reviewable conventional work-unit commit.
   Open issue-linked PRs, require repository CI and review, and merge only
   under separate delivery authority.
5. Never print secret values. Audit reports record locations, classifications,
   hashes where safe, coverage, and owner dispositions—not credentials.
6. A failed or incomplete check stops its dependent tasks. Recovery uses a new
   scoped corrective work unit; it does not silently weaken acceptance.
7. The mandatory human-confirmation task is a hard stop. No visibility, tag,
   release, or upload command may run before its explicit recorded answer.
8. Publication artifacts are always rebuilt from the final tagged commit.
   Rehearsal outputs and hashes are documentary only and must never be reused.

## Dependency path

```text
owner-gate audit and disposition
  -> tracking issue
  -> go-ast verify/archive PR
  -> bootstrap reconciliation/archive PR
  -> release-preparation PR and CI merge
  -> exact final commit + refreshed audit
  -> HUMAN CONFIRMATION STOP
  -> public visibility
  -> exact tag
  -> fresh unsigned assets
  -> GitHub Release upload
  -> public verification
  -> optional post-preview skill follow-up
```

## Phase A — Audit before repository publication work

- [x] 1. **Run the automated owner-gate audit first.** In a fresh read-only worktree at the current candidate commit, inventory the current tree, full reachable Git history, local and remote refs visible to the authorized account, large/binary files, generated metadata, local paths, dependency/license/NOTICE inputs, likely secrets, confidential identifiers, personal data indicators, and repository-host settings. Inspect visibility, branch protection, Actions permissions/history, security features, access/deploy-key metadata, existing tags/releases/packages, and unpublished refs without reading or reproducing secret values. Record tools, versions, exact scope, exclusions, findings, false-positive handling, hashes, and unavailable checks in a dedicated audit artifact. **Evidence:** audit report `docs/publication-audit-v0.1.0-preview.1.md` SHA-256 `3604c20cfdc124c444ec1b965036eea9e1b6d80b80f8c979cc768664941944e9`; branch `audit/v0.1.0-preview.1`, commit `2536397213efd973b1dd83f2c96d683a9b58a9f8`, tree `2cba92cb6c1f2b1bc6be411f3c55a6e22e3d0f58`; clean pre-audit index/worktree; all local refs covered (83 commits, 601 objects, 265 blobs); read-only remote/GitHub inventory recorded; fallback secret scan found no credential/private-key pattern but has bounded limitations; independent verification PASS confirmed the corrected 24-worktree and 9-unavailable-endpoint census, exact scanner recipe, report hash, safe output, and two-file scope. Owner disposition remains required for F-01 through F-06 and unavailable checks; Task 2 disposition may start, but downstream lifecycle/publication remains stopped until publication-relevant findings are resolved or accepted. **Recovery:** discard incomplete reports and rerun from a clean worktree; any credential, legal, confidentiality, or history-remediation action requires separate explicit authority.

- [ ] 2. **Obtain owner/legal disposition for every audit finding.** Classify each finding as cleared, accepted limitation, remediation required, or unresolved. Attach evidence without exposing sensitive content. Create separately scoped remediation plans for required fixes; a history rewrite, credential action, deletion, legal determination, or confidentiality response is never implied by this plan. Do not proceed while a publication-relevant finding is unresolved. This is finding disposition, not the final visibility confirmation. **Depends on:** 1. **Evidence:** pending signed/recorded disposition matrix and links to any authorized remediation. **Recovery:** return to Task 1 after remediation and issue a new audit revision; never edit away the earlier report.

- [ ] 3. **Create and approve the publication tracking issue.** After audit disposition, use the repository issue-first workflow to create one issue for the Public Preview publication program. Link this plan, the preparation plan, audit revision, both lifecycle debts, required PRs, the human stop, exact six-asset release set, and post-publication checks. The issue must state that technical evidence does not authorize publication. **Depends on:** 2. **Evidence:** pending issue number, URL, approval/triage result, and immutable audit references. **Recovery:** correct the issue through ordinary comments/edits; do not open implementation PRs against an unapproved or materially incomplete issue.

## Phase B — Canonical lifecycle reconciliation

- [ ] 4. **Freshly verify `go-ast-census` in an isolated branch/worktree.** Start from the then-current target branch, confirm 9/9 task state, preserve both historical FAIL reports, and run canonical verification against the current residual-C3 correction. Require a new report that distinguishes current behavior from historical failures and validates all 10 requirements and 30 scenarios. Do not manually edit a verdict or reuse the stale canonical FAIL as a pass. **Depends on:** 3. **Evidence:** pending branch/worktree, native lineage/evidence revision, report hash/verdict, focused/full/format results, and unchanged historical-report hashes. **Recovery:** retain the failed report, correct only through a separately authorized bounded work unit, and re-verify; no archive until this program's canonical-verify requirement is met.

- [ ] 5. **Canonically archive `go-ast-census` and merge its lifecycle PR.** Use canonical archive tooling only, with fresh collision/delta validation and complete historical preservation. Confirm the active source moved to the dated archive, canonical specs were updated as reported, and an archive report records the current verification plus historical failures. Commit the archive as one conventional lifecycle unit; open the issue-linked PR, require CI/review, and merge under separate authority. **Depends on:** 4. **Evidence:** pending archive path/report/hash, commit, PR, CI runs, review outcome, merge commit, and proof the active source no longer remains. **Recovery:** if canonical archive does not execute or validate, leave the active change untouched and stop; never perform a manual move.

- [ ] 6. **Map bootstrap tasks to current implementation without changing code.** In a separate isolated bootstrap-reconciliation worktree, build a requirement/task-to-evidence matrix for all 12 tasks. Re-run acceptance evidence for materially present 2.2, 2.3, 3.1, 3.2, and 3.3 behavior; identify any semantic gaps rather than assuming file presence proves completion. Confirm that 3.4 supplied forecast, 4.1 CNSIC profile, and 4.2 handoff are genuinely absent, and classify 4.3's preview distribution decisions versus its post-profile/pinned-Git decision. Preserve the original 3/12 ledger until this matrix is independently reviewed. **Depends on:** 5. **Evidence:** pending matrix path/hash, commands, exact source/spec links, independent review, and explicit present/partial/absent results. **Recovery:** correct the matrix, not production code; any discovered product gap becomes a separate bounded implementation task.

- [ ] 7. **Canonically reconcile bootstrap and separate post-preview work.** Through canonical SDD apply authority, update bootstrap task/apply evidence only for requirements independently proven by Task 6. Preserve provenance from later archived changes and do not manufacture historical TDD chronology. Create a separate active post-preview SDD change for genuinely pending supplied forecast, CNSIC profile, reversible handoff/observation/deletion, and the later distribution/pinned-Git decision. That residual change must remain explicitly nonblocking for this CLI-focused preview and evidence-only. **Depends on:** 6. **Evidence:** pending reconciled task counts, provenance map, residual change name/artifacts, collision/relationship results, and independent review. **Recovery:** revert only the reconciliation work unit if authority or mapping is invalid; never duplicate current production merely to satisfy an old checkbox.

- [ ] 8. **Verify/archive the reconciled bootstrap and merge its lifecycle PR.** Run canonical verification against the reconciled scope, preserving incomplete/post-preview separation as explicit evidence. Canonically archive the reconciled bootstrap only after validation; do not archive or complete the new residual change. Commit one lifecycle work unit, open the issue-linked PR, require CI/review, and merge under separate authority. **Depends on:** 7. **Evidence:** pending verify/archive reports and hashes, archived task state, residual active status, commit, PR, CI, review, and merge commit. **Recovery:** leave bootstrap active if canonical verification/archive fails; fix through a new scoped work unit, never a manual move or task-state rewrite.

## Phase C — Release-preparation integration

- [ ] 9. **Refresh the release-preparation branch onto the reconciled base.** From the updated default branch after both lifecycle merges, create a fresh isolated release worktree. Reapply or merge the reviewed preparation work without rewriting published history, resolve only observed conflicts, and confirm the changelog, runbook, ODD evidence, Apache-2.0 license, canonical/legacy CLI behavior, active residual SDD change, and evidence-only language remain truthful. **Depends on:** 5 and 8. **Evidence:** pending branch base/head/tree, retained work-unit identities or replacement mapping, changed paths/LOC, and clean status. **Recovery:** discard the release worktree and repeat from the merged base; never force-update a shared ref without explicit authority.

- [ ] 10. **Run final release-preparation verification.** Execute `go test ./...`, `go test -race ./...`, the exact configured check-only formatting command, `git diff --check`, documentation links, license/version checks, root-Go inventory, secret/portable-path scan, candidate-residue scan, lifecycle-state readback, and realistic negative-coverage review. Confirm no release asset is present and no publication claim is premature. **Depends on:** 9. **Evidence:** pending exact commands, outputs/hashes, task/path census, and independent read-only PASS. **Recovery:** correct the candidate in a new work-unit commit and repeat the complete verification.

- [ ] 11. **Open the issue-linked release-preparation PR.** Use the repository PR workflow, include the exact audit revision and both lifecycle merge commits, explain the six-asset unsigned release contract, list deferred post-preview work, and keep visibility/tag/release/upload unchecked. Protect reviewer workload with coherent commits or a separately approved size exception. **Depends on:** 10. **Evidence:** pending PR URL/number, base/head commits, changed paths/LOC, review instructions, and issue linkage. **Recovery:** update the branch with ordinary corrective commits and rerun affected checks; do not hide review history.

- [ ] 12. **Require CI, review, and merge for release preparation.** Require all configured PR checks and policy checks, independent review, resolved findings, and exact merge-commit capture. Merge only under separate repository authority. After merge, verify the default branch contains both lifecycle archives and release-preparation evidence and has no unexpected files. **Depends on:** 11. **Evidence:** pending CI run URLs/conclusions, review disposition, merge method, merge commit/tree, and clean default-branch readback. **Recovery:** stop on failed or missing checks; correct through a new PR/work-unit commit rather than bypassing policy.

## Phase D — Final owner stop and publication

- [ ] 13. **Freeze the exact final release commit and refresh automated audits.** Select one immutable default-branch commit as the only candidate for tag `v0.1.0-preview.1`. Re-run Task 1's automated audit and Task 10's verification against that exact commit, including hosting settings and tag/release/package absence. Produce a final publication manifest binding commit, tree, audit revision, CI, lifecycle archives, asset contract, limitations, and owner decisions. Keep the repository private. **Depends on:** 12. **Evidence:** pending final commit/tree, audit/report hashes, CI, manifest, and clean reproducibility readback. **Recovery:** any source or setting change invalidates the freeze and returns to the applicable earlier task.

- [ ] 14. **STOP for final human confirmation.** Present the final audit findings/dispositions, unresolved limitations, exact repository, current visibility, final commit/tree, proposed tag, exact six upload filenames, unsigned status, release notes, and rollback/incident options. Ask the human to explicitly confirm or decline public visibility and the subsequent tag/release/upload sequence. No automation, prior plan, issue approval, CI pass, or technical report may substitute for this response. **Depends on:** 13. **Evidence:** pending exact confirmation prompt, human answer, timestamp, actor, target binding, and any conditions. **Recovery:** a decline or ambiguous answer stops the program without visibility/tag/release/upload; changed inputs require a new Task 13 audit and confirmation.

- [ ] 15. **Change repository visibility to public and verify access.** Only after Task 14's explicit confirmation, perform the single authorized visibility change against the exact repository. Verify unauthenticated metadata/read access, default branch, license, README, security-sensitive settings, and absence of unintended releases/packages/assets. Record the before/after setting and command/API response without exposing credentials. **Depends on:** 14. **Evidence:** pending authorization binding, visibility operation result, public URL, unauthenticated checks, and settings readback. **Recovery:** stop on uncertainty or mismatch; any visibility reversal or settings remediation requires fresh human authority.

- [ ] 16. **Create the exact immutable preview tag.** Confirm the frozen commit/tree still match Task 13, the tag does not exist locally or remotely, and all prior tasks remain valid. Create `v0.1.0-preview.1` pointing exactly to the frozen commit using the separately authorized tag method; never move or reuse the tag. Push only that tag under explicit publication authority and verify its remote target. **Depends on:** 15. **Evidence:** pending tag object/type, target commit/tree, local/remote identities, command result, and public readback. **Recovery:** stop on collision or target mismatch; never force-update a public tag. Escalate any correction as a new owner decision.

- [ ] 17. **Build fresh unsigned assets from the tag.** In a new owner-only temporary worktree/directory checked out at the exact tag, build deterministic source, Linux `amd64`, and Linux `arm64` archives. Each Linux archive contains exactly `gce` and `git-change-evidence`. Generate unmodified per-target CycloneDX JSON SBOMs from each exact `gce` binary with pinned `cyclonedx-gomod@v1.10.0` `bin -json -noserial -notimestamp -version v0.1.0-preview.1`; do not request unavailable main-module license download and do not post-process JSON. Generate `SHA256SUMS` over exactly the other five upload assets. Upload set: source archive, two binary archives, two target SBOMs, and `SHA256SUMS`—six unsigned files, no rehearsal artifact or signature/provenance attachment. Validate members, architectures, aliases, SBOM linkage, version/help/census parity, checksums, local-path absence, and cleanup boundaries before upload. **Depends on:** 16. **Evidence:** pending temp path, tool versions, filenames, sizes, hashes, runtime/architecture results, SBOM summaries, checksum verification, and pre-upload manifest retained as evidence but not uploaded. **Recovery:** delete the complete temporary candidate set and rebuild from the tag; never repair or reuse a partial build.

- [ ] 18. **Create the GitHub Release and upload the exact assets.** Create release `v0.1.0-preview.1` from the exact public tag with notes derived from the reviewed changelog and explicit preview limitations. Upload only the six Task 17 files, verify each returned asset name/size/hash against local evidence, and record the public release URL. Do not publish a Go proxy/package, container, installer, signature, provenance bundle, or additional platform asset. **Depends on:** 17. **Evidence:** pending release ID/URL, tag target, notes hash, six asset API records, local/remote size/hash comparison, and command results. **Recovery:** stop on partial/ambiguous upload. Any deletion, replacement, re-upload, or release-state change requires explicit human authorization and a documented incident/correction plan.

- [ ] 19. **Verify the public preview from an unauthenticated environment.** From a fresh external directory/session, resolve the public repository and tag, download all six assets, verify `SHA256SUMS`, inspect archive members/architectures and CycloneDX structure, run supported `amd64` version/help/census checks, verify documented links/license, and confirm no unintended package/release/asset exists. Record arm64 runtime as observed or explicitly unavailable. Compare public commit/tree/tag/release identities to the frozen manifest. **Depends on:** 18. **Evidence:** pending unauthenticated commands, URLs, downloaded hashes, runtime output, identity matrix, and final PASS/FAIL. **Recovery:** report discrepancies immediately; do not delete, replace, hide, or republish without new human authority.

## Phase E — Explicit post-preview deferral

- [ ] 20. **Open a separate post-preview follow-up for the formal agent skill.** Only after Task 19 records the public preview outcome, evaluate whether a dedicated agent skill is warranted. Use a separate issue/SDD scope for skill triggers, authority boundaries, workflow, tests, documentation, registry update, and distribution. Do not add a skill to this publication program or make it a retroactive preview prerequisite. **Depends on:** 19. **Evidence:** pending follow-up decision/issue or explicit no-action record. **Recovery:** keep the skill deferred; no preview artifact or release metadata changes are needed.

## Completion statement

This program is complete only when Tasks 1–19 have observed evidence, Task 2's
findings are dispositioned, Task 14 has an explicit human confirmation, the
public tag/release/assets match the frozen final commit, and Task 19 verifies
them externally. Task 20 remains a post-preview follow-up and does not change
the release outcome. Completion records facts; it does not retroactively grant
publication authority to any agent or tool.
