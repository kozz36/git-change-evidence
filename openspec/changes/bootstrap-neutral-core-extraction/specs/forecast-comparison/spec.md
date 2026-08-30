# Forecast Comparison Specification

## Purpose

Compare supplied forecasts to measured evidence while preserving forecast ownership and delivery neutrality.

## Requirements

### Requirement: Supplied forecast ingestion

The system MUST validate and ingest a supplied forecast as input to comparison. It MUST NOT author, infer, complete, mutate, or overwrite a forecast, including when evidence is missing or a profile supplies thresholds.

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
