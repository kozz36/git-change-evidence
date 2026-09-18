# Publication path classification — `v0.1.0-preview.1`

## Purpose and bounded result

This evidence-only record classifies the local-path pattern results in the
candidate range from `origin/main@1a431c576b277496152a3781c421d412663064e2`
through correction commit `1e655a46f04a993017a446b57548199a2d405edd`.
It creates no owner, product, legal, publication, visibility, release, or
history-rewrite decision.

Read this with the [publication audit](publication-audit-v0.1.0-preview.1.md)
and the [publication disposition](publication-disposition-v0.1.0-preview.1.md).
It binds F-01 to the owner disposition **`Aceptar exposición`** / **Accepted
limitation**: current and historical local-environment path exposure is
accepted, history is preserved, and no rewrite is authorized.

**Bound scan date:** 2026-09-18. **Candidate pre-write context:** `HEAD` was
`1e655a46f04a993017a446b57548199a2d405edd`, `origin/main` resolved to the
bound base above, and `git status --porcelain=v1` was empty before this artifact
was created. The input range contains 15 commits.

## Result at a glance

| Surface | Absolute Unix/macOS-path occurrences | Windows/additional root-home occurrences | Recognized credential/private-key/token occurrences |
| --- | ---: | ---: | ---: |
| Changed-file bytes | 35 | 0 | 0 |
| Fifteen commit patches | 84 | 0 | 0 |

The scan reports pattern occurrences, not unique values. It retains no matched
value, username, or local directory string. Every digest below is the full
SHA-256 of one matched byte sequence, computed immediately after matching.

## Exhaustive source census and classification

| Source file | Changed-file count | Patch count | Lifecycle category | Semantics | Sensitivity assessment | Acceptance binding |
| --- | ---: | ---: | --- | --- | --- | --- |
| `.atl/skill-registry.md` | 34 | 68 removed | deleted generated registry | generated `.atl` registry install locations | Local-environment structure can be inferred; deleted metadata is not runtime input in this candidate. | F-01 Accepted limitation; preserve history, no rewrite. |
| `openspec/changes/archive/2026-09-17-public-preview-readiness/apply-progress.md` | 1 | 0 | current historical archive record | historical apply-progress evidence | One retained documentary local-environment reference; no executable behavior. | F-01 Accepted limitation. |
| `openspec/changes/archive/2026-09-18-go-ast-census/archive-report.md` | 0 | 1 added | historical added archive/verification evidence | SDD/worktree/action-context/tool provenance | Historical verification context can disclose environment structure. | F-01 Accepted limitation; preserve history. |
| `openspec/changes/archive/2026-09-18-go-ast-census/verify-report.failed-docs.md` | 0 | 6 added | historical added archive/verification evidence | SDD/worktree/action-context/tool provenance | Historical verification context can disclose environment structure. | F-01 Accepted limitation; preserve history. |
| `openspec/changes/archive/2026-09-18-go-ast-census/verify-report.md` | 0 | 7 added | historical added archive/verification evidence | SDD/worktree/action-context/tool provenance | Historical verification context can disclose environment structure. | F-01 Accepted limitation; preserve history. |
| `openspec/changes/go-ast-census/verify-report.md` | 0 | 2 removed | historical removed active evidence | SDD/worktree/action-context/tool provenance | Removed active evidence remains observable in patch history only. | F-01 Accepted limitation; preserve history. |
| **Total** | **35** | **84** | — | — | — | — |

All 35 changed-file occurrences are either deleted generated registry metadata
(34) or a current historical archive record (1). All 84 patch occurrences are
historical add/remove evidence: 68 generated-registry removals, 14 archived
verification additions, and 2 removals from active verification evidence.

## Value-free SHA-256 digest census

The group labels below are only shorthand for the full SHA-256 values in the
registry. `×N` is the multiplicity of that value in the named surface/source.
The accounting is exhaustive and value-free.

| Surface/source and action | Digest groups with multiplicity | Occurrences |
| --- | --- | ---: |
| Changed-file: `.atl/skill-registry.md`, base-deleted | `C01`–`C34` ×1 each | 34 |
| Changed-file: archived apply-progress record, current | `C35` ×1 | 1 |
| Patch: `.atl/skill-registry.md`, removed | `C01`–`C34` ×2 each | 68 |
| Patch: archived archive report, added | `P01` ×1 | 1 |
| Patch: archived failed-docs verification, added | `P01`–`P06` ×1 each | 6 |
| Patch: archived verification report, added | `P01`–`P03`, `P07`–`P10` ×1 each | 7 |
| Patch: active verification report, removed | `P11`–`P12` ×1 each | 2 |
| **Reconciled totals** | **35 changed-file; 84 patch** | **119** |

