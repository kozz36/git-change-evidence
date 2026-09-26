# Proposal: Repo-local agent guidance for published-preview GCE

## Intent

Give agents consuming GCE evidence a concise, reusable runtime skill for choosing the supported surface, retaining and validating evidence, reproducing it faithfully, and explaining technical limits without assigning delivery authority. Issue [#129](https://github.com/kozz36/git-change-evidence/issues/129) is separate post-preview work, not a condition on the already published Public Preview.

The parent-confirmed release record and this change's [exploration](exploration.md) establish that [v0.1.0-preview.1](https://github.com/kozz36/git-change-evidence/releases/tag/v0.1.0-preview.1) was published at **2026-09-19T07:09:00Z**, with `isDraft:false`, `isPrerelease:true`, target commit `939785dfe122d82022df54929d8f4578b56373fd`, and **six assets**. Issue #129 describes it as externally verified. This phase does not claim a new independent release verification. Local planned/unpublished wording in `AGENTS.md`, README, and configuration is stale status context, not a reason to reopen publication.

This is a clean change. `openspec/changes/post-preview-agent-skill/` remains intact as a historical failed attempt; its conclusions are not authority for this proposal.

## Scope

### In scope

- One repo-local, not-packaged agent skill for consuming, validating, reproducing, regenerating, and technically interpreting existing GCE evidence.
- Exact activation and nonactivation rules, including ambiguous and mixed requests.
- CLI-first guidance preserving canonical and pre-v1 legacy command compatibility and continued public root Go API and public `census` subpackage support.
- Evidence-only output, error interpretation, transparent recovery, and exact-input reproduction requirements.
- Minimal local references, representative scenario fixtures/checks, and documentation/skill-registry integration through a narrow additive root `AGENTS.md` skill registration in a later authorized implementation slice; an optional consumer-guide link supplements, not replaces, registration.
- Explicit separation of skill metadata from product, extractor, and wire-contract versions; no invented distribution or compatibility policy.

### Out of scope

- Any retroactive preview gate, publication re-verification prerequisite, or change to the published tag, release, assets, or target commit.
- Editing any file except this proposal during this phase. Later implementation may add the scoped root `AGENTS.md` skill registration while preserving all existing bytes and governance, including the stale preview paragraph. README edits, historical failed artifact edits, and stale-status reconciliation remain outside this change.
- CLI, Go API, schema, contract, or runtime behavior changes; a new generic evidence-generation CLI or policy layer.
- Approval, rejection, gating, blocking, merge, deployment, release, or delivery-authority decisions.
- Packaging, installers, external registries, broad distribution, signing/provenance, release automation, or automatic global skill installation.
- Multilingual census implementation. Its separate roadmap remains Python → JavaScript/TypeScript → Vue → Java → Rust.
- Spec, design, tasks, implementation, commit, or push in this proposal phase.

## Capabilities

### New capabilities

- `gce-evidence-agent-skill`: Repo-local runtime guidance covering activation boundaries, supported consumer surfaces, evidence preservation and interpretation, reproduction/recovery, local discoverability, and honest validation coverage.

### Modified capabilities

None. Existing census and other public evidence contracts remain unchanged; the skill consumes them rather than redefining them.

## Approach

Use the exploration's consumer workflow: **consume → validate → reproduce → interpret**. Keep the skill imperative and compact, with frontmatter and the ordered Activation Contract, Hard Rules, Decision Gates, Execution Steps, Output Contract, and References sections. Use the loaded bundled skill style guide because no repository `docs/skill-style-guide.md` was available. Link to local maintained consumer documentation instead of embedding parallel contracts.

### Activation contract

| Request | Skill behavior |
| --- | --- |
| Explicitly consume, validate, reproduce, regenerate, or technically interpret GCE evidence, including GCE Go AST census | Activate for the requested evidence task and supplied scope. |
| Bare mention of Go, AST, census, review, or GCE; generic programming; ordinary code changes; unrelated review | Do not activate merely on these words. |
| Author or modify GCE implementation, specifications, or contracts | Do not activate as an implementation or authoring workflow. |
| Solely request a pass/fail delivery verdict, approval/rejection, gate/block, merge, deploy, or release decision | Do not activate to provide that verdict; explain the evidence-only boundary. Technical validity observations remain distinct from delivery verdicts. |
| Request both GCE evidence work and a delivery verdict | Separate and perform only the authorized technical evidence subtask; leave delivery decisions with the consumer. |
| Evidence intent is plausible but family, scope, inputs, or authorization is unclear | Request the missing clarification before execution; do not invent scope or authority. |

### Supported surfaces and output

- Prefer `gce census go-ast`; preserve `git-change-evidence census-go-ast` as equivalent and supported throughout pre-v1. Require explicit absolute source root, query identifiers, `--`, and selected relative Go paths; no implicit discovery, Git revisions, or shell expressions for census.
- Keep the root Go API and `census.GoASTV1` available for applicable integrations. Use existing contract decoders and provenance validation, including `ValidateReportV1ResultProvenance` where applicable; do not imply all evidence families have a CLI assembler.
- Preserve canonical bytes including final LF, applicable digests, exact antecedents, raw-path Base64 identities, argv, stdout, stderr, and exit status. Validate against the declared existing contract rather than inventing a universal envelope.
- Describe census matches as exact textual syntax candidates in the supplied validated scope. Zero matches is not semantic absence, whole-repository completeness, safety, or approval. Selected files are not an atomic tree snapshot, and input bounds do not promise parser isolation.
- Return the technical observation, provenance/scope, validation or reproduction status, limitations, and any recovery performed or missing inputs. Never fill missing evidence with invented success or translate technical findings into delivery authority.

### Errors and recovery

| Observation | Required handling |
| --- | --- |
| Census exit 0, including zero matches | Interpret the complete document within its exact supplied scope. |
| Census exit 2 | Record invalid input, including command/path/parse failures; obtain a transparent correction. |
| Census exit 3 | Record unavailable content, including missing, symlinked, nonregular, or non-Linux access; do not bypass confinement. |
| Census exit 4 | Record exceeded bounds; do not silently narrow scope or raise limits. |
| Nonzero exit, interruption, partial write, malformed/truncated bytes, or stale/mismatched antecedents | Preserve failure observations; do not treat a prefix or missing result as empty success. |

Use command-specific documented exits for other commands, not the census mapping universally. Recovery requires transparent correction or authorized reacquisition; changed inputs produce **new evidence**, not reproduction. Avoid indefinite retries and unauthorized filesystem access. Census reproduction needs the same captured source bytes, selected paths, query, limits, and compatible extractor behavior; revision-bound evidence needs the exact immutable revisions and canonical antecedents. If unavailable, report that reproduction cannot be established.

### Documentation, distribution, and versioning

Propose `skills/gce-evidence/SKILL.md` as the checkout-local entry point, with supporting files only when needed. Plan a narrow additive registration in root `AGENTS.md` during a later authorized implementation slice, satisfying issue #129's documentation/skill-registry integration and the loaded skill-creator requirement to register project skills there. The entry should identify the skill, activation purpose, and local `SKILL.md` path. Preserve all existing bytes and governance; do not copy the stale publication claim into new guidance or silently rewrite the existing preview paragraph. No `AGENTS.md` edit occurs in this proposal phase.

An optional concise link from `docs/consumer-integration-guide.md` may improve human discoverability but is not a substitute for project skill registration. No separate maintained repo-local registry was identified by exploration; do not revive historical registries or machine-local paths. Root registration provides repository-local discovery, not automatic installation, external registry publication, or guaranteed activation by every agent runtime.

The skill is available only as repository content in this slice, not as an addition to the published preview assets or an installable package. Required `metadata.version` identifies skill content, not the CLI release, extractor, or schema version. The exact initial metadata value and any future skill-version compatibility/distribution policy remain unresolved; do not claim a release cadence or independently versioned distribution. These do not reopen the confirmed repo-local, not-packaged product scope.

### Tests and fixtures

Plan fixtures for positive evidence-use requests, negative generic/implementation/verdict requests, mixed requests, and clarification cases. Cover canonical/legacy CLI guidance, both public Go surfaces, zero matches, exits 2/3/4, partial writes, stale/mismatched inputs, raw-byte identity, exact-input reproduction, and no-authority language. Static checks can verify frontmatter, local reference resolution, the additive root registration and preservation of existing `AGENTS.md` bytes, and scenario expectations; they cannot prove actual agent behavior.

Use existing co-located root, CLI, and `census` tests as behavioral anchors, without changing runtime behavior for this skill. If Go checks are introduced later, follow the repository testing rubric and co-located layout. Report static inspection, executable contract tests, and actual agent evaluation as separate evidence levels. No tests or agent evaluation were run in this proposal phase.

## Affected areas

| Area | Intended impact in later phases |
| --- | --- |
| `skills/gce-evidence/SKILL.md` | Proposed new runtime skill; exact naming finalized in design. |
| `skills/gce-evidence/references/`, `skills/gce-evidence/assets/` | Optional local supporting references and scenario fixtures, only where necessary. |
| `AGENTS.md` | Later authorized implementation: append a narrow project skill registration, preserving all existing bytes and governance, including the stale preview paragraph. No edit in this phase. |
| `docs/consumer-integration-guide.md` | Optional supplemental discoverability link, not a substitute for root registration; retain primary consumer contract meaning. |
| `docs/census.md`, `openspec/specs/go-ast-census/spec.md` | Read-only contract anchors. |
| Root, `cmd/git-change-evidence/`, and `census/` tests | Existing behavior references; placement of any new checks remains a design decision. |
| README, failed change, published preview | Unchanged. |

## Risks

| Risk | Likelihood | Mitigation |
| --- | --- | --- |
| Stale preview text or failed artifacts contaminate new guidance | High | Use confirmed publication facts and primary contracts; preserve but do not adopt failed conclusions. |
| Overbroad activation or evidence interpreted as delivery authority | Medium | Explicit positive/negative/mixed scenarios and invariant no-authority output. |
| Duplicated guidance drifts from CLI/API contracts | Medium | Local authoritative references, compatibility examples, and narrowly scoped fixtures. |
| Repo-local registration is mistaken for automatic installation | Medium | Include the required root skill registration while distinguishing local discovery from automatic activation, packaging, and external registry publication. |
| Static checks are presented as demonstrated agent reliability | Medium | Report each validation evidence level separately; identify unexecuted evaluations. |
| Later docs, fixtures, and specs exceed the review budget | Medium | Forecast in tasks; under ask-on-risk, return delivery choice to the owner rather than infer chaining or a size exception. |

## Rollback plan

If later skill guidance proves misleading, remove or revert only the new skill/support files, the additive root `AGENTS.md` registration, and any newly added consumer-guide link. Preserve all pre-existing `AGENTS.md` bytes and unrelated edits. No runtime, evidence migration, or release rollback is required. Preserve this change's audit trail, the historical failed attempt, and all published preview objects. This proposal phase writes only `proposal.md`.

## Dependencies and remaining decisions

- Use the new exploration, parent-confirmed issue/release context, consumer guide, census documentation/specification, and unchanged evidence-only and public API rules.
- No new external service, runtime dependency, release operation, or stale-status cleanup is a prerequisite.
- Later design must finalize skill naming, minimal support files, check placement, explicit discovery wording, and initial skill metadata without inventing policy. Return any substantive distribution/versioning policy choice to the orchestrator rather than infer consent.

## Success criteria

- [ ] Later authorized implementation provides a compact repo-local entry point registered through a narrow additive root `AGENTS.md` entry; any consumer-guide link is supplemental.
- [ ] Registration preserves all existing `AGENTS.md` bytes and governance, including the stale preview paragraph; this proposal phase changes only `proposal.md`, and README remains unchanged.
- [ ] Activation, nonactivation, mixed requests, and missing-input clarification are explicitly scenario-covered.
- [ ] Canonical/legacy CLI and public root/`census` APIs retain their established meanings and support.
- [ ] Evidence preservation, zero-match limits, exits, partial failures, exact reproduction, and transparent recovery are covered without delivery verdicts.
- [ ] Local references and fixtures are checked, and validation claims distinguish static checks from executable and agent-evaluation evidence.
- [ ] Skill metadata does not imply packaging, new public compatibility policy, or modification of the published preview.
- [ ] Published release facts remain correct, historical failed artifacts remain intact, and multilingual work remains separate.
