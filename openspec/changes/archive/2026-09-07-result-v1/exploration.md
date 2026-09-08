# Exploration: Result/Accounting V1

## Authority and current state

`docs/product-charter.md` is authoritative. The product is a neutral, evidence-only Git evidence library and CLI. It may report facts, diagnostics, warnings, and unavailable measurements, but must not approve, reject, gate, block, merge, deploy, release, or assign delivery authority.

This worktree was created from clean `origin/main=a2b9654`. The investigation boundary is `origin/main=a2b9654`, and the only intended delta is this exploration artifact. The active-path historical/drifted `bootstrap-neutral-core-extraction` change remains context only and is not current Result V1 requirements.

The worktree already contains:

- Go 1.25.10 module `github.com/kozz36/git-change-evidence`.
- Public neutral APIs at the repository root.
- `internal/` shell packages and `cmd/git-change-evidence/`.
- Merged Policy V1 and Inventory V1 contracts.
- Shipped accounting and report behavior.

## Source-of-truth hierarchy

1. `docs/product-charter.md` — product boundary and authority model.
2. Canonical OpenSpec Inventory V1 specifications and archived closure evidence.
3. Shipped Policy V1, Inventory V1, accounting, snapshot, and report implementations/tests.
4. Active-path historical/drifted `bootstrap-neutral-core-extraction` artifacts — historical planning and predecessor context only.
5. No authoritative Result V1 specification currently exists.

The archived Inventory V1 contract identifies the intended sequence:

```text
Accounting Policy V1 → Inventory V1 → Accounting Result V1
```

It also explicitly leaves Result V1 outside Inventory V1.

## Existing Policy V1

Policy V1 is implemented in `policy_contract.go`, `policy_decode.go`, and their tests. It provides:

- `PolicyDocumentV1`;
- `NewAccountingPolicyV1` and `DecodeAccountingPolicyV1`;
- canonical bytes and a system-derived SHA-256 document digest;
- immutable/defensive-copy views;
- ordered categories and raw-byte path globs;
- a default category and decimal ratio limits;
- strict closed-shape decoding and canonical-byte validation.

Policy construction normalizes caller input for canonical serialization; decoding requires the already canonical representation. Policy identity is the SHA-256 digest of exact canonical bytes.

Result V1 must consume the exact validated Policy V1 document, not a caller-supplied policy digest or an unvalidated `AccountingPolicy`.

## Existing Inventory V1

Inventory V1 is the canonical second provenance document. It is implemented in `inventory_v1.go`, `inventory_v1_build.go`, `inventory_v1_decode.go`, and their tests. It provides:

- `InventoryDocumentV1`;
- `NewInventoryV1` and `DecodeInventoryV1`;
- canonical bytes and document digest;
- derived `accounting_policy_sha256`;
- raw path-byte preservation;
- exact content SHA-256 and byte length;
- deterministic raw-byte path ordering;
- strict decoding and defensive ownership.

Inventory V1 is an untracked/content inventory. It does not perform filesystem discovery, Git acquisition, line accounting, or Result V1 assembly. Empty inventory means a verified empty scope, not unavailable acquisition.

Result V1 must verify that the supplied Inventory V1 is valid and linked to the exact supplied Policy V1 document.

## Existing accounting behavior

`accounting.go` and `accounting_test.go` provide shipped in-memory accounting behavior:

- `AccountingPolicy` and `Account`;
- exclusive first-matching-category classification;
- configured default-category fallback;
- additions, deletions, and non-countable totals;
- line-threshold and ratio observations;
- unavailable observations when non-countable data prevents measurement;
- invalid-policy rejection before totals are emitted;
- overflow rejection with an empty result;
- raw-byte glob matching through string-to-byte conversion.

The current result is not a canonical versioned document and has no document digest. It also consumes the separate `AccountingPolicy` shape rather than `PolicyDocumentV1`.