### Digest registry

```text
C01=09b0b0e73ad0c254959f4a0acf47bb19ebdb5a7b90eb05b114ae7f8a07850a2d
C02=13ad506167d2bc963f1594315d5498e1e6d87ad5c45a5713571e9035e1fc7982
C03=30dfa4e17371d185e8282065ff88fb18a2bf5c8dc7c78ccf8e6c1f896d61dc49
C04=39d3e822463acce1ab97b3b77bf924ee09614488874feb0c72a6588674ee65d0
C05=40faaa5b75acbbae8257c78b58a92983a786425b66514427977f9747d4486279
C06=418f63e3b2023f938142a8f53fb29828fce1a9fb0db021ae368814c57411669e
C07=452070606bf64bdd91f7d44f66548cf9a49536781ca8fd670675b08db970f2a6
C08=4a4aca4af9f41745487cfd8aab8946a11119467b4ee459f9937eaf663081bec1
C09=513fa7c82da6ef78c893b4a3ec07c2b44c28908f1eb1100067a90ace2105c3fc
C10=697baebbc1301a41d468a3d9972ea66043903c20e47c3a4a639b52efff6ff53c
C11=72c9d3de9d35809bd8ccf67e348ea9f1164743941397168608b3ce1ec035023d
C12=73441fe579197152e87d279bde2a7804ec59184f46f34dc05fa207a5e1a0c55f
C13=7d161d4e75c4e417e6f4e871e46c8d81e76a72e4f08d71e70bacad6fede1aad2
C14=834a2858e068e574b489edcf5c38fb2fb535aa70607cdd37082b56ad82fcbb8b
C15=89f5dad7ce3945c4312a6b5419f64e950fde9b569903527411ca7df4cf8ebefb
C16=8a0650a4fa00041d51ab81730165b94b9a119b580758a1db8536d69b8b9d6378
C17=8adbc936e35b871a99df069f1fbc3d4ce38b0ff42115bafec3b72d5a174076c6
C18=8e363f84c96eea3ce0d9e2b42c002dfebe4f97af868e02f33c1a450daad02f87
C19=a16ca09812b6cdc41f82709d898a5378d6d99c86b6522a54213e1f096bd6d97a
C20=a2d427fade828f2774bb74668e5424fb350a22bde298eb08c38b600b3b386fc5
C21=a2f9547c2aa883a1bacb97042f22c2dd6f4a7e553041a4fc40fe304490be0393
C22=a459e002c034e82801e86e46fa85c268de7bc69e2e230e59d85804aec9a7691a
C23=a6f3244dcbb3afdd244450b80a6a991202dc06128c3eb4d4be802c33a74095ab
C24=b91fb1c41a608b504c69ef89ba166cfe923b6159c74a8c6dba5e58aa9af62da8
C25=bab2064ff94d577486f4ab775924a719e9d69f4287b65d6b161346e1ea9f950d
C26=bbace5795ae00645f3b9868166dc761d317de9ee3d8664b61c61ec121cb725b3
C27=c0802843e8d02ada666fab3a27091d1186d38b076f256253361faf96d155c345
C28=c6c61a415860dec207e2b23bd49c812105053c0b381708684f9debe5766d00af
C29=d40ddd9e781e5eb6553f4a49d5706e7327cddc7148f66cf64c1a66772d9bc4e4
C30=dd136607060b8a8ffb261fea613bcc0a69a369324a748ffcb19cbf6f7e5b0b4a
C31=df4e4d4e430d6cfadb9d686d259c6a7678f7c77ccb9bf85662a8de13d4e96203
C32=e7c8effc3bf36801bcc9b9c7bf0e84d697a10c721b6b2c9c487325838024ad1f
C33=edf9a6159182403a6d251b79e656c9a76c99d37a057e370d5d9fbcc51264a2ad
C34=f7a04e78c7d9ef71a8610588c59c80837a8419f73dfb7ec4028d8904470a6516
C35=b0d2ad20d727caefbf17d474817f4216c770e59d21669535fccfbd900ffe49de
P01=88ca17bd91bfcf9a29c5b5267beb65cdb09d5b28fe2a7d4abd7a3c90e64e7177
P02=237e9398f79a8d6fd43ff82c5c025e9812f99a16a225fe485be70871a819efe0
P03=4d874f36ba484c9d07383cb77ab2d3759e614c363b0b006b32f8193d4082028c
P04=7f1790ebd7b1f9dc611776835742b3ada80480936f6ed7ab90ce20ec14b64647
P05=8b174660f6beda586a948015d1028ea079e4f960f53bcf56325f8be67e9db7a8
P06=94cc84699341dc2fdbca4b73ff1f8c9d4bb50e0b0b7783ecea8f94ef93a5f5ac
P07=021084e996906ae6b8d1fc514a18b57241d0502e7010e0f1918002d25b537c20
P08=3432971cb59f9b89c7b60db360df563b03ed2fa7330218ead4950a5afb04ceef
P09=4284316be4b06e7c86413ac7a1f140688842c036126fa122fae5ecc9960a172a
P10=4cded6f6a7d5d5baca20abe1cfa984af73b57b86f919ba7357979498f2deca2b
P11=06da4ddc3cecaadef15099fac3e0098fce07d5d685b0a53d934deb9b68321d10
P12=0dae4e9c34434849a3fc2185e27186aac73f3cd1c97cb477b54081d145f349fe
```

