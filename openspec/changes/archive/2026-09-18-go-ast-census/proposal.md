# Propose an exact textual Go AST selector-call census

## Decision summary

Plan a future public `changeevidence/census` subpackage that locates every syntactic selector-call candidate matching an exact textual query across a complete, caller-supplied set of Go files bound to an existing `InventoryDocumentV1`.

This is a neutral, evidence-only census. It guarantees completeness only for the exact supplied inventory scope and exact syntactic query. It does not resolve imports, types, aliases, shadowing, dependencies, or semantic identity, and it has no authority to approve, reject, gate, block, merge, deploy, or release anything.

This change is planning-only. It authorizes this proposal, not source implementation, tests, commands, delivery, or edits to repository policy.

## Intent

Provide consumers with reproducible source-location evidence for questions such as “where does the supplied Go source contain a call whose textual receiver and selector are exactly `os.Open`?” The future census will bind each answer to the supplied inventory, query, extractor version, file content, exact call fragment, and physical byte ranges so downstream users can distinguish exact locations rather than relying on counts or semantic assumptions.

## Users and situations

The intended users are evidence collectors and maintainers who need a bounded, deterministic inventory of candidate call sites before applying project-specific interpretation elsewhere. They need to:

- inspect repeated, legacy, or differently used textual call sites across an explicitly supplied file set;
- retain exact source evidence even when equal counts conceal changed locations;
- distinguish syntactic candidate discovery from import, type, or policy classification; and
- reject incomplete or drifted inputs rather than publish a partial census.

## Current-state gap

The repository has a canonical `InventoryDocumentV1`, including raw paths, file digests, and byte lengths, but no public facility that accepts the corresponding full file bytes and produces an exact AST-based selector-call candidate census. Consumers would otherwise need ad hoc scanning, could accidentally treat aliases or shadowing as resolved identities, and might publish results from a partial or mismatched inventory.

## Proposed scope

### Public boundary

A future public `changeevidence/census` subpackage will define owned-data APIs for:

- a supplied file record containing raw `Path []byte` and complete `Content []byte`;
- an exact, versioned selector-call query containing textual receiver and selector identifiers;
- positive limits `MaxFiles`, `MaxFileBytes`, `MaxTotalBytes`, and `MaxMatches`;
- a census result, match records, physical spans, and stable typed validation errors where appropriate; and
- an implementation-owned extractor-version constant.

The census entry point will accept:

1. an existing `changeevidence.InventoryDocumentV1`;
2. the complete caller-supplied raw path/content file set;
3. the exact versioned query; and
4. explicit positive limits, with an API-provided exact default set.

It will not accept project policy. It will not decode inventory serialization, duplicate `InventoryDocumentV1` decoding, or add `PolicyDocumentV1` as an input.

### Inventory and supplied-scope validation

Before parsing any file, the future implementation must:

- reject a zero-value inventory antecedent (the public input is an `InventoryDocumentV1` value; nil is not representable);
- validate the inventory through its owned public accessor state, including canonical-byte SHA-256 identity against its document digest;
- reject duplicate supplied paths; retain Inventory V1 constructor/decoder ownership of duplicate inventory rejection, without fabricating invalid private document state;
- require every inventory path to have exactly one supplied file and every supplied file to have exactly one inventory entry;
- compare paths byte-for-byte as raw paths;
- require `sha256(content)` to equal the inventory entry’s `ContentSHA256`;
- require `uint64(len(content))` to equal the inventory entry’s `ByteLength`; and
- enforce file-count, per-file-byte, and aggregate-byte limits before parsing.

A constructed valid empty inventory paired with an empty supplied set succeeds with an empty result. A zero-value inventory is invalid.

These rules make completeness exact for the full supplied inventory scope. The API is not hunk-based and must not accept a partial changed-line view as if it were a complete census.

### Exact query and candidate predicate

The query identity is exact and versioned, is retained verbatim in the result, and names a textual receiver identifier and selector identifier.

A candidate is an `*ast.CallExpr` whose callee becomes an `*ast.SelectorExpr` after repeatedly unwrapping only:

