package gitadapter

import (
	"bytes"
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
)

// SnapshotMetadataRecord preserves parsed raw metadata and exact numstat tokens.
// It is unmeasured adapter data, not a validated complete snapshot document.
type SnapshotMetadataRecord struct {
	Change  evidence.CommittedChange
	Numstat evidence.SnapshotNumstatV1
}

type snapshotMetadataKey struct{ previous, path string }

// ParseSnapshotMetadata pairs complete NUL-framed Git responses exactly once.
// It validates framing, tokens and pairing, not Git authenticity or side metrics.
func ParseSnapshotMetadata(raw, numstat []byte) ([]SnapshotMetadataRecord, error) {
	changes, valid := parseSnapshotRawMetadata(raw)
	measurements, countsValid := parseSnapshotNumstatMetadata(numstat)
	if !valid || !countsValid || len(changes) != len(measurements) {
		return nil, &evidence.SnapshotError{Code: evidence.SnapshotRacing}
	}
	records := make([]SnapshotMetadataRecord, len(changes))
	for i, change := range changes {
		key := snapshotMetadataKey{change.PreviousPath, change.Path}
		measurement, exists := measurements[key]
		if !exists {
			return nil, &evidence.SnapshotError{Code: evidence.SnapshotRacing}
		}
		delete(measurements, key)
		change.Binary = measurement.Kind == "binary"
		if !change.Binary {
			change.Lines = evidence.CommittedLineCounts{Additions: *measurement.Added, Deletions: *measurement.Deleted, Countable: true}
		}
		records[i] = SnapshotMetadataRecord{Change: change, Numstat: measurement}
	}
	return records, nil
}

func parseSnapshotRawMetadata(raw []byte) ([]evidence.CommittedChange, bool) {
	if len(raw) == 0 {
		return nil, true
	}
	if raw[len(raw)-1] != 0 {
		return nil, false
	}
	parts := bytes.Split(raw, []byte{0})
	last := len(parts) - 1
	changes := []evidence.CommittedChange{}
	for i := 0; i < last; {
		fields := strings.Fields(string(parts[i]))
		if len(fields) != 5 || len(fields[0]) < 2 || fields[0][0] != ':' || !validSnapshotMetadataStatus(fields[4]) {
			return nil, false
		}
		change := evidence.CommittedChange{
			Status: fields[4][:1], OldMode: fields[0][1:], NewMode: fields[1],
			OldObject: evidence.GitObjectID(fields[2]), NewObject: evidence.GitObjectID(fields[3]),
		}
		change.OldKind, change.NewKind = kind(change.OldMode), kind(change.NewMode)
		i++
		pathCount := 1
		if change.Status == "R" || change.Status == "C" {
			pathCount = 2
		}
		if i+pathCount > last || len(parts[i]) == 0 {
			return nil, false
		}
		if pathCount == 2 {
			if len(parts[i+1]) == 0 {
				return nil, false
			}
			change.PreviousPath = string(parts[i])
		}
		change.Path = string(parts[i+pathCount-1])
		i += pathCount
		changes = append(changes, change)
	}
	return changes, true
}

func validSnapshotMetadataStatus(status string) bool {
	if len(status) == 0 || !strings.ContainsRune("ADMRCT", rune(status[0])) {
		return false
	}
	if len(status) == 1 {
		return true
	}
	if !strings.ContainsRune("MRC", rune(status[0])) || len(status) > 4 {
		return false
	}
	score := 0
	for i := 1; i < len(status); i++ {
		if status[i] < '0' || status[i] > '9' {
			return false
		}
		score = score*10 + int(status[i]-'0')
	}
	return score <= 100
}
