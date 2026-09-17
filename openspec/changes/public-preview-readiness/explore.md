# Public Preview readiness exploration

## Outcome

The current implementation is materially ahead of the public presentation, but it is **not ready for the stated Public Preview** without a focused documentation, licensing, release-identity, and CLI-compatibility slice. The evidence supports a small, reviewable readiness change; it must preserve the evidence-not-authority boundary and must not fold multilingual census into this change.

The owner ruling is the decision authority for this exploration:

- License: Apache-2.0.
- Primary surface: CLI and machine-readable contracts; supported Go API is secondary.
- Canonical command: `gce census go-ast`.
- Pre-v1 compatibility command: `git-change-evidence census-go-ast`.
- First version: `v0.1.0-preview.1`.
- English-only technical documentation, comments, and artifacts.
- Multilingual census is a separate post-preview roadmap program.

No readiness finding grants delivery, publication, merge, release, or other authority. GCE may emit evidence and technical statuses only; consumers own delivery decisions.

## Inputs and method

Read: `AGENTS.md`, `openspec/config.yaml`, `README.md`, `docs/product-charter.md`, `docs/census.md`, active OpenSpec specs, active `go-ast-census` materials, current CLI files, public contracts, and both workflows. CodeGraph was initialized and used for CLI/public-contract call paths. The non-authoritative input at `/home/kozz36/Downloads/preparar git-change-evidence para Public Preview.md` was read as requested.

The supplied execution environment exposes no Git command tool, so `origin/main@1f7e815` could not be directly read as a Git object. This exploration therefore models the predecessor gap from current-tree evidence and retained historical OpenSpec statements; it does not assert an unobserved commit diff. The historical materials explicitly record the earlier no-source/no-module/no-CI state, while the current tree contains `go.mod`, Go implementation, tests, workflows, and `census`.

Publication audit is limited to readable current-tree evidence. A pattern scan found no common credential/token/private-key signatures. It did find intentional CNSIC provenance and historical local paths in OpenSpec artifacts. Git-history-only secrets, unpublished refs, repository settings, and GitHub metadata cannot be established with the available read-only tools and require a later owner-controlled audit.

## Evidence snapshot

| Area | Current evidence | Readiness implication |
|---|---|---|
| Product presentation | `README.md` says implementation has not started and that source/module/CI do not exist. | Pre-preview blocker: it is materially false. |
| Product charter | `docs/product-charter.md` retains bootstrap tense and deferred license/CI/release decisions despite current implementation. | Pre-preview blocker: it cannot remain an apparently current authority without an explicit historical/current distinction. |
| Module/runtime | `go.mod` declares `github.com/kozz36/git-change-evidence`, Go `1.25.10`, and `golang.org/x/sys`. | README/charter must no longer claim absent module/source. |
| Public contracts | Root Policy/Inventory/Result/Report V1 documents are canonical, digest-bound, strict, and evidence-only; public `census` is syntax-only and bounded. | Describe actual contract scope and non-guarantees; do not advertise authority. |
| Census CLI | `cmd/git-change-evidence/census.go` implements only `git-change-evidence census-go-ast --source-root … --receiver … --selector … -- <paths...>`. It outputs `git-change-evidence.census-go-ast-cli/v1`. | Owner-required canonical `gce census go-ast` is absent; legacy command exists and must remain. |
| CLI UX | `main.go` has positional legacy analysis plus `project`, `publish`, and `census-go-ast`; it has no `--help`, `--version`, `gce` executable identity, `census` command group, or capabilities/schema discovery command. | Help/version and canonical command compatibility are pre-preview concerns; capabilities deserves separate contract review. |
| Determinism/diagnostics | Census succeeds with one JSON document and newline; its error paths produce generic stderr diagnostics and documented exits 2/3/4. Other commands have stdout/stderr separation and exit mapping in source. | Preserve these machine contracts; document only observed behavior. |
| Census confinement | Linux uses descriptor-relative `openat`, no-follow checks, regular-file checks, bounded reads, and generic diagnostics. Non-Linux returns unavailable. | Worth accurately documenting as an established trust boundary, not calling it sandboxing or a security certification. |
| CI | `go-pr-checks.yml` runs check-only gofmt, `go test ./...`, and `go test -race ./...`; actions are full-SHA pinned, credentials disabled, minimal read permissions, timeout and concurrency set. `pr-policy.yml` uses `pull_request_target` but checks out the trusted base SHA, not PR code. | No workflow blocker observed; retain this hardening. |
| License | No `LICENSE*` file was found; `AGENTS.md`, charter, and config contain obsolete deferred-license language. | Pre-preview blocker under the owner ruling. |
| Release/versioning | No observed release workflow, release metadata, binary packaging, or CLI version output. | `v0.1.0-preview.1` needs an explicit release/versioning plan and truthful installation path before publication. |