There are 35 changed-file digest groups and 46 patch digest groups. The patch
registry reuses `C01`–`C34` because the same deleted registry values occur
again in merge-parent patch views; it reuses `P01`–`P03` across archival reports
when the same historical value appears in more than one patch source.

## Reproducible safe scanner

Run the following command from the repository root. It is read-only: Git object
reads and patch generation occur in memory, and the program prints only
repository source path, surface/action, occurrence count, digest multiplicity,
and recognized-sensitive category counts. It never prints or writes a match.
The hexadecimal root fragments construct the canonical Unix/macOS and Windows
patterns without embedding a local directory string in this record.

```bash
python3 - <<'PY'
import collections, hashlib, re, subprocess
base = '1a431c576b277496152a3781c421d412663064e2'
head = '1e655a46f04a993017a446b57548199a2d405edd'

def git(*args, text=False):
    return subprocess.check_output(('git',) + args, text=text)

def literal(hex_text):
    return bytes.fromhex(hex_text)

assert git('rev-parse', 'origin/main', text=True).strip() == base
assert git('rev-parse', 'HEAD', text=True).strip() == head
path_roots = b'|'.join(map(literal, (
    '686f6d65', '5573657273', '64617461', '746d70', '7661722f746d70')))
windows_roots = b'|'.join(map(literal, (
    '5573657273', '446f63756d656e747320616e642053657474696e6773', '686f6d65')))
primary_path_rules = {
    'unix_macos': re.compile(rb'(?<![A-Za-z0-9_])/(?:' + path_roots + rb')/[^\s"\']+'),
}
secondary_path_rules = {
    'windows': re.compile(rb'(?i)\b[A-Z]:\\(?:' + windows_roots + rb')\\[^\s"\']+'),
    'root_home': re.compile(rb'(?<![A-Za-z0-9_])/' + literal('726f6f74') + rb'/[^\s"\']+'),
}
secret_rules = {
    'private_key_marker': re.compile(rb'-----BEGIN [A-Z0-9 ]*PRIVATE KEY-----'),
    'aws_access_key_id': re.compile(rb'\b(?:AKIA|ASIA)[A-Z0-9]{16}\b'),
    'github_token_shape': re.compile(rb'\b(?:ghp_[A-Za-z0-9]{36}|github_pat_[A-Za-z0-9_]{20,})\b'),
    'gitlab_token_shape': re.compile(rb'\bglpat-[A-Za-z0-9_-]{20,}\b'),
    'slack_token_shape': re.compile(rb'\bxox[baprs]-[A-Za-z0-9-]{10,}\b'),
    'credential_assignment': re.compile(rb'(?i)\b(?:api[_-]?key|secret|token|password|passwd|client[_-]?secret|access[_-]?key)\b\s*[:=]\s*["\']?(?!\$\{\{|\$[A-Z_]+)[^\s"\']{8,}'),
    'credential_uri': re.compile(rb'(?i)\b(?:https?|postgres(?:ql)?|mysql|mongodb)://[^\s/@:]+:[^\s/@]+@'),
}
records = collections.defaultdict(collections.Counter)
secondary_path_counts = collections.Counter()
secret_counts = collections.Counter()

def absorb(surface, path, action, data):
    bucket = records[(surface, path, action)]
    for rule in primary_path_rules.values():
        for match in rule.finditer(data):
            bucket[hashlib.sha256(match.group(0)).hexdigest()] += 1
    for category, rule in secondary_path_rules.items():
        secondary_path_counts[(surface, category)] += sum(1 for _ in rule.finditer(data))
    for category, rule in secret_rules.items():
        secret_counts[(surface, category)] += sum(1 for _ in rule.finditer(data))

def blob(revision, path):
    return git('show', revision + ':' + path)

for row in git('diff', '--name-status', 'origin/main..HEAD', text=True).splitlines():
    status, path = row.split('\t', 1)
    if status.startswith('D'):
        absorb('changed-file', path, 'base-deleted', blob(base, path))
    else:
        absorb('changed-file', path, 'current', blob(head, path))

commits = git('rev-list', '--reverse', 'origin/main..HEAD', text=True).splitlines()
assert len(commits) == 15
for commit in commits:
    patch = git('show', '-m', '--format=', '--no-ext-diff', '--unified=0', commit)
    path = None
    for line in patch.splitlines():
        if line.startswith(b'diff --git a/'):
            path = line.split(b' ')[2][2:].decode('utf-8', 'surrogateescape')
        elif path and line.startswith(b'+') and not line.startswith(b'+++'):
            absorb('patch', path, 'added', line[1:])
        elif path and line.startswith(b'-') and not line.startswith(b'---'):
            absorb('patch', path, 'removed', line[1:])

for (surface, path, action), digests in sorted(records.items()):
    count = sum(digests.values())
    if count:
        groups = ','.join('%s:%d' % pair for pair in sorted(digests.items()))
        print(surface, path, action, count, groups)
for surface in ('changed-file', 'patch'):
    secondary = sum(secondary_path_counts[(surface, category)] for category in secondary_path_rules)
    print(surface, 'path_windows_or_additional_root_home', secondary)
    for category in sorted(secret_rules):
        print(surface, category, secret_counts[(surface, category)])
print('changed-file_total', sum(sum(v.values()) for k, v in records.items() if k[0] == 'changed-file'))
print('patch_total', sum(sum(v.values()) for k, v in records.items() if k[0] == 'patch'))
PY
```

