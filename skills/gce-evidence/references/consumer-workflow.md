# GCE consumer workflow

Use **consume → validate → reproduce → interpret**. This is checkout-local guidance, not a packaged skill or a change to the published preview. Skill `metadata.version` identifies skill content only, not the CLI release, extractor, schema, or a compatibility/distribution policy. Start with the [consumer integration guide](../../../docs/consumer-integration-guide.md); use each family's existing contract, not an invented universal decoder or assembler. Local CLI `publish` stores evidence, not a product release.

## Select and capture

| Surface | Contract and inputs |
| --- | --- |
| CLI Go AST census | Prefer `gce census go-ast --source-root <absolute-root> --receiver <identifier> --selector <identifier> -- <explicit-relative-go-paths...>`; `git-change-evidence census-go-ast` is equivalent throughout pre-v1. Require each flag once, `--`, and at least one selected relative `.go` path; no implicit root, Git revision, shell expansion, discovery, or invented CLI limit flags. Check local executable availability before use. See [census command](../../../docs/census.md) and [census spec](../../../openspec/specs/go-ast-census/spec.md). |
| Public Go API | Use `github.com/kozz36/git-change-evidence` family-specific constructors/decoders, and `ValidateReportV1ResultProvenance(report, binding)` for applicable retained Report bindings ([root contract](../../../contract.go)). The public `github.com/kozz36/git-change-evidence/census` `GoASTV1(inventory, files, query, limits)` accepts complete matching raw path/content records, valid inventory/query, and positive limits ([query API](../../../census/query.go)); a valid explicitly empty inventory is allowed by API, unlike CLI selection. Not every family has a CLI assembler. |

Retain original canonical document bytes through final LF, applicable digest and its byte source, immutable base/head Git IDs and canonical Policy/Inventory/Result antecedents for revision-bound evidence. Preserve raw-path standard Base64, not normalized/display text; non-UTF-8 paths are possible. Capture argument vector, stdout, stderr, exit status separately for processes; for API calls retain typed inputs and result/error instead of fabricating streams. Record the census absolute source root alongside the invocation: it is absent from output. Validate declared schema and family-specific provenance, not just JSON parseability or digest syntax. Selected files can be captured at different instants; CLI bounds do not guarantee parser isolation.

## Reproduce and interpret

| Evidence | Reproduction requirement |
| --- | --- |
| Census | Original captured source bytes and raw selected paths, query, effective limits, compatible extractor behavior/identity; compare output bytes and applicable digests. File hashes alone cannot restore lost bytes. A later reread is not automatically identical. |
| Revision-bound | Exact immutable base/head IDs, canonical antecedents and their digests, and supported API inputs; compare regenerated canonical bytes/digests. Moving refs or a convenient checkout are not substitutes. |
| Originals missing or changed | Reproduction cannot be established. An authorized correction/reacquisition with changed inputs produces **new evidence**, even if counts match. |

Census candidates are exact textual syntax matches, including shadowed spellings or ambiguous indexed forms, only in the validated supplied scope. Zero matches is scoped syntax-only success, not semantic absence, whole-tree completeness, safety, or approval. Do not assign any delivery decision from technical evidence.

## Failures and recovery

| Observation | Technical response |
| --- | --- |
| Census exit 0 | Validate the complete single JSON document with final LF and empty stderr; zero matches is valid. |
| Census exit 2 | Invalid input (including parse/path errors); retain diagnostic, request transparent correction. |
| Census exit 3 | Content unavailable (missing, symlinked, nonregular, non-Linux), or output-write failure; exit alone does not prove missing source. Do not bypass Linux confinement. |
| Census exit 4 | Resource bound exceeded; do not silently split scope or increase limits. CLI defaults: 256 files, 1 MiB/file, 8 MiB total, 10,000 matches. |
| Interruption, nonzero with write prefix, malformed/truncated/noncanonical bytes, stale binding | Preserve raw failure and mismatched antecedents; no prefix or missing result is empty success. Seek exact originals or authorized new acquisition. |

Other CLI commands have their own documented exit mappings; do not generalize census exits. Name the original failure, correction, changed inputs, authorization, and observed result. Do not retry indefinitely or access unauthorized files. Static reading of this guidance is not an agent-behavior evaluation; actual skill loading and observed traces require separate acceptance work.
