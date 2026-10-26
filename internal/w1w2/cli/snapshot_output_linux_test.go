//go:build linux

package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

func TestSnapshotCLIFailureMappingAndProfile(t *testing.T) {
	for _, tc := range []struct {
		err   error
		code  int
		token string
	}{
		{&evidence.SnapshotError{Code: evidence.SnapshotForeign}, 2, "invalid_arguments"},
		{&evidence.SnapshotError{Code: evidence.SnapshotMissing}, 3, "git_unavailable"},
		{&evidence.SnapshotError{Code: evidence.SnapshotRacing}, 7, "acquisition_race"},
		{gitadapter.ErrSnapshotResourceLimit, 4, "resource_limit"},
		{gitadapter.ErrSnapshotInvalidBounds, 2, "invalid_arguments"},
		{context.Canceled, 6, "interrupted"},
		{context.DeadlineExceeded, 6, "interrupted"},
		{errors.New("private path/content must not escape"), 3, "git_unavailable"},
		{nil, 2, "invalid_document"},
	} {
		t.Run(tc.token, func(t *testing.T) {
			calls := 0
			source := func(ctx context.Context, request gitadapter.Request, bounds gitadapter.SnapshotBounds) (evidence.SnapshotFactsV1, error) {
				calls++
				deadline, ok := ctx.Deadline()
				if !ok || time.Until(deadline) <= 0 || time.Until(deadline) > contract.OperationTimeout || bounds != (gitadapter.SnapshotBounds{MaxResponseBytes: contract.MaxGitResponseBytes, MaxTotalBytes: contract.MaxGitBytes, MaxEntries: contract.MaxDiffEntries}) || request != (gitadapter.Request{Repository: "/fixture", Base: "--literal-base", Head: "--literal-head"}) {
					t.Fatal("literal selectors, single operation deadline or shared profile changed")
				}
				return evidence.SnapshotFactsV1{}, tc.err
			}
			var out, diagnostic bytes.Buffer
			code := runWithSources(context.Background(), snapshotArguments(), unusedInput{}, &out, &diagnostic, nil, source)
			if code != tc.code || out.Len() != 0 || diagnostic.String() != "w1w2 snapshot: "+tc.token+"\n" || calls != 1 {
				t.Fatalf("code%d out%q stderr%q calls%d", code, out.Bytes(), diagnostic.String(), calls)
			}
		})
	}
}

func TestSnapshotCLIOutputLimitBeforeEmission(t *testing.T) {
	zero := uint64(0)
	facts, err := evidence.NewSnapshotFactsV1(snapshotFacts(t).Revisions(), []evidence.SnapshotEntryV1{{Status: "A", New: &evidence.SnapshotSideV1{Path: []byte(strings.Repeat("x", contract.MaxOutputBytes)), Mode: "100644", Kind: "file", ObjectID: evidence.GitObjectID(strings.Repeat("c", 40)), Content: &evidence.SnapshotContentV1{}}, Numstat: evidence.SnapshotNumstatV1{Kind: "text", AddedRaw: "0", DeletedRaw: "0", Added: &zero, Deleted: &zero}}})
	if err != nil {
		t.Fatal(err)
	}
	source := func(context.Context, gitadapter.Request, gitadapter.SnapshotBounds) (evidence.SnapshotFactsV1, error) {
		return facts, nil
	}
	var out, diagnostic bytes.Buffer
	code := runWithSources(context.Background(), snapshotArguments(), unusedInput{}, &out, &diagnostic, nil, source)
	if code != 4 || out.Len() != 0 || diagnostic.String() != "w1w2 snapshot: resource_limit\n" {
		t.Fatalf("%d %q %q", code, out.Bytes(), diagnostic.String())
	}
}
