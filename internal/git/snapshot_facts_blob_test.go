package gitadapter_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func TestSnapshotMetadataAcquireFactsRejectsIncompleteObjects(t *testing.T) {
	for _, tc := range []struct {
		name   string
		at     int
		output []byte
	}{
		{"wrong_type", 9, []byte("tree\n")}, {"malformed_type", 9, []byte(" blob\n")},
		{"negative_size", 10, []byte("-1\n")}, {"padded_size", 10, []byte("04\n")},
		{"unterminated_size", 10, []byte("4")}, {"overflow_size", 10, []byte("18446744073709551616\n")},
		{"short_blob", 11, []byte("abc")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			outputs, _ := snapshotSideScript(40, metadataRawFixture("M", "a", "", 40), metadataNumstatFixture("0", "0", "a", ""), true, true)
			outputs[tc.at] = tc.output
			calls := 0
			runner := func(context.Context, string, int, ...string) ([]byte, error) {
				out := outputs[calls]
				calls++
				return out, nil
			}
			a, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(context.Background(), runner, gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}, metadataAcquireBounds())
			if err != nil {
				t.Fatal(err)
			}
			before := a.Records()
			facts, err := a.AcquireFacts()
			var failed *evidence.SnapshotError
			if !errors.As(err, &failed) || failed.Code != evidence.SnapshotRacing || err.Error() != "snapshot acquisition: racing_object" || facts.Revisions() != (evidence.RevisionIdentity{}) || len(facts.Entries()) != 0 || calls != tc.at+1 || !reflect.DeepEqual(before, a.Records()) {
				t.Fatalf("incomplete object accepted or state changed: %v %v calls%d", facts, err, calls)
			}
		})
	}
	for failAt := 9; failAt < 19; failAt++ {
		t.Run(fmt.Sprintf("command_error_%d", failAt), func(t *testing.T) {
			outputs, _ := snapshotSideScript(40, metadataRawFixture("M", "a", "", 40), metadataNumstatFixture("0", "0", "a", ""), true, true)
			calls := 0
			runner := func(context.Context, string, int, ...string) ([]byte, error) {
				at := calls
				calls++
				if at == failAt {
					return nil, errors.New("private runner stderr")
				}
				return outputs[at], nil
			}
			a, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(context.Background(), runner, gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}, metadataAcquireBounds())
			if err != nil {
				t.Fatal(err)
			}
			facts, err := a.AcquireFacts()
			var failed *evidence.SnapshotError
			if !errors.As(err, &failed) || failed.Code != evidence.SnapshotRacing || err.Error() != "snapshot acquisition: racing_object" || len(facts.Entries()) != 0 || facts.Revisions() != (evidence.RevisionIdentity{}) || calls != failAt+1 {
				t.Fatalf("side error leaked or partial facts returned: %v %v calls%d", facts, err, calls)
			}
		})
	}
}

func TestSnapshotMetadataAcquireFactsKeepsSharedBoundsAndContext(t *testing.T) {
	for _, scenario := range []string{"exact_response", "over_response", "aggregate_exhaustion", "canceled"} {
		t.Run(scenario, func(t *testing.T) {
			outputs, _ := snapshotSideScript(40, metadataRawFixture("M", "a", "", 40), metadataNumstatFixture("0", "0", "a", ""), true, true)
			bounds, calls := metadataAcquireBounds(), 0
			if scenario == "exact_response" || scenario == "over_response" {
				bounds.MaxResponseBytes = 256
				outputs[10], outputs[11] = []byte("256\n"), bytes.Repeat([]byte("x"), 256)
				if scenario == "over_response" {
					outputs[10] = []byte("257\n")
				}
			}
			if scenario == "aggregate_exhaustion" {
				bounds.MaxTotalBytes = 0
				for _, out := range outputs[:11] {
					bounds.MaxTotalBytes += len(out)
				}
				bounds.MaxTotalBytes += 3
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			runner := func(context.Context, string, int, ...string) ([]byte, error) {
				at := calls
				calls++
				if scenario == "canceled" && at == 9 {
					cancel()
				}
				return outputs[at], nil
			}
			a, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(ctx, runner, gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}, bounds)
			if err != nil {
				t.Fatal(err)
			}
			before := a.Records()
			facts, err := a.AcquireFacts()
			if !reflect.DeepEqual(before, a.Records()) {
				t.Fatal("measurement mutated metadata session")
			}
			if scenario == "exact_response" {
				if err != nil || len(facts.Entries()) != 1 || *facts.Entries()[0].Old.Content != (evidence.SnapshotContentV1{ByteLength: 256, LogicalLines: 1}) || calls != len(outputs) {
					t.Fatalf("exact bound: %v %v calls%d", facts, err, calls)
				}
				return
			}
			wantErr, wantCalls := gitadapter.ErrSnapshotResourceLimit, 11
			if scenario == "canceled" {
				wantErr, wantCalls = context.Canceled, 10
			}
			if !errors.Is(err, wantErr) || len(facts.Entries()) != 0 || facts.Revisions() != (evidence.RevisionIdentity{}) || calls != wantCalls {
				t.Fatalf("shared admission/context: %v %v calls%d", facts, err, calls)
			}
		})
	}
}
