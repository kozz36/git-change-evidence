# Design: Exact supplied-scope Go AST census

## Technical approach and evidence

Implement a future pure `github.com/kozz36/git-change-evidence/census` package: validate an immutable Inventory V1 value against complete supplied bytes, parse every supplied file, locate every exact textual selector-call candidate, and publish owned evidence atomically. Completeness covers exactly that validated scope and syntactic predicate—not semantic Go identity or strict language validity. The parser is permissive; parser success is not semantic certification.

This artifact implements the planning intent of [proposal.md](proposal.md) and all 10 requirements/30 scenarios in [spec.md](specs/go-ast-census/spec.md). Directly read research/preproposal revision 3 supports wrappers (Q1-A), physical coordinates (Q2-A), whole-file parsing (Q3-A), and resource caveats (Q4-A). CodeGraph had already failed; source reads used the authorized fallback without initialization/retry.

Direct source anchors: `inventory_v1.go`, `inventory_v1_build.go`, `inventory_v1_decode.go`, their contract/ownership tests, `contract.go`, `policy_contract.go`, `result_v1.go`, `result_v1_antecedents.go`, `go.mod`, `AGENTS.md`, and `openspec/config.yaml`. The actual module is `github.com/kozz36/git-change-evidence`, Go 1.25.10. Bootstrap-era statements in AGENTS/README about absent source/module are stale; they are not changed here.

## Architecture decisions

| Option | Tradeoff | Decision and rationale |
|---|---|---|
| Public value API in `census` versus root expansion | Requires approved narrow layout exception | Use `census`; depend one-way on root contracts, preserving Functional Core / Imperative Shell and Hexagonal boundaries. No shell effects are needed. |
| Accessor validation versus policy-dependent re-decoding | Cannot independently certify fabricated private state | Use owned accessors and canonical-byte digest identity. Root construction owns validity, canonical encoding, path rules, and duplicate-inventory rejection. |
| AST syntax versus semantic resolution | Includes shadowed and ambiguous indexed-value calls | Use only the approved syntax predicate; no imports, types, aliases, shadows, dependencies, or invariant classification. |
| Private result state versus exported mutable aggregate | Requires defensive accessors | Follow root inventory/result ownership conventions; no serialization or result digest. |
| Explicit input/count limits versus isolation machinery | No heap/stack/CPU isolation | Keep the four approved limits only; no depth controls, subprocesses, or new concurrent surface. |

## Public contracts

The following declarations specify the future API, not implementation. `ce` denotes the existing root import; digest values remain lowercase SHA-256 hexadecimal strings, matching root conventions.

```go
const SelectorCallQueryV1 ce.ContractVersion = "git-change-evidence.selector-call-query/v1"
const ExtractorVersion ce.ContractVersion = "git-change-evidence.go-ast-census/v1"

type File struct { Path, Content []byte }
type QueryV1 struct {
    Version ce.ContractVersion
    Receiver, Selector string
}
type Limits struct {
    MaxFiles, MaxFileBytes, MaxTotalBytes, MaxMatches uint64
}
func DefaultLimits() Limits
func GoASTV1(ce.InventoryDocumentV1, []File, QueryV1, Limits) (ResultV1, error)

type Position struct { Offset, Line, Column int }
type Span struct { Start, End Position }
type MatchV1 struct {
    Path []byte
    FileSHA256, FragmentSHA256 ce.Digest
    Call, Selector Span
}
type ResultV1 struct { /* private bindings and []MatchV1 */ }
func (ResultV1) InventoryDigest() ce.DocumentDigest
func (ResultV1) Query() QueryV1
func (ResultV1) ExtractorVersion() ce.ContractVersion
func (ResultV1) Matches() []MatchV1
```

Query strings are retained verbatim: no trimming, case folding, Unicode normalization, patterns, or DSL. Require the exact supported version and `token.IsIdentifier` for each name (including rejection of empty strings and keywords). The query version identifies the predicate contract; the implementation-owned extractor version identifies extraction behavior. Consumers cannot select an extractor. Future behavior changes must review these identities; neither is a serialization schema or sibling V2 dependency.

