# Go AST Census Specification

## Purpose

Provide neutral, evidence-only discovery of every exact textual Go selector-call candidate in one validated, caller-supplied `InventoryDocumentV1` scope. The census SHALL make no semantic, policy, or delivery-authority determination.

## Requirements

### Requirement: GAC-001 — Neutral versioned census contract

The system MUST expose a neutral, versioned public census contract that accepts an `InventoryDocumentV1` value, complete caller-supplied raw path/content file records, an exact versioned selector-call query, and explicit positive limits. Every successful result MUST bind and retain the exact inventory digest, verbatim query, and implementation-owned extractor version. The contract MUST NOT accept a policy document, hunk view, filesystem-discovery input, inventory decoder/serialization round trip, DSL, or nil-inventory alternative API. A zero-value or otherwise invalid inventory value MUST fail.

#### Scenario: GAC-001-S01 — Valid contract binding

- GIVEN a valid inventory value, matching supplied files, a valid versioned query, and positive limits
- WHEN a census succeeds
- THEN its result contains the exact inventory digest, verbatim query, and extractor version

#### Scenario: GAC-001-S02 — Invalid query is rejected without a DSL

- GIVEN an empty receiver, empty selector, or unsupported query version
- WHEN a caller requests a census
- THEN the system MUST return no result without accepting a DSL or other query form

### Requirement: GAC-002 — Exact supplied-scope validation

Before parsing, the system MUST validate inventory identity through owned public accessor state, including canonical-byte SHA-256 identity against the inventory document digest. It MUST require a one-to-one, byte-for-byte raw-path correspondence between inventory entries and supplied files; reject duplicate supplied paths, missing supplied files, and supplied extras; and require each supplied content SHA-256 and `uint64` byte length to equal its inventory entry. `InventoryDocumentV1` is a value with private fields: its root constructor/decoder owns rejection of duplicate inventory paths antecedent to census input, and the census MUST NOT duplicate that decoder responsibility or fabricate private document state. A valid explicitly constructed empty inventory paired with an empty supplied set MUST succeed with zero matches. Completeness MUST be claimed only for this exact validated supplied inventory scope, never for a partial changed-line view.

#### Scenario: GAC-002-S01 — Duplicate supplied path fails atomically

- GIVEN two supplied file records with the same raw path
- WHEN a census is requested
- THEN the system MUST return no result before parsing

#### Scenario: GAC-002-S02 — Missing supplied path fails atomically

- GIVEN an inventory path without exactly one supplied file
- WHEN a census is requested
- THEN the system MUST return no result before parsing

#### Scenario: GAC-002-S03 — Extra supplied path fails atomically

- GIVEN a supplied file whose raw path has no inventory entry
- WHEN a census is requested
- THEN the system MUST return no result before parsing

#### Scenario: GAC-002-S04 — Drifted supplied content fails atomically

- GIVEN a supplied path present in the inventory whose content digest or byte length differs from its inventory entry
- WHEN a census is requested
- THEN the system MUST return no result before parsing

#### Scenario: GAC-002-S05 — Zero-value inventory fails through the value API

- GIVEN a zero-value inventory value
- WHEN a caller requests a census
- THEN the system MUST return no result; census coverage uses no nil or pointer alternative, unsafe, reflection, forged private state, extra policy, or decoder API

#### Scenario: GAC-002-S06 — Complete empty scope succeeds

- GIVEN a valid explicitly constructed empty inventory and no supplied files
- WHEN a census is requested with valid query and limits
- THEN the system MUST succeed with an empty match collection

### Requirement: GAC-003 — Exact syntactic candidate predicate

The system MUST derive candidates only by complete-file `go/parser` parsing and `go/ast` traversal. It MUST NOT use textual scanning, filesystem or runtime discovery, import loading, type checking, object resolution, alias resolution, shadow resolution, dependency traversal, or semantic classification. A candidate MUST be an outer `CallExpr` whose callee, after repeatedly unwrapping only `ParenExpr`, `IndexExpr`, and `IndexListExpr`, is a selector with an identifier receiver and selector whose spellings exactly equal the query receiver and selector. The system MUST report every such candidate in the validated scope and MUST NOT claim semantic completeness.

#### Scenario: GAC-003-S01 — Exact textual selector sites are candidates

- GIVEN validated Go source containing `os.Open(...)`, `(os.Open)(...)`, and `os.Open[T](...)` for query `os.Open`
- WHEN a census succeeds
- THEN the system MUST report each call as a candidate

