# Approved — Neutral committed-snapshot requirements

**Contract approved for implementation under S5–S6, not proof of runtime behavior or consumer adoption.** See [apply progress](../../apply-progress.md) for per-unit executed evidence. See [design](../../design.md) and [plan](../../tasks.md). Every scenario below is prospective and **NOT RUN**.

## ADDED Requirements

### Requirement: NGE-1 — Preserve an independent additive boundary

The proposed export MUST remain independent of policy, inventory, categories and delivery authority. It MUST preserve existing SDK/CLI routing, legacy triples and V1 bytes/semantics. The approved grammar/schema in design MUST NOT be treated as available or verified until implemented and checked; contract approval is separate from runtime evidence.

#### Scenario: Existing consumers retain behavior

- GIVEN existing commands, legacy triples and canonical V1 documents
- WHEN the future additive export is introduced
- THEN their outputs/errors remain unchanged except explicitly approved additive help text.

#### Scenario: Facts do not decide delivery

- GIVEN any successful export or technical failure
- WHEN another system consumes it
- THEN it receives facts/errors, not approval, rejection, gating or delivery authority.

### Requirement: NGE-2 — Acquire an explicit whole immutable range

The candidate stage MUST require the design's exact explicit repository/base/head arguments, resolve immutable native IDs and cover the full range without policy/inventory filtering. It MUST NOT read unused stdin, substitute index/worktree contents, implicitly fetch or follow host symlink targets/gitlinks for side measurement.

#### Scenario: Dirty state and selection do not filter facts

- GIVEN fixed commits with independently changed index/worktree and unrelated untracked paths
- WHEN the range is exported
- THEN committed membership/identities remain unchanged and no inventory selection is consumed.

#### Scenario: Invalid or unavailable acquisition remains failure

- GIVEN malformed arguments, missing commits/required blobs or unsupported capabilities
- WHEN export is attempted
- THEN the appropriate technical failure is explicit, not empty evidence or a fallback acquisition.

### Requirement: NGE-3 — Retain complete typed entry and side facts

Entries MUST preserve native revisions, raw paths, status, modes/kinds, immutable object IDs, valid absent sides and numstat semantics as proposed in design. New validation MUST cover their cross-field constraints rather than relying on accounting-result path validation. Regular-side measurements MUST be complete; unavailable required facts MUST NOT become null/zero success. No similarity-score or expanded copy-detection capability is implied.

#### Scenario: Type, mode, rename and absent sides remain distinguishable

- GIVEN real immutable add/delete, rename, mode-only, type, binary, symlink and gitlink cases
- WHEN the new document is constructed
- THEN valid old/new facts remain distinguishable without replacing missing data with a synthetic file or zero count.

#### Scenario: Inconsistent supplied facts are rejected structurally

- GIVEN incompatible status/side/mode/kind/ID or duplicate effective paths
- WHEN a constructor or decoder validates them
- THEN it rejects the document without claiming Git authenticity or performing repository I/O.

### Requirement: NGE-4 — Define side byte measurements without accounting claims

Required file-side measurements MUST use the proposed explicit byte-length/NUL/LF-line basis in design, with bounded immutable blob reads. They MUST NOT silently become production LOC, binary policy or accounting ownership. Raw numstat tokens MUST be retained before loss, validated against numeric/null values and kept distinct from per-side content metrics.

#### Scenario: Empty, LF, unterminated, CRLF and NUL bytes stay explicit

- GIVEN complete bounded file blobs with those byte forms
- WHEN measurements are acquired
- THEN byte lengths/NUL flags/logical-line values follow the declared basis without Unicode/newline normalization or automatic accounting.

#### Scenario: Git binary numstat is not a side-length substitute

- GIVEN binary numstat or different content characteristics on the two sides
- WHEN facts are exported
- THEN entry numstat and each file-side metric remain separate and no side length is inferred from additions/deletions.

### Requirement: NGE-5 — Preserve canonical byte identity

The candidate schema MUST use the closed ordered objects, Base64 raw paths, deterministic effective-path ordering and exact JSON-plus-LF form in design. Document digests MUST be SHA-256 of those bytes, separate from native Git IDs. Constructor/decoder results MUST be immutable; malformed or noncanonical input MUST NOT be normalized into accepted evidence.

#### Scenario: Equivalent valid inputs produce one canonical document

- GIVEN the same validated facts in different input ordering
- WHEN canonical assembly is performed
- THEN the exact declared bytes/digest match and caller mutation cannot alter the returned document.

#### Scenario: Malformed encodings remain errors

- GIVEN duplicate/unknown keys, invalid Base64, unsafe paths, wrong numeric/token forms or extra trailing data
- WHEN decoding is attempted
- THEN failure is explicit rather than silently repaired.

### Requirement: NGE-6 — Bound complete acquisition and output

The shell MUST enforce the proposed individual/aggregate/entry/output/deadline ceilings across all acquisition phases, including new side reads. It MUST preserve explicit unavailable/resource/interruption/inconsistency/incomplete-output outcomes and MUST NOT silently raise bounds, truncate, clamp or return partial success. Linux restriction and cooperative limits MUST NOT be presented as all-platform or hard resource/atomicity certificates.

#### Scenario: A later phase exhausts the shared budget

- GIVEN individually admissible responses whose combined bytes exceed admission
- WHEN metadata plus side reads execute
- THEN acquisition fails explicitly rather than resetting phase counters or emitting success.

#### Scenario: Output delivery or cancellation fails

- GIVEN interrupted acquisition or incomplete stdout delivery
- WHEN the operation ends
- THEN its actual technical outcome is retained and partial bytes are not accepted as evidence.

### Requirement: NGE-7 — Separate association from report provenance

A proposed association check MUST compare native revisions and complete effective-path membership without altering Result/Report V1. Exact new-document bytes/digest MUST be retained separately. Association MUST NOT claim authentication, semantic count/category equality, consumer compatibility or existing-report coverage of the new antecedent.

#### Scenario: Range mismatch is detectable

- GIVEN valid documents with different revisions or committed path membership
- WHEN association is checked
- THEN mismatch is reported without rewriting any antecedent or reporting provenance success.

#### Scenario: Matching range does not imply covered measurements

- GIVEN matching revisions/membership and an existing V1 report
- WHEN the new facts document is retained beside it
- THEN the report retains its old binding only and the caller explicitly identifies the additional evidence and its limitations.

### Requirement: NGE-8 — Keep design approval separate from delivery

Implementation MUST wait for concrete contract/domain approval, frozen vectors, scoped layout review and a calibrated work-unit plan. Distribution, pin, release and consumer adoption MUST remain separate owner decisions. This proposal and its synthetic example MUST NOT be presented as runtime or parity evidence.

#### Scenario: A documented candidate is not an available capability

- GIVEN only this Markdown proposal
- WHEN availability or adoption is evaluated
- THEN implementation, acquisition, release and consumer parity remain unclaimed.

#### Scenario: Compatibility gaps survive export delivery

- GIVEN a future available export but unresolved consumer policy/annotation/error/byte/domain mapping
- WHEN a consumer considers adopting it
- THEN those independent prerequisites remain its responsibility and no automatic adoption is inferred.