The current accounting implementation is therefore shipped implementation evidence, not the Result V1 contract. The smallest safe boundary appears additive: preserve `Account` behavior and tests, then consider a canonical Result V1 constructor/decoder that uses validated Policy V1 and existing immutable snapshot semantics. This is an exploratory recommendation only; no implementation contract is frozen.

## Existing snapshot behavior


`snapshot.go` defines `CommittedSnapshot`, `CommittedChange`, `CommittedLineCounts`, revision identifiers, file/status/path/object/mode/kind metadata, and countable versus non-countable line measurements.

Result V1 should account committed snapshot entries only. It must not infer committed accounting from index or working-tree state, and it must not merge untracked Inventory V1 entries into committed line totals unless an explicit product decision defines that behavior.

## Existing report behavior

`contract.go`, `decode.go`, `report.go`, and their tests provide the shipped `git-change-evidence.evidence/v1` report contract with subject, base/head revisions, input-policy/inventory/accounting digests, canonical JSON, strict decoding, projections, and sanitized errors.

`ReportInput` currently carries only provenance digests and revisions; it does not carry accounting totals, classifications, or diagnostics. Because CLI/publication expansion is outside Result V1, this change should not expand projections or publication. The existing report can remain an outer evidence envelope whose `accounting_sha256` may point to the Result V1 document digest, subject to an explicit compatibility decision.

## Smallest coherent Result V1 boundary

The following is an exploratory boundary for proposal consideration, not an approved implementation contract. No later SDD phase has run, and all MUST-level requirements remain subject to proposal and spec approval.

Result V1 should be considered as a pure public-root contract that:

1. Accepts an exact validated `PolicyDocumentV1`.
2. Accepts an exact validated `InventoryDocumentV1` linked to that policy.
3. Accepts an immutable committed snapshot or equivalent typed committed-change input.
4. Applies exclusive policy classification exactly once.
5. Emits canonical totals and diagnostic observations.
6. Preserves explicit non-countable/unavailable measurements.
7. Includes immutable revision identity.
8. Derives a Result V1 document digest from canonical bytes.
9. Provides strict canonical decoding requiring the exact Policy V1 and Inventory V1 antecedents.
10. Contains no Git, filesystem, CLI, publication, moving-reference, CNSIC, or delivery-authority behavior.

The likely provenance chain is:

```text
Policy V1 canonical bytes
  └─ policy digest
      └─ Inventory V1, linked to policy digest
          └─ inventory digest
              └─ Result V1, linked to policy digest + inventory digest
                  └─ optional outer Evidence/report envelope
```

## Candidate Result V1 contents

The exact wire shape is not authoritative yet. A coherent minimal shape would include:

- schema/version;
- policy digest;
- inventory digest;
- base/head revisions;
- per-category exclusive totals;
- non-countable counts;
- threshold and ratio observations;
- diagnostics;
- optionally, per-entry classifications when canonical auditability requires them;
- a Result V1 document digest outside the wire body.

Totals-only is smaller. Per-entry classification provides stronger auditability and makes exclusivity directly inspectable, but increases the canonical document and requires decisions about path/status/object metadata representation.

## Diagnostics

Result V1 should distinguish at least:

- invalid policy;
- invalid or mismatched inventory antecedent;
- invalid snapshot/revision input;
- overflow;
- non-countable line measurements;
- unavailable threshold or ratio measurements;
- threshold or ratio exceeded observations.

“Exceeded” must remain an observation, never a verdict. Malformed antecedents and arithmetic overflow are candidate fatal errors; non-countable measurements and unavailable thresholds are candidate non-fatal evidence records.

## Open product decisions

