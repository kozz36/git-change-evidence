package input_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
	"github.com/kozz36/git-change-evidence/internal/w1w2/input"
)

func failure(t *testing.T, err error, status int, token string) {
	t.Helper()
	var f *contract.Failure
	if !errors.As(err, &f) || f.Status != status || f.Token != token {
		t.Fatalf("got %v; want %d %s", err, status, token)
	}
}

func TestPublicInput(t *testing.T) {
	r, err := input.Parse([]string{"inventory", "--paths", "-", "--policy", "policy", "--source-root", "/root"})
	if err != nil || r.Stage != "inventory" || r.Paths != "-" {
		t.Fatalf("%+v %v", r, err)
	}
	paths, err := input.DecodeSelection([]byte(`{"schema":"git-change-evidence.cli-path-selection/v1","paths_b64":["YS50eHQ="]}`))
	if err != nil || len(paths) != 1 || string(paths[0]) != "a.txt" {
		t.Fatalf("%q %v", paths, err)
	}
	docs, err := input.ReadDocuments(context.Background(), contract.Request{Policy: "-"}, strings.NewReader("raw\n"))
	if err != nil || string(docs.Policy) != "raw\n" {
		t.Fatalf("%+v %v", docs, err)
	}
}
