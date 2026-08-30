# Project Profile Adaptation Specification

## Purpose

Adapt project policy and compatibility needs at a one-way edge around a neutral evidence core.

## Requirements

### Requirement: One-way profile boundary

The system MUST accept project vocabulary, categories, thresholds, and compatibility requirements through a profile edge. The neutral core MUST NOT depend on, embed, or expose CNSIC-specific vocabulary, review rules, or delivery authority; a profile MUST NOT change neutral contract meaning.

#### Scenario: CNSIC profile translation
- GIVEN a CNSIC profile supplies its vocabulary and threshold values
- WHEN neutral evidence is produced
- THEN core evidence remains neutral while the profile may translate it at the edge

#### Scenario: Distinct project profile
- GIVEN another project supplies different vocabulary and limits
- WHEN it uses the same neutral core
- THEN its profile values do not alter neutral contract semantics

#### Scenario: Attempted authority injection
- GIVEN a profile attempts to inject consumer review rules, governance semantics, or delivery authority into neutral core behavior or contracts
- WHEN neutral evidence is requested
- THEN the attempted injection is rejected and is absent from neutral evidence

### Requirement: Parity-gated reversible cutover

A consumer cutover MUST occur only after the profile demonstrates parity against the frozen compatibility vectors and applicable real-Git evidence. The cutover MUST remain reversible to the predecessor path on parity failure, and MUST NOT treat parity as delivery authority.

#### Scenario: Parity demonstrated
- GIVEN a consumer profile passes the frozen vectors and applicable real-Git parity evidence
- WHEN a cutover is proposed
- THEN the consumer is eligible for a reversible cutover

#### Scenario: Parity failure
- GIVEN a profile fails a frozen vector or applicable real-Git parity case
- WHEN cutover is evaluated
- THEN source-of-truth transition is withheld and the predecessor path remains available
