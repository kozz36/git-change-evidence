# Proposal: Public Preview readiness

## Intent

Prepare an accurate, usable Public Preview entrypoint for `git-change-evidence`, with Apache-2.0 licensing, CLI-first integration, compatible command discovery, and owner-visible publication hygiene. The first release identity is `v0.1.0-preview.1`; this change prepares that identity without publishing a tag, release, or repository.

Authors, reviewers, maintainers, and coding-agent integrators need documentation that matches the implemented product and commands that can be used without guessing. The exploration found implemented Go APIs, machine contracts, census, tests, and CI alongside bootstrap-era documentation claiming those surfaces do not exist. This mismatch creates adoption friction and misleading expectations.

GCE remains evidence-only: it may report measurements, observations, and technical states, but never approve, reject, gate, block, merge, deploy, release, or assign delivery authority. Readiness findings are evidence for the owner, not product verdicts or publication authorization.

## Authority and inputs

- The owner's confirmed ruling selects Apache-2.0; CLI and machine contracts as primary with the Go API supported; canonical `gce census go-ast`; legacy `git-change-evidence census-go-ast` throughout pre-v1; `v0.1.0-preview.1`; English technical material; and a separate post-preview multilingual census program.
- `explore.md` provides the current-state gap, compatibility risks, and audit limitations. No research artifact exists in this change directory.
- `AGENTS.md`, `openspec/config.yaml`, `README.md`, and `docs/product-charter.md` establish current guidance and the contradictions to reconcile.
- `.atl/skill-registry.md` was read as a generated local-path audit input, not as deletion authority or a source of product policy.
- Existing architecture, neutral contracts, Go 1.25.10 baseline, adapter boundaries, and immutable CNSIC predecessor provenance remain intact.

## Scope

| Area | Proposed outcome |
|---|---|
| Apache-2.0 | Add the unmodified standard license text and make current-facing licensing statements consistent. Classify third-party attribution and legal/IP questions for owner review; do not assume the repository license settles third-party rights. |
| Supported surfaces | Present the CLI and machine-readable contracts as the primary integration path. Explicitly retain support for the existing root Go API and public `census` subpackage, without removal or implied deprecation. |
| Canonical command | Make `gce census go-ast` executable through a documented, reproducible local build/install path, not merely aspirational README wording. Prefer a real binary named `gce` rather than requiring an interactive shell alias. The later design must settle the build/name mechanism within current architecture constraints. |
| Pre-v1 compatibility | Retain `git-change-evidence census-go-ast` throughout pre-v1. Both census entrypoints must invoke equivalent behavior, preserving existing inputs, JSON bytes for equivalent inputs, schema identifier, diagnostics policy, and exit semantics. Other existing commands remain supported. |
| Preview identity | Prepare version reporting and current-facing preview documentation for `v0.1.0-preview.1`. Distinguish development/source builds from an actually published release; do not claim download assets or distribution channels exist. |
| English technical material | Use English for all newly authored or updated technical docstrings, comments, and artifacts. Audit current-facing material for consistency; classify historical or generated non-English material without silently rewriting provenance. Preserve literal fixtures, identifiers, and quoted evidence when translation would change meaning or behavior. |
| Governance reconciliation | Reconcile current governance, `README.md`, `AGENTS.md`, `openspec/config.yaml`, and current OpenSpec guidance with implementation and the owner ruling. Distinguish current decisions from historical bootstrap records, preserving the latter. Keep genuinely deferred matters deferred, including concrete configuration syntax and broader release/distribution automation. |
| Publication audit/hygiene | Inventory and classify current-tree publication surfaces, generated/local-path artifacts, provenance, potential sensitive data, license/attribution concerns, and audit coverage gaps. Perform only justified, non-destructive current-facing hygiene; uncertain confidentiality or legal findings go to the owner before remediation. |
| Agent-consumer guidance | Add a consumer-oriented guide distinct from contributor `AGENTS.md`, covering applicable immutable revision binding, digest preservation, independent regeneration, stdout/stderr and exit interpretation, stale/incomplete/unavailable evidence, and evidence-not-approval rules. Describe actual census scope and limitations rather than implying semantic analysis or whole-tree coverage. |
| Help/version | Provide useful explicit help and version discovery for both executable identities, including the canonical census hierarchy and legacy compatibility usage. Discovery must work without evidence inputs, avoid evidence-production side effects, and not contaminate normal machine-output streams. Exact UX and development-version handling belong to later specification/design. |

Agent-consumer guidance is included in this preview slice by the explicit requested scope, even though exploration allowed it as a near-term follow-up.

## Non-goals

- Multilingual census implementation: roadmap mention only, as a separate post-preview program.
- A capabilities schema or discovery protocol, canonical schema changes, or new evidence semantics.
- Removal or deprecation of the supported Go API.
- Repository visibility changes, tag creation/publication, release publication, or uploading release assets.
- Git history rewriting, automatic deletion of local-path surfaces, or rewriting archival OpenSpec evidence.
- CNSIC cutover, importing CNSIC-specific policy into the neutral core, or MEH integration.
- Broad packaging, signing/provenance, release automation, new platform support, or unrelated workflow redesign.

## Affected areas

