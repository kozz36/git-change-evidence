# Proposed W1/W2 CLI requirements

**DRAFT · PROPOSED · NOT IMPLEMENTED.** These are future requirements, not current behavior or implementation authorization. [Design](../../design.md) supplies proposed grammar, transport, diagnostics and the selected initial Linux resource/threat profile; [tasks](../../tasks.md) supplies conditional work. Existing canonical contracts remain unchanged.

## ADDED Requirements

### Requirement: R1 — Explicit selected-path acquisition

The proposed inventory command SHALL accept the closed UTF-8 `git-change-evidence.cli-path-selection/v1` request described in the design, acquire caller-selected raw paths beneath an independently authorized absolute root, and construct existing Inventory V1. It SHALL NOT discover additional paths or treat selection as authorization/confidentiality proof. Initial acquisition SHALL use caller-controlled local roots and a stable mount namespace as described in the selected design profile; name-based Git-metadata rejection includes ASCII case variants of `.git`, not certification of all physical aliases/hardlink provenance. It SHALL hash/measure captured bytes, without claiming an atomic snapshot or detection of every same-inode content mutation.

#### Scenario: Raw paths and safe content

GIVEN an authorized supported root, valid policy and unique Base64-encoded relative paths including spaces, newlines, leading hyphens or non-UTF-8 bytes
WHEN inventory acquires regular files or leaf symlink payloads, including dangling links
THEN raw path bytes and acquired content bytes reach `NewInventoryV1` unchanged, symlink targets are not read, and no unselected path is inventoried.

#### Scenario: Empty requested selection

GIVEN a validated accessible authorized root and `paths_b64: []`
WHEN inventory completes
THEN it returns valid empty Inventory V1 bound to policy, not a claim of repository emptiness or complete untracked coverage.

#### Scenario: Unsafe or malformed selection

GIVEN duplicate/unknown JSON keys, duplicate decoded paths, noncanonical Base64, invalid UTF-8 transport, trailing JSON, traversal, absolute/NUL/empty paths, `.git` components or their ASCII case variants, intermediate symlinks or nonregular leaves
WHEN inventory validates or acquires the selection
THEN it fails without following unsafe targets, enumerating directories, normalizing Unicode, expanding globs/tilde, or returning synthetic empty success.

### Requirement: R2 — Exact versioned transport and output

Each stage SHALL decode canonical antecedents using the actual versioned APIs and emit exactly the existing canonical document including LF. File/stdin selectors SHALL be explicit and separately authorized; at most one SHALL use stdin. The shell SHALL buffer and validate bounded output before normal emission.

#### Scenario: Exact success streams

GIVEN valid stage inputs
WHEN a stage completes successfully
THEN exit is `0`, stderr is empty, and stdout is exactly its Inventory V1, Accounting Result V1 or Report V1 canonical bytes, with no extra digest, newline or wrapper; consumers can retain SHA-256 of those exact bytes.

#### Scenario: Ambiguous stdin or noncanonical input

GIVEN two stdin selectors or a noncanonical versioned document
WHEN the command validates its inputs
THEN exit is `2`, stdout is empty and the stable corresponding diagnostic is emitted; unused stdin is never read.

#### Scenario: Output failure

GIVEN valid buffered evidence and a short/failed stdout writer
WHEN emission fails
THEN exit is `9` with `output_incomplete`, already-written bytes may remain, and the command never claims success or retractability.

### Requirement: R3 — Immutable full range and retained antecedents

Result SHALL resolve explicit base/head refs to immutable commits and acquire the entire committed range without inventory-path filtering. Repository/ref arguments SHALL remain data passed to a bounded sanitized Git adapter, not shell commands; missing objects SHALL NOT trigger implicit fetching or network access. Report SHALL validate retained exact policy/inventory/result binding and derive existing Report V1 provenance using the versioned binding APIs.

#### Scenario: Change outside inventory

GIVEN inventory selects only `a.txt` and the committed range changes both `a.txt` and `b.txt`
WHEN result is constructed
THEN both committed changes remain in scope; worktree edits and inventory selection do not narrow or substitute for the range.

#### Scenario: Literal refs and unavailable Git evidence

GIVEN valid refs containing shell metacharacters, an invalid/foreign root, or an unresolved otherwise-valid ref/missing object
WHEN result attempts acquisition
THEN refs remain literal data without shell execution or implicit network fetch; invalid root/request fails with exit `2`/`invalid_arguments`, missing evidence with exit `3`/`git_unavailable`, and pre-emission stdout remains empty on failure.

#### Scenario: Binding mismatch

GIVEN a result and a different retained policy or inventory, even if each is individually valid
WHEN report decodes and binds them
THEN it fails with exit `2` and `antecedent_mismatch`, stdout empty; an isolated digest is not accepted as retained provenance.

