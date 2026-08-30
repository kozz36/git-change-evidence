# Tasks: Bootstrap Neutral Core Extraction

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | production 1,310; tests 1,330; mechanical 90; artifacts 140 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes — autonomous stacked slices to `main` |
| Delivery / chain | `auto-chain` / `stacked-to-main` |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

Layout thresholds and the initial file-crossing baseline are authoritative in `openspec/config.yaml`: PR 1 begins at zero Go files/tests and zero Go-file change crossings. Later slices must record first crossings in `apply-progress.md` before exceeding the selected repository-development thresholds. Those thresholds are not generic-core evidence or governance policy.

| PR | Test; harness | Rollback | P/T/M/A |
|---|---|---|---|
| 1 Decision | check-only Go format gate; N/A — planning-only with no module/runtime boundary | `AGENTS.md`, OpenSpec/configuration, and decision docs | 0/0/197/202 |
| 2 W1/W2 | `go test ./...`; compatibility-vector harness | contracts/vectors | 190/190/10/10 |
| 3 Git/inventory | `go test ./...`; real temporary Git | shell | 200/210/0/0 |
| 4 Accounting | `go test ./...`; policy fixtures | accounting core | 150/160/0/0 |
| 5 W3 | `go test ./...`; blob/dirty-tree | core/edge | 180/180/0/0 |
| 6 W4 assembly | `go test ./...`; repeat | report core | 150/140/0/0 |
| 7 W4 shell | `go test ./...`; conflict/race | CLI/shell | 180/190/0/0 |
| 8 W5 | `go test ./...`; supplied/missing | comparison | 120/120/0/0 |
| 9 CNSIC profile | `go test ./...`; parity corpus | standalone profile | 140/140/0/20 |
| 10 Handoff/distribution | N/A — decision artifacts; N/A | decision records | 0/0/50/90 |

## Phase 1: Bootstrap Decision (planning; PR 1)

- [x] 1.1 Record the minimal root `AGENTS.md`, Go 1.25.10 target, public root API/`cmd/git-change-evidence/`/`internal/` layout, co-located `*_test.go` tests, primary `go test ./...` gate, concurrency `go test -race ./...` gate, check-only formatting gate, and immutable typed Go configuration boundary. Defer module/publication identity, packaging, CI, release, license, distribution, and concrete configuration-file syntax.
- [x] 1.2 Refresh `openspec/config.yaml` as the authoritative Go strict-TDD/testing rubric; record exact commands, layout thresholds, and the zero-file initial crossing baseline without fabricating Go source or tests.

## Phase 2: Neutral W1/W2 (implementation; PRs 2–4)

- [x] 2.1 RED then GREEN closed/versioned contracts, canonical bytes/digest/provenance, authority rejection, and frozen vectors/manifest with origin identities.
- [ ] 2.2 RED then GREEN immutable Git/inventory: SHA-1/256, hostile paths, rename/binary/symlink/gitlink, missing/racing objects, dirty/index isolation, argv/closed-env/replacement isolation, and no-follow races in real Git.
- [ ] 2.3 RED then GREEN injected policy: overlap/invalid/default classification, non-countable data, text totals, non-gating thresholds.

## Phase 3: Staged Capability (implementation; PRs 5–8)

- [ ] 3.1 RED then GREEN W3 bounded immutable carveout: invalid bound, unavailable, timeout, deterministic match, dirty-tree isolation; positional legacy edge only.
- [ ] 3.2 RED then GREEN W4 report/projection and diagnostics redaction: names, secrets, environment, paths, command output never leak or imply authority.
- [ ] 3.3 RED then GREEN W4 CLI/publication: technical outcomes; destination/interruption/content identity/moving-ref race proofs; atomic exclusive writes.
- [ ] 3.4 RED then GREEN W5 supplied forecast: valid, missing, divergent, unavailable, threshold; never author, mutate, infer, or gate.

## Phase 4: Profile and Handoff (planning + implementation; PRs 9–10)

- [ ] 4.1 RED then GREEN standalone CNSIC profile: one-way vocabulary/authority rejection, distinct-profile neutrality, frozen-vector and real-Git parity.
- [ ] 4.2 Create separate CNSIC SDD handoff: dependency integration, reversible cutover, owner observation period, local-code deletion; no CNSIC edit here.
- [ ] 4.3 Decide distribution after core/profile stability; evaluate pinned-Git consumption without commitment.
