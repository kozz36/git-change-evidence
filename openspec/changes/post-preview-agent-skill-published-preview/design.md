# Design: Repo-local GCE evidence consumption skill

## Technical approach

Add `skills/gce-evidence/SKILL.md` as a compact runtime instruction contract for **consume → validate → reproduce → interpret**. Keep detailed recovery guidance in one local reference and representative prompt/observation fixtures in one local asset. A later authorized implementation appends a narrow root `AGENTS.md` registration. No Go production code, public contract, release object, or distribution mechanism changes.

This design implements the corrected [proposal](proposal.md) and [skill specification](specs/gce-evidence-agent-skill/spec.md), using the new [exploration](exploration.md). This phase writes only this file; it does not create tasks or implement the skill.

### Authority and publication baseline

The supplied live release record establishes `v0.1.0-preview.1` as **published**, at `2026-09-19T07:09:00Z`, with `isDraft:false`, `isPrerelease:true`, six assets, and target commit `939785dfe122d82022df54929d8f4578b56373fd`. Issue #129 describes the preview as externally verified. These are inherited observations, not an independent verification performed by this design phase. Preserve the published tag, release, assets, and target commit; no re-verification prerequisite or retroactive preview condition is introduced.

Existing planned/unpublished assertions in `AGENTS.md`, README, and configuration are stale status context. Preserve them without importing them into new guidance. The failed draft at `openspec/changes/post-preview-agent-skill/` exists in the original authoring worktree, not in the reviewed Git tree; this PR does not add, edit, or delete that path or adopt its conclusions as design authority. Stale-status reconciliation is a separate owner-directed change.

### Inspected contract anchors

| Anchor | Design consequence |
| --- | --- |
| `docs/consumer-integration-guide.md` | Preserve bytes, antecedents, surface compatibility, and consumer-owned decisions. |
| `docs/census.md`; `openspec/specs/go-ast-census/spec.md` | Explicit input scope, Linux confinement, syntax-only meaning, bounds, output, and census-specific exits. |
| `cmd/git-change-evidence/main.go` | `normalizeCensusArgs` maps canonical census arguments to the retained implementation before CWD/Git-dependent dispatch. Source-build version text is not release provenance. |
| `cmd/git-change-evidence/census.go` | Required flags and paths; fixed default CLI limits; identical captured bytes feed inventory and extractor; a short or failing output write returns exit 3. |
| `census/query.go`; `census/go_ast.go` | Public value API, explicit positive limits, exact validated scope, defensive ownership, no semantic lookup. |
| `contract.go` | `ValidateReportV1ResultProvenance` validates canonical Report bytes and the retained Result/Policy/Inventory binding, not delivery readiness. |
| `go.mod`; `openspec/config.yaml` | GCE Go module, Go 1.25.10, co-located Go tests, existing testing rubric. Not `packages/coding-agent`. |

The supplied worktree root was used directly. Shell and CodeGraph tools were unavailable in this executor; targeted file reads and limited source searches substituted for a fresh graph query. No tests, shell commands, or agent evaluations ran in this phase.

## Architecture decisions

### 1. Keep the skill outside the evidence implementation

**Choice:** Repository-local Markdown guidance only, with no executable helper, new Go API, schema, dependency, configuration syntax, or runtime hook.

**Rationale:** GCE already has the required consumer surfaces. The skill coordinates their use; it must not move filesystem/process effects into the Functional Core or introduce policy through an adapter. Existing Functional Core / Imperative Shell and Hexagonal boundaries remain unchanged.

**Rejected:** A generic evidence assembler CLI, a new agent integration package, automatic file discovery, or a policy engine. These would expand scope and risk changing neutral contracts.

### 2. Use a compact entry point with two purpose-specific supports

**Choice:** Three new files: `SKILL.md`, `references/consumer-workflow.md`, and `assets/scenarios.md`. The reference maps evidence families to maintained local contracts and records recovery decisions. The asset supplies labeled scenario inputs and expected/forbidden observations, not executable automation or replacement wire schemas.

**Rationale:** A single long skill would exceed the skill style budget; copying full consumer documentation would create drift. One reference and one fixture document are sufficient. Do not add tutorials, example Go programs, binary fixtures, goldens, or a checker framework without a demonstrated later need and authorization.

