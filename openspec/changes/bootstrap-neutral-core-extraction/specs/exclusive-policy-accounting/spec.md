# Exclusive Policy Accounting Specification

## Purpose

Classify changed entries with injected policy and expose factual accounting observations.

## Requirements

### Requirement: Injected exclusive classification

The system MUST receive categories, precedence, and default treatment from a profile-supplied policy. Unless that profile explicitly specifies another default, an unmatched path MUST receive its production-equivalent category. It MUST validate that policy before emitting totals, and MUST assign every changed entry to exactly one category without embedding project vocabulary or policy values.

#### Scenario: Overlapping category rules
- GIVEN a valid policy with two matching categories and explicit precedence
- WHEN a changed entry matches both
- THEN the entry is assigned only to the highest-precedence category

#### Scenario: Invalid policy
- GIVEN a policy with ambiguous precedence or invalid category definition
- WHEN accounting is requested
- THEN accounting fails with typed policy evidence and emits no totals

#### Scenario: Unmatched path default
- GIVEN a changed path matches no explicit category and the profile specifies no alternate default
- WHEN accounting is requested
- THEN it receives the profile's production-equivalent category exactly once

### Requirement: Explicit measurements and observations

The system MUST NOT invent totals for binary, non-countable, unavailable, or untrusted data. Configurable thresholds and ratios MAY produce non-gating observations or warnings, but MUST NOT approve, reject, block, or otherwise decide delivery.

#### Scenario: Non-countable changed entry
- GIVEN an entry without a trustworthy line measurement
- WHEN totals are assembled
- THEN its measurement remains explicit and no fictional line total is added

#### Scenario: Trustworthy text measurement
- GIVEN a text entry has declared trustworthy additions and deletions
- WHEN totals are assembled
- THEN its exact additions and deletions contribute only to its exclusive category total

#### Scenario: Threshold crossing
- GIVEN a profile threshold is crossed
- WHEN accounting is reported
- THEN the report records an observation or warning and no delivery verdict
