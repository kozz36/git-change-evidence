package gitadapter_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func TestAcquireSnapshotMetadataBoundedBudgets(t *testing.T) {
	raw, numstat := metadataRawFixture("M", "a", "", 40), metadataNumstatFixture("2", "1", "a", "")
	outputs, _ := metadataAcquireScript(40, raw, numstat)
	total := 0
	for _, out := range outputs {
		total += len(out)
	}
	for _, tc := range []struct {
		name            string
		response, total int
		limited         bool
	}{
		{"exact", len(raw), total, false}, {"response_over", len(raw) - 1, total, true}, {"aggregate_over", len(raw), total - 1, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls, remaining := 0, tc.total
			runner := func(_ context.Context, _ string, cap int, _ ...string) ([]byte, error) {
				if calls >= len(outputs) || cap != min(tc.response, remaining) {
					t.Fatal("unexpected call or reset aggregate budget")
				}
				out := outputs[calls]
				calls++
				remaining -= len(out)
				return out, nil
			}
			got, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(context.Background(), runner, gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}, gitadapter.SnapshotBounds{MaxResponseBytes: tc.response, MaxTotalBytes: tc.total, MaxEntries: 1})
			if tc.limited {
				if !errors.Is(err, gitadapter.ErrSnapshotResourceLimit) || got != nil {
					t.Fatalf("limit: %v %v", got, err)
				}
			} else if err != nil || got == nil || calls != len(outputs) || len(got.Records()) != 1 {
				t.Fatalf("exact budget: %v %v", got, err)
			}
		})
	}
	outputs[5] = append(append([]byte(nil), raw...), metadataRawFixture("M", "b", "", 40)...)
	calls := 0
	runner := func(context.Context, string, int, ...string) ([]byte, error) {
		out := outputs[calls]
		calls++
		return out, nil
	}
	bounds := metadataAcquireBounds()
	bounds.MaxEntries = 1
	got, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(context.Background(), runner, gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}, bounds)
	if !errors.Is(err, gitadapter.ErrSnapshotResourceLimit) || got != nil || calls != 6 {
		t.Fatalf("entry limit must precede numstat: %v %v calls%d", got, err, calls)
	}
}

func TestAcquireSnapshotMetadataBoundedFailures(t *testing.T) {
	codes := []evidence.SnapshotErrorCode{evidence.SnapshotForeign, evidence.SnapshotUnresolved, evidence.SnapshotMissing, evidence.SnapshotUnresolved, evidence.SnapshotMissing, evidence.SnapshotRacing, evidence.SnapshotRacing, evidence.SnapshotRacing, evidence.SnapshotRacing}
	for failAt, code := range codes {
		t.Run(string(code)+string(rune('0'+failAt)), func(t *testing.T) {
			outputs, _ := metadataAcquireScript(40, metadataRawFixture("M", "a", "", 40), metadataNumstatFixture("2", "1", "a", ""))
			calls := 0
			runner := func(context.Context, string, int, ...string) ([]byte, error) {
				n := calls
				calls++
				if n == failAt {
					return nil, errors.New("private raw stderr")
				}
				return outputs[n], nil
			}
			got, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(context.Background(), runner, gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}, metadataAcquireBounds())
			var snapshotError *evidence.SnapshotError
			if !errors.As(err, &snapshotError) || snapshotError.Code != code || strings.Contains(err.Error(), "private") || got != nil || calls != failAt+1 {
				t.Fatalf("failure: %v %v calls%d", got, err, calls)
			}
		})
	}
}

func TestAcquireSnapshotMetadataBoundedInvalidAndCancelled(t *testing.T) {
	request := gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}
	for _, tc := range []struct {
		name        string
		request     gitadapter.Request
		bounds      gitadapter.SnapshotBounds
		code        evidence.SnapshotErrorCode
		boundsError bool
	}{
		{"invalid_bounds", request, gitadapter.SnapshotBounds{}, "", true},
		{"zero_response", request, gitadapter.SnapshotBounds{MaxResponseBytes: 0, MaxTotalBytes: 1, MaxEntries: 1}, "", true},
		{"zero_aggregate", request, gitadapter.SnapshotBounds{MaxResponseBytes: 1, MaxTotalBytes: 0, MaxEntries: 1}, "", true},
		{"zero_entries", request, gitadapter.SnapshotBounds{MaxResponseBytes: 1, MaxTotalBytes: 1, MaxEntries: 0}, "", true},
		{"relative_root", gitadapter.Request{Repository: "relative", Base: "base", Head: "head"}, metadataAcquireBounds(), evidence.SnapshotForeign, false},
		{"nul_root", gitadapter.Request{Repository: "/a\x00b", Base: "base", Head: "head"}, metadataAcquireBounds(), evidence.SnapshotForeign, false},
		{"empty_base", gitadapter.Request{Repository: "/fixture", Head: "head"}, metadataAcquireBounds(), evidence.SnapshotUnresolved, false},
		{"nul_head", gitadapter.Request{Repository: "/fixture", Base: "base", Head: "a\x00b"}, metadataAcquireBounds(), evidence.SnapshotUnresolved, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			runner := func(context.Context, string, int, ...string) ([]byte, error) { calls++; return nil, nil }
			got, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(context.Background(), runner, tc.request, tc.bounds)
			var snapshotError *evidence.SnapshotError
			valid := tc.boundsError && errors.Is(err, gitadapter.ErrSnapshotInvalidBounds) || !tc.boundsError && errors.As(err, &snapshotError) && snapshotError.Code == tc.code
			if !valid || got != nil || calls != 0 {
				t.Fatalf("invalid input executed commands: %v %v %d", got, err, calls)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	runner := func(context.Context, string, int, ...string) ([]byte, error) {
		calls++
		cancel()
		return []byte("sha1\n"), nil
	}
	got, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(ctx, runner, request, metadataAcquireBounds())
	if !errors.Is(err, context.Canceled) || got != nil || calls != 1 {
		t.Fatalf("mid-command cancellation: %v %v %d", got, err, calls)
	}
	got, err = gitadapter.AcquireSnapshotMetadataBoundedWithRunner(ctx, runner, request, metadataAcquireBounds())
	if !errors.Is(err, context.Canceled) || got != nil || calls != 1 {
		t.Fatalf("pre-command cancellation: %v %v %d", got, err, calls)
	}
	for _, tc := range []struct {
		ctx    context.Context
		runner gitadapter.Runner
	}{{nil, runner}, {context.Background(), nil}} {
		got, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(tc.ctx, tc.runner, request, metadataAcquireBounds())
		var snapshotError *evidence.SnapshotError
		if !errors.As(err, &snapshotError) || snapshotError.Code != evidence.SnapshotForeign || got != nil || calls != 1 {
			t.Fatalf("nil execution dependency: %v %v %d", got, err, calls)
		}
	}
	expired, closeDeadline := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer closeDeadline()
	got, err = gitadapter.AcquireSnapshotMetadataBoundedWithRunner(expired, runner, request, metadataAcquireBounds())
	if !errors.Is(err, context.DeadlineExceeded) || got != nil || calls != 1 {
		t.Fatalf("caller deadline: %v %v %d", got, err, calls)
	}
}