**Style:** The repository has no `docs/skill-style-guide.md`; use the loaded skill-creator bundled guide. Required section order is Activation Contract, Hard Rules, Decision Gates, Execution Steps, Output Contract, References. Target 180–450 body tokens; recommended maximum 700, hard maximum 1000. Keep the evidence-only boundary, intent routing, and failure/reproduction constraints in the entry point, not only in linked detail.

### 3. Register by appending, not rewriting governance

**Choice:** In later implementation, append this standalone discovery entry to root `AGENTS.md`:

```markdown
## Repo-local skills

- `gce-evidence`: For explicit requests to consume, validate, reproduce, regenerate, or technically interpret GCE evidence, read `skills/gce-evidence/SKILL.md`. Technical evidence only; no delivery authority. Repository-local guidance, not automatic installation.
```

**Rationale:** Project skills require root registration under both skill-creator and the new specification. Preserve every pre-existing byte, including the preview paragraph; capture the pre-edit bytes and require the resulting file to have that exact prefix plus only the reviewed append. Do not reflow, normalize line endings, or reconcile unrelated guidance. If the registration already exists when implementation starts, inspect it rather than duplicate it.

**Rejected:** A separate registry, machine-local installation, package-level `AGENTS.md`, or a consumer-guide link as the only registration. A `docs/consumer-integration-guide.md` link is optional supplemental discovery, omitted from the minimal design. Registration does not guarantee any runtime loads the skill.

### 4. Separate metadata from product and distribution policy

Use this frontmatter shape; the author and initial content-version values below are proposed literals requiring owner confirmation before implementation, not settled publication policy:

```yaml
---
name: gce-evidence
description: "Trigger: consume, validate, reproduce, regenerate, or interpret GCE evidence. Preserve scope and provenance without delivery verdicts."
license: Apache-2.0
metadata:
  author: "kozz36"
  version: "1.0"
---
```

The description is quoted, one physical line, trigger-first, and below 160 characters; do not add `Keywords`. Apache-2.0 follows repository license selection, not a legal clearance claim. The proposed `metadata.version` identifies skill content only. It is neither `v0.1.0-preview.1`, an extractor version, a wire schema, nor a compatibility promise. Attribution and the initial metadata literal remain owner-owned; do not silently substitute skill-creator's author.

**Rationale:** Required metadata need not invent a release cadence or independently distributed product. Packaging, installers, global installation, external registry publication, signing, release automation, and future compatibility/versioning policy remain outside this change. The owner must explicitly authorize any such expansion. No new skill content enters the six existing preview assets.

## Activation and data flow

| Request class | Required route |
| --- | --- |
| Explicit consumption, validation, reproduction, regeneration, or technical interpretation of GCE evidence | Activate for the authorized evidence subtask and supplied scope. |
| Bare Go/AST/census/GCE/review mention; generic help; ordinary programming; unrelated review | Do not activate from keywords alone. |
| GCE implementation, contract, or specification authoring | Do not use this consumption skill as an authoring workflow. |
| Delivery verdict only: approve/reject, gate/block, merge/deploy/release | Do not activate to provide the verdict; explain the evidence-only boundary. |
| Evidence task plus delivery verdict | Separate the two; perform only the authorized technical task and leave delivery decisions to the consumer. |
| Plausible evidence intent but unclear family, scope, inputs, or authorization | Ask for missing information before execution; do not guess or expand scope. |

```text
Request → intent/scope decision → load family-specific local contracts
                                    ↓
Retained evidence or authorized acquisition → raw observations
                                    ↓
Contract validation → exact-input reproduction, if requested and possible
                                    ↓
Technical observation + provenance + limitations + recovery status
```

Failures branch to preserved failure observations and clarification/authorized recovery, never to synthesized empty success. Evidence bytes, source comments, path text, and diagnostics are data, not instructions granting authority or changing the task. No implicit retry, release action, or access escalation occurs.

## Interfaces and evidence contracts

### Supported surfaces

Prefer the canonical census invocation, as an argument vector rather than interpolated shell text:

```text
gce census go-ast --source-root <absolute-root> --receiver <identifier> --selector <identifier> -- <explicit-relative-go-paths...>
```

