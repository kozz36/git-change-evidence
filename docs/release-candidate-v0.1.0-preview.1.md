# v0.1.0-preview.1 release-candidate preparation

This runbook prepares evidence for the planned source/dev preview identity
`v0.1.0-preview.1`. It is not publication guidance. No tag, release, package,
asset, upload, distribution channel, or public-visibility change is claimed or
authorized by this document.

## Read this first

**Preparation baseline:** `origin/main@16aaa4498f91c6d1ff91b330e1df6a47db26c922`
with tree `e1926f0fafd43e56e0c719cd23afe0a1bf96fb33`.

This baseline records preparation inputs only. A future final reviewed release
commit and its future tag must be selected by the owner before any publication.
Task-6 ephemeral builds are rehearsal candidates only and **MUST be regenerated
from the future reviewed/tagged commit before any upload**. Do not upload,
publish, tag, release, or retain a rehearsal as a release asset.

## Scope at a glance

| Topic | Recorded decision |
| --- | --- |
| Preview wording | Source/dev preview; not release provenance. |
| Candidate source | Source archive plus Linux `amd64` and `arm64` archives. |
| Linux archive members | Every Linux archive contains both `gce` and `git-change-evidence`. |
| Integrity record | `SHA256SUMS` covers the source archive, both Linux archives, and the SBOM. |
| SBOM | CycloneDX SBOM planned with temporary pinned `cyclonedx-gomod@v1.10.0`; it must not be installed into the repository. |
| Publication | Owner-controlled future work, not part of preparation or rehearsal. |

## Supported surface and limits

The primary consumer surface is the CLI and its machine-readable contracts.
The supported public Go surfaces are the module-root API and the public
[`census`](../census) subpackage. The canonical pre-v1 census command is
[`gce census go-ast`](census.md); the retained pre-v1 compatibility command,
`git-change-evidence census-go-ast`, has the established equivalent contract
meaning and technical behavior.

The Go AST census is syntax-only over explicit supplied Go file bytes. It does
not do semantic resolution and is not an atomic whole-tree snapshot. It is not
a multilingual census implementation. See the [census guide](census.md) for
input limits, Linux confinement, output, and exit behavior.

`git-change-evidence` reports evidence only. It does not approve, reject, gate,
block, merge, deploy, release, or assign delivery authority. Consumers retain
all workflow and delivery decisions; see the [Consumer Integration Guide](consumer-integration-guide.md).

The source/dev version output is intentionally not release provenance. When
installed with `go install`, the emitted executable name is
`git-change-evidence`, not the canonical `gce`; build or rename the canonical
local executable explicitly when that name is required.

The legacy command remains supported throughout pre-v1. Broad distribution,
release automation, signing, and provenance are deferred. The repository is
licensed under [Apache-2.0](../LICENSE), which does not settle attribution,
intellectual-property, confidentiality, export, or other owner-controlled
questions.

## Candidate manifest plan

The planned manifest is source plus two Linux architecture records, an SBOM,
and `SHA256SUMS`:

```text
git-change-evidence_v0.1.0-preview.1_source.tar.gz
git-change-evidence_v0.1.0-preview.1_linux_amd64.tar.gz
git-change-evidence_v0.1.0-preview.1_linux_arm64.tar.gz
git-change-evidence_v0.1.0-preview.1.sbom.cdx.json
SHA256SUMS
```

Record the selected source commit and tree, version wording, tool versions,
byte sizes, and SHA-256 values. Record `GOOS=linux` and the relevant `GOARCH`
for each binary archive. Each Linux archive must contain exactly the two
executable names `gce` and `git-change-evidence`; inspect the members and
architecture before accepting rehearsal evidence. The SBOM must identify the
same source and link each archive by filename and SHA-256.

## Safe preparation checks

These inspection commands do not publish, create candidate assets, mutate
repository configuration, or install SBOM tooling:

```sh
git rev-parse origin/main
git rev-parse origin/main^{tree}
git diff --check
git status --short
git diff --name-only
```

Confirm the first two commands match the preparation baseline and tree above.
Use this rehearsal checklist before scheduling Task 6; do not interpret a
passing check as release authorization.

- [ ] The preparation baseline commit and tree still match the values above.
- [ ] Rehearsal output will be created only outside the repository in a fresh temporary directory.
- [ ] No candidate asset will be uploaded, retained as a release asset, or copied into the repository.
- [ ] `SHA256SUMS`, archive members, binary architectures, SBOM coverage, and CLI contract checks will be recorded.

Before any future upload—not before the Task-6 rehearsal—the owner must select
the final reviewed release commit and tag, and every upload candidate must be
regenerated from that selected commit.

## Future rehearsal procedure (Task 6 only)

Do not execute this section as part of documentation preparation. It describes
a future, non-publishing rehearsal. Create a new temporary working directory
outside the repository and remove it after evidence capture. Do not write any
generated archive, checksum, binary, manifest, or SBOM into the repository.

