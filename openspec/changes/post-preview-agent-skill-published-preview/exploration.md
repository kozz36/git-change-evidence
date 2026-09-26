# Exploration: post-preview GCE agent skill

## Scope and publication facts

This is a clean exploration for GCE issue [#129](https://github.com/kozz36/git-change-evidence/issues/129), explicitly authorized after publication. The issue describes the Public Preview as “externally verified.” The live release record, checked 2026-09-25 with `gh release view`, is [v0.1.0-preview.1](https://github.com/kozz36/git-change-evidence/releases/tag/v0.1.0-preview.1): `isDraft=false`, `isPrerelease=true`, `publishedAt=2026-09-19T07:09:00Z`, target commit `939785dfe122d82022df54929d8f4578b56373fd`, and six published assets. These supplied live facts supersede neither historical records nor public contracts, but they contradict the stale preview-status text in the current `AGENTS.md` and README. A future authorized documentation change may reconcile that conflict; this exploration must not rewrite it.

The prior attempt at `openspec/changes/post-preview-agent-skill/` is a failed, nonauthoritative attempt and remains wholly intact. Its proposal, exploration, design, and spec include the earlier “planned/unpublished” status and design choices made under that incorrect status; those are factual reasons not to carry its conclusions forward. This new change has a distinct name so the failed record is preserved. Design inputs here come from issue #129, the live release facts above, current project policy, and primary public consumer/CLI/API contracts—not from the failed draft. This phase produces exploration notes only; no preview gate, retrospective approval, or release action is in scope.

## Consumer surfaces and contract anchors

- `docs/consumer-integration-guide.md` declares the CLI and machine-readable contracts primary, while keeping the public root Go API and public `census` subpackage supported. It frames use as consume → validate → reproduce → interpret and explicitly reserves workflow decisions for consumers.
- `docs/census.md` specifies canonical `gce census go-ast` and supported pre-v1 legacy `git-change-evidence census-go-ast`, required explicit root/query/paths, no Git or discovery, Linux confinement, limits, wire output, syntax-only scope, and exits 0/2/3/4.
- `openspec/specs/go-ast-census/spec.md` is the current census behavior anchor. Root inventory/result contracts and the consumer guide ground other evidence families and public API validation. CLI-first skill guidance must not deprecate the public Go surfaces or invent a parallel contract.
- `README.md` and `AGENTS.md` state the primary consumer surface and evidence-only boundary, but each currently repeats stale “planned/unpublished” preview language. Their stale assertions are noted, not copied as truth or modified here.
- CodeGraph exploration located the canonical command normalization and legacy implementation in `cmd/git-change-evidence/main.go` / `census.go`, CLI tests in the same command package, and public `census.GoASTV1` plus its API tests in `census/`. No implementation change is proposed.

## Trigger and non-trigger boundary

The potential skill is for an agent explicitly asked to consume, validate, reproduce, regenerate, or technically interpret GCE evidence. Census usage is a prominent case, but activation should be intent-based, not keyed on a bare mention of “Go”, “review”, “census”, or “GCE.” It should help establish requested evidence family and scope, identify authoritative local contracts, preserve raw observations, validate and reproduce only with required original inputs, and state technical findings and limits.

Non-triggers include generic Go/AST help; ordinary code changes; GCE implementation, contract, or specification authoring; unrelated review; and requests solely for a pass/fail, approve/reject, gate/block, merge, deploy, or release verdict. For mixed requests, any technical evidence subtask can be separated from consumer-owned delivery choices. Clarification is appropriate when the requested family, scope, authorization, or inputs are ambiguous; the skill must not expand scope by guessing.

## Evidence-only interpretation and recovery

The skill must inherit GCE's evidence-only charter: it reports technical states, measurements, and observations and never assigns delivery authority. Census matches remain exact textual syntax candidates in the validated supplied scope, not semantic identity, whole-repository completeness, safety, or approval. Zero matches means no candidate for that query in that scope. Fixed input bounds do not promise parser resource isolation.

Preserve canonical output bytes (including final LF), applicable digests and exact antecedents, raw-path Base64 identities, argv, stdout, stderr, and exit status as separate observations. A nonzero exit, interruption, write-prefix, malformed/truncated document, unavailable input, or stale/mismatched antecedent is not empty success. Recovery should record the failure, make only a transparent correction or authorized reacquisition, label changed inputs as new evidence rather than reproduction, and avoid silent narrowing, limit changes, confinement bypass, or indefinite retries. Exact census reproduction requires the same selected source bytes, query, paths, limits, and compatible extractor behavior; revision-bound APIs require their exact immutable antecedents and revisions.

## Repo-local skill, documentation, tests, and registry questions

Repository inventory found no existing `skills/` tree or maintained repo-local skill registry. Existing contributor discovery is `AGENTS.md`; consumer contract documentation lives under `docs/`. A likely direction for later design is a checkout-relative skill and a concise `AGENTS.md` discovery entry linking to local maintained contracts, but this is not finalized by exploration. Do not revive generated historical registries, machine-local paths, or external registry assumptions. Initial scope is repo-local and not packaged; repository presence is not a new release asset, install channel, publication event, or compatibility/versioning policy.

Future planning should tie scenario coverage to the actual activation and safety requirements: positive evidence-use prompts; negative generic/programming/implementation/review/verdict prompts; canonical/legacy CLI compatibility; continued public root and `census` API support; zero-match meaning; exits 2/3/4 and incomplete writes; stale/mismatched inputs; exact-input reproduction; and explicit no-authority language. Existing co-located tests in `cmd/git-change-evidence/`, root, and `census/` are behavioral anchors. New static/fixture checks may verify local references and trigger examples, but must not claim actual agent reliability; agent execution, if performed later, is a separate evidence level. Documentation should link to the primary consumer guide and census/API contracts rather than duplicate them. No tests or agent evaluation were run for this exploration.

## Out of scope and open planning questions

Out of scope: retroactive preview gate or re-verification; changing the live release/tag/assets; rewriting `AGENTS.md`, README, or the failed attempt; changing CLI/API/contracts or schemas; multilingual census; semantic analysis or automatic discovery; packaging, external distribution, installers, signing, or release automation; and any delivery verdict.

For subsequent planning, establish precise trigger wording and mixed-request behavior; select the repo-local skill path and discovery/registration mechanism; decide whether reference material or fixtures are needed and what checks can honestly validate them; and scope only the docs/tests/registry changes needed for this consumer guidance. Treat the stale preview-status conflict as a separate owner-directed documentation issue, not as a prerequisite or retroactive gate for #129.