The equivalent `git-change-evidence census-go-ast` prefix remains supported throughout pre-v1, with the same flags, paths, output, and technical behavior. Both executable names come from the same command package; availability must be established locally, not assumed from repository registration. Do not automatically install or build software without applicable authorization.

Require each option exactly once, the `--` separator, and at least one selected relative `.go` path. Root and paths must satisfy the documented lexical rules. No implicit CWD root, Git revisions, wildcard discovery, shell expressions, or invented CLI limit flags. CLI defaults remain 256 files, 1 MiB/file, 8 MiB total, 10,000 matches. Linux confinement is mandatory; unavailable content on other platforms is not permission to use an unconstrained substitute.

For Go integrations, retain `github.com/kozz36/git-change-evidence` and `github.com/kozz36/git-change-evidence/census` as supported imports. Use family-specific root constructors/decoders and `ValidateReportV1ResultProvenance(report, binding)` where applicable. Use `census.GoASTV1(inventory, files, query, limits)` for supplied-byte census, with a valid `InventoryDocumentV1`, matching complete raw path/content records, `QueryV1`, and positive `Limits`. The API permits a valid explicitly empty inventory; the CLI requires selected paths. Do not erase this distinction.

Not every family has a CLI assembler. `DecodeCanonical` is not a universal census decoder, and the skill adds no universal evidence envelope. Local CLI `publish` stores evidence; it does not publish a product release.

### Capture, validation, and reproduction

| Stage | Required retained inputs and observations |
| --- | --- |
| Capture | Exact canonical bytes including final LF; applicable published digests and their source; argv; stdout, stderr, and exit status separately. For API use retain typed inputs/result or error; do not fabricate process observations. |
| Raw identity | Preserve standard Base64 raw-path identities; never replace them with normalized/display paths or assume UTF-8. Preserve source bytes and exact immutable native Git IDs where applicable. |
| Validate | Identify the declared family/schema; use its existing decoder/contract. Check applicable digests over unmodified canonical bytes and retained antecedent bindings. JSON parseability or digest syntax alone is insufficient. |
| Census reproduction | Same captured selected file bytes, raw paths, query, effective limits, and compatible extractor behavior; compare resulting output bytes and applicable digests. Record extractor identity. The emitted hashes do not recover missing source bytes. |
| Revision-bound reproduction | Exact immutable base/head IDs, canonical Policy/Inventory/Result antecedents and digests, and the same supported API inputs; compare canonical bytes/digests. Moving refs or similar checkout state are not substitutes. |
| Missing or changed inputs | Report that reproduction cannot be established. Authorized regeneration on changed inputs yields new evidence, even if counts happen to match. |

Record the absolute source root with the invocation because the census document deliberately omits it. Do not claim that a later filesystem reread captured the original bytes unless verified against their recorded identities; files may have been captured at different instants. Do not introduce source copying, snapshots, or writes beyond the consumer's authorization.

Census matches are exact textual syntax candidates, including shadowed spellings and ambiguous indexed forms, within the validated supplied scope. Zero matches means no candidate for that exact query in that scope, not semantic absence, repository-wide completeness, safety, or approval. Input bounds do not guarantee parser heap, stack, time, or process isolation.

### Failure and recovery decisions

| Observation | Report and permitted next action |
| --- | --- |
| Census exit 0 | Validate the complete document; zero matches remains scoped success. Normal stdout is one JSON document with final LF and stderr is empty. |
| Census exit 2 | Invalid input, including malformed command/path or parse failure; retain `invalid input` diagnostic and seek a transparent correction. |
| Census exit 3 | Content unavailable, including missing/symlinked/nonregular sources or non-Linux access; preserve `content unavailable`. An output-write failure also maps here, so do not diagnose missing source from the exit alone. |
| Census exit 4 | Bounds exceeded; retain `resource bound exceeded`. Do not silently drop files, split the scope, or raise limits. |
| Nonzero exit, interruption, truncated/malformed/noncanonical bytes, or write prefix | Retain raw streams and observed failure; a prefix is not a successful document or empty result. Do not infer a completed status after interruption. |
| Stale/mismatched antecedents or unavailable originals | Preserve the mismatch or missing-input observation; request exact inputs or authorization for new acquisition. Never repair provenance by substituting convenient antecedents. |

