package gitadapter_test

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func snapshotSideScript(width int, raw, numstat []byte, oldFile, newFile bool) ([][]byte, [][]string) {
	outputs, commands := metadataAcquireScript(width, raw, numstat)
	for i, present := range []bool{oldFile, newFile} {
		if !present {
			continue
		}
		id, content := strings.Repeat("d", width), []byte{0, '\r', '\n', 'x'}
		if i == 1 {
			id, content = strings.Repeat("c", width), []byte{}
		}
		outputs = append(outputs, []byte("blob\n"), []byte(fmt.Sprintf("%d\n", len(content))), content, nil)
		commands = append(commands, []string{"cat-file", "-t", id}, []string{"cat-file", "-s", id}, []string{"cat-file", "blob", id}, []string{"cat-file", "-e", id + "^{blob}"})
	}
	outputs = append(outputs, nil, nil)
	commands = append(commands, []string{"cat-file", "-e", strings.Repeat("a", width) + "^{commit}"}, []string{"cat-file", "-e", strings.Repeat("b", width) + "^{commit}"})
	return outputs, commands
}

func TestSnapshotMetadataAcquireFactsCompleteSides(t *testing.T) {
	for _, width := range []int{40, 64} {
		for _, tc := range []struct {
			name, status, modes string
			oldFile, newFile    bool
		}{
			{"addition", "A", "", false, true}, {"deletion", "D", "", true, false},
			{"binary", "M", "", true, true}, {"rename", "R100", "", true, true},
			{"copy", "C100", "", true, true}, {"type", "T", "", true, false},
			{"symlink", "M", "120000 120000", false, false},
			{"gitlink", "M", "160000 160000", false, false},
			{"executable", "M", "100644 100755", true, true},
		} {
			t.Run(fmt.Sprintf("%d/%s", width, tc.name), func(t *testing.T) {
				path, previous := "new\t\n\xff", ""
				if tc.status[0] == 'R' || tc.status[0] == 'C' {
					previous = "old\t\n\xff"
				}
				raw := metadataRawFixture(tc.status, path, previous, width)
				if tc.modes != "" {
					raw = bytes.Replace(raw, []byte("100644 100644"), []byte(tc.modes), 1)
				}
				added, deleted := "0", "0"
				if tc.name == "binary" {
					added, deleted = "-", "-"
				}
				outputs, commands := snapshotSideScript(width, raw, metadataNumstatFixture(added, deleted, path, previous), tc.oldFile, tc.newFile)
				ctx, calls, remaining := context.Background(), 0, 0
				for _, out := range outputs {
					remaining += len(out)
				}
				bounds := metadataAcquireBounds()
				bounds.MaxTotalBytes = remaining
				runner := func(gotCtx context.Context, repo string, cap int, args ...string) ([]byte, error) {
					if calls >= len(commands) || gotCtx != ctx || repo != "/fixture" || cap != min(bounds.MaxResponseBytes, remaining) || !reflect.DeepEqual(args, commands[calls]) {
						t.Fatalf("shared owner/immutable command mismatch at%d: %v cap%d remaining%d", calls, args, cap, remaining)
					}
					out := outputs[calls]
					calls++
					remaining -= len(out)
					return out, nil
				}
				a, err := gitadapter.AcquireSnapshotMetadataBoundedWithRunner(ctx, runner, gitadapter.Request{Repository: "/fixture", Base: "base", Head: "head"}, bounds)
				if err != nil {
					t.Fatal(err)
				}
				before := a.Records()
				facts, err := a.AcquireFacts()
				if err != nil || len(facts.Entries()) != 1 || facts.Revisions() != a.Revisions() || calls != len(commands) || remaining != 0 || !reflect.DeepEqual(before, a.Records()) {
					t.Fatalf("complete facts: %v %v calls%d remaining%d", facts, err, calls, remaining)
				}
				entry := facts.Entries()[0]
				if entry.Status != tc.status[:1] || entry.Numstat.AddedRaw != added || entry.Numstat.DeletedRaw != deleted {
					t.Fatal("status or independent numstat evidence lost")
				}
				for i, side := range []*evidence.SnapshotSideV1{entry.Old, entry.New} {
					absent := i == 0 && tc.status == "A" || i == 1 && tc.status == "D"
					if absent {
						if side != nil {
							t.Fatal("absent side manufactured")
						}
						continue
					}
					wantPath := path
					if i == 0 && previous != "" {
						wantPath = previous
					}
					if side == nil || string(side.Path) != wantPath {
						t.Fatal("side path lost")
					}
					wantMode, wantKind, wantID := before[0].Change.OldMode, before[0].Change.OldKind, before[0].Change.OldObject
					if i == 1 {
						wantMode, wantKind, wantID = before[0].Change.NewMode, before[0].Change.NewKind, before[0].Change.NewObject
					}
					if side.Mode != wantMode || side.Kind != wantKind || side.ObjectID != wantID {
						t.Fatal("captured mode, kind or object identity changed")
					}
					if ![]bool{tc.oldFile, tc.newFile}[i] {
						if side.Content != nil {
							t.Fatal("non-file content measured")
						}
						continue
					}
					want := evidence.SnapshotContentV1{ByteLength: 4, ContainsNUL: true, LogicalLines: 2}
					if i == 1 {
						want = evidence.SnapshotContentV1{}
					}
					if side.Content == nil || *side.Content != want {
						t.Fatalf("byte metric: %+v", side.Content)
					}
				}
				doc, err := evidence.NewSnapshotV1(facts)
				if err != nil {
					t.Fatal(err)
				}
				for _, side := range []*evidence.SnapshotSideV1{entry.Old, entry.New} {
					if side == nil {
						continue
					}
					side.Path[0] = '!'
					if side.Content != nil {
						side.Content.ByteLength++
					}
				}
				again, err := evidence.NewSnapshotV1(facts)
				if err != nil || !bytes.Equal(doc.CanonicalBytes(), again.CanonicalBytes()) {
					t.Fatal("facts alias returned paths or content metrics")
				}
			})
		}
	}
}
