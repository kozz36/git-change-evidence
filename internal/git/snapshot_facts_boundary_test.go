package gitadapter_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func TestSnapshotMetadataAcquireFactsRejectsUninitializedOwner(t *testing.T) {
	for _, a := range []*gitadapter.SnapshotMetadataAcquisition{nil, new(gitadapter.SnapshotMetadataAcquisition)} {
		facts, err := a.AcquireFacts()
		var failed *evidence.SnapshotError
		if !errors.As(err, &failed) || failed.Code != evidence.SnapshotForeign || facts.Revisions() != (evidence.RevisionIdentity{}) || len(facts.Entries()) != 0 {
			t.Fatalf("uninitialized owner: %v %v", facts, err)
		}
	}
}

func TestSnapshotMetadataAcquireFactsRetainsOwnerAcrossCalls(t *testing.T) {
	for _, scenario := range []string{"empty", "already_canceled", "repeat_exhaustion"} {
		t.Run(scenario, func(t *testing.T) {
			raw, numstat := metadataRawFixture("M", "a", "", 40), metadataNumstatFixture("0", "0", "a", "")
			file := scenario != "empty"
			if !file {
				raw, numstat = nil, nil
			}
			outputs, commands := snapshotSideScript(40, raw, numstat, file, file)
			bounds, calls := metadataAcquireBounds(), 0
			if scenario == "repeat_exhaustion" {
				bounds.MaxTotalBytes = 0
				for _, out := range outputs {
					bounds.MaxTotalBytes += len(out)
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			runner := func(gotCtx context.Context, repository string, cap int, args ...string) ([]byte, error) {
				at := calls
				calls++
				if gotCtx != ctx || repository != "/fixture" {
					t.Fatal("owner context or repository changed")
				}
				if at >= len(outputs) {
					if scenario != "repeat_exhaustion" || cap != 0 || !reflect.DeepEqual(args, commands[9]) {
						t.Fatal("budget reset or a different object command was used")
					}
					return []byte("blob\n"), nil
				}
				return outputs[at], nil
			}
			a, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(ctx, runner, gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}, bounds)
			if err != nil {
				t.Fatal(err)
			}
			before := a.Records()
			if scenario == "already_canceled" {
				cancel()
			}
			facts, err := a.AcquireFacts()
			if scenario == "already_canceled" {
				if !errors.Is(err, context.Canceled) || calls != 9 || facts.Revisions() != (evidence.RevisionIdentity{}) || len(facts.Entries()) != 0 {
					t.Fatalf("pre-side cancellation: %v %v calls%d", facts, err, calls)
				}
				return
			}
			if err != nil || facts.Revisions() != a.Revisions() || calls != len(outputs) || !reflect.DeepEqual(before, a.Records()) {
				t.Fatalf("completed owner: %v %v calls%d", facts, err, calls)
			}
			if scenario == "empty" {
				if len(facts.Entries()) != 0 {
					t.Fatal("empty range invented entries")
				}
				if _, err := evidence.NewSnapshotV1(facts); err != nil {
					t.Fatal(err)
				}
				return
			}
			second, err := a.AcquireFacts()
			if !errors.Is(err, gitadapter.ErrSnapshotResourceLimit) || calls != len(outputs)+1 || second.Revisions() != (evidence.RevisionIdentity{}) || len(second.Entries()) != 0 || !reflect.DeepEqual(before, a.Records()) || len(facts.Entries()) != 1 {
				t.Fatalf("repeat acquisition reset budget or rewrote prior facts: %v %v calls%d", second, err, calls)
			}
		})
	}
}
