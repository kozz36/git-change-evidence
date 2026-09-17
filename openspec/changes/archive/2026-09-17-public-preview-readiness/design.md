# Design: Public Preview readiness

## Decision and scope

Implement a thin CLI compatibility/discovery boundary and reconcile the public entrypoint with existing implementation. Keep one executable source package, all current evidence constructors and wire contracts, and all delivery authority outside GCE. This design prepares `v0.1.0-preview.1`; it does not publish anything.

Inputs read directly: this change's exploration, proposal and specification; `AGENTS.md`; `openspec/config.yaml`; CLI implementation and co-located tests; `census/go_ast.go`; `go.mod`; README, charter and census documentation; active Inventory V1 and Evidence V1 provenance specs. `.atl/skill-registry.md` was read solely as a publication-audit input, not to select skills.

Execution is auto, artifact storage is OpenSpec, delivery is ask-on-risk, and the review budget is 400 changed lines. Only this design is written in this phase. The supplied package-owned `sdd-design.md` agent definition is the fallback phase instruction (`skill_resolution: fallback-path`); no standalone skill or model selection is invented. CodeGraph/command execution tools are unavailable here: supplied CodeGraph call-flow context was followed by targeted source reads. No Git diff, test run, fresh security scan, or external audit is claimed.

## Architecture and data flow

All new runtime behavior belongs to the existing imperative-shell package `cmd/git-change-evidence`. The public root API, public `census` package, and internal Git/publication adapters remain unchanged. No configuration decoding, command naming, version reporting, or filesystem effects enter the functional core.

```text
one source package -> locally built gce OR git-change-evidence
  main -> runCommand
    explicit discovery recognition -> stdout discovery text -> exit
    census command normalization -> existing census dispatch
      runCensus -> runCensusWith -> parseCensusConfig
        confined source reads -> exact captured bytes + Inventory V1
        -> census.GoASTV1 -> existing census envelope writer
    otherwise -> existing project/publish dispatch or getwd -> positional analysis
```

`runCommand` currently dispatches census before `getwd`; `runCLI` is also directly exercised by tests and independently dispatches census. Both boundaries must use the same pure preparation helpers so their accepted command shapes do not diverge. Normalization is idempotent, so repeating it in the nested path is harmless. Do not move the existing census parser or reads into a new router.

## Runtime decisions

### D1. One source package, two local output names

Document these exact builds from the repository root with Go 1.25.10:

```sh
go build -o gce ./cmd/git-change-evidence
go build -o git-change-evidence ./cmd/git-change-evidence
./gce --version
./git-change-evidence --version
./gce census go-ast --source-root "$(pwd -P)" \
  --receiver os --selector Open -- cmd/git-change-evidence/main.go
```

The last command supplies a physical absolute root and explicit repository-relative source path; zero matches is a valid outcome. All root components still must satisfy census confinement. Use `./gce` in checkout examples; users may place that real executable in their own PATH to invoke `gce`. Do not claim `go install` automatically produces a binary named `gce`: the package basename remains `git-change-evidence`.

Both binaries accept both command spellings; executable basename does not select behavior or alter diagnostics. No new `cmd/gce`, wrapper process, shell alias requirement, symlink installation, Makefile, binary artifact, or distribution channel is necessary. Separate main packages would duplicate entrypoints or force an unnecessary shared CLI package; a shell alias alone would not satisfy reproducible executable discovery.

### D2. Pure, early census normalization

Add `normalizeCensusArgs(args []string) []string` in the existing CLI source. Whenever `len(args) >= 2`, `args[0] == "census"`, and `args[1] == "go-ast"`, produce a newly allocated slice containing `census-go-ast <tail>`, preserving every tail token and byte. Recognize this exact command-family prefix regardless of total arity or whether the tail is valid, including an empty tail and every three-argument vector. Do not mutate the caller's slice/backing array, sort inputs, interpret flags, strip `--`, or modify path/identifier bytes. Leave legacy census arguments unchanged. Keep `parseCensusConfig` accepting its existing normalized grammar only.

Apply discovery first, then normalization before `isCensusShaped` in both `runCommand` and `runCLI`. Existing census-shaped malformed input must still reach the census parser, not positional Git analysis. Valid canonical invocations have at least the three required option/value pairs, separator and path; they normalize directly to the existing accepted legacy vector.

Reserve only the exact two-token prefix `census go-ast` for canonical census, after explicit discovery. Thus `census go-ast --invalid` normalizes to `census-go-ast --invalid`: both must emit empty stdout, `invalid input\n` on stderr, and exit 2 without CWD lookup, Git, stdin, publication or source-root access. `census go-ast` with no tail and `census go-ast <new-file>` likewise reach census validation, never positional analysis. The exact `census go-ast --help` vector remains the explicit discovery path shared with legacy help.

