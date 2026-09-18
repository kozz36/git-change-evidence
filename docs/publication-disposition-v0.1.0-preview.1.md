# Publication finding disposition — `v0.1.0-preview.1`

**Outcome:** all Task 1 findings have an owner disposition. Task 3 may start,
but this record grants no issue, remediation, visibility, tag, release, upload,
or other publication authority. The third-party notice remediation and the
post-public security readback below remain hard dependencies for release
completion.

## Evidence binding

| Evidence | Bound value |
| --- | --- |
| Task 1 audit | `docs/publication-audit-v0.1.0-preview.1.md` |
| Audit SHA-256 | `3604c20cfdc124c444ec1b965036eea9e1b6d80b80f8c979cc768664941944e9` |
| Audited candidate | commit `2536397213efd973b1dd83f2c96d683a9b58a9f8`, tree `2cba92cb6c1f2b1bc6be411f3c55a6e22e3d0f58` |
| Disposition scope | F-01 through F-06 and the audit's unavailable checks |

The Task 1 audit is preserved unchanged. This document records later owner
classifications, targeted technical review, and hosting-control remediation.

## Disposition census

| Finding | Owner disposition | Evidence and consequence |
| --- | --- | --- |
| F-01 — local paths | **Accepted limitation** | Accept current and historical local-environment path exposure. Preserve history; no path removal or history rewrite is authorized. |
| F-02 — email and identity metadata | **Cleared / accepted limitation** | Clear email-shaped test values as synthetic fixtures. Accept public author and committer attribution, including the immutable-history exposure. |
| F-03 — confidentiality words | **Cleared** | Context review found policy and audit prose, not an actual confidential or proprietary marking. |
| F-04 — entropy candidates | **Cleared for this audit scope** | The bounded fallback review classified the candidates as identifiers, fixtures, path-like terms, schema values, or checksums. A pinned Gitleaks history scan found zero findings; details are below. |
| F-05 — dependency notices | **Remediation required** | Before release completion, add root `THIRD_PARTY_NOTICES` reproducing the `golang.org/x/sys v0.38.0` `LICENSE` and `PATENTS` obligations. This is a technical notice requirement, not legal clearance. |
| F-06 — GitHub controls | **Remediated / deferred by availability** | Private-repository controls were inspected and the available controls below were configured. Public-only controls remain a mandatory Task 15 readback and enablement step after visibility changes. |

No finding remains unclassified. “Cleared” here means the owner accepted the
recorded technical context; it is not legal, credential, confidentiality, or
publication clearance beyond the stated scope.

## F-04 specialist scan

The owner authorized one temporary, read-only, pinned Gitleaks scan of all
history reachable from local refs.

| Item | Result |
| --- | --- |
| Tool | Gitleaks `v8.30.1` |
| Source | Official `gitleaks/gitleaks` GitHub release |
| Linux archive SHA-256 | `551f6fc83ea457d62a0d98237cbad105af8d557003051f41f3e7ca7b3f2470eb` |
| Checksum verification | Archive matched the official release checksum manifest before execution |
| Scope | 84 unique reachable commits; 608 unique reachable objects; `git` scan with `--log-opts="--all"` |
| Safety | Default embedded rules, 100% redaction, report and cache confined to an owner-only temporary directory |
| Result | 0 findings; 0 finding categories |
| Cleanup | Tool, archive, checksum manifest, report, logs, cache, and temporary directory removed |
| State check | HEAD, HEAD tree, index tree, worktree status, refs, and reflogs unchanged |

The zero-finding result applies to that pinned tool, ruleset, and reachable
local-ref state. It does not promise that every possible scanner would find
nothing.

## F-05 required notice and asset contract

The local module cache for `golang.org/x/sys v0.38.0` contains its upstream
`LICENSE` and `PATENTS`. The future remediation must create a root
`THIRD_PARTY_NOTICES` that reproduces those obligations. Do not create that
file as part of Task 2; implement and review it as a separate work unit.

The external upload set remains exactly six unsigned files:

1. one deterministic source archive;
2. one Linux `amd64` archive;
3. one Linux `arm64` archive;
4. one CycloneDX JSON SBOM for Linux `amd64`;
5. one CycloneDX JSON SBOM for Linux `arm64`; and
6. `SHA256SUMS`, covering exactly the other five upload files.

Each Linux binary archive must contain exactly these four members:

```text
gce
git-change-evidence
LICENSE
THIRD_PARTY_NOTICES
```

The source archive must include the reviewed root `LICENSE` and
`THIRD_PARTY_NOTICES` from the final tagged source tree. No notice file is a
seventh upload asset.

## F-06 hosting-control evidence

Authenticated web inspection observed a private repository with no initial
ruleset or classic branch-protection rule, no packages, and Pages unavailable
while private. The owner then created classic protection for `main` with:

- pull requests required;
- required checks `Go checks` and `Validate PR policy`;
- conversation resolution required;
- no bypass;
- force pushes and branch deletion disabled; and
- approvals not required.

The rule read back as **Not enforced** while the repository remains private. It
must be read back as enforcing after public visibility is established.

Dependency Graph, Dependabot alerts, Dependabot security updates, and grouped
security updates were enabled and read back as enabled. Advanced Security,
private vulnerability reporting, secret scanning, and code scanning were not
available on the private repository and are deferred to Task 15. After the
visibility change, Task 15 must inspect and enable every available deferred
security control, verify the `main` protection is enforced, and record the
readback before release completion. An unavailable control must be reported as
such; it must not be inferred enabled or disabled.

## Remaining hard dependencies

- Add and independently review root `THIRD_PARTY_NOTICES` before the final
  release candidate can complete verification.
- Keep every Linux archive at exactly four members and the external upload set
  at exactly six unsigned files.
- After visibility changes in Task 15, verify branch-protection enforcement and
  enable/read back the newly available security controls before tag, asset,
  release, or completion work proceeds.
- Re-run the final audit and release-preparation verification against the exact
  frozen commit as required by the publication plan.

## Next step

Task 2 is complete because every finding is classified and required remediation
has an explicit bounded plan. Task 3 may create and seek approval for the
publication tracking issue. This does not start Task 3 and does not authorize
any later lifecycle or publication action.

## F-05 remediation completion — 2026-09-18

The separately bounded Task-9 remediation added root `THIRD_PARTY_NOTICES`.
It identifies `golang.org/x/sys v0.38.0` and clearly delimits verbatim copies
of the cached upstream `LICENSE` and `PATENTS` bytes. This records technical
notice evidence only; it does not provide legal clearance or change the
original F-05 finding/table status.

The current release-preparation integration is merge commit
`181fa986865d1f5a7d255db445bb70504b2d7753`, tree
`befc256d7b3c504f1472e2db4ff29ebb7d23e63d`, with parents
`1a431c576b277496152a3781c421d412663064e2` and
`6531661033d0be3fc56236fe83dc155dc7de080a`; no conflicts were observed.
The current Linux archive contract is exactly `gce`, `git-change-evidence`,
`LICENSE`, and `THIRD_PARTY_NOTICES`; the external set remains exactly six
unsigned assets. Final assets must be rebuilt from the final tagged commit.

This completion record does not create a tag, release, upload, package,
distribution channel, or public-visibility change, and does not authorize any
of them.
