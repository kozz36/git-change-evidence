package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

func outputFixture(t *testing.T, nameBytes, pathBytes int) (evidence.PolicyDocumentV1, evidence.InventoryDocumentV1, evidence.CommittedSnapshot, []byte) {
	t.Helper()
	name := strings.Repeat("c", nameBytes)
	p, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{Categories: []evidence.CategoryInput{{Name: name, LineThreshold: 1}}, Default: name})
	if err != nil {
		t.Fatal(err)
	}
	i, err := evidence.NewInventoryV1(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	snapshot := evidence.NewCommittedSnapshot(evidence.GitObjectID(strings.Repeat("a", 40)), evidence.GitObjectID(strings.Repeat("b", 40)), []evidence.CommittedChange{{
		Status: "A", Path: strings.Repeat("p", pathBytes), OldMode: "000000", NewMode: "100644",
		OldObject: evidence.GitObjectID(strings.Repeat("0", 40)), NewObject: evidence.GitObjectID(strings.Repeat("c", 40)),
		OldKind: evidence.GitAbsent, NewKind: evidence.GitFile, Lines: evidence.CommittedLineCounts{Additions: 2, Countable: true},
	}})
	result, err := evidence.NewAccountingResultV1(p, i, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	return p, i, snapshot, result.CanonicalBytes()
}

func TestCanonicalOutputBudgetBelowExactlyAbove(t *testing.T) {
	pathBytes, nameBytes := 0, 0
	for n := 1; n <= 20; n++ {
		_, _, _, a := outputFixture(t, 1, n)
		_, _, _, b := outputFixture(t, 2, n)
		step := len(b) - len(a)
		if step > 0 && (contract.MaxOutputBytes-len(a))%step == 0 {
			pathBytes, nameBytes = n, 1+(contract.MaxOutputBytes-len(a))/step
			break
		}
	}
	if pathBytes == 0 {
		t.Fatal("no admissible exact canonical output fixture")
	}
	for _, delta := range []int{-1, 0, 1} {
		p, i, snapshot, want := outputFixture(t, nameBytes+delta, pathBytes)
		if len(p.CanonicalBytes()) > contract.MaxDocumentBytes || len(p.CanonicalBytes())+len(i.CanonicalBytes()) > contract.MaxInputBytes {
			t.Fatal("output fixture exceeds input admission bounds")
		}
		args := []string{"result", "--repo", t.TempDir(), "--base", "base", "--head", "head", "--policy", documentFile(t, p.CanonicalBytes()), "--inventory", documentFile(t, i.CanonicalBytes())}
		source := func(context.Context, gitadapter.Request, gitadapter.SnapshotBounds) (evidence.CommittedSnapshot, error) {
			return snapshot, nil
		}
		var out, diagnostic bytes.Buffer
		code := runWithSource(context.Background(), args, unusedInput{}, &out, &diagnostic, source)
		if delta > 0 {
			if len(want) <= contract.MaxOutputBytes || code != 4 || out.Len() != 0 || diagnostic.String() != "w1w2 result: resource_limit\n" {
				t.Fatalf("over limit: length=%d code=%d stdout=%d stderr=%q", len(want), code, out.Len(), diagnostic.String())
			}
		} else if code != 0 || diagnostic.Len() != 0 || !bytes.Equal(out.Bytes(), want) || delta == 0 && out.Len() != contract.MaxOutputBytes {
			t.Fatalf("admitted output: delta=%d length=%d code=%d stderr=%q", delta, len(want), code, diagnostic.String())
		}
	}
}
