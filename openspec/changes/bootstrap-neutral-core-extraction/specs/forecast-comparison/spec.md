# Forecast Comparison Specification

## Purpose

Compare supplied forecasts to measured evidence while preserving forecast ownership and delivery neutrality.

## Requirements

### Requirement: Supplied forecast ingestion

The system MUST validate and ingest a supplied forecast as input to comparison. It MUST NOT author, infer, complete, mutate, or overwrite a forecast, including when evidence is missing or a profile supplies thresholds.

Ingestion uses root typed `IngestForecastV1(ForecastInputV1) (ForecastV1, error)`: declared policy/inventory `DocumentDigest` and coherent `RevisionIdentity`, with closed per-category additions/deletions (lines) and non-countable (entries), all optional `uint64` values. Malformed identity/keys and duplicate keys return `ContractError`; valid MaxUint64 is retained.
Input and getter slices/pointers are independently copied without sorting or completion. Known zero, explicit unknown, omitted metrics, and nil versus empty slices remain distinct; an empty metric list infers nothing.
Root typed `CompareForecastV1(ResultDocumentV1, *ForecastV1, []ForecastThresholdV1) (ForecastComparisonV1, error)` retains independent copies: nil whole forecast is absent; policy/inventory/base/head binding mismatch is incompatible without numeric comparison. Observations use the measured/supplied-key union sorted by category then kind; a missing measured category is unavailable, not zero. Non-countable >0 makes additions/deletions unavailable but its own count remains known.
Deltas use measured-minus-supplied direction and ordered-subtraction uint64 magnitude. Optional per-key absolute thresholds share the metric unit and use strict greater-than; equality is not exceeded. Optional measured/forecast/delta/limit/exceeded pointers and forecast presence preserve unavailable, unknown, omitted, and absent-threshold distinctions; incompatible or unavailable comparisons compute neither delta nor exceeded.

#### Scenario: Valid supplied forecast
- GIVEN a valid supplied forecast and measured evidence
- WHEN comparison is requested
- THEN the forecast is retained as supplied and comparison evidence is produced

#### Scenario: Missing forecast
- GIVEN measured evidence without a supplied forecast
- WHEN comparison is requested
- THEN the system reports absent forecast input and does not create one

### Requirement: Neutral comparison observations

The system MUST report comparison differences, unknown values, and configured threshold observations as evidence. It MUST NOT represent any comparison outcome as approval, rejection, blocking, or delivery authority.

#### Scenario: Divergent forecast
- GIVEN a supplied forecast differs from measured evidence
- WHEN comparison is produced
- THEN the difference is explicit in the comparison evidence

#### Scenario: Unavailable measurement
- GIVEN a supplied forecast and an unavailable measured value
- WHEN comparison is produced
- THEN the affected comparison is explicitly unavailable rather than inferred or verdict-bearing

#### Scenario: Configured threshold observation
- GIVEN a valid comparison crosses a profile-configured threshold
- WHEN comparison is produced
- THEN the threshold observation is explicit evidence and does not gate or decide delivery
