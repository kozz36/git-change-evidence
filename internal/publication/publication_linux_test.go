//go:build linux

package publicationadapter

import (
	"context"
	"errors"
	evidence "github.com/kozz36/git-change-evidence"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
	"testing"
)

func TestPublishWritesCanonicalEvidenceAtDigestIdentity(t *testing.T) {
	root := t.TempDir()
	document := testEvidence(t)
	got, err := Publish(context.Background(), root, document)
	must(t, err)
	want := "sha256/" + string(document.Digest()) + ".json"
	require(t, got.Identity == want, "identity")
	bytes, err := os.ReadFile(filepath.Join(root, got.Identity))
	must(t, err)
	require(t, string(bytes) == string(document.CanonicalBytes()), "bytes")
	info, err := os.Stat(filepath.Join(root, "sha256"))
	must(t, err)
	require(t, info.Mode().Perm() == 0o700, "sha256 mode")
	info, err = os.Stat(filepath.Join(root, got.Identity))
	must(t, err)
	require(t, info.Mode().Perm() == 0o600, "artifact mode")
}
func testEvidence(t *testing.T) evidence.Evidence {
	document, err := evidence.NewReportV1(evidence.ReportInput{Subject: "publication-test", Provenance: evidence.Provenance{AccountingDigest: evidence.Digest("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), InputPolicyDigest: evidence.Digest("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), InventoryDigest: evidence.Digest("cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"), Revisions: evidence.RevisionIdentity{Base: "1111111111111111111111111111111111111111", Head: "2222222222222222222222222222222222222222"}}})
	must(t, err)
	return document
}
func TestPublishRejectsExistingEntriesAndSymlinkRoots(t *testing.T) {
	document := testEvidence(t)
	for _, kind := range []string{"regular", "symlink", "directory"} {
		root := t.TempDir()
		must(t, os.Mkdir(filepath.Join(root, "sha256"), 0o700))
		final := filepath.Join(root, "sha256", string(document.Digest())+".json")
		switch kind {
		case "regular":
			must(t, os.WriteFile(final, []byte("existing"), 0o600))
		case "symlink":
			must(t, os.Symlink("target", final))
		case "directory":
			must(t, os.Mkdir(final, 0o700))
		}
		before, err := os.Lstat(final)
		must(t, err)
		_, err = Publish(context.Background(), root, document)
		assertCode(t, err, Conflict)
		after, err := os.Lstat(final)
		must(t, err)
		require(t, os.SameFile(before, after), "existing entry")
		if kind == "regular" {
			bytes, err := os.ReadFile(final)
			must(t, err)
			require(t, string(bytes) == "existing", "regular entry")
		}
	}
	base := t.TempDir()
	actual := filepath.Join(base, "actual")
	must(t, os.MkdirAll(filepath.Join(actual, "child"), 0o700))
	must(t, os.Symlink(actual, filepath.Join(base, "root")))
	must(t, os.Symlink(actual, filepath.Join(base, "link")))
	for _, root := range []string{filepath.Join(base, "root"), filepath.Join(base, "link", "child")} {
		_, err := Publish(context.Background(), root, document)
		assertCode(t, err, Unavailable)
	}
}
func TestPublishCancellationAndExclusivity(t *testing.T) {
	document, root := testEvidence(t), t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Publish(ctx, root, document)
	assertCode(t, err, Interrupted)
	_, err = os.Stat(filepath.Join(root, "sha256"))
	require(t, os.IsNotExist(err), "canceled effects")
	start, results := make(chan struct{}), make(chan error, 2)
	for range 2 {
		go func() {
			<-start
			_, err := Publish(context.Background(), root, document)
			results <- err
		}()
	}
	close(start)
	first, second := <-results, <-results
	require(t, (first == nil) != (second == nil), "concurrency")
	if first != nil {
		assertCode(t, first, Conflict)
	} else {
		assertCode(t, second, Conflict)
	}
	bytes, err := os.ReadFile(filepath.Join(root, "sha256", string(document.Digest())+".json"))
	must(t, err)
	require(t, string(bytes) == string(document.CanonicalBytes()), "concurrent bytes")
}
func TestSyscallSeams(t *testing.T) {
	for _, fail := range []string{"tmpfile", "write", "sync"} {
		assertCode(t, stageFailure(fail), Unavailable)
	}
	assertCode(t, finalize(context.Background(), 1, "x", func(int, string) error { return unix.EOPNOTSUPP }, nil, nil), Unavailable)
	fd, err := unix.Open(t.TempDir(), unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	must(t, err)
	defer unix.Close(fd)
	err = finalize(context.Background(), fd, "rollback", linkEntry, func(int) error { return unix.EIO }, func(parent int, name string) error { return unix.Unlinkat(parent, name, 0) })
	assertCode(t, err, Unavailable)
	var stat unix.Stat_t
	require(t, errors.Is(unix.Fstatat(fd, "rollback", &stat, unix.AT_SYMLINK_NOFOLLOW), unix.ENOENT), "rollback entry")
	must(t, unix.Mkdirat(fd, "sha256", 0o700))
	_, err = openDirectory(fd, "sha256", func() {
		must(t, unix.Renameat(fd, "sha256", fd, "moved"))
		must(t, unix.Mkdirat(fd, "sha256", 0o700))
	})
	assertCode(t, err, Unavailable)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	linked := false
	assertCode(t, finalize(ctx, 1, "x", func(int, string) error { linked = true; return nil }, nil, nil), Interrupted)
	require(t, !linked, "canceled finalization")
}
func TestWriteAll(t *testing.T) {
	type result struct {
		n   int
		err error
	}
	for _, test := range []struct {
		name    string
		sizes   []int
		results []result
		fail    bool
	}{
		{"short", []int{4, 3}, []result{{1, nil}, {3, nil}}, false},
		{"EINTR", []int{4, 4}, []result{{0, unix.EINTR}, {4, nil}}, false},
		{"zero", []int{4}, []result{{0, nil}}, true},
		{"negative", []int{4}, []result{{-1, nil}}, true},
		{"oversized", []int{4}, []result{{5, nil}}, true},
		{"error", []int{4}, []result{{4, unix.EIO}}, true},
	} {
		calls := 0
		err := writeAll(1, []byte("four"), func(_ int, bytes []byte) (int, error) {
			if len(bytes) != test.sizes[calls] {
				return -1, unix.EIO
			}
			result := test.results[calls]
			calls++
			return result.n, result.err
		})
		if (err != nil) != test.fail || calls != len(test.results) {
			t.Fatalf("%s: err, calls = %v, %d", test.name, err, calls)
		}
	}
}
func stageFailure(fail string) error {
	_, err := stage(1, []byte("x"), func(int) (int, error) { return 7, map[string]error{"tmpfile": unix.EOPNOTSUPP}[fail] }, func(int) error { return nil }, func(_ int, bytes []byte) (int, error) { return len(bytes), map[string]error{"write": unix.EIO}[fail] }, func(int) error { return map[string]error{"sync": unix.EIO}[fail] }, func(int) error { return nil })
	return err
}
func linkEntry(parent int, name string) error {
	fd, err := unix.Openat(parent, name, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_CLOEXEC, 0o600)
	if err != nil {
		return err
	}
	return unix.Close(fd)
}
func must(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
func require(t *testing.T, condition bool, message string) {
	if !condition {
		t.Fatal(message)
	}
}
func assertCode(t *testing.T, err error, want ErrorCode) {
	var got *Error
	if !errors.As(err, &got) || got.Code != want || err.Error() != "evidence publication: "+string(want) {
		t.Fatalf("error = %v, want %s", err, want)
	}
}
