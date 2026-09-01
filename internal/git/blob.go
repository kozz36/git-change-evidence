package gitadapter

import (
	"bytes"
	"errors"
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
)

type BlobRequest struct{ Repository, Revision, Path string }

type Blob struct {
	Commit, Object evidence.GitObjectID
	Mode, Type     string
	Kind           evidence.GitEntryKind
	Content        string
}

type BlobErrorCode string

const (
	BlobInvalidPath BlobErrorCode = "invalid_path"
	BlobUnresolved  BlobErrorCode = "unresolved_revision"
	BlobForeign     BlobErrorCode = "foreign_repository"
	BlobMissing     BlobErrorCode = "missing_object"
	BlobNotBlob     BlobErrorCode = "non_blob"
	BlobRacing      BlobErrorCode = "racing_object"
)

type BlobError struct{ Code BlobErrorCode }

func (e *BlobError) Error() string { return "blob acquisition: " + string(e.Code) }

func AcquireBlob(request BlobRequest) (Blob, error) { return acquireBlob(request, nil) }

func acquireBlob(request BlobRequest, afterResolve func()) (Blob, error) {
	if !validBlobPath(request.Path) {
		return Blob{}, blobFailure(BlobInvalidPath)
	}
	format, err := gitOutput(request.Repository, "rev-parse", "--show-object-format")
	if err != nil {
		return Blob{}, blobFailure(BlobForeign)
	}
	length := 40
	if strings.TrimSpace(string(format)) == "sha256" {
		length = 64
	}
	commit, err := resolve(request.Repository, request.Revision, length)
	if err != nil {
		return Blob{}, mapBlobResolveFailure(err)
	}
	if afterResolve != nil {
		afterResolve()
	}
	entry, err := blobTreeEntry(request.Repository, commit, request.Path, length)
	if err != nil {
		return Blob{}, err
	}
	if entry.typ != "blob" {
		return Blob{}, blobFailure(BlobNotBlob)
	}
	content, err := gitOutput(request.Repository, "cat-file", "blob", string(entry.object))
	if err != nil || !exists(request.Repository, commit) || !blobExists(request.Repository, entry.object) {
		return Blob{}, blobFailure(BlobRacing)
	}
	return Blob{Commit: commit, Object: entry.object, Mode: entry.mode, Type: entry.typ, Kind: kind(entry.mode), Content: string(content)}, nil
}

type blobTreeRecord struct {
	mode, typ string
	object    evidence.GitObjectID
}

func blobTreeEntry(repository string, commit evidence.GitObjectID, path string, length int) (blobTreeRecord, error) {
	output, err := gitOutput(repository, "ls-tree", "-z", "--full-tree", string(commit), "--", ":(top,literal)"+path)
	if err != nil {
		return blobTreeRecord{}, blobFailure(BlobRacing)
	}
	records := bytes.Split(output, []byte{0})
	if len(records) != 2 || len(records[0]) == 0 || len(records[1]) != 0 {
		return blobTreeRecord{}, blobFailure(BlobMissing)
	}
	separator := bytes.IndexByte(records[0], '\t')
	if separator < 0 || string(records[0][separator+1:]) != path {
		return blobTreeRecord{}, blobFailure(BlobMissing)
	}
	fields := strings.Fields(string(records[0][:separator]))
	if len(fields) != 3 || !validBlobMode(fields[0]) || !isObjectID(fields[2], length) {
		return blobTreeRecord{}, blobFailure(BlobRacing)
	}
	return blobTreeRecord{mode: fields[0], typ: fields[1], object: evidence.GitObjectID(fields[2])}, nil
}

func validBlobPath(path string) bool {
	if path == "" || strings.HasPrefix(path, "/") || strings.IndexByte(path, 0) >= 0 {
		return false
	}
	for _, element := range strings.Split(path, "/") {
		if element == "" || element == "." || element == ".." {
			return false
		}
	}
	return true
}

func validBlobMode(mode string) bool {
	if len(mode) != 6 {
		return false
	}
	for index := range mode {
		if mode[index] < '0' || mode[index] > '7' {
			return false
		}
	}
	return true
}

func blobExists(repository string, object evidence.GitObjectID) bool {
	_, err := gitOutput(repository, "cat-file", "-e", string(object)+"^{blob}")
	return err == nil
}

func blobFailure(code BlobErrorCode) error { return &BlobError{Code: code} }

func mapBlobResolveFailure(err error) error {
	var snapshot *evidence.SnapshotError
	if !errors.As(err, &snapshot) {
		return blobFailure(BlobRacing)
	}
	switch snapshot.Code {
	case evidence.SnapshotUnresolved:
		return blobFailure(BlobUnresolved)
	case evidence.SnapshotForeign:
		return blobFailure(BlobForeign)
	case evidence.SnapshotMissing:
		return blobFailure(BlobMissing)
	default:
		return blobFailure(BlobRacing)
	}
}
