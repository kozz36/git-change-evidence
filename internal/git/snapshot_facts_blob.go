package gitadapter

import (
	"strconv"

	evidence "github.com/kozz36/git-change-evidence"
)

func (acquisition *SnapshotMetadataAcquisition) measureSideBlob(id evidence.GitObjectID) (evidence.SnapshotContentV1, error) {
	commands := &acquisition.commands
	typeBytes, err := commands.output(evidence.SnapshotRacing, "cat-file", "-t", string(id))
	if err != nil {
		return evidence.SnapshotContentV1{}, err
	}
	if string(typeBytes) != "blob\n" {
		return evidence.SnapshotContentV1{}, racing()
	}
	sizeBytes, err := commands.output(evidence.SnapshotRacing, "cat-file", "-s", string(id))
	if err != nil {
		return evidence.SnapshotContentV1{}, err
	}
	size, valid := snapshotBlobSize(sizeBytes)
	if !valid {
		return evidence.SnapshotContentV1{}, racing()
	}
	if size > uint64(commands.response) || size > uint64(commands.remaining) {
		return evidence.SnapshotContentV1{}, ErrSnapshotResourceLimit
	}
	content, err := commands.output(evidence.SnapshotRacing, "cat-file", "blob", string(id))
	if err != nil {
		return evidence.SnapshotContentV1{}, err
	}
	if uint64(len(content)) != size {
		return evidence.SnapshotContentV1{}, racing()
	}
	if _, err := commands.output(evidence.SnapshotRacing, "cat-file", "-e", string(id)+"^{blob}"); err != nil {
		return evidence.SnapshotContentV1{}, err
	}
	return evidence.MeasureSnapshotContentV1(content), nil
}

func snapshotBlobSize(raw []byte) (uint64, bool) {
	if len(raw) < 2 || raw[len(raw)-1] != '\n' || len(raw) > 2 && raw[0] == '0' {
		return 0, false
	}
	for _, digit := range raw[:len(raw)-1] {
		if digit < '0' || digit > '9' {
			return 0, false
		}
	}
	size, err := strconv.ParseUint(string(raw[:len(raw)-1]), 10, 64)
	return size, err == nil
}
