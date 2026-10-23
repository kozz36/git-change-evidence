# Approved design — Independent committed-snapshot contract

**Concrete contract approved for isolated implementation under S5–S6; not yet a delivered CLI.** The names, field/byte-metric semantics, association limits and initial admission profile below are the implementation target. See [apply progress](apply-progress.md) for actual evidence. See [requirements](specs/w1-w2-neutral-git-facts-export/spec.md) and [future work](tasks.md).

## Source facts and alternatives

Baseline is main `edeffba9dcada6c5512e3f7d623d72bbe39adfc5`. Relevant sources: [snapshot type](../../../snapshot.go), [bounded acquisition](../../../internal/git/snapshot_bounded.go), [result validation](../../../result_v1_antecedents.go), [result wire](../../../result_v1_wire.go), [report binding](../../../contract.go), [parser](../../../internal/w1w2/input/parse.go), [dispatcher](../../../cmd/git-change-evidence/main.go) and [profile](../../../internal/w1w2/contract/profile.go).

| Established fact | Design consequence |
|---|---|
| Snapshot/acquisition retains status, paths, modes/kinds, IDs and numeric numstat; raw tokens and side LOC are missing. | Add bounded side measurement/token retention; do not claim direct serialization suffices. |
| Result validation checks revision/path/accounting projection, not every side fact. | A complete new facts validator is necessary; old validation is not a certificate. |
| Result/report V1 omit side facts; report binding contains policy/inventory/result/revisions only. | Add a distinct contract, preserve V1 and describe association without inventing report coverage. |
| Parser/root dispatcher has exact stage grammar and retained legacy triples. | Add only the proposed stage; no reinterpretation of old invocations or stage-help triples. |

A separate document avoids coupling Git facts to accounting or policy. A new result version was the alternative, not the selected direction. Existing public Go types are useful inputs, not existing CLI export contracts. Copy status is representable when actually acquired; parser recognition does not promise expanded copy detection or similarity-score preservation.

## Approved CLI

Here `gce` denotes an eventual approved executable, not installed availability.

```text
gce w1w2 snapshot --repo <absolute-repository> --base <revision> --head <revision>
```

Exactly-once separate flag/value pairs, any order; literal values may begin with `--`. No defaults, positionals, equals-form flags, policy, inventory, path selection, output file or stdin document. Unused stdin is not read. Family help remains `gce w1w2 --help`; root legacy triples, including existing help-shaped triples, keep their routing.

Resolve explicit refs once to native immutable commits and acquire the **whole range**. Reuse the bounded Git adapter semantics: no implicit fetch, external diff/textconv, worktree/index substitution or selected-path filtering. Side reads use immutable blob IDs, not working-tree paths. Do not dereference gitlinks or follow symlink targets.

## Approved document and API

Approved schema: `git-change-evidence.committed-snapshot/v1`. Constructor/decoder and immutable document views belong at the module root; exact Go identifiers/signatures remain design candidates, not current exports. Use a new full-facts input/view rather than adding required fields or changing semantics of existing public snapshot/result types.

| Object | Ordered keys and semantics |
|---|---|
| Document | `schema, revisions, entries`. No policy/category/inventory fields or self-referential digest. |
| Revisions | `base, head`: native lowercase 40/64-character IDs, same object format. |
| Entry | `status, old, new, numstat`. Status is raw Git code `A/D/M/R/T/C`, not a project-specific classification; no similarity score. |
| Side | `path_b64, mode, kind, object_id, content`. Canonical standard Base64 raw relative path, raw mode, native nonzero object ID. |
| File side | Modes `100644/100755`, kind `file`; content is a measured object. |
| Other side | Mode `120000`/`160000`, kind `symlink`/`gitlink`; content is null, not zero or a followed-target measurement. |
| Content | `byte_length, contains_nul, logical_lines`: unsigned integers/bool describing bounded blob bytes, not accounting LOC. |
| Numstat | `kind, added_raw, deleted_raw, added, deleted`. Text: canonical decimal tokens and matching integers; binary: both tokens `-` and both numeric values null. |

Approved byte metric: empty content has zero logical lines; otherwise count LF bytes plus one iff the last byte is not LF. Do not normalize CRLF or decode Unicode. Record NUL independently; even a NUL-containing file has a byte metric, **not a claim of countable text LOC**. Git numstat binary state is an entry observation, not proof of each side's content classification. Consumer LOC/classification equivalence remains unproven.

Numeric measurements are JSON integers in 0..2^64−1; reject booleans, fractions, exponents, negatives and overflow. Text numstat tokens are `0` or a nonzero ASCII decimal digit followed by decimal digits, and equal their numeric values. Native side IDs have the revision format's width and are not all-zero placeholders. Empty file content requires zero bytes/lines and contains_nul=false; nonempty content requires 1 <= logical_lines <= byte_length. Structural consistency cannot prove these measurements were honestly acquired.