1. `*ast.ParenExpr`;
2. `*ast.IndexExpr`; and
3. `*ast.IndexListExpr`.

The normalized selector must have an identifier receiver whose spelling exactly equals the query receiver, and its selector identifier must exactly equal the query selector.

Examples admitted by a query for `os.Open` include `os.Open(...)`, `(os.Open)(...)`, `os.Open[T](...)`, and `(os.Open[T])(...)`. Because a single `IndexExpr` is syntactically ambiguous between instantiation and value indexing, normalization intentionally may admit an indexed-value call such as `obj.Funcs[i]()` as a candidate for textual selector `obj.Funcs`. This false-positive class is required and must not be silently narrowed using semantic reasoning.

A locally shadowed `os.Open` remains a candidate. `alias.Open` does not match `os.Open`. Non-selector callees do not match. Comments and string literals are not candidates.

The parser must consume each complete supplied file and candidates must come only from `go/parser` and `go/ast` traversal, not textual scanning. Object resolution is skipped; no import loading or type checking is performed.

### Exact source evidence

Each match must retain owned copies or immutable values for:

- the raw path;
- the supplied file’s SHA-256 digest;
- the SHA-256 digest of the exact full-call fragment `content[callStart:callEnd]`;
- the outer `CallExpr` half-open byte span `[start,end)`; and
- the normalized underlying `SelectorExpr` half-open byte span `[start,end)`.

Both spans are over the exact supplied content. Each endpoint reports a zero-based byte offset and a physical 1-based line and byte-column position. Position calculation must ignore `//line` remapping and use the owning parser file’s unadjusted physical positions. Endpoint ownership and bounds must be checked rather than relying on clamping behavior.

Offsets and byte columns must remain exact for UTF-8, CRLF, multiline calls, and files containing line directives.

### Result binding and ordering

The result must bind exactly to:

- the `InventoryDocumentV1` digest;
- the exact versioned query, retained verbatim; and
- the implementation-owned extractor version.

Matches must be sorted by raw path bytes, then outer call start byte offset, then a documented deterministic final position tie-breaker. Counts alone are not evidence of equal results; spans and fragment hashes must distinguish different locations.

All input byte slices must be defensively copied before retention. All exported output byte slices and match collections must be defensive copies so callers cannot mutate retained census state.

### Limits and atomic failure

The approved default limits are exactly:

| Limit | Default |
|---|---:|
| `MaxFiles` | 256 files |
| `MaxFileBytes` | 1 MiB |
| `MaxTotalBytes` | 8 MiB |
| `MaxMatches` | 10,000 matches |

Every explicit limit must be positive. File and byte limits are enforced before parsing; the match limit is enforced while collecting.

Any inventory, supplied-set, digest, length, query, limit, parser, span, or match-limit error returns no result. A partial parser AST is never publishable after a parse error, and an accumulated partial match set is discarded on every failure.

These limits bound admitted input and published match count. They do not claim to bound parser heap amplification, stack use, CPU time, or process isolation.

## Affected future areas

No area is edited by this proposal beyond this file. A future authorized implementation is expected to affect only these surfaces:

| Future surface | Intended responsibility |
|---|---|
| `census/query.go` | Public query, input, limits/defaults, result/span, stable error, defensive-accessor, and extractor-version contracts. |
| `census/go_ast.go` | Inventory/file validation coordination, complete-file parsing, bounded callee normalization, exact physical spans, collection, and sorting. |
| `census/go_ast_test.go` | Candidate, validation, bounds, position, binding, sorting, and atomicity coverage. |
| `census/external_api_test.go` | External-consumer API and defensive-ownership coverage. |
| `docs/census.md` | Consumer-facing syntax-only scope, exact query semantics, bindings, limits, and non-goals. |
| `AGENTS.md` | A future narrow exception permitting this public subpackage while preserving the root-only default for other neutral public APIs. |
| `openspec/config.yaml` | Future alignment of `layout.public_api` only with the same narrow exception. |

The future `AGENTS.md` and `openspec/config.yaml` exception is approved only for the public `census` subpackage. Neither file is changed now, and the testing policy and rubric in `openspec/config.yaml` must remain untouched.

