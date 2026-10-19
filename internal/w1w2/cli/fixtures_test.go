package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func documents(t *testing.T) (evidence.PolicyDocumentV1, evidence.InventoryDocumentV1, evidence.ResultDocumentV1) {
	t.Helper()
	policy, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{
		Categories: []evidence.CategoryInput{{Name: "all", PathGlobs: [][]byte{[]byte("**")}, LineThreshold: 1}},
		Default:    "all",
	})
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := evidence.NewInventoryV1(policy, nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := evidence.NewAccountingResultV1(policy, inventory, emptySnapshot())
	if err != nil {
		t.Fatal(err)
	}
	return policy, inventory, result
}

func emptySnapshot() evidence.CommittedSnapshot {
	return evidence.NewCommittedSnapshot(evidence.GitObjectID(strings.Repeat("a", 40)), evidence.GitObjectID(strings.Repeat("b", 40)), nil)
}

func documentFile(t *testing.T, raw []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "document")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func reportArguments(t *testing.T) []string {
	t.Helper()
	p, i, result := documents(t)
	return []string{"report", "--policy", documentFile(t, p.CanonicalBytes()), "--inventory", documentFile(t, i.CanonicalBytes()), "--result", documentFile(t, result.CanonicalBytes())}
}

type unusedInput struct{}

func (unusedInput) Read([]byte) (int, error) { panic("unused stdin was read") }

func invoke(t *testing.T, args []string) (int, []byte, string) {
	t.Helper()
	var out, diagnostic bytes.Buffer
	code := Run(context.Background(), args, unusedInput{}, &out, &diagnostic)
	return code, out.Bytes(), diagnostic.String()
}

func requireFailure(t *testing.T, args []string, code int, diagnostic string) {
	t.Helper()
	got, out, err := invoke(t, args)
	if got != code || len(out) != 0 || err != diagnostic {
		t.Fatalf("code=%d stdout=%q stderr=%q", got, out, err)
	}
}
