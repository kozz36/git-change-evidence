package gitadapter

import (
	"bytes"
	"os/exec"
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
)

type Request struct{ Repository, Base, Head string }

func Acquire(request Request) (evidence.CommittedSnapshot, error) { return acquire(request, nil) }

func acquire(request Request, afterRead func()) (evidence.CommittedSnapshot, error) {
	format, err := gitOutput(request.Repository, "rev-parse", "--show-object-format")
	if err != nil {
		return evidence.CommittedSnapshot{}, failure(evidence.SnapshotForeign)
	}
	length := 40
	if strings.TrimSpace(string(format)) == "sha256" {
		length = 64
	}
	base, err := resolve(request.Repository, request.Base, length)
	if err != nil {
		return evidence.CommittedSnapshot{}, err
	}
	head, err := resolve(request.Repository, request.Head, length)
	if err != nil {
		return evidence.CommittedSnapshot{}, err
	}
	raw, err := gitOutput(request.Repository, "diff", "--raw", "-z", "-M", "--no-abbrev", "--no-ext-diff", "--no-textconv", string(base), string(head), "--")
	if err != nil {
		return evidence.CommittedSnapshot{}, racing()
	}
	numstat, err := gitOutput(request.Repository, "diff", "--numstat", "-z", "-M", "--no-ext-diff", "--no-textconv", string(base), string(head), "--")
	if err != nil {
		return evidence.CommittedSnapshot{}, racing()
	}
	entries, ok := parseRaw(raw)
	binary, binaryOK := parseNumstat(numstat)
	if !ok || !binaryOK {
		return evidence.CommittedSnapshot{}, racing()
	}
	for index := range entries {
		entries[index].Binary = binary[entryKey(entries[index])]
	}
	if afterRead != nil {
		afterRead()
	}
	if !exists(request.Repository, base) || !exists(request.Repository, head) {
		return evidence.CommittedSnapshot{}, racing()
	}
	return evidence.NewCommittedSnapshot(base, head, entries), nil
}

func resolve(repository, revision string, length int) (evidence.GitObjectID, error) {
	output, err := gitOutput(repository, "rev-parse", "--verify", "--end-of-options", revision+"^{commit}")
	if err != nil {
		code := evidence.SnapshotUnresolved
		if isObjectID(revision, length) {
			code = evidence.SnapshotMissing
		}
		return "", failure(code)
	}
	object := strings.TrimSpace(string(output))
	if !isObjectID(object, length) {
		return "", failure(evidence.SnapshotUnresolved)
	}
	id := evidence.GitObjectID(object)
	if !exists(repository, id) {
		return "", failure(evidence.SnapshotMissing)
	}
	return id, nil
}

func exists(repository string, id evidence.GitObjectID) bool {
	_, err := gitOutput(repository, "cat-file", "-e", string(id)+"^{commit}")
	return err == nil
}
func racing() error                                 { return failure(evidence.SnapshotRacing) }
func failure(code evidence.SnapshotErrorCode) error { return &evidence.SnapshotError{Code: code} }

func gitOutput(repository string, args ...string) ([]byte, error) {
	arguments := append([]string{"--no-replace-objects", "-C", repository}, args...)
	command := exec.Command("git", arguments...)
	command.Env = []string{"GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_ATTR_NOSYSTEM=1", "GIT_CONFIG_COUNT=0", "GIT_NO_REPLACE_OBJECTS=1", "LC_ALL=C", "LANG=C", "PATH=/usr/bin:/bin"}
	return command.Output()
}

func parseRaw(raw []byte) ([]evidence.CommittedChange, bool) {
	records, entries := bytes.Split(raw, []byte{0}), []evidence.CommittedChange{}
	for index := 0; index+1 < len(records) && len(records[index]) != 0; {
		fields := strings.Fields(string(records[index]))
		if len(fields) != 5 || len(fields[4]) == 0 || !strings.HasPrefix(fields[0], ":") {
			return nil, false
		}
		entry := evidence.CommittedChange{Status: fields[4][:1], OldMode: fields[0][1:], NewMode: fields[1], OldObject: evidence.GitObjectID(fields[2]), NewObject: evidence.GitObjectID(fields[3])}
		entry.OldKind, entry.NewKind = kind(entry.OldMode), kind(entry.NewMode)
		index++
		if entry.Status == "R" || entry.Status == "C" {
			if index+1 >= len(records) {
				return nil, false
			}
			entry.PreviousPath, entry.Path, index = string(records[index]), string(records[index+1]), index+2
		} else {
			entry.Path, index = string(records[index]), index+1
		}
		entries = append(entries, entry)
	}
	return entries, true
}

func parseNumstat(raw []byte) (map[string]bool, bool) {
	binary, records := map[string]bool{}, bytes.Split(raw, []byte{0})
	for index := 0; index+1 < len(records) && len(records[index]) != 0; index++ {
		fields := strings.SplitN(string(records[index]), "\t", 3)
		if len(fields) != 3 {
			return nil, false
		}
		key := fields[2]
		if key == "" {
			if index+2 >= len(records) {
				return nil, false
			}
			key, index = string(records[index+1])+"\x00"+string(records[index+2]), index+2
		}
		if fields[0] == "-" && fields[1] == "-" {
			binary[key] = true
		}
	}
	return binary, true
}

func entryKey(entry evidence.CommittedChange) string {
	if entry.PreviousPath != "" {
		return entry.PreviousPath + "\x00" + entry.Path
	}
	return entry.Path
}
func kind(mode string) evidence.GitEntryKind {
	switch {
	case mode == "000000":
		return evidence.GitAbsent
	case strings.HasPrefix(mode, "100"):
		return evidence.GitFile
	case mode == "120000":
		return evidence.GitSymlink
	case mode == "160000":
		return evidence.GitGitlink
	default:
		return evidence.GitOther
	}
}
func isObjectID(value string, length int) bool {
	if len(value) != length {
		return false
	}
	for _, character := range value {
		if !(character >= '0' && character <= '9' || character >= 'a' && character <= 'f') {
			return false
		}
	}
	return true
}
