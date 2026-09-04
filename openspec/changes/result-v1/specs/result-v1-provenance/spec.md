# Result V1 Provenance Specification

## Purpose

Define the exact Policy V1 → Inventory V1 → Result V1 provenance chain and prohibit asserted, substituted, or consumer-specific authority.

## ADDED Requirements

### Requirement: Exact validated antecedent chain

Result V1 construction and strict decoding MUST require the exact canonical `PolicyDocumentV1` and exact canonical `InventoryDocumentV1` documents, not merely values that look like their digests. For every accepted result, the policy link MUST equal the SHA-256 digest of the exact Policy V1 canonical bytes, the inventory link MUST equal the SHA-256 digest of the exact Inventory V1 canonical bytes, and the inventory's `accounting_policy_sha256` MUST equal that same policy digest.

A Result V1 document MUST serialize its `accounting_policy_sha256` from the exact validated Policy V1 canonical bytes and its `inventory_sha256` from the exact validated Inventory V1 canonical bytes. The policy and inventory digests in the Result V1 body are links, not caller assertions. A valid digest from a different policy or inventory, a digest that is syntactically well-formed but has no document, a noncanonical document value, a zero document value, or an inventory linked to another policy MUST NOT establish provenance.

The Result V1 decoder MUST not accept a result against an antecedent solely because the candidate result carries matching-looking digest text. Acceptance MUST require the exact supplied antecedents and a valid complete policy → inventory → result chain. A successful decode establishes only that the self-described result is canonically consistent with its validated antecedents and embedded revisions/entries; it MUST NOT claim Git acquisition, object existence, ref stability, filesystem state, or delivery authority.

#### Scenario: Byte-identical antecedents reproduce provenance
- **GIVEN** two independently constructed Policy V1 documents with byte-identical canonical bytes, two Inventory V1 documents with byte-identical canonical bytes linked to those policies, and equivalent valid committed snapshots
- **WHEN** each antecedent pair constructs Result V1
- **THEN** both results contain the same derived policy and inventory links, canonical bytes, and Result V1 digest

#### Scenario: Inventory policy link is wrong but all digests are well formed
- **GIVEN** a valid policy `P1`, a valid inventory `I2` linked to different valid policy `P2`, and a candidate Result V1 body whose policy and inventory digest strings are each valid lowercase SHA-256 values
- **WHEN** construction or strict decoding is requested with `P1` and `I2`
- **THEN** it fails fatally and emits or returns no Result V1 document

#### Scenario: Same accounting with a different antecedent is not interchangeable
- **GIVEN** two policies or inventories with different canonical bytes but committed entries that would yield identical classifications and totals
- **WHEN** Result V1 is constructed for each valid antecedent chain
- **THEN** the results retain the distinct derived links and MUST NOT be treated as interchangeable evidence documents

### Requirement: Result identity is system-derived and complete

The Result V1 document digest MUST be exactly the SHA-256 digest of the complete canonical Result V1 bytes, including schema, both antecedent links, revisions, every classification, all totals, and every observation plus the final LF. Every constructed or decoded Result V1 value MUST expose that exact digest. A purported Result V1 value whose exposed digest differs from its canonical-byte SHA-256 MUST be rejected even when its body is otherwise decodable.

Changing any Result V1 evidence that is validly representable—including either antecedent link, either revision, a raw path byte, a classification, a countable value, a non-countable marker, a total, an observation, or ordering—MUST change the canonical bytes and therefore produce a distinct Result V1 identity. Equivalent valid input ordering is the sole permitted construction-time normalization and MUST converge before identity derivation; strict decode MUST reject instead of normalize noncanonical order.

#### Scenario: One raw path byte changes the identity
- **GIVEN** two otherwise equivalent Result V1 inputs differ only by one committed primary-path byte, such as LF versus CRLF-adjacent byte sequence or `0xff` versus `0xfe`
- **WHEN** each valid input is constructed
- **THEN** each path is Base64-encoded from its exact bytes and the two canonical documents and Result V1 digests differ

#### Scenario: Reordered candidate is not repaired
- **GIVEN** canonical Result V1 bytes are modified only by swapping two raw-byte-sorted entries or two policy-ordered totals
- **WHEN** strict decoding is requested with the exact antecedents
- **THEN** decoding fails without sorting, recomputing a replacement digest, or returning a repaired document

### Requirement: Result V1 remains neutral evidence-only provenance

Result V1 provenance fields, classifications, totals, measurements, and observations MUST describe technical evidence only. Result V1 field names and structural metadata MUST NOT introduce approval, rejection, denial, pass/fail, gating, blocking, merge authorization, deployment authorization, release authorization, delivery authority, CNSIC compatibility, a consumer profile vocabulary, a project-specific hash, or a review decision. A permitted field value, including an exact Policy V1-derived category or other policy-provided string, MUST NOT be rejected merely because its content contains words such as `merge` or `release`.

The Result V1 contract MUST NOT add Git acquisition, moving-reference transactions, retry/revalidation, filesystem or untracked-file acquisition, CLI commands or output, report projection expansion, publication, configuration syntax, packaging, distribution, release mechanics, delivery mechanics, or CNSIC integration. These concerns remain outside the pure canonical contract.

#### Scenario: Authority-bearing candidate provenance
- **GIVEN** a producer adds an object field such as `approved`, `review_gate`, `release`, `deliveryAuthority`, `cnsic_status`, or `project_compatibility_hash` at any Result V1 object level
- **WHEN** the candidate is strictly decoded
- **THEN** decoding rejects the unknown or authority-bearing field name or structural metadata and does not expose it as Result V1 evidence

#### Scenario: Policy-derived authority-like value remains evidence
- **GIVEN** a validated Policy V1 document permits a category or other policy-provided string value containing `merge` or `release`, and an otherwise valid Result V1 uses that value only in its permitted closed-shape field
- **WHEN** Result V1 is constructed or strictly decoded with the exact antecedents
- **THEN** the value is evaluated under the Policy V1 and Result V1 semantic constraints and is not rejected solely because of that vocabulary

#### Scenario: Out-of-scope integration request
- **GIVEN** a caller requests Result V1 to acquire a moving ref, read an untracked file, invoke a CLI, publish a result, select a configuration format, or decide whether a change may be delivered
- **WHEN** the request is evaluated
- **THEN** it is outside Result V1 and no such behavior or authority is supplied
