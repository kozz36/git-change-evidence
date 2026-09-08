# Proposal: Result/Accounting V1 Canonical Contract

## Intent

Add **Result/Accounting V1** as the next independent neutral-core provenance unit after **Accounting Policy V1 → Inventory V1**. The new public root-package contract will turn an exact validated `PolicyDocumentV1`, an exact validated and policy-linked `InventoryDocumentV1`, and an immutable committed snapshot into deterministic, auditable accounting evidence with a canonical document and system-derived digest.

The user and product outcome is a reproducible accounting result that authors, reviewers, and automation can independently inspect and recompute: every committed entry has one explicit classification, totals are exclusive and recomputable, non-countable or unavailable measurements remain visible, and the full policy/inventory/result provenance chain is validated rather than asserted.

`docs/product-charter.md` is authoritative for this change. The active-path historical `bootstrap-neutral-core-extraction` artifacts are drifted historical context only; they do not define Result V1 requirements or authority.

## Current-State Gap

The repository already ships the required antecedents and accounting behavior, but not a canonical Result V1 contract:

- Policy V1 provides strict canonical policy documents, raw-byte path globs, ordered categories, thresholds and ratios, and a system-derived document digest.
- Inventory V1 provides the canonical second provenance document, bound to the exact Policy V1 antecedent, with a strict decoder and system-derived digest.
- `Account` provides exclusive first-match/default classification, per-category totals, non-countable counts, threshold and ratio observations, and overflow rejection over committed snapshots. Its public ratio behavior retains its historical `float64` comparison semantics.
- The outer Evidence V1 report carries an `accounting_sha256` value but currently validates only its digest syntax; it does not formally establish that the value identifies a canonical Result V1 document.

The shipped `AccountingResult` is an unversioned in-memory value. It does not preserve explicit per-entry classifications, canonical bytes, a document digest, or the complete validated antecedent chain. Its public input is the separate `AccountingPolicy` shape rather than the exact `PolicyDocumentV1`. Consequently, consumers cannot exchange or strictly validate one canonical, audit-ready accounting result.

## Scope

### In Scope

- A public, neutral, versioned Result V1 document/API in the repository root.
- Construction from all of the following exact typed inputs:
  - an exact validated `PolicyDocumentV1`;
  - an exact validated `InventoryDocumentV1` whose policy link and document digest are verified against that policy; and
  - an immutable committed snapshot with valid immutable base/head revision identities.
- Direct public consumption of `PolicyDocumentV1`. Private implementation code may adapt that document to reuse shipped classification, totals, and availability semantics, while Result V1 applies exact Policy V1 decimal ratio comparison without changing legacy `Account` behavior.
- Accounting of committed snapshot entries only. Inventory V1 is a required provenance antecedent and does not contribute untracked entries, content lengths, or other values to committed accounting totals.
- An immutable canonical result containing explicit per-entry exclusive classifications and per-category exclusive totals that can be recomputed from those entries.
- Exact committed-path byte preservation, with canonical wire representation using Base64 rather than UTF-8 assumptions or path normalization.
- Threshold and ratio evidence governed by exact Policy V1 decimal semantics, while preserving established legacy accounting compatibility for classifications, totals, availability, and ordinary ratio cases.
- A mixed diagnostic model:
  - malformed or mismatched antecedents, invalid snapshot or revision input, and arithmetic overflow are fatal and emit no valid Result V1 document;
  - non-countable values and unavailable threshold or ratio measurements remain embedded evidence observations rather than construction failures.
- Deterministic canonicalization and ordering based on:
  - Policy V1 category order for category-derived result material;
  - raw path-byte order for committed-entry classifications; and
  - one closed contractual order for observation variants.
- A strict validating decoder that requires the exact Policy V1 and Inventory V1 antecedent documents, validates the complete policy → inventory → result chain, rejects non-canonical alternatives, and returns no partial valid document on fatal failure.
- A system-derived Result V1 digest over the exact canonical Result V1 bytes.
- Formal binding of the existing outer report `accounting_sha256` meaning to the digest of the exact canonical Result V1 document, without expanding report projections, CLI output, or publication behavior.
- Focused root-package tests for accounting compatibility, classifications, recomputable totals, path-byte fidelity, observations, fatal failures, provenance linkage, canonicalization, strict decoding, immutable ownership, and report-digest binding.

The detailed Result V1 field names, object layout, diagnostic encoding, and public Go type signatures remain for specification and design. This proposal freezes the product semantics above, not a prematurely detailed wire layout.

### Out of Scope

