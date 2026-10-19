# Proposed W1/W2 shell design

> Historical design record retained below. For the local source implementation and unfinished qualification, see [current progress](apply-progress.md). This is not a release or adoption claim.

**DRAFT · PROPOSED · NOT IMPLEMENTED.** All command names, request identifiers, diagnostics and new shell behavior below are proposals. Existing public Go APIs are source facts, not a new publication claim. See [proposal](proposal.md), [requirements](specs/w1-w2-cli/spec.md) and [conditional tasks](tasks.md).

## Grammar and dispatch

```text
gce w1w2 --help
gce w1w2 inventory --source-root <absolute-directory> --policy <document> --paths <selection>
gce w1w2 result --repo <absolute-repository> --base <revision> --head <revision> --policy <document> --inventory <document>
gce w1w2 report --policy <document> --inventory <document> --result <document>
```

Flags are required exactly once; order is immaterial. Unknown flags, positional extras, missing values and duplicate flags are invalid. Selectors are explicit authorized files or `-` for stdin, including `--paths`; at most one selector may use stdin. Never read unused stdin. No output file, implicit current-directory root, implicit policy, or implicit HEAD is supplied. Report uses the proposed fixed subject `committed change evidence`, satisfying the existing required subject without adding another input.

Retain existing discovery first, then census normalization/routing, then every existing legacy triple, before admitting the new non-triple family. Family help has exactly two arguments. Do not introduce `w1w2 <stage> --help`: such a triple must retain legacy routing. This includes hypothetical valid refs named `w1w2`, `project` or `publish`. Census discovery exceptions remain. Current ordering and project/publish triple exclusions: `cmd/git-change-evidence/main.go:65–85,112–164`; legacy three-argument measurement: `main.go:247–260` in that same directory. New names must not extend existing project/publish accepted document kinds. The one proposed help-text addition is a root `--help` entry advertising the new family; preserve existing entries verbatim. Root `--version` and existing command-specific help remain unchanged. Differential checks must allow only that declared additive help entry, not arbitrary help/output drift.

## Closed selection request

Proposed transport identifier, not a core evidence/configuration schema:

```json
{"schema":"git-change-evidence.cli-path-selection/v1","paths_b64":["YS50eHQ="]}
```

Require one UTF-8 JSON object, exactly `schema` and `paths_b64`, no duplicate/unknown keys, no trailing JSON values; allow JSON whitespace. `paths_b64` is an array of canonical standard Base64 strings (decode/re-encode equality), conveying raw relative path bytes. Reject duplicate decoded paths, empty paths, absolute paths, NUL, empty components, `.`/`..` components and any `.git` component, including ASCII case variants. No glob, tilde, Unicode normalization or implicit recursion. Spaces, newlines, leading hyphens and non-UTF-8 path bytes survive unchanged; transport text itself remains UTF-8. Apply the same name-based Git-metadata exclusion to root components. Do not reinterpret raw names to bypass these rules. This exclusion is not certification of physical hardlink provenance or the absence of every possible metadata alias; paths are not authorization or confidentiality proof.

`[]` is an EMPTY REQUESTED selection, successfully verified only after root authorization/access and input validation succeed. It says nothing about the directory's other contents. Canonical document selectors are separately authorized explicit inputs, not directory-discovery requests.

## Safe acquisition proposal

The harness authorizes the source root and selected scope separately; GCE validates and reads within that scope. Open the authorized root as a directory descriptor and resolve components descriptor-relatively without following symlinks. Reject unsafe intermediate symlinks, metadata paths, traversal and nonregular entries; permit only regular-file leaves or leaf symlinks. Read leaf symlink payload bytes, including dangling links, without dereferencing the target. Directory leaves do not trigger recursion. Bound link payloads as content.

Selected initial DESIGN profile: execute this new family only on Linux, using local roots controlled by the invoker and a stable mount namespace. The required descriptor-relative no-follow and handle-based leaf-link payload operations must be enforceable; other platforms or unavailable required capabilities return exit 3/`unsupported_platform` before acquisition, without a weaker pathname-based fallback. This does not restrict the public Go API, census or legacy commands, and does not claim that any runtime/platform qualification has run.

Confinement is relative to the captured directory descriptors and selected names, not an inode-provenance or confidentiality certificate. Privileged actors changing mount/namespace topology, and adversaries controlling the invoking process, are outside this initial model. Do not claim exclusion of every physical metadata alias or protection against every hardlink relationship. Detectable acquisition races fail with exit 7/`acquisition_race`; hash and measure the captured byte stream. Retained descriptors, before/after metadata checks and captured bytes do not establish an atomic whole-tree snapshot or detection of every in-place write on the same inode.

