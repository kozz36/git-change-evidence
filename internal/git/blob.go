package gitadapter

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
)

type BlobRequest struct{ Repository, Revision, Path string }

const MaxAcquiredBlobBytes = 16 << 20

type BlobBounds struct {
	// ByteCap is per acquisition; callers may pass a decreasing W3a aggregate remainder.
	ByteCap int
}

type Blob struct {
	Commit, Object evidence.GitObjectID
	Mode, Type     string
	Kind           evidence.GitEntryKind
	Content        string
}

type BlobErrorCode string

const (
	BlobInvalidPath  BlobErrorCode = "invalid_path"
	BlobUnresolved   BlobErrorCode = "unresolved_revision"
	BlobForeign      BlobErrorCode = "foreign_repository"
	BlobMissing      BlobErrorCode = "missing_object"
	BlobNotBlob      BlobErrorCode = "non_blob"
	BlobRacing       BlobErrorCode = "racing_object"
	BlobInvalidInput BlobErrorCode = "invalid_input"
	BlobOversized    BlobErrorCode = "oversized"
	BlobCanceled     BlobErrorCode = "canceled"
	BlobDeadline     BlobErrorCode = "deadline_exceeded"
)

type BlobError struct{ Code BlobErrorCode }

func (e *BlobError) Error() string { return "blob acquisition: " + string(e.Code) }

func AcquireBlob(ctx context.Context, runner Runner, request BlobRequest, bounds BlobBounds) (Blob, error) {
	if ctx == nil || runner == nil || bounds.ByteCap <= 0 || bounds.ByteCap > MaxAcquiredBlobBytes {
		return Blob{}, blobFailure(BlobInvalidInput)
	}
	if err := blobContextFailure(ctx); err != nil {
		return Blob{}, err
	}
	return acquireBlob(ctx, runner, request, bounds, nil)
}

func acquireBlob(ctx context.Context, runner Runner, request BlobRequest, bounds BlobBounds, afterResolve func()) (Blob, error) {
	if !validBlobPath(request.Path) {
		return Blob{}, blobFailure(BlobInvalidPath)
	}
	format, err := runner(ctx, request.Repository, maxGitOutputBytes, "rev-parse", "--show-object-format")
	if err != nil {
		return Blob{}, blobCommandFailure(ctx, BlobForeign)
	}
	length := 40
	if strings.TrimSpace(string(format)) == "sha256" {
		length = 64
	}
	commit, err := blobResolve(ctx, runner, request.Repository, request.Revision, length)
	if err != nil {
		return Blob{}, err
	}
	if afterResolve != nil {
		afterResolve()
	}
	entry, err := blobTreeEntry(ctx, runner, request.Repository, commit, request.Path, length)
	if err != nil {
		return Blob{}, err
	}
	if entry.typ != "blob" {
		return Blob{}, blobFailure(BlobNotBlob)
	}
	declared, err := runner(ctx, request.Repository, maxGitOutputBytes, "cat-file", "-s", string(entry.object))
	if err != nil {
		return Blob{}, blobCommandFailure(ctx, BlobRacing)
	}
	if len(declared) < 2 || declared[len(declared)-1] != '\n' || (len(declared) > 2 && declared[0] == '0') {
		return Blob{}, blobFailure(BlobRacing)
	}
	for _, digit := range declared[:len(declared)-1] {
		if digit < '0' || digit > '9' {
			return Blob{}, blobFailure(BlobRacing)
		}
	}
	size64, err := strconv.ParseUint(string(declared[:len(declared)-1]), 10, strconv.IntSize)
	if err != nil {
		return Blob{}, blobFailure(BlobRacing)
	}
	if size64 > uint64(MaxAcquiredBlobBytes) || size64 > uint64(bounds.ByteCap) {
		return Blob{}, blobFailure(BlobOversized)
	}
	size := int(size64)
	content, err := runner(ctx, request.Repository, bounds.ByteCap, "cat-file", "blob", string(entry.object))
	if err != nil || len(content) != size || !blobCommitExists(ctx, runner, request.Repository, commit) || !blobExists(ctx, runner, request.Repository, entry.object) {
		return Blob{}, blobCommandFailure(ctx, BlobRacing)
	}
	return Blob{Commit: commit, Object: entry.object, Mode: entry.mode, Type: entry.typ, Kind: kind(entry.mode), Content: string(content)}, nil
}

type blobTreeRecord struct {
	mode, typ string
	object    evidence.GitObjectID
}

func blobTreeEntry(ctx context.Context, runner Runner, repository string, commit evidence.GitObjectID, path string, length int) (blobTreeRecord, error) {
	output, err := runner(ctx, repository, maxGitOutputBytes, "ls-tree", "-z", "--full-tree", string(commit), "--", ":(top,literal)"+path)
	if err != nil {
		return blobTreeRecord{}, blobCommandFailure(ctx, BlobRacing)
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

func blobResolve(ctx context.Context, runner Runner, repository, revision string, length int) (evidence.GitObjectID, error) {
	output, err := runner(ctx, repository, maxGitOutputBytes, "rev-parse", "--verify", "--end-of-options", revision+"^{commit}")
	if err != nil {
		code := BlobUnresolved
		if isObjectID(revision, length) {
			code = BlobMissing
		}
		return "", blobCommandFailure(ctx, code)
	}
	object := strings.TrimSpace(string(output))
	if !isObjectID(object, length) {
		return "", blobFailure(BlobUnresolved)
	}
	id := evidence.GitObjectID(object)
	if !blobCommitExists(ctx, runner, repository, id) {
		return "", blobCommandFailure(ctx, BlobMissing)
	}
	return id, nil
}

func blobCommitExists(ctx context.Context, runner Runner, repository string, commit evidence.GitObjectID) bool {
	_, err := runner(ctx, repository, maxGitOutputBytes, "cat-file", "-e", string(commit)+"^{commit}")
	return err == nil
}

func blobExists(ctx context.Context, runner Runner, repository string, object evidence.GitObjectID) bool {
	_, err := runner(ctx, repository, maxGitOutputBytes, "cat-file", "-e", string(object)+"^{blob}")
	return err == nil
}

func blobContextFailure(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return blobFailure(BlobCanceled)
	case context.DeadlineExceeded:
		return blobFailure(BlobDeadline)
	default:
		return nil
	}
}

func blobCommandFailure(ctx context.Context, fallback BlobErrorCode) error {
	if err := blobContextFailure(ctx); err != nil {
		return err
	}
	return blobFailure(fallback)
}

func blobFailure(code BlobErrorCode) error { return &BlobError{Code: code} }
