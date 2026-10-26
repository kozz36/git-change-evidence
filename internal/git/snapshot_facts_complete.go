package gitadapter

import (
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
)

// AcquireFacts reads required immutable file-side objects and validates full facts.
// Calls must be serialized; repeated reads charge the original invocation budget.
// The original caller context and deadline remain in force. Metadata is unchanged.
func (acquisition *SnapshotMetadataAcquisition) AcquireFacts() (evidence.SnapshotFactsV1, error) {
	if acquisition == nil || acquisition.commands.ctx == nil || acquisition.commands.runner == nil {
		return evidence.SnapshotFactsV1{}, failure(evidence.SnapshotForeign)
	}
	entries := make([]evidence.SnapshotEntryV1, len(acquisition.records))
	for i, record := range acquisition.records {
		if err := acquisition.commands.ctx.Err(); err != nil {
			return evidence.SnapshotFactsV1{}, err
		}
		change := record.Change
		oldPath := change.Path
		if change.Status == "R" || change.Status == "C" {
			oldPath = change.PreviousPath
		}
		oldSide, err := acquisition.factSide(oldPath, change.OldMode, change.OldKind, change.OldObject)
		if err != nil {
			return evidence.SnapshotFactsV1{}, err
		}
		newSide, err := acquisition.factSide(change.Path, change.NewMode, change.NewKind, change.NewObject)
		if err != nil {
			return evidence.SnapshotFactsV1{}, err
		}
		entries[i] = evidence.SnapshotEntryV1{Status: change.Status, Old: oldSide, New: newSide, Numstat: record.Numstat}
	}
	for _, id := range []string{acquisition.revisions.Base, acquisition.revisions.Head} {
		if _, err := acquisition.commands.output(evidence.SnapshotRacing, "cat-file", "-e", id+"^{commit}"); err != nil {
			return evidence.SnapshotFactsV1{}, err
		}
	}
	facts, err := evidence.NewSnapshotFactsV1(acquisition.revisions, entries)
	if err != nil {
		return evidence.SnapshotFactsV1{}, racing()
	}
	if err := acquisition.commands.ctx.Err(); err != nil {
		return evidence.SnapshotFactsV1{}, err
	}
	return facts, nil
}

func (acquisition *SnapshotMetadataAcquisition) factSide(path, mode string, entryKind evidence.GitEntryKind, id evidence.GitObjectID) (*evidence.SnapshotSideV1, error) {
	if mode == "000000" {
		return nil, nil
	}
	switch mode {
	case "100644", "100755", "120000", "160000":
	default:
		return nil, racing()
	}
	if entryKind != kind(mode) || !isObjectID(string(id), len(acquisition.revisions.Base)) || strings.Trim(string(id), "0") == "" {
		return nil, racing()
	}
	side := &evidence.SnapshotSideV1{Path: []byte(path), Mode: mode, Kind: entryKind, ObjectID: id}
	if entryKind == evidence.GitFile {
		content, err := acquisition.measureSideBlob(id)
		if err != nil {
			return nil, err
		}
		side.Content = &content
	}
	return side, nil
}