`ResultV1` privately stores inventory digest, query, extractor, and matches. Failure returns exactly `ResultV1{}, err`: empty bindings and nil matches. Successful empty scope retains all bindings and returns a non-nil empty match slice. `Matches()` clones the collection and every path. Other accessors return values. No AST, source content, token file, or caller file slice survives in the result.

The synchronous operation copies retained input paths/content before hashing/parsing; result paths have no caller-owned aliases. Strings, digests, and positions are value data. Mutating input/output slices after return cannot alter later observations. Callers must not concurrently mutate arguments during the call; defensive copying does not promise a race-free snapshot under concurrent mutation. No synchronization API is introduced.

## Deterministic data flow

`values → query/limits → antecedent → exact file set/bounds → owned content identity → whole-file AST → spans/fragments → sorted owned result`

1. **Query and limits.** Validate version, receiver, selector, then the four limits in declaration order. `DefaultLimits()` returns a fresh value: **256**, **1,048,576**, **8,388,608**, **10,000**. Every zero is invalid; unsigned negatives are unrepresentable. Compare actual counts as `uint64`; never convert caller limits to `int` for allocation.
2. **Inventory antecedent.** Read `CanonicalBytes()`, `Digest()`, `AccountingPolicyDigest()`, and `Entries()` once. Require nonempty canonical bytes, nonempty policy link, non-nil entries, and `DocumentDigest(hex(sha256(canonical))) == Digest()`. This rejects the zero value while admitting constructor-created empty inventories. Root private fields and defensive accessors preserve consistency between canonical bytes and entry state. Do not call `DecodeInventoryV1(raw, policy)`, duplicate JSON/path/schema validation, reconstruct serialization, accept policy, or add an injectable/forgeable document interface. Unlike `result_v1_antecedents.go`, this operation has no policy-dependent round trip.
3. **Scope and admission before any parse.** Check both inventory/supplied file counts against `MaxFiles`. Copy path records and sort supplied records with `bytes.Compare`. Reject adjacent supplied duplicates first; compare against canonical inventory entry order, detecting extra/missing paths byte-for-byte. No extension filtering, UTF-8 requirement, path cleaning, case folding, or filename interpretation: every valid inventory path is eligible, including newline, backslash, and high-bit bytes. Inventory constructor/decoder already owns duplicate antecedents; census does not recreate that invariant.
4. **Byte budgets and identity before any parse.** For each matched record require `uint64(len(Content)) <= MaxFileBytes`. With invariant `total <= MaxTotalBytes`, reject `n > MaxTotalBytes-total` before `total += n`; equality succeeds without overflow. Finish these checks across the entire set, then copy content, require length equal to `ByteLength`, and compare computed SHA-256 with `ContentSHA256`. Check length before hash for deterministic drift classification. Store owned validated records; parsing cannot begin until all records pass.
5. **Complete-file parsing.** In raw-path order create a fresh `token.FileSet` per file and call `parser.ParseFile(fset, "supplied.go", content, parser.SkipObjectResolution)`, where `content` is always statically `[]byte`, including typed nil/empty bytes. The interface argument is therefore never nil, so the parser never opens the diagnostic filename. No `PackageClauseOnly`, `ImportsOnly`, import loading, object resolution, or build-tag/file-name selection. Reject every non-nil parse error before walking even a non-nil partial AST; treat an unexpected nil AST as parse failure.
6. **Candidate traversal.** Walk the complete AST with `ast.Inspect`, continuing into children of matched calls to discover nested calls. For each `*ast.CallExpr`, iteratively replace only callee `ParenExpr.X`, `IndexExpr.X`, or `IndexListExpr.X`; stop at any other node. Require the terminal `SelectorExpr.X` itself to be `*ast.Ident` and both identifier spellings to match exactly. Do not unwrap inside the receiver: `(os).Open()` and `obj.field.Open()` do not satisfy receiver `os`. `((os.Open[T,U]))()` does; `obj.Funcs[i]()` intentionally does for query `obj.Funcs`, irrespective of generic-versus-indexed-value ambiguity.
7. **Atomic collection and publication.** Before each append, if collected count equals `MaxMatches`, fail and discard everything; otherwise append and increment. Equality at completion succeeds; no overflowing increment or `MaxMatches+1` arithmetic is needed. A traversal error stops further useful collection and returns the zero result. After all files succeed, sort and construct the private result; no callbacks, iterators, or partial accumulators escape.

