package changeevidence

// SnapshotContentV1 records byte facts, not accounting LOC or a binary verdict.
type SnapshotContentV1 struct {
	ByteLength   uint64
	ContainsNUL  bool
	LogicalLines uint64
}

// MeasureSnapshotContentV1 measures bytes without decoding or modifying them.
// LogicalLines counts LF bytes plus a nonempty unterminated tail. NUL is an
// independent observation; these facts do not decide text countability.
func MeasureSnapshotContentV1(content []byte) SnapshotContentV1 {
	result := SnapshotContentV1{ByteLength: uint64(len(content))}
	for _, value := range content {
		if value == 0 {
			result.ContainsNUL = true
		}
		if value == '\n' {
			result.LogicalLines++
		}
	}
	if len(content) > 0 && content[len(content)-1] != '\n' {
		result.LogicalLines++
	}
	return result
}
