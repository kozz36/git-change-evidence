package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

type incompleteWriter struct {
	bytes.Buffer
	fail bool
}

func (w *incompleteWriter) Write(raw []byte) (int, error) {
	n, _ := w.Buffer.Write(raw[:len(raw)/2])
	if w.fail {
		return n, io.ErrClosedPipe
	}
	return n, nil
}

func TestOutputShortOrFailedWritesAreNeverSuccess(t *testing.T) {
	for _, failed := range []bool{false, true} {
		var out incompleteWriter
		out.fail = failed
		var diagnostic bytes.Buffer
		code := Run(context.Background(), reportArguments(t), unusedInput{}, &out, &diagnostic)
		if code != 9 || out.Len() == 0 || diagnostic.String() != "w1w2 report: output_incomplete\n" {
			t.Fatalf("%d %q %q", code, out.Bytes(), diagnostic.String())
		}
	}
}

func TestCancellationBeforeComposition(t *testing.T) {
	args := reportArguments(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out, diagnostic bytes.Buffer
	if code := Run(ctx, args, unusedInput{}, &out, &diagnostic); code != 6 || out.Len() != 0 || diagnostic.String() != "w1w2 report: interrupted\n" {
		t.Fatalf("%d %q %q", code, out.Bytes(), diagnostic.String())
	}
}

func TestGitFailuresMapWithoutRawDetails(t *testing.T) {
	p, inventory, _ := documents(t)
	args := []string{"result", "--repo", t.TempDir(), "--base", "$(literal);base", "--head", "head", "--policy", documentFile(t, p.CanonicalBytes()), "--inventory", documentFile(t, inventory.CanonicalBytes())}
	for _, tc := range []struct {
		err   error
		code  int
		token string
	}{
		{&evidence.SnapshotError{Code: evidence.SnapshotForeign}, 2, "invalid_arguments"},
		{&evidence.SnapshotError{Code: evidence.SnapshotMissing}, 3, "git_unavailable"},
		{&evidence.SnapshotError{Code: evidence.SnapshotUnresolved}, 3, "git_unavailable"},
		{&evidence.SnapshotError{Code: evidence.SnapshotRacing}, 7, "acquisition_race"},
		{gitadapter.ErrSnapshotResourceLimit, 4, "resource_limit"},
		{gitadapter.ErrSnapshotInvalidBounds, 2, "invalid_arguments"},
		{context.Canceled, 6, "interrupted"},
		{context.DeadlineExceeded, 6, "interrupted"},
		{errors.New("private path/content must not escape"), 3, "git_unavailable"},
	} {
		var out, diagnostic bytes.Buffer
		source := func(ctx context.Context, request gitadapter.Request, bounds gitadapter.SnapshotBounds) (evidence.CommittedSnapshot, error) {
			if request.Base != "$(literal);base" || request.Head != "head" {
				t.Fatal("revision arguments changed")
			}
			if _, ok := ctx.Deadline(); !ok || bounds.MaxResponseBytes != contract.MaxGitResponseBytes || bounds.MaxTotalBytes != contract.MaxGitBytes || bounds.MaxEntries != contract.MaxDiffEntries {
				t.Fatal("context/profile not propagated")
			}
			return evidence.CommittedSnapshot{}, tc.err
		}
		code := runWithSource(context.Background(), args, unusedInput{}, &out, &diagnostic, source)
		if code != tc.code || out.Len() != 0 || diagnostic.String() != "w1w2 result: "+tc.token+"\n" {
			t.Fatalf("%d %q %q", code, out.Bytes(), diagnostic.String())
		}
	}
}

func TestCancellationAfterResultConstruction(t *testing.T) {
	p, inventory, _ := documents(t)
	args := []string{"result", "--repo", t.TempDir(), "--base", "base", "--head", "head", "--policy", documentFile(t, p.CanonicalBytes()), "--inventory", documentFile(t, inventory.CanonicalBytes())}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	source := func(context.Context, gitadapter.Request, gitadapter.SnapshotBounds) (evidence.CommittedSnapshot, error) {
		cancel()
		return emptySnapshot(), nil
	}
	var stdout, stderr bytes.Buffer
	code := runWithSource(ctx, args, unusedInput{}, &stdout, &stderr, source)
	if code != 6 || stdout.Len() != 0 || stderr.String() != "w1w2 result: interrupted\n" {
		t.Fatalf("%d %q %q", code, stdout.Bytes(), stderr.String())
	}
}

func TestFailedDiagnosticDeliveryStillReturnsFailure(t *testing.T) {
	var stdout bytes.Buffer
	stderr := incompleteWriter{fail: true}
	if code := Run(context.Background(), []string{"unknown"}, unusedInput{}, &stdout, &stderr); code != 2 || stdout.Len() != 0 {
		t.Fatalf("code=%d stdout=%q", code, stdout.Bytes())
	}
}
