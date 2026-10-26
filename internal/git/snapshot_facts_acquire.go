package gitadapter

import (
	"context"
	"path/filepath"
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
)

// AcquireSnapshotMetadataBounded starts a whole-range metadata acquisition.
// Its caller owns the cooperative deadline; regular side content is not measured.
func AcquireSnapshotMetadataBounded(ctx context.Context, request Request, bounds SnapshotBounds) (*SnapshotMetadataAcquisition, error) {
	return AcquireSnapshotMetadataBoundedWithRunner(ctx, snapshotCommand, request, bounds)
}

// AcquireSnapshotMetadataBoundedWithRunner retains one context and aggregate
// budget across all commands. Injected runners must honor context and avoid effects.
// A failure returns no session, and the existing acquisition APIs are unchanged.
func AcquireSnapshotMetadataBoundedWithRunner(ctx context.Context, runner Runner, request Request, bounds SnapshotBounds) (*SnapshotMetadataAcquisition, error) {
	if bounds.MaxResponseBytes <= 0 || bounds.MaxTotalBytes <= 0 || bounds.MaxEntries <= 0 {
		return nil, ErrSnapshotInvalidBounds
	}
	if ctx == nil || runner == nil {
		return nil, failure(evidence.SnapshotForeign)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !filepath.IsAbs(request.Repository) || strings.ContainsRune(request.Repository, 0) {
		return nil, failure(evidence.SnapshotForeign)
	}
	if request.Base == "" || request.Head == "" || strings.ContainsRune(request.Base, 0) || strings.ContainsRune(request.Head, 0) {
		return nil, failure(evidence.SnapshotUnresolved)
	}
	acquisition := &SnapshotMetadataAcquisition{
		commands: snapshotCommands{ctx, runner, request.Repository, bounds.MaxResponseBytes, bounds.MaxTotalBytes},
	}
	commands := &acquisition.commands
	format, err := commands.output(evidence.SnapshotForeign, "rev-parse", "--show-object-format")
	if err != nil {
		return nil, err
	}
	length := 40
	switch strings.TrimSpace(string(format)) {
	case "sha1":
	case "sha256":
		length = 64
	default:
		return nil, failure(evidence.SnapshotForeign)
	}
	base, err := commands.resolve(request.Base, length)
	if err != nil {
		return nil, err
	}
	head, err := commands.resolve(request.Head, length)
	if err != nil {
		return nil, err
	}
	raw, err := commands.output(evidence.SnapshotRacing, "diff", "--raw", "-z", "-M", "--no-abbrev", "--no-ext-diff", "--no-textconv", "--no-relative", string(base), string(head), "--")
	if err != nil {
		return nil, err
	}
	changes, valid := parseSnapshotRawMetadata(raw)
	if !valid || base == head && len(changes) != 0 {
		return nil, racing()
	}
	if len(changes) > bounds.MaxEntries {
		return nil, ErrSnapshotResourceLimit
	}
	numstat, err := commands.output(evidence.SnapshotRacing, "diff", "--numstat", "-z", "-M", "--no-ext-diff", "--no-textconv", "--no-relative", string(base), string(head), "--")
	if err != nil {
		return nil, err
	}
	records, err := ParseSnapshotMetadata(raw, numstat)
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		if !isObjectID(string(record.Change.OldObject), length) || !isObjectID(string(record.Change.NewObject), length) {
			return nil, racing()
		}
	}
	for _, id := range []evidence.GitObjectID{base, head} {
		if _, err := commands.output(evidence.SnapshotRacing, "cat-file", "-e", string(id)+"^{commit}"); err != nil {
			return nil, err
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	acquisition.revisions = evidence.RevisionIdentity{Base: string(base), Head: string(head)}
	acquisition.records = records
	return acquisition, nil
}