All precedence above is deterministic for permuted equivalent inputs. Within one validation stage, inspect sorted paths. Limits bound admitted content and published match count only; path/canonical-accessor allocations, parser amplification, CPU, stack, panic recovery, and process isolation are not bounded guarantees.

## Physical spans, fragment identity, and order

For outer call and normalized selector use their own `Pos()`/`End()`. Obtain the owning token file from the parsed file position and require its size to equal source length. For each endpoint, reject `NoPos`, positions below file base, positions whose base-relative distance exceeds file size, or positions not owned by this file in the FileSet. Perform these checks **before** `Offset`/`PositionFor`, because offset conversion clamps invalid positions. EOF endpoints equal to file size are allowed.

Require `0 <= start < end <= len(content)` for both spans and selector containment inside the call. Only then use the owning file's `Offset(p)` and `PositionFor(p, false)`; validate consistent offsets and positive line/column values. Offsets are zero-based; physical line and byte-column are one-based. UTF-8 occupies its actual byte width, CR remains a byte in CRLF, LF determines lines, and `//line` never remaps these coordinates. Follow token.File's physical EOF convention; do not synthesize a trailing empty line.

Compute `FragmentSHA256` over exactly `content[Call.Start.Offset:Call.End.Offset]`, including callee wrappers, arguments, and internal whitespace/comments—not normalized selector text or result serialization. `FileSHA256` covers the complete validated content.

Ordering is lexicographic **raw path bytes, call start offset, call end offset, selector start offset, selector end offset**. The latter three form the documented final position tie-breaker. Distinct nodes normally differ earlier; exact equal keys are observationally identical. No map iteration order or locale influences results.

## Error taxonomy

Reuse public `*ce.ContractError{Field, Code}` and `errors.As`, rather than adding a parallel hierarchy. These pairs are stable; strings rendered by `Error()` are diagnostic, not authority. No error exposes matches, ASTs, or source slices.

| Field | Code(s) |
|---|---|
| `query.version` | `unsupported_version` |
| `query.receiver`, `query.selector` | `invalid_identifier` |
| `limits.max_files`, `limits.max_file_bytes`, `limits.max_total_bytes`, `limits.max_matches` | `invalid_limit` (zero), `limit_exceeded` (breach) |
| `inventory` | `invalid_document` |
| `files.path` | `duplicate_path`, `missing_path`, `extra_path` |
| `files.byte_length` | `length_mismatch` |
| `files.content_sha256` | `digest_mismatch` |
| `source` | `parse_error` |
| `span` | `invalid_span` |

Root duplicate inventory errors remain root `ContractError` results, not census failure injection. Publicly unreachable corrupt-private-state branches need no reflection/unsafe tests. Span safety can be tested with private token-file helper inputs without forging inventory or exposing a testing API.

## Testing strategy and complete scenario mapping

Future tests follow RED → GREEN → TRIANGULATE → REFACTOR, co-located and table-driven. `census/go_ast_test.go` covers behavior/private span and arithmetic boundaries; `census/external_api_test.go` uses `package census_test` and only exported root/census APIs. Fixture policy construction uses `NewAccountingPolicyV1(PolicyInput{Categories: []CategoryInput{{Name: "go"}}, Default: "go"})` solely to construct the existing inventory; policy is never passed to census. No subprocess, real repository, filesystem fixture, or golden-update machinery is needed.

