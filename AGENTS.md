# Agent Guide

## Authority

`git-change-evidence` is evidence-only. It may report technical states, measurements, and observations, but it must never approve, reject, gate, block, merge, deploy, release, or assign delivery authority.

## Current public-preview governance

The repository contains its Go module, implementation, tests, and CI. The planned source/dev preview identity is `v0.1.0-preview.1`; it is not a published tag, release, package, asset, distribution channel, or visibility change. The repository license selection is Apache-2.0, without a claim that this settles third-party attribution, intellectual-property, confidentiality, export, or other legal questions.

The CLI and its machine-readable contracts are the primary consumer surface. The canonical pre-v1 census command is `gce census go-ast`; `git-change-evidence census-go-ast` remains supported throughout pre-v1 with the established contract meaning and technical behavior. The public root Go API and public `census` subpackage also remain supported.

Multilingual census is a separate post-preview roadmap, not a preview blocker. Its language sequence is Python → JavaScript/TypeScript → Vue → Java → Rust. Do not implement that behavior in the preview work.

Historical bootstrap/OpenSpec/CNSIC/local-path records, literal fixtures, and the immutable predecessor reference are evidence, not current instructions. Preserve them without silent rewriting or deletion.

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
- Target Go baseline: **Go 1.25.10**. The established module identity is `github.com/kozz36/git-change-evidence`; its `go.mod`, Go source files, and co-located Go tests are present in this repository.
- For Go implementation work, follow RED → GREEN → TRIANGULATE → REFACTOR. The primary gate is `go test ./...`; use `go test -race ./...` for concurrency-bearing surfaces; run the check-only formatting gate recorded in `openspec/config.yaml`.

## Deferred Decisions

Concrete configuration-file syntax, broad distribution, packaging, signing/provenance, and release automation remain deferred. CNSIC W1/W2 at merge `a0fc7b26ff8a0e0a61baa586b32c46841611c806` remains the immutable predecessor oracle; do not copy CNSIC-specific rules into the neutral core.

## Repo-local skills

- `gce-evidence`: For explicit requests to consume, validate, reproduce, regenerate, or technically interpret GCE evidence, read `skills/gce-evidence/SKILL.md`. Technical evidence only; no delivery authority. Repository-local guidance, not automatic installation.
