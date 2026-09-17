# Consumer Integration Guide

This guide is for consumers of `git-change-evidence` (GCE), not contributors.
The CLI and its machine-readable contracts are the primary integration surface;
the public root Go API and public `census` subpackage remain supported. GCE
reports evidence only. It never approves, rejects, gates, blocks, merges,
deploys, releases, or assigns delivery authority.

Use this sequence: **consume → validate → reproduce → interpret**. A technical
result is not delivery approval, and consumer-owned workflow decisions must not
be attributed to GCE.

## 1. Consume

Capture an evidence result as bytes, rather than as reformatted JSON. Preserve
its exact canonical document bytes, including the final LF, and retain the
published SHA-256 digest with the source of the bytes. Do not pretty-print,
transcode, trim, or append data before validating a digest.

Capture stdout, stderr, and the process exit status separately. A successful
machine-readable result is written to stdout; diagnostics belong on stderr.
A nonzero exit, an interrupted write, or a truncated byte stream is not an
empty successful result and must not be filled in with invented evidence.

Paths such as `path_b64` use standard Base64 of raw path bytes. Preserve those
Base64 values as wire data until a byte-aware consumer decodes them. Do not
assume that a path is UTF-8, normalize it, or substitute a display string for
its raw-byte identity.

A census document has the established envelope and provenance fields; this is
an example fragment, not a new capabilities schema:

```json
{
  "schema": "git-change-evidence.census-go-ast-cli/v1",
  "extractor_version": "git-change-evidence.go-ast-census/v1",
  "inventory": {
    "schema": "git-change-evidence.inventory/v1",
    "sha256": "<inventory-sha256>",
    "accounting_policy_sha256": "<fixed-census-policy-sha256>"
  },
  "source_scope": [{
    "path_b64": "ZXhhbXBsZS5nbw==",
    "content_sha256": "<file-sha256>",
    "byte_length": 42
  }]
}
```

For the complete current census command and JSON shape, use
[`docs/census.md`](census.md). This guide does not add a CLI command that
assembles Policy, Inventory, Result, or Report documents.

## 2. Validate

Validate canonical documents against their declared existing schema and
contract-specific decoder before relying on their contents. For Policy V1,
Inventory V1, Result V1, and Report V1, retain the exact canonical bytes and
check the SHA-256 digest over those bytes. A report's digest syntax alone does
not prove its antecedents or provenance.

For revision-bound Result and Report evidence, retain the immutable native Git
object IDs for both `base` and `head`, without shortening, replacing, or
rehashing them. Keep the exact canonical Policy, Inventory, and Result
antecedents and their digests. Report provenance binds the accounting-result,
input-policy, and inventory digests to those revisions; validate that binding
against the retained antecedents rather than trusting an isolated report.

Use the public contract APIs when integrating in Go. In particular, a consumer
can validate a Report V1 and its retained `ResultV1ReportBinding` with
`ValidateReportV1ResultProvenance`. This validates existing provenance; it does
not make the report a delivery decision.

## 3. Reproduce

Independently regenerate revision-bound evidence from the same immutable base
and head IDs, the same policy, and the same supported API inputs. Compare the
resulting canonical bytes and digests with the retained evidence. Do not infer
that a matching branch name, moving ref, checkout state, or a standalone report
is an equivalent antecedent.

Census is different: it accepts explicit supplied source bytes and has no Git
revision input. Regenerate it with the same query, limits, selected paths, and
captured file bytes. Rerunning against a changed checkout is not reproduction
of prior census evidence.

`gce census go-ast` is the canonical command, and
`git-change-evidence census-go-ast` remains the equivalent pre-v1 command.
Census is syntax-only over the supplied file scope. It performs no semantic
resolution and is not an atomic whole-tree snapshot: selected files can be
captured at different instants. Its inventory, file hashes, byte lengths, query,
extractor version, and Base64 raw paths identify the supplied scope.

## 4. Interpret

Interpret exit status and streams as technical observations, not workflow
verdicts:

| Census condition | Exit | Meaning |
| --- | ---: | --- |
| successful match or zero matches | 0 | A census document was produced for the supplied scope. Zero matches is not semantic absence, safety, security certification, or approval. |
| malformed command, invalid path, or parse failure | 2 | Input is invalid; stderr is `invalid input`. |
| unavailable, missing, symlinked, or nonregular source | 3 | Content is unavailable; stderr is `content unavailable`. |
| file, byte, or match bound exceeded | 4 | Resource bound reached; stderr is `resource bound exceeded`. |

Census is Linux-confined; on non-Linux platforms its source content is
unavailable rather than read with a different confinement model. Its normal
successful stdout is one JSON document with a final LF and empty stderr.

Other existing CLI commands have command-specific technical exits. Where their
respective command path documents them, exit 5 means invalid UTF-8, exit 6
means operation interrupted, exit 7 means a Git object race, and exit 8 means a
publication conflict. Do not apply those mappings to census or assume that an
exit code from one command creates a universal evidence schema. Local
`publish`, where implemented, stores evidence; it is not product-release
publication.

Treat stale, incomplete, mismatched, unavailable, noncanonical, or nonzero
results as a reason to retain the observed state and independently reacquire or
regenerate applicable evidence. The consumer—not GCE—decides whether to stop,
escalate, request fresh inputs, continue a review, merge, deploy, or release.
No evidence result is approval, rejection, gating, blocking, or any other
assignment of delivery authority.

## Post-preview roadmap

Multilingual census is a separate post-preview roadmap and is not a preview
blocker. The planned language sequence is: Python → JavaScript/TypeScript → Vue
→ Java → Rust. This guide does not implement or imply multilingual behavior.
