package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
	publicationadapter "github.com/kozz36/git-change-evidence/internal/publication"
)

func TestRunCLIDispatchesExactProjectArity(t *testing.T) {
	input := projectInput(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI(context.Background(), "repository", []string{"project", "canonical-json"}, bytes.NewReader(input), &stdout, &stderr, nil, limits{}, nil); code != 0 || !bytes.Equal(stdout.Bytes(), input) || stderr.Len() != 0 {
		t.Fatalf("two-argument project exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
	}
	for _, args := range [][]string{{"project"}, {"project", "canonical-json", "extra", "extra"}} {
		stdout.Reset()
		stderr.Reset()
		if code := runCLI(context.Background(), "repository", args, bytes.NewReader(input), &stdout, &stderr, nil, limits{}, nil); code != exitInvalid || stdout.Len() != 0 || stderr.String() != "invalid input\n" {
			t.Fatalf("args %q: exit, stdout, stderr = %d, %q, %q", args, code, stdout.String(), stderr.String())
		}
	}
}

func TestRunCLIPreservesThreeArgumentSpecialBaseReferences(t *testing.T) {
	for _, reference := range []string{"project", "publish"} {
		t.Run(reference, func(t *testing.T) {
			fixture := newFixture(map[string]string{"base/source": "keep\nmove", "head/source": "keep", "head/new": "move\nfresh"})
			fixture.baseAlias = reference
			calls := 0
			publisher := func(context.Context, string, evidence.Evidence) (publicationadapter.Result, error) {
				calls++
				return publicationadapter.Result{}, nil
			}
			var stdout, stderr bytes.Buffer

			if code := runCLI(context.Background(), "repository", []string{reference, "source", "new"}, bytes.NewReader(nil), &stdout, &stderr, fixture.runner, standardLimits(), publisher); code != 0 {
				t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
			}
			if got, want := stdout.String(), "moved=1 new=1 additions=2\n"; got != want || stderr.Len() != 0 || calls != 0 {
				t.Fatalf("stdout, stderr, calls = %q, %q, %d; want %q, empty, zero", got, stderr.String(), calls, want)
			}
		})
	}
}

func TestRunPublishSuccess(t *testing.T) {
	root := "/publication-root/../unchanged"
	input := projectInput(t)
	want := publicationadapter.Result{Identity: "sha256/identity.json"}
	calls := 0
	publisher := func(_ context.Context, gotRoot string, document evidence.Evidence) (publicationadapter.Result, error) {
		calls++
		if gotRoot != root {
			t.Errorf("root = %q, want %q", gotRoot, root)
		}
		if got := document.CanonicalBytes(); !bytes.Equal(got, input) {
			t.Errorf("canonical bytes = %q, want %q", got, input)
		}
		return want, nil
	}
	var stdout, stderr bytes.Buffer
	if code := runPublish(context.Background(), []string{root}, bytes.NewReader(input), &stdout, &stderr, publisher); code != 0 || stdout.String() != want.Identity+"\n" || stderr.Len() != 0 || calls != 1 {
		t.Fatalf("exit, stdout, stderr, calls = %d, %q, %q, %d", code, stdout.String(), stderr.String(), calls)
	}
}

func TestRunPublishMapsFailuresAndInput(t *testing.T) {
	input := projectInput(t)
	for _, test := range []struct {
		name       string
		input      io.Reader
		failure    error
		wantCode   int
		wantStderr string
	}{
		{"invalid", bytes.NewReader(input), &publicationadapter.Error{Code: publicationadapter.Invalid}, exitInvalid, "invalid input\n"},
		{"unavailable", bytes.NewReader(input), &publicationadapter.Error{Code: publicationadapter.Unavailable}, exitAbsent, "content unavailable\n"},
		{"interrupted", bytes.NewReader(input), &publicationadapter.Error{Code: publicationadapter.Interrupted}, exitStopped, "operation interrupted\n"},
		{"conflict", bytes.NewReader(input), &publicationadapter.Error{Code: publicationadapter.Conflict}, exitConflict, "publication conflict\n"},
		{"unknown", bytes.NewReader(input), errors.New("unknown publisher failure"), exitAbsent, "content unavailable\n"},
		{"noncanonical", bytes.NewReader(append(input, ' ')), nil, exitInvalid, "invalid input\n"},
		{"read failure", failingReader{errors.New("read failure")}, nil, exitAbsent, "content unavailable\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			publisher := func(context.Context, string, evidence.Evidence) (publicationadapter.Result, error) {
				calls++
				return publicationadapter.Result{}, test.failure
			}
			var stdout, stderr bytes.Buffer
			if code := runPublish(context.Background(), []string{"/root"}, test.input, &stdout, &stderr, publisher); code != test.wantCode || stdout.Len() != 0 || stderr.String() != test.wantStderr {
				t.Fatalf("exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
			}
			wantCalls := 1
			if test.failure == nil {
				wantCalls = 0
			}
			if calls != wantCalls {
				t.Fatalf("publisher calls = %d, want %d", calls, wantCalls)
			}
		})
	}
	oversized := &countingReader{reader: bytes.NewReader(bytes.Repeat([]byte("x"), int(projectInputCap)+2))}
	var stdout, stderr bytes.Buffer
	if code := runPublish(context.Background(), []string{"/root"}, oversized, &stdout, &stderr, func(context.Context, string, evidence.Evidence) (publicationadapter.Result, error) {
		t.Fatal("publisher called")
		return publicationadapter.Result{}, nil
	}); code != exitBound || stdout.Len() != 0 || stderr.String() != "resource bound exceeded\n" {
		t.Fatalf("oversized exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
	}
	if got, want := oversized.read, int(projectInputCap+1); got != want {
		t.Fatalf("oversized bytes read = %d, want %d", got, want)
	}
}

func TestRunCommandDispatchesPublishBeforeCWD(t *testing.T) {
	input := projectInput(t)
	cwdCalls, publisherCalls := 0, 0
	publisher := func(context.Context, string, evidence.Evidence) (publicationadapter.Result, error) {
		publisherCalls++
		return publicationadapter.Result{Identity: "sha256/identity.json"}, nil
	}
	unavailableCWD := func() (string, error) {
		cwdCalls++
		return "", errors.New("cwd unavailable")
	}
	var stdout, stderr bytes.Buffer
	if code := runCommand(context.Background(), []string{"publish", "/root"}, bytes.NewReader(input), &stdout, &stderr, unavailableCWD, nil, limits{}, publisher); code != 0 || stdout.String() != "sha256/identity.json\n" || stderr.Len() != 0 || cwdCalls != 0 || publisherCalls != 1 {
		t.Fatalf("exit, stdout, stderr, cwd, publisher = %d, %q, %q, %d, %d", code, stdout.String(), stderr.String(), cwdCalls, publisherCalls)
	}
	for _, args := range [][]string{{"publish"}, {"publish", "/root", "extra", "extra"}} {
		stdout.Reset()
		stderr.Reset()
		if code := runCommand(context.Background(), args, bytes.NewReader(input), &stdout, &stderr, unavailableCWD, nil, limits{}, publisher); code != exitInvalid || stdout.Len() != 0 || stderr.String() != "invalid input\n" || cwdCalls != 0 || publisherCalls != 1 {
			t.Fatalf("args %q: exit, stdout, stderr, cwd, publisher = %d, %q, %q, %d, %d", args, code, stdout.String(), stderr.String(), cwdCalls, publisherCalls)
		}
	}
}

func TestRunProjectCanonicalJSON(t *testing.T) {
	input := projectInput(t)
	var stdout, stderr bytes.Buffer

	if code := runProject([]string{"canonical-json"}, bytes.NewReader(input), &stdout, &stderr); code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if got := stdout.Bytes(); !bytes.Equal(got, input) || stderr.Len() != 0 {
		t.Fatalf("stdout, stderr = %q, %q; want supplied bytes, empty", got, stderr.String())
	}
}

func projectInput(t *testing.T) []byte {
	t.Helper()
	document, err := evidence.NewReportV1(evidence.ReportInput{
		Subject: "projection fixture",
		Provenance: evidence.Provenance{
			AccountingDigest:  evidence.Digest(strings.Repeat("a", 64)),
			InputPolicyDigest: evidence.Digest(strings.Repeat("b", 64)),
			InventoryDigest:   evidence.Digest(strings.Repeat("c", 64)),
			Revisions: evidence.RevisionIdentity{
				Base: strings.Repeat("d", 40),
				Head: strings.Repeat("e", 40),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return document.CanonicalBytes()
}

func TestRunProjectHumanTextAndRepeatedInput(t *testing.T) {
	input := projectInput(t)
	document, err := evidence.DecodeCanonical(input)
	if err != nil {
		t.Fatal(err)
	}
	want, err := evidence.ProjectEvidence(document, evidence.ProjectionHumanText)
	if err != nil {
		t.Fatal(err)
	}
	project := func() []byte {
		var stdout, stderr bytes.Buffer
		if code := runProject([]string{"human-text"}, bytes.NewReader(input), &stdout, &stderr); code != 0 || stderr.Len() != 0 {
			t.Fatalf("exit, stderr = %d, %q", code, stderr.String())
		}
		return stdout.Bytes()
	}
	first, second := project(), project()
	if !bytes.Equal(first, want) || !bytes.Equal(second, first) {
		t.Fatalf("outputs = %q, %q; want byte-identical core output %q", first, second, want)
	}
	if bytes.Contains(first, []byte("projection fixture")) {
		t.Fatalf("human output includes subject: %q", first)
	}
}

func TestRunProjectRejectsInvalidInputWithoutWritingStdout(t *testing.T) {
	input := projectInput(t)
	for _, test := range []struct {
		name  string
		args  []string
		input []byte
	}{
		{"missing format", nil, input},
		{"extra argument", []string{"canonical-json", "extra"}, input},
		{"unknown format", []string{"other"}, input},
		{"noncanonical document", []string{"canonical-json"}, append(input, ' ')},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := runProject(test.args, bytes.NewReader(test.input), &stdout, &stderr); code != exitInvalid || stdout.Len() != 0 || stderr.String() != "invalid input\n" {
				t.Fatalf("exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestRunProjectBoundsAndReadFailures(t *testing.T) {
	oversized := &countingReader{reader: bytes.NewReader(bytes.Repeat([]byte("x"), int(projectInputCap)+2))}
	var stdout, stderr bytes.Buffer
	if code := runProject([]string{"canonical-json"}, oversized, &stdout, &stderr); code != exitBound || stdout.Len() != 0 || stderr.String() != "resource bound exceeded\n" {
		t.Fatalf("oversized exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
	}
	if got, want := oversized.read, int(projectInputCap+1); got != want {
		t.Fatalf("oversized bytes read = %d, want %d", got, want)
	}

	stdout.Reset()
	stderr.Reset()
	if code := runProject([]string{"canonical-json"}, failingReader{errors.New("read failure")}, &stdout, &stderr); code != exitAbsent || stdout.Len() != 0 || stderr.String() != "content unavailable\n" {
		t.Fatalf("read error exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
	}
}

type countingReader struct {
	reader *bytes.Reader
	read   int
}

func (r *countingReader) Read(value []byte) (int, error) {
	count, err := r.reader.Read(value)
	r.read += count
	return count, err
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

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
	if code := run(context.Background(), repository, []string{"HEAD~1", "source", "new"}, &stdout, &stderr, gitRunner, standardLimits()); code != 0 {
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
	ctx, cancel := context.WithDeadline(context.Background(), limit)
	defer cancel()
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
		got := run(context.Background(), "repository", args, &stdout, &stderr, runner, standardLimits())
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
	if got := run(canceled, "repository", []string{"base", "source", "new"}, &bytes.Buffer{}, &bytes.Buffer{}, valid.runner, standardLimits()); got != exitStopped {
		t.Fatalf("canceled exit = %d", got)
	}
}

type fixture struct {
	blobs, objects                 map[string]string
	base, head, baseAlias, nonBlob string
	headResolves                   int
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
		if revision == "base" || revision == f.baseAlias || revision == f.base {
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