This intentionally resolves the narrow collision formerly interpreted as base reference `census`, source file `go-ast`, and a third positional file in favor of the required canonical command family. Preserving that collision would violate malformed-input parity; do not claim it remains compatible. Preserve unrelated three-argument analysis, including `project <source> <new>`, `publish <source> <new>`, and `census <source-other-than-go-ast> <new>`. Do not reserve all arguments starting with `census`, guess intent from flag-like tails, or fall back to Git after census validation fails. Existing fallback behavior outside the exact prefix and explicit discovery shapes remains unchanged. Document the narrow reservation in command guidance; callers needing the colliding positional analysis can supply the base revision's immutable object ID instead of the literal reference `census`.

No census envelope field or default changes: preserve `git-change-evidence.census-go-ast-cli/v1`, extractor/query identifiers, fixed neutral policy, source scope, canonical inventory digest, order-independent output, raw-path Base64, spans and final LF. Byte parity includes stderr and exit status for equivalent census inputs, not merely decoded JSON equality. Existing downstream write failures can leave an accepted byte prefix; normalization must not claim to retract output or change that limitation.

Rejected: exempting three-argument vectors or recognizing only valid/flag-shaped census tails, because malformed canonical commands would diverge from legacy diagnostics/exits or accidentally invoke Git. Also reject reserving the entire `census` first-token namespace, which would unnecessarily break unrelated positional analysis. Teaching the parser two grammars, duplicating `runCensusWith`, renaming schema IDs to `gce`, or adding a CLI framework creates unnecessary compatibility and review surfaces.

### D3. Exact discovery before all evidence paths

Use a small pure `discoveryText(args []string) (string, bool)` helper and a shared writer path. Recognize only these complete explicit vectors, independently of executable name:

| Arguments | stdout content on success |
| --- | --- |
| `--help` | Root usage: both census forms, positional analysis, project, publish, version |
| `census --help` | Census group and supported `go-ast` command |
| `census go-ast --help` | Canonical census options, explicit scope, legacy compatibility |
| `census-go-ast --help` | Same census guidance, including retained legacy invocation |
| `--version` | The version line specified below |

Successful discovery writes only the requested text with a terminal LF, stderr empty, exit 0. It must precede `getwd`, stdin reads, Git, census root opening and publication. Do not probe the environment, read configuration, inspect files, resolve Git state, or perform network access for discovery. Constants and formatting in the shell suffice.

Do not search arbitrary arguments for `--help` or `--version`. Mixed evidence/discovery vectors, extra tokens, option values, and tokens after `--` retain normal validation/legacy handling; no automatic success or banner is added. No-argument behavior remains existing usage/error behavior. Short flags, `help` commands, subcommand versions, and `capabilities --json` are not added.

On a failed/short discovery stdout write, return exit 3 with existing generic `content unavailable\n` stderr rather than claiming successful discovery. This does not change evidence-path writer behavior. Help output is human text, not a new machine schema. Use stable product/canonical names in the text, never echo arbitrary executable paths.

### D4. Honest source-build version

Define one CLI-local constant `previewVersion = "v0.1.0-preview.1"`. Version output is exactly:

```text
git-change-evidence v0.1.0-preview.1 (source/dev build; not release provenance)
```

This identifies the planned preview while refusing to equate source bytes with an owner-published distribution. Do not add a date, hostname, absolute path, runtime Git lookup, environment override, or invented commit. Do not change the census extractor or wire versions. Documentation repeats the planned identity as a human-facing fact; the constant is the runtime source of truth.

Reject bare release-looking output, unconditional linker-injected release claims, and a new root `VERSION` file read at runtime. A later owner-authorized release/distribution design may add trustworthy build metadata; this slice needs no release mode, tag, upload, workflow, or packaging automation.

## Documentation, licensing and governance

### D5. Standard license, narrow legal claims

Create root `LICENSE` using the complete official Apache License 2.0 text, including its appendix, unmodified. Verify against the official text during implementation; do not synthesize a shortened license or substitute project-specific prose into it. Put the repository's Apache-2.0 identification in README and current governance. Do not invent a copyright owner, third-party clearance, or blanket NOTICE obligation; classify actual dependency/attribution obligations for owner review. Any required attribution addition is evidence-driven and separately reviewed.

### D6. Current authority versus preserved history

