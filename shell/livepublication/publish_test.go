package livepublication

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	publication "github.com/kozz36/git-change-evidence/internal/publication"
)

func git(t *testing.T, repo string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_AUTHOR_NAME=fixture", "GIT_AUTHOR_EMAIL=fixture@example.invalid", "GIT_COMMITTER_NAME=fixture", "GIT_COMMITTER_EMAIL=fixture@example.invalid")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fixture git: %v: %s", err, out)
	}
	return strings.TrimSpace(string(out))
}

func fixture(t *testing.T) (string, string, string, string) {
	t.Helper()
	if testing.Short() {
		t.Skip("real Git/filesystem integration")
	}
	repo := t.TempDir()
	git(t, repo, "init", "--quiet")
	ids := make([]string, 3)
	for i := range ids {
		if err := os.WriteFile(filepath.Join(repo, "file"), []byte(strings.Repeat("line\n", i+1)), 0600); err != nil {
			t.Fatal(err)
		}
		git(t, repo, "add", "file")
		git(t, repo, "commit", "--quiet", "-m", "fixture")
		ids[i] = git(t, repo, "rev-parse", "HEAD")
	}
	git(t, repo, "update-ref", "refs/heads/live", ids[1])
	return repo, ids[0], ids[1], ids[2]
}

func report(t *testing.T, base, head evidence.GitObjectID) evidence.Evidence {
	t.Helper()
	digest := evidence.Digest(strings.Repeat("a", 64))
	doc, err := evidence.NewReportV1(evidence.ReportInput{Subject: "fixture", Provenance: evidence.Provenance{
		AccountingDigest: digest, InputPolicyDigest: digest, InventoryDigest: digest,
		Revisions: evidence.RevisionIdentity{Base: string(base), Head: string(head)},
	}})
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestPublishStableAndRebuild(t *testing.T) {
	for _, move := range []bool{false, true} {
		name := "stable"
		if move {
			name = "rebuild"
		}
		t.Run(name, func(t *testing.T) {
			repo, base, a, b := fixture(t)
			destination, builds := t.TempDir(), 0
			var accepted evidence.Evidence
			result, err := Publish(context.Background(), repo, base, "refs/heads/live", destination, func(gotBase, gotHead evidence.GitObjectID) (evidence.Evidence, error) {
				builds++
				want := a
				if builds == 2 {
					want = b
				}
				if string(gotBase) != base || string(gotHead) != want {
					t.Fatal("builder IDs")
				}
				accepted = report(t, gotBase, gotHead)
				if move && builds == 1 {
					git(t, repo, "update-ref", "refs/heads/live", b)
				}
				return accepted, nil
			})
			if err != nil {
				t.Fatal(err)
			}
			wantBuilds := 1
			if move {
				wantBuilds = 2
			}
			if builds != wantBuilds {
				t.Fatalf("builds = %d, want %d", builds, wantBuilds)
			}
			if result.Identity != "sha256/"+string(accepted.Digest())+".json" {
				t.Fatal("digest identity")
			}
			raw, err := os.ReadFile(filepath.Join(destination, result.Identity))
			if err != nil || !bytes.Equal(raw, accepted.CanonicalBytes()) {
				t.Fatalf("canonical artifact: %v", err)
			}
			decoded, err := evidence.DecodeCanonical(raw)
			if err != nil || decoded.Provenance() != accepted.Provenance() {
				t.Fatal("artifact provenance")
			}
			entries, err := os.ReadDir(filepath.Join(destination, "sha256"))
			if err != nil || len(entries) != 1 {
				t.Fatal("only accepted artifact")
			}
		})
	}
}

func TestRejectBeforePublisher(t *testing.T) {
	for _, mode := range []string{"moving", "divergent", "zero", "builder-error"} {
		t.Run(mode, func(t *testing.T) {
			repo, base, a, b := fixture(t)
			destination, builds, calls := t.TempDir(), 0, 0
			builderErr := &evidence.ContractError{Field: "subject", Code: "required"}
			_, err := publish(context.Background(), repo, base, "refs/heads/live", destination, func(gotBase, gotHead evidence.GitObjectID) (evidence.Evidence, error) {
				builds++
				if mode == "zero" {
					return evidence.Evidence{}, nil
				}
				if mode == "builder-error" {
					return evidence.Evidence{}, builderErr
				}
				if mode == "divergent" {
					return report(t, gotBase, evidence.GitObjectID(base)), nil
				}
				next := b
				if string(gotHead) == b {
					next = a
				}
				git(t, repo, "update-ref", "refs/heads/live", next)
				return report(t, gotBase, gotHead), nil
			}, func(ctx context.Context, root string, doc evidence.Evidence) (Result, error) {
				calls++
				return publication.Publish(ctx, root, doc)
			})
			if mode == "moving" {
				var typed *evidence.SnapshotError
				if !errors.As(err, &typed) || typed.Code != evidence.SnapshotRacing || builds != 3 {
					t.Fatalf("racing/builds: %v/%d", err, builds)
				}
			} else if mode == "builder-error" {
				if !errors.Is(err, builderErr) {
					t.Fatalf("builder error: %v", err)
				}
			} else {
				var typed *Error
				if !errors.As(err, &typed) || typed.Code != Invalid {
					t.Fatalf("invalid: %v", err)
				}
			}
			entries, readErr := os.ReadDir(destination)
			if calls != 0 || readErr != nil || len(entries) != 0 {
				t.Fatal("publisher called or artifact created")
			}
		})
	}
}

func TestMovementAfterCheckPublishesObservedPair(t *testing.T) {
	repo, base, a, b := fixture(t)
	destination := t.TempDir()
	var accepted evidence.Evidence
	_, err := publish(context.Background(), repo, base, "refs/heads/live", destination, func(base, head evidence.GitObjectID) (evidence.Evidence, error) {
		accepted = report(t, base, head)
		return accepted, nil
	}, func(ctx context.Context, root string, doc evidence.Evidence) (Result, error) {
		git(t, repo, "update-ref", "refs/heads/live", b)
		return publication.Publish(ctx, root, doc)
	})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(destination, "sha256", string(accepted.Digest())+".json"))
	if err != nil || !bytes.Equal(raw, accepted.CanonicalBytes()) || accepted.Provenance().Revisions.Head != a {
		t.Fatal("observed snapshot lost")
	}
	if git(t, repo, "rev-parse", "refs/heads/live") != b {
		t.Fatal("fixture did not move")
	}
}
