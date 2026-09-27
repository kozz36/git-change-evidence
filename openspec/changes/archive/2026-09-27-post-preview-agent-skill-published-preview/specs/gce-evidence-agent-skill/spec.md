# GCE Evidence Agent Skill Specification

## Purpose

Define repo-local agent guidance for consuming, validating, reproducing, regenerating, and technically interpreting Git Change Evidence (GCE) without changing GCE runtime contracts or assigning delivery authority.

## Requirements

### Requirement: Intent-based skill activation

The repo-local skill MUST activate for explicit requests to consume, validate, reproduce, regenerate, or technically interpret GCE evidence, including Go AST census evidence. It MUST NOT activate solely because a request mentions Go, AST, census, GCE, or review. It MUST distinguish evidence work from implementation, contract, or specification authoring.

#### Scenario: Positive evidence-use request

- GIVEN an agent is asked to validate or reproduce specified GCE evidence
- WHEN it evaluates whether the skill applies
- THEN it MUST activate for that evidence task and supplied scope

#### Scenario: Negative generic or authoring request

- GIVEN a request asks only for generic Go help, ordinary code changes, unrelated review, or GCE implementation/specification authoring
- WHEN the agent evaluates whether the skill applies
- THEN it MUST NOT treat a bare keyword match as sufficient reason to activate the evidence-consumption workflow

#### Scenario: Mixed evidence and delivery-verdict request

- GIVEN a request combines a GCE evidence task with a request to approve, reject, gate, block, merge, deploy, or release
- WHEN the agent responds
- THEN it MUST separate and perform only the authorized technical evidence task and MUST leave the delivery decision to the consumer

#### Scenario: Ambiguous evidence scope

- GIVEN evidence intent is plausible but the evidence family, scope, inputs, or authorization is unclear
- WHEN the agent cannot proceed without assuming missing information
- THEN it MUST request clarification and MUST NOT invent scope or authority

### Requirement: Technical observations confer no delivery authority

The skill MUST constrain outputs to technical states, measurements, observations, provenance, validation or reproduction status, limitations, and recovery performed or required. It MUST NOT approve, reject, gate, block, merge, deploy, release, or otherwise assign delivery authority. It MUST describe Go AST census matches as exact textual syntax candidates within the validated supplied scope; zero matches MUST NOT be represented as semantic absence, whole-repository completeness, safety, or approval.

#### Scenario: Technical result is reported without a verdict

- GIVEN GCE evidence has been technically validated
- WHEN the agent reports its result
- THEN it MUST state the observed evidence and its scope and limitations without converting the result into a delivery verdict

#### Scenario: Census reports zero matches

- GIVEN a census completes successfully with zero matches
- WHEN the agent interprets the result
- THEN it MUST report no candidate for that query in the supplied scope and MUST NOT claim semantic absence or repository-wide completeness

### Requirement: Supported consumer surfaces remain compatible

The skill MUST identify the CLI and machine-readable contracts as the primary consumer surface while preserving the public root Go API and public `census` subpackage as supported surfaces. It MUST present `gce census go-ast` as the canonical pre-v1 census command and `git-change-evidence census-go-ast` as supported with the established contract meaning throughout pre-v1. The skill MUST NOT imply that it changes or supersedes these contracts.

#### Scenario: Canonical and legacy CLI guidance

- GIVEN an agent is directed to consume Go AST census evidence through the CLI
- WHEN it selects a command
- THEN it MUST prefer `gce census go-ast` and MUST recognize `git-change-evidence census-go-ast` as a supported compatible pre-v1 command

#### Scenario: Public Go API integration

- GIVEN an integration consumes GCE through Go rather than the CLI
- WHEN the agent identifies supported consumer surfaces
- THEN it MUST retain both the public root Go API and public `census` subpackage as supported and MUST NOT claim that all evidence families have a CLI assembler

### Requirement: Evidence capture and reproduction are faithful

The skill MUST require preservation of applicable canonical output bytes including the final LF, digests, exact antecedents, raw-path Base64 identities, argv, stdout, stderr, and exit status as distinct observations. It MUST require validation against the declared existing contract rather than an invented universal envelope. It MUST distinguish reproduction from reacquisition: changed inputs produce new evidence, not reproduction. Census reproduction MUST be described as requiring the same captured source bytes, selected paths, query, limits, and compatible extractor behavior; revision-bound evidence MUST use its exact immutable revisions and canonical antecedents.

#### Scenario: Evidence is captured and validated

- GIVEN an agent obtains GCE evidence for a requested scope
- WHEN it records and validates the result
- THEN it MUST preserve the applicable raw observations and antecedents and validate them against the evidence family's declared contract

#### Scenario: Exact census reproduction