### Requirement: R4 — Neutral authority and scope

The family SHALL report technical evidence only. It SHALL support only neutral Policy V1 semantics, without a project configuration format, profile/plugin loader, complete-inventory claim or delivery authority.

#### Scenario: Observations are data

GIVEN valid zero-change evidence or threshold-exceeded/unavailable observations
WHEN result/report succeeds
THEN exit remains `0`; observations do not approve, reject, gate, block, merge or release anything, or decide delivery.

#### Scenario: No implicit discovery or effects

GIVEN unrelated untracked files and a selected-path request
WHEN any stage runs
THEN it does not consult Git status/ignore for inventory completeness, modify sources/repository/index, publish output, create temporary evidence files or load consumer profiles.

### Requirement: R5 — Explicit versions, failures and finite bounds

The family SHALL implement the design's command-specific stable technical errors, distinct from existing commands. It SHALL enforce the owner-selected [initial Linux resource profile](../../design.md#initial-linux-resource-profile), including its numeric input/document/path/content/aggregate/Git/output budgets and cooperative deadline. These are new-family admission limits, not limits of existing Go APIs/commands, a hard RSS/wall-time certificate or calibrated workload coverage. Unsupported request/core versions SHALL fail explicitly without fallback; source/dev identity SHALL NOT imply release provenance.

#### Scenario: Unsafe, unavailable or bounded acquisition

GIVEN an unsafe path, unavailable input, unsupported platform, crossed resource budget or cancelled operation
WHEN acquisition cannot complete
THEN it fails respectively with the design's exit `2`, `3`, `3`, `4` or `6` and corresponding token; pre-emission stdout is empty, never truncated/empty-success evidence.

#### Scenario: Version separation

GIVEN an unsupported selection-request or canonical-document version
WHEN decoding occurs
THEN exit is `2` with `unsupported_version`, no guessing or legacy decoder fallback occurs, and no preview identity changes the schema interpretation.

#### Scenario: Detectable acquisition race

GIVEN a detectable Git-object or selected-source race under the approved acquisition model
WHEN the operation cannot retain its declared evidence binding
THEN exit is `7` with `acquisition_race` and pre-emission stdout empty; this does not claim detection of every concurrent content change or an atomic whole-tree snapshot.

#### Scenario: Initial Linux capability boundary

GIVEN a non-Linux operation environment or unavailable required descriptor/no-follow/leaf-link capabilities
WHEN a new-family operation is admitted
THEN exit is `3`/`unsupported_platform` before acquisition, no weaker pathname fallback occurs, and existing Go APIs/census/legacy commands are not given a new platform restriction.

#### Scenario: Individual and aggregate budget boundaries

GIVEN each selected numeric size/count budget exercised below, exactly at and above its limit, or individually admissible documents whose combined input exceeds 32 MiB
WHEN the operation validates or acquires otherwise-valid inputs
THEN within-budget cases remain eligible subject to all other rules; any size/count overrun returns `4`/`resource_limit` with pre-emission stdout empty, without truncation, deduplication that bypasses aggregate charging or synthetic empty success.

#### Scenario: Cooperative deadline and external termination

GIVEN the 30 s cooperative deadline expires before emission, a write fails during emission, or an external supervisor forcibly terminates the process
WHEN interruption is handled or the caller collects the actual outcome
THEN detected pre-emission cancellation uses `6`/`interrupted`, failed emission uses `9`/`output_incomplete`, and external termination has no promised family exit/diagnostic/complete stdout; neither a hard RSS limit nor hard completion time under blocked OS/stdio is certified.

### Requirement: R6 — Existing command compatibility

Existing discovery/census priority, legacy triples, streams and errors SHALL remain unchanged except one declared additive root `--help` entry advertising the new family; existing help entries remain verbatim. Root `--version` and existing command-specific help SHALL remain unchanged. The new family SHALL admit only non-triple operation forms and two-argument family help, without extending legacy project/publish document acceptance.

#### Scenario: Legacy triple collisions

GIVEN triples beginning with `project`, `publish` or hypothetical valid ref `w1w2`, including `w1w2 inventory --help`
WHEN dispatch occurs
THEN they retain their prior legacy measurement route rather than new stage/help routing; existing census discovery/normalization exceptions remain unchanged.

#### Scenario: Differential preservation

GIVEN existing root/version/project discovery, canonical/compatibility census, legacy measurement, projection and publication cases
WHEN compared before and after future implementation
THEN dispatch, streams, exits and accepted input kinds are unchanged except the single declared additive root-help entry; no other help/output drift is allowed and a different versioned document is not silently accepted by project/publish.
