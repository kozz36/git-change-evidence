package gitadapter

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
	"time"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestBoundedSnapshotRunningGitCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("uses a real Git process and owned FIFO")
	}
	repo := repository(t, false)
	write(t, repo, "file", "content\n")
	head := commit(t, repo)
	fifo := filepath.Join(t.TempDir(), "include")
	if err := syscall.Mkfifo(fifo, 0600); err != nil {
		t.Fatal(err)
	}
	git(t, repo, "config", "include.path", fifo)
	before := boundedFixtureState(t, repo)
	ctx, cancel := context.WithCancel(context.Background())
	type outcome struct {
		snapshot evidence.CommittedSnapshot
		err      error
	}
	result := make(chan outcome, 1)
	go func() {
		snapshot, err := AcquireBounded(ctx, Request{repo, head, head}, snapshotTestBounds())
		result <- outcome{snapshot, err}
	}()
	type opened struct {
		file *os.File
		err  error
	}
	ready := make(chan opened, 1)
	go func() {
		file, err := os.OpenFile(fifo, os.O_WRONLY, 0600)
		ready <- opened{file, err}
	}()
	var writer *os.File
	joined, received := false, false
	t.Cleanup(func() {
		cancel()
		if !joined {
			// Release a failed-readiness blocking open without polling processes.
			rescue, err := os.OpenFile(fifo, os.O_RDWR|syscall.O_NONBLOCK, 0600)
			if err != nil {
				t.Errorf("unblock FIFO writer: %v", err)
			} else {
				defer rescue.Close()
			}
			select {
			case opened := <-ready:
				writer = opened.file
			case <-time.After(time.Second):
				t.Error("FIFO opener did not join")
			}
		}
		if writer != nil {
			writer.Close()
		}
		if !received {
			select {
			case <-result:
			case <-time.After(time.Second):
				t.Error("owned acquisition did not join")
			}
		}
	})
	select {
	case opened := <-ready:
		joined, writer = true, opened.file
		if opened.err != nil {
			t.Fatal(opened.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Git did not open the owned FIFO")
	}
	// No other reader exists: the blocking writer open proves Git reached its
	// config include. Keeping this writer open without bytes blocks its read.
	select {
	case got := <-result:
		received = true
		t.Fatalf("Git returned before cancellation: %v", got.err)
	default:
	}
	cancel()
	select {
	case got := <-result:
		received = true
		if !errors.Is(got.err, context.Canceled) || !reflect.DeepEqual(got.snapshot, evidence.CommittedSnapshot{}) {
			t.Fatalf("cancellation: %#v, %v", got.snapshot, got.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("running Git did not return after cancellation")
	}
	// snapshotCommand returns only after exec.Cmd.Run's Wait boundary. This
	// observes direct-child collection, not PID inspection or grandchild safety.
	if after := boundedFixtureState(t, repo); !reflect.DeepEqual(before, after) {
		t.Fatal("canceled acquisition changed the repository")
	}
	t.Log("FIFO rendezvous observed; canceled native acquisition returned after Run/Wait with empty snapshot and unchanged repository hashes")
}