- GIVEN the original source bytes, selected paths, query, limits, and compatible extractor behavior are available
- WHEN the agent attempts to reproduce census evidence
- THEN it MUST use those same inputs and MUST identify the result as reproduction only when they match

#### Scenario: Required reproduction inputs are unavailable or changed

- GIVEN any required original input or antecedent is unavailable or differs
- WHEN the agent attempts to repeat the evidence operation
- THEN it MUST report that reproduction cannot be established and MUST label any operation on changed inputs as new evidence

### Requirement: Failures and recovery remain transparent

The skill MUST treat nonzero exits, interruption, partial writes, malformed or truncated bytes, unavailable inputs, and stale or mismatched antecedents as failures or unresolved evidence, never as empty success. For Go AST census it MUST preserve the documented meanings of exits 0, 2, 3, and 4, including zero-match success and input, unavailable-content, and bounds failures. It MUST prohibit silent scope narrowing, unapproved limit changes, confinement bypass, and indefinite retries; recovery MUST be a transparent correction or authorized reacquisition.

#### Scenario: Census exit codes are interpreted in context

- GIVEN a Go AST census returns exit 0, 2, 3, or 4
- WHEN the agent reports the outcome
- THEN it MUST apply the documented census-specific meaning and MUST NOT generalize that exit mapping to unrelated commands

#### Scenario: Partial or invalid result is observed

- GIVEN execution is interrupted, output is partial or malformed, or antecedents are stale or mismatched
- WHEN the agent evaluates the evidence
- THEN it MUST preserve the failure observation and MUST NOT interpret missing or prefix data as empty success

#### Scenario: Recovery changes inputs

- GIVEN recovery requires correcting inputs or reacquiring evidence
- WHEN the agent proceeds with appropriate authorization
- THEN it MUST disclose the correction or reacquisition and MUST identify changed-input output as new evidence

### Requirement: Skill discoverability is repo-local and additive

The skill MUST be available as repository-local guidance through a future additive project-root `AGENTS.md` entry that identifies the skill's purpose and local skill path. Registration MUST preserve all pre-existing `AGENTS.md` content and governance, including stale preview wording, without repeating stale publication assertions as current facts. Repo-local registration MUST NOT imply automatic installation, guaranteed activation by every runtime, packaging, external registry publication, or inclusion in preview assets.

#### Scenario: Project guidance discovers the skill

- GIVEN the repo-local skill is present and its registration is added
- WHEN an agent reads project-root `AGENTS.md`
- THEN it MUST be able to discover the skill purpose and follow its local path without requiring a machine-local or external registry

#### Scenario: Registration preserves existing project guidance

- GIVEN an additive skill registration is made in project-root `AGENTS.md`
- WHEN the resulting guidance is reviewed
- THEN all pre-existing bytes and governance MUST remain preserved, and the new entry MUST NOT adopt stale preview-status wording as fact

### Requirement: Validation claims distinguish evidence levels

Skill validation MUST distinguish static inspection and fixture checks from executable contract tests and actual agent evaluation. Static or fixture checks MUST NOT be represented as proof of actual agent behavior. Scenario coverage MUST include positive, negative, mixed, and clarification triggers; supported CLI and Go API surfaces; zero-match limits; census exits 2/3/4; incomplete writes; stale or mismatched inputs; exact-input reproduction; and no-delivery-authority language.

#### Scenario: Static or fixture checks pass without agent evaluation

- GIVEN static checks or representative fixtures pass but no actual agent evaluation was performed
- WHEN validation is reported
- THEN it MUST identify only the checks performed and MUST NOT claim demonstrated agent reliability

#### Scenario: Validation evidence levels are reported separately

- GIVEN static checks, executable contract tests, and/or actual agent evaluation have been performed
- WHEN results are summarized
- THEN each performed evidence level MUST be identified separately from levels not performed

### Requirement: Skill scope does not alter release or language plans

The skill MUST remain repo-local and not packaged, and its metadata MUST NOT imply a CLI, extractor, or schema version, distribution, or compatibility policy. It MUST NOT require or cause a change to the published `v0.1.0-preview.1` tag, release, assets, or target commit, and MUST NOT introduce a retroactive preview gate or re-verification prerequisite. Multilingual census MUST remain separate from this capability.

#### Scenario: Skill guidance is added after preview publication

- GIVEN the preview release is already published and the agent skill is added as repository content
- WHEN the skill's scope or status is described
- THEN it MUST NOT imply a new preview publication, retroactive release condition, preview asset, or change to existing published objects

#### Scenario: Multilingual census is mentioned

- GIVEN the skill or its scope discusses future multilingual census
- WHEN languages are described
- THEN it MUST identify multilingual census as a separate roadmap and MUST NOT prescribe its implementation as part of this capability
