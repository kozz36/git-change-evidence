package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func TestRunUsesCommittedBlobsAndFormatsExactly(t *testing.T) {
	repository := t.TempDir()
	git(t, repository, "init", "-q")
	git(t, repository, "config", "user.email", "test@example.invalid")
	git(t, repository, "config", "user.name", "Test")
	write(t, repository, "source", "keep\nmove\n")
	git(t, repository, "add", "source")
	git(t, repository, "commit", "-qm", "base")
	write(t, repository, "source", "keep\n")
	write(t, repository, "new", "move\nfresh\n")
	git(t, repository, "add", "source", "new")
	git(t, repository, "commit", "-qm", "head")
	write(t, repository, "source", "indexed\n")
	write(t, repository, "new", "indexed\n")
	git(t, repository, "add", "source", "new")
	write(t, repository, "source", "dirty\n")
	write(t, repository, "new", "dirty\n")
	var stdout, stderr bytes.Buffer
	if code := run(context.Background(), repository, []string{"HEAD~1", "source", "new"}, &stdout, &stderr, gitRunner, standardLimits); code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if got, want := stdout.String(), "moved=1 new=1 additions=2\n"; got != want || stderr.Len() != 0 {
		t.Fatalf("stdout, stderr = %q, %q; want %q, empty", got, stderr.String(), want)
	}
}
func TestRunRejectsResourceCaps(t *testing.T) {
	assert := func(blobs map[string]string, bound limits) {
		fixture := newFixture(blobs)
		var stdout, stderr bytes.Buffer
		code := run(context.Background(), "repository", []string{"base", "source", "new"}, &stdout, &stderr, fixture.runner, bound)
		if code != exitBound || stdout.Len() != 0 || stderr.String() != "resource bound exceeded\n" {
			t.Fatalf("exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
		}
	}
	assert(map[string]string{"base/source": "aa", "head/source": "b", "head/new": "cc"}, limits{perBlob: 3, aggregate: 4, lines: 8, lineBytes: 8, pairs: 64})
	blobs := map[string]string{"base/source": "keep\nmove", "head/source": "keep", "head/new": "move\nfresh"}
	assert(blobs, limits{16, 32, 3, 16, 64})
	assert(blobs, limits{16, 32, 8, 16, 1})
}

type testClock []time.Time

func (c *testClock) Now() time.Time {
	now := (*c)[0]
	*c = (*c)[len(*c)-1:]
	return now
}
func TestMatchedSharesAbsoluteDeadline(t *testing.T) {
	bound, now := limits{16, 32, 8, 16, 64}, time.Unix(0, 0)
	limit := now.Add(time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), limit)
	deadline := matcherDeadline(ctx, now)
	clock, pairs := &testClock{limit.Add(-time.Millisecond)}, bound.pairs
	source, _ := prepared("move", "move", &bound.lines, bound)
	if result, code := matched(context.Background(), source, &pairs, deadline, clock.Now); code != 0 || result == nil {
		t.Fatalf("source = %v, %d", result, code)
	}
	*clock = []time.Time{limit.Add(-time.Millisecond), limit}
	if result, code := matched(context.Background(), source, &pairs, deadline, clock.Now); code != exitStopped || result != nil {
		t.Fatalf("destination = %v, %d", result, code)
	}
}
func TestRunClassifiesInputAndContentFailures(t *testing.T) {
	valid := newFixture(map[string]string{"base/source": "x\nx", "head/source": "", "head/new": "x\nx\nx"})
	invalid := newFixture(map[string]string{"base/source": string([]byte{0xff}), "head/source": "", "head/new": "x"})
	missing := newFixture(map[string]string{"base/source": "x", "head/source": "", "head/new": "x"})
	nonBlob := newFixture(map[string]string{"base/tree": "x"})
	nonBlob.nonBlob = "tree"
	assert := func(args []string, runner gitadapter.Runner, want int) {
		var stdout, stderr bytes.Buffer
		got := run(context.Background(), "repository", args, &stdout, &stderr, runner, standardLimits)
		if got != want || (got == 0 && stdout.String() != "moved=2 new=1 additions=3\n") || (got != 0 && stdout.Len() != 0) || (got != 0 && strings.Count(stderr.String(), "\n") != 1) {
			t.Fatalf("exit, stdout, stderr = %d, %q, %q", got, stdout.String(), stderr.String())
		}
	}
	assert([]string{"base"}, nil, exitInvalid)
	assert([]string{"base", "source", "new"}, invalid.runner, exitUTF8)
	assert([]string{"base", "missing", "new"}, missing.runner, exitAbsent)
	assert([]string{"base", "tree", "new"}, nonBlob.runner, exitAbsent)
	assert([]string{"base", "source", "new"}, valid.runner, 0)
	if valid.headResolves != 1 {
		t.Fatalf("HEAD resolutions = %d, want 1", valid.headResolves)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if got := run(canceled, "repository", []string{"base", "source", "new"}, &bytes.Buffer{}, &bytes.Buffer{}, valid.runner, standardLimits); got != exitStopped {
		t.Fatalf("canceled exit = %d", got)
	}
}

type fixture struct {
	blobs, objects      map[string]string
	base, head, nonBlob string
	headResolves        int
}

func newFixture(blobs map[string]string) *fixture {
	f := &fixture{blobs: blobs, objects: map[string]string{}, base: strings.Repeat("a", 40), head: strings.Repeat("b", 40)}
	for _, content := range blobs {
		f.objects[fmt.Sprintf("%040x", len(f.objects)+1)] = content
	}
	return f
}
func (f *fixture) runner(_ context.Context, _ string, _ int, args ...string) ([]byte, error) {
	switch args[0] {
	case "rev-parse":
		if args[1] == "--show-object-format" {
			return []byte("sha1\n"), nil
		}
		revision := strings.TrimSuffix(args[len(args)-1], "^{commit}")
		if revision == "HEAD" {
			f.headResolves++
			return []byte(f.head + "\n"), nil
		}
		if revision == "base" || revision == f.base {
			return []byte(f.base + "\n"), nil
		}
		if revision == f.head {
			return []byte(f.head + "\n"), nil
		}
	case "ls-tree":
		commit, path := args[3], strings.TrimPrefix(args[5], ":(top,literal)")
		key := map[string]string{f.base: "base/", f.head: "head/"}[commit] + path
		if path == f.nonBlob {
			return []byte("040000 tree " + strings.Repeat("c", 40) + "\t" + path + "\x00"), nil
		}
		if content, ok := f.blobs[key]; ok {
			for object, value := range f.objects {
				if value == content {
					return []byte("100644 blob " + object + "\t" + path + "\x00"), nil
				}
			}
		}
		return []byte{0}, nil
	case "cat-file":
		if args[1] == "-e" {
			return nil, nil
		}
		for object, content := range f.objects {
			if args[len(args)-1] == object {
				if args[1] == "-s" {
					return []byte(fmt.Sprintf("%d\n", len(content))), nil
				}
				return []byte(content), nil
			}
		}
	}
	return nil, errors.New("missing fixture object")
}
func git(t *testing.T, repository string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repository}, args...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", args, err, output)
	}
}
func write(t *testing.T, repository, name, contents string) {
	t.Helper()
	if err := os.WriteFile(repository+"/"+name, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
