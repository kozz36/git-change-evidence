# git-change-evidence

`git-change-evidence` is a planned generic, multi-project Git evidence CLI and Go package. It will produce deterministic change evidence from immutable Git revisions and project-specific profiles.

## Status

**Bootstrap decision complete; implementation has not started.** The repository intentionally contains no Go source, module identity, packaging, CI, or release automation yet.

The product is evidence-only: it may report observations and warnings, but it will never approve, reject, block, or otherwise own delivery decisions.

## Bootstrap decisions

- Go 1.25.10 is the target baseline. The future public neutral API is at the repository/module root, the CLI is at `cmd/git-change-evidence/`, and non-public shell/adapters are under `internal/`.
- The future primary gate is `go test ./...`; concurrency-bearing surfaces also run `go test -race ./...`; formatting is check-only as recorded in `openspec/config.yaml`.
- The core receives immutable typed Go structs and interfaces. TOML, YAML, and JSON decoding remain adapter concerns; no concrete configuration-file syntax is selected.
- Module/publication identity, packaging, CI, license, release, and distribution remain deferred.

## Product charter

The authoritative scope, architecture, governance, and milestones are in [docs/product-charter.md](docs/product-charter.md). The first intended consumer is a CNSIC profile; its policy is not part of this product's generic core.
