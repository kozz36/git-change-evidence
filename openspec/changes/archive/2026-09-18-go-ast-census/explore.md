# Go AST selector-call census exploration

## Outcome

Plan a future, public `changeevidence/census` subpackage that performs an exact textual Go AST census of selector-call candidates across a fully supplied raw path/content set. It is a locator only: it does not resolve imports, aliases, shadowing, types, dependencies, or semantics. Semantic completeness is excluded; exact completeness for the supplied full-file scope and syntactic query is required.

The API must accept an `InventoryDocumentV1`, supplied file bytes, an exact versioned query, and explicit positive limits. It must validate the complete supplied set against the inventory before parsing any file and return an error with no partial result for invalid inventory, missing, extra, duplicate, digest/length drift, parsing failure, or limit breach. A valid empty inventory produces an empty result; a zero-value inventory is invalid.

## Repository evidence

| Evidence | Location | Observation |
|---|---|---|
| Inventory API | `inventory_v1.go:5-36` | `InventoryEntryV1` contains `Path []byte`, `ContentSHA256 Digest`, and `ByteLength uint64`; `InventoryDocumentV1` exposes defensive `CanonicalBytes`, `Digest`, `AccountingPolicyDigest`, and `Entries` copies. Inventory itself has no content bytes. |
| Inventory construction | `inventory_v1_build.go:23-95` | Inventory entries are sorted by raw byte path, duplicate paths are rejected, source bytes are SHA-256 hashed, and an empty explicit entry slice is a valid canonical document. |
| Inventory decoding | `inventory_v1_decode.go:11-51` | Decoding validates canonical form, paths, entry ordering, and duplicate paths, but does not supply file content. |
| Ownership and empty behavior | `inventory_v1_ownership_test.go:10-44` | Existing tests distinguish valid empty from zero/failed inventory and assert that entry paths and canonical bytes cannot escape mutable state. |
| Digest/error conventions | `contract.go:14-16,48-50`; `policy_contract.go:19-21` | `Digest` and `DocumentDigest` are distinct string types; contract failures use `*ContractError` with stable field/code values. |
| Layout guidance | `AGENTS.md:18-30`; `openspec/config.yaml:35-42,70-84` | Public neutral APIs are currently directed to the module root, co-located tests are required, and Go changes require strict TDD plus `go test ./...` and check-only gofmt. |

## Proposed future contract

### Inputs

Define public, owned-data types under `census` for:

- A raw file record: raw `Path []byte` plus full `Content []byte`.
- A versioned selector-call query, with exact textual receiver and selector identifiers. The query is an exact candidate predicate, not an import or type identity.
- Explicit limits: `MaxFiles`, `MaxFileBytes`, `MaxTotalBytes`, and `MaxMatches`.

The constructor/entry point should receive `changeevidence.InventoryDocumentV1`, the raw file records, query, and limits, then return a result or error. The approved default limits are exactly 256 files, 1 MiB per file, 8 MiB total input, and 10,000 matches. Any configurable limit must be positive; no parser-memory isolation guarantee is made beyond these input/count bounds.

Validate the inventory document as an antecedent using its owned public accessor state: distinguish the zero value from a constructed document and verify the canonical-byte SHA-256 identity against its digest. Do not introduce canonical round-trip decoding here: `DecodeInventoryV1(raw, policy)` requires a `PolicyDocumentV1`, which is not an approved census input. The document fields are unexported, so normal external Go consumers obtain nonzero instances through existing validated constructors/decoders; census must not duplicate that decoder or add a policy input. Then validate both directions of the path set using byte-wise raw paths: each inventory entry has exactly one supplied record and each supplied record has exactly one inventory entry. For every pair require `sha256(content) == ContentSHA256` and `uint64(len(content)) == ByteLength`. Reject nil/zero-value or otherwise invalid inventory, duplicate supplied paths, duplicate inventory paths, absent inventory paths, supplied extras, and any drift before AST parsing or match publication.

### Query and matching rules

A query must have a versioned exact identity, retained verbatim in the result. It matches only an AST `*ast.CallExpr` whose function is an `*ast.SelectorExpr` after repeatedly unwrapping only these syntactic wrappers around the callee expression:

1. Parentheses (`*ast.ParenExpr`).
2. Index forms (`*ast.IndexExpr` and `*ast.IndexListExpr`). An `IndexExpr` is syntactically ambiguous: it can represent type instantiation or value indexing. Unwrapping may therefore include `obj.Funcs[i]()` as a candidate for `obj.Funcs`; no type-based disambiguation or silent normalization narrowing is allowed.

After normalization, the selector expression must have an identifier receiver whose spelling equals the query receiver and whose selector spelling equals the query selector. This admits text such as `os.Open`, `(os.Open)`, `os.Open[T]`, and `(os.Open[T])`; it rejects calls requiring semantic interpretation. Thus a locally shadowed `os.Open` may appear, while `alias.Open` does not match a query for `os.Open`; those facts are explicitly candidate limitations, not errors.

