package gitadapter

import (
	"context"
	"path/filepath"
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
)

// AcquireBounded acquires an immutable whole committed range without fetching,
// worktree/index writes, external diff, textconv, or inventory-path filtering.
func AcquireBounded(ctx context.Context, request Request, bounds SnapshotBounds) (evidence.CommittedSnapshot, error) {
	return AcquireBoundedWithRunner(ctx, snapshotCommand, request, bounds)
}

// AcquireBoundedWithRunner applies budgets and cancellation even to injected
// runner responses. A runner must honor context and avoid effects itself.
// Failures return no partial snapshot; existing Acquire remains unchanged.
func AcquireBoundedWithRunner(ctx context.Context, runner Runner, request Request, bounds SnapshotBounds) (evidence.CommittedSnapshot, error) {
	empty := evidence.CommittedSnapshot{}
	if bounds.MaxResponseBytes <= 0 || bounds.MaxTotalBytes <= 0 || bounds.MaxEntries <= 0 {
		return empty, ErrSnapshotInvalidBounds
	}
	if ctx == nil || runner == nil {
		return empty, failure(evidence.SnapshotForeign)
	}
	if err := ctx.Err(); err != nil {
		return empty, err
	}
	if !filepath.IsAbs(request.Repository) || strings.ContainsRune(request.Repository, 0) {
		return empty, failure(evidence.SnapshotForeign)
	}
	if request.Base == "" || request.Head == "" || strings.ContainsRune(request.Base, 0) || strings.ContainsRune(request.Head, 0) {
		return empty, failure(evidence.SnapshotUnresolved)
	}
	commands := snapshotCommands{ctx, runner, request.Repository, bounds.MaxResponseBytes, bounds.MaxTotalBytes}
	format, err := commands.output(evidence.SnapshotForeign, "rev-parse", "--show-object-format")
	if err != nil {
		return empty, err
	}
	length := 40
	switch strings.TrimSpace(string(format)) {
	case "sha1":
	case "sha256":
		length = 64
	default:
		return empty, failure(evidence.SnapshotForeign)
	}
	base, err := commands.resolve(request.Base, length)
	if err != nil {
		return empty, err
	}
	head, err := commands.resolve(request.Head, length)
	if err != nil {
		return empty, err
	}
	raw, err := commands.output(evidence.SnapshotRacing, "diff", "--raw", "-z", "-M", "--no-abbrev", "--no-ext-diff", "--no-textconv", "--no-relative", string(base), string(head), "--")
	if err != nil {
		return empty, err
	}
	entries, ok := parseRaw(raw)
	if !ok {
		return empty, racing()
	}
	if len(entries) > bounds.MaxEntries {
		return empty, ErrSnapshotResourceLimit
	}
	numstat, err := commands.output(evidence.SnapshotRacing, "diff", "--numstat", "-z", "-M", "--no-ext-diff", "--no-textconv", "--no-relative", string(base), string(head), "--")
	if err != nil {
		return empty, err
	}
	measurements, ok := parseNumstat(numstat)
	if !ok || len(measurements) != len(entries) {
		return empty, racing()
	}
	for index := range entries {
		measurement, found := measurements[entryKey(entries[index])]
		if !found {
			return empty, racing()
		}
		entries[index].Binary, entries[index].Lines = measurement.Binary, measurement.Lines
	}
	for _, id := range []evidence.GitObjectID{base, head} {
		if _, err := commands.output(evidence.SnapshotRacing, "cat-file", "-e", string(id)+"^{commit}"); err != nil {
			return empty, err
		}
	}
	if err := ctx.Err(); err != nil {
		return empty, err
	}
	return evidence.NewCommittedSnapshot(base, head, entries), nil
}

func (commands *snapshotCommands) resolve(revision string, length int) (evidence.GitObjectID, error) {
	code := evidence.SnapshotUnresolved
	if isObjectID(revision, length) {
		code = evidence.SnapshotMissing
	}
	output, err := commands.output(code, "rev-parse", "--verify", "--end-of-options", revision+"^{commit}")
	if err != nil {
		return "", err
	}
	object := strings.TrimSpace(string(output))
	if !isObjectID(object, length) {
		return "", failure(evidence.SnapshotUnresolved)
	}
	if _, err := commands.output(evidence.SnapshotMissing, "cat-file", "-e", object+"^{commit}"); err != nil {
		return "", err
	}
	return evidence.GitObjectID(object), nil
}
