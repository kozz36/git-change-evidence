# Publication Audit: Public Preview Readiness

## Scope and method

This bounded, current-tree audit records inspected facts, retained provenance, and
coverage limits for task 8. It is evidence only: it does not authorize, approve,
gate, block, publish, remediate, or decide delivery.

| Surface | Inspection method / scope / revision | Observed fact | Classification | Coverage limitation | Owner follow-up |
|---|---|---|---|---|---|
| Current guidance and examples | Direct full-text reads of `README.md`, `AGENTS.md`, `docs/census.md`, `docs/consumer-integration-guide.md`, `docs/product-charter.md`, and both `docs/validation/*.md` files; `examples/` directory check at current worktree revision `a189e2e179c53022b84dd8251638ab6f0513b2ae` | Current guidance states source/dev preview status, CLI/machine-contract primacy, supported Go APIs, and evidence-only boundaries; no `examples/` directory is present. | Inspected current-facing documentation; examples absent. | This is a documentation inspection, not a public-site, package-registry, or external-rendering inspection. | Owner decides any later public presentation or visibility change after separately reviewing the intended external surface. |
| Current repository and CI context | Direct reads of `go.mod`, `.github/workflows/go-pr-checks.yml`, `.github/workflows/pr-policy.yml`, and `openspec/config.yaml` at the current worktree revision | The module identity is present; Go source/tests and two GitHub Actions workflow files are present; the configured context still has stale deferred CI/license/publication wording. | Inspected current-facing repository context; narrow reconciliation justified. | Workflow presence does not establish execution history, external settings, permissions, or publication readiness. | Owner retains decisions about CI administration, external configuration, and any delivery action. |
| Active OpenSpec specifications | Direct full-text reads of exactly `openspec/specs/evidence-v1-provenance/spec.md`, `openspec/specs/inventory-v1-content-identity/spec.md`, `openspec/specs/inventory-v1-contract/spec.md`, `openspec/specs/result-v1-accounting-evidence/spec.md`, `openspec/specs/result-v1-contract/spec.md`, and `openspec/specs/result-v1-provenance/spec.md` | Only Inventory V1's provenance-sequence and WU-4C wording is stale; its acquisition exclusion remains normative. Evidence V1 provenance has no concrete current-facing contradiction and is retained byte-unchanged. The remaining four active specs are retained unchanged. | Inspected active contracts; one narrow prose correction, otherwise retained. | This does not revalidate implementation or change schemas, APIs, scenarios, or contract semantics. | Owner may authorize a separate contract change only when a concrete semantic contradiction is established. |
| `.atl/skill-registry.md` | Direct full-text read at the current worktree revision | The file identifies itself as generated metadata and contains machine-local skill paths and descriptions. | Generated local-path metadata and audit input; not deletion authority or product policy. | Local-path portability and any future tracking treatment were not decided here. | Owner must choose any regeneration, exclusion, removal, or confidentiality handling with explicit scope; this audit authorizes none. |
| Historical OpenSpec, CNSIC, and local-path provenance | Current guidance and validation-document reads; predecessor object resolution attempted without history traversal | Historical records and local paths are evidence, not current instructions. The referenced CNSIC predecessor object was not locally resolved in this worktree, and no historical artifact was edited. | Preserved provenance; historical/current distinction required. | Git history was not traversed; archive bytes, predecessor contents, and unpublished refs were not inspected. | Owner decides any archival correction, predecessor investigation, or migration under separate authority; preserve the audit trail. |
| License, dependency, and third-party questions | Direct read of `LICENSE` and `go.mod` at the current worktree revision | Apache License 2.0 text is present; the module declares `golang.org/x/sys v0.38.0`. | Inspected repository declarations; legal and third-party questions remain open. | This is not a license-compliance, attribution, NOTICE, intellectual-property, export, or dependency-security review. | Owner obtains appropriate legal, attribution, dependency, and security review before any distribution or other external use decision. |
| Potential sensitive findings | Scoped prose and generated-metadata inspection; no secret value is copied into this record | No secret value is reproduced. The inspected material includes local paths, which alone are not a secret finding. | Non-secret-bearing audit record; potentially sensitive material requires restricted handling. | No comprehensive credential, history, external-secret, or metadata audit was performed. | Owner chooses any restricted review, credential action, or destructive remediation after explicit authorization; GCE supplies no such authority. |
| Git history and unpublished refs | Not inspected | Git history was not traversed and unpublished refs were not inspected. | Uninspected. | Current-tree reads cannot establish historical or unpublished-ref content. | Owner authorizes any scoped history/ref review and determines any response to its findings. |
| External settings, secrets, and metadata | Not inspected | Hosting settings, CI secrets, organization metadata, access controls, and external service state were not inspected. | Uninspected. | No external-system access was used. | Owner controls access, review, and remediation decisions for external systems. |
| Outside-tree generated and release assets | Not inspected | Generated outputs and release assets outside the tracked worktree were not inspected. | Uninspected. | Repository-tree inspection cannot establish outside-tree asset content or distribution state. | Owner decides whether to inventory, retain, remove, or publish such assets under explicit authority. |
| Legal and intellectual-property clearance | Not inspected | Apache-2.0 selection is a repository fact, not legal, confidentiality, export, attribution, or intellectual-property clearance. | Uninspected owner/legal determination. | This task has no legal-review authority or evidence sufficient for clearance. | Owner obtains the required legal or confidentiality determination before external exposure or remediation. |
| Visibility, release, and destructive-action boundaries | Current task, proposal, design, and guidance review | Future visibility changes, tags, releases, assets, destructive cleanup, credential actions, and legal/confidentiality determinations remain human-controlled decisions. | Owner-only consent and decision boundaries. | These boundaries are informational; they are not GCE approvals, rejections, gates, or automated blockers. | Owner gives explicit consent and selects the applicable process before any such action. |

## Retained active-spec results

- `openspec/specs/evidence-v1-provenance/spec.md` was fully read and retained
  byte-unchanged because no concrete current-facing contradiction was found.
- `openspec/specs/inventory-v1-content-identity/spec.md`,
  `openspec/specs/result-v1-accounting-evidence/spec.md`,
  `openspec/specs/result-v1-contract/spec.md`, and
  `openspec/specs/result-v1-provenance/spec.md` were fully read and retained
  unchanged because task 8 found no justified prose correction.
- `openspec/specs/inventory-v1-contract/spec.md` retains every normative
  requirement, scenario, and acquisition exclusion; only its historical
  provenance and WU-4C reference are reconciled.

## Preservation and correction policy

No archive, predecessor hash, literal historical evidence, or generated registry
was changed. If this audit later requires correction, append a dated correction,
supersession, or withdrawal record here instead of removing this audit history.