## Acceptance criteria

### Contract and scope

- [ ] The public API lives in a narrowly authorized `changeevidence/census` subpackage and remains neutral and evidence-only.
- [ ] Inputs are an existing `InventoryDocumentV1`, complete supplied raw path/content records, an exact versioned query, and positive limits.
- [ ] No policy input, inventory decoder duplication, inventory serialization round trip, hunk input, filesystem discovery, import loading, or dependency traversal is introduced.
- [ ] Exact completeness is claimed only for the supplied inventory scope and exact syntactic query; semantic completeness is explicitly disclaimed.
- [ ] The result retains the exact inventory digest, verbatim versioned query, and extractor version.

### Candidate behavior

- [ ] Repeated and legacy textual sites are all reported within the validated supplied scope.
- [ ] Comments and string literals containing matching text do not create candidates.
- [ ] A locally shadowed receiver spelling is admitted as a textual candidate without semantic classification.
- [ ] An alias spelling is excluded when it does not exactly equal the query receiver.
- [ ] The same textual callee used under different downstream invariants is reported without inferred classification.
- [ ] Non-selector callees do not match.

### AST normalization

- [ ] Repeated callee parentheses normalize to the underlying selector.
- [ ] Single-index and multi-index callee forms normalize only through `IndexExpr` and `IndexListExpr` around the callee.
- [ ] Ambiguous single-index value calls, including the `obj.Funcs[i]()` class, remain admitted textual candidates when their normalized selector spelling matches.
- [ ] No type, import, alias, shadow, dependency, or other semantic resolution narrows or expands the predicate.

### Physical coordinates and fragment binding

- [ ] Outer call and normalized selector ranges are exact physical half-open byte spans.
- [ ] Each endpoint has an exact zero-based byte offset and physical 1-based line and byte-column position.
- [ ] UTF-8 before or within a call preserves byte-based offsets and columns.
- [ ] CRLF source preserves exact byte accounting and LF-based physical lines.
- [ ] Multiline calls preserve exact call and selector boundaries.
- [ ] `//line` directives do not alter reported physical positions.
- [ ] Position ownership and range bounds are validated rather than accepted through offset clamping.
- [ ] Each match binds the raw path, file SHA-256, and SHA-256 of the exact full-call source fragment.
- [ ] Tests distinguish equal match counts with different locations, ranges, or fragments.

### Inventory and all-or-nothing behavior

- [ ] Missing supplied files, supplied extras, and duplicate supplied paths each fail with no result; existing V1 constructor/decoder rejection establishes duplicate-inventory antecedent invalidity, without reflection/unsafe or a widened census API.
- [ ] Content-hash drift and byte-length drift each fail with no result before parsing.
- [ ] A zero-value inventory fails with no result; no nil/pointer alternative or forged-private-state test is required.
- [ ] A valid explicitly constructed empty inventory with an empty supplied set succeeds with zero matches.
- [ ] Malformed Go returns no result even if the parser provides a partial AST.
- [ ] No validation, parsing, span, or limit failure exposes a partial match set.

### Bounds

- [ ] Defaults are exactly 256 files, 1 MiB per file, 8 MiB total input, and 10,000 matches.
- [ ] Every configurable limit rejects zero; signed representations must also reject negatives if selected, but unsigned limits have no representable negative case and require no widened input API.
- [ ] File-count, per-file-byte, and aggregate-byte breaches fail atomically before parsing.
- [ ] Match-limit breach fails atomically while collecting and publishes no accumulator.
- [ ] Documentation does not imply parser memory, stack, CPU, or process-isolation guarantees from these limits.

### Determinism and ownership

- [ ] Matches are deterministically sorted by raw path bytes, call start byte offset, and a documented final position tie-breaker.
- [ ] Input paths, input content retained by the operation, output paths, and output match collections cannot be mutated through caller-owned aliases.
- [ ] External-package coverage confirms public API usability and defensive input/output ownership.

### Repository alignment and delivery controls

