# Evidence CLI Publication Specification

## Purpose

Expose deterministic evidence projections and publish immutable report content safely.

## Requirements

### Requirement: Deterministic CLI projections

The CLI MUST project the same canonical evidence into deterministic machine-readable and human-readable outputs for identical inputs. Projections MUST preserve evidence-only semantics and MUST NOT emit approval, denial, gating, merge, deployment, or release decisions.

#### Scenario: Repeated projection
- GIVEN identical canonical evidence and projection options
- WHEN the CLI runs twice
- THEN each selected projection is byte-identical

#### Scenario: Attention observation
- GIVEN evidence contains a threshold attention observation
- WHEN the CLI projects it
- THEN output presents it as evidence without a delivery verdict

### Requirement: Safe report diagnostics

Human-readable report projections and diagnostics MUST redact or omit sensitive repository paths, command output, credentials, and environment values. They MUST preserve technical usefulness without exposing those values and MUST NOT echo sensitive Subject or raw input, error, process, or environment details. This privacy boundary does not sanitize canonical JSON: it retains literal caller-supplied Subject under the canonical report contract. Evidence-only semantics remain unchanged; this clarification changes no formats, transport, export, or publication behavior.

#### Scenario: Sensitive diagnostic input
- GIVEN an execution failure contains sensitive paths, command output, credentials, or environment values
- WHEN a human-readable report projection or diagnostic is emitted
- THEN the emitted output omits or redacts every sensitive value

### Requirement: Technical exit semantics

The CLI MUST use exits only for technical execution outcomes and MUST assign distinct classifications to successful evidence production, invalid input, unavailable evidence, and publication conflict. A threshold observation MUST NOT cause a delivery-gating exit.

#### Scenario: Threshold attention
- GIVEN evidence production succeeds with a threshold attention observation
- WHEN the CLI completes
- THEN it returns the successful technical outcome

#### Scenario: Invalid input
- GIVEN CLI input violates a contract
- WHEN the CLI is invoked
- THEN it returns a technical invalid-input outcome without publishing evidence

#### Scenario: Unavailable evidence
- GIVEN required evidence is unavailable
- WHEN the CLI is invoked
- THEN it returns the unavailable-evidence technical classification without a delivery verdict

#### Scenario: Publication conflict
- GIVEN evidence is valid but publication conflicts with an existing destination
- WHEN the CLI is invoked
- THEN it returns the publication-conflict technical classification without overwriting evidence

### Requirement: Safe exclusive publication

The opt-in Go `shell/livepublication.Publish` API MUST perform at most three complete acquire/build/validate/reacquire attempts against the caller's original repository and refs. It MUST validate canonical evidence and provenance against the acquired base/head IDs before invoking the existing publisher. A changed pair MUST discard all evidence and trigger a complete rebuild; persistent movement MUST return existing typed `SnapshotRacing` without a publisher call or artifact. Invalid or divergent evidence MUST return a typed technical error before publication. The frozen publisher and CLI remain independent of Git.

The accepted snapshot MUST mean the last observed matching pair before publication, not an atomic read of both refs or a Git/filesystem transaction. No post-check ref-stability guarantee, Git metadata locks, or cooperative-writer assumption is provided. The caller supplies the destination; artifact identity derives from canonical content. Existing publication remains content-addressed, exclusive, and atomic at its filesystem boundary; an existing artifact MUST NOT be overwritten.

#### Scenario: Moving reference stabilizes
- GIVEN a symbolic reference changes before a revalidation attempt then stabilizes within the bound
- WHEN publication is retried
- THEN only evidence for the stabilized immutable identity is published

#### Scenario: Publication conflict or persistent movement
- GIVEN the destination exists or the reference remains unstable beyond the bound
- WHEN publication is attempted
- THEN no artifact is overwritten or partially published and a typed technical outcome is returned

#### Scenario: Canonical content publication identity
- GIVEN canonical content is selected for publication
- WHEN it is published below the caller-supplied root
- THEN the caller-supplied publication root is passed unchanged, the relative artifact identity derives from canonical content, and any existing artifact remains unchanged