- Bounded moving-reference transactions, retries, or revalidation; those belong to a separate later OpenSpec change.
- Git invocation or committed-snapshot acquisition.
- Filesystem discovery, untracked-file acquisition, or adding Inventory V1 entries to committed accounting.
- CLI commands, CLI/projection expansion, publication expansion, human-output redesign, or archive/synchronization behavior.
- CNSIC profile integration or cutover, CNSIC vocabulary, compatibility hashes, review semantics, or project-specific thresholds.
- Configuration-file syntax or adapter decoding choices.
- Module/publication identity, packaging, licensing, CI, release automation, signing, or distribution.
- Approval, rejection, gating, blocking, merge, deployment, release, delivery authority, or interpretation of threshold observations as verdicts.
- Changes to the historical/drifted `bootstrap-neutral-core-extraction` artifacts or any claim under their authority.
- Delivery mechanics, issue creation, commit strategy, PR slicing, pushing, or pull-request creation.

## Capabilities

### New Capabilities

- `result-v1-contract`: Construct, inspect, canonically serialize, digest, and strictly decode a neutral Result V1 document from exact validated Policy V1, Inventory V1, and committed-snapshot inputs.
- `result-v1-accounting-evidence`: Preserve explicit exclusive committed-entry classifications, recomputable category totals, non-countable evidence, and deterministic threshold/ratio observations without introducing delivery authority.
- `result-v1-provenance`: Validate the complete antecedent chain and bind Result V1 identity to both the exact policy and exact inventory documents.

### Modified Capabilities

- `evidence-v1-provenance`: Define `accounting_sha256` as the SHA-256 digest of the exact canonical Result V1 document. The outer report wire shape, projections, publication, and CLI remain otherwise unchanged in this work.

The existing `Account` API retains its public API and observable historical behavior. Internal factoring is permitted only when it does not change that behavior; Result V1 exact decimal ratio comparisons intentionally differ only at documented `float64` precision edges.

## Approach

Keep Result V1 in the functional core. Validate the exact Policy V1 and Inventory V1 document values first, independently derive and compare their canonical identities, then validate the committed snapshot and revisions before producing classifications or totals. Classify every committed entry exactly once using Policy V1 precedence and default behavior. Preserve each committed path as exact bytes, order classifications by those bytes, derive category totals from the classifications in policy category order, and assemble observations in the contractually closed order.

Reuse established classification, totals, non-countable, threshold, and availability semantics. A private adapter from `PolicyDocumentV1` to established accounting internals is acceptable, but Result V1 must not expose the legacy `AccountingPolicy` as its public policy boundary or maintain a second divergent classifier/totaller. Result V1 compares every available ratio as an exact rational value against the exact Policy V1 decimal limit; legacy `Account` retains its historical `float64` comparison behavior. Compatibility vectors cover the shared domain where those comparisons have the same mathematical result, and the documented precision-edge divergence is intentional. Non-countable and unavailable values are retained in the canonical evidence; antecedent, snapshot/revision, and overflow failures abort construction and decoding.

The strict decoder will accept candidate canonical bytes plus the exact Policy V1 and Inventory V1 antecedents. It will validate closed shapes and values, provenance links, classifications, recomputed totals and observations, deterministic ordering, and byte-for-byte canonical re-encoding. It will reject repairable-looking alternatives rather than normalize them.

Report integration is deliberately narrow: the existing `accounting_sha256` field gains an exact Result V1 identity meaning and any validating construction path needed to establish that binding, but Result V1 content is not added to report projections, CLI output, or publication in this change.

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| Public root Go API | New | Immutable Result V1 value, constructor, views, canonical bytes, digest, and antecedent-aware strict decoder |
| Accounting core | Reused / possibly refactored | Share existing exclusive classification, totals, non-countable, threshold, ratio-availability, and overflow semantics; keep Result V1 exact decimal ratio comparison distinct without changing `Account` behavior |
| Snapshot boundary | Validated | Consume immutable committed entries and validate revision/path/count semantics required by Result V1 |
| Policy V1 | Antecedent | Consume the exact document directly and preserve category order, raw-byte glob meaning, thresholds, ratios, and digest identity |
| Inventory V1 | Antecedent | Require the exact document, verify its policy link and digest, and keep its untracked entries outside committed totals |
| Evidence V1 report | Narrow semantic binding | Bind `accounting_sha256` to the exact canonical Result V1 digest without wire, projection, CLI, or publication expansion |
| Root contract tests | New / extended | Canonical, provenance, strict-decoding, accounting-compatibility, path-byte, observation, fatal-error, and report-binding vectors |
| Shell and adapters | Unchanged | No Git, filesystem, CLI, publication, configuration, profile, or moving-reference work |

Likely implementation areas are additive `result_v1*` root-package files and co-located tests, with only narrow behavior-preserving updates to private accounting helpers and report contract validation where required by the formal digest binding.

## Risks and Mitigations