Use documented command-specific exits for other CLI paths; census meanings are not universal. Recovery must name the original failure, proposed correction, any changed inputs, authorization, and result. No confinement bypass, silent narrowing, unapproved limit change, or indefinite retry. If correction cannot preserve the original inputs, classify its output as new evidence, not reproduction. Report unresolved evidence honestly without issuing a delivery verdict.

### Agent response contract

Return a concise technical summary with: requested family/scope and surface; observations/provenance; validation performed and result; reproduction status and byte/digest comparison if attempted; limitations/missing inputs; failures and recovery performed or required. Link retained artifacts rather than substituting summaries for bytes. Distinguish not attempted, unavailable, mismatched, and reproduced; these are descriptive statuses, not a new machine-readable contract. End without pass/fail delivery language or workflow authorization.

### Local reference resolution

From `skills/gce-evidence/SKILL.md`, use `references/consumer-workflow.md` and `assets/scenarios.md`. From either supporting directory, repository-root links begin with `../../../`: for example `../../../docs/consumer-integration-guide.md`, `../../../docs/census.md`, and `../../../openspec/specs/go-ast-census/spec.md`. Include family/API pointers to `../../../contract.go`, `../../../census/query.go`, and the existing test anchors where useful. Resolve each path from the containing file; no machine-local paths, external URLs as primary references, or failed-change links as runtime authority.

## File changes in later authorized implementation

| File | Action | Responsibility |
| --- | --- | --- |
| `skills/gce-evidence/SKILL.md` | Create | Activation, hard constraints, routing, workflow, output, local references. |
| `skills/gce-evidence/references/consumer-workflow.md` | Create | Family/contract pointers, capture/reproduction distinctions, compact recovery table. |
| `skills/gce-evidence/assets/scenarios.md` | Create | Representative prompts/observations with required and forbidden behavior; validation instructions. |
| `AGENTS.md` | Append only | Narrow discovery registration; preserve all pre-existing bytes. |

No other implementation file is required. The optional consumer-guide link is not part of this minimal file set. Existing source, tests, README, preview-status paragraphs, failed artifacts, published objects, and multilingual behavior remain unchanged. Python → JavaScript/TypeScript → Vue → Java → Rust remains a separate roadmap.

## Validation strategy

Choose static inspection plus scenario fixture review for the new Markdown surface, not a bespoke executable checker. Existing Go tests remain contract anchors. The default documentation-only rubric is `skip` for Go implementation testing, which does not remove the specification's static/scenario validation requirements. Any later authorized Go changes require RED → GREEN → TRIANGULATE → REFACTOR, co-located tests, `go test ./...`, and the exact check-only formatting command in `openspec/config.yaml`; race testing applies to concurrency-bearing changes. No new Go changes are designed here.

Each fixture in `assets/scenarios.md` has a stable ID, supplied prompt and observations, expected activation/clarification route, required observations, forbidden claims/actions, and local contract anchor. Synthetic observations are labeled; they are not claimed as captured execution. At minimum cover these independently assessable cases:

| Cases | Expected assertions |
| --- | --- |
| Explicit consume / validate / reproduce / regenerate / interpret | Each intent activates only for requested scope. |
| Generic Go/AST; ordinary edit; unrelated review; GCE authoring; bare keyword | No consumption workflow activation. |
| Verdict-only; mixed evidence/verdict; missing family/scope/input/authorization | Boundary-only response; technical subtask separation; clarification respectively. |
| Canonical and legacy commands; root API and public census API | Supported compatibility retained; no universal assembler or decoder invented. |
| Valid zero-match document | Syntax-only supplied-scope interpretation, not semantic absence or delivery approval. |
| Census exits 2, 3, 4; non-Linux/symlink; output-writer failure | Correct contextual interpretation and no bypass, automatic limit increase, or silent narrowing. |
| Nonzero with prefix; interruption; malformed/truncated bytes | Failure preserved; no synthesized empty success. |
| Final-LF mutation; raw non-UTF-8 path; stale Report binding | Byte/digest/path identity retained; mismatch exposed rather than normalized away. |
| Exact originals; changed checkout with same count; missing originals; authorized correction | Reproduction only with original inputs and comparison; otherwise unresolved or explicitly new evidence. |
| Stale preview assertion; request to install/publish skill or release | Preserve confirmed publication facts and owner boundaries; no release modification or implied distribution. |

