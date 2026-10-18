package gitadapter

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestBoundedSnapshotExactAndExceededBudgets(t *testing.T) {
	id := strings.Repeat("a", 40)
	raw := ":000000 100644 " + strings.Repeat("0", 40) + " " + id + " A\x00a\x00"
	outputs := []string{"sha1\n", id + "\n", "", id + "\n", "", raw, "1\t0\ta\x00", "", ""}
	total := 0
	for _, out := range outputs {
		total += len(out)
	}
	for _, tc := range []struct {
		name            string
		response, total int
		want            error
	}{
		{"exact", len(raw), total, nil},
		{"response below", len(raw) - 1, total, ErrSnapshotResourceLimit},
		{"aggregate below", len(raw), total - 1, ErrSnapshotResourceLimit},
		{"above", len(raw) + 1, total + 1, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls, remaining := 0, tc.total
			runner := func(_ context.Context, repo string, cap int, args ...string) ([]byte, error) {
				if cap != min(tc.response, remaining) {
					t.Fatalf("cap %d remaining %d", cap, remaining)
				}
				if repo != "/fixture" || calls >= len(outputs) {
					t.Fatal("unexpected command")
				}
				out := []byte(outputs[calls])
				calls++
				remaining -= len(out)
				return out, nil
			}
			got, err := AcquireBoundedWithRunner(context.Background(), runner, Request{"/fixture", "base", "head"}, SnapshotBounds{tc.response, tc.total, 1})
			if !errors.Is(err, tc.want) {
				t.Fatalf("error %v want %v", err, tc.want)
			}
			if err != nil && !reflect.DeepEqual(got, evidence.CommittedSnapshot{}) {
				t.Fatalf("partial %#v", got)
			}
			if err == nil && (len(got.Entries()) != 1 || calls != len(outputs)) {
				t.Fatalf("incomplete %#v calls %d", got, calls)
			}
		})
	}
}

func TestBoundedSnapshotImmutableArgumentsAndFailureMapping(t *testing.T) {
	id := strings.Repeat("b", 40)
	secret := errors.New("private raw subprocess stderr")
	for _, tc := range []struct {
		name     string
		failAt   int
		revision string
		code     evidence.SnapshotErrorCode
	}{
		{"foreign", 0, "base", evidence.SnapshotForeign},
		{"unresolved", 1, "$(literal);ref", evidence.SnapshotUnresolved},
		{"missing id", 1, id, evidence.SnapshotMissing},
		{"missing resolved", 2, "base", evidence.SnapshotMissing},
		{"diff race", 5, "base", evidence.SnapshotRacing},
		{"post read race", 7, "base", evidence.SnapshotRacing},
		{"ref movement harmless", -1, "$(literal);ref", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			runner := func(_ context.Context, _ string, _ int, args ...string) ([]byte, error) {
				n := calls
				calls++
				if n == tc.failAt {
					return nil, secret
				}
				switch n {
				case 0:
					return []byte("sha1\n"), nil
				case 1, 3:
					revision := tc.revision
					if n == 3 {
						revision = "head"
					}
					if !reflect.DeepEqual(args, []string{"rev-parse", "--verify", "--end-of-options", revision + "^{commit}"}) {
						t.Fatalf("revision not literal: %q", args)
					}
					return []byte(id + "\n"), nil
				case 5, 6:
					if args[len(args)-3] != id || args[len(args)-2] != id || args[len(args)-1] != "--" || !strings.Contains(strings.Join(args, " "), "--no-relative") {
						t.Fatalf("range %q", args)
					}
				}
				return nil, nil
			}
			got, err := AcquireBoundedWithRunner(context.Background(), runner, Request{"/fixture", tc.revision, "head"}, snapshotTestBounds())
			if tc.failAt >= 0 {
				assertFailure(t, got, err, tc.code)
				if strings.Contains(err.Error(), "private") {
					t.Fatal("raw stderr exposed")
				}
			} else if err != nil || got.Base != evidence.GitObjectID(id) {
				t.Fatalf("%#v %v", got, err)
			}
		})
	}
}

func TestBoundedSnapshotMidRunnerCancellation(t *testing.T) {
	for _, deadline := range []bool{false, true} {
		inner, cancel := context.WithCancel(context.Background())
		var ctx context.Context = inner
		want := context.Canceled
		if deadline {
			ctx = deadlineContext{inner}
			want = context.DeadlineExceeded
		}
		calls := 0
		got, err := AcquireBoundedWithRunner(ctx, func(context.Context, string, int, ...string) ([]byte, error) {
			calls++
			cancel()
			return []byte("sha1\n"), nil
		}, Request{"/fixture", "base", "head"}, snapshotTestBounds())
		cancel()
		if !errors.Is(err, want) || calls != 1 || !reflect.DeepEqual(got, evidence.CommittedSnapshot{}) {
			t.Fatalf("%#v %v calls=%d", got, err, calls)
		}
	}
}