| Surface | Expected later impact |
|---|---|
| `LICENSE`, `README.md`, `docs/product-charter.md` | License addition, truthful preview entrypoint, and explicit current/historical governance boundary. |
| `AGENTS.md`, `openspec/config.yaml`, current OpenSpec guidance/specs | Remove operational contradictions while retaining architecture, test discipline, provenance, and undecided owner choices. Active contract semantics remain unchanged. |
| `cmd/git-change-evidence/`, narrowly necessary internal CLI adapters, co-located tests | Canonical/legacy command dispatch and executable naming, help/version, and compatibility coverage; no effects move into the functional core. Final file layout is a design responsibility. |
| `docs/census.md` and a dedicated agent-consumer document | Actual invocation/install guidance, machine contract usage, error interpretation, Linux confinement and non-Linux unavailability, scope and non-guarantees. |
| Publication audit records | Classification of current-tree and local/generated surfaces, historical evidence, attribution, and uninspected owner-controlled surfaces, without implying legal clearance. |
| `.atl/skill-registry.md`, historical OpenSpec, local-path references, CNSIC provenance | Audit inputs only by default. A local path or provenance reference alone is not evidence of a secret and does not authorize deletion. |

These are proposed later changes, not files authorized for modification in this proposal phase. This phase writes only `openspec/changes/public-preview-readiness/proposal.md`.

## Publication audit boundary

The exploration's current-tree pattern scan found no common credential signatures, but that is not a comprehensive security or confidentiality finding. This proposal does not repeat the scan or claim additional coverage. Git history, unpublished refs, GitHub settings/secrets/metadata, generated or release assets outside the tree, and legal/IP clearance remain owner-controlled audit surfaces unless later inspected with explicit authority.

The audit should distinguish inspected evidence, actionable findings, benign local/generated context, preserved historical provenance, and not-inspected surfaces. `.atl/skill-registry.md` contains machine-local skill paths and generated metadata: classify its publication relevance and portability before proposing any treatment. Do not delete it merely because paths are local. Historical CNSIC references and the immutable predecessor merge remain evidence, not implied integration commitments.

New public instructions should use portable repository-relative paths and executable examples. Any finding requiring destructive remediation, confidentiality judgment, credential action, or external settings changes must return to the owner; none is authorized by this proposal.

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Renaming breaks agents or changes machine bytes | Keep both entrypoints on equivalent behavior and require parity evidence for successful and failing invocations, with no schema rename. |
| CLI-first positioning implies Go API abandonment | State API support explicitly and retain existing public API coverage. |
| Preview language overclaims maturity, security, or availability | Document observed behavior and limitations; distinguish source-build readiness from release publication. |
| Reconciliation erases historical evidence | Update current guidance and add historical/current distinctions; preserve archive records and CNSIC provenance. |
| Hygiene causes unjustified deletion or confidentiality exposure | Classify before remediation, report coverage gaps, and retain owner authority over sensitive or destructive decisions. |
| Help/version affects existing output or exit handling | Isolate explicit discovery paths and preserve ordinary command behavior and stream separation. |
| Review surface exceeds 400 changed lines | The standard license plus documentation, CLI, and tests creates credible implementation-size risk. Measure before implementation and pause for an owner delivery decision under `ask-on-risk`; do not silently chain or accept `size:exception`. |

The 400-line review budget is a delivery preference, not a generic-core threshold. No chain strategy or exception is selected here. No Git diff measurement or implementation estimate is claimed as observed evidence.

## Rollback

This phase adds a proposal only; withdrawing it requires no runtime rollback. Later readiness implementation must keep the legacy command usable while canonical dispatch and discovery are introduced. Before any external publication, defective implementation changes can be reverted through ordinary forward commits while retaining audit history, existing machine contracts, and supported APIs. Correct misleading current-facing claims rather than restoring known-false bootstrap statements.

Do not treat license removal as a revocation of rights already granted. Any licensing correction requires owner/legal review. History rewriting, unpublishing releases, or changing repository visibility is not a rollback mechanism authorized here.

## Success criteria

These are intended outcomes for later implementation, not claims established by this proposal.

1. Apache-2.0 text is present and current license statements agree with the owner ruling; attribution uncertainties are explicit.
2. README, governance, contributor guidance, config, and current OpenSpec consistently describe the implemented product, preview identity, supported Go API, and remaining deferred decisions without altering historical evidence.
3. A documented local build/install path produces a usable `gce census go-ast`; the legacy `git-change-evidence census-go-ast` remains usable throughout pre-v1 with equivalent evidence outputs and existing error/exit semantics.
4. Explicit help/version discovery works for both executable identities without evidence inputs or unintended effects, and preview version reporting identifies `v0.1.0-preview.1` without asserting a release has been published.
5. Existing canonical schemas, digest semantics, normal command output contracts, Go APIs, and evidence-only meaning remain unchanged.
6. English technical material and portable current-facing examples describe actual behavior, including census platform and syntax-only limitations; historical language/path exceptions remain classified rather than silently erased.
7. Agent consumers can find instructions for revision/digest handling, independent regeneration, stale or unavailable evidence, and exits without interpreting technical success as delivery approval.
8. An owner-visible publication audit distinguishes inspected findings from uninspected history/settings/assets/legal surfaces and records local-path/provenance classifications without unauthorized deletion or external action.
9. Later Go changes follow the configured RED → GREEN → TRIANGULATE → REFACTOR discipline, `go test ./...`, check-only formatting, and race testing for concurrency-bearing surfaces. Command parity and discovery evidence accompany the behavioral changes.
10. No excluded implementation, integration, visibility change, history rewrite, tag, or release publication occurs. Any implementation review-budget risk receives an explicit owner delivery decision before oversized work proceeds.

## Phase boundary and next decision

The confirmed product decisions are accepted without another interview; execution is auto but authorization is proposal-only. No specification, design, task planning, implementation, verification phase, or archive work is performed here. A parent-authorized later phase may refine executable naming and discovery details without reopening the owner's fixed compatibility commitments. The parent must resolve the identified delivery-size risk before implementation; this proposal grants no publishing or destructive-action consent.