#### Scenario: GAC-003-S02 — Nested parentheses and multi-index callees normalize

- GIVEN validated Go source containing nested `ParenExpr` wrappers and multi-index `IndexListExpr` wrappers around a matching selector callee
- WHEN a census succeeds
- THEN the system MUST normalize each only through the approved callee wrappers and report the matching calls

#### Scenario: GAC-003-S03 — Syntax-only decoys and index ambiguity remain explicit

- GIVEN a locally shadowed `os.Open`, `alias.Open`, `obj.Funcs[i]()`, comments, strings, and non-selector calls for query `os.Open` or `obj.Funcs` as applicable
- WHEN a census succeeds
- THEN shadowed exact spelling and matching ambiguous single-index forms MUST be candidates, while alias spelling, comments, strings, and non-selector calls MUST NOT be candidates

#### Scenario: GAC-003-S04 — Repeated and legacy sites have no invariant classification

- GIVEN repeated or legacy matching textual sites used under different downstream invariants
- WHEN a census succeeds
- THEN the system MUST report the sites as candidates without classifying their invariants or semantic identities

### Requirement: GAC-004 — Exact physical source evidence

For every match, the system MUST retain the raw path, supplied-file SHA-256, and SHA-256 of the exact full-call fragment. It MUST report both the outer call and normalized selector as exact half-open byte spans `[start,end)` over supplied content. Each endpoint MUST include a zero-based byte offset and an unadjusted physical one-based line and byte-column position. The system MUST preserve exact byte accounting for UTF-8, CRLF, multiline calls, and `//line` directives, and MUST validate endpoint ownership and bounds rather than accept clamped offsets.

#### Scenario: GAC-004-S01 — Physical coordinates remain unadjusted

- GIVEN a validated file containing UTF-8 bytes, CRLF lines, a multiline matching call, and a `//line` directive
- WHEN a census succeeds
- THEN call and selector half-open spans, endpoint offsets, lines, and byte columns MUST identify the exact physical supplied bytes without line-directive remapping

#### Scenario: GAC-004-S02 — Fragment binding distinguishes locations

- GIVEN two census inputs with equal match counts but distinct matching locations, ranges, or full-call source fragments
- WHEN censuses succeed
- THEN their raw paths, file digests, fragment digests, or spans MUST distinguish the results

### Requirement: GAC-005 — Deterministic defensive result ownership

The system MUST defensively own all input bytes retained by the operation and MUST return defensive copies of exported output byte slices and match collections. It MUST order matches by raw path bytes, then outer-call start byte offset, then a documented deterministic final position tie-breaker.

#### Scenario: GAC-005-S01 — Caller mutation cannot alter evidence

- GIVEN a successful census result
- WHEN a caller mutates original input slices or exported result slices and collections
- THEN retained census evidence and later exported observations MUST remain unchanged

#### Scenario: GAC-005-S02 — Stable raw-path order

- GIVEN matching calls across raw paths and positions supplied in different input orders
- WHEN censuses succeed
- THEN their match order MUST be identical under the documented raw-path, call-start, and final-tie-breaker ordering

#### Scenario: GAC-005-S03 — External value API is usable and owned

- GIVEN an external-package consumer importing the public API and passing inventory and query values
- WHEN it obtains and mutates exported observations
- THEN the public value API MUST be usable and defensive input and output ownership MUST hold

### Requirement: GAC-006 — Positive bounded inputs and matches

The API-provided default limits MUST be exactly 256 files, 1 MiB per file, 8 MiB total supplied bytes, and 10,000 matches. Every explicit limit MUST be positive. Limits MAY use `uint64`; zero is invalid, and no impossible negative-`uint64` case requires coverage. The system MUST enforce file-count, per-file-byte, and aggregate-byte limits before parsing, and MUST enforce the match limit while collecting. These limits MUST NOT be represented as guarantees of parser heap, stack, CPU, or process isolation.

#### Scenario: GAC-006-S01 — File-count equality succeeds and excess fails atomically

- GIVEN otherwise valid inputs with exactly `MaxFiles` files and then with one more file
- WHEN a census is requested
- THEN the equality case MUST succeed and the excess case MUST return no result before parsing

#### Scenario: GAC-006-S02 — Per-file-byte equality succeeds and excess fails atomically