Existing census acquisition is a useful pattern, not a conforming inventory reader: `cmd/git-change-evidence/census_open_linux.go:93–123` rejects leaf symlinks and checks device/inode/type identity, not stable contents. The proposed leaf-payload mechanism and all asserted protection need future implementation and tests; no existing support is inferred.

Hash and byte length identify captured bytes within the requested scope. Unsafe, missing, unreadable, interrupted or bounded acquisition is a technical failure, never synthetic empty success. No source/repository writes, temporary output files, Git index refresh, or publication are part of these operations.

## Existing neutral API mapping

Source citations below are repository-relative; implementations remain unchanged by this proposal.

| Stage | Existing public API / contract |
|---|---|
| Policy input | `DecodeAccountingPolicyV1(raw []byte) (PolicyDocumentV1, error)`; `policy_decode.go:11–41`. Schema `git-change-evidence.accounting-policy/v1`; `policy_contract.go:16`. |
| Inventory construction | `NewInventoryV1(policy PolicyDocumentV1, entries []InventoryEntryInput) (InventoryDocumentV1, error)`; `inventory_v1.go:3–26`. Schema `git-change-evidence.inventory/v1`. |
| Inventory input | `DecodeInventoryV1(raw []byte, policy PolicyDocumentV1) (InventoryDocumentV1, error)`; `inventory_v1_decode.go:11–52`. |
| Result construction | `NewAccountingResultV1(policy PolicyDocumentV1, inventory InventoryDocumentV1, snapshot CommittedSnapshot) (ResultDocumentV1, error)`; `result_v1_build.go:8–42`. Schema `git-change-evidence.accounting-result/v1`; `result_v1.go:3`. |
| Result input | `DecodeAccountingResultV1(raw []byte, policy PolicyDocumentV1, inventory InventoryDocumentV1) (ResultDocumentV1, error)`; `result_v1_decode.go:5–26`. |
| Report construction | `NewReportV1FromResult(subject string, binding ResultV1ReportBinding) (Evidence, error)`; binding fields in `contract.go:36–40`, constructor at `contract.go:72`. Report schema is `git-change-evidence.evidence/v1`, kind `report`, not a new report schema identifier; `contract.go:14–17`. |
| Report binding | `ValidateReportV1ResultProvenance(report Evidence, binding ResultV1ReportBinding) error`; `contract.go:80`. |

The inventory core hashes supplied contents, sorts raw paths and rejects duplicates; it performs no filesystem discovery (`inventory_v1_build.go:23–52`). Canonical Inventory V1 includes final LF and exact-byte SHA-256 (`inventory_v1_build.go:91–105`); Policy V1 likewise (`policy_contract.go:117–123`), Report V1 likewise (`contract.go:56–67`). Preserve each versioned document's existing canonical bytes, LF and digest semantics, without reserialization or an extra newline.

Result shell acquisition must resolve both refs once to immutable commit identities and acquire the full committed range against those identities, never inventory-path-filtered or worktree-derived. Treat repository/ref values as data through a bounded, sanitized Git adapter, never shell-interpolated commands. Do not fetch missing objects or enable implicit network access. Invalid/foreign roots or syntactically invalid requests use exit 2/`invalid_arguments`; unresolved otherwise-valid refs or missing objects use exit 3/`git_unavailable`, detectable races exit 7/`acquisition_race`, and cancellation exit 6/`interrupted`. Validate the snapshot through the versioned constructor. Report decodes retained exact policy, inventory and result bytes in order, constructs their binding, derives provenance, and validates it. Result decoding rebuilds and compares exact bytes; an isolated digest is not provenance or authentication. Retain antecedents alongside exact stdout bytes externally. Do not substitute legacy `Account` or use report-oriented `DecodeCanonical` as a policy/inventory/result decoder. Internal project adapters are not public plugin admission.

## Streams and technical outcomes

Success: exit `0`, stdout exactly one existing canonical stage document including its LF, stderr empty. No digest line, wrapper, progress output or authority field. Consumer hashes exact stdout including LF. Buffer and validate the complete bounded output before normal emission. Family help is discovery text, not evidence output.

Proposed family-only stable failures: stderr is exactly `w1w2 <stage>: <token>\n` (use `command` if no stage is recognized), no raw paths/content or subprocess stderr. Select the first failure in order: grammar, platform/budget readiness, input decoding, acquisition/binding, output construction/write.

