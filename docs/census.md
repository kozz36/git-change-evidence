# Go AST selector-call census

`census.GoASTV1` finds syntax-only textual selector-call candidates in one complete, validated `InventoryDocumentV1` scope. It reports technical source evidence only; it does not approve, reject, gate, block, merge, deploy, release, or assign delivery authority.

## Use the supplied scope

Construct the inventory with the root API, then provide exactly one complete byte record for every inventory path:

```go
result, err := census.GoASTV1(inventory, []census.File{
    {Path: []byte("source.go"), Content: source},
}, census.QueryV1{
    Version: census.SelectorCallQueryV1,
    Receiver: "os",
    Selector: "Open",
}, census.DefaultLimits())
```

The operation validates the inventory's canonical digest through its public accessors, raw path membership in both directions, content byte lengths, and content SHA-256 digests before parsing. Duplicate supplied paths, missing records, extra records, drift, invalid queries or limits, parser failures, invalid spans, and match-limit failures return no result. A constructor-created empty inventory and empty supplied set succeeds with an empty, non-nil match collection.

## Candidate predicate

A candidate is a Go `CallExpr` whose callee becomes a selector after unwrapping only callee `ParenExpr`, `IndexExpr`, and `IndexListExpr` nodes. The selector receiver must be an identifier exactly spelling `QueryV1.Receiver`, and its selector must exactly spell `QueryV1.Selector`.

For query `os.Open`, `os.Open(...)`, `(os.Open)(...)`, and `os.Open[T](...)` are candidates. A shadowed local spelling remains a candidate. The intentionally ambiguous `obj.Funcs[i]()` form remains a candidate for `obj.Funcs`, because the syntax alone cannot distinguish value indexing from instantiation. Aliases with another spelling, receiver wrappers such as `(os).Open()`, comments, strings, and non-selector calls are not candidates.

The census parses every complete supplied file with `go/parser` and traverses its AST. It does not scan text, load files, resolve imports, types, aliases, shadows, dependencies, or downstream invariants. Results therefore claim exact completeness only for the validated supplied inventory and exact textual query, never semantic package or method identity.

## Evidence and ownership

Each ordered match retains a defensively owned raw path, whole-file SHA-256, SHA-256 of the exact full call fragment, and half-open call and normalized-selector spans. Span endpoints contain zero-based byte offsets and physical one-based line and byte-column values. UTF-8, CRLF, multiline calls, and `//line` directives retain supplied physical byte coordinates; directives do not remap reported locations.

Matches sort by raw path bytes, call start offset, call end offset, selector start offset, then selector end offset. `ResultV1.Matches` returns a new collection and new path bytes every time. Callers may mutate input or returned slices after the operation without changing retained evidence; callers must still avoid concurrent argument mutation during the call.

## Bounds

`DefaultLimits` sets 256 files, 1 MiB per file, 8 MiB total input, and 10,000 matches. Every explicit limit must be positive. File-count and byte limits are checked before parsing; match limits are checked while collecting; equality succeeds.

These limits bound admitted input bytes and published match count only. They do not guarantee bounds for parser heap use, stack use, CPU time, panic recovery, or process isolation.
