# Result V1 Archive-Time Specification Synchronization Report

## Authorization and input evidence

This synchronization uses the explicitly authorized archive-time fallback because the native status supplied for `result-v1` has no synchronization field. It is a canonical-specification merge only; it does not publish, deliver, commit, push, create a PR, or alter implementation files.

- Host archive date: `2026-09-07`.
- Branch and code baseline: `test/result-v1-unit9-mutation-census` at `17854e7539ca571107a01b35907695a979f1295c`.
- Accepted verification input: `verify-report.md`, SHA-256 `b2d41ca63bef0d6fa9e378a7d64d12ebe021fa708b895b35a0c0df1067f0d2ae`.
- Accepted verification verdict: **PASS WITH WARNINGS**; 13/13 requirements, 38/38 scenarios, 12/12 tasks, 60 focused tests, 124 full-suite top-level tests, 928 subtest PASS events, and 2 CLI tests passed. Historical FAIL findings F1–F3 are closed and retained verbatim in the report appendix.
- Required archive validation policy file `.claude/rules/sdd-validation.md` was requested before archive work but is absent from this worktree (the `.claude` directory and tracked path are absent). The injected archive instructions and this report therefore record the required preflight, tracked-file, scenario-census, and validation-stamp evidence directly.

## Preflight and canonical operations

The complete change-local delta set was preflighted before writing. Every delta contains only `## ADDED Requirements`; there are no `MODIFIED`, `REMOVED`, or `RENAMED` sections. The four canonical destinations were absent, so the operation creates their initial canonical specifications. No destructive synchronization approval was needed or consumed.

| Capability | Operation | Requirements | Scenarios | Canonical SHA-256 |
|---|---|---:|---:|---|
| `evidence-v1-provenance` | ADD | 1 | 4 | `9d3e413d1ce67bd4bd6382fa51a61bbe57ad7586c10973dd588530ab967aed5d` |
| `result-v1-accounting-evidence` | ADD | 5 | 14 | `b05fa9fb96ff10f23fb77071a0516885ac43bf605add1f1c4aa0f1b2a70a6c5c` |
| `result-v1-contract` | ADD | 4 | 12 | `a8eba847394d7323937431e80000e0698e01a43ee8399816cbe590af49f1e658` |
| `result-v1-provenance` | ADD | 3 | 8 | `b0913ccec71378cccc8d004b61dc0c615fd76561370a105850203504fc5f76b3` |
| **Total** | **ADD** | **13** | **38** | — |

The canonical copies mechanically retain every requirement and scenario body. Only delta-to-canonical framing changed (`Delta Specification`/`ADDED Requirements` to `Specification`/`Requirements` where applicable), followed by a `Validated against code` stamp that points to the archived formal verification report, exact branch/HEAD, and passing census. There are **0 modified requirements**, **0 removed requirements**, and **0 replaced requirements**.

## Validation and integrity

- The formal report was read before sync and is the accepted current verification evidence. This docs-only archive phase does not rerun tests.
- The authoritative `tasks.md` was reread immediately before this synchronization and has all 12/12 implementation checkboxes complete; SHA-256 `475b01e9d8523ed21ced86680be28de55a4ee7200c26c543fc852c8306e9778a`.
- `apply-progress.md` SHA-256 is `6d438276a6c309ae61e62ddfcc9df2d3d2dc99cae309384c2d65612c7e85146d`.
- The protected original tracked-Go manifest is 82 files with SHA-256 `9e95f953e6ebe929d35ab2ff1fca0f91f8af794d6d937eacce8bb3916cdf7794`; the current 83-file Go manifest, including the untracked supplemental test, is SHA-256 `f6d6ccb3ac59c548d95f2e665a6e7303f504a341c12c97545afb5c961964e0e8`. The supplemental test SHA-256 is `41d5e0e5374cb18b20ae4e94dc2eb7244dc3c3bd55dde9d25bfde833ff4544cc`.
- The roadmap sweep found only `bootstrap-neutral-core-extraction` and this change as active changes (beside the archive container). No tracked active-path references to `result-v1` exist outside this change, and no preexisting canonical reference existed. No dangling active-path reference requires an out-of-scope edit.

## Retained limitations and delivery boundary

The canonical contract is synchronized despite the formal report's nonblocking warnings: W1 absent optional strict-TDD verifier support; W2 native-SHA-256/hostile-path coverage is composed rather than one public exact-fixture golden; and W3 retains optional invalid-binding/overflow vectors. Historical TDD recovery depends on a local session artifact and is not repository-portable. These are disclosed limitations, not a delivery claim.

Units 0–8 are recorded as shipped baseline history, while Unit 9 evidence, the supplemental test, the current tasks/progress changes, and the formal report are local and uncommitted. Lifecycle archive closes the specification change only; it does not deliver or adopt the implementation, publish a contract, or confer consumer authority.

## Next archive action

Synchronization completed without destructive deltas. Preserve all active artifacts byte-for-byte and move the complete `result-v1` change directory, including this report and the archive report, to `openspec/changes/archive/2026-09-07-result-v1` only after final destination and artifact-integrity assertions pass.
