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

## Guidance

Contributors should read [AGENTS.md](AGENTS.md). Consumer integrators should read [the Consumer Integration Guide](docs/consumer-integration-guide.md).

## License

This repository is licensed under [Apache-2.0](LICENSE). That selection does
not settle separate attribution, intellectual-property, confidentiality, or
other legal questions; those remain owner-controlled review matters.