README becomes a concise implemented-product entrypoint: evidence-only purpose, planned source/dev preview status, local builds, canonical/legacy census, supported root and census Go APIs, limitations, license, and links to consumer/contributor guidance. Do not imply the CLI currently offers a universal command assembling every public Go contract.

Update `AGENTS.md` and `openspec/config.yaml` to remove false no-source/no-module and deferred-license statements. Preserve Go 1.25.10, architecture, public layout, test commands, thresholds, evidence authority and the exact CNSIC predecessor identity. Keep concrete configuration syntax, broad distribution, signing and release automation deferred. Existing CI is an observed fact, not a deferred implementation decision.

In `docs/product-charter.md`, add a dated/versioned current-preview decision section with a small supersession table: module/source/tests/CI now exist; Apache-2.0 and preview identity are selected; CLI-first positioning retains Go API support. Label the retained bootstrap milestones, old deferred choices and old immediate-next-step section as historical, not current instructions. Correct current-tense contradictions outside the historical section. Do not mechanically convert planned CNSIC profile milestones into completed product capabilities.

Historical OpenSpec records, validation receipts, literal fixtures, local paths and CNSIC provenance remain unchanged. All new prose/comments are English; generated registry descriptions and historical language remain classified exceptions rather than translation targets.

### D7. Agent-consumer guide and precise census examples

Create `docs/agent-consumers.md`, linked from README and census guidance, separate from contributor `AGENTS.md`. Organize it as consume -> validate -> reproduce -> interpret:

- Preserve exact canonical document bytes, terminal LF and digests; do not pretty-print before digest validation. Preserve Base64 raw paths rather than coercing them to UTF-8.
- For revision-bound Report/Result evidence, retain immutable native Git IDs and exact Policy/Inventory/Result antecedents. Standalone report parsing/digest syntax is not provenance proof. Regenerate with the same revisions, policy, scope and supported API; do not invent an absent CLI assembly command.
- For census, capture explicit selected source bytes, hashes, query, limits and extractor identity. It has no Git revision input and is not an atomic whole-tree snapshot. Independent regeneration requires the same supplied bytes and scope, not merely rerunning against a changed checkout.
- Capture exit status and streams separately. Census exits 0/2/3/4 mean technical success/invalid/unavailable/bound; other CLI mappings remain command-specific (5 UTF-8, 6 interrupted, 7 Git race, 8 publication conflict where implemented). Local `publish` stores evidence; it is not product-release publication.
- Nonzero/truncated/unavailable output is not a valid empty result. Zero matches means no syntax candidates in the supplied scope, not semantic absence, safety, security certification or approval. Stale or mismatched inputs require independent reacquisition/regeneration or a consumer-owned decision, never fabricated evidence.

Update `docs/census.md` to lead with the canonical invocation and exact legacy equivalence, linking local builds rather than repeating installation fiction. Retain Linux descriptor confinement, non-Linux unavailability, bounds, supplied-scope semantics, stream/write caveats and schema examples unchanged.

### D8. Publication audit and owner stops

Create `openspec/changes/public-preview-readiness/publication-audit.md` only in a later authorized implementation phase. Each row records surface, inspection method/revision or scope, observation, classification, coverage limitation and owner follow-up. Distinguish inspected facts from inherited exploration findings; do not promote the prior pattern scan to a fresh scan or comprehensive clearance.

| Surface | Default treatment and decision boundary |
| --- | --- |
| README/current docs/config/CLI examples | Inspect accuracy, English and portability; fix only justified current-facing defects |
| `.atl/skill-registry.md` | Generated tooling index containing machine-local paths, metadata and multilingual descriptions; portability/publication relevance requires owner decision; preserve as-is by default |
| Historical OpenSpec/CNSIC predecessor/local paths | Preserved provenance; not a secret finding solely because local/project-specific |
| License/dependencies/historical third-party material | Record attribution and IP questions; Apache-2.0 alone is not clearance |
| Potential sensitive findings | Record a non-secret-bearing summary and owner-controlled review; do not reproduce credentials into the public audit |
| Git history, unpublished refs, external settings/secrets, metadata, outside-tree assets, legal/IP clearance | Explicitly not inspected unless actual authorized evidence is obtained |

No deletion, registry regeneration, `.gitignore` change, history rewrite or removal from tracking is implicit. If the owner later chooses portable regeneration, exclusion, remediation or confidential handling, obtain exact scope/authority first and preserve the audit trail. Human decisions remain mandatory before any visibility change, tag, release, asset publication, destructive cleanup, credential action or legal/confidentiality determination. These are workflow consent boundaries, not automated GCE verdicts.

