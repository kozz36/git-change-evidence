# Neutral Evidence Contracts Specification

## Purpose

Define strict, versioned, neutral evidence documents with stable canonical bytes.

## Requirements

### Requirement: Closed versioned contracts

The system MUST validate public evidence documents against declared contract versions, exact types, closed shapes, and discriminated variants. It MUST reject malformed, ambiguous, consumer-specific, or authority-bearing fields rather than silently normalizing them.

#### Scenario: Valid neutral document
- GIVEN a document conforming to a declared contract version
- WHEN it is validated
- THEN validation succeeds without adding or changing fields

#### Scenario: Unknown or authority field
- GIVEN a document containing an unknown field or a delivery-authority claim
- WHEN it is validated
- THEN validation fails with a typed contract error

### Requirement: Canonical content identity

The system MUST serialize equivalent valid documents to identical canonical JSON bytes. Where byte identity is contractual, it MUST reject non-canonical input representations. A content digest MUST be exactly 64 hexadecimal characters and MUST NOT accept SHA-1-length or alternative-length values where a content digest is required.

#### Scenario: Equivalent documents
- GIVEN two equivalent valid documents with different input ordering
- WHEN each is canonically serialized
- THEN the resulting bytes and 64-character content digest are identical

#### Scenario: Invalid digest length
- GIVEN a document containing a 40-character or non-64-character content digest
- WHEN it is validated
- THEN validation fails before the document is accepted

#### Scenario: Non-canonical contractual input
- GIVEN a byte-identity contractual document uses a non-canonical representation
- WHEN it is validated
- THEN validation fails without normalizing the representation

### Requirement: Provenance and reproducibility

Every report MUST retain immutable revision identity plus input-policy, inventory, and accounting provenance. The same immutable revisions and profile MUST produce identical canonical report bytes in independent runs.

#### Scenario: Independent reproducible report
- GIVEN two independent runs use the same immutable revisions and profile
- WHEN each assembles a report
- THEN both canonical report bytes match and each report retains all required provenance

### Requirement: Neutral Go package API

The system MUST provide a neutral, versioned public Go package API at the module root alongside document contracts for validation, canonical evidence construction, and report use. The API MUST accept immutable typed Go values at its core boundary. It MUST NOT select a module/publication identity or concrete configuration-file syntax as part of that API behavior.

#### Scenario: Versioned Go package use
- GIVEN a caller uses a declared Go package API version with valid immutable neutral inputs
- WHEN it constructs and validates evidence
- THEN it receives the corresponding versioned neutral contract behavior without a module/publication identity or configuration-file syntax selection

### Requirement: Product-surface neutrality

The product MUST NOT define universal changed-lines limits, review rubrics, forecast vocabulary, or governance models as core defaults. It MUST NOT provide remote delivery orchestration, pull-request management, CI administration, UI, database, or RBAC product surfaces.

#### Scenario: Attempted universal governance default
- GIVEN a caller requests a universal limit, rubric, forecast vocabulary, or governance model
- WHEN neutral core behavior is configured
- THEN no core default is supplied and profile-owned policy remains external

#### Scenario: Attempted excluded product surface
- GIVEN a caller requests remote delivery, pull-request, CI, UI, database, or RBAC behavior
- WHEN the product surface is evaluated
- THEN the request is outside the product and no such surface is provided
