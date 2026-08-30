# Design: Bootstrap Neutral Core Extraction

## Technical Approach

Use **Branch by Abstraction**: immutable CNSIC W1/W2 at merge `a0fc7b26ff8a0e0a61baa586b32c46841611c806` remains the predecessor oracle while `git-change-evidence` builds the neutral successor. No mechanical copy or simultaneous consumer rewrite.

## Architecture and Ownership

| Logical boundary | Owns | Excludes |
|---|---|---|
| Functional Core | Neutral contracts/canonical bytes, normalization, exclusive injected-policy accounting, bounded carveout, reports/projections, forecast comparison | Effects, CNSIC vocabulary, delivery authority |
| Imperative Shell | Immutable Git acquisition, separate inventory, CLI adaptation, bounded retries, diagnostics, exclusive atomic publication | Classification or governance decisions |
| Compatibility edge | Frozen vectors and legacy-carveout translation to immutable inputs | A second matcher/accounting implementation |
| CNSIC profile | CNSIC policy, vocabulary, compatibility, and workflow integration | Neutral contract or Git semantics |

Dependencies flow shell/profile/edge → neutral interfaces → core; concrete layout remains undecided.

## Interfaces / Contracts

| Contract | Boundary |
|---|---|
| Bounded carveout | Invalid bounds return typed invalid-bound evidence and no analysis result. Unavailable content and timeout remain distinct typed outcomes, never inferred matches, mismatches, or measurements. |
| Python library | A neutral, versioned API exposes validation, evidence construction, and report use without selecting packaging, concrete layout, or runtime policy. |
| CLI outcomes | Map success, invalid input, unavailable evidence, and publication conflict distinctly. Threshold attention remains evidence-only within successful technical completion. |
| Forecast comparison | Validate supplied forecasts but never author, infer, complete, mutate, or overwrite them, including when evidence is unavailable. |

## Transition Plan and Gates

1. **Successor appears:** implementation first appears in standalone `git-change-evidence` during neutral-core extraction; CNSIC W1/W2 remains the oracle.
2. **Generic source of truth changes:** only after frozen vectors and real-Git cases pass independently and byte-for-byte where contractual. Failure preserves predecessor authority.
3. **Consumer cutover:** CNSIC consumption and repo-local deletion occur only in a later CNSIC SDD. `/arranque tools` remains its orchestration/handoff track; implementation runs from `git-change-evidence`.

Serial order: W1/W2 freeze → **W3** bounded core + shell blob acquisition + legacy edge → **W4** report/projection core + CLI/publication shell → **W5** neutral forecast comparison + CNSIC profile → later cutover.

Import both frozen vectors and manifest into a versioned corpus recording origin, merge SHA, source/blob identities, and digest. Preserve bytes; they pin observed behavior, not generic names, count, or schema. Mirrors retain CNSIC archive provenance.

No-dual-evolution gate: before transition, successor changes require parity; afterward, generic semantics evolve only standalone. Rollback selects the predecessor; later cutover retains a reversible adapter until owner-approved deprecation. Delete CNSIC-local generic code only after parity, observation, and rollback-window closure.

## Verification Strategy

Use pure contract/accounting tests and real temporary Git repositories for both object formats, hostile byte paths, kinds/modes, missing objects, mutable-state isolation, and races. Prove publication identity/exclusivity/cleanup and profile one-way neutrality.

## Threat Matrix

| Boundary | Status and response | Planned RED proof |
|---|---|---|
| Git repository/commit state | Applicable: explicit root/native objects; ignore staged, `commit -a`, empty index | Relative/absolute/foreign roots; clean/dirty/staged/empty-index parity |
| Shell/process | Applicable: argv-only, closed environment, no shell | Metacharacter input and hostile Git config cannot execute or alter output |
| Filesystem | Applicable: bounded descriptor-based reads and atomic exclusive writes | Missing/replaced/special files fail typed with no partial artifact |
| Publication | Applicable: canonical-content identity, bounded retry, no overwrite | Existing destination and interrupted write preserve prior bytes |
| Paths/document-like names | Applicable: raw bytes are data; names never imply execution | `requirements.txt`, `CMakeLists.txt`, executable MD/MDX, `README.sh`, non-UTF-8 paths |
| Refs | Applicable: resolve before evidence and revalidate moving symbolic refs | Hostile, missing, foreign, replaced, and persistently moving refs fail safely |
| Symlink races | Applicable: no-follow every component; descriptor binding | Intermediate/final pre/post-open replacement cannot escape root |
| Diagnostics | Applicable: typed, sanitized, evidence-only | Secrets, environment, output, and paths are absent/redacted |
| Consumer boundary | Applicable: one-way profiles, no authority injection | Different profiles preserve neutral meaning; gate/verdict fields reject |
| Push state / PR commands | N/A: no push, refspec destination, `gh`, or PR automation is in product scope | None |

## Repository Readiness

| State | Items |
|---|---|
| Present | `.git`, `.gitignore`, `.codegraph`, `.atl/skill-registry.md`, `openspec/config.yaml` |
| Deferred to owning decisions | AGENTS/project conventions; packaging/project marker; runtime pin; formatter, linter, typechecker, and test configuration; CI; license; release automation; configuration syntax |

Before the first implementation slice, owning decisions must establish public module boundaries, source/test roots, a supported runtime, and a reproducible test command/provider. Their filenames, tools, packaging, and syntax are not selected here.

## File Changes

This phase creates only this `design.md`; implementation paths remain unselected.

## Open Questions

- Which owning decisions will define source layout, runtime support, test provider, and profile-configuration syntax?
- What owner-approved observation period closes the later CNSIC rollback/deprecation window?
