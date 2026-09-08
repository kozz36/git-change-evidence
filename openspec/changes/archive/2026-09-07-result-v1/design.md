# Technical Design: Result/Accounting V1 Canonical Contract

## Status and design boundary

This design implements the approved Result V1 proposal and the four independently passing specifications. It is a pure, additive root-package contract built on the existing Policy V1, Inventory V1, committed-snapshot, accounting, and Report V1 contracts at `origin/main=a2b9654`.

Result V1 is evidence only. It reports classifications, measurements, totals, and factual `exceeded` observations; it does not decide approval, rejection, gating, blocking, merge, deployment, release, or delivery.

The design deliberately excludes moving refs and retries, Git or filesystem acquisition, inventory acquisition, untracked accounting, CLI/projection/publication expansion, CNSIC, configuration formats, module/release/distribution choices, and delivery mechanics or authority.

## Architecture decisions

### 1. A versioned immutable root document

Add the following public root-package contract, using `DocumentDigest` consistently with Policy V1 and Inventory V1:

```go
const AccountingResultContractV1 ContractVersion = "git-change-evidence.accounting-result/v1"

type ResultEntryMeasurementKind string
const (
    ResultEntryCountableV1    ResultEntryMeasurementKind = "countable"
    ResultEntryNonCountableV1 ResultEntryMeasurementKind = "non_countable"
)

type ResultObservationKind string
const (
    ResultLineThresholdObservationV1 ResultObservationKind = "line_threshold"
    ResultRatioObservationV1         ResultObservationKind = "ratio"
)

type ResultEntryMeasurementV1 struct {
    Kind                 ResultEntryMeasurementKind
    Additions, Deletions uint64 // both zero for non_countable
}
type ResultEntryV1 struct {
    Path        []byte
    Category    string
    Measurement ResultEntryMeasurementV1
}
type ResultCategoryTotalV1 struct {
    Category                         string
    Additions, Deletions, NonCountable uint64
}
type ResultObservationV1 struct {
    Kind               ResultObservationKind
    Category           string
    Reference          string // ratio only
    LineThresholdLimit uint64 // line_threshold only
    RatioLimit         string // ratio only; exact Policy V1 decimal
    Available          bool
    Actual             uint64 // available line_threshold only
    Numerator          uint64 // available ratio only
    Denominator        uint64 // available ratio only
    Exceeded           bool
}
type ResultDocumentV1 struct { /* canonical, digest, links, revisions, views are private */ }

func NewAccountingResultV1(policy PolicyDocumentV1, inventory InventoryDocumentV1, snapshot CommittedSnapshot) (ResultDocumentV1, error)
func DecodeAccountingResultV1(raw []byte, policy PolicyDocumentV1, inventory InventoryDocumentV1) (ResultDocumentV1, error)

func (r ResultDocumentV1) CanonicalBytes() []byte
func (r ResultDocumentV1) Digest() DocumentDigest
func (r ResultDocumentV1) AccountingPolicyDigest() DocumentDigest
func (r ResultDocumentV1) InventoryDigest() DocumentDigest
func (r ResultDocumentV1) Revisions() RevisionIdentity
func (r ResultDocumentV1) Entries() []ResultEntryV1
func (r ResultDocumentV1) Totals() []ResultCategoryTotalV1
func (r ResultDocumentV1) Observations() []ResultObservationV1
```

`ResultObservationV1` is a closed discriminated view, not a generic diagnostic/status model. Its irrelevant fields are zero for its `Kind`; the canonical encoder emits only the fields authorized for that variant. This makes an available zero measurement distinguishable from an unavailable measurement without exposing a nullable, caller-mutable union.

`ResultDocumentV1` owns copies of canonical bytes, classifications (including each path byte slice), totals, and observations. Every byte/slice-returning method returns a new copy; `Entries` deep-copies every `Path`; scalar strings, digests, and revisions are safe value copies. A zero or error-returned value has nil canonical bytes and views, empty digest/links, and zero revisions.

Construction accepts exact antecedent documents, not their digests and not `AccountingPolicy`. The legacy types remain available solely to legacy callers of `Account`.

### 2. Validate the antecedent chain before accounting

A private `validateResultAntecedents(policy, inventory)` routine is the sole antecedent boundary for both construction and decoding:

1. Obtain `policy.CanonicalBytes()`, strictly decode it with `DecodeAccountingPolicyV1`, recompute SHA-256, and require byte equality and equality with the supplied value's `Digest`. Retain the decoded immutable `PolicyView`, not caller-owned internal state.
2. Obtain `inventory.CanonicalBytes()`, strictly decode it with `DecodeInventoryV1(raw, decodedPolicy)`, recompute SHA-256, and require equality with the supplied inventory's bytes, digest, and policy link.
3. Require the inventory's derived `AccountingPolicyDigest` to equal the exact validated policy digest.

This revalidation prevents a zero, stale, internally inconsistent, substituted, or merely digest-shaped typed input from being treated as provenance. It also establishes the only permissible links serialized in Result V1. `InventoryDocumentV1` is never converted to committed entries or line data.

### 3. One accounting foundation, two policy adapters

Introduce a private `accountingDomain` in the root package. It owns only the common classification and accumulation mechanics:

- ordered category names and raw-byte glob slices;
- default-category index and category-name index;
- `classifyAndTotal([]accountingEntry)`, which applies first matching raw-byte glob precedence, uses the default only when no glob matches, produces one classification per supplied entry, and accumulates additions, deletions, and non-countable counts with checked `uint64` addition.

The existing `matchPathGlob`, `addLines`, and raw-byte semantics remain the primitive matcher/overflow behavior. `accountingDomain` does not sort legacy input and does not introduce a new public policy type.

Two private adapters feed the same domain:

- `legacyAccountingDomain(AccountingPolicy)` preserves `validateAccountingPolicy`, its error fields/codes, legacy ordering, and the current public `Account` surface.
- `resultAccountingDomain(PolicyView)` copies Policy V1 category names and raw glob bytes directly, preserves Policy V1 category order/default, and retains its uint64 thresholds and exact ratio-decimal strings alongside the domain for Result-only observations.

`Account` will delegate classification and category-total accumulation to this domain, then retain its present legacy observation conversion, `float64` comparison, return types, ordering, and errors. Result V1 calls the same domain and builds its per-entry evidence from the returned classifications. Thus precedence, default choice, countable treatment, and category totals have one implementation foundation rather than two drifting classifiers/totallers.

Result V1 does **not** adapt Policy V1 through public `AccountingPolicy`: that would lose Policy V1's uint64 threshold range and exact decimal wire value. The private Result adapter is a lossless view-to-domain projection only.

### 4. Snapshot and path boundary

`validateResultSnapshot(snapshot)` copies `snapshot.Entries()` and validates before classification:

- base and head are each exactly 40 or 64 lowercase hexadecimal bytes and use the same width;
- equal base/head is accepted only for an empty entry list;
- every current `CommittedChange.Path` is converted with `[]byte(change.Path)`, preserving all Go-string bytes exactly;
- each path is nonempty, relative, slash separated, lacks NUL, empty, `.` and `..` components, and is not normalized, case-folded, UTF-8 checked, locale sorted, separator-rewritten, or cleaned;
- paths are byte-identical unique; and
- only `Path` and `Lines` feed accounting. `PreviousPath`, status, object IDs, modes, kinds, and binary metadata are deliberately retained outside the Result V1 evidence calculation.

The builder retains only derived, copied evidence. This validates a supplied snapshot boundary but does not claim that it was acquired from Git, belongs to a repository, has existing objects, or was stable under a moving reference.

### 5. Totals and observations

The Result builder obtains classifications and category totals once from `accountingDomain`, then sorts the entry view by `bytes.Compare(Path)` for canonical output. Category totals remain in Policy V1 category order, including all-zero categories.

For a non-countable source entry, the result measurement is exactly `{"kind":"non_countable"}` and no input addition/deletion value is used. It increments only its category's `non_countable` count. Countable entries serialize their exact additions and deletions.

Result-specific observation derivation has a checked `resultLineTotal` helper:

- `non_countable != 0` makes a measurement unavailable;
- otherwise additions plus deletions must be computed with checked `uint64` addition;
- an overflow in a stored category sum, non-countable count, threshold actual, or ratio numerator/denominator is fatal, returning `ResultDocumentV1{}`.

Line-threshold observations are emitted first for nonzero Policy V1 thresholds in category order. Available observations contain `actual` and use strict `actual > limit`; unavailable observations omit `actual` and set `exceeded:false`.

Ratio observations follow in declared Policy V1 ratio order. They are available only when both category line totals are available and the reference total is nonzero. Their numerator and denominator are the exact line totals. `exactRatioExceeded(numerator, denominator, policyDecimal)` compares