| Risk | Likelihood | Mitigation |
|---|---|---|
| Result V1 duplicates and later diverges from shipped `Account` classification, totals, or availability semantics | Medium | Route shared semantics through private helpers and retain compatibility vectors for the legacy API; keep the exact Policy V1 rational ratio comparator distinct so `Account` retains its historical `float64` behavior |
| Policy adaptation loses raw-byte glob semantics, order, decimal precision, or canonical identity | Medium | Require `PolicyDocumentV1` publicly, derive from its validated immutable view only, and test hostile bytes, exact rational ratio behavior, and the intentional `9007199254740993 / 18014398509481984 > 0.5` legacy precision-edge divergence |
| Untracked Inventory V1 values leak into committed totals | Medium | Treat Inventory V1 only as a validated provenance antecedent and test that changing inventory entries changes provenance identity but not committed accounting totals |
| Per-entry evidence and totals disagree | Medium | Recompute totals and observations from classifications during construction and strict decode; reject discrepancies |
| String handling corrupts hostile committed paths | High | Establish exact-byte input semantics, encode paths canonically as Base64, order by raw bytes, and add non-UTF-8/control-byte fixtures |
| Canonical bytes vary because ordering is underspecified | Medium | Contractually fix policy category order, raw path-byte order, and a closed observation order; require canonical re-encoding equality |
| Fatal and embedded conditions are confused | Medium | Specify the mixed model explicitly and test no-document fatal failures separately from embedded non-countable/unavailable evidence |
| Report `accounting_sha256` remains a free-form digest assertion | Medium | Require exact Result V1 identity at the validating report boundary and test mismatched but well-formed digests |
| Threshold terminology is interpreted as approval or gating | Medium | Keep observations factual and evidence-only; reject authority-bearing contract fields and avoid verdict semantics |
| Scope expands into acquisition, moving references, CLI, or publication | Medium | Keep all effects and integrations outside this pure contract and defer them to separately authorized changes |

## Dependencies

- The authoritative product and authority boundary in `docs/product-charter.md`.
- Shipped Policy V1 canonical construction/decoding, immutable views, raw-byte globs, category order, threshold/ratio definitions, and document digest.
- Canonical Inventory V1 specifications and archived closure evidence, including exact Policy V1 binding, strict canonical decode, raw-byte identity, deterministic ordering, and document digest.
- Shipped committed snapshot types and immutable base/head revision identities.
- Shipped accounting behavior and tests as the compatibility baseline for first-match/default classification, totals, non-countable values, observations, and overflow handling.
- Shipped report construction/decoding and tests as the compatibility baseline for the outer `accounting_sha256` field.

Safe Git/filesystem acquisition and the later bounded moving-reference transaction are downstream integration concerns, not prerequisites for this pure neutral-core contract.

## Rollback Plan

Remove the Result V1 root contract, decoder, and focused tests; revert only any private behavior-preserving accounting factoring and narrow report-binding API/validation introduced for this change. Policy V1, Inventory V1, legacy `Account`, committed snapshot acquisition, report wire shape, CLI, projections, and publication remain independently usable.

Because this change performs no consumer cutover, acquisition integration, publication expansion, or delivery automation, rollback requires no data migration and does not alter the historical/drifted bootstrap artifacts. After rollback, the outer report's `accounting_sha256` returns to its prior syntactically validated but not Result-V1-established state until a replacement contract is approved.

## Success Criteria

- [ ] A public neutral Result V1 API constructs a valid document only from the exact validated `PolicyDocumentV1`, exact validated and correctly policy-linked `InventoryDocumentV1`, and valid immutable committed snapshot/revisions.
- [ ] Inventory V1 is required and its exact digest and policy link are validated, while its untracked entries never contribute to committed classifications or totals.
- [ ] Every committed entry appears exactly once with an explicit classification, and category totals are exclusively and deterministically recomputable from those entries in Policy V1 category order.
- [ ] Committed paths round-trip as exact bytes, use canonical Base64 in the wire document, and are ordered by raw path bytes without UTF-8, Unicode, locale, or path normalization.
- [ ] Result V1 preserves established first-match/default classification, totals, non-countable, threshold, ratio availability, and overflow semantics without changing observable legacy `Account` behavior; available ratios use exact Policy V1 decimal rational comparison, with compatibility required only where legacy `float64` comparison has the same mathematical result and with the intentional `9007199254740993 / 18014398509481984 > 0.5` precision-edge vector covered.
- [ ] Non-countable values and unavailable threshold/ratio measurements are embedded evidence observations; malformed or mismatched antecedents, invalid snapshot/revision input, and overflow produce a fatal error and no valid result document.
- [ ] Equivalent valid inputs produce byte-identical canonical Result V1 documents and the same system-derived digest, including deterministic policy-category, raw-path, and closed observation ordering.
- [ ] Strict decode requires the exact Policy V1 and Inventory V1 antecedents, validates the full provenance chain and recomputed accounting, and rejects malformed, ambiguous, unknown-field, duplicate-field, authority-bearing, mismatched, reordered, or otherwise non-canonical input.
- [ ] The outer report's `accounting_sha256` is formally and testably the digest of the exact canonical Result V1 document, with no CLI, projection, or publication expansion.
- [ ] Result V1 remains evidence-only and contains no CNSIC policy, delivery verdict, Git/filesystem effect, moving-reference transaction, configuration syntax, release/distribution choice, or delivery mechanic.