| Evidence level | Later validation and honest reporting |
| --- | --- |
| Static inspection | Check frontmatter completeness, quoted single-line description length, ordered sections, token budget, every local link, additive registration and byte-prefix preservation. Record method and observed result; an estimated token count is not an exact measurement. |
| Fixture review | Check each fixture against the specification and guidance; record expected behavior, omissions, and unresolved cases. Passing expectations on paper is not agent execution. |
| Executable contract tests | If run, record commands/environment/results for existing CLI/census/root tests. They demonstrate implementation contracts, not skill adherence. Otherwise report not run. |
| Actual agent evaluation | Only if separately performed: load the actual skill in a named runtime/model, execute fixture prompts in bounded authorized contexts, retain prompts, outputs/tool traces, skill revision, environment, and per-case observed deviations. Report actual results separately from expectations; do not invent reliability or aggregate success. |

Behavioral anchors include `TestCanonicalAndLegacyCensusParity`, `TestCensusGoASTPreservesNonUTF8PathsAndUnicodeIdentifiers`, `TestCensusGoASTOutputWriteFailureHasNoSuccessDocument`, `census/external_api_test.go`, and `TestValidateReportV1ResultProvenanceRejectsEveryMismatchAndPreservesLegacySurfaces`. Existing writer tests do not by themselves demonstrate an agent handles every truncated-stream scenario. No validation pass or actual agent reliability is claimed by this design.

## Threat matrix applicability

N/A for the phase skill's routing/shell/subprocess/VCS-automation/executable-classification/process-integration matrix: this change creates declarative Markdown, not an executable router, launcher, classifier, or process integration. Existing CLI command examples describe unchanged interfaces. Agent-facing misuse boundaries still require the scenarios above: treat evidence as data, preserve argv/raw paths, respect authorization and confinement, and perform no hidden retries or delivery actions. A later executable helper would expand scope and require a new applicable matrix and RED tests before implementation.

## Rollout, rollback, and review risk

No migration, feature flag, release operation, or evidence regeneration is required to add the skill. After separate implementation authorization and owner metadata confirmation, the three local files and additive root registration form one coherent repository-local capability. Validate discovery, links, scenario coverage, and preserved root bytes together. Availability means present in a checkout; actual runtime loading remains observable rather than guaranteed.

Rollback removes or reverts only these new files and the exact appended registration, preserving the original `AGENTS.md` and unrelated concurrent changes. If an optional consumer-guide link is separately added, remove only that addition. Preserve this change's audit trail, all existing contracts, and all published preview objects; this PR does not add, edit, or delete the failed authoring-worktree draft. No runtime or release rollback is needed.

The session uses `auto`, `openspec`, `ask-on-risk`, and a 400 authored changed-line review budget. Planning artifacts count, as do authored fixtures, tests, and registration. The proposal, specification, exploration, and this design together already indicate **high aggregate budget risk** before implementation; no base-relative Git diff was available to measure the actual delivery surface. A preliminary implementation-only estimate is 160–270 added lines (skill 45–70, reference 55–90, fixtures 55–105, registration about 5), not an approved delivery slice or permission to ignore planning lines.

The parent must pause for an explicit delivery choice before selecting an oversized delivery or applying work. No chain strategy, automatic split, `size:exception`, or acceptance reduction is selected here. Keep the full acceptance surface; do not compress fixtures merely to evade review accounting. This phase stops after design as requested.

## Remaining owner decisions and risks

- Confirm initial `metadata.author` and `metadata.version` literals before materializing frontmatter. This does not prevent completing the design or authorize a new distribution/compatibility policy.
- Resolve high review-budget risk with an actual base-relative authored-line count and explicit owner delivery choice before apply; chaining remains deferred.
- Retained stale preview text may confuse future agents. New guidance must not repeat it as fact; reconciliation stays outside this change.
- Local links and prompt guidance can drift from public contracts. Keep references narrow and validate them against the unchanged source/test anchors.
- Static checks and Go tests cannot establish agent adherence. Actual evaluation remains unperformed unless separately evidenced.