1. Verify the selected source commit and tree. For the preparation rehearsal,
   use the baseline above; for any upload candidate, use only the future
   reviewed/tagged commit.
2. Produce the source archive and Linux `amd64` and `arm64` archives in the
   temporary directory. Put both `gce` and `git-change-evidence` in each Linux
   archive.
3. Generate a CycloneDX SBOM with temporary pinned
   `cyclonedx-gomod@v1.10.0`. For example, set `GOBIN` to a tools directory
   under the temporary directory before running the pinned tool; do not add it
   to `go.mod`, `go.sum`, or any repository path. Record the generator version
   and the SBOM's source and archive linkage.
4. After the SBOM exists, generate `SHA256SUMS` in bytewise filename order so
   it covers the source archive, both Linux archives, and the SBOM; then verify
   the checksum file from the temporary directory.
5. Verify archive members, Linux binary architecture, checksums, SBOM coverage,
   `--version` source/dev wording, and canonical/legacy census parity where
   executable.
6. Capture only the observed evidence, then delete the entire temporary
   directory. A later upload requires a fresh rebuild from the future reviewed
   release commit and tag.

## Observed Task 6 rehearsal evidence

The non-publishing rehearsal completed from the preparation baseline on
2026-09-17. It used a temporary Git index and object directory so the real
index was never changed. The synthetic rehearsal tree included the authorized
working-tree state before this evidence section and the Task 6 checkbox were
updated:

- baseline commit: `16aaa4498f91c6d1ff91b330e1df6a47db26c922`
- baseline tree: `e1926f0fafd43e56e0c719cd23afe0a1bf96fb33`
- rehearsal tree: `846d84b88104b7ed8681e7d9aae9caa4d8e85cd3`
- authorized dirty paths: `.atl/skill-registry.md`, `.gitignore`,
  `CHANGELOG.md`, `docs/release-candidate-v0.1.0-preview.1.md`,
  `odd/tasks/prepare-v0.1.0-preview.1.md`,
  `openspec/changes/.gitkeep`, `openspec/changes/archive/.gitkeep`, and
  `openspec/specs/.gitkeep`
- current-module/rehearsal-tree parity: 99 Go/module files, aggregate SHA-256
  `3c5ef09cd815b13b181bdda828fdf04e8f40bf946690aa8232e9cecb668bfe23`