- [ ] Any future policy alignment changes only `AGENTS.md` and `openspec/config.yaml` `layout.public_api` to express the narrow `census` exception.
- [ ] The `openspec/config.yaml` testing rubric remains unchanged.
- [ ] Novel production code across `census/query.go` and `census/go_ast.go` does not exceed the explicit 400-line ceiling.
- [ ] The complete acceptance surface is preserved even though tests are separately estimated at 450–700 lines.
- [ ] If the canonical 400 changed-line review budget is at risk, `ask-on-risk` pauses for a delivery decision; no chain strategy, automatic delivery-size exception, `size:exception`, or acceptance reduction is inferred.

## Explicit non-goals

This proposal does not authorize or design:

- semantic completeness, import resolution, type checking, alias resolution, shadow resolution, or dependency resolution;
- stores, DSLs, CLIs, serialization formats, adapters, mutations, runners, or receipts;
- Python work, a pilot, U10/U11 work, or any source implementation;
- runtime commands, tests, commits, pull requests, publication, or delivery actions;
- lifecycle or registry changes, including `.atl/skill-registry.md`, `.gitignore`, or `.pi` preflight state; or
- dependencies on the sibling V2 or `cnsic-profile-parity` change.

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| Users mistake a textual candidate for a resolved package or method call. | Name the output “candidate,” retain the exact query, document admitted shadowing/index ambiguity, and prohibit semantic claims. |
| A partial or drifted file set produces misleading completeness. | Validate exact bidirectional inventory membership, uniqueness, lengths, and digests before any parse. |
| Parser errors leak partial evidence. | Treat every non-nil parse error and every limit/validation error as atomic failure with no result. |
| `IndexExpr` ambiguity creates false positives. | Preserve the approved ambiguity explicitly and never imply type-instantiation certainty. |
| Position APIs apply line directives or conceal invalid endpoints. | Use unadjusted physical positions and independently validate ownership and bounds. |
| Mutable slices corrupt retained evidence. | Apply defensive ownership on all retained inputs and exported outputs, with external-package acceptance coverage. |
| Positive input limits are mistaken for parser isolation. | State that no parser heap, stack, CPU, or process-isolation bound is established. |
| The public subpackage conflicts with current layout policy. | Require a future narrow `layout.public_api` exception in both policy surfaces without changing the testing rubric. |
| Full acceptance exceeds the selected review shape. | Preserve the 400-line production ceiling and 450–700-line test estimate; stop under `ask-on-risk` rather than dropping acceptance or inferring an exception. |

## Rollback

Because this change creates only a planning artifact, rollback is deletion of `openspec/changes/go-ast-census/proposal.md`. No source, test, configuration, registry, lifecycle, or runtime state is changed.

If a future implementation is separately authorized and later withdrawn, remove the `census` implementation, its tests and documentation, and the narrow public-layout exception together. Existing inventory contracts remain unchanged because this proposal does not modify or replace them.

## Success criteria

This proposal succeeds when reviewers can verify, without reconstructing product decisions, that the future slice:

1. produces a complete and deterministic AST census for exactly the validated supplied inventory and exact textual query;
2. binds every result to inventory, query, extractor, file content, exact call fragment, and physical call/selector locations;
3. remains syntax-only, neutral, defensively owned, bounded, and all-or-nothing;
4. preserves every acceptance category and all explicit non-goals;
5. keeps novel production code at or below 400 lines while treating the separate 450–700 test estimate as an unresolved review-budget risk; and
6. requires an explicit future delivery decision under `ask-on-risk` rather than silently granting a size exception or reducing acceptance.

## Evidence basis

The proposal preserves the confirmed product decisions in `explore.md` and `preproposal.md` revision 3. Its Go parser assumptions are supported by `research.md` revision 3 claims Q1-A through Q4-A, based only on official Go `go1.25.10` source for `go/ast`, `go/parser`, `go/token`, and `go/scanner`. The evidence supports syntactic node shapes, single-index ambiguity, unadjusted physical byte positions, complete-file parser operation, skipped object resolution, partial-AST rejection, and the absence of caller-configurable parser resource isolation. These references support planning only and do not authorize implementation or delivery.
