# git-change-evidence

`git-change-evidence` is an evidence-only tool for reporting technical states,
measurements, and observations from supplied inputs. It does not approve,
reject, gate, block, merge, deploy, release, or assign delivery authority.

## Preview status

`v0.1.0-preview.1` is the planned preview identity for this source/dev build.
It is not a published distribution, release, tag, package, or download asset.

The CLI and its machine-readable contracts are the primary integration surface.
The public Go API at the module root and the public [`census`](census) subpackage
remain supported.

## Build locally

From the repository root with Go 1.25.10:

```sh
go build -o gce ./cmd/git-change-evidence
go build -o git-change-evidence ./cmd/git-change-evidence
```

Use the canonical census form `./gce census go-ast`; the legacy
`./git-change-evidence census-go-ast` form remains equivalent during pre-v1.
See [the census guide](docs/census.md) for executable commands and limits.

## W1/W2 selected-path evidence (local source)

The local source CLI provides three Linux-only stages over the existing neutral
V1 contracts. Use caller-controlled local roots and a stable mount namespace;
selection is not authorization, confidentiality or complete-inventory proof.
Final qualification is tracked in [implementation progress](openspec/changes/w1-w2-cli-public-contract/apply-progress.md), not implied by this source build.

```sh
./gce w1w2 --help
./gce w1w2 inventory --source-root /authorized/source --policy policy.json --paths selection.json > inventory.json
./gce w1w2 result --repo /authorized/repository --base HEAD~1 --head HEAD --policy policy.json --inventory inventory.json > result.json
./gce w1w2 report --policy policy.json --inventory inventory.json --result result.json > report.json
```

Supply each required `--name value` pair exactly once. Document selectors are
explicit files or `-` for stdin; at most one selector may use stdin. Unused stdin
is not read. The selection request carries raw relative names as canonical Base64:

```json
{"schema":"git-change-evidence.cli-path-selection/v1","paths_b64":["YS50eHQ="]}
```

Inventory reads selected regular files or leaf-link payloads without reading
symlink targets or discovering other paths. Result uses the full immutable
committed range, not the inventory selection or worktree. Retain exact policy,
inventory and result bytes with the emitted report; hashes alone are not binding proof.

Success is code `0`, one canonical document including LF on stdout, and empty
stderr. Failure uses `w1w2 <stage>: <token>` on stderr; codes `2/3/4/6/7/9` denote
invalid input, unavailable acquisition, resource limit, interruption, detectable
race and incomplete output. A failed write can leave a partial prefix. See the
[grammar, fixed limits and threat boundaries](openspec/changes/w1-w2-cli-public-contract/design.md).
There is no delivery authority, hard RSS/blocked-I/O deadline certificate or
published distribution. Three-argument legacy forms, including `w1w2 inventory --help`,
remain legacy measurement; only `w1w2 --help` is family discovery.

## Guidance

Contributors should read [AGENTS.md](AGENTS.md). Consumer integrators should read [the Consumer Integration Guide](docs/consumer-integration-guide.md).

## License

This repository is licensed under [Apache-2.0](LICENSE). That selection does
not settle separate attribution, intellectual-property, confidentiality, or
other legal questions; those remain owner-controlled review matters.
