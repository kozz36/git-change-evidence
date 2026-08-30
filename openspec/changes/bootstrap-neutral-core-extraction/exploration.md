## Exploration: bootstrap-neutral-core-extraction

### Current State

The product charter defines an evidence-only `git-change-evidence` product using Functional Core / Imperative Shell and Hexagonal Architecture. The repository is pre-bootstrap and contains no implementation; packaging, source layout, license, runtime matrix, CI, release automation, and configuration-file syntax are explicitly open decisions.

CNSIC W1/W2 is the immutable extraction baseline, not a new implementation target. The baseline is preserved in the CNSIC `origin/main` Git objects at merge commit `a0fc7b26ff8a0e0a61baa586b32c46841611c806` and contains:

- `scripts/sd7/models.py` — strict, frozen Pydantic v1 contracts: native 40/64-character Git `Sha`, 64-character `Digest`, snapshot entry unions, accounting totals, thresholds, inventory, forecast comparison, annotations, and evidence-only report documents.
- `scripts/sd7/contracts.py` — runtime/schema validation, forbidden-authority-key checks, canonical JSON bytes, and schema snapshot verification.
- `scripts/sd7/git_snapshot.py` — imperative Git acquisition with argv-only execution, closed Git environment, immutable ref resolution, raw path preservation, object-kind detection, isolated numstat, and native SHA-1/SHA-256 identity.
- `scripts/sd7/inventory.py` — separate untracked inventory with descriptor-relative no-follow traversal and link-payload hashing.
- `scripts/sd7/policy.py` — data-only policy loading; CNSIC globs and precedence remain profile data.
- `scripts/sd7/accounting.py` — pure exclusive accounting: `artifacts → tests → mechanical → production`, validation before precedence, text-only totals, and non-gating threshold observation.
- `tests/fixtures/sd7/compatibility/v1/` — the two frozen portable vectors (`contracts-v1.json`, `accounting-v1.json`) and exact SHA-256 membership manifest. These are compatibility evidence, not a generic-core filename/count/schema contract.

The archived W1/W2 artifacts record 152 focused tests, 4/4 requirements, 13/13 scenarios, approved Judgment Day after the descriptor-relative race correction, and no remaining severe findings. Historical W3/W4/W5 plans are therefore transfer provenance. They must be re-derived against the charter and the frozen vectors rather than copied as active tasks or forecasts.

### Current-State / Module Ownership Map

| Boundary | Owns | Must not own | Evidence / successor responsibility |
|---|---|---|---|
| Functional Core | Normalization, immutable-input domain models, validation, exclusive accounting, carveout matching, canonical serialization, report assembly, projections, forecast comparison, and deterministic manifest calculations | Git/process execution, filesystem effects, project policy, publication effects, or delivery authority | W1 contracts and W2 accounting are the semantic starting point; successor interfaces must consume typed immutable values/bytes rather than recreate DTOs |
| Imperative Shell | Git argv calls, ref/object acquisition, environment isolation, filesystem reads, untracked discovery, CLI parsing/adaptation, retry orchestration, output streams, and local atomic publication | Policy classification, accounting rules, compatibility semantics, or approval/gating decisions | `git_snapshot.py` and `inventory.py` are behavioral evidence for the shell boundary, not a mandated future source layout |
| CNSIC profile adapter | CNSIC globs, categories, thresholds, forecast vocabulary, OpenSpec/CHANGELOG integration, review interpretation, and profile-specific compatibility behavior | Generic Git semantics, neutral contract meaning, or delivery authority | Current `policy-v1.json` and frozen accounting vector define profile evidence; future integration belongs here |
| Compatibility edge | Legacy `scripts/sd7_move_carveout.py` positional interface and one-line output, W1/W2 vector execution, translation between legacy/CNSIC shapes and neutral contracts | A second accounting/matching algorithm, unsafe working-tree semantics, or neutral-core policy | Preserve legacy behavior at the edge while delegating to a bounded immutable-bytes matcher; the frozen vectors pin only W1/W2 boundaries |
| Eventual CNSIC consumer cutover | CNSIC callers, review/document integration, migration from legacy commands, and independent interpretation of evidence | Generic-core ownership, package/release decisions, or approval/gate authority | Cut over only after the profile consumes equivalent canonical evidence and compatibility vectors; no CNSIC code is modified by this exploration |

### Affected Areas

