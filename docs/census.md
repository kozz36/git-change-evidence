# Go AST census CLI

`git-change-evidence census-go-ast` is a local, evidence-only adapter over the
public `census.GoASTV1` API. It neither discovers files nor invokes Git.

## Command

```text
git-change-evidence census-go-ast \
  --source-root <absolute-source-root> \
  --receiver <identifier> \
  --selector <identifier> \
  -- <explicit-relative-go-files...>
```

The `--` separator and at least one selected path are required. `--source-root`,
`--receiver`, and `--selector` are each required exactly once. Identifiers use
Go's `token.IsIdentifier` validation. The command has no implicit current
working directory root and does not accept Git revisions, shell expressions, or
discovery patterns.

For example, after creating `/opt/census-fixture/example.go` with an `os.Open`
call, a local smoke command is:

```sh
git-change-evidence census-go-ast \
  --source-root /opt/census-fixture \
  --receiver os --selector Open -- example.go
```

## Input confinement and limits

The source root must be an absolute, lexical-canonical path. Selected paths must
be unique, lexical-canonical, relative `.go` paths. Absolute paths, empty paths,
`.` or `..` components, repeated separators, traversal, duplicates, and
non-Go paths are rejected before source-root access.

On Linux, the adapter opens `/` and every supplied root/path component through
file descriptors using `openat`, `O_NOFOLLOW`, `O_NONBLOCK`, and post-open
`fstat`/`fstatat(AT_SYMLINK_NOFOLLOW)` identity checks. It rejects symlinks at
the root, intermediate, and selected-file positions; rejects nonregular files;
and reads only the opened regular file descriptor. This prevents path-component
symlink traversal and avoids FIFO/device open hangs, including replacement races
between metadata checks and open.

The Linux implementation is deliberate: another platform reports source content
as unavailable rather than weakening the no-follow confinement guarantee.

`census.DefaultLimits()` is applied before parsing source content:

- `max_files`: 256
- `max_file_bytes`: 1 MiB
- `max_total_bytes`: 8 MiB
- `max_matches`: 10,000

Every selected regular file receives one bounded content read. The CLI builds
`InventoryDocumentV1` from those exact captured bytes and passes the same bytes
to `census.GoASTV1`.

This is a per-file captured-byte snapshot, not an atomic snapshot of the whole
source tree. A concurrently changing tree can therefore contain selected files
captured at different instants; each emitted file hash, byte length, inventory
identity, and match span still bind to the exact bytes supplied to the extractor.

## JSON output

A successful invocation writes exactly one JSON document and a newline. Output
is deterministic for the same selected byte content and query; input path order
does not affect it. Errors exit nonzero, report only a generic diagnostic on
stderr, and do not intentionally emit a success document.

The versioned envelope is `git-change-evidence.census-go-ast-cli/v1`:

```json
{
  "schema": "git-change-evidence.census-go-ast-cli/v1",
  "extractor_version": "git-change-evidence.go-ast-census/v1",
  "query": {
    "version": "git-change-evidence.selector-call-query/v1",
    "receiver": "os",
    "selector": "Open"
  },
  "limits": {
    "max_files": 256,
    "max_file_bytes": 1048576,
    "max_total_bytes": 8388608,
    "max_matches": 10000
  },
  "inventory": {
    "schema": "git-change-evidence.inventory/v1",
    "sha256": "<inventory-sha256>",
    "accounting_policy_sha256": "<fixed-census-policy-sha256>"
  },
  "source_scope": [
    {
      "path_b64": "ZXhhbXBsZS5nbw==",
      "content_sha256": "<file-sha256>",
      "byte_length": 42
    }
  ],
  "matches": [
    {
      "path_b64": "ZXhhbXBsZS5nbw==",
      "file_sha256": "<file-sha256>",
      "fragment_sha256": "<call-fragment-sha256>",
      "call": {
        "start": {"offset": 25, "line": 3, "column": 12},
        "end": {"offset": 37, "line": 3, "column": 24}
      },
      "selector": {
        "start": {"offset": 25, "line": 3, "column": 12},
        "end": {"offset": 32, "line": 3, "column": 19}
      }
    }
  ]
}
```

`source_scope` is complete, including selected files with zero matches. The
inventory uses a fixed neutral census policy only because the public
`NewInventoryV1` constructor requires a policy antecedent; the CLI performs no
accounting classification.

`path_b64` is standard Base64 of the raw selected path bytes. It is the wire
representation for every source-scope and match path, so paths with non-UTF-8
bytes round-trip losslessly without placing invalid UTF-8 in JSON strings. The
source root is deliberately absent from the output. Query identifiers are JSON
strings because invalid Go identifiers are rejected before output.

The `matches` array is empty (`[]`) when no candidate is found. Positions are
physical source-byte offsets plus one-based line and byte column values from the
published census API.

The document is first encoded in memory and then attempted as one write. If a
failing downstream writer retains a prefix of that attempted write, the process
still returns nonzero; a streaming CLI cannot retract bytes already accepted by
such a writer.

## Exit behavior

| Condition | Exit | stderr |
| --- | ---: | --- |
| malformed command, invalid paths, or parse failure | 2 | `invalid input` |
| unavailable, missing, symlinked, or nonregular source | 3 | `content unavailable` |
| a configured file, byte, or match limit is exceeded | 4 | `resource bound exceeded` |
| successful zero-match or nonzero-match census | 0 | empty |
