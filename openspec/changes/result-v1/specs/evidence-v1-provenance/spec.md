# Evidence V1 Provenance Delta Specification

## ADDED Requirements

### Requirement: Provenance and reproducibility

Every Report V1 document MUST retain immutable revision identity plus input-policy, inventory, and accounting provenance. For Result V1-aware report construction and provenance validation, `provenance.accounting_sha256` MUST be the SHA-256 digest of the exact canonical Result V1 document, not a free-form digest assertion. The report boundary MUST require the exact Result V1 document and its exact Policy V1 and Inventory V1 antecedents, validate that complete chain, and derive rather than trust all of the following report provenance values:

- `provenance.accounting_sha256` from the Result V1 canonical bytes;
- `provenance.input_policy_sha256` from the exact Policy V1 canonical bytes linked by Result V1;
- `provenance.inventory_sha256` from the exact Inventory V1 canonical bytes linked by Result V1; and
- `provenance.revisions.base` and `provenance.revisions.head` from the immutable revisions embedded in Result V1.

The binding boundary MUST reject construction or provenance validation when any supplied report provenance value differs from those derived values, when the Result V1 document is invalid or noncanonical, when either antecedent is invalid or mismatched, or when a syntactically valid 64-character accounting digest identifies a different Result V1 document. A report parser that receives only standalone report bytes MUST NOT represent syntactic digest validation as proof of Result V1 provenance; a successful Result V1 provenance validation requires the exact Result V1 document and exact antecedents.

The existing Report V1 wire object shape, field names, canonical ordering, schema, human and machine projections, CLI behavior, and publication behavior MUST NOT expand in this change. Result V1 content, classifications, totals, observations, diagnostics, and antecedent documents MUST NOT be embedded in Report V1 merely to establish this binding. The report remains an outer evidence envelope and MUST NOT add approval, rejection, gating, blocking, merge, deployment, release, or delivery authority.

Given the same exact Policy V1, Inventory V1, Result V1, subject, and immutable revisions, independent bound report construction MUST produce byte-identical canonical Report V1 bytes and the same report digest.

#### Scenario: Bound report derives all provenance
- **GIVEN** an exact validated Policy V1 document, its exact linked Inventory V1 document, an exact validated Result V1 document, and a report subject
- **WHEN** a Result V1-aware Report V1 construction path creates the outer report
- **THEN** `accounting_sha256` equals exactly `SHA-256(ResultV1.CanonicalBytes())`, the policy, inventory, and revisions equal the values derived through Result V1, and the Report V1 wire shape contains no Result V1 body

#### Scenario: Well-formed but mismatched accounting digest
- **GIVEN** a valid Result V1 document and a caller-supplied report provenance value containing the 64-lowercase-hex digest of a different valid Result V1 document
- **WHEN** bound report construction or provenance validation is requested
- **THEN** it fails without emitting or accepting a provenance-validated report

#### Scenario: Standalone report bytes are not provenance proof
- **GIVEN** canonical Report V1 bytes with syntactically valid policy, inventory, accounting digests, and revisions but without the exact Result V1 and antecedent documents at validation time
- **WHEN** a consumer asks whether `accounting_sha256` is bound to Result V1
- **THEN** the consumer receives no binding proof from digest syntax alone; only the antecedent-aware validation path may establish it

#### Scenario: No projection or authority expansion
- **GIVEN** a valid Result V1-bound Report V1 contains an exceeded accounting observation
- **WHEN** the report is projected or otherwise handled by existing Report V1 behavior
- **THEN** its existing wire and projection surfaces remain unchanged and no approval, denial, gate, block, merge, deploy, release, or delivery verdict is added
