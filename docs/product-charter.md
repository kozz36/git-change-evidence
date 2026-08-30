# Product Charter: git-change-evidence

`git-change-evidence` is a generic multi-project Git evidence CLI and Go package. It creates reproducible, machine-readable evidence about a change without deciding whether that change may be delivered.

**Status:** the Bootstrap Decision is complete. This charter is the authoritative product boundary until superseded by a versioned governance decision.

## Product intent

Manual reconstruction of changed-file inventories, line accounting, thresholds, forecasts, and review evidence is slow and can produce inconsistent or stale records. `git-change-evidence` will derive those observations from immutable Git snapshots, explicit policy, and versioned contracts so that authors, reviewers, and automation can inspect the same evidence.

The product is **evidence-only**. It MUST NOT approve, reject, gate, block, merge, deploy, release, assign delivery authority, or substitute for a project's review process. A threshold crossing is an observation or warning, never a product verdict.

## Users and use cases

| User | Use case | Expected outcome |
|---|---|---|
| Change author | Generate a local evidence packet for a base/head revision pair. | Deterministic inventory, accounting, thresholds, and projections. |
| Reviewer | Independently regenerate or inspect evidence cited by a change. | Comparable canonical data without reconstructing Git behavior manually. |
| Project maintainer | Define a project profile with categories, thresholds, and vocabulary. | Project policy remains explicit, versioned, and outside the generic core. |
| Automation integrator | Consume stable machine-readable contracts from a CLI or versioned Go package API. | Neutral, versioned data that can be adapted to a project workflow. |

## Value proposition

- **Reproducibility:** fixed immutable revisions and the same profile yield the same canonical evidence.
- **Traceability:** reports retain revision identity, input-policy identity, and inventory/accounting provenance.
- **Portability:** one core supports multiple projects without embedding any one project's governance rules.
- **Review clarity:** summaries and projections expose evidence while preserving human or project-owned decision authority.

## Product boundaries

### In scope

- Immutable Git snapshot acquisition, including commit identity, change entries, file kinds, and byte-preserving paths.
- Separate untracked-file inventory where a profile requests it; index and working-tree state do not alter committed-change accounting.
- Policy-driven, exclusive change accounting.
- Neutral, versioned canonical JSON documents and a versioned Go package API.
- Deterministic report assembly, projections, local publication, and CLI behavior.
- Project/profile adapters that translate project configuration and compatibility requirements into immutable typed core inputs.

### Non-goals

- Approval, rejection, gating, blocking, merge, deployment, release, or delivery authority.
- A universal changed-lines limit, review rubric, forecast vocabulary, or governance model.
- Remote delivery orchestration, pull-request management, CI administration, UI, database, or RBAC features.
- Replacing project review protocols, policy owners, or release processes.
- Selecting or implementing a module/publication identity, packaging, release automation, licenses, CI, distribution, or a concrete configuration-file syntax during this bootstrap stage.

## Authority model

The CLI and Go package may emit factual statuses such as an unavailable measurement, an incomplete input, or a policy attention condition. Their meaning is limited to evidence production.

Every project retains responsibility for interpreting that evidence. Consumer integrations MUST preserve this separation: they may display or transport an observation, but they must not represent a core output as approval, denial, or a delivery gate.

## Architecture

The product uses **Functional Core / Imperative Shell** with **Hexagonal Architecture** project-profile adapters.

| Boundary | Responsibility | Must not own |
|---|---|---|
| Functional core | Normalization, validation, exclusive policy accounting, contract construction, canonical serialization, report assembly, projections, and comparisons over immutable typed Go inputs. | Process execution, filesystem effects, configuration parsing, project-specific policy, or delivery authority. |
| Imperative shell | Git invocation, filesystem access, input/output, retry orchestration, local publication, and CLI adaptation. | Policy decisions or non-deterministic business rules. |
| Project-profile adapter | Project globs, category names, threshold values, configuration decoding, framework integration, compatibility schemas, and review integration. | Generic Git semantics, core contract semantics, or delivery authority. |

The public neutral Go API belongs at the repository/module root; the CLI belongs at `cmd/git-change-evidence/`; non-public shell and adapters belong under `internal/`. Tests are co-located as `*_test.go`. The repository deliberately has no `go.mod` in this planning slice because module/publication identity remains deferred.

The core accepts immutable typed Go structs and interfaces. TOML, YAML, JSON, and any other external configuration decoding belong only in adapters. No concrete configuration-file syntax is selected.

The generic product owns Git snapshot, inventory, policy-driven accounting, neutral versioned contracts, canonical JSON, report assembly/publication, and the CLI. A profile supplies configuration and compatible vocabulary; it cannot make a core evidence result authoritative.

### First consumer example

CNSIC is the named first consumer/profile. Its OpenSpec and CHANGELOG conventions, forecast vocabulary, compatibility schemas, review integration, legacy carveout compatibility, and current `400` changed-lines default remain CNSIC concerns. They are not generic-core policy.

CNSIC W1/W2 at merge `a0fc7b26ff8a0e0a61baa586b32c46841611c806` is the immutable Branch by Abstraction predecessor oracle. It supplies compatibility evidence only; CNSIC code and vocabulary are not a source layout or semantic dependency for the neutral core.

## Contract principles