### D9. Reconcile active OpenSpec without expanding contracts

Use `openspec/config.yaml` for current repository context and keep this change's spec as the additive preview requirement. During later application, inspect active `openspec/specs/**/spec.md` for operationally stale prose, recording changed versus intentionally retained passages in the audit. Do not treat every historical future-tense phrase as a semantic defect.

Concrete likely surface: `openspec/specs/inventory-v1-contract/spec.md` says “before any future Accounting Result V1” and calls adapter work deferred to WU-4C. Update the purpose to describe the existing provenance sequence; qualify the WU-4C reference as historical while preserving the normative acquisition exclusion. Inventory's “Result V1 outside this contract” remains correct package/contract scope and must stay. `openspec/specs/evidence-v1-provenance/spec.md` preserves Report wire/CLI behavior within its original provenance change; that is not a ban on additive shell aliases here, so leave it unchanged. Do not alter its historical validation receipt.

Other active specs require inspection before any targeted prose-only correction; no speculative contract edits are authorized. Do not edit `openspec/changes/archive/**`, bootstrap/go-ast closure artifacts, predecessor hashes or old progress/verification records. Normal later spec synchronization is a separate lifecycle phase, not work performed by this design.

## Exact likely implementation surfaces

| Path | Intended change |
| --- | --- |
| `cmd/git-change-evidence/main.go` | Small pure normalization/discovery helpers, constant/help text, and early calls in both dispatch boundaries |
| `cmd/git-change-evidence/main_test.go` | Discovery stream/effect tests, version literal, positional collision regressions |
| `cmd/git-change-evidence/census_test.go` | Paired canonical/legacy vectors using existing helpers; normalization immutability and error parity |
| `LICENSE` | New unmodified standard license |
| `README.md` | Replace obsolete entrypoint with actual source-build usage and links |
| `AGENTS.md` | Current implementation/license facts, unchanged architecture/testing authority |
| `openspec/config.yaml` | Current preview decisions and remaining deferrals; no threshold relaxation |
| `docs/product-charter.md` | Current decision/supersession section and explicit bootstrap-history boundary |
| `docs/census.md` | Canonical/legacy invocation and discovery/build links |
| `docs/agent-consumers.md` | New consumer contract-use guide |
| `openspec/specs/inventory-v1-contract/spec.md` | Narrow purpose/historical milestone wording only |
| `openspec/changes/public-preview-readiness/publication-audit.md` | Later owner-visible audit findings and coverage gaps |

No expected changes to `census.go`, platform openers, root Go APIs, `internal/`, `go.mod`, dependencies, workflow files or schemas. Reuse current files instead of adding a discovery/version package. If keeping helpers in `main.go` crosses its changed-line threshold, stop for the configured review rather than moving files merely to disguise the diff.

## Validation and strict TDD

For each behavior: RED with focused tests before production changes; GREEN with the smallest shell change; TRIANGULATE using success/failure and boundary cases; REFACTOR only with focused reruns. Planned coverage:

| Test boundary | Required evidence |
| --- | --- |
| Normalizer | Exact-prefix rewrite at every arity (empty tail, `--invalid`, plain third token, complete options), legacy unchanged, idempotence, backing-array nonmutation, raw bytes/`--` preserved; unrelated prefixes unchanged |
| Both dispatch entrypoints | Pair `census go-ast --invalid` with `census-go-ast --invalid`, empty tails, and plain third-token tails: empty stdout, exact `invalid input\n` stderr, exit 2, and no CWD/Git/stdin/publication/source-root effects; canonical malformed input never falls into positional analysis |
| Narrow positional reservation | `census go-ast <new-file>` reaches census validation even with a usable Git fixture; unrelated three-argument `project`, `publish`, and `census` references with source other than `go-ast` retain existing analysis bytes/exits; an immutable base ID retains positional access to a source named `go-ast` |
| Census success | Byte-identical stdout/stderr/exit for matches, zero matches, reversed path order, Unicode identifiers and non-UTF-8 paths; existing schema/default assertions retained |
| Census failures | Paired missing/duplicate flags, separator/path/identifier errors, missing file, parser error, bounded input, unavailable root, and writer failure with exact generic diagnostics |
| Resource/platform seams | Normalize into `runCensusWith` with injected small limits/openRoot to test all bounds without huge fixtures; retain Linux hostile-file tests; verify non-Linux unavailable behavior on a suitable runner without claiming cross-compilation executes it |
| Discovery | Every exact supported shape: only applicable stdout, final LF, empty stderr, exit 0; failing/short writer handling |
| No evidence effects | Panic/counting stdin reader, getwd, Git runner and publisher; no calls. Discovery never reaches census dispatch by control flow; binary smoke from an empty non-repository directory with no source/publication tree |
| Negative discovery | Extra tokens/mixed evidence flags and `--` values do not trigger help; normal empty/invalid/project/publish and unrelated positional output remains unchanged; the exact canonical prefix follows D2 even when malformed |
| Built executable parity | Build both named binaries from the one package into test-owned temporary locations; invoke each for help/version, canonical and legacy census success/failure and compare bytes/status |
| Existing contracts | Full suite remains green; no golden/schema/digest updates justified by a command rename |

