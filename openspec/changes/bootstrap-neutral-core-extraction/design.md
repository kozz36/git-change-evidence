# Design: Bootstrap Neutral Core Extraction

## Technical Approach

Use **Branch by Abstraction**: immutable CNSIC W1/W2 at merge `a0fc7b26ff8a0e0a61baa586b32c46841611c806` remains the predecessor oracle while `git-change-evidence` builds the neutral Go successor. No mechanical copy or simultaneous consumer rewrite.

Go 1.25.10 is the target baseline. This PR records the runtime, API/layout, typed configuration, and testing decisions without adding a Go module identity, Go source, vectors, or production implementation.

## Architecture and Ownership

| Logical boundary | Owns | Excludes |
|---|---|---|
| Functional Core | Neutral contracts/canonical bytes, normalization, exclusive injected-policy accounting, bounded carveout, reports/projections, forecast comparison over immutable typed Go structs/interfaces | Effects, configuration parsing, CNSIC vocabulary, delivery authority |
| Imperative Shell | Immutable Git acquisition, separate inventory, CLI adaptation, bounded retries, diagnostics, exclusive atomic publication | Classification or governance decisions |
| Compatibility edge | Frozen vectors and legacy-carveout translation to immutable inputs | A second matcher/accounting implementation |
| CNSIC profile | CNSIC policy, vocabulary, compatibility, and workflow integration | Neutral contract or Git semantics |

Dependencies flow shell/profile/edge → neutral interfaces → core. The public neutral Go package API is at the repository/module root; the CLI is at `cmd/git-change-evidence/`; non-public shell/adapters are under `internal/`. Tests are co-located as `*_test.go`. `pkg/`, a Python-style `src/` plus `tests/` split, and package-level `AGENTS.md` files are prohibited.

Core configuration inputs are immutable typed Go structs/interfaces. TOML, YAML, JSON, or other configuration decoding belongs only to adapters. A concrete configuration-file syntax remains deferred.

## Interfaces / Contracts

| Contract | Boundary |
|---|---|
| Bounded carveout | Invalid bounds return typed invalid-bound evidence and no analysis result. Unavailable content and timeout remain distinct typed outcomes, never inferred matches, mismatches, or measurements. |
| Versioned Go package API | A neutral, versioned public Go API exposes validation, evidence construction, and report use alongside canonical document contracts. It does not select module/publication identity or configuration-file syntax. |
| CLI outcomes | Map success, invalid input, unavailable evidence, and publication conflict distinctly. Threshold attention remains evidence-only within successful technical completion. |
| Forecast comparison | Validate supplied forecasts but never author, infer, complete, mutate, or overwrite them, including when evidence is unavailable. |
| Configuration boundary | Adapters decode external configuration and construct immutable typed core inputs. The core neither parses external formats nor chooses a format. |

## Transition Plan and Gates

1. **Successor appears:** implementation first appears in standalone `git-change-evidence` during neutral-core extraction; CNSIC W1/W2 remains the oracle.
2. **Generic source of truth changes:** only after frozen vectors and real-Git cases pass independently and byte-for-byte where contractual. Failure preserves predecessor authority.
3. **Consumer cutover:** CNSIC consumption and repo-local deletion occur only in a later CNSIC SDD. `/arranque tools` remains its orchestration/handoff track; implementation runs from `git-change-evidence`.

Serial order: PR 1 bootstrap decision → **W1/W2** contracts, immutable Git/inventory, and accounting → **W3** bounded core + shell blob acquisition + legacy edge → **W4** report/projection core + CLI/publication shell → **W5** neutral forecast comparison + CNSIC profile → later cutover.

Import both frozen vectors and manifest into a versioned corpus recording origin, merge SHA, source/blob identities, and digest. Preserve bytes; they pin observed behavior, not generic names, count, or schema. Mirrors retain CNSIC archive provenance.

No-dual-evolution gate: before transition, successor changes require parity; afterward, generic semantics evolve only standalone. Rollback selects the predecessor; later cutover retains a reversible adapter until owner-approved deprecation. Delete CNSIC-local generic code only after parity, observation, and rollback-window closure.

## Verification Strategy

Future Go production changes follow RED → GREEN → TRIANGULATE → REFACTOR. The primary reproducible gate is `go test ./...`; concurrency-bearing surfaces additionally run `go test -race ./...`; the Go formatting gate is check-only. These gates become executable once an authorized later work unit establishes the deferred module identity. This PR contains no Go source or module and therefore has no RED/GREEN production-test cycle.

Use pure contract/accounting tests and real temporary Git repositories for both object formats, hostile byte paths, kinds/modes, missing objects, mutable-state isolation, and races. Prove publication identity/exclusivity/cleanup and profile one-way neutrality.

## Layout Thresholds and Baseline

The repository-development thresholds below limit package growth; they are not generic-core accounting, governance, or evidence policy.

| Surface | Threshold | PR 1 baseline |
|---|---:|---:|
| Public root Go files | 12 | 0 |
| CLI Go files under `cmd/git-change-evidence/` | 6 | 0 |
| Go files in an internal leaf package | 15 | 0 |
| Changed lines in one Go file | 160 | 0 |
| Changed Go files in one package | 8 | 0 |

The first crossing of any threshold must be recorded in `apply-progress.md` and trigger an explicit package-boundary review before further files are added beyond that threshold.

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
| Present after PR 1 | `.git`, `.gitignore`, `.codegraph`, `.atl/skill-registry.md`, root `AGENTS.md`, authoritative `openspec/config.yaml`, and planning artifacts |
| Selected by PR 1 | Go 1.25.10 target, root public API, `cmd/git-change-evidence/`, `internal/`, co-located tests, strict TDD test commands, typed configuration boundary, and thresholds/baseline |
| Intentionally absent | Go source files, tests, `go.mod`, vectors, generated artifacts, CI, and release automation |
| Deferred to owning decisions | Module/publication identity, packaging, CI, license, release, distribution, concrete configuration-file syntax, and CNSIC rollback-window observation period |

No source or tests are fabricated by this planning slice. A later authorized implementation slice must establish the module identity before it can execute the selected Go test commands.

## File Changes

PR 1 creates `AGENTS.md` and `apply-progress.md`, and updates this design, proposal, neutral-contract spec, tasks, OpenSpec configuration, README, charter, and exploration status. It creates no production implementation path.

## Open Questions

- Which owner-approved module/publication identity will allow a later implementation slice to establish `go.mod`?
- What concrete configuration-file syntax should adapters support after a consumer need exists?
- What owner-approved observation period closes the later CNSIC rollback/deprecation window?