Every positive row compares full ordered records, query/extractor/inventory bindings, literal fragments, digests, and ranges as applicable—not just counts. Every negative row asserts typed field/code and all zero-result accessors. Concrete independent coordinate oracles include:

- `"package p\nfunc f(){ (os.Open)(\"x\") }\n"`: call `[20,34)`, physical `(2,11)→(2,25)`; selector `[21,28)`, `(2,12)→(2,19)`; fragment `(os.Open)("x")`.
- `"package p\r\n//line fake.go:900\r\nfunc f(){ _ = \"é\"; os.Open(\r\n\"x\") }\r\n"`: call `[51,65)`, `(3,21)→(4,5)`; selector `[51,58)`, `(3,21)→(3,28)`; fragment `"os.Open(\r\n\"x\")"` as a Go string literal. Add Unicode identifiers inside matching selectors, not only before them.

Expected hashes come from those literal full files/fragments with standard SHA-256, never by reusing extractor offsets or normalization helpers. Include equal-count changed-location and changed-argument fixtures, plus EOF, nested-call and nested-wrapper ranges.

Each row below maps one specification scenario; identifiers are exact.

| Scenario | Planned assertion / evidence |
|---|---|
| GAC-001-S01 | External successful value call retains exact three bindings and full expected match. |
| GAC-001-S02 | Empty receiver/selector and unsupported version fail independently; lexical invalidity table, no query rewriting. |
| GAC-002-S01 | Duplicate supplied paths, equal/different contents, fail before deliberately malformed source is parsed. |
| GAC-002-S02 | Missing supplied path fails before parsing another malformed file. |
| GAC-002-S03 | Extra supplied path fails before parsing; do not rely solely on unequal counts. |
| GAC-002-S04 | Same-length content substitution isolates hash drift; unequal length isolates length drift; both precede parsing. |
| GAC-002-S05 | Literal zero inventory fails; root duplicate constructor/decoder tests remain antecedent evidence, no private forging. |
| GAC-002-S06 | Explicit constructed empty scope succeeds with bindings and non-nil empty matches. |
| GAC-003-S01 | Exact direct, parenthesized, single-index candidates have expected literal call/selector fragments. |
| GAC-003-S02 | Nested parentheses and multi-index wrappers, including combinations and nested calls, retain exact outer ranges. |
| GAC-003-S03 | Shadow candidate and indexed-value candidate included; alias/comments/strings/non-selector/receiver-wrapper decoys excluded. |
| GAC-003-S04 | Repeated/legacy calls under distinct usages retain every separate location without classification. |
| GAC-004-S01 | Assert literal UTF-8/CRLF/multiline/directive coordinates above plus Unicode selector and endpoint cases. |
| GAC-004-S02 | Equal-count shifted locations and changed arguments differ in full records and exact fragment hashes. |
| GAC-005-S01 | Mutate original file/path/content slices, returned paths and collections; re-read unchanged evidence. |
| GAC-005-S02 | Permute supplied files with newline/high-bit/backslash/case-distinct paths; assert byte order and positional ties. |
| GAC-005-S03 | External package constructs real inventory/query values and repeats binding and defensive-accessor assertions. |
| GAC-006-S01 | Exactly MaxFiles succeeds; one extra fails before parsing, using matching inventories. |
| GAC-006-S02 | Exact MaxFileBytes succeeds; one byte extra fails, with parsable whitespace-padded source. |
| GAC-006-S03 | Exact MaxTotalBytes succeeds; one byte extra fails; private arithmetic test exercises near-uint64 overflow without huge allocation. |
| GAC-006-S04 | Exact MaxMatches returns expected records; one extra discards accumulated candidates across files. |
| GAC-006-S05 | Zero MaxFiles fails with other limits positive; also assert exact default tuple. |
| GAC-006-S06 | Zero MaxFileBytes fails independently. |
| GAC-006-S07 | Zero MaxTotalBytes fails independently. |
| GAC-006-S08 | Zero MaxMatches fails independently, including empty scope. |
| GAC-007-S01 | Matching call followed by malformed syntax yields no evidence despite partial AST; repeat after an earlier valid file. |
| GAC-007-S02 | Shared failure assertions cover every taxonomy pair; private span tests cover NoPos, foreign file, out-of-bounds, reversed/noncontained ranges. |
| GAC-008-S01 | Review consumer docs against shadow/index/repeated/legacy fixtures and explicit semantic/resource/authority non-guarantees. |
| GAC-009-S01 | Future diff review permits only census layout exception; existing testing rubric remains byte-for-byte unchanged. |
| GAC-010-S01 | Review production and gross authored-line accounting separately; record explicit delivery decision before oversized work. |