Use `go test ./cmd/git-change-evidence -run 'Test.*(Census|Discovery|Version|Normalize|RunCommand|RunCLI)'` with actual new test names matching this focus, then `go test ./...` and the exact check-only formatting command in `openspec/config.yaml`. Existing concurrency/confinement coverage is retained; run `go test -race ./...` for concurrency-bearing changes and retain the existing CI race run. Build both executables explicitly; `go test ./...` alone does not establish their documented output names. Keep temporary binaries outside the tracked checkout.

Documentation verification checks executable examples on Linux, README links, license identity/text, current-versus-historical claims, active-spec semantic invariance, English/portable new prose and audit coverage labels. Do not mark owner-controlled external/legal surfaces inspected merely because automated tests pass. No test execution occurred during this design phase.

## Review surface, ordering and rollout

The existing CLI directory contains seven Go files including tests/platform files (four production, three test), already above the configured `cli_go_files: 6` under an all-Go-files count. This design adds no Go files. Do not reinterpret the threshold as production-only without an explicit ruling. Before implementation, check existing crossing receipts; document a missing/current crossing in later `apply-progress.md` and obtain the required package-boundary review before adding any file beyond it.

Other configured observations remain: public root 12 Go files, internal leaf 15, 160 changed lines per Go file and eight changed Go files per package. Planned runtime/test changes touch three existing CLI Go files and no public/internal files. Measure each actual changed-file total; adding a test file to evade 160 lines would worsen the existing file-count crossing and needs review. These development thresholds never become generic-core policy.

Estimated implementation review surface, not a measured diff: standard license about 200 lines; CLI/tests about 180–300; current docs/governance/audit/spec prose about 200–350. Total about 580–850 additions/deletions, before planning artifacts or any extra attribution. Therefore a single 400-line implementation review is at credible risk. Count additions plus deletions, including license/docs/tests and already-branch-local artifacts; do not omit generated/legal text or hide formatting churn. The parent must pause under ask-on-risk for an explicit delivery decision before implementation. No chain strategy or `size:exception` is selected here.

Dependency order, not a selected PR chain or task list:

1. Resolve delivery-size treatment and package-boundary review needs; confirm authorized current-tree audit scope.
2. Establish normalization/discovery RED tests, implement the shared shell boundary and source/dev version, then prove both executable names and legacy regressions.
3. Update executable documentation against tested behavior; add license/current-governance reconciliation and the consumer guide with evidence-only wording.
4. Complete current-tree publication audit and narrow active-spec prose reconciliation, retaining all historical records; leave owner-controlled unknowns explicit.
5. Collect focused/full/build/format and applicable race evidence, measure actual review surfaces, and return findings to the owner. Do not tag, upload, change visibility or publish.

Legal/governance/audit work can be prepared independently of CLI changes, but public command instructions must not claim implementation before it exists. If the owner chooses a chain, later planning must keep behavior, tests and matching command docs coherent in each review unit. Auto execution does not bypass this decision or any publication/destructive consent.

## Rollback and exclusions

Before external publication, revert defective shell changes through ordinary forward commits while retaining legacy census, evidence schemas and supported APIs. Correct canonical-build documentation if an unshipped alias must be withdrawn; do not leave commands advertised but absent. Once a public compatibility promise exists, removing the canonical path requires a separate owner compatibility decision; prefer a corrective fix. Revert misleading new claims to accurate wording, not to known-false no-implementation claims.

Keep audit/history records even when a finding is corrected. License removal does not revoke previously granted rights; licensing corrections require owner/legal review. No history rewriting, release unpublishing, visibility changes or artifact deletion is a rollback mechanism authorized here.

Excluded: multilingual census, semantic analysis, whole-tree snapshots, capabilities protocol, new evidence schemas/digests, Go API removal, profile/CNSIC cutover, CNSIC policy in the core, MEH integration, packaging/signing/release automation, new platform support, unrelated workflows and any delivery authority. No tasks, apply, verify or archive phase is executed by this artifact.
