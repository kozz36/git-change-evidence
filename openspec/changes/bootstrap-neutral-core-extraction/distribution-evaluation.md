# Pinned-Git provenance evaluation — no distribution commitment

**Selected direction: evaluate an immutable Git SHA for reproducible source
provenance without a package registry.**
Access and API/CLI consumption remain future consumer-SDD decisions.
This is evaluation-only provenance, not adoption or delivery approval.
No concrete SHA, tag, consumption pin, mechanism, or publication contract is
selected; a development revision is not a consumption pin.
Preparing this evaluation does not require consumer adoption or creation of its
SDD. A separately human-approved consumer SDD remains a future adoption prerequisite.

## Why evaluate Git SHA provenance?

An immutable commit identifier identifies source content and known checkout inputs
without requiring a package registry. Unlike a moving branch or tag, a recorded
commit identifier supports comparison against the same identified source revision.
Repository access and continued object availability remain prerequisites; identity
alone does not ensure retrievability.
A source pin does not freeze transitive dependencies, toolchains, or the build
environment. It does not automatically establish security, support, license
compliance, completeness, liveness, or consumer parity.
Revision ownership, change handling, and reproducibility verification remain future
consumer requirements, not validated packaging.
No network experiment, package infrastructure, or consumer consumption route is
implemented or validated by this evaluation.

## Available surfaces and evidence limits

The neutral root Go API and CLI are available technical surfaces, not a selected
consumer boundary. `internal/cnsicprofile.DecodePolicyV1` is NONPUBLIC: it is not
a consumer import API. No bridge, public profile exposure, or consumer-owned
adapter implementation is selected.
Existing neutral acceptance of 40/64-character Git object identities is not proof
that a specific consumer route works.
The existing frozen standalone W2 profile and real SHA-1/SHA-256 Git evidence
cover **change-only evidence with empty inventory**. They do not establish general
CNSIC parity, full-tree inventory, consumer parity, cutover, or delivery authority.
These limits describe prior technical work, not new verification results from
preparing this evaluation. Technical parity is not human transition approval.

## Alternatives considered — declined for now, not forever

| Option | Why not selected now |
|---|---|
| Versioned Go module | Tags and version/release policy need decisions; the Go-consumer boundary is undefined and the internal profile stays private. |
| Versioned CLI binaries with checksums | Platform coverage, builds, distribution, and maintenance have not been approved. |
| Defer all distribution evaluation | No implementation now remains true, but Git-SHA evaluation is the selected direction rather than postponing evaluation. |

Module and binary options remain later choices if Git consumption proves
unsuitable; neither is declared unsupported.

## Future consumer decisions — all UNDECIDED

| Decision | Required future human-owned resolution |
|---|---|
| Access and credentials | Repository access and credential requirements, without credential values here. |
| Consumption boundary | API versus CLI, integration mechanism, and consumer-owned adapter responsibilities. |
| Concrete revision pin | Revision ownership, pin selection, change handling, and reproducibility verification. |
| Actual distribution | Packaging, distribution, release policy, license obligations, platforms, and support if later chosen. |
| Owner approval and observations | Consumer parity/acceptance criteria, observation owner and duration, reversal verification, deprecation approval, and rollback-window closure. |

- [ ] Obtain separate human approval for the consumer SDD, its boundary, access, and ownership.
- [ ] Resolve the undecided pin and consumption mechanism; select an alternative only through a later human decision if needed.
- [ ] Define and observe consumer parity and rollback criteria before requesting any reversible cutover.
- [ ] Resolve actual distribution approval separately; no default owner, duration, pin, or environment is assigned here.

The [consumer handoff](cnsic-handoff.md#1-dependency-integration) retains the
separate dependency-integration, reversible-cutover, owner-observation, and
conditional local-code-deletion stages. Completion of documentary planning does
not complete these future stages or authorize adoption.

## Next decision and separate delivery authority

The next decision belongs to a separately human-approved consumer SDD and actual
distribution proposal, not this evaluation.
No packaging, publication, release, integration, credentials, runtime change,
cutover, or local deletion is performed or authorized here.
Push, PR creation, merge, and release each require distinct delivery authorization;
this document grants no operational or shipping authority. GCE remains evidence-only.