- GIVEN otherwise valid inputs with a file exactly `MaxFileBytes` long and then one byte longer
- WHEN a census is requested
- THEN the equality case MUST succeed and the excess case MUST return no result before parsing

#### Scenario: GAC-006-S03 — Aggregate-byte equality succeeds and excess fails atomically

- GIVEN otherwise valid inputs totaling exactly `MaxTotalBytes` and then one byte more
- WHEN a census is requested
- THEN the equality case MUST succeed and the excess case MUST return no result before parsing

#### Scenario: GAC-006-S04 — Match equality succeeds and excess fails atomically

- GIVEN validated parsable inputs producing exactly `MaxMatches` candidates and then one more candidate
- WHEN a census is requested
- THEN the equality case MUST succeed and the excess case MUST return no result with no accumulated matches published

#### Scenario: GAC-006-S05 — Zero file-count limit is invalid

- GIVEN `MaxFiles` equal to zero
- WHEN a census is requested
- THEN the system MUST return no result

#### Scenario: GAC-006-S06 — Zero per-file-byte limit is invalid

- GIVEN `MaxFileBytes` equal to zero
- WHEN a census is requested
- THEN the system MUST return no result

#### Scenario: GAC-006-S07 — Zero aggregate-byte limit is invalid

- GIVEN `MaxTotalBytes` equal to zero
- WHEN a census is requested
- THEN the system MUST return no result

#### Scenario: GAC-006-S08 — Zero match limit is invalid

- GIVEN `MaxMatches` equal to zero
- WHEN a census is requested
- THEN the system MUST return no result

### Requirement: GAC-007 — Atomic syntax and validation failure

Any inventory, supplied-set, digest, length, query, limit, parser, span, or match-limit error MUST return no result. A parser error MUST invalidate any partial AST and accumulated matches; no partial match set is publishable after any failure.

#### Scenario: GAC-007-S01 — Malformed source after an earlier candidate

- GIVEN a supplied file with a matching call followed by malformed Go syntax
- WHEN parsing reports an error
- THEN the system MUST return no result and MUST expose no earlier match

#### Scenario: GAC-007-S02 — Every failure class discards partial evidence

- GIVEN an operation that has encountered any validation, parsing, span, or limit failure
- WHEN the operation returns
- THEN it MUST return no result and expose no partial match set

### Requirement: GAC-008 — Evidence-only scope and documented non-guarantees

The census and its documentation MUST describe matches as syntax-only textual candidates, not resolved package, method, import, type, alias, shadow, dependency, or invariant identities. Exact completeness applies only to the validated supplied inventory scope and exact syntactic query; semantic completeness is expressly not claimed. The system MUST remain evidence-only and MUST NOT approve, reject, gate, block, merge, deploy, release, or assign delivery authority. Documentation MUST state that input limits do not establish parser resource isolation.

#### Scenario: GAC-008-S01 — Consumer can identify the non-guarantees

- GIVEN a result containing a shadowed, ambiguous indexed, repeated, or legacy textual candidate
- WHEN a consumer reads the contract documentation
- THEN it MUST be able to determine that semantic identity, invariant classification, and parser isolation are not guaranteed

### Requirement: GAC-009 — Narrow public-layout exception

The public census contract MUST be confined to the narrowly authorized `changeevidence/census` subpackage. A future repository-guidance and `layout.public_api` change MUST express only this layout exception while preserving the existing testing rubric and all other public-layout constraints; it MUST NOT be made by this planning slice.

#### Scenario: GAC-009-S01 — Future policy alignment review

- GIVEN a future change that introduces the public census subpackage
- WHEN repository guidance and layout configuration are reviewed
- THEN they MUST permit only the `census` public-subpackage exception and MUST leave the testing rubric unchanged

### Requirement: GAC-010 — Review-budget preservation

The complete acceptance surface in this specification MUST be preserved. Novel production code for the census query and Go-AST behavior MUST NOT exceed 400 lines. Tests are separately estimated at 450–700 lines and remain part of the complete acceptance surface. If the canonical 400 changed-line review budget is at risk, delivery selection MUST pause under `ask-on-risk`; no chain strategy, automatic size exception, or acceptance reduction MAY be inferred.

#### Scenario: GAC-010-S01 — Review-budget risk

- GIVEN an authorized implementation whose projected changed lines exceed the review budget
- WHEN delivery planning reaches that risk
- THEN planning MUST pause for an explicit delivery decision without removing a requirement or granting an exception
