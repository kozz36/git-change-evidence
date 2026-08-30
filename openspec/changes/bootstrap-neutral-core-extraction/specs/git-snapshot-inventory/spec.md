# Git Snapshot Inventory Specification

## Purpose

Produce reproducible committed-change evidence and scoped untracked inventory.

## Requirements

### Requirement: Native immutable snapshot identity

The system MUST resolve each revision to its immutable native Git object ID before deriving committed evidence. It MUST preserve SHA-1 and SHA-256 identities without rehashing, truncation, or conversion, and MUST derive committed changes only from those immutable objects.

#### Scenario: SHA-1 snapshot
- GIVEN a repository whose resolved revisions use SHA-1 IDs
- WHEN committed evidence is derived
- THEN the report retains the native SHA-1 IDs unchanged

#### Scenario: SHA-256 snapshot
- GIVEN a repository whose resolved revisions use SHA-256 IDs
- WHEN committed evidence is derived
- THEN the report retains the native SHA-256 IDs unchanged

#### Scenario: Unresolved or invalid object identity
- GIVEN a selected revision is unresolved, missing, invalid, or foreign to the repository
- WHEN committed evidence is requested
- THEN a typed resolution failure is returned and no snapshot, entries, or totals are invented

#### Scenario: Repository race during acquisition
- GIVEN a required immutable object becomes unavailable or changes during committed-snapshot acquisition
- WHEN committed evidence is requested
- THEN the system fails safely with typed race evidence and emits no partial or invented snapshot

### Requirement: Faithful and isolated acquisition

The system MUST preserve paths as bytes and retain entry file kinds and modes, including renames, binaries, symlinks, and gitlinks. It MUST invoke Git only with explicit argument vectors in a closed environment, and MUST disable ambient diff, text-conversion, replacement-reference, and inherited-configuration behavior that could alter evidence.

#### Scenario: Special entries and paths
- GIVEN immutable changes include a byte-preserving path, rename, binary, symlink, and gitlink
- WHEN snapshot entries are acquired
- THEN every entry retains its native path bytes, kind, mode, and applicable rename identity

#### Scenario: Ambient Git behavior is present
- GIVEN ambient Git configuration, diff drivers, text conversion, or replacement references would alter a normal Git result
- WHEN committed evidence is acquired
- THEN acquisition uses only its explicit arguments and controlled environment

### Requirement: Separate untracked inventory

The system MUST keep an untracked-file inventory separate from committed-change accounting. Index and working-tree state MUST NOT alter committed evidence; unavailable inventory data MUST be reported explicitly rather than inferred.

#### Scenario: Untracked file beside committed change
- GIVEN an immutable revision pair and an untracked working-tree file
- WHEN evidence is generated with inventory requested
- THEN committed accounting excludes the untracked file and inventory records it separately

#### Scenario: Inventory unavailable
- GIVEN inventory acquisition cannot be completed
- WHEN evidence is generated
- THEN the system reports typed unavailable inventory evidence without changing committed totals

#### Scenario: Staged or dirty tracked state
- GIVEN the selected immutable revision pair and tracked files with staged or unstaged changes
- WHEN committed evidence is generated
- THEN snapshot entries and committed accounting remain identical to a clean working tree

### Requirement: Safe untracked traversal

The system MUST traverse requested untracked inventory descriptor-relatively without following links. It MUST prevent intermediate or final symlink escape and replacement races from including data outside the requested root or inventing inventory evidence.

#### Scenario: Symlink escape or replacement race
- GIVEN untracked traversal encounters an intermediate or final symlink escape, or a path is replaced during traversal
- WHEN inventory is generated
- THEN no escaped content is inventoried and the affected result is safely omitted or typed unavailable

### Requirement: Real Git verification

The system MUST verify snapshot and inventory behavior against real temporary Git repositories, not mocks, for SHA-1 and SHA-256 formats, hostile paths, file kinds, missing objects, dirty-state isolation, and relevant repository or traversal races.

#### Scenario: Real repository coverage
- GIVEN temporary repositories containing each required format and adversarial condition
- WHEN the verification suite runs
- THEN it exercises real Git behavior for every listed condition without substituting mocks