1. **Neutral and versioned:** public documents and Go APIs declare a schema or contract version and avoid consumer-specific governance terms.
2. **Canonical:** equivalent valid inputs serialize to stable canonical JSON bytes; parsers reject non-canonical representations where byte identity is part of the contract.
3. **Strict:** contracts validate exact types, closed shapes, discriminated variants, and immutable revision identity without silently normalizing inputs.
4. **Exclusive accounting:** a changed entry belongs to exactly one profile-defined category according to explicit precedence. Unknown paths default conservatively to the profile's production-equivalent category unless that profile specifies otherwise.
5. **No invented measurements:** binary, non-countable, unavailable, or untrusted data remains explicit rather than being converted into a line count or inferred value.
6. **Profile configuration:** every threshold, including changed-lines limits, is configurable per project/profile. Crossings become evidence or warnings only.
7. **Compatibility at the edge:** consumer compatibility schemas and legacy behavior live in adapters, not in neutral core contracts.

## Security and determinism expectations

- Resolve the selected revision pair to immutable native Git object IDs before deriving evidence. Support for both SHA-1 and SHA-256 repository object formats must preserve native identifiers without rehashing or transformation.
- Treat paths, refs, configuration, repository metadata, and profile input as untrusted boundaries. Preserve raw path bytes internally where required; never interpolate them into shell commands.
- Invoke Git with explicit argument vectors and a controlled environment. Disable external diff/text conversion, inherited configuration, replacement references, and other ambient behavior that could alter evidence.
- Derive committed-change accounting from immutable objects, not from index or working-tree state. Keep untracked inventory separate and explicitly scoped.
- Fail safely: malformed input, missing objects, repository races, or publication conflicts must not produce invented evidence or overwrite an existing artifact.
- Make local publication content-addressed, exclusive, and atomic when publication is implemented. A moving symbolic reference requires bounded revalidation before publication.
- Ensure report status and diagnostic output do not leak sensitive repository paths, command output, credentials, or environment values.

## Ownership and governance

Jose is the single initial owner for releases and compatibility decisions. He owns the decision to freeze, evolve, or deprecate public contracts and consumer compatibility vectors.

Project maintainers own their profiles, policy vocabulary, thresholds, and downstream interpretation. A profile change is not a generic-core governance change unless it alters a neutral contract or core Git/accounting semantic.

## Initial milestones

| Milestone | Exit outcome |
|---|---|
| 0. Freeze extraction inputs | CNSIC W1/W2 at merge `a0fc7b26ff8a0e0a61baa586b32c46841611c806` remains the immutable predecessor oracle and compatibility evidence boundary. |
| 1. Bootstrap design | Complete: Go 1.25.10, public/internal layout, typed configuration boundary, and test strategy are recorded without selecting module/publication identity or release automation. |
| 2. Neutral core extraction | Establish the Git snapshot, inventory, policy accounting, contracts, canonical JSON, and report-assembly boundaries behind real-Git verification. |
| 3. CLI and publication | Add deterministic CLI projections and safe local publication, with evidence-only statuses and bounded race behavior. |
| 4. First consumer profile | Introduce the CNSIC adapter/profile and compatibility vectors without importing CNSIC policy into the generic core. |
| 5. Distribution decision | Evaluate pinned Git consumption and choose packaging, versioning, and release automation only after the core and first profile are stable. |

## Acceptance outcomes

The product is ready for initial consumer adoption only when it can demonstrate all of the following:

- Given the same immutable revision pair and profile, independent runs produce the same canonical report bytes and projections.
- Real Git repositories exercise SHA-1 and SHA-256, hostile path bytes, renames, binary files, symlinks, gitlinks, missing objects, and index/working-tree isolation.
- Each accounting entry has one exclusive category; policy validation fails before totals are emitted; non-countable data is not converted into fictional line counts.
- Thresholds and ratios are visible as configurable evidence/warnings and never change the tool into an approval or blocking mechanism.
- Contracts reject malformed, ambiguous, and non-canonical boundary inputs, while compatibility concerns remain in the consumer adapter.
- Publication never overwrites an existing artifact and leaves no partial artifact when it fails.
- CNSIC can consume a frozen compatibility vector through its profile while another project can define different vocabulary and limits without forking the core.

## Bootstrap decisions and deferred decisions

The following decisions are authoritative for the first implementation work units:

- Go 1.25.10 is the target baseline. `go test ./...` is the primary reproducible gate; concurrency-bearing surfaces also require `go test -race ./...`; formatting is check-only.
- Public neutral APIs are at the repository/module root, the CLI is at `cmd/git-change-evidence/`, non-public shell/adapters are under `internal/`, and tests are co-located as `*_test.go`.
- Core configuration is immutable typed Go structs/interfaces. External format parsing is adapter-only, and TOML/YAML/JSON selection remains deferred.

The following are intentionally unresolved and MUST NOT be selected or implemented in this planning slice:

- Module/publication identity, packaging and distribution model, including whether initial consumption is only a pinned Git dependency or evolves to a package registry.
- Public versioning and compatibility/deprecation policy beyond the owner role defined above.
- Release automation, CI, publication channels, and signing/provenance mechanism.
- License and concrete configuration-file format.

## Immediate next step

Implement the neutral W1/W2 semantic boundary in its separately scoped chained PR. Preserve the CNSIC predecessor oracle, begin with RED tests, establish a module identity only when the owning work unit is authorized to do so, and do not begin CNSIC consumer cutover.