## Gap model: predecessor/current state to owner ruling

### Reconciled current facts

The retained bootstrap-era README/charter and `bootstrap-neutral-core-extraction` materials describe a planning-only repository with deferred module, CI, license, and release decisions. The current implementation supersedes that operational state: it has a Go module, root public APIs, `census` API, CLI, internal adapters, tests, and hardened PR workflows. Historical OpenSpec records remain useful provenance and must not be deleted or rewritten merely because they are stale.

The owner ruling supersedes remaining deferred presentation/release decisions for this preview only: Apache-2.0, pre-v1 version identity, CLI-first positioning, canonical/legacy command intent, and English-only public technical material. It does not authorize a change to evidence semantics, delivery authority, canonical V1 schemas, or historical artifacts.

### Required compatibility decision

The current code cannot satisfy `gce census go-ast`: it recognizes only a first argument of `census-go-ast`. The owner ruling requires both surfaces throughout pre-v1, with `gce census go-ast` canonical and `git-change-evidence census-go-ast` retained. A proposal must specify the executable naming/install model and a single dispatch compatibility design before implementation. It must also decide whether `gce` is an installed binary name, a documented alias, or both; documentation alone cannot make an absent command work.

Do not remove or reinterpret `git-change-evidence census-go-ast`. Do not introduce a new census JSON schema merely to rename the command unless an explicit compatibility analysis says the machine contract needs it.

## Classified work

### Pre-preview blockers

1. **Apache-2.0 licensing is absent.** Add the unmodified standard Apache-2.0 license text and reconcile current-facing documentation/configuration that still says license is deferred.
2. **The README is factually obsolete.** Replace the no-implementation/no-module/no-CI presentation with a concise English public entrypoint grounded in current behavior and OpenSpec contracts.
3. **The charter needs a current-status boundary.** Preserve bootstrap history and CNSIC provenance, but mark historical/deferred statements that the owner ruling or current tree supersedes. Do not silently rewrite archival OpenSpec.
4. **The canonical command is unimplemented.** Provide and document `gce census go-ast` while retaining `git-change-evidence census-go-ast` for pre-v1 compatibility. Preserve the existing machine-readable census document unless a separately justified contract change is made.
5. **Public CLI guidance is incomplete.** Add useful `--help` and `--version` behavior or explicitly establish an alternate, actually executable discovery path; document meaningful exits and non-interactive semantics. The owner-selected first release must identify itself as `v0.1.0-preview.1` without inventing packaging that does not exist.
6. **Publication hygiene needs an owner-visible audit.** The current scan found no common secret patterns, but history/settings/release assets were not directly inspectable. Before changing visibility, a human-controlled audit must cover Git history, GitHub secrets/settings, generated assets, and legal/IP review of historical material.

### Near-term work after preview or in a tightly bounded follow-up

1. Add an agent-consumer guide separate from contributor `AGENTS.md`: immutable revision binding, digest preservation, independent regeneration, exit interpretation, stale-evidence handling, and evidence-not-approval rules.
2. Evaluate capability/schema discovery separately. `capabilities --json` would freeze a new public machine contract; it is useful for agents but must not be added casually.
3. Clarify supported Go API versus preferred CLI integration in README and contract documentation without deprecating existing exported APIs.
4. Decide distribution and installation details for the preview release, release notes, tag process, repository description, and accurate GitHub topics.
5. Consider proportionate security scanning after ownership/tooling choices are explicit; preserve existing workflow hardening and do not add a large platform by default.

### Post-preview roadmap

1. **Multilingual census is a separate program.** It is not part of CLI renaming, Go AST census, or Public Preview readiness. It needs independent language scope, parser/contract/versioning, resource-bound, and documentation decisions.
2. Expand capability discovery only after a stable contract design and consumer need are demonstrated.
3. Evaluate packaging, signing/provenance, release automation, and additional platforms as separate owner-approved work.
4. Revisit first-profile/CNSIC cutover only under its existing parity and authority boundaries; do not import profile policy into neutral contracts.

## Public documentation and contract direction

