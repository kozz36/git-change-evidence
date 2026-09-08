# W2 frozen compatibility harness

## Scope and provenance

This is test-only, supplemental compatibility evidence. It does not reopen or archive a
Generic Validation SDD, change production behavior, or make an approval, gate, or delivery
claim. The closed translator is deliberately not a general CNSIC profile or a W1-to-neutral
converter.

The imported corpus is the byte-frozen predecessor oracle from CNSIC commit
`a0fc7b26ff8a0e0a61baa586b32c46841611c806`. `TestImportedPredecessorCorpus` gates the
existing relative fixtures and their SHA-256 values:

- `accounting-v1.json`: `f265bc48750ccd68e41dcb64ec68134b9fb7db69655e2c215e39f8c6c8794c3f`
- `contracts-v1.json`: `84bba5b7b1471a1c12477de6a46a2f6f175fd09bdd175ecaf81c698902776ca5`
- `manifest.sha256`: `f8db149ef966e1eeab42e36e42e299ba4b197baeac22bc75cfb4b9cd55b8caac`

Those fixture bytes were unchanged before and after this work. The immutable `expected`
content is read as an oracle and is never regenerated or rewritten.

## W2 neutral observation

The external `changeevidence_test` harness reads only the frozen accounting fixture, validates
its closed five-added-entry W2 shape, and uses only this public chain:

```text
NewAccountingPolicyV1 -> NewInventoryV1(policy, nil) -> NewCommittedSnapshot -> NewAccountingResultV1
```

The translated policy is fixed to `openspec/changes/**`, `tests/**`, a validated literal
mechanical path derived from the `generated` annotation, and default `production`, in that
category order. Result categories are checked by each raw path against the fixture's
`expected.categories`, not against the package matcher or legacy `Account`; result ordering is
checked independently as normalized raw-path order.

The five entries produce known totals of artifacts 5, tests 3, mechanical 2, and production 4;
all known deletions are zero. The binary production entry remains a non-countable entry and
production total has `NonCountable: 1`. Consequently the neutral production-400 observation is
unavailable (`Available`, `Actual`, and `Exceeded` are all false/zero), while the separately
computed legacy known-text sum is 4 and `4 > 400` is false. This is semantic partial
compatibility, not a blanket compatibility pass.

A separate literal synthetic copy uses production known sum 401 plus the binary entry. It proves
that the legacy known-text comparison is true while the neutral observation remains unavailable
and false because non-countable evidence is retained. The fixture's side object IDs `a` through
`e` are validated as declared placeholders, not Git object IDs. The harness intentionally uses
its own distinct valid 40-hex base/head metadata and an empty inventory; it is not Git acquisition
proof.

## Closed-domain rejection coverage

`TestW2TranslatorRejectsClosedDomain` uses RawMessage-backed closed-object decoding so missing,
null, and wrong-typed values are not silently defaulted. Its 19 negative cases reject missing or
extra annotation keys; boolean, fractional, string, null, and negative counts; blank evidence;
bad basis; absent, duplicate, partial, binary, metacharacter, and noncanonical-base64 paths;
unsupported policy and status; a non-null old side; and side LOC/count disagreement. The harness
also requires canonical base64 raw paths, unique paths and placeholder IDs, regular-text known
counts, a whole-entry annotation, and a constrained policy/entry/status shape before any
public classification runs. This validation does not require every value to equal the frozen
fixture: exact corpus identity is enforced separately by its hash test and observation assertions.

This is intentionally not a full SD-7/CNSIC document validator. In particular, W1 canonical
contract hashes remain evidence from the frozen Python runner below; they are not neutral wire
parity and W1 is not converted into a neutral document.

## TDD and execution evidence

Strict TDD was active for the Go test surfaces.

1. **RED:** the smallest W2 behavior test first called undefined `w2FrozenResult`; the focused
   command failed with `undefined: w2FrozenResult` as intended.
