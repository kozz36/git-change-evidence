# Bounded Carveout Analysis Specification

## Purpose

Compare immutable content within an explicit bounded carveout without mutable-state substitution.

## Requirements

### Requirement: Immutable bounded comparison

The system MUST analyze only supplied immutable content and an explicit carveout bound. It MUST produce deterministic matching evidence for identical inputs and MUST NOT substitute index, working-tree, or ambient filesystem content for immutable input.

#### Scenario: Matching immutable content
- GIVEN identical immutable content and the same valid carveout bound
- WHEN analysis is repeated
- THEN both runs produce identical carveout evidence

#### Scenario: Working-tree divergence
- GIVEN immutable content and different working-tree content at the same path
- WHEN analysis is requested
- THEN the result is derived from immutable content only

#### Scenario: Available content completes within bound
- GIVEN valid immutable content and a carveout analysis that completes within its declared bound
- WHEN analysis is requested
- THEN the result contains completed comparison evidence and no unavailable or timeout evidence

### Requirement: Safe carveout bounds

The system MUST reject absent, malformed, non-positive, or unsafe carveout bounds before analysis starts. It MUST return typed invalid-bound evidence and MUST NOT produce a comparison result or invented measurement.

#### Scenario: Invalid carveout bound
- GIVEN a carveout bound that is absent, malformed, non-positive, or unsafe
- WHEN analysis is requested
- THEN no analysis starts and the result records typed invalid-bound evidence

### Requirement: Typed incomplete evidence

The system MUST represent unavailable content and exceeded analysis time as distinct typed evidence states. It MUST NOT convert either state into a match, mismatch, zero measurement, or inferred result.

#### Scenario: Content unavailable
- GIVEN required immutable content cannot be obtained
- WHEN analysis is requested
- THEN the result records typed unavailable evidence

#### Scenario: Bounded timeout
- GIVEN analysis exceeds its declared bound
- WHEN the bound is reached
- THEN the result records typed timeout evidence and no inferred comparison result
