# Python syntax census — design proposal (#141)

Status: planning only. This is not an OpenSpec specification, parser selection, implementation commitment, or delivery authority. The active post-preview neutral-capabilities work excludes multilingual census; this separate issue must be decided on its own merits. The existing [Go census contract](census.md), [`census.QueryV1`](../census/query.go), and [`GoASTV1`](../census/go_ast.go) are the local baseline, not evidence that a Python parser meets it.

## Boundary and predicate

**OWNER-DECIDED:** For a query with exact receiver `os` and selector `open`, report precisely an AST `Call` whose `func` is `Attribute(value=Name(id="os"), attr="open")`. This is syntax-only: neither imports nor bindings are resolved. Select `.py` and `.pyi` paths explicitly; reject invalid or unsupported syntax atomically. The executable must be a standalone Go binary requiring neither installed Python nor cgo.

| Source expression (query `os.open`) | Candidate? | Reason |
| --- | --- | --- |
| `os.open("x")` | Yes | Direct `Name.Attribute` call. |
| `os.open` | No | Attribute access, not a call. |
| `os.path.open("x")` | No | Attribute receiver is another attribute. |
| `os.open[0]("x")` | No | Call function is an index expression. |
| `os.open("x").read()` | Yes, inner call only | Outer call has a different function shape. |
| `os = fake; os.open("x")` | Yes | Shadowing does not suppress a textual candidate. |

Parenthesized, decorated, async, and other grammar forms require explicit fixtures before freezing extractor behavior; do not silently widen the exact function-node predicate. [Python's `ast` reference](https://docs.python.org/3/library/ast.html) describes `Call`, `Attribute`, and `Name` nodes, but CPython behavior is reference material, not proof of a Go backend.

## Evidence contract to preserve

**PROPOSED invariants reused from Go:** Accept a complete, explicit inventory of raw paths, content hashes, and byte lengths, including zero-match files; never discover files or invoke Git. The API accepts typed path/content bytes and validates exact inventory-to-content correspondence before interpreting syntax. Only a CLI adapter accesses a caller-supplied confined filesystem root and captures bounded bytes. An inventory is not a whole-tree snapshot. Provide deterministic, versioned query/extractor/output identity, raw-path-safe wire identity, ordered matches bound to file and fragment hashes, and validated half-open physical source-byte spans. Reject mismatches, duplicate/invalid paths, limits exceeded, invalid/unsupported syntax, and unavailable inputs without a partial success result. The CLI should serialize a complete result only after success and report failure nonzero; a downstream writer that retains a partial write cannot be rolled back. Candidates do not establish semantic identity or authorize gates, merge, or release. See [Go CLI scope, errors, and limits](census.md) and [typed bytes and positions](../census/query.go).

**PROPOSED, provisional:** Define a Python-specific versioned envelope and query version rather than silently reusing the Go extractor version. Bind every match to raw path bytes (for example Base64 on JSON), exact file hash, fragment hash, and call/function spans; sort by raw path bytes and offsets independent of input order. Reuse validated inventory machinery where compatible. The schema name, CLI spelling, limits/defaults, grammar target/version, and line/column unit remain provisional. Physical source **byte offsets** are proposed as required evidence, but backend proof of offsets is pending; do not infer byte columns from Unicode character columns. Existing Go defaults bound selected file counts/bytes and match count, not serialized output bytes or parser heap, stack, or time.

## FEASIBILITY-BLOCKED: parser proof before dependency choice

No dependency is selected. [Gopapy's upstream README](https://github.com/tamnd/gopapy#readme) claims pure-Go parsing and Python grammar coverage; it does not establish the byte spans, `.pyi` fixtures, atomic errors, or resource behavior required here. Those need independent tests before adoption. The [Python AST reference](https://docs.python.org/3/library/ast.html) supplies node vocabulary, not guarantees about a Go backend.

| Proof axis | Required independent demonstration on a pinned candidate |
| --- | --- |
| `.py` / `.pyi` grammar | Parse representative supported-version scripts and stubs (annotations, overloads, newer constructs); name grammar target and reject unsupported forms instead of silently producing partial matches. |
| Unicode and physical spans | Compare offsets against the original byte buffer for multibyte identifiers/text, mixed line endings, non-ASCII encodings if supported, and nested calls; verify half-open slices reproduce exact fragments and hashes. |
| Invalid syntax / atomicity | Malformed and unsupported files mixed with valid files must fail the whole API result; error recovery nodes cannot yield success candidates. |
| Large / deep inputs | Measure worst-case depth, broad trees, huge tokens, file counts and match counts under limits; establish cancellation, time, heap/stack controls or explicit isolation, not merely byte-size bounds. |
| Supply chain / standalone build | Pin version and grammar artifacts; independently inspect license, transitive dependencies, and reproducible `CGO_ENABLED=0` build/run with no installed Python or dynamic parser requirement. |

**FEASIBILITY-BLOCKED exit:** If no candidate passes every applicable proof axis, stop and return the failed evidence and alternatives to the owner for a new scope decision. Do not relax atomicity, byte evidence, syntax coverage, or standalone constraints implicitly. No implementation, dependency change, release, or CNSIC pilot follows from this document.
