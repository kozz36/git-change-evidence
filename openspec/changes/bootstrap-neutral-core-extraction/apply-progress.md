# Apply Progress: Bootstrap Neutral Core Extraction

## PR 1 Historical Slice
**Completed:** PR 1 / Phase 1 — Bootstrap Decision (`auto-chain`, `stacked-to-main`).

This planning/configuration slice establishes Go-first boundaries only. It adds no Go source, tests, `go.mod`, compatibility vectors, CNSIC changes, commits, pushes, or pull requests.
## Completed Tasks

- [x] **1.1** Recorded the single root `AGENTS.md`, Go 1.25.10 target, root public API/`cmd/git-change-evidence/`/`internal/` layout, co-located tests, Go gates, and immutable typed configuration boundary.
- [x] **1.2** Replaced the pre-bootstrap skip rubric with the authoritative Go strict-TDD rubric, exact commands, layout thresholds, and a zero-file initial crossing baseline.

## Files Changed

- `AGENTS.md`
- `README.md`
- `docs/product-charter.md`
- `openspec/config.yaml`
- `openspec/changes/bootstrap-neutral-core-extraction/exploration.md`
- `openspec/changes/bootstrap-neutral-core-extraction/proposal.md`
- `openspec/changes/bootstrap-neutral-core-extraction/design.md`
- `openspec/changes/bootstrap-neutral-core-extraction/specs/neutral-evidence-contracts/spec.md`
- `openspec/changes/bootstrap-neutral-core-extraction/tasks.md`
- `openspec/changes/bootstrap-neutral-core-extraction/apply-progress.md`

## Verification Evidence

| Check | Result |
|---|---|
| `git diff --check` | PASS |
| `python3` with PyYAML: parse `openspec/config.yaml` | PASS |
| `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 \| xargs -0 -r gofmt -l)"` | PASS; no Go files exist and the check performed no writes |
| Search for active legacy-language library requirement in `README.md`, `docs/`, and `openspec/` | PASS; none remains |
| Root `AGENTS.md` inventory | PASS; exactly one file at `./AGENTS.md` |
| `go test ./...` | N/A for PR 1; the approved no-module boundary means no `go.mod` or Go package exists yet |
| `go test -race ./...` | N/A for PR 1; no concurrency-bearing Go surface exists |
| Runtime harness | N/A; this slice has no runtime boundary |

The executor host reports `go version go1.27.0-X:nodwarf5 linux/amd64`, not the selected Go 1.25.10 baseline. The host toolchain was not used as evidence for the target baseline.

## TDD Cycle Evidence

| Task | RED | GREEN | TRIANGULATE | REFACTOR | Result |
|---|---|---|---|---|---|
| 1.1 | N/A — planning documentation only; no production behavior | N/A | N/A | N/A | Complete |
| 1.2 | N/A — test-policy configuration only; no production behavior | N/A | N/A | N/A | Complete |

Strict TDD is active for later Go production work. The required support file `.pi/gentle-ai/support/strict-tdd.md` is absent, so later implementation must follow the RED → GREEN → TRIANGULATE → REFACTOR contract directly and retain this risk until the support file is supplied.

## Design Conformance and Deviations

- The runtime is Go-first; all active Python-library requirements were replaced with a neutral versioned Go package API and canonical-document contracts.
- The Functional Core accepts immutable typed Go structs/interfaces. Configuration parsing remains adapter-only; TOML, YAML, JSON, and other syntax choices remain deferred.
- The public/internal layout and thresholds are development guidance only, never core evidence or governance policy.
- **Intentional deviation from future implementation:** no `go.mod` or Go files were added. Module/publication identity is still unresolved and outside PR 1.

## Workload and PR Boundary

```text
main
└── PR 1 — Bootstrap Decision 📍
    └── PR 2 — Neutral W1/W2 implementation
```

- **Strategy:** `auto-chain` with `stacked-to-main`.
- **Current boundary and review budget:** only tasks 1.1 and 1.2; 285 additions + 114 deletions = 399 changed lines.
- **Dependency:** immutable CNSIC W1/W2 predecessor oracle at merge `a0fc7b26ff8a0e0a61baa586b32c46841611c806`.
- **Rollback boundary:** remove the listed PR 1 decision/configuration artifacts or restore their pre-PR-1 wording; no production or consumer behavior is affected.
- **Follow-up:** PR 2 may begin task 2.1 only after its owner-authorized module identity is available. Do not begin any later task in this PR.

## Remaining Tasks

- [ ] 2.1 through 2.3 — Neutral W1/W2 implementation.
- [ ] 3.1 through 3.4 — Staged carveout, report, CLI/publication, and forecast capabilities.
- [ ] 4.1 through 4.3 — CNSIC profile, separate cutover handoff, and distribution decision.

## Unresolved Decisions and Risks

