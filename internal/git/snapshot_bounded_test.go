package gitadapter

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func snapshotTestBounds() SnapshotBounds { return SnapshotBounds{16 << 20, 64 << 20, 10000} }

func TestBoundedSnapshotWholeRange(t *testing.T) {
	if testing.Short() {
		t.Skip("uses fixture Git")
	}
	repo := repository(t, false)
	write(t, repo, "a.txt", "before\n")
	write(t, repo, "b.txt", "before\n")
	base := commit(t, repo)
	write(t, repo, "a.txt", "after\n")
	write(t, repo, "b.txt", "after\n")
	head := commit(t, repo)
	request := Request{repo, base, head}
	want, err := Acquire(request)
	if err != nil {
		t.Fatal(err)
	}
	write(t, repo, "a.txt", "worktree only\n")
	git(t, repo, "config", "diff.relative", "true")
	got, err := AcquireBounded(context.Background(), request, snapshotTestBounds())
	if err != nil || !reflect.DeepEqual(got, want) || len(got.Entries()) != 2 {
		t.Fatalf("got %#v, %v; want %#v", got, err, want)
	}
	request.Head = base
	got, err = AcquireBounded(context.Background(), request, snapshotTestBounds())
	if err != nil || len(got.Entries()) != 0 || got.Base != got.Head {
		t.Fatalf("zero range: %#v, %v", got, err)
	}
}

func TestBoundedSnapshotRejectsLimitsAndCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, tc := range []struct {
		name   string
		ctx    context.Context
		bounds SnapshotBounds
		want   error
	}{
		{"canceled", ctx, snapshotTestBounds(), context.Canceled},
		{"invalid", context.Background(), SnapshotBounds{}, ErrSnapshotInvalidBounds},
		{"zero response", context.Background(), SnapshotBounds{0, 1, 1}, ErrSnapshotInvalidBounds},
		{"zero total", context.Background(), SnapshotBounds{1, 0, 1}, ErrSnapshotInvalidBounds},
		{"zero entries", context.Background(), SnapshotBounds{1, 1, 0}, ErrSnapshotInvalidBounds},
		{"negative response", context.Background(), SnapshotBounds{-1, 1, 1}, ErrSnapshotInvalidBounds},
		{"negative total", context.Background(), SnapshotBounds{1, -1, 1}, ErrSnapshotInvalidBounds},
		{"negative entries", context.Background(), SnapshotBounds{1, 1, -1}, ErrSnapshotInvalidBounds},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			got, err := AcquireBoundedWithRunner(tc.ctx, func(context.Context, string, int, ...string) ([]byte, error) { called = true; return nil, nil }, Request{}, tc.bounds)
			if !errors.Is(err, tc.want) || called || !reflect.DeepEqual(got, evidence.CommittedSnapshot{}) {
				t.Fatalf("%#v %v called=%v", got, err, called)
			}
		})
	}
	for _, tc := range []struct {
		name      string
		response  []byte
		runnerErr error
	}{
		{"oversized mock", []byte(strings.Repeat("x", 9)), nil},
		{"typed overflow", nil, ErrSnapshotResourceLimit},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := AcquireBoundedWithRunner(context.Background(), func(context.Context, string, int, ...string) ([]byte, error) { return tc.response, tc.runnerErr }, Request{Repository: "/fixture", Base: "base", Head: "head"}, SnapshotBounds{8, 32, 1})
			if !errors.Is(err, ErrSnapshotResourceLimit) || !reflect.DeepEqual(got, evidence.CommittedSnapshot{}) {
				t.Fatalf("%#v %v", got, err)
			}
		})
	}
}
