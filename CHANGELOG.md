# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [0.1.0-preview.1] - Unreleased

> **Planned and unpublished.** `0.1.0-preview.1` is a planned source/dev preview
> identity. It does not claim a tag, release, package, asset, distribution
> channel, or public-visibility change.

### Added

- The canonical CLI and machine-readable contract surface, with
  `gce census go-ast` as the canonical pre-v1 census command.
- Retained compatibility for `git-change-evidence census-go-ast` throughout
  pre-v1, with the established equivalent contract meaning and technical
  behavior.
- Supported public Go APIs at the module root and in the public `census`
  subpackage.
- Apache-2.0 repository licensing, root third-party notices for
  `golang.org/x/sys v0.38.0`, and public-preview documentation and governance
  guidance.
- Release-preparation evidence that `go-ast-census` is archived, while
  `bootstrap-neutral-core-extraction` remains active and unarchived at 3/12
  tasks and `post-preview-neutral-capabilities` remains an active, proposal-only
  change with zero tasks and next recommended step `spec`.

### Limitations

- The Go AST census is syntax-only over explicitly supplied Go file bytes; it
  does not perform semantic resolution or provide an atomic whole-tree snapshot.
- Multilingual census remains a post-preview roadmap item, not a preview
  capability or blocker.
- `git-change-evidence` is evidence-only: it reports technical states,
  measurements, and observations and does not approve, reject, gate, block,
  merge, deploy, release, or assign delivery authority.
- Broad distribution, release automation, signing, and provenance remain
  deferred. Apache-2.0 does not settle third-party attribution,
  intellectual-property, confidentiality, export, or other owner-controlled
  questions.