| Exit | Tokens and meaning |
|---|---|
| 2 | `invalid_arguments`, `invalid_document`, `unsupported_version`, `unsafe_path`, `antecedent_mismatch`: invalid request/contract |
| 3 | `input_unavailable`, `git_unavailable`, `unsupported_platform`: acquisition unavailable |
| 4 | `resource_limit`: a declared finite budget exceeded |
| 6 | `interrupted`: cancellation/deadline |
| 7 | `acquisition_race`: detectable Git-object or selected-source race under the approved acquisition model |
| 9 | `output_incomplete`: short/failed stdout write |

Before emission, any failure leaves stdout empty. A short/failed write may already have emitted bytes: they cannot be undone; exit 9 means incomplete technical output, never success. Diagnostic write failure cannot guarantee stderr delivery but must not produce exit 0. Valid zero-change/zero-selection and threshold-exceeded or unavailable observations remain data with exit 0; unavailable observation data is distinct from acquisition failure. Existing commands keep their existing streams, diagnostics and exit meanings except the declared additive root-help entry; this is not a universal remapping of existing 2–8 codes.

## Initial Linux resource profile

The owner selected these initial DESIGN limits; none is an implemented or calibrated W1/W2 default. They govern only the proposed shell, not core accounting thresholds or existing APIs/commands. One MiB is 1,048,576 bytes; one KiB is 1,024 bytes.

| Surface | Selected initial limit |
|---|---:|
| Path-selection request | 1 MiB |
| Each canonical input document | 16 MiB |
| All input documents/requests per operation, including selection | 32 MiB |
| Selected paths | 128 |
| Each decoded relative path | 4095 bytes |
| All decoded paths | 512 KiB |
| Each raw path component | 255 bytes |
| Admitted regular-file or leaf-link content per entry | 1 MiB |
| Aggregate admitted content | 8 MiB |
| Each Git response | 16 MiB |
| Aggregate Git responses per operation | 64 MiB |
| Full committed-diff entries | 10,000 |
| Canonical stdout document | 16 MiB |
| Cooperative operation deadline | 30 s |

Count actual document/request transport bytes, including JSON whitespace, Base64 overhead and canonical LF, toward input limits. Path limits instead count decoded raw bytes, including separators for a whole path and excluding them for each component. Charge acquired content per selected path, including link payloads and repeated content through distinct hardlinks; never deduplicate inode/content accounting. Overflow detection and read-ahead must themselves remain bounded and accounted, not bypass an aggregate. The table bounds admission, not physical disk traffic or process RSS. Input aggregates combine: three individually admissible documents can still exceed 32 MiB; individual maxima do not guarantee combined admissibility.

Exceeding a size/count budget fails with exit 4/`resource_limit`, before normal emission, with no truncation or synthetic empty success. Validate bounds before unbounded reads/allocations; calibration and allocation/overflow checks remain future work. No adaptive limits, new configuration/profile loader or limit-override flags are introduced. The public census defaults are 256 files, 1 MiB/file, 8 MiB total and 10,000 matches (`census/query.go:24–32`); they explicitly do not bound parser resources. Only the content scale is reused as a conservative design choice. Git blob/response caps of 16 MiB exist at `internal/git/blob.go:14–18` and `internal/git/snapshot.go:93`; legacy projection's 16 MiB cap is at `cmd/git-change-evidence/main.go:22`. These different scopes neither validate this new profile nor turn its 10,000 diff-entry cap into a census-match limit.

Start the cooperative deadline before document I/O, propagate cancellation through owned acquisition/Git subprocesses, and test the cancellation strategy before claiming conformance. A detected pre-emission cancellation/deadline uses exit 6/`interrupted`; failed emission remains exit 9/`output_incomplete`. This is not a hard wall-time guarantee under blocked OS/stdio, nor a hard RSS limit: decoding, copies and allocation overhead can exceed admitted byte counts. The caller/harness must supervise any hard time/memory budget it requires; no new supervisor infrastructure is supplied here. Forced external termination cannot promise a CLI diagnostic, stable family exit or complete emission, even if captured bytes look canonical. Record the actual outcome rather than confirmed CLI success.

## Versions and implementation readiness

Platform, initial limits and the scoped threat model are selected for this design. Qualification, representative-load calibration and concrete acquisition/cancellation proof are not supplied by that selection. Local implementation has been separately authorized; the [implementation boundaries](apply-progress.md) carry initial, uncalibrated work-unit forecasts. This is not release, distribution, publication or consumer-adoption authority.

Core document versions, proposed CLI request/protocol v1 and planned source/dev preview identity are independent. Unsupported document/request versions fail explicitly; no guessing, downgrade or legacy projection fallback. Breaking request/stream semantics require a separately identified protocol revision. This draft creates no schema file, published package, release or install channel.
