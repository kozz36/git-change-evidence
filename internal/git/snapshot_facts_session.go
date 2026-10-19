package gitadapter

import evidence "github.com/kozz36/git-change-evidence"

// SnapshotMetadataAcquisition owns one invocation's commands, budget and context.
// Records remain unmeasured metadata, not a complete validated facts document.
// Subsequent acquisition phases must use this same owner and be serialized.
type SnapshotMetadataAcquisition struct {
	commands  snapshotCommands
	revisions evidence.RevisionIdentity
	records   []SnapshotMetadataRecord
}

// Revisions returns the native immutable commit identities resolved once.
func (acquisition *SnapshotMetadataAcquisition) Revisions() evidence.RevisionIdentity {
	if acquisition == nil {
		return evidence.RevisionIdentity{}
	}
	return acquisition.revisions
}

// Records returns independently owned metadata values and numeric pointers.
func (acquisition *SnapshotMetadataAcquisition) Records() []SnapshotMetadataRecord {
	if acquisition == nil {
		return nil
	}
	records := make([]SnapshotMetadataRecord, len(acquisition.records))
	for i, record := range acquisition.records {
		record.Numstat.Added = copySnapshotMetadataCount(record.Numstat.Added)
		record.Numstat.Deleted = copySnapshotMetadataCount(record.Numstat.Deleted)
		records[i] = record
	}
	return records
}

func copySnapshotMetadataCount(count *uint64) *uint64 {
	if count == nil {
		return nil
	}
	value := *count
	return &value
}