The future README should lead with: reproducible change evidence for software and coding-agent workflows; current Public Preview/pre-v1 status; and the evidence-not-authority boundary. It should identify the CLI and machine documents as primary, point contributors to `AGENTS.md`, point consumers to a dedicated agent guide, and link detailed behavior to OpenSpec.

Only documented/implemented claims should appear:

- canonical/digest-bound Policy, Inventory, Result, and Report V1 documents;
- immutable revision and provenance bindings where the relevant contract establishes them;
- deterministic projections and local content-addressed publication where the current CLI/API provides them;
- syntax-only, supplied-scope `GoASTV1` candidates and the CLI's explicit filesystem scope;
- selected-file base64 path representation, source hashes, spans, defaults, and non-guarantees for census.

Do not claim a general source-tree snapshot, semantic call resolution, approval/gating, security certification, sandboxing, production maturity, packaging/release availability, or broad platform support. Census documentation currently accurately limits Linux source confinement and says non-Linux is unavailable; retain that clarity.

## Local/historical artifact handling

Do not delete or rewrite:

- `openspec/changes/archive/**` and historical bootstrap/go-ast materials;
- CNSIC predecessor references and immutable merge identity, which document provenance and a compatibility boundary;
- local paths embedded in historical verification/progress artifacts, including worktree and prior agent-session locations.

These paths are historical evidence rather than current public instructions. New public-facing documentation should avoid adding local paths. A human/legal review should decide whether any existing historical path or CNSIC reference is confidential before visibility changes; this exploration found no direct proof of confidential content and does not authorize removal or history rewriting.

## Risk register

| Risk | Evidence | Mitigation for later proposal |
|---|---|---|
| Canonical command breaks existing agents | Current parser only accepts `census-go-ast`. | Specify dual command dispatch and fixture-level compatibility before code changes. |
| README/charter overclaim or contradict | Both retain bootstrap claims disproved by the tree. | Ground every public statement in current code/current specs and label history. |
| License ambiguity | No LICENSE file; documents say deferred. | Add Apache-2.0 text only under owner ruling and record scope. |
| New discovery command freezes an unstable schema | No capabilities contract exists. | Keep it near-term until separately designed/versioned. |
| Authority leakage through release language | Charter and contracts prohibit delivery authority. | Treat preview readiness as evidence, never a gate or approval. |
| Review size exceeds 400 lines | README, license, command compatibility, docs, and release mechanics may cross the budget. | Under `ask-on-risk`, estimate before implementation and request a delivery shape rather than reducing acceptance or inferring an exception. |
| Historical information exposure | Historical artifacts include CNSIC and local paths. | Owner-controlled legal/confidentiality audit; preserve unless an authorized finding requires remediation. |

## Suggested implementation boundary for a later proposal

A proposal may group only these directly coupled changes: Apache license, README/charter status reconciliation, consumer guide, canonical-plus-legacy CLI compatibility/help/version, and narrowly scoped release metadata/documentation for `v0.1.0-preview.1`. It must explicitly measure the review surface first; the session delivery strategy is `ask-on-risk` with a 400 changed-line budget.

It must exclude multilingual census, capabilities schema design, new canonical evidence schemas, API removal, profile policy/CNSIC cutover, repository visibility changes, release publication, history rewriting, and deleting historical artifacts.

## Result Contract

- **status:** success
- **executive_summary:** Current implementation and CI support a Public Preview direction, but the repository is not ready under the owner ruling because license, public presentation, canonical CLI compatibility, discovery/version UX, and owner-visible publication audit work remain unresolved. The evidence-only boundary remains intact and must remain non-authoritative.
- **detailed_report:** This exploration records current evidence, a predecessor-to-ruling gap model, blocker/near-term/post-preview classification, CLI/public-contract findings, workflow audit, historical-artifact handling, and a bounded future-proposal boundary.
- **artifacts:** `openspec/changes/public-preview-readiness/explore.md`
- **paths_created_or_modified:** `openspec/changes/public-preview-readiness/explore.md` (created)
- **next_recommended:** `sdd-propose` only after owner review of the pre-preview blocker set and any delivery-size risk; no later phase was performed.
- **risks:** Canonical command compatibility, stale public documentation, licensing ambiguity, new-contract freezing, authority leakage, review-budget risk, and historical-information exposure.
- **skill_resolution:** `paths-injected` for `/home/kozz36/.agents/skills/cognitive-doc-design/SKILL.md`; `version-management` unavailable; `security-ops` unavailable. No contents were invented for unavailable skills.
