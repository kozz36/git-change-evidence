package changeevidence

import "bytes"

// CheckSnapshotResultRangeV1 compares both native revisions and complete raw
// effective-path membership of initialized immutable documents, without mutation.
// Success does not authenticate Git acquisition, compare measurements/categories,
// or extend the existing Result/Report V1 provenance binding.
func CheckSnapshotResultRangeV1(snapshot SnapshotDocumentV1, result ResultDocumentV1) error {
	if len(snapshot.canonical) == 0 {
		return &ContractError{"snapshot", "invalid_document"}
	}
	if len(result.canonical) == 0 {
		return &ContractError{"result", "invalid_document"}
	}
	if snapshot.Revisions() != result.Revisions() {
		return &ContractError{"revisions", "mismatch"}
	}
	entries := snapshot.facts.entries
	if len(entries) != len(result.entries) {
		return &ContractError{"entries", "mismatch"}
	}
	for i, entry := range entries {
		if !bytes.Equal(snapshotFactPathV1(entry), result.entries[i].Path) {
			return &ContractError{"entries", "mismatch"}
		}
	}
	return nil
}