Use `go/parser`/`go/ast` on each complete supplied file. Comments and string literals are not AST selector-call candidates. Parse errors abort the operation without a partial result. Traverse only the parsed AST, record the outer `CallExpr` and its normalized underlying selector, and do not use textual scanning to create candidates.

### Result and binding

The result must retain:

- the inventory document digest;
- the exact, versioned query;
- an implementation-owned extractor-version constant; and
- matches sorted by raw path bytes, then `CallExpr` start byte offset, then a deterministic final position tie-breaker.

Each match retains defensive copies of the raw path, file SHA-256, and a SHA-256 digest of the exact full-call source fragment. It also retains two half-open byte spans: one for the outer `CallExpr`, and one for its normalized selector expression. A span is `[start,end)` over the supplied content and reports physical, 1-based byte line and column positions. Position computation must ignore `//line` remapping and therefore use physical parser file positions/offsets, not adjusted position APIs. Byte offsets and byte columns must remain correct for UTF-8, CRLF, and multiline source.

The full-call fragment digest is `sha256(content[callStart:callEnd])`; it is not a digest of a canonical/global result serialization. All exported byte slices and match collections require defensive ownership on input and output, including result accessors.

### Error and atomicity policy

Use typed, stable validation errors consistent with the repository’s `ContractError` pattern where public contract fields/codes are useful. Do not expose a result on any error; errors from validation, parsing, or limits are atomic failures. Validate file and total byte limits before parse; enforce match limit while collecting and discard any accumulator on breach. The acceptance surface should assert that equal match counts cannot hide incorrect locations/ranges and that returned data cannot mutate retained state.

## Future implementation surfaces

| Surface | Purpose |
|---|---|
| `census/query.go` | Public versioned query, limits/defaults, input/result/error types, validation boundary, defensive accessors, and extractor version ownership. |
| `census/go_ast.go` | Internal AST parsing, bounded wrapper normalization, physical byte-span calculation, candidate extraction, and deterministic sorting. |
| `census/go_ast_test.go` | Unit/triangulation coverage for candidate behavior, validation, limits, positions, sorting, and fragment binding. |
| `census/external_api_test.go` | External-package API and defensive-ownership coverage. |
| `docs/census.md` | Consumer-facing locator scope, query syntax/identity, input binding, result semantics, limits, and explicit non-goals. |
| `AGENTS.md` | A narrow, future owner-approved exception permitting this public `census` subpackage while retaining root-only public neutral APIs otherwise. |
| `openspec/config.yaml` | Future owner-approved alignment of `layout.public_api` only with the same narrow exception; testing policy/rubric remains untouched. |

No store, DSL, CLI, serialization format, adapter, mutation, runner, receipt, Python pilot, U10/U11 work, or dependency-resolution facility is in scope.

## Acceptance matrix

| Area | Required future examples |
|---|---|
| Candidate behavior | Repeated and legacy sites; comment/string decoys; `os.Open` shadowing admitted as a candidate; aliases excluded; one callee spelling associated with different downstream invariants without inferred classification. |
| AST normalization | Parenthesized callees and single/multi generic index forms normalize only around the callee; non-selector expressions do not match. |
| Source coordinates | Unicode, CRLF, multiline calls, and `//line` directives yield exact physical byte offsets and 1-based byte line/columns; both selector and call spans are half-open and exact. |
| Input binding | Missing, extra, duplicate, content-hash drift, byte-length drift, invalid/zero inventory, and malformed Go all error with no result; empty valid inventory returns zero matches. |
| Bounds | Defaults are exact; every explicit limit is positive; file, total, and match limits fail atomically. |
| Determinism and ownership | Raw-byte path then position ordering is stable; ranges remain exact even when counts are equal; per-fragment hashes bind source; external consumers cannot mutate inputs or returned bytes/slices. |

## Delivery sizing and risk

The production ceiling is 400 novel lines across `census/query.go` and `census/go_ast.go`; this is not a native gross review budget. The requested tests are independently estimated at 450–700 lines across the two test files, driven by the acceptance matrix. The proposed vertical slice therefore has a review-budget risk under the session’s 400-line canonical threshold, even if production stays within its stated ceiling. Because delivery strategy is `ask-on-risk`, future planning/apply must stop and request a delivery decision if the acceptance suite cannot be organized within the selected review shape; it must not silently reduce acceptance or infer `size:exception`.

## Configuration/layout issue

The approved public `census` subpackage conflicts with the current root-only public API language in both `AGENTS.md` and `openspec/config.yaml` (`layout.public_api.directory: .`). The approved product calls for a deliberately narrow exception, but the configuration has no documented mechanism for a package-specific public-subpackage exception. The user confirmed a future task aligning both `AGENTS.md` and `openspec/config.yaml` `layout.public_api` only, without testing-policy/rubric changes. Neither surface may be edited during formalization.

## Out of scope

This exploration does not implement code or tests, execute runtime commands, alter lifecycle markers, create a commit, push, or modify `.atl/skill-registry.md` or `.gitignore`. It does not depend on the sibling `cnsic-profile-parity` change or V2 work.
