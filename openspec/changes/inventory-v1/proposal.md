# Proposal: Inventory V1 Canonical Contract

## Intent

Add the neutral **Untracked Inventory V1** canonical document and public Go API as the second provenance document in the sequence **Accounting Policy V1 → Inventory V1 → Accounting Result V1**. The contract will preserve raw path identity, derive each entry's SHA-256 content digest and exact byte length from content, and bind the inventory to the exact antecedent Accounting Policy V1 document.

This change is additive. The shipped `UntrackedRecord` and `UntrackedInventory` API and its acquisition-facing semantics remain unchanged.

## Scope

### In Scope

- A public root-package Inventory V1 value/document API with immutable ownership and defensive copies.
- A neutral, versioned canonical JSON document containing:
  - the Inventory V1 schema identity;
  - a system-derived SHA-256 link to the exact canonical Accounting Policy V1 document;
  - entries preserving caller-provided raw path bytes;
  - a system-derived SHA-256 digest of each entry's exact content bytes; and
  - each entry's exact byte length.
- Construction from the exact caller-provided path set and content, with deterministic entry ordering and rejection of duplicate or invalid paths.
- Strict decoding that rejects malformed, ambiguous, unknown-field, duplicate-field, non-canonical, or incorrectly linked input rather than repairing it.
- Canonical bytes and a system-derived Inventory V1 document digest suitable for the later Accounting Result V1 provenance link.
- Explicit empty-inventory semantics: an empty document represents a successfully verified empty requested scope, never unavailable acquisition.
- Focused root-contract tests for canonicalization, provenance, raw-path fidelity, digest and length derivation, strict rejection, immutability, and additive legacy coexistence.

### Out of Scope

- Filesystem discovery, path enumeration, file opening, safe descriptor-relative reads, race detection, or integration with the existing internal inventory adapter; that acquisition work is deferred to `WU-4C`.
- Any definition of discovery semantics beyond consuming the exact path set supplied by the caller.
- Changes to `UntrackedRecord`, `UntrackedInventory`, `InventoryUnavailableError`, or existing internal Unix/non-Unix acquisition behavior.
- Caller-supplied content digests, opaque compatibility hashes, or CNSIC/project-specific inventory hash conventions.
- Accounting Result V1, accounting behavior, Git acquisition, CLI behavior, publication, consumer-profile adaptation, or delivery authority.
- Changes to the historical `bootstrap-neutral-core-extraction` authority or its worktree.

## Capabilities

### New Capabilities

- `inventory-v1-contract`: Construct, validate, strictly decode, and canonically serialize a neutral Inventory V1 document linked to Accounting Policy V1.
- `inventory-v1-content-identity`: Derive exact-byte SHA-256 identity and byte length for every caller-provided inventory entry while preserving raw path bytes.

### Modified Capabilities

None. Legacy untracked-inventory behavior continues unchanged.

## Approach

Keep Inventory V1 entirely in the functional core. Constructors receive the exact validated Accounting Policy V1 document and caller-provided entry path/content values, copy mutable inputs, validate the complete set, derive policy and content identities internally, and emit one canonical representation. The validating decoder requires the exact antecedent policy value, applies closed-shape and canonical-byte checks, and verifies the document-derived policy link.

The pure contract does not discover files or interpret acquisition failure. Acquisition remains an adapter concern: unavailable acquisition must remain an unavailable outcome, while an empty Inventory V1 may be created only when the caller has established that the requested scope was successfully verified as empty.

Preliminary implementation forecasting is approximately **450 changed lines**: about **170 production**, **260 tests**, and **20 artifacts**. The likely delivery shape is one cohesive **2A/2B** contract-and-test increment; this proposal does not decide PR splitting, chaining, or other delivery mechanics.

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| Public root Go API | New | Inventory V1 immutable value, constructor, strict decoder, canonical bytes, and document digest |
| Root contract tests | New | Positive and negative vectors for identity, ordering, linkage, strictness, empty semantics, and defensive copying |
| Legacy inventory API | Coexists | Existing path-only records and acquisition semantics remain behaviorally unchanged |
| Provenance chain | Extended | Inventory V1 consumes Accounting Policy V1 and becomes the required second antecedent for later Accounting Result V1 |

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| A caller-provided digest is accidentally treated as authoritative | Medium | Accept exact content and antecedent document values; derive all SHA-256 values internally; add negative API/decoder tests |
| Raw paths are corrupted through string or UTF-8 normalization | Medium | Preserve copied raw bytes, define ordering over raw bytes, and test hostile/non-UTF-8 path values |
| Permissive decoding admits multiple byte representations | Medium | Use strict closed-shape decoding, duplicate-key detection, semantic validation, canonical re-encoding, and byte equality |
| Empty inventory is confused with acquisition failure | Medium | Specify empty as verified empty scope only; keep unavailable acquisition outside the V1 document and test the distinction |
| New V1 work changes legacy semantics | Low | Add rather than replace APIs and run legacy tests unchanged |
| Acquisition integration expands this contract increment | Medium | Keep filesystem and adapter surfaces untouched and defer them explicitly to `WU-4C` |
| Policy linkage is syntactically valid but not authoritative | Medium | Require the exact Policy V1 document for construction and validating decode; compute and compare its canonical digest internally |

## Rollback Plan

Remove only the new root Inventory V1 contract and its focused tests. No consumer migration or acquisition integration occurs in this change, so rollback does not require modifying the legacy `UntrackedInventory` path or the historical bootstrap authority. Accounting Policy V1 remains independently valid, and Accounting Result V1 must not proceed until a replacement second-document contract is approved.

## Dependencies

- Accounting Policy V1 must provide a validated canonical antecedent document and system-derived document digest.
- Existing canonical-document and strict-decoding conventions in the neutral root package are the implementation baseline.
- `WU-4C` remains a later dependency for safe filesystem acquisition and adapter integration, not a prerequisite for this pure contract.

## Success Criteria

- [ ] Inventory V1 has a neutral, versioned, canonical root-package document/API and does not contain CNSIC or project-specific compatibility identity.
- [ ] Every entry preserves exact raw path bytes and derives SHA-256 content identity plus exact byte length from copied content, never from a caller assertion.
- [ ] Construction and validating decode require the exact Accounting Policy V1 document and verify the policy link from its canonical bytes.
- [ ] Repeated equivalent valid inputs produce byte-identical canonical Inventory V1 documents and the same system-derived document digest.
- [ ] The decoder rejects duplicate or unknown fields, invalid types or values, non-canonical ordering/encoding, duplicate paths, and incorrect policy linkage.
- [ ] Empty Inventory V1 is documented and tested as verified empty requested scope; unavailable acquisition cannot be encoded as empty inventory.
- [ ] The API consumes only the exact caller-provided path set and defines no filesystem discovery behavior.
- [ ] Existing `UntrackedRecord`, `UntrackedInventory`, defensive-copy behavior, unavailable outcomes, and internal acquisition tests remain unchanged and passing.
- [ ] Filesystem acquisition and `WU-4C` integration, Accounting Result V1, and delivery mechanics remain outside this change.