- `docs/product-charter.md` — authoritative product boundary, architecture, milestones, and open decisions.
- `openspec/config.yaml` — confirms documentation-only bootstrap state and unresolved implementation/testing decisions.
- CNSIC `scripts/sd7/{models,contracts,git_snapshot,inventory,policy,accounting}.py` at the immutable baseline — source behavior to preserve or translate, not source files to modify here.
- CNSIC `tests/fixtures/sd7/compatibility/v1/*` at the immutable baseline — W1/W2 compatibility vectors and manifest boundary.
- Archived CNSIC W1/W2 proposal, design, tasks, specs, review, Judgment Day, and verification artifacts — historical evidence for what shipped, what was corrected, and which W3/W4/W5 assumptions are stale.
- CNSIC `scripts/sd7_move_carveout.py` — legacy compatibility surface requiring an adapter strategy in a later implementation phase.
- Future standalone core/profile/CLI modules — ownership to be decided by bootstrap design; no source layout is selected here.

### Approaches

1. **Mechanical copy / big-bang move** — copy the CNSIC package and tests into the standalone repository, then remove or rewrite CNSIC-specific behavior in place.
   - Pros: lowest initial planning overhead; preserves many existing imports and test shapes.
   - Cons: imports CNSIC vocabulary and governance into the neutral core; makes compatibility and semantic cleanup simultaneous; creates a large rollback boundary; encourages copying stale W3/W4/W5 assumptions; weakly separates shell effects from pure logic.
   - Effort: High risk despite apparently low initial effort.

2. **Package-first extraction** — establish a standalone neutral package boundary first, then move W1/W2 semantics and add W3/W4/W5 capabilities behind that boundary before CNSIC integration.
   - Pros: clear ownership and dependency direction; preserves Functional Core / Imperative Shell; allows real-Git and compatibility proof at each boundary; keeps the CNSIC profile outside neutral contracts.
   - Cons: requires an explicit bootstrap design; temporary adapters and duplicate boundary representations may exist; eventual CNSIC cutover is later.
   - Effort: Medium.

3. **Branch by Abstraction with compatibility vectors and later CNSIC cutover** — retain the CNSIC implementation as the observed predecessor, introduce a neutral abstraction and profile adapter, execute the frozen W1/W2 vectors against the successor, then add W3/W4/W5 serially and switch CNSIC consumers only after parity is demonstrated.
   - Pros: smallest reversible semantic step; compatibility vectors provide an external parity oracle; isolates profile policy at the edge; supports staged consumer cutover and rollback; makes divergence visible before CNSIC adoption.
   - Cons: temporarily maintains predecessor and successor paths; requires disciplined vector/version ownership; cutover and deprecation work remain future tasks.
   - Effort: Medium, with the lowest migration risk.

### Recommendation

Use **Branch by Abstraction**, implemented as a package-first neutral boundary with the CNSIC profile and compatibility edge outside the Functional Core. Do not mechanically copy the CNSIC package or treat its filenames as the standalone layout.

The smallest safe route is:

1. Freeze the W1/W2 vector set and its manifest as the extraction gate; this is already the charter's Milestone 0 prerequisite.
2. In bootstrap design, decide only the neutral module boundaries, profile boundary, contract ownership, test strategy, and compatibility-vector execution model. Leave packaging, source layout, license, runtime matrix, CI, release automation, and config syntax open.
3. Establish the neutral W1/W2 semantic boundary from the immutable baseline: contracts/canonical bytes, Git snapshot inputs, inventory separation, policy injection, and accounting.
4. Add W3 as a pure bounded carveout core plus a legacy adapter. The adapter preserves positional arguments and one-line output but cannot read the working tree as the source of immutable evidence.
5. Add W4 as report assembly, projections, CLI adaptation, publication, exit mapping, and moving-head retry. Keep all Git and filesystem effects in the shell.
6. Add W5 as forecast ingestion/comparison and pilot evidence. It may observe thresholds and ratios but cannot author forecasts, approve changes, or become a gate.
7. Run a later CNSIC profile/cutover work unit: prove frozen-vector parity, switch consumers, retain rollback compatibility, and remove legacy use only under an explicit compatibility decision.

### Historical W3/W4/W5 to Product Milestone Mapping