```text
numerator / denominator > exact-policy-decimal
```

with `math/big.Int` cross multiplication (including the decimal scale), never `float64` or a reformatted decimal. The output `limit` is the exact canonical Policy V1 decimal string. This supports the required `9007199254740993 / 18014398509481984 > 0.5` vector.

The shared compatibility domain is intentionally tested only where legacy `float64` and exact rational comparison have the same mathematical result. The one intended difference for an available ratio is the documented floating-point precision boundary. `Account` remains untouched externally and retains its current float64 behavior. Separately, Result V1's specified fatal required-line-total overflow is enforced before a ratio comparison; legacy `Account` retains its established unavailable-observation behavior for that legacy-only edge. This is not a changed legacy API and is tested explicitly as a Result boundary rather than represented as a ratio-comparison equivalence vector.

Non-countable entries, unavailable observations, and `exceeded:true` are embedded factual evidence. Antecedent mismatch, invalid snapshot/revision/path, inconsistent decoded evidence, and every required arithmetic overflow are fatal `ContractError` failures with no partial result.

### 6. Canonical wire encoder and identity

Private, declaration-ordered wire structs encode exactly the specified seven root fields:

```json
{"schema":"git-change-evidence.accounting-result/v1","accounting_policy_sha256":"…","inventory_sha256":"…","revisions":{"base":"…","head":"…"},"entries":[…],"totals":[…],"observations":[…]}
```

The encoder uses `encoding/json` only over concrete structs and concrete per-variant observation structs (marshaled into validated `json.RawMessage` elements), never maps. This fixes root, revision, entry, measurement, total, and observation field order. It uses `base64.StdEncoding.EncodeToString` for padded RFC 4648 standard Base64 path values, emits non-nil slices as `[]`, appends exactly one LF, and derives `DocumentDigest` as lowercase-hex SHA-256 over those final bytes. The digest is never a wire field or constructor input.

Construction normalizes only equivalent caller entry order by raw byte sort after validation. It does not normalize any path byte, policy string, decimal, observation, or antecedent. Any valid evidence-byte change, including antecedent link, revision, path, classification, measurement, total, observation, or order, changes canonical bytes and therefore identity.

### 7. Duplicate-aware strict Result decoding

`DecodeAccountingResultV1` first runs `validateResultAntecedents`, rejects non-UTF-8 input, then parses every object level with a Result-specific wrapper around the existing token-based `decodePolicyObject` helper. The wrapper:

- requires every field and allows only the appropriate closed field set;
- catches duplicate decoded keys, including escaped-key duplicates, before `json.Unmarshal` can collapse them;
- recognizes authority-bearing unknown key names after case/separator normalization (`approval`, `approve`, `rejection`, `reject`, `denial`, `deny`, `blocking`, `block`, `gating`, `gate`, `merge`, `deploy`, `release`, `delivery`, `deliver`, and `delivery_authority` variants); and
- applies this check only to keys and discriminators. Exact Policy-derived values in permitted `category`, `reference`, and `limit` fields remain valid evidence even when they contain words such as `merge` or `release`.

The shape parser validates root/revisions/entries/measurements/totals/observations, exact discriminator variants, primitive JSON types, lower-hex digests and revisions, padded standard Base64, raw result-path validity, unsigned integer lexical form and range, and arrays rather than `null`. It accepts no trailing JSON value.

After structural parsing, decoding reconstructs a minimal copied `CommittedSnapshot` from the candidate revisions and entry path/measurement data, invokes the same Result builder, and requires byte-for-byte equality between the candidate and rebuilt canonical bytes. That single equality check additionally rejects alternate string escapes, whitespace, field order, array order, omitted/extra observations, duplicate paths, incorrect classification, totals, availability, actuals, ratio values, and exceeded flags. It repairs nothing and returns the zero result on every failure.

### 8. Narrow Report V1 provenance binding

Keep `NewReportV1`, `DecodeCanonical`, `ReportInput`, the report wire object, projections, CLI, and publication unchanged for legacy/standalone behavior. Add an antecedent-carrying binding type and two additive APIs in `contract.go`:

```go
type ResultV1ReportBinding struct {
    Policy    PolicyDocumentV1
    Inventory InventoryDocumentV1
    Result    ResultDocumentV1
}

func NewReportV1FromResult(subject string, binding ResultV1ReportBinding) (Evidence, error)
func ValidateReportV1ResultProvenance(report Evidence, binding ResultV1ReportBinding) error
```

