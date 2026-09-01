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
	"time"
	"unicode/utf8"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

const exitInvalid, exitAbsent, exitBound, exitUTF8, exitStopped, exitRace = 2, 3, 4, 5, 6, 7

type limits struct{ perBlob, aggregate, lines, lineBytes, pairs int }

var standardLimits = limits{gitadapter.MaxAcquiredBlobBytes, gitadapter.MaxAcquiredBlobBytes, evidence.CarveoutSafeMaxLines, evidence.CarveoutSafeMaxLineBytes, 1 << 20}

type absoluteClock struct {
	now             func() time.Time
	start, deadline time.Time
}

func (c *absoluteClock) Now() time.Time {
	if start := c.start; !start.IsZero() {
		c.start = time.Time{}
		return start
	}
	now := c.now()
	if now.After(c.deadline) {
		return c.deadline
	}
	return now
}
func matcherDeadline(ctx context.Context, now time.Time) time.Time {
	if limit, ok := ctx.Deadline(); ok && limit.Before(now.Add(time.Second)) {
		return limit
	}
	return now.Add(time.Second)
}
func main() {
	repository, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, "invalid input\n")
		os.Exit(exitInvalid)
	}
	os.Exit(run(context.Background(), repository, os.Args[1:], os.Stdout, os.Stderr, gitRunner, standardLimits))
}
func run(ctx context.Context, repository string, args []string, stdout, stderr io.Writer, runner gitadapter.Runner, bound limits) int {
	if len(args) != 3 {
		fmt.Fprint(stderr, "usage: git-change-evidence <base-ref> <source-file> <new-file>\n")
		return exitInvalid
	}
	moved, fresh, code := measure(ctx, repository, args, runner, bound)
	if code != 0 {
		fmt.Fprint(stderr, diagnostic(code))
		return code
	}
	fmt.Fprintf(stdout, "moved=%d new=%d additions=%d\n", moved, fresh, moved+fresh)
	return 0
}
func measure(ctx context.Context, repository string, args []string, runner gitadapter.Runner, bound limits) (int, int, int) {
	if runner == nil || bound.perBlob <= 0 || bound.aggregate <= 0 || bound.lines <= 0 || bound.lineBytes <= 0 || bound.pairs <= 0 {
		return 0, 0, exitBound
	}
	remaining, lines, pairs := bound.aggregate, bound.lines, bound.pairs
	deadline := matcherDeadline(ctx, time.Now())
	base, code := acquire(ctx, runner, repository, args[0], args[1], &remaining, bound.perBlob)
	if code != 0 {
		return 0, 0, code
	}
	afterSource, code := acquire(ctx, runner, repository, "HEAD", args[1], &remaining, bound.perBlob)
	if code != 0 {
		return 0, 0, code
	}
	afterNew, code := acquire(ctx, runner, repository, string(afterSource.Commit), args[2], &remaining, bound.perBlob)
	if code != 0 {
		return 0, 0, code
	}
	raw := []string{base.Content, afterSource.Content, afterNew.Content}
	for index := range raw {
		if !utf8.ValidString(raw[index]) {
			return 0, 0, exitUTF8
		}
		raw[index] = normalized(raw[index])
	}
	source, code := prepared(raw[0], raw[1], &lines, bound)
	if code != 0 {
		return 0, 0, code
	}
	retained, code := matched(ctx, source, &pairs, deadline, time.Now)
	if code != 0 {
		return 0, 0, code
	}
	used := make([]bool, source.BeforeLineCount())
	for index := 0; index < retained.MatchCount(); index++ {
		match, _ := retained.MatchAt(index)
		used[match.BeforeIndex] = true
	}
	deleted := make([]string, 0, len(used))
	for index := range used {
		if !used[index] {
			line, _ := source.BeforeLine(index)
			deleted = append(deleted, line)
		}
	}
	destination, code := prepared(strings.Join(deleted, "\n"), raw[2], &lines, bound)
	if code != 0 {
		return 0, 0, code
	}
	moved, code := matched(ctx, destination, &pairs, deadline, time.Now)
	if code != 0 {
		return 0, 0, code
	}
	return moved.Moved(), destination.AfterLineCount() - moved.Moved(), 0
}
func normalized(raw string) string {
	lines := make([]string, 0)
	for _, line := range strings.Split(raw, "\n") {
		if fields := strings.Fields(line); len(fields) != 0 {
			lines = append(lines, strings.Join(fields, " "))
		}
	}
	return strings.Join(lines, "\n")
}
func prepared(before, after string, lines *int, bound limits) (evidence.PreparedCarveout, int) {
	prepared := evidence.PrepareCarveout(evidence.CarveoutInput{Before: evidence.CarveoutContent{Available: true, Raw: before}, After: evidence.CarveoutContent{Available: true, Raw: after}}, evidence.CarveoutBounds{ByteCap: bound.aggregate, LineCap: *lines, MaxLineBytes: bound.lineBytes})
	if prepared.State != evidence.CarveoutPrepared {
		return evidence.PreparedCarveout{}, exitBound
	}
	*lines -= prepared.Prepared.BeforeLineCount() + prepared.Prepared.AfterLineCount()
	return prepared.Prepared, 0
}
func matched(ctx context.Context, prepared evidence.PreparedCarveout, pairs *int, deadline time.Time, now func() time.Time) (*evidence.CarveoutMatchResult, int) {
	start := now()
	if ctx.Err() != nil || !start.Before(deadline) {
		return nil, exitStopped
	}
	analysis := evidence.MatchPreparedCarveout(prepared, evidence.CarveoutMatcherOptions{PairCap: *pairs, TimeCap: deadline.Sub(start), Clock: &absoluteClock{now, start, deadline}})
	*pairs -= analysis.Provenance.PairInspections
	if ctx.Err() != nil || analysis.Provenance.Reason == evidence.CarveoutMatchDeadlineExceeded {
		return nil, exitStopped
	}
	if analysis.State == evidence.CarveoutMatchCompleted {
		return analysis.Result, 0
	}
	return nil, exitBound
}
func acquire(ctx context.Context, runner gitadapter.Runner, repository, revision, path string, remaining *int, cap int) (gitadapter.Blob, int) {
	if ctx.Err() != nil {
		return gitadapter.Blob{}, exitStopped
	}
	if *remaining <= 0 {
		return gitadapter.Blob{}, exitBound
	}
	if cap > *remaining {
		cap = *remaining
	}
	blob, err := gitadapter.AcquireBlob(ctx, runner, gitadapter.BlobRequest{Repository: repository, Revision: revision, Path: path}, gitadapter.BlobBounds{ByteCap: cap})
	if err != nil {
		return gitadapter.Blob{}, blobExit(err)
	}
	*remaining -= len(blob.Content)
	return blob, 0
}
func blobExit(err error) int {
	var failure *gitadapter.BlobError
	if !errors.As(err, &failure) {
		return exitRace
	}
	switch failure.Code {
	case gitadapter.BlobInvalidPath, gitadapter.BlobInvalidInput, gitadapter.BlobUnresolved:
		return exitInvalid
	case gitadapter.BlobMissing, gitadapter.BlobNotBlob, gitadapter.BlobForeign:
		return exitAbsent
	case gitadapter.BlobOversized:
		return exitBound
	case gitadapter.BlobCanceled, gitadapter.BlobDeadline:
		return exitStopped
	default:
		return exitRace
	}
}
func diagnostic(code int) string {
	return map[int]string{exitInvalid: "invalid input\n", exitAbsent: "content unavailable\n", exitBound: "resource bound exceeded\n", exitUTF8: "invalid UTF-8\n", exitStopped: "operation interrupted\n", exitRace: "git object race\n"}[code]
}
func gitRunner(ctx context.Context, repository string, cap int, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", append([]string{"--no-replace-objects", "-C", repository}, args...)...)
	command.Env = []string{"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_ATTR_NOSYSTEM=1", "GIT_CONFIG_COUNT=0", "GIT_NO_REPLACE_OBJECTS=1", "LC_ALL=C", "LANG=C", "PATH=/usr/bin:/bin"}
	output := bounded{cap: cap}
	command.Stdout, command.Stderr = &output, io.Discard
	if err := command.Run(); err != nil {
		return nil, err
	}
	if output.exceeded {
		return nil, errors.New("git output exceeded cap")
	}
	return output.Bytes(), nil
}

type bounded struct {
	bytes.Buffer
	cap      int
	exceeded bool
}

func (output *bounded) Write(value []byte) (int, error) {
	left := output.cap - output.Len()
	if len(value) > left {
		output.Buffer.Write(value[:left])
		output.exceeded = true
		return len(value), nil
	}
	return output.Buffer.Write(value)
}
