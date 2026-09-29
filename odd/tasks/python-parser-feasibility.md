# Probe pure-Go Python parser feasibility (issue #143)

Status: authorized test-only feasibility work, not census implementation or dependency adoption.
Branch: `feat/python-parser-feasibility`; base: `7c096f6f58e83bb3135afb7e98569ebe4c69b757`.
Recovery mirror: Engram `odd/python-parser-feasibility/tasks` in the CNSIC-bound orchestration session (#19719); a GCE-project mirror was refused by session binding.

## Scope and route

- F1 — Pin one Gopapy candidate revision and build an isolated probe with `.py`/`.pyi`, exact call-shape, byte-span and atomic-error fixtures. Route: delegated writer (multi-file trigger); native strict TDD RED → GREEN → TRIANGULATE → REFACTOR. Characterization complete; span and Go-baseline axes fail. Commit: `4f63a3f`; see `docs/python-parser-feasibility.md`.
- F2 — Execute bounded resource, supply-chain and standalone `CGO_ENABLED=0` checks; publish reproducible pass/fail/unknown evidence. Route: delegated writer plus focused verification (external and multi-file triggers). Characterization complete; 10,000-match timeout and unproved hard resource controls recorded. Commit: recorded in handoff; see `docs/python-parser-feasibility.md`.
- F3 — Independent diff/evidence review, PR and owner-authorized merge of the evidence work unit only. Route: read-only reviewer and parent delivery. Pending.

Forecast: 350–500 authored changed lines, subject to measurement before delivery; avoid broadening into the production census. One work-unit commit per completed implementation task, with focused tests and documentation beside behavior. Do not check off an axis without observed evidence. A failed parser axis remains a reported failure, not a loosened requirement.

## Verification source

GCE `AGENTS.md` and `openspec/config.yaml`: strict TDD for Go source/tests; primary `go test ./...`; check-only format gate `test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"`; `go test -race ./...` only if concurrency-bearing. The fixed `.pi/gentle-ai/project-policy.yaml` is absent, so this is the verified native route, not an active optional rubric.

## Hard exclusions

No public API, product CLI, production parser dependency selection, release, CNSIC pilot/integration, semantic-resolution claim, or delivery authority. Parser feasibility is not established by the existence of this task file.
