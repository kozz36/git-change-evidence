# Publication audit — `v0.1.0-preview.1`

**Audit kind:** evidence-only, read-only technical audit plus this documentation record. It grants no approval, clearance, publication, visibility, release, remediation, or history-change authority.

## Executive summary

This report binds its observations to local branch `audit/v0.1.0-preview.1`,
commit `2536397213efd973b1dd83f2c96d683a9b58a9f8`, and tree
`2cba92cb6c1f2b1bc6be411f3c55a6e22e3d0f58`. Before evidence collection, `HEAD`,
the index tree, and the working tree matched that commit/tree and `git status
--porcelain=v1` was empty. The audit examined all 83 commits and 601 reachable
objects from all local refs, the current tree, and read-only remote/GitHub
surfaces available to the authenticated account.

No credential-shaped pattern, private-key marker, or binary blob was detected
by the bounded fallback scan. This is **not credential clearance**: specialist
secret scanners and their rulepacks were unavailable, and the fallback has
explicit exclusions below. The scan did find local-path, email-shaped,
confidentiality-word, and high-entropy candidates. Their meanings have not been
determined and no candidate values are recorded here.

The GitHub repository metadata endpoint reported a private repository, 0 remote
tags, and 0 remote releases. The candidate tree contains no release-like asset
path and no release/publish/deploy workflow. Several hosting-control and
security-analysis endpoints were denied or unsupported, so their state is not
attested.

### Findings census