2. **GREEN:** after adding the closed test-only translator and public chain, the focused W2 run
   passed. During implementation, the string-count negative case exposed an initial helper gap;
   adding the raw-number lexical guard made that focused run pass.
3. **TRIANGULATE:** the 19 rejection cases and the independent 401-plus-binary synthetic case
   passed, proving the text-only legacy comparison differs from neutral availability.
4. **REFACTOR:** the three files were formatted without production changes, then the focused W2
   and corpus command passed again.

The writer did not locate the installed Go 1.25.10 executable and used
`go1.27.1-X:nodwarf5` for the historical TDD cycle and commands below. Its initial claim that
Go 1.25.10 was unavailable was incorrect; the independent verification addendum records the
supported-toolchain rerun without relabeling that historical cycle.

```text
go test . -run '^(TestW2.*|TestImportedPredecessorCorpus)$' -count=1 -v
# PASS: imported corpus, two W2 observation tests, and all 19 negative subtests

go test ./... -count=1
# PASS: root, CLI, internal/git, internal/inventory, internal/publication

test -z "$(find . -path './.git' -prune -o -path './.codegraph' -prune -o -type f -name '*.go' -print0 | xargs -0 -r gofmt -l)"
# PASS: no unformatted Go files
```

The original predecessor runner was exported only from the cited Git blobs into a mode-0700
owned temporary root. Every copied file's Git blob ID was checked after export; no current source,
root `conftest`, pytest configuration, application import, or database bootstrap was copied. The
temporary root was removed on process exit. Its standalone execution intentionally used `/dev/null`
as pytest configuration and `--confcutdir=<frozen-root>`, excluding ancestor conftests.

```text
PYTHONDONTWRITEBYTECODE=1 PYTEST_DISABLE_PLUGIN_AUTOLOAD=1 PYTHONPATH=<frozen-root> \
  <frozen-venv-python> -m pytest -c /dev/null --confcutdir=<frozen-root> \
  -p no:cacheprovider tests/unit/scripts/test_sd7_compatibility_vectors.py -v
# PASS: manifest, W1 canonical vectors, W2 accounting vector (3 passed; 3 model-field warnings)
```

## Independent Go 1.25.10 verification

A separate verifier located the existing Go 1.25.10 installation, selected its `go` and
`gofmt` binaries, and confirmed `go version go1.25.10 linux/amd64`. It independently reran:

```text
go test . -run '^(TestW2.*|TestImportedPredecessorCorpus)$' -count=1 -v
go test ./... -count=1
go test -race ./...
```

The focused corpus and W2 cases, full five-package suite, race checks, and configured
check-only formatting passed. The verifier also confirmed unchanged tracked files and
frozen fixture hashes, reviewed the independent oracles and annotation validation, and
found no severe issue within the bounded harness scope. It did not independently execute
the Python leg or witness the historical RED/GREEN cycle; those results remain attributed
to the writer. The current Go toolchain blocker is resolved, not retroactively erased.

## Layout, size, and remaining limits

The three new root files are co-located external-package tests because the asserted boundary is
public API consumption. They contain 99, 160, and 156 lines respectively (415 test-helper lines
total), all at or below the 160-line per-Go-file limit. The root Go census moved from 69 tracked
files to 72 with these three untracked test files; the configured root threshold of 12 predates
this work. Production LOC changed by zero against the 400-production-LOC budget.

Remaining limits are deliberate: this only covers the immutable five-entry W2 fixture and its
specified synthetic threshold triangulation; it neither proves real Git acquisition nor general
CNSIC policy/configuration compatibility. Independent Go verification completed as recorded
above. The test-local known-production addition is safe for the fixed vectors exercised here;
it does not establish a general overflow-safe profile implementation.

## Key Learnings.

- W2's exclusive categories and known text totals can be observed through the neutral public API
  without carrying CNSIC policy behavior into production.
- A retained binary entry makes neutral line-threshold availability materially different from a
  legacy known-text threshold comparison.
- Frozen W1 vector success is provenance evidence only; it is not neutral canonical-wire parity.