| Historical work | New product milestone | Re-derived ownership |
|---|---|---|
| W3 — carveout and legacy compatibility | **Milestone 2: Neutral core extraction**, then **Milestone 4: First consumer profile** | Pure bounded matcher and carveout result belong to the Functional Core; Git/blob acquisition belongs to the shell; legacy command translation belongs to the compatibility edge; CNSIC-specific fallback interpretation belongs to the profile |
| W4 — report CLI, projections, publication, exit semantics, moving-head safety | **Milestone 3: CLI and publication** | Report construction/projections are core; Git invocation, retry, CLI adaptation, and atomic local publication are shell responsibilities; statuses remain evidence-only and technical exits are not delivery verdicts |
| W5 — forecast, thresholds, ratios, pilot/governance evidence | **Milestone 4: First consumer profile**, with later evaluation under **Milestone 5: Distribution decision** only if distribution evidence is needed | Forecast comparison is neutral evidence machinery; forecast vocabulary, CNSIC thresholds, OpenSpec/CHANGELOG integration, and review interpretation belong to the CNSIC profile; governance authority remains outside the product |

The historical W2 work is the extraction input, not a new milestone: it satisfies the W1/W2 freeze gate and supplies the compatibility boundary for Milestone 2. Historical W3/W4/W5 forecasts are stale; `sdd-tasks` must measure new production/test/mechanical/artifact addends from its own item table and current bytes.

### Risks

- **Semantic contamination:** copying `policy-v1.json`, CNSIC thresholds, forecast vocabulary, or review terminology into the neutral core would violate the charter. Mitigation: inject policy and keep CNSIC translation in the profile adapter.
- **Compatibility drift:** replacing the legacy carveout behavior with a new algorithm can silently change output. Mitigation: preserve the legacy edge and add divergent compatibility vectors/tests before cutover.
- **Unsafe evidence source:** using index or working-tree state would destroy reproducibility. Mitigation: shell acquisition must resolve immutable objects and keep untracked inventory separate.
- **Path and process hazards:** hostile raw paths, refs, environment variables, and symlink races can alter evidence or escape intended roots. Mitigation: retain argv-only execution, closed Git environment, raw-byte handling, and descriptor-relative no-follow acquisition.
- **Moving-reference publication race:** a symbolic head can move between measurement and publication. Mitigation: bounded re-resolution and retry, then a technical repository/race failure without publishing a mismatched identity.
- **God-module growth:** combining Git, policy, accounting, report formatting, and publication would cross ownership boundaries. Mitigation: enforce core/shell/profile/compatibility interfaces and deliver W3, W4, and W5 as separate work units.
- **False authority:** threshold attention, pilot observations, or CLI statuses could be interpreted as approval or blocking. Mitigation: keep evidence-only field vocabulary and explicitly test absence of delivery-authority semantics.
- **Premature packaging decisions:** selecting package format, source layout, license, runtime matrix, CI, release automation, or config syntax now would exceed the charter. Mitigation: record these as open decisions for later design or Milestone 5.

### Ready for Proposal

**Yes.** The proposal should state that this is a standalone bootstrap, that CNSIC W1/W2 is the immutable extraction baseline, that Branch by Abstraction is preferred over a big-bang copy, and that W3/W4/W5 are re-derived product milestones rather than copied CNSIC tasks. It should explicitly preserve the open decisions listed in the charter and defer consumer cutover until compatibility parity is demonstrated.

### Result Contract

- **status:** success
- **executive_summary:** The immutable CNSIC W1/W2 baseline supplies frozen contract/accounting compatibility evidence, while the standalone product remains documentation-only. Branch by Abstraction with a package-first neutral boundary, an edge compatibility adapter, and a later CNSIC profile cutover is the smallest safe route; historical W3/W4/W5 map to neutral extraction, CLI/publication, and profile milestones.
- **detailed_report:** Current-state/module ownership map, three migration approaches, recommendation, milestone mapping, risks, and proposal readiness are recorded above. Packaging, source layout, license, runtime matrix, CI, release automation, and configuration syntax remain explicitly undecided.
- **artifacts:** `openspec/changes/bootstrap-neutral-core-extraction/exploration.md`
- **next_recommended:** `sdd-propose`
- **risks:** Compatibility drift, policy contamination, unsafe mutable-state acquisition, path/process races, moving-reference publication races, god-module growth, false delivery authority, and premature resolution of charter-open decisions.
- **skill_resolution:** `paths-injected` — loaded `/home/kozz36/.agents/skills/sdd-explore/SKILL.md` and `/home/kozz36/.config/opencode/skills/cognitive-doc-design/SKILL.md`; no delegation, implementation, commit, push, PR, or archive performed.

## Key Learnings

1. CNSIC W1/W2 is a frozen compatibility baseline rather than a standalone product implementation.
2. CNSIC policy and review vocabulary must remain profile-owned outside neutral core contracts.
3. Branch by Abstraction provides the safest reversible path for later consumer cutover.
4. Historical W3/W4/W5 plans require remeasurement instead of direct task-text reuse.
