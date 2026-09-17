# Public Preview Readiness Specification

## Purpose

Define the truthful, evidence-only public-preview behavior and documentation for `git-change-evidence` while preserving established evidence contracts and pre-v1 compatibility.

## Requirements

### Requirement: Current governance is distinct from historical provenance

Current-facing governance, contributor guidance, configuration guidance, and public documentation MUST describe the implemented product and the owner-approved preview decisions truthfully. They MUST distinguish those current decisions from retained bootstrap-era, archived OpenSpec, CNSIC-provenance, and other historical records. Historical records, including literal local-path evidence and immutable predecessor references, MUST remain preserved and MUST NOT be silently rewritten merely because they are no longer current guidance.

The system and all updated guidance MUST preserve the evidence-only boundary: they MUST report technical evidence, measurements, and observations without approving, rejecting, gating, blocking, merging, deploying, releasing, or assigning delivery authority. Concrete configuration syntax, broader distribution, release automation, and other genuinely deferred matters MUST remain identified as deferred.

#### Scenario: Current guidance reconciles an obsolete bootstrap claim

- GIVEN a current-facing document or guidance record contains a bootstrap-era statement contradicted by the implemented module, APIs, CLI, tests, or workflows
- WHEN the public-preview update is applied
- THEN current guidance MUST state the current verified position and applicable preview decisions, while the historical source remains identifiable as historical provenance rather than deleted or represented as current authority

#### Scenario: Evidence does not become delivery authority

- GIVEN a consumer reads a preview-readiness document or receives an evidence result with an attention observation
- WHEN the consumer interprets that material
- THEN it MUST find no approval, rejection, gating, blocking, merge, deployment, release, or delivery decision attributed to `git-change-evidence`

### Requirement: Apache-2.0 is present without a legal overclaim

The repository MUST contain the complete, unmodified standard Apache License 2.0 text. Current-facing licensing statements MUST identify Apache-2.0 consistently with that text. The repository license and its documentation MUST NOT claim that third-party attribution, intellectual-property ownership, confidentiality, export, or other legal questions have been cleared merely by the presence of the license; unresolved legal or attribution matters MUST be identified for owner-controlled review.

#### Scenario: Repository license is identifiable

- GIVEN a reader inspects the repository's current licensing materials
- WHEN it looks for the applicable repository license
- THEN it MUST find the complete unmodified Apache License 2.0 text and current-facing statements that consistently identify Apache-2.0

#### Scenario: License scope is not overstated

- GIVEN attribution or intellectual-property status for a historical or third-party artifact has not been independently established
- WHEN current-facing preview materials describe the repository license
- THEN they MUST NOT represent that unresolved status as legally cleared or settled by Apache-2.0 presence alone

### Requirement: CLI is primary and public Go APIs remain supported

Current-facing integration guidance MUST present the CLI and its machine-readable contracts as the primary integration surface. It MUST also state that the existing public root Go API and public `census` subpackage remain supported, without removal, deprecation, or an implication that consumers must migrate away from them. The positioning MUST preserve the existing evidence-only meaning, schema identifiers, canonical-byte and digest semantics, and documented non-guarantees of established contracts.

#### Scenario: A Go API consumer retains a supported path

- GIVEN a consumer currently integrates through the root public Go API or the public `census` subpackage
- WHEN it reads the preview integration guidance
- THEN it MUST be able to determine that the CLI is primary while its existing Go API surface remains supported and is not deprecated or removed

### Requirement: Canonical and legacy census commands have pre-v1 parity

Throughout pre-v1, `gce census go-ast` MUST be the canonical documented census invocation, and `git-change-evidence census-go-ast` MUST remain a supported compatibility invocation. A documented reproducible local build or install path MUST provide a usable `gce` executable; an undocumented interactive shell alias alone MUST NOT satisfy this requirement.

