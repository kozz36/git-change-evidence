package gitadapter

import (
	"bytes"
	"strconv"
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
)

func parseSnapshotNumstatMetadata(raw []byte) (map[snapshotMetadataKey]evidence.SnapshotNumstatV1, bool) {
	measurements := make(map[snapshotMetadataKey]evidence.SnapshotNumstatV1)
	if len(raw) == 0 {
		return measurements, true
	}
	if raw[len(raw)-1] != 0 {
		return nil, false
	}
	parts := bytes.Split(raw, []byte{0})
	last := len(parts) - 1
	for i := 0; i < last; i++ {
		fields := strings.SplitN(string(parts[i]), "\t", 3)
		if len(fields) != 3 {
			return nil, false
		}
		key := snapshotMetadataKey{path: fields[2]}
		if key.path == "" {
			if i+2 >= last || len(parts[i+1]) == 0 || len(parts[i+2]) == 0 {
				return nil, false
			}
			key = snapshotMetadataKey{previous: string(parts[i+1]), path: string(parts[i+2])}
			i += 2
		}
		if _, duplicate := measurements[key]; duplicate {
			return nil, false
		}
		measurement, valid := snapshotMetadataCounts(fields[0], fields[1])
		if !valid {
			return nil, false
		}
		measurements[key] = measurement
	}
	return measurements, true
}

func snapshotMetadataCounts(added, deleted string) (evidence.SnapshotNumstatV1, bool) {
	measurement := evidence.SnapshotNumstatV1{AddedRaw: added, DeletedRaw: deleted}
	if added == "-" || deleted == "-" {
		measurement.Kind = "binary"
		return measurement, added == "-" && deleted == "-"
	}
	additions, addedValid := snapshotMetadataCount(added)
	deletions, deletedValid := snapshotMetadataCount(deleted)
	if !addedValid || !deletedValid {
		return evidence.SnapshotNumstatV1{}, false
	}
	measurement.Kind, measurement.Added, measurement.Deleted = "text", &additions, &deletions
	return measurement, true
}

func snapshotMetadataCount(token string) (uint64, bool) {
	if len(token) == 0 || len(token) > 1 && token[0] == '0' {
		return 0, false
	}
	for i := range token {
		if token[i] < '0' || token[i] > '9' {
			return 0, false
		}
	}
	value, err := strconv.ParseUint(token, 10, 64)
	return value, err == nil
}