Future validation: focused `go test ./census`, then `go test ./...` and the unchanged check-only formatting command in `openspec/config.yaml`. Race testing remains conditional on actual concurrency-bearing changes; this design creates none. No tests or runtime commands were run for this artifact.

## Future file changes, feasibility, and rollout

Only this `design.md` is written now. Future authorization would cover:

| File | Future action / responsibility |
|---|---|
| `census/query.go` | Create public contracts, defaults, private result and defensive accessors. |
| `census/go_ast.go` | Create complete validation-to-publication pipeline, errors, wrappers, spans, hashes, sorting. |
| `census/go_ast_test.go` | Create scenario tables and private safety-boundary tests. |
| `census/external_api_test.go` | Create external consumer and ownership coverage. |
| `docs/census.md` | Create usage, exact predicate, bindings, limits, and non-guarantees. |
| `AGENTS.md` | Modify only narrow approved public-census layout exception. |
| `openspec/config.yaml` | Modify `layout.public_api` only for that same exception; no testing/rubric edits. |

**Feasibility is tight, not demonstrated.** Ordinary formatted production estimate: API/defaults/accessors 85–105 lines; scope/identity/budget validation 100–120; parsing/collection 65–75; spans 45–55; normalization/order/hash helpers 30–40: **325–395 novel production lines total** across the two production files. This includes the complete vertical and all defensive checks, not a toy extractor. It leaves at most 5 lines at the upper estimate; do not compress formatting, move production into tests, or omit acceptance to claim compliance. If implementation cannot fit **400**, stop and ask before proceeding.

Tests remain **450–700** lines; consumer documentation approximately **60–100**, layout edits approximately **6–12** changed lines. Gross future implementation forecast is therefore **841–1,207 authored changed lines**, excluding planning artifacts, which also count when included in review. Fitting the complete work into one gross400 review is impossible even at the test-only minimum. The existing 160 changed-lines-per-Go-file threshold will also require its recorded crossing/package-boundary review during future apply; it is not silently waived or edited.

**Pause under ask-on-risk now:** which delivery shape should a later authorized implementation use—explicitly selected reviewable chaining, or explicitly accepted `size:exception`? Neither is selected here; chain strategy remains deferred. The production ceiling remains 400 under either choice. No implementation or delivery proceeds from this design alone.

No migration, feature flag, new store, DSL, CLI, serialization, adapter, mutation, runner, receipt, Python work, U10/U11, pilot, or sibling V2 dependency is introduced. Threat matrix: N/A—no routing, shell, subprocess, VCS/PR automation, executable classification, or process integration. The census remains evidence-only and cannot approve, reject, gate, block, merge, deploy, release, or assign delivery authority; validation errors are technical input failures only.

If separately authorized, introduce implementation/tests/docs and the narrow layout exception together in the explicitly chosen delivery shape. Rollback removes those additions and only that exception, leaving Inventory V1 unchanged. Product choices are settled; the remaining decision is delivery sizing, and production-line feasibility must be measured during future authorized work.
