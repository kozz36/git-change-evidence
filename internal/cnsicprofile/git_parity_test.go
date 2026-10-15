package cnsicprofile_test

import (
	"bytes"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	"github.com/kozz36/git-change-evidence/internal/cnsicprofile"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

// Frozen W2 semantics: testdata/compatibility/v1/accounting-v1.json,
// predecessor a0fc7b26ff8a0e0a61baa586b32c46841611c806, SHA256
// f265bc48750ccd68e41dcb64ec68134b9fb7db69655e2c215e39f8c6c8794c3f.
// TestImportedPredecessorCorpus checks integrity in the full suite. This is
// change-only parity with empty inventory, not fresh predecessor execution.
func TestCNSICProfileRealGitSemanticParity(t *testing.T) {
	if testing.Short() {
		t.Skip("uses real local Git repositories")
	}
	for _, format := range []string{"sha1", "sha256"} {
		t.Run(format, func(t *testing.T) {
			repo := t.TempDir()
			git := func(args ...string) string {
				t.Helper()
				cmd := exec.Command("git", append([]string{"--no-replace-objects", "-C", repo, "-c", "core.hooksPath=/dev/null", "-c", "user.name=fixture", "-c", "user.email=fixture@example.invalid"}, args...)...)
				cmd.Env = []string{"PATH=/usr/bin:/bin", "HOME=" + os.Getenv("HOME"), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_ATTR_NOSYSTEM=1", "GIT_CONFIG_COUNT=0", "GIT_TERMINAL_PROMPT=0", "LC_ALL=C"}
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("git %v: %v: %s", args, err, out)
				}
				return strings.TrimSpace(string(out))
			}
			write := func(path, content string) {
				t.Helper()
				target := filepath.Join(repo, path)
				if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(target, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
			}
			git("init", "--object-format="+format)
			git("commit", "--allow-empty", "-m", "baseline")
			base := git("rev-parse", "HEAD")
			paths := []string{"binary.bin", "generated.json", "openspec/changes/vector.md", "src/core.py", "tests/vector.py"}
			categories := []string{"production", "mechanical", "artifacts", "production", "tests"}
			additions := []uint64{0, 2, 5, 4, 3}
			for i, path := range paths {
				content := strings.Repeat("line\n", int(additions[i]))
				if i == 0 {
					content = "binary\x00data"
				}
				write(path, content)
			}
			git("add", "--", ".")
			git("commit", "-m", "frozen semantic vector")
			head := git("rev-parse", "HEAD")
			length := 40
			if format == "sha256" {
				length = 64
			}
			validID := func(id string) bool { _, err := hex.DecodeString(id); return len(id) == length && err == nil }
			if !validID(base) || !validID(head) || base == head {
				t.Fatal("invalid committed revisions")
			}
			input := control()
			// Frozen generated annotation is represented solely by this policy selector.
			input.Categories[2].PathGlobs = [][]byte{[]byte("generated.json")}
			neutral, err := evidence.NewAccountingPolicyV1(input)
			if err != nil {
				t.Fatal(err)
			}
			policy, err := cnsicprofile.DecodePolicyV1(neutral.CanonicalBytes())
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(policy.CanonicalBytes(), neutral.CanonicalBytes()) || policy.Digest() != neutral.Digest() {
				t.Fatal("policy identity")
			}
			inventory, err := evidence.NewInventoryV1(policy, nil)
			if err != nil {
				t.Fatal(err)
			}
			acquire := func() (evidence.CommittedSnapshot, evidence.ResultDocumentV1) {
				t.Helper()
				snapshot, err := gitadapter.Acquire(gitadapter.Request{Repository: repo, Base: base, Head: head})
				if err != nil {
					t.Fatal(err)
				}
				result, err := evidence.NewAccountingResultV1(policy, inventory, snapshot)
				if err != nil {
					t.Fatal(err)
				}
				return snapshot, result
			}
			snapshot, result := acquire()
			entries := snapshot.Entries()
			if len(entries) != 5 || snapshot.Base != evidence.GitObjectID(base) || snapshot.Head != evidence.GitObjectID(head) {
				t.Fatal("snapshot identity/count")
			}
			wantEntries := make([]evidence.ResultEntryV1, len(paths))
			for i, entry := range entries {
				if entry.Path != paths[i] || entry.Status != "A" || !validID(string(entry.NewObject)) || string(entry.NewObject) != git("rev-parse", head+":"+paths[i]) {
					t.Fatalf("committed path/object: %#v", entry)
				}
				if entry.Binary != (i == 0) || entry.Lines != (evidence.CommittedLineCounts{Additions: additions[i], Countable: i != 0}) {
					t.Fatalf("measurement: %#v", entry)
				}
				kind := evidence.ResultEntryCountableV1
				if i == 0 {
					kind = evidence.ResultEntryNonCountableV1
				}
				wantEntries[i] = evidence.ResultEntryV1{Path: []byte(paths[i]), Category: categories[i], Measurement: evidence.ResultEntryMeasurementV1{Kind: kind, Additions: additions[i]}}
			}
			wantTotals := []evidence.ResultCategoryTotalV1{{Category: "artifacts", Additions: 5}, {Category: "tests", Additions: 3}, {Category: "mechanical", Additions: 2}, {Category: "production", Additions: 4, NonCountable: 1}}
			// Available=false means UNKNOWN; Actual/Exceeded zero values are not verdicts.
			wantObservations := []evidence.ResultObservationV1{{Kind: evidence.ResultLineThresholdObservationV1, Category: "production", LineThresholdLimit: 400}}
			if !reflect.DeepEqual(result.Entries(), wantEntries) || !reflect.DeepEqual(result.Totals(), wantTotals) || !reflect.DeepEqual(result.Observations(), wantObservations) {
				t.Fatalf("semantic parity: %v / %v / %v", result.Entries(), result.Totals(), result.Observations())
			}
			if result.Revisions() != (evidence.RevisionIdentity{Base: base, Head: head}) || result.AccountingPolicyDigest() != policy.Digest() || result.InventoryDigest() != inventory.Digest() {
				t.Fatal("result identity/binding")
			}
			for _, path := range paths {
				write(path, strings.Repeat("dirty\n", 17))
			}
			git("add", "--", ".")
			write("src/core.py", "unstaged\n")
			if git("status", "--porcelain") == "" {
				t.Fatal("dirty-domain control is empty")
			}
			again, repeated := acquire()
			if !reflect.DeepEqual(snapshot, again) || repeated.Revisions() != result.Revisions() || !bytes.Equal(result.CanonicalBytes(), repeated.CanonicalBytes()) || repeated.Digest() != result.Digest() {
				t.Fatal("mutable index/worktree affected immutable profile result")
			}
		})
	}
}