| Category | Count | Evidence-only result |
| --- | ---: | --- |
| Confirmed credential/private-key patterns | 0 | No match in the fallback scan; not a clearance. |
| Confirmed binary blobs | 0 | 265 reachable blobs and 182 distinct current-tree blobs were NUL-free. |
| Confirmed local-path indicators | 60 occurrences / 16 blobs | Potential local-environment disclosure; 10 blobs are in the candidate tree. |
| Email-shaped indicators | 9 occurrences / 9 blobs | Potential PII; 3 blobs are in the candidate tree. |
| Commit identity metadata | 83 commits | 2 distinct email hashes and 4 distinct name hashes. |
| Confidentiality/proprietary-word indicators | 4 occurrences / 3 current blobs | Heuristic context review required. |
| High-entropy candidates | 1,996 occurrences / 109 blobs | Heuristic only; no value was retained. |
| Legal/dependency evidence limitations | 1 | License selection is visible; third-party license metadata was not available offline. |
| Hosting-control/security-analysis unavailable checks | 9 | See [Unavailable checks](#unavailable-checks). |

## Finding register

The locators below are Git blob object IDs and one-way path digests, never
matched values. `current` means the blob is reachable from the bound candidate
tree; `history-only` means it is reachable only through another local ref.

| ID | Surface | Safe evidence locator | Classification / confidence | Publication relevance | Required owner disposition |
| --- | --- | --- | --- | --- | --- |
| F-01 | Source and documentation history | `absolute_unix_path`: 60 occurrences in 16 blobs (10 current, 6 history-only); sample locators `c6da622d…:15f30ca66fbb`, `27ace203…:de6a66c9cb10` | Confirmed pattern; semantic meaning unreviewed / high | May disclose local environment or reduce portability | Decide whether each current-tree occurrence is acceptable, must be removed under separate authority, or is historical evidence to retain. |
| F-02 | Source/history and commit metadata | Email-shaped: 9 occurrences in 9 blobs (3 current, 6 history-only); commit metadata: 83 records, 2 distinct email hashes, 4 distinct name hashes | Confirmed pattern/metadata; whether it is personal data is unreviewed / high | Potential privacy and attribution exposure if repository becomes public | Determine privacy/attribution treatment, including immutable-history implications; any rewrite or deletion needs separate authority. |
| F-03 | Current documentation | `confidentiality_marker`: 4 occurrences in 3 current blobs; locators `7f5f3375…:f6fd844e8941`, `3d06eb5a…:2f5cc76e4602`, `cc231e33…:5b01cb20e385` | Heuristic candidate / medium | Could describe confidential/proprietary status or merely discuss audit policy | Review context without treating the keyword hit as proof; classify cleared, accepted limitation, remediation required, or unresolved. |
| F-04 | All reachable blobs | 1,996 high-entropy base64-shaped candidates in 109 blobs (71 current, 38 history-only); by extension: `.go` 195, `.json` 4, `.md` 1,795, `.sum` 2 | Heuristic candidate; no credential-specific rule matched / low | Any real credential would be publication-relevant; the aggregate contains likely checksums/fixtures/document identifiers | Review with an approved specialist scanner or targeted owner review. Do not infer a credential from entropy alone. |
| F-05 | Dependency/license/NOTICE evidence | `go.mod` has one direct requirement, `golang.org/x/sys v0.38.0`; root `LICENSE` is Apache-2.0 text; candidate contains no `NOTICE`, `THIRD_PARTY`, `COPYING`, or vendored dependency file | Confirmed repository facts; license sufficiency unreviewed / high | Attribution and third-party licensing are owner/legal matters | Obtain owner/legal disposition; offline dependency-license/security metadata was not installed or fetched. |
| F-06 | GitHub controls and analysis | Rulesets and default-branch protection returned HTTP 403; code-scanning and Dependabot alerts returned HTTP 403; other unsupported endpoints are listed below | Confirmed access limitation / high | Cannot attest controls or alert state before publication | An authorized owner must supply/read the settings or accept the documented limitation; this report does not infer a disabled feature from denial. |

## Bound repository state and coverage

### Candidate and cleanliness

| Check | Observed result |
| --- | --- |
| Branch | `audit/v0.1.0-preview.1` |
| `HEAD` / requested candidate | `2536397213efd973b1dd83f2c96d683a9b58a9f8` / match |
| Candidate and `HEAD` tree | `2cba92cb6c1f2b1bc6be411f3c55a6e22e3d0f58` / match |
| Parent | `ba59a284ebdfa20bb100962aa0ffd8ca45dceac0` |
| Candidate subject | `docs: plan v0.1.0-preview.1 publication` |
| Pre-audit index tree | `2cba92cb6c1f2b1bc6be411f3c55a6e22e3d0f58` |
| Pre-audit index/worktree | clean; no staged, modified, or untracked path |
| Object connectivity | `git fsck --connectivity-only --no-dangling --no-reflogs` succeeded with no diagnostic |
| Candidate parent diff whitespace | `git diff --check <parent> <candidate>` passed; one changed path |

The clean identity above is the audit input identity. After this report and the
Task 1 evidence update are written, the index must remain at the candidate tree
and the worktree should differ only in the two authorized documentation paths;
that deliberate documentation delta is recorded by the post-write verification
command in [Verification recommendations](#verification-recommendations).

### Reachable-object and file census

| Surface | Observed result |
| --- | --- |
| Reachability scope | `git rev-list --objects --all` across every local ref |
| Commits / unique reachable objects | 83 / 601 |
| Reachable blobs / trees | 265 / 253 |
| Current tree files / distinct blobs / bytes | 185 / 182 / 1,456,451 |
| Reachable blob bytes | 2,890,533 |
| Binary classification | 0 NUL-containing reachable blobs; all 265 text-like by this limited test |
| Current size buckets | 17 below 1 KiB; 125 at 1–<10 KiB; 43 at 10–<100 KiB; 0 at or above 100 KiB |
| Largest current blobs | 86,435 `b0418a10…`; 73,854 `c6da622d…`; 44,217 `5ad84707…`; 30,446 `27ace203…`; 29,427 `d81e5cae…` |
| Largest reachable-history blobs | 86,435 `b0418a10…`; 73,854 `c6da622d…`; 71,730 `09e1eea4…`; 67,159 `12838ab9…`; 58,311 `965f4700…` |
| Generated/metadata files in current tree | GitHub workflows/templates/scripts, `go.mod`, `go.sum`, `openspec/config.yaml`, and compatibility JSON fixtures; no generated release asset was found |
| Release-like paths across reachable history | 0 for common archive/package/executable/SBOM extensions |
| Candidate release/publish/deploy workflows | 0 |

Current-tree file extensions were: 97 `.go`, 72 `.md`, 4 `.json`, 4 `.yml`, 2
`.mjs`, and one each of `.gitignore`, `.mod`, `.sha256`, `.sum`, `.yaml`, and
an extensionless `LICENSE`.

### Local ref, remote, and worktree inventory

There is one configured remote, `origin`, with sanitized host `github.com` and
no embedded URL credential detected. The local ref census had 30 heads, 21
remote-tracking refs (including `origin/HEAD`), and 0 tags. No other local ref
namespace was returned by `git for-each-ref`.

**Local heads (all):**

```text
audit/v0.1.0-preview.1
ci/7-go-pr-checks
docs/cnsic-profile-parity
docs/go-ast-census
docs/prepare-v0.1.0-preview.1
docs/public-preview-audit
docs/public-preview-entrypoint
docs/public-preview-governance
docs/public-preview-planning-design
docs/public-preview-planning-foundations
feat/11-policy-accounting
feat/34-policy-contract-v1
feat/5-immutable-git-snapshot
feat/9-untracked-inventory
feat/census-cli-agents
feat/inventory-v1
feat/public-preview-cli
feat/result-v1
feat/result-v1-pr10-constructor
feat/result-v1-pr13-entry-shape
feat/result-v1-pr14-observation-shape
feat/result-v1-unit7a-values
feat/result-v1-unit7b-decoder
feat/result-v1-unit8-report-binding
feat/result-v1-pr14b-revisions-totals-shape
feat/result-v1-pr14b-root-envelope-shape
main
test/result-v1-external-validation-pr83
test/result-v1-unit9-mutation-census
test/w2-frozen-compatibility
```

Read-only `git ls-remote --heads --tags origin` returned 20 heads and 0 tags:

```text
docs/go-ast-census
docs/public-preview-audit
docs/public-preview-entrypoint
docs/public-preview-governance
docs/public-preview-planning-design
docs/public-preview-planning-foundations
feat/census-cli-agents
feat/public-preview-cli
feat/result-v1-pr10-constructor
feat/result-v1-pr13-entry-shape
feat/result-v1-pr14-observation-shape
feat/result-v1-unit7a-values
feat/result-v1-unit7b-decoder
feat/result-v1-unit8-report-binding
feat/result-v1-pr14b-revisions-totals-shape
feat/result-v1-pr14b-root-envelope-shape
main
test/result-v1-external-validation-pr83
test/result-v1-unit9-mutation-census
test/w2-frozen-compatibility
```

The candidate/audit branch and several local heads are therefore locally visible
but not in the read-only remote head inventory. That is a publication-surface
observation, not authority to push, delete, or change any ref.

`git worktree list --porcelain` reported 24 worktrees: the main checkout,
`11-policy-accounting`, `34-policy-contract-v1`, `5-immutable-git-snapshot`,
`7-go-pr-checks`, `9-untracked-inventory`, `census-cli-agents`,
`cnsic-profile-parity`, `go-ast-census`, `inventory-v1`, `public-preview`,
`publication-audit`, `result-v1`, `result-v1-external-validation-pr83`,
`result-v1-pr10-constructor`, `result-v1-pr13-entry-shape`,
`result-v1-pr14-observation-shape`, `result-v1-pr14b-revisions-totals-shape`,
`result-v1-pr14b-root-envelope-shape`, `result-v1-unit7a-values`,
`result-v1-unit7b-decoder`, `result-v1-unit8-report-binding`,
`result-v1-unit9-mutation-census`, and `w2-frozen-compatibility`. The command
reported each attached to its corresponding local branch; no detached worktree
was reported.

## Secret, privacy, confidentiality, and path scan

Specialist scanners were inventoried first: `gitleaks`, `trufflehog`,
`detect-secrets`, and `semgrep` were unavailable. `rg` was available but was
not used to print matches. A bounded in-memory fallback scanned every reachable
blob once using private-key, cloud/token-shape, credential-assignment,
credential-URI, email, absolute Unix/Windows path, local-endpoint,
confidentiality-word, and base64 entropy patterns. It printed only aggregate
counts plus safe blob/path-digest locators.

The fallback scope included 265 reachable blobs: 182 in the candidate tree and
83 history-only. It did not decrypt, decode binary formats, invoke network
validation, use vendor rulepacks, inspect inaccessible objects, or retain,
print, hash, or store matched secret candidate values. It is consequently a
bounded safety check, not an exhaustive or specialist secret scan.

No private-key, AWS access-key, GitHub/GitLab/Slack token-shape,
credential-assignment, credential-URI, Windows-path, or local-network-endpoint
pattern matched. The local-path, email-shaped, confidentiality-word, and entropy
results are F-01 through F-04. In particular, entropy is a finding candidate,
not evidence of a credential. The `.sum` hits are consistent with dependency
integrity metadata but were not treated as a legal or security disposition.

## Fallback scanner reproducibility appendix

This executable recipe preserves the original regex literals, accounting, entropy threshold, and safe-output rules. It prints only aggregate counts and `blob_oid:path_sha256_12` locators, never matched values; it was not rerun for this correction.

```bash
python3 - <<'PY'
import subprocess,re,math,collections,hashlib
candidate='2536397213efd973b1dd83f2c96d683a9b58a9f8'; paths={}
for line in subprocess.check_output(['git','rev-list','--objects','--all'],text=True).splitlines():
    oid,*rest=line.split(' ',1); paths.setdefault(oid,rest[0] if rest else '')
ids='\n'.join(paths)+'\n'; types=subprocess.check_output(['git','cat-file','--batch-check=%(objectname) %(objecttype)'],input=ids,text=True).splitlines()
blob_ids=[x.split()[0] for x in types if x.split()[1]=='blob']
current={x.split()[2] for x in subprocess.check_output(['git','ls-tree','-r',candidate],text=True).splitlines()}
patterns={
'private_key_marker':re.compile(rb'-----BEGIN [A-Z0-9 ]*PRIVATE KEY-----'),'aws_access_key_id':re.compile(rb'\b(?:AKIA|ASIA)[A-Z0-9]{16}\b'),
'github_token_shape':re.compile(rb'\b(?:ghp_[A-Za-z0-9]{36}|github_pat_[A-Za-z0-9_]{20,})\b'),'gitlab_token_shape':re.compile(rb'\bglpat-[A-Za-z0-9_-]{20,}\b'),
'slack_token_shape':re.compile(rb'\bxox[baprs]-[A-Za-z0-9-]{10,}\b'),'credential_assignment':re.compile(rb'(?i)\b(?:api[_-]?key|secret|token|password|passwd|client[_-]?secret|access[_-]?key)\b\s*[:=]\s*["\']?(?!\$\{\{|\$[A-Z_]+)[^\s"\']{8,}'),
'credential_uri':re.compile(rb'(?i)\b(?:https?|postgres(?:ql)?|mysql|mongodb)://[^\s/@:]+:[^\s/@]+@'),'email_indicator':re.compile(rb'(?i)\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}\b'),
'absolute_unix_path':re.compile(rb'(?<![A-Za-z0-9_])/(?:home|Users|data|tmp|var/tmp)/[^\s"\']+'),'absolute_windows_path':re.compile(rb'(?i)\b[A-Z]:\\(?:Users|Documents and Settings|home)\\[^\s"\']+'),
'local_network_endpoint':re.compile(rb'(?i)\b(?:localhost|127\.0\.0\.1|0\.0\.0\.0)\b'),'confidentiality_marker':re.compile(rb'(?i)\b(?:confidential|proprietary|internal use only|do not distribute)\b')}
base64=re.compile(rb'(?<![A-Za-z0-9+/=])[A-Za-z0-9+/]{32,}={0,2}(?![A-Za-z0-9+/=])')
def entropy(b):
 c=collections.Counter(b); n=len(b); return -sum((v/n)*math.log2(v/n) for v in c.values())
results=collections.defaultdict(list)
for oid in blob_ids:
 data=subprocess.check_output(['git','cat-file','blob',oid])
 for cls,pat in patterns.items():
  n=len(pat.findall(data)); results[cls].append((oid,n)) if n else None
 n=sum(1 for token in base64.findall(data) if entropy(token)>=4.0)
 if n: results['high_entropy_base64_candidate'].append((oid,n))
for cls,rows in sorted(results.items()):
 total=sum(n for _,n in rows); now=sum(1 for oid,_ in rows if oid in current)
 print(cls,total,len(rows),now,len(rows)-now,','.join('%s:%s'%(oid,hashlib.sha256(paths[oid].encode()).hexdigest()[:12]) for oid,_ in rows[:20]))
PY
```

`--all` defines local-ref reachability; each unique object is type-filtered to a blob. A hit is current when its blob ID appears in `ls-tree -r <candidate>` and history-only otherwise. Counts are occurrences, matching blobs, current matching blobs, and history-only matching blobs. Exclusions: no decryption, binary-format parsing, network verification, specialist rulepacks, inaccessible objects, or retained/printed/hashed matched values.

## Dependencies, license, and notices

The bound `go.mod` declares module `github.com/kozz36/git-change-evidence`, Go
baseline `1.25.10`, and one direct dependency: `golang.org/x/sys v0.38.0`.
`go.sum` contains the two expected integrity entries for that module/version.
The root `LICENSE` contains Apache License 2.0 text. The candidate tree has no
root `NOTICE`, `THIRD_PARTY`, `COPYING`, or vendored dependency tree.

No dependency installation, `go mod download`, package-manager operation,
module-file mutation, or network license/security lookup was run. Local module
cache inspection was deliberately excluded. Automated inventory does not
resolve third-party attribution, copyright, trademark, export, confidentiality,
or other legal questions.

## CI and publication surfaces

The candidate has two GitHub workflows:

- `Go PR checks`: pull-request trigger to `main`; repository content permission
  is read-only; checkout persists no credentials; it runs formatting, `go test
  ./...`, and `go test -race ./...`.
- `PR policy`: `pull_request_target` trigger restricted to `main`; `contents`
  and `issues` permissions are read-only; it checks out the trusted base SHA,
  not the pull-request head.

No current workflow name/path indicated release, publish, or deploy behavior.
No release asset was created or inspected beyond the path census.

Read-only GitHub API observations for `repos/{owner}/{repo}`:

| Endpoint | HTTP status | Safe result |
| --- | ---: | --- |
| `GET /repos/{owner}/{repo}` | 200 | `visibility=private`, `private=true`, `archived=false`, `disabled=false`, default branch `main`, Pages flag `false`; `security_and_analysis` was not returned. |
| `GET /branches?per_page=100` | 200 | 20 branches returned. |
| `GET /tags?per_page=100` | 200 | 0 tags. |
| `GET /releases?per_page=100` | 200 | 0 releases. |
| `GET /actions/permissions` | 200 | Actions enabled; allowed actions `all`. |
| `GET /actions/permissions/workflow` | 200 | Default workflow permission `read`; pull-request review approval `false`. |
| `GET /actions/workflows?per_page=100` | 200 | 2 workflows. |
| `GET /actions/runs?per_page=100` | 200 | 146 total runs; newest 100 were completed: 71 success, 28 cancelled, 1 failure. The older 46 were not enumerated. |
| `GET /actions/secrets` | 200 | Secret count 0; no names or values requested. |
| `GET /actions/variables` | 200 | Variable count 0; no names or values requested. |
| `GET /collaborators?per_page=100` | 200 | 1 collaborator returned, role count `admin:1`; no identities recorded. |
| `GET /keys?per_page=100` | 200 | 0 deploy keys; 0 write-enabled keys. |
| `GET /environments?per_page=100` | 200 | 0 environments. |
| `GET /hooks?per_page=100` | 200 | 0 webhooks. |
| `GET /pages` | 404 | No Pages state inferred beyond the repository metadata flag. |
| `GET /packages?per_page=100` | 404 | Endpoint unsupported/unavailable; not proof of package absence. |

## Unavailable checks

| Surface | Endpoint/status | Limitation |
| --- | --- | --- |
| Rulesets | `GET /rulesets` — 403 | Ruleset state could not be read. |
| Default-branch protection | `GET /branches/main/protection` — 403 | Branch-protection state could not be read. |
| Code scanning | `GET /code-scanning/alerts?per_page=1` — 403 | Alert state could not be read. |
| Dependabot | `GET /dependabot/alerts?per_page=1` — 403 | Alert state could not be read. |
| Secret scanning | `GET /secret-scanning/alerts?per_page=1` — 404 | Feature/endpoint state unavailable; no conclusion inferred. |
| Dependabot vulnerability alerts | `GET /vulnerability-alerts` — 404 | Feature/endpoint state unavailable; no conclusion inferred. |
| Private vulnerability reporting | `GET /private-vulnerability-reporting` — 404 | Feature/endpoint state unavailable; no conclusion inferred. |
| Repository packages | `GET /packages?per_page=100` — 404 | Repository package inventory unavailable; no conclusion inferred. |
| Pages | `GET /pages` — 404 | Endpoint/feature state unavailable; the 404 was not treated as proof that Pages is disabled. |
| Specialist secret scanners | Local commands unavailable | Fallback scan has materially narrower coverage. |
| Offline dependency-license/security metadata | Deliberately not installed/fetched | No automated third-party legal/security conclusion. |

The repository metadata response omitted `security_and_analysis`. This is an
additional field-level limitation, not a tenth unavailable API endpoint; no
security-feature state is inferred from omission.

## Commands, tools, and safe result record

Tool versions: Git `2.55.0`; Go `go1.27.1-X:nodwarf5 linux/amd64`; GitHub CLI
`2.101.0 (2026-09-15)`; Python 3 was used for non-mutating aggregate scripts.
The repository target Go baseline is separately declared as 1.25.10; no build
or dependency install was run for this audit.

The following exact commands or safe command descriptions were used without
fetching/updating refs, altering settings, or writing repository data:

```text
git status --short
git branch --show-current
git rev-parse HEAD HEAD^{tree} <candidate> <candidate>^{tree}
git write-tree
git for-each-ref --format='%(refname) %(objecttype) %(objectname)'
git worktree list --porcelain
git remote
git remote get-url <remote>                 # host only was retained
git ls-remote --heads --tags origin
git rev-list --objects --all
git cat-file --batch-check=...
git ls-tree -rl -r <candidate>
git fsck --connectivity-only --no-dangling --no-reflogs
git diff --check <candidate-parent> <candidate>
git log --all --format=...                  # counts/hashes only retained
git log --all --format= --name-only         # release-like path count only
git show -s --format=... <candidate>
gh api -i repos/{owner}/{repo}[...read-only endpoints above]
```

The in-memory Python helpers consumed Git object streams and GitHub CLI JSON,
then emitted only counts, object IDs, path digests, endpoint/status, and allowed
metadata. They did not create files, modify refs, write the index, install
modules, or print candidate secret/PII values. The report's reproducibility
inputs are the branch, commit, tree, parent, object/reachability counts, and
command scope recorded above. Its final SHA-256 is intentionally computed only
after this file is complete.

## Owner-only decisions and next task

Task 2 owner disposition **may start**: Task 1 has complete evidence coverage
for the accessible local and remote surfaces, including explicit limitations.
Task 2 must classify F-01 through F-06 and every unavailable publication-relevant
check as cleared, accepted limitation, remediation required, or unresolved.

This report does not declare any finding cleared and does not make legal,
confidentiality, privacy, credential, history-rewrite, repository-policy, or
publication decisions. Under the program plan, unresolved publication-relevant
findings must stop downstream lifecycle/publication tasks; that is a plan
constraint, not an authority asserted by this report.

## Verification recommendations

Run these read-only checks before an independent verifier records Task 1:

```text
# Bind the audit artifact to the intended candidate and verify the only worktree delta.
git rev-parse HEAD HEAD^{tree}
git diff --check
git diff --name-only
git diff -- odd/tasks/publish-v0.1.0-preview.1.md docs/publication-audit-v0.1.0-preview.1.md
sha256sum docs/publication-audit-v0.1.0-preview.1.md

# Reconfirm local/remote publication surfaces without updating local refs.
git status --porcelain=v1
git ls-remote --heads --tags origin
gh api -i repos/{owner}/{repo}/tags?per_page=100
gh api -i repos/{owner}/{repo}/releases?per_page=100
```

A post-write local-relative Markdown-link check scanned 74 Markdown files and
21 local links. It reported three unresolved pre-existing targets outside this
report: two literal `...` references under `openspec/changes/go-ast-census/` and
one missing `docs/consumer-integration-guide.md` target from archived readiness
evidence. This report itself contributed no unresolved local file link.

The verifier should additionally inspect F-01 through F-06 with authorized
owner/legal/security access, confirm no matched value is exposed in this report,
and leave Publication Task 1 unchecked until that independent verification is
recorded.