Both APIs strictly validate the exact Policy → Inventory → Result chain by canonical Result decode/revalidation. `NewReportV1FromResult` derives the Report V1 `Provenance` from the validated result: accounting digest from Result bytes, policy/inventory links from Result's exact antecedents, and revisions from Result. It then calls the unchanged `NewReportV1`, so Report V1 bytes, field order, schema, projections, and publication behavior do not change.

`ValidateReportV1ResultProvenance` first canonical-validates the supplied `Evidence`, validates the binding, then compares all five report provenance values to the derived values. A legacy report made with well-formed but different digests fails this validator. Conversely, `DecodeCanonical` alone remains a syntactic/canonical Report V1 parser and must not claim Result V1 binding proof because it has no exact Result and antecedents.

## Data flow

### Construction

```text
PolicyDocumentV1 + InventoryDocumentV1
  -> strict canonical revalidation and exact policy->inventory link
  -> immutable PolicyView + derived links
CommittedSnapshot
  -> revision/path/duplicate validation; copied accounting entries
  -> shared accountingDomain classification and checked category totals
  -> raw-byte entry sort + Policy-order totals
  -> exact Result observations
  -> concrete ordered wire structs -> canonical JSON + LF -> SHA-256
  -> immutable ResultDocumentV1
```

Inventory contributes only the validated `inventory_sha256` link. Its paths, content hashes, lengths, and entry count do not enter the accounting branch.

### Strict decode

```text
candidate bytes + exact policy + exact inventory
  -> antecedent revalidation
  -> UTF-8/closed-shape/duplicate-aware parse
  -> primitive, path, link, and variant validation
  -> reconstruct minimal snapshot from candidate entries
  -> shared builder recomputation
  -> byte-for-byte canonical equality
  -> immutable ResultDocumentV1 or exact zero value
```

### Bound report creation/validation

```text
exact result + exact antecedents -> strict Result validation
  -> derive Report V1 provenance -> unchanged NewReportV1
existing Evidence + exact binding -> canonical report validation + derived provenance comparison
```

## Validation boundaries and errors

| Boundary | Accepted input | Fatal failure | Nonfatal evidence |
|---|---|---|---|
| Antecedents | exact canonical Policy V1 and linked Inventory V1 values | zero/noncanonical/substituted document, digest mismatch, wrong inventory policy link | none |
| Snapshot | immutable IDs and copied committed entries | malformed/cross-width IDs, equal revisions with entries, invalid/duplicate current path | none |
| Accounting | one classification for each committed entry | checked accumulation or required line-total overflow | explicit non-countable classification |
| Observations | Policy V1 thresholds and ratios only | required arithmetic overflow or inconsistent decoded observation | unavailable threshold/ratio; exceeded factual value |
| Result wire | one exact canonical closed JSON document | invalid UTF-8/JSON/type/key/variant/value/order/link/recomputation | none |
| Report binding | exact Result plus antecedents and canonical Evidence | any derived report provenance mismatch | standalone report parsing remains unbound, not failed evidence |

Use `ContractError` consistently for new Result and binding validation, with stable field paths such as `accounting_policy`, `inventory`, `snapshot.base`, `snapshot.head`, `entries.path_b64`, `totals`, `observations`, `document`, and `provenance.*`. Do not add a generic Result diagnostics array or expose internal parse state.

## Requirement traceability

| Requirement source | Design mechanism | Verification focus |
|---|---|---|
| `result-v1-contract`: typed construction, zero-on-error, ownership | `ResultDocumentV1`, exact constructor inputs, private state and cloning accessors | zero/error, caller/view mutation, no digest-only or legacy-policy construction |
| `result-v1-contract`: exact wire and identity | ordered wire structs, padded Base64, final LF, SHA-256 | hostile-path golden bytes, empty arrays, digest vectors |
| `result-v1-contract`: strict decode | recursive token-based closed parser plus rebuild/equality | duplicate/escaped keys, unknown/authority fields, alternate encodings/order |
| `result-v1-accounting-evidence`: snapshot scope | `validateResultSnapshot` and copied current-path entries | SHA-1/SHA-256 IDs, equal revisions, hostile/invalid paths, rename behavior |
| `result-v1-accounting-evidence`: exclusive/recomputable totals | shared `accountingDomain`; policy order and raw-byte sort | overlap/default/non-countable/duplicate/overflow/order vectors |
| `result-v1-accounting-evidence`: observations | checked line totals and `big.Int` exact rational comparison | threshold/ratio order, unavailable variants, strict greater-than, precision edge |
| `result-v1-accounting-evidence`: legacy compatibility | legacy adapter uses same domain; legacy observations remain float64 | shared-domain differential vectors and legacy-only preservation vectors |
| `result-v1-provenance` | antecedent revalidation; links inside Result canonical body | same/different antecedents, wrong inventory link, any evidence-byte identity change |
| `evidence-v1-provenance` | `ResultV1ReportBinding`, derived builder, validator | all report links/revisions derived, mismatched well-formed digest rejected, standalone decode not proof |

