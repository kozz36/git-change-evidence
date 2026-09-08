# Result V1 external validation

## Scope and boundary

This is supplemental test evidence for delivered Result V1 behavior, not a production
change or a reopening of its archived delivery. The root test is `package
changeevidence_test`: it imports only `github.com/kozz36/git-change-evidence` and
uses exported typed constructors/accessors. The Git acquisition composition test is
in the existing internal adapter package because `Acquire` is intentionally not a
public API.

This does not demonstrate an independently released module, a public Git acquirer,
or a CLI feature. It makes no performance or resource-bound claim: `Acquire` uses its
existing unbounded byte-cap path, and the 256-line text fixture is only a bounded test
fixture. A SHA-256 Git-object composition vector is deferred; adapter-specific
SHA-256 acquisition coverage already exists separately.

## Oracle matrix

| Case | Input boundary | Literal oracle |
| --- | --- | --- |
| Public typed consumer | Two separately constructed Policy → Inventory → committed Snapshot → Result → Report chains | Equal result/report canonical bytes and digests; revisions are `0123456789abcdef0123456789abcdef01234567` to `89abcdef0123456789abcdef0123456789abcdef`; entries are `assets/logo.bin` non-countable `other`, `docs/guide.md` `docs` +3/-1, and `src/new.go` `source` +4/-0; totals are docs +3/-1, source +4/-0, other non-countable 1; threshold and ratio observations are asserted literally. |
| Public negative boundary | A nonempty valid `other` policy and its inventory are deliberately different from the valid chain | Result construction rejects the wrong antecedent and returns zero Result; result/report binding rejects and returns zero Evidence; empty policy returns `categories/required`, and `../invalid` inventory path returns `entries.path/invalid_path`. |
| Committed adapter composition | Real temporary Git repository: exact rename `docs/old.txt` to `src/new.go`, binary NUL change, raw `raw/0xff+LF+path` change, and 256-line text replacement | Snapshot has exactly four entries. Rename is `R`, old/new paths are literal, and is +0/-0 countable. Binary is `M` and non-countable. Raw path is +1/-1; bulk text is +256/-256. Result categories/totals are asserted independently: docs 0, source 0, binary non-countable 1, raw +1/-1, bulk +256/-256, other 0. |
| Provenance and committed-only evidence | The same acquired snapshot is composed into Policy/Inventory/Result/Report; then index, worktree, and untracked content are dirtied without a commit | Strict report-result provenance validation, revision identity, and all digest links succeed. Reacquisition and recomposition have equal canonical result/report bytes and digests. `Base == Head` acquisition returns zero entries. |

The raw path fixture is asserted as bytes and is skipped on Windows, where the needed
Unix filename semantics are unavailable. It passed on the observed Unix test host.

## Commands and observed results

Before adding these tests, the focused baseline was:

```text
go test . ./internal/git -count=1 -v
```

It passed for the root and internal Git packages. After the additions, the two focused
commands each selected one new top-level test (`TestExternalPublicResultV1ChainIsReproducibleAndBound`, with two negative subtests; and
`TestAcquireComposesCommittedResultAndReportV1`) and passed:

```text
go test . -run '^TestExternalPublicResultV1ChainIsReproducibleAndBound$' -count=1 -v
go test ./internal/git -run '^TestAcquireComposesCommittedResultAndReportV1$' -count=1 -v
```

The broader observed commands also passed for the root, CLI, internal Git, internal
inventory, and internal publication packages:

```text
go test ./... -count=1
go test ./... -short -count=1
```

The integration test intentionally reports `SKIP: uses a real Git repository` under
`-short`; it was observed with:

```text
go test ./internal/git -short -run '^TestAcquireComposesCommittedResultAndReportV1$' -count=1 -v
```

The configured PR workflow uses Go 1.25.10 and runs `go test ./...`, so these normal
co-located tests are included in that future CI command. Its separate race command is
not newly required: this unit adds no concurrent behavior.

## Layout and size census

The post-change physical census is 69 root Go files (68 before this new co-located
external test), and the internal leaf packages contain `git=5`, `inventory=5`, and
`publication=4` Go files. The configured public-root threshold is 12 and therefore
predates this test; the suggested root count of 84 was not observed in this checkout.
The external test remains at the root to preserve the actual public API boundary,
while the adapter integration test remains co-located with the internal adapter.

`external_api_test.go` is 109 lines and
`internal/git/snapshot_result_integration_test.go` is 121 lines, both below the
160-line configured per-Go-file limit. This unit adds zero production lines against
the 400-production-LOC budget; test and documentation lines are separate. No
production, pre-existing test, configuration, or archive file was changed.

## Key Learnings

1. A public typed consumer can reproduce canonical Result and Report documents without internal imports.
2. Result/report provenance rejects a valid document chain when its antecedents do not match.
3. The Git adapter composes committed rename, binary, hostile-path, and bounded text evidence without observing dirty state.
4. Raw non-UTF-8 newline filenames require Unix filesystem semantics and are deliberately skipped on Windows.
5. The repository PR workflow already executes `go test ./...` with Go 1.25.10.