- Module/publication identity, packaging, CI, license, release, distribution, and concrete configuration-file syntax remain unresolved by design.
- The executor's installed Go toolchain is not Go 1.25.10.
- `.pi/gentle-ai/support/strict-tdd.md` is absent. This did not affect planning-only work, but it is a strict-TDD risk for the first Go implementation slice.

## PR 2 / Task 2.1 — Neutral W1/W2 contracts
- [x] **2.1** Closed/versioned contracts, canonical bytes/SHA-256/provenance, authority rejection, and the corrected predecessor corpus pass all acceptance criteria.
- **Initial record (corrected):** the prior “frozen W1/W2 vectors” claim was overbroad: both vectors were successor-authored, not predecessor bytes.
- **Independent rejection / correction:** replaced them with exact predecessor `contracts-v1.json`, `accounting-v1.json`, and `manifest.sha256`; deleted both misleading vectors.
- **Files:** `go.mod`, `contract.go`, `decode.go`, `contract_test.go`, `decode_test.go`, `compatibility_test.go`, the four corpus/manifest files, `openspec/config.yaml`, `tasks.md`, and this record.
- **Verification:** focused corpus RED/GREEN and decoder characterization, `go test .`, `go test ./...`, `go vet ./...`, and `git diff --check origin/main` → PASS.
- **Format:** check-only `gofmt -l` → PASS; **Race:** N/A (no concurrency).
### TDD Cycle Evidence
| Task | RED | GREEN | TRIANGULATE | REFACTOR |
|---|---|---|---|---|
| 2.1 initial | Public-symbol compile failure | Constructor/canonical pass | Earlier rejection proof, now narrowed | Decoder simplification pass |
| 2.1 correction | Corpus test failed on false manifest | Exact predecessor imports pass | Negative decoder cases already passed (characterization) | `slices.Contains` pass |
### Design Conformance and Deviations
- Functional Core only; the standalone manifest is successor-authored metadata binding predecessor identities, not a predecessor artifact.
### Workload / PR Boundary
- **Budget:** 393 additions + 4 deletions = 397 changed lines relative to `origin/main` (hard cap: 400).
- **Boundary:** stacked-to-main PR 2 task 2.1 only; rollback removes PR 2 contracts, tests, and corpus without affecting PR 1 planning.
### Remaining Tasks / Risks
- [ ] 2.2 onward remain outside this slice; strict-TDD support was read from the predecessor checkout because it is absent here.

## Public Preview Task 7 — Bounded Reconciliation Decision
- User-resolved decision: **Preservar y separar**.
- Preserve active `bootstrap-neutral-core-extraction` at exactly 3/12; no reconciliation, archive, or sync is performed.
- The final matrix census is 4 proven-present / 4 partial / 3 absent / 1 superseded-decision; this evidence census does not alter task chronology or checkbox status. See [reconciliation-matrix.md](reconciliation-matrix.md).
- The matrix-identified residual capabilities are separated into active proposal-level change `post-preview-neutral-capabilities`; see [explore.md](../post-preview-neutral-capabilities/explore.md) and [proposal.md](../post-preview-neutral-capabilities/proposal.md).
- That residual change contains exactly matrix residuals 2.2, 2.3, 3.2, 3.3, 3.4, 4.1, 4.2, 4.3 at proposal level only; no specs, design, tasks, implementation, or verification exist yet.
- Every future phase is post-preview and nonblocking for planned source/dev `v0.1.0-preview.1`.
- This is evidence-only and conveys no authority.
- Planning-only strict TDD and test execution are N/A. No production, tests, or external actions occurred.
- Task 7 creates no authority for Task 8 or later publication work.

## Task 2.2 — Bounded Git/inventory reconciliation
- Local commits: `2bc0fb587e476b1d378937822581e843c5f77e5a`, `4362b2aa456311f2118ed9a866b90e119a48f0d0`, `8bc302b8c00f23764326b7e53fb1fcec19809bba`, `a83346e3b6825ea9fc120b71c43628fff3635d95`.
- Immutable Git/inventory coverage proof, not a behavioral-defect fix: SHA-1/256, hostile paths, committed kinds, missing/racing objects, dirty/index and argv/closed-environment/replacement isolation.
- Real-Git inventory evidence: exact separate record and unchanged snapshot; static intermediate/final symlink escapes, intermediate/final pre-open replacement, and final post-open replacement.
- Internal hook is invocation-local; normal acquisition uses nil callbacks. The literal intermediate/final pre/post-open requirement is preserved: intermediate post-open is justified by shared-guard inspection, NOT an executed fixture; final post-open, pre-open replacements, and static escapes were observed. No exhaustive component/timing/format or cross-platform claim.
- Forwarded independent final review found no severe issue or missing required behavior in its bounded Linux/amd64 interpretation; overseer accepted task 2.2 closure on that bounded proof, overall 4/12, with 2.3 and later unchanged.
- Forwarded Go 1.25.10/Linux amd64 checks: `go test -count=1 -v ./internal/git ./internal/inventory` — 21 top-level PASS, 30 subtest PASS, 0 skips, 0 failures; ordinary helper return is not independent coverage.
- `go test ./...` — six packages PASS; `go test -race ./internal/git ./internal/inventory` — both PASS; no whole-suite race or full/race individual skip-count claim.
- Exact check-only format gate PASS: `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"`; `git diff --check` PASS in forwarded verification.
- Coverage-only TDD exception: initial missing internal signature compile failure was wiring RED only; semantic/historical behavioral RED was not rederived.
- Optional configured external-diff/textconv fixtures remain unverified. Private evidence: `<private-evidence>/gce-tools-2-2-final-20261009`; no logs copied or tests rerun here.
- Clone RDD off; technical evidence only, no native delivery authority or publication. No task 2.3 implementation or completion.