Decoded paths are nonempty relative slash-separated bytes: no NUL, leading/trailing slash, empty component, `.` or `..` component. Backslash is a literal byte, not a separator; non-UTF-8 bytes are preserved in Base64. No case folding, path expansion or host-path resolution occurs in core validation.

An addition has old=null/new present; a deletion the reverse. M/R/T/C require both sides. M/T retain the same raw path; R/C retain distinct paths. Modes/kinds, object-ID widths, status and absence must agree. Validate zero-count mode-only cases without inventing a separate status. Unsupported kinds or inconsistent acquisition must fail, not be coerced to a file. Valid absent sides/non-file content are different from missing required measurements.

Canonical order is ascending decoded effective path (new path if present, otherwise old), with unique effective paths; shared old paths are not globally forbidden. Raw identities are not display strings and are never Unicode-normalized or used as host paths by decoding.

Synthetic schema example only, **NOT acquired or executed**; IDs are illustrative:

```json
{"schema":"git-change-evidence.committed-snapshot/v1","revisions":{"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},"entries":[{"status":"A","old":null,"new":{"path_b64":"YS50eHQ=","mode":"100644","kind":"file","object_id":"cccccccccccccccccccccccccccccccccccccccc","content":{"byte_length":2,"contains_nul":false,"logical_lines":1}},"numstat":{"kind":"text","added_raw":"1","deleted_raw":"0","added":1,"deleted":0}}]}
```

Canonical UTF-8 JSON uses the declared key order, compact encoding and exactly one LF. Reject duplicate/unknown keys, noncanonical encodings/tokens/numbers, invalid Base64, unsafe relative names, duplicate effective paths and cross-field inconsistency. SHA-256 covers exact document bytes including LF and is exposed separately. Constructors/decoders are pure, deep-copy mutable data and validate structure/consistency; they do not prove supplied facts came from Git or open repositories.

## Acquisition and resource proposal

New work must acquire every required regular-side blob measurement completely, retaining raw numstat tokens before parsing loses them. It must share invocation-local budgets across raw diff, numstat, object checks and side reads; independently resetting a budget per phase is not aggregate enforcement. Repeated reads charge actual bytes; no new cache or daemon is proposed.

| Approved initial shell ceiling, not workload-calibrated | Value |
|---|---:|
| Each Git response, including a side blob | 16 MiB |
| Combined Git response bytes across all phases | 64 MiB |
| Committed entries | 10,000 |
| Canonical stdout, including LF | 16 MiB |
| Cooperative operation deadline | 30 s |

These reuse existing ceilings without claiming workload calibration. Selected-inventory limits are not committed-range entry limits. Pure core validation is separate from shell admission limits. Initial CLI capability policy retains the Linux family restriction, not certifying other platforms. Missing objects/capabilities or budgets must be explicit failures; no truncation, clamping, partial success, hidden limit override or weaker acquisition.

Success is exit0, empty stderr and one canonical document. Proposed family technical exits remain 2 invalid input/document, 3 unavailable/platform, 4 resource limit, 6 interrupted, 7 inconsistent acquisition and 9 incomplete output, with sanitized `w1w2 snapshot: <token>\n` diagnostics. Pre-emission failure has empty stdout; exit9 bytes are incomplete and unusable. External termination has no guaranteed family diagnostic. Exact existing token reuse must be verified before implementation.

No hard RSS, blocked-I/O/stdio deadline, atomic multi-read snapshot, all-alias provenance, every-write detection or descendant-reaping guarantee is added. Output excludes full blob payloads, but measuring sides still requires bounded blob **reads**.

## Association, not inherited report provenance

Retain exact snapshot bytes/digest alongside any policy/inventory/result/report documents. A proposed pure association check compares both native revision IDs and complete effective-path membership against the result; it does not mutate or add a snapshot digest to Result/Report V1.

This is range association only: not proof of acquisition authenticity, category/count equivalence, consumer parity or report coverage of snapshot measurements. Existing report validation still checks its existing antecedents only. A combined signed/bundle/provenance contract is not proposed. Callers must retain and identify the additional antecedent explicitly; matching revisions or an isolated digest is not an automatic provenance certificate.

## Implementation review boundaries

The owner approved this grammar/schema, field constraints, raw-token fidelity, byte-metric basis, association limits and initial domain/budgets under S5–S6; any later semantic change needs a new decision. Consumer LOC equivalence remains unproven. Review future root-file/layout threshold crossings before adding Go files; no new public subpackage or infrastructure is selected. Consumer adaptation, packaging/distribution, pin, release and adoption remain separate approvals. Current implementation progress is scoped in apply progress: the pure byte-measurement API is tested, while the complete document/CLI remains pending. Existing public API/CLI/V1 semantics and canonical specifications remain unchanged.
