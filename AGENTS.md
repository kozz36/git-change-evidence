# Agent Guide

## Authority

`git-change-evidence` is evidence-only. It may report technical states, measurements, and observations, but it must never approve, reject, gate, block, merge, deploy, release, or assign delivery authority.

## Architecture

Use **Functional Core / Imperative Shell** with **Hexagonal Architecture** adapters.

- The Functional Core accepts immutable typed Go structs and interfaces, validates and assembles neutral contracts, and has no process, filesystem, configuration-parsing, or project-policy effects.
- The Imperative Shell owns Git and filesystem effects, CLI adaptation, retries, diagnostics, and local publication.
- Profile adapters translate project-specific vocabulary, categories, thresholds, and compatibility needs one way into neutral core inputs. They cannot alter neutral contract meaning or introduce delivery authority.
- Configuration decoding belongs only in adapters. Core inputs are already typed and immutable; no configuration-file syntax is selected yet.

## Go Layout and Testing

- Public neutral Go APIs belong at the repository/module root, except for the narrowly approved public `census` subpackage.
- The CLI belongs in `cmd/git-change-evidence/`; non-public shell and adapters belong in `internal/`.
- Tests are co-located as `*_test.go` files. Do not add `pkg/`, a Python-style `src/` plus `tests/` layout, or package-level `AGENTS.md` files.
- Target Go baseline: **Go 1.25.10**. Until a later work unit establishes a module identity, this planning slice intentionally contains no `go.mod` or Go source files.
- For Go implementation work, follow RED → GREEN → TRIANGULATE → REFACTOR. The primary gate is `go test ./...`; use `go test -race ./...` for concurrency-bearing surfaces; run the check-only formatting gate recorded in `openspec/config.yaml`.

## Deferred Decisions

Module/publication identity, packaging, CI, license, release automation, distribution, and concrete configuration-file syntax remain deferred. CNSIC W1/W2 at merge `a0fc7b26ff8a0e0a61baa586b32c46841611c806` remains the immutable predecessor oracle; do not copy CNSIC-specific rules into the neutral core.
