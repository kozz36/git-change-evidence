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

Per-file crossings: unavailable until PR 1 selects layout/thresholds; later slices record post-PR counts and new crossings.

| PR | Test; harness | Rollback | P/T/M/A |
|---|---|---|---|
| 1 Decision | pending bootstrap decision; pending bootstrap decision | decision files | 0/0/30/20 |
| 2 W1/W2 | pending bootstrap decision; pending bootstrap decision | contracts/vectors | 190/190/10/10 |
| 3 Git/inventory | pending bootstrap decision; real temporary Git | shell | 200/210/0/0 |
| 4 Accounting | pending bootstrap decision; policy fixtures | accounting core | 150/160/0/0 |
| 5 W3 | pending bootstrap decision; blob/dirty-tree | core/edge | 180/180/0/0 |
| 6 W4 assembly | pending bootstrap decision; repeat | report core | 150/140/0/0 |
| 7 W4 shell | pending bootstrap decision; conflict/race | CLI/shell | 180/190/0/0 |
| 8 W5 | pending bootstrap decision; supplied/missing | comparison | 120/120/0/0 |
| 9 CNSIC profile | pending bootstrap decision; parity corpus | standalone profile | 140/140/0/20 |
| 10 Handoff/distribution | N/A—decision artifacts; N/A | decision records | 0/0/50/90 |

## Phase 1: Bootstrap Decision (planning; PR 1)

- [ ] 1.1 Decide minimal `AGENTS` boundary, source/test layout, runtime, reproducible test provider/command, and due configuration representation; defer packaging, CI, release, license, distribution.
- [ ] 1.2 Refresh `openspec/config.yaml` testing rubric; record layout thresholds and crossing baseline.

## Phase 2: Neutral W1/W2 (implementation; PRs 2–4)

- [ ] 2.1 RED then GREEN closed/versioned contracts, canonical bytes/digest/provenance, authority rejection, and frozen vectors/manifest with origin identities.
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