## Alternatives and tradeoffs

1. **Call `Account` from Result V1. Rejected.** It loses per-entry classifications, uses the wrong public policy type, cannot retain exact Policy V1 decimals/uint64 thresholds, and hardcodes float64 ratio semantics.
2. **Implement a separate Result classifier/totaller. Rejected.** It invites drift from first-match/default/raw-byte semantics. A private shared domain retains one foundation while preserving `Account`'s public API.
3. **Use UTF-8 JSON paths directly. Rejected.** Go strings can contain hostile bytes; Base64 preserves those bytes and avoids normalization assumptions.
4. **Use maps or general JSON re-marshaling for the result. Rejected.** Neither fixes field order nor detects duplicate keys before Go's JSON decoder collapses them.
5. **Treat non-countable/unavailable as errors or a generic diagnostics list. Rejected.** The specifications require visible closed evidence variants and no authority-bearing status vocabulary.
6. **Change `NewReportV1` or `DecodeCanonical` to require Result antecedents. Rejected.** That breaks existing constructors/decoders and existing report consumers. An additive bound builder/validator establishes the new proof boundary.
7. **Add snapshot acquisition or inventory content to the result. Rejected.** It violates the pure functional-core boundary and the explicit provenance-only Inventory V1 requirement.

## Threat model and mitigations

| Threat | Mitigation | Explicit non-claim |
|---|---|---|
| Tampered, duplicate, escaped, or ambiguous JSON keys | token-level duplicate-aware closed parsing and canonical rebuild equality | no acceptance of repairable alternatives |
| Digest substitution across valid documents | exact canonical antecedent decode and complete link comparison | SHA-256 identity does not prove acquisition or object existence |
| Hostile/non-UTF-8 path bytes | byte-preserving Go-string conversion, byte validation/sort, padded Base64 | no path normalization or filesystem interpretation |
| Alias mutation | copies at every ingress/egress and private document state | no claim callers cannot mutate their own inputs after return |
| Numeric wrap/fictional measurement | checked uint64 accumulation and `math/big.Int` ratio comparison | no saturation, estimation, or float fallback |
| Authority vocabulary injection | closed key/discriminator validation with normalized authority-key recognition | exact permitted policy values are not censored |
| Semantic drift from legacy accounting | shared classification/totals domain and differential tests | legacy float64 precision behavior remains legacy behavior |
| Resource exhaustion from huge untrusted documents/decimals | pure parser/encoder has no I/O and allocates only from supplied values | no size/quota policy is invented here; shell-level budgets are separate work |

## Test strategy (strict TDD in the later apply phase)

All Go changes follow RED → GREEN → TRIANGULATE → REFACTOR. Add focused failing tests before each production unit, rerun focused tests after refactoring, then run the configured `go test ./...` and check-only `gofmt` gate. No concurrency-bearing behavior is added, so `go test -race ./...` is not required by this slice.

1. **Contract/ownership golden tests:** exact hostile-path JSON and digest; empty snapshot totals; defensive copying at constructor and all views; exact antecedent links; inventory-only identity change; failed values expose nothing.
2. **Accounting/domain tests:** first-match, default, rename/current-path-only, raw byte ordering, every zero total, non-countable records, threshold and ratio observation order/variants, all checked overflow locations, and no inventory contribution.
3. **Differential compatibility tests:** build equivalent legacy and Policy V1 policies. Compare categories/totals, availability, thresholds, ordinary ratio vectors, and per-field accumulation overflow; preserve legacy caller behavior. Isolate the required float precision vector as the only available-ratio comparator difference, and separately assert Result's required-line-total overflow failure without altering legacy `Account`.
4. **Decoder adversarial tests:** non-UTF-8, trailing values, null/wrong types, missing/unknown/duplicate/escaped-duplicate fields at every depth, authority-shaped keys/discriminators, invalid base64/digest/revision/uint/bool, invalid and duplicate paths, all ordering mutations, bad links, noncanonical escape/whitespace/LF variants, and each recomputation mismatch. Every failure must return the zero Result value.
5. **Report binding tests:** derived provenance and byte-identical Report V1 output; one differing but well-formed policy/inventory/accounting/revision value fails binding validation; invalid Result/antecedents fail; `NewReportV1` and `DecodeCanonical` retain their legacy syntactic behavior; projection/CLI/publication golden tests remain unchanged.