1. Is Inventory V1 only a required provenance antecedent, or does Result V1 account inventory entries alongside committed changes?
2. Must Result V1 consume `PolicyDocumentV1` directly, or may a typed adapter translate it into existing `AccountingPolicy`?
3. Does the canonical result contain totals only or explicit per-entry classifications?
4. What exact revision and committed-entry metadata belongs in the Result V1 wire document?
5. Which conditions are fatal construction errors versus embedded diagnostics?
6. Does Result V1 repeat both policy and inventory digests and independently verify Inventory V1’s policy link?
7. Is existing `Provenance.AccountingDigest` formally bound to the Result V1 digest, or is report integration deferred?
8. Do threshold and ratio observations remain unavailable when any relevant category has non-countable entries, matching current behavior?
9. What is the canonical category, entry, observation, and diagnostic ordering?
10. Does strict decoding require exact Policy V1 and Inventory V1 antecedent documents?
11. Must committed paths become raw bytes before Result V1, since `CommittedChange.Path` is currently a Go string?
12. Does legacy `Account` remain unchanged with additive Result V1 behavior, or may internals be reorganized without observable change?

## Non-goals

This change must not include:

- bounded Git moving-reference transactions or retry/revalidation;
- Git acquisition changes;
- filesystem or untracked-inventory acquisition;
- CLI or publication expansion;
- human projection redesign;
- CNSIC integration or profile cutover;
- CNSIC vocabulary, thresholds, compatibility hashes, or review semantics;
- changes to legacy `Account` behavior except a behavior-preserving internal refactor required for reuse;
- approval, rejection, gating, blocking, merge, deployment, release, or delivery authority;
- package/release/CI/licensing/distribution decisions;
- concrete external configuration syntax;
- changes to the historical `bootstrap-neutral-core-extraction` artifacts or authority.

## Likely affected files

Likely additions:

- `result_v1.go`
- `result_v1_build.go`
- `result_v1_decode.go`
- `result_v1_contract_test.go`
- `result_v1_decode_test.go`

Possible narrow updates:

- `accounting.go` or private accounting helpers, only to share behavior without changing `Account`;
- `contract.go`/`decode.go` only if Result V1 is explicitly integrated into the outer report contract.

Expected non-changes:

- `internal/git/`
- `internal/inventory/`
- `internal/publication/`
- `cmd/git-change-evidence/`
- legacy inventory APIs
- historical bootstrap artifacts

## Risks

- **Semantic duplication:** a second accounting algorithm could diverge from shipped `Account`.
- **Policy conversion drift:** translating Policy V1 could lose raw-byte glob semantics or canonical identity.
- **Inventory misinterpretation:** treating untracked inventory as committed line-accounting input could violate separation.
- **Insufficient auditability:** totals alone may not prove each entry was classified exactly once.
- **Canonical instability:** unspecified ordering or diagnostics could produce different bytes.
- **Authority leakage:** exceeded/status vocabulary could be interpreted as a gate.
- **Report ambiguity:** reusing `accounting_sha256` without defining its meaning could weaken provenance.
- **Scope expansion:** CLI, publication, or moving-reference integration would combine multiple work units.

## Proposal readiness

**Conditionally ready for proposal consideration.** The product boundary and provenance chain appear clear, but this exploration does not approve implementation or freeze requirements. A later proposal would still require owner decisions on wire shape, inventory role, diagnostic model, direct Policy V1 integration, per-entry auditability, and report linkage. No later phase has run, and all MUST-level requirements await proposal/spec approval.

## Result Contract

- **status:** success
- **executive_summary:** Exploration identifies a possible pure additive Result V1 document consuming validated Policy V1, validated Inventory V1, and an immutable committed snapshot; performing evidence-only accounting; emitting totals and diagnostics; and deriving a canonical digest. These are exploratory findings and recommendations only: no implementation contract is frozen, no later SDD phase has run, and all MUST-level requirements await proposal/spec approval.
- **artifacts:** `openspec/changes/result-v1/exploration.md`
- **next_recommended:** Conditionally proceed to `sdd-propose` after owner decisions on wire shape, inventory role, diagnostics, and report linkage; this recommendation does not constitute proposal or implementation approval.
- **risks:** Accounting duplication or drift, Policy V1 adaptation loss, inventory/committed-scope confusion, insufficient auditability, canonical instability, authority leakage, report provenance ambiguity, and scope expansion.
- **skill_resolution:** `paths-injected`