For equivalent census inputs, both invocations MUST have equivalent successful and failing behavior: they MUST preserve accepted inputs, successful JSON bytes, the existing census schema identifier, diagnostics policy, and technical exit semantics. The command rename MUST NOT create a new evidence schema or reinterpret existing census evidence. Other existing commands MUST retain their supported behavior.

#### Scenario: Equivalent successful census invocation

- GIVEN equivalent valid census options and paths are supplied through `gce census go-ast` and `git-change-evidence census-go-ast`
- WHEN each invocation completes successfully
- THEN both MUST emit byte-identical existing census JSON and the same successful technical exit outcome

#### Scenario: Equivalent failing census invocation

- GIVEN equivalent invalid or unavailable census inputs are supplied through both census invocations
- WHEN each invocation fails
- THEN both MUST preserve equivalent diagnostics policy and the same existing technical exit classification without emitting a new schema or partial successful evidence

### Requirement: Help and version discovery are side-effect-free

Both executable identities MUST provide explicit help and version discovery, including help for the canonical census hierarchy and the retained legacy census invocation. A successful explicit help or version request MUST require no evidence inputs, produce no evidence document, make no filesystem, Git, publication, or delivery decision side effect, write only its requested discovery text to standard output, write nothing to standard error, and exit successfully.

Normal evidence-producing and error command paths MUST retain their established standard-output, standard-error, JSON, diagnostics, and technical-exit behavior. Discovery text MUST NOT be interleaved with or otherwise contaminate a normal machine-readable evidence stream.

#### Scenario: Canonical help is discoverable without evidence inputs

- GIVEN a user invokes explicit help for `gce`, `gce census`, or `gce census go-ast` without census evidence inputs
- WHEN the help request is processed
- THEN the command MUST write only applicable help text to standard output, write nothing to standard error, exit successfully, and create, publish, or alter no evidence or repository state

#### Scenario: Version identity is separately discoverable

- GIVEN a user invokes explicit version discovery through either executable identity
- WHEN the request is processed
- THEN the command MUST write only version information to standard output, write nothing to standard error, exit successfully, and leave normal machine-readable command output behavior unchanged

### Requirement: README and agent-consumer guidance describe actual use

The README MUST be a concise English public entry point grounded in implemented behavior and MUST direct contributors to `AGENTS.md` and agent consumers to a dedicated consumer-oriented guide. The consumer guide MUST be distinct from contributor guidance and MUST explain applicable immutable revision binding, digest preservation, independent regeneration, standard-output versus standard-error interpretation, technical exit interpretation, and how to handle stale, incomplete, unavailable, or syntax-only evidence without treating it as approval.

Current-facing examples MUST use portable repository-relative paths and executable invocations. Census guidance MUST accurately state its supplied scope, syntax-only nature, Linux confinement behavior, and non-Linux unavailability; it MUST NOT imply whole-tree coverage, semantic resolution, security certification, sandboxing, or delivery authority.

#### Scenario: Agent consumer can independently interpret evidence

- GIVEN an agent consumer follows the dedicated guidance for a revision-bound evidence document or census result
- WHEN it needs to reproduce, assess freshness, or interpret an unavailable or nonzero technical outcome
- THEN it MUST find instructions for preserving digests, independently regenerating applicable evidence, interpreting streams and exits, and escalating its own delivery decision rather than treating the evidence as approval

### Requirement: Publication audit records evidence and owner-controlled stops

The public-preview audit record MUST classify inspected current-tree publication surfaces, generated or local-path material, historical provenance, potential sensitive-data or attribution concerns, and all uninspected audit surfaces. It MUST distinguish observed evidence, actionable findings, benign or preserved context, and coverage gaps. In particular, Git history, unpublished refs, external settings and secrets, generated or release assets outside the inspected tree, and legal or IP clearance MUST remain explicitly uninspected unless actually reviewed.

