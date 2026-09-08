package gitadapter

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestAcquireComposesCommittedResultAndReportV1(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git repository")
	}
	if runtime.GOOS == "windows" {
		t.Skip("the raw non-UTF-8 newline path fixture requires Unix filesystem semantics")
	}
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_ATTR_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_COUNT", "0")
	t.Setenv("GIT_NO_REPLACE_OBJECTS", "1")
	repo := repository(t, false)
	hostile := string([]byte("raw/\xff\npath"))
	for _, directory := range []string{"docs", "raw", "bulk"} {
		if err := os.MkdirAll(filepath.Join(repo, directory), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(t, repo, "docs/old.txt", "rename\n")
	write(t, repo, "binary", "before\x00")
	write(t, repo, hostile, "before\n")
	write(t, repo, "bulk/256.txt", strings.Repeat("before\n", 256))
	base := commit(t, repo)
	if err := os.MkdirAll(filepath.Join(repo, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "mv", "docs/old.txt", "src/new.go")
	write(t, repo, "binary", "after\x00")
	write(t, repo, hostile, "after\n")
	write(t, repo, "bulk/256.txt", strings.Repeat("after\n", 256))
	head := commit(t, repo)

	snapshot, err := Acquire(Request{Repository: repo, Base: base, Head: head})
	if err != nil {
		t.Fatal(err)
	}
	entries := map[string]evidence.CommittedChange{}
	for _, entry := range snapshot.Entries() {
		entries[entry.Path] = entry
	}
	if len(entries) != 4 || snapshot.Base != evidence.GitObjectID(base) || snapshot.Head != evidence.GitObjectID(head) {
		t.Fatalf("snapshot revisions/entries = %#v", snapshot)
	}
	if got := entries["src/new.go"]; got.Status != "R" || got.PreviousPath != "docs/old.txt" || got.Path != "src/new.go" || got.Lines != (evidence.CommittedLineCounts{Countable: true}) {
		t.Fatalf("rename = %#v", got)
	}
	if got := entries["binary"]; got.Status != "M" || !got.Binary || got.Lines.Countable {
		t.Fatalf("binary = %#v", got)
	}
	if got := entries[hostile]; got.Status != "M" || !bytes.Equal([]byte(got.Path), []byte("raw/\xff\npath")) || got.Lines != (evidence.CommittedLineCounts{Additions: 1, Deletions: 1, Countable: true}) {
		t.Fatalf("hostile path = %#v", got)
	}
	if got := entries["bulk/256.txt"]; got.Status != "M" || got.Lines != (evidence.CommittedLineCounts{Additions: 256, Deletions: 256, Countable: true}) {
		t.Fatalf("bounded text = %#v", got)
	}

	policy, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{Categories: []evidence.CategoryInput{
		{Name: "docs", PathGlobs: [][]byte{[]byte("docs/**")}}, {Name: "source", PathGlobs: [][]byte{[]byte("src/**")}},
		{Name: "binary", PathGlobs: [][]byte{[]byte("binary")}}, {Name: "raw", PathGlobs: [][]byte{[]byte("raw/**")}},
		{Name: "bulk", PathGlobs: [][]byte{[]byte("bulk/**")}}, {Name: "other"},
	}, Default: "other"})
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := evidence.NewInventoryV1(policy, nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := evidence.NewAccountingResultV1(policy, inventory, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	wantEntries := []evidence.ResultEntryV1{
		{Path: []byte("binary"), Category: "binary", Measurement: evidence.ResultEntryMeasurementV1{Kind: evidence.ResultEntryNonCountableV1}},
		{Path: []byte("bulk/256.txt"), Category: "bulk", Measurement: evidence.ResultEntryMeasurementV1{Kind: evidence.ResultEntryCountableV1, Additions: 256, Deletions: 256}},
		{Path: []byte("raw/\xff\npath"), Category: "raw", Measurement: evidence.ResultEntryMeasurementV1{Kind: evidence.ResultEntryCountableV1, Additions: 1, Deletions: 1}},
		{Path: []byte("src/new.go"), Category: "source", Measurement: evidence.ResultEntryMeasurementV1{Kind: evidence.ResultEntryCountableV1}},
	}
	wantTotals := []evidence.ResultCategoryTotalV1{{Category: "docs"}, {Category: "source"}, {Category: "binary", NonCountable: 1}, {Category: "raw", Additions: 1, Deletions: 1}, {Category: "bulk", Additions: 256, Deletions: 256}, {Category: "other"}}
	if !reflect.DeepEqual(result.Entries(), wantEntries) || !reflect.DeepEqual(result.Totals(), wantTotals) || result.Revisions() != (evidence.RevisionIdentity{Base: base, Head: head}) || result.AccountingPolicyDigest() != policy.Digest() || result.InventoryDigest() != inventory.Digest() {
		t.Fatalf("public composition = (%#v, %#v, %#v)", result.Entries(), result.Totals(), result.Revisions())
	}
	binding := evidence.ResultV1ReportBinding{Policy: policy, Inventory: inventory, Result: result}
	report, err := evidence.NewReportV1FromResult("acquired committed change", binding)
	if err != nil || evidence.ValidateReportV1ResultProvenance(report, binding) != nil || report.Provenance().Revisions != result.Revisions() || report.Provenance().AccountingDigest != evidence.Digest(result.Digest()) || report.Provenance().InputPolicyDigest != evidence.Digest(policy.Digest()) || report.Provenance().InventoryDigest != evidence.Digest(inventory.Digest()) {
		t.Fatalf("report provenance = (%#v, %v)", report.Provenance(), err)
	}

	write(t, repo, "binary", "indexed only\x00")
	git(t, repo, "add", "binary")
	write(t, repo, hostile, "worktree only\n")
	write(t, repo, "untracked", "ignored\n")
	again, err := Acquire(Request{Repository: repo, Base: base, Head: head})
	if err != nil || !reflect.DeepEqual(snapshot, again) {
		t.Fatalf("dirty acquisition = (%#v, %v)", again, err)
	}
	repeated, err := evidence.NewAccountingResultV1(policy, inventory, again)
	repeatedReport, reportErr := evidence.NewReportV1FromResult("acquired committed change", evidence.ResultV1ReportBinding{Policy: policy, Inventory: inventory, Result: repeated})
	if err != nil || reportErr != nil || !bytes.Equal(result.CanonicalBytes(), repeated.CanonicalBytes()) || result.Digest() != repeated.Digest() || !bytes.Equal(report.CanonicalBytes(), repeatedReport.CanonicalBytes()) || report.Digest() != repeatedReport.Digest() {
		t.Fatalf("dirty reproducibility = (%#v, %v, %v)", repeated, err, reportErr)
	}
	empty, err := Acquire(Request{Repository: repo, Base: head, Head: head})
	if err != nil || empty.Base != evidence.GitObjectID(head) || empty.Head != evidence.GitObjectID(head) || len(empty.Entries()) != 0 {
		t.Fatalf("empty committed diff = (%#v, %v)", empty, err)
	}
}
