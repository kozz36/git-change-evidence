# Proposed selected-path W1/W2 CLI contract

> Historical design record retained below. For the local source implementation and unfinished qualification, see [current progress](apply-progress.md). This is not a release or adoption claim.

**DRAFT · PROPOSED · NOT IMPLEMENTED.** Decision proposed: add a staged `gce w1w2` shell over existing versioned neutral contracts, with caller-selected paths and tool-owned safe content acquisition. This is technical documentation, not implementation authorization or an availability claim.

## Selected scope

1. `inventory` acquires only explicitly selected raw relative paths beneath an authorized absolute source root and constructs Inventory V1.
2. `result` acquires the entire immutable committed base/head range, independently of selected inventory paths, and constructs Accounting Result V1.
3. `report` validates retained policy/inventory/result antecedents and constructs Report V1.

The caller need not supply file contents. A successful empty selection proves only that the requested selection was empty and verified; it does not establish complete inventory, repository emptiness, or complete untracked coverage. Consumers requiring full-untracked coverage must establish it separately.

## Motivation and ownership

Existing pure Go constructors accept typed inputs and produce canonical evidence. The proposed imperative shell supplies safe acquisition, document transport, resource bounds and diagnostics; it does not change neutral meaning. Authorization and confidentiality decisions remain with the caller/harness. Selection is not permission proof. No evidence output approves, rejects, gates, blocks, merges or releases anything, or assigns delivery authority.

## Non-goals

No automatic Git-status/ignore discovery, recursive enumeration, consumer implementation, consumer configuration syntax, profile/plugin loader, executable schemas, core schema change, or silent broadening of legacy `project`/`publish` inputs. Supported neutral Policy V1 does not promise every project-specific adapter semantic.

No signing, installation, distribution, telemetry, multilingual census, release, publication, consumption pin or repository cleanup. Planned source/dev identity `v0.1.0-preview.1` is not a published tag, package or release. No current tests, functional parity, review outcome or implementation approval is claimed.

## Reading path and readiness

- [Design and public source mapping](design.md): grammar, acquisition, streams and the owner-selected initial Linux resource/threat profile, not runtime support.
- [Proposed requirements](specs/w1-w2-cli/spec.md): behavior and negative scenarios, R1–R6.
- [Conditional work units](tasks.md): entirely future, not apply-ready.
- [Implementation boundaries](apply-progress.md): public package review, planned work units and verification scope; private execution records are not part of this contract.

The owner selected an initial Linux-only profile for this new family: caller-controlled local roots/stable mount namespace, bounded selected-path acquisition, detected-race failures but no atomic snapshot, physical hardlink-provenance certificate, hard RSS or blocked-OS/stdio wall-time guarantee. See the [exact limits](design.md#initial-linux-resource-profile). The initial selection authorized specification updates only; local implementation has since been separately authorized. Platform qualification, calibration and executable conformance remain unperformed, and no publication or adoption is implied. Existing Go APIs, census and legacy CLI behavior are not restricted by these new-family limits.