The audit record MUST identify mandatory owner-controlled human decision points before any later repository-visibility change, tag, release, asset publication, destructive remediation, credential action, or legal/confidentiality determination. Recording these stops MUST NOT make the tool approve, reject, block, or otherwise control those actions. A local path, generated metadata, or historical CNSIC reference alone MUST NOT be treated as proof of a secret or as authorization to delete or rewrite material.

#### Scenario: Audit discloses coverage boundaries

- GIVEN the audit has inspected the current repository tree but not GitHub settings, Git history, or external release assets
- WHEN it records publication-readiness evidence
- THEN it MUST classify the inspected findings separately from those uninspected surfaces and MUST NOT claim comprehensive security, confidentiality, or legal clearance

#### Scenario: A sensitive remediation question is found

- GIVEN audit evidence raises a possible credential, confidentiality, attribution, or destructive-cleanup concern
- WHEN the concern is recorded
- THEN the record MUST identify an owner-controlled human decision point and MUST NOT delete material, alter external settings, or report automated approval or rejection

### Requirement: New technical material is English with historical and literal exceptions

All newly authored or updated technical documentation, code comments, and OpenSpec artifacts in this change MUST be English. Current-facing technical examples MUST avoid machine-local paths. Historical, archival, generated, or quoted evidence MAY retain its original language or literal form when translation or normalization would alter provenance, identifier spelling, fixture bytes, quoted evidence, or behavior; such material MUST be classified rather than silently rewritten. This requirement MUST NOT introduce multilingual census implementation.

#### Scenario: Historical literal evidence is preserved

- GIVEN an archival artifact contains a non-English phrase, a machine-local path, or a literal identifier whose translation would alter historical evidence
- WHEN current-facing English material is added or updated
- THEN the historical artifact MUST remain preserved and distinguishable from the new material, while the new material uses English and portable examples

### Requirement: Preview identity is truthful and non-publishing

Current-facing preview documentation and successful version discovery MUST identify the planned first preview identity as `v0.1.0-preview.1`. They MUST distinguish a source or development build from a published distribution and MUST NOT claim a download asset, tag, release, package, repository visibility change, or publication exists unless it has separately occurred under owner control.

This change MUST NOT create or publish a tag or release, upload assets, change repository visibility, or perform release publication.

#### Scenario: Preview version is reported without a release claim

- GIVEN a user reads current preview documentation or requests the program version from a source build
- WHEN it encounters `v0.1.0-preview.1`
- THEN the material MUST identify it as the preview identity without asserting that a tag, release, downloadable asset, or distribution channel has been published

### Requirement: Existing evidence contracts and excluded work remain unchanged

This change MUST preserve all existing evidence schemas and their established public behavior, including Policy V1, Inventory V1, Result V1, Report V1, and the existing Go AST census schema. It MUST preserve canonical bytes, digest bindings, closed contract shapes, evidence-only semantics, supported public Go APIs, normal CLI machine-output contracts, and documented technical exit behavior except for the explicitly added compatible discovery and command paths.

The change MUST NOT implement multilingual census, a capabilities schema or discovery protocol, canonical evidence-schema changes, supported Go API removal, repository visibility changes, tag or release publication, history rewriting, CNSIC cutover, CNSIC-specific policy in the neutral core, or MEH integration.

#### Scenario: Existing evidence remains contract-compatible

- GIVEN an existing consumer validates a pre-existing Policy V1, Inventory V1, Result V1, Report V1, or Go AST census document
- WHEN public-preview readiness behavior is added
- THEN the consumer MUST continue to observe the established schema identifier, canonical and digest semantics, and evidence-only meaning without a migration to a new evidence contract

#### Scenario: An excluded request is proposed during implementation

- GIVEN an implementation request would add multilingual census, a capabilities schema, API removal, repository visibility change, tag or release publication, history rewrite, CNSIC cutover, CNSIC-specific neutral-core policy, or MEH integration
- WHEN that request is evaluated as part of public-preview readiness
- THEN it MUST be treated as out of scope and no such behavior or change MUST be supplied by this change