## Likely file plan and layout controls

| File | Change | Responsibility |
|---|---|---|
| `accounting_domain.go` | add, private root helper | shared raw-byte classification and checked category accumulation |
| `accounting.go` | narrow refactor | adapt legacy policy to the shared domain while preserving `Account` API/observations |
| `result_v1.go` | add | exported Result V1 types, constants, accessors, cloning |
| `result_v1_build.go` | add | antecedent/snapshot validation, Policy V1 adapter, construction, observations |
| `result_v1_wire.go` | add | concrete canonical wire variants, JSON+LF encoder, digest |
| `result_v1_decode.go` | add | duplicate-aware strict Result parser and recomputation validation |
| `contract.go` | narrow additive update | `ResultV1ReportBinding` and bound report builder/validator |
| `result_v1_contract_test.go` | add | API, canonical, ownership, provenance, and report-binding tests |
| `result_v1_accounting_test.go` | add | domain/observation/compatibility/overflow/path tests |
| `result_v1_decode_test.go` | add | strict decode adversarial tests |
| existing co-located tests as needed | narrow additions only | regression assertions for unchanged legacy behavior |

No change is planned for Policy V1, Inventory V1, `snapshot.go`, Report V1 wire decoding/projection, `internal/`, `cmd/`, publication, or historical bootstrap artifacts.

The root package already exceeds the configured 12-public-root-Go-file baseline. This design also likely touches more than the configured eight Go files in that package because strict decoder and adversarial tests must remain separately readable. Apply must therefore record the first applicable threshold crossing in `apply-progress.md` and obtain the explicit package-boundary review required by `openspec/config.yaml` before adding files beyond that threshold. Keep each new or materially edited Go file under the 160 changed-line threshold by splitting only along the responsibilities above; do not evade the review by creating `pkg/`, a separate test tree, or a package-level guide.

## Migration and rollback

The rollout is additive: callers that need canonical accounting construct Policy V1 and linked Inventory V1, supply an already acquired committed snapshot, and retain the Result V1 bytes/digest with the exact antecedents. Callers that need a Report V1 provenance proof use `NewReportV1FromResult` and retain the binding material for later `ValidateReportV1ResultProvenance`. Existing callers may continue using `Account`, `NewReportV1`, and `DecodeCanonical`; standalone Report V1 parsing remains syntactic rather than Result-bound proof.

No data migration, configuration migration, CLI migration, publication change, or consumer cutover is part of this change. Rollback removes the Result V1 files, bound report APIs/tests, and only behavior-preserving accounting-domain refactor. Existing Policy/Inventory documents, legacy accounting, report wire bytes, decoders, projections, CLI, and publication remain usable; no stored artifact is rewritten.

## Unresolved architecture decisions

None remain at the Result V1 product-contract level. The specifications fix the wire shape, antecedent role, accounting scope, ordering, fatal-versus-embedded outcomes, exact ratio semantics, and narrow Report V1 proof boundary. Implementation may select only private helper/error-code organization consistent with the mechanisms above; it must not introduce a new diagnostic vocabulary, policy/configuration format, acquisition behavior, or delivery semantics.

## Result contract

- **status:** success
- **artifact:** `openspec/changes/result-v1/design.md`
- **summary:** Designed an immutable canonical Result V1 root contract with strict Policy→Inventory→Result validation, a shared private accounting domain, exact byte paths and ratio arithmetic, duplicate-aware canonical decoding, and additive Report V1 provenance binding while preserving legacy public behavior and report surfaces.
- **next_recommended:** tasks
- **skill_resolution:** injected
- **persistence:** OpenSpec artifact written; no external memory persistence available
- **notes:** CodeGraph/CLI tools were unavailable in this executor environment, so source inspection used targeted filesystem reads after the required project material was read.