## Task 2.3 — Existing neutral injected-policy reconciliation
- Reconciled already-implemented neutral policy; no new implementation, defect fix, or semantic change. Only 2.3 checked: local 5/12; later tasks unchanged.
- Neutral canonical construction: exclusive first-match overlap (raw 4/source 0), injected default `other` (7/1); 24 invalid-policy contract cases, typed `invalid_document` antecedents/no partial validation, constructor zero/fatal result.
- Non-countable addition 99 ignored in text totals; checked text totals and unavailable/overflow cases verified. Threshold 4 > 3 and ratio 3/4 > 0.5 produce observations only; equality does not exceed.
- Neutral/compatibility shared vectors are supplementary, not a replacement for neutral construction evidence.
- Forwarded independent Go 1.25.10/Linux amd64 verification: root 93 top-level + 909 subtests PASS, 0 failures/0 skips; primary `go test ./...` six packages PASS, not cached.
- Forwarded exact check-only format gate and `git diff --check` PASS; commands executed by verifier, not rerun here. Private evidence: `<private-evidence>`.
- Historical RED not rederived; documentation-only reconciliation, current semantic change 0; no race run for pure policy. Technical evidence only, no delivery authority or publication.

## Task 3.1 — Existing bounded immutable carveout reconciliation
- Independent final PASS: existing bounded W3 core, immutable blob acquisition and legacy three-positional CLI verified; no required gap, new implementation or semantic change. Only 3.1 checked: 6/12; later tasks unchanged; this latest record supersedes historical remaining-task counts.
- Nine invalid typed bounds/zero preparation and side-specific unavailable; seven matcher deadline fixtures yield nil result/fallback; deterministic four ordered matches and immutable result API verified.
- Committed blob identity/content isolated from dirty/index/divergent worktrees; exact positional output moved=1, new=1, additions=2 and failure exits 2/3/5 verified.
- CLI deadline nil result/exit 6 executed; stdout suppression guard INSPECTED, not a full deadline-output fixture; no acquisition during deletion-only interval.
- Forwarded Go 1.25.10 Linux/amd64: root 93 top/909 subtests, CLI 33/72, Git 14/23 PASS; all FAIL/SKIP 0; primary six-package suite PASS, not cached; exact check-only format and diff checks PASS.
- Verifier executed commands, not rerun here; private evidence: `<private-evidence>`. Historical RED not rederived; malformed bound text outside int-typed API; optional exact unavailable stderr not a required gap.
- Proof limits: no race (unchanged concurrency), cross-platform or exhaustive claim; RDD off, technical evidence only, no delivery authority/publication; no new-behavior CHANGELOG/release entry.

## Task 3.2 — Canonical Subject fidelity and private human projection
- Only 3.2 checked: 7/12; 3.3 onward unchanged; this latest record supersedes historical remaining-task counts.
- Jose's explicit `preserve_canonical_subject_fidelity` choice, two-spec clarification `f09adab` and baseline `ca2db5e` govern: canonical Subject bytes/hash/schema/validation unchanged; JSON may contain sensitive caller content, so emitters must avoid secrets and handle evidence as sensitive. Human/diagnostic privacy remains required; no canonical sanitization/privacy guarantee.
- Two new synthetic boundary tests: 90 additions (46 root, 44 CLI), production 0; direct decoded Subject/canonical bytes/digest/provenance fidelity, exact human omission of all five sentinels, core/CLI dispatch, success and empty application stderr.
- Forwarded writer GREEN and final functional Go 1.25.10 Linux/amd64: root two top-level PASS (new case + fixed oracle), CLI one top/two subtests PASS, focused FAIL/SKIP 0; six packages uncached PASS; exact format/diff PASS; all five final checks PASS, exit 0.
- Commands forwarded, not rerun here; private evidence `<private-evidence>`. Coverage-only/passive documentation TDD exception; no semantic/historical RED, no race (no concurrency change), no production-source change.
- ASSESS medium/runtime Large; RDD off, writer self-verification, no separate architectural review or consumed authority. Technical evidence only, no approval/gating or publication; no new-behavior CHANGELOG/release entry.
