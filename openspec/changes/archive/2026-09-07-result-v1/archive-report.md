# Result V1 Archive Report

## Archive record

- **Change:** `result-v1`
- **Archive date:** `2026-09-07` (host date)
- **Branch / HEAD:** `test/result-v1-unit9-mutation-census` / `17854e7539ca571107a01b35907695a979f1295c`
- **Destination:** `openspec/changes/archive/2026-09-07-result-v1`
- **Destination precondition:** absent before this operation; this report is written only after that refusal check and successful canonical sync.
- **Lifecycle meaning:** archived specification-change audit trail only. This is not delivery, adoption, publication, commit, push, PR, merge, deployment, release, or consumer authority.

## Accepted closure evidence

The current formal verification report, SHA-256 `b2d41ca63bef0d6fa9e378a7d64d12ebe021fa708b895b35a0c0df1067f0d2ae`, is **PASS WITH WARNINGS**: 13/13 requirements, 38/38 scenarios, 12/12 tasks, no current blockers, and closed historical F1–F3 findings retained verbatim. It records 60 focused tests, 124 full-suite top-level PASS events, 928 subtest PASS events, and 2 CLI tests passing. The report was read before archive; no source or test command is rerun in this docs-only phase.

`tasks.md` was reread immediately before reporting/archive and remains 12/12 complete at SHA-256 `475b01e9d8523ed21ced86680be28de55a4ee7200c26c543fc852c8306e9778a`. `apply-progress.md` remains SHA-256 `6d438276a6c309ae61e62ddfcc9df2d3d2dc99cae309384c2d65612c7e85146d`.

## Archive-time canonical synchronization

Archive-time fallback synchronization was explicitly authorized because native status supplied no sync field. `sync-report.md` SHA-256 `3d283458555fc3923ee2a7729396f8b9fde9e4fd65f7c3c5bc86389dda564b54` records the non-destructive preflight and canonical merge:

- 4 new canonical capability specs;
- 13 added requirements and 38 preserved scenarios;
- 0 modified, removed, or replaced requirements;
- no `MODIFIED`, `REMOVED`, or `RENAMED` delta blocks; and
- one `Validated against code` stamp in each canonical spec, bound to the formal verification report and exact code baseline.

No preexisting canonical Result V1 spec was overwritten. No active-path reference to this change was found outside its own pre-archive directory during the roadmap sweep, so no dangling reference requires an out-of-scope repair.

## Audit preservation and workload

The archive relocation moves 12 documentation artifacts as one directory: the ten preexisting change artifacts (exploration, proposal, design, tasks, apply progress, verification report, and four delta specs) plus `sync-report.md` and this report. The ten preexisting artifacts were hash-verified immediately before the move and are not rewritten. In particular, the verification report's historical FAIL appendix remains byte-for-byte intact.

Workload is separated from delivery: canonical synchronization added four canonical documents and stamps; archive administration added two audit reports and moved 12 documentation artifacts; no Go source/test/config file was written. The protected tracked-Go manifest is 82 files / SHA-256 `9e95f953e6ebe929d35ab2ff1fca0f91f8af794d6d937eacce8bb3916cdf7794`; the current 83-file Go manifest including the supplemental test is SHA-256 `f6d6ccb3ac59c548d95f2e665a6e7303f504a341c12c97545afb5c961964e0e8`. The supplemental test remains SHA-256 `41d5e0e5374cb18b20ae4e94dc2eb7244dc3c3bd55dde9d25bfde833ff4544cc`.

## Retained risks and handoff

- Formal warnings W1–W3 remain: absent optional strict-TDD verifier support, composed rather than one public native-SHA-256/hostile-path golden, and optional invalid-binding/overflow vectors.
- Historical TDD recovery is local-session evidence and not repository-portable.
- Units 0–8 are historical shipped baseline; Unit 9 evidence/test/report and current documentation changes remain local and uncommitted. Archive does not settle those delivery concerns.
- The required `.claude/rules/sdd-validation.md` file is absent in this worktree, so this archive record preserves the injected-rule preflight and validation evidence explicitly.

After the immediately following directory rename, parent readback should verify the archived path, canonical specifications, artifact hashes, and repository status. No additional lifecycle or delivery action is recommended by this executor.