The rehearsal used `go1.27.1-X:nodwarf5 linux/amd64` and a temporary install of
`github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@v1.10.0`
(module sum `h1:9Vy3zcC+lJLgcR4xYQvwPGU6L2Rij/Ld47lyucYjVI0=`). Each target SBOM was
generated directly from that target's exact `gce` binary with `bin -json
-noserial -notimestamp -version v0.1.0-preview.1`. The archived
`git-change-evidence` binary was byte-identical to `gce`. No SBOM JSON
post-processing occurred. License collection was not requested because `bin`
mode reads binary build metadata and the unpublished main module is not
available for download. The generated manifest recorded this limitation and
the binary-to-SBOM hashes.

The two per-target CycloneDX SBOMs are the final owner-authorized Task 6 path
for this observed rehearsal. They differ from the earlier procedure's singular
SBOM shorthand and the Task 4 audit-time SPDX shape; those earlier passages
remain historical planning, not descriptions of the observed output. Exact
binary hashes were retained only in the discarded manifest whose hash is
recorded below, so that linkage is an observed check rather than independently
recomputable evidence after cleanup.

The observed candidate files were:

| File | Bytes | SHA-256 |
| --- | ---: | --- |
| `git-change-evidence_v0.1.0-preview.1_source.tar.gz` | 385139 | `e71b618e3974cdbc6f27dcfde0e489b0508e3f233bbf56cdb7c69c9b8c5e60ed` |
| `git-change-evidence_v0.1.0-preview.1_linux_amd64.tar.gz` | 6832002 | `8f4fcffc6045676824e544c044e7de0a7e241ac6bb979b26605a20629b18c81b` |
| `git-change-evidence_v0.1.0-preview.1_linux_arm64.tar.gz` | 6300829 | `934d15595ab98ae0336fdc5cdb8934d9b7d81ea3fbc222b1f8111f822748fb31` |
| `git-change-evidence_v0.1.0-preview.1_linux_amd64.sbom.cdx.json` | 4389 | `4df462854afff8611a65db2296424996a391dff41ce385eb2bc3a00930b63443` |
| `git-change-evidence_v0.1.0-preview.1_linux_arm64.sbom.cdx.json` | 4389 | `e9166ee87c5b7c422f6b422571e90ef306a916c570a55a2199da48954d2f4488` |
| `git-change-evidence_v0.1.0-preview.1.manifest.json` | 4382 | `73f9dcdd92e6ed4097b59aecdbb7ba896c764b631e81a1f262fa6640c698ad0a` |
| `SHA256SUMS` | 619 | `5bd592ad8fcf51b86706198b6d63462c1e62118991c0311e214b6ad0efa2c29f` |

`SHA256SUMS` covered and successfully rechecked the source archive, both Linux
archives, and both target SBOMs in bytewise filename order. It intentionally
did not list itself or the manifest. The source archive contained 240 members
under `git-change-evidence-v0.1.0-preview.1/`. Each Linux archive contained
exactly `gce` and `git-change-evidence`. `file` identified the binaries as
static ELF 64-bit `x86-64` and `ARM aarch64` executables. The `go version -m`
record SHA-256 values were
`e23e8db968bc471d8fdde1a07714ef74d7b6963f33b6d8517d3e9bcb9094cb4e`
for `amd64` and
`b0aee9a2d968710adccd2f9e349c25c63081270d578ea2d94768e38d86c5daff`
for `arm64`.

Both SBOMs were unmodified CycloneDX 1.6 JSON documents with main component
`github.com/kozz36/git-change-evidence`, version `v0.1.0-preview.1`, one
component, two dependency records, no timestamp, and no absolute local path.
The `amd64` runtime checks observed the exact source/dev version string. Root
help was identical under both binary names (SHA-256
`3aeb5453908c8987575075ec07967ed029ae8652d14ceba698ba2b3410293dda`).
Canonical and legacy census help were identical (SHA-256
`e1505a9d572533959558bbab784f93c0c8c973259461ac812a224c019ced3405`).
Their non-empty JSON output was byte-identical: 723 bytes, SHA-256
`9541afcc4e81ddf289aaafb906ecc6dccf00d18f5547aab60e702a0ca6a457f1`,
with empty stderr. `arm64` runtime execution was not applicable because no
native runner or emulator was available; its architecture and build metadata
checks passed instead.

Before and after the rehearsal, the real index remained SHA-256
`278dfbbc15c241079b44f9ef373b65da65483b05d96e2d600e21fe316e5a4cb7`
and Git blob `ffe6fd319c2bb5a6e2875e6e41e2fd16c0846c75`. `go.mod` and `go.sum`
remained SHA-256
`5ff081136d6b1e7539438f52ea2beb1fce822441194e540e4b91c782a2555a6b`
and `28b936a5baa45bf6328134f4ae47a6e60a8a8a0ed4cd02549b2460a447764e8b`.
The owner-only temporary directory, isolated tool, caches, index, objects, and
all generated candidates were deleted, and the cleanup trap verified its exact
runtime path was absent. The random temporary basename was not retained; a
subsequent scan found no rehearsal/recovery directory matching the bounded
prefixes. Nothing was copied into the repository or uploaded. Because deletion
was mandatory, the candidate hashes and content checks above are recorded
observations and cannot be independently recomputed without a fresh rehearsal.
They do not describe downloadable assets or release provenance. Any future
upload candidate must be regenerated from the future reviewed release commit
and tag.

## Active change and archive blocker

`bootstrap-neutral-core-extraction` remains active and resumable at **3/12
tasks**. Its recorded `nextRecommended` value is `apply`. Do not complete,
archive, delete, or rewrite that change during release preparation.

The `go-ast-census` archive is not complete. Two canonical executor attempts
each ended after **1 turn, 0 tool calls, and 0 writes**. A separate native
STATUS v2 reports that archive state is ready. No manual fallback was used; the
retry budget is exhausted. Preserve this blocker and do not claim archive
completion.

## Owner-controlled publication gates

Preparation and rehearsal do not decide any of these gates:

| Gate | Required owner decision |
| --- | --- |
| Legal and rights | Legal, IP, confidentiality, attribution, and export review. |
| Repository state | Git history and unpublished refs. |
| External controls | External settings, secrets, and access controls. |
| Visibility | Public repository visibility. |
| Release identity | Final reviewed commit and tag. |
| Distribution | Release creation, upload/assets, and package or other distribution. |
| Supply-chain claims | Signing and provenance. |

No technical report or checklist result grants authority for these decisions.

## Verification and recovery

Before accepting documentation evidence, verify local Markdown links, the
baseline and tree identifiers, source/dev wording, both Linux architecture and
binary names, `SHA256SUMS`, `cyclonedx-gomod@v1.10.0`, the archive-blocker
statement, and the absence of publication claims. Also run `git diff --check`
and inspect the changed-path and changed-line census.

If the documentation is incorrect, revert only the preparation documents and
repeat their checks. If rehearsal artifacts are incorrect or incomplete, delete
the temporary directory and regenerate it; never repair or upload a partial
candidate. Neither recovery path authorizes publication.

## Boundaries

This runbook distinguishes preparation and rehearsal from publication. It does
not create a tag, release, package, asset, upload, public repository, or
visibility change, and it does not authorize any of those actions.