### Expected safe output accounting

The non-zero path rows reconcile exactly to the six rows in the source census:
changed-file `34 + 1 = 35`; patch `68 + 1 + 6 + 7 + 2 = 84`. The command prints
zero for the Windows/additional-root-home aggregate and zero for each recognized
private-key, cloud-key, token-shape, credential-assignment, and credential-URI
category on both surfaces. It uses `-m` for the one merge commit so patch
occurrences are counted per parent comparison; that is why the deleted registry
has 68 patch removals while its selected base blob has 34 matches.

## Interpretation and boundaries

- **Current versus deleted versus historical.** The current candidate retains
  one historical apply-progress reference. The registry is deleted in the
  candidate but remains in the base blob and patch history. Archive and removed
  verification records are documentary history, not current runtime inputs.
- **Patch semantics.** A patch count is an occurrence in an added or removed
  hunk line, including a merge-parent comparison; it is not a count of distinct
  values or currently reachable bytes. Repeated values therefore have separate
  occurrence multiplicities.
- **Sensitivity.** Recognized secrets are zero in both bounded surfaces. The
  accepted path disclosures can still reveal local-environment structure, but
  their retention is bound to the stated F-01 owner acceptance.
- **Runtime effect.** No runtime portability behavior is introduced: the
  occurrences are deleted registry metadata or documentary/historical evidence.
  This artifact adds no production or runtime path behavior.
- **Authority limit.** This is not legal, credential, confidentiality, privacy,
  or publication clearance beyond the existing disposition. It does not alter
  audit, archive, or history bytes; it only records a classification of them.

## Verification record

The read-only reproduction was run against the bound base and correction commit
until it returned the required 35 changed-file and 84 patch occurrences. Local
links in this artifact target the two existing publication records. After the
artifact is written, verify its whitespace, links, safe-content constraints,
line count, digest, and single-path worktree delta without changing history.
