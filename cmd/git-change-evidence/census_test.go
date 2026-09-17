package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/kozz36/git-change-evidence/census"
)

type censusWireDocument struct {
	Schema    string `json:"schema"`
	Extractor string `json:"extractor_version"`
	Query     struct {
		Version  string `json:"version"`
		Receiver string `json:"receiver"`
		Selector string `json:"selector"`
	} `json:"query"`
	Limits struct {
		MaxFiles      uint64 `json:"max_files"`
		MaxFileBytes  uint64 `json:"max_file_bytes"`
		MaxTotalBytes uint64 `json:"max_total_bytes"`
		MaxMatches    uint64 `json:"max_matches"`
	} `json:"limits"`
	SourceScope []struct {
		PathB64       string `json:"path_b64"`
		ContentSHA256 string `json:"content_sha256"`
		ByteLength    uint64 `json:"byte_length"`
	} `json:"source_scope"`
	Matches []struct {
		PathB64 string `json:"path_b64"`
	} `json:"matches"`
}

type censusErrorWriter struct{}

func (censusErrorWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

type censusBoundedReaderRoot struct {
	contents  map[string][]byte
	readBytes []int
	closed    bool
}

func (root *censusBoundedReaderRoot) ReadFile(name string, maxBytes uint64) ([]byte, error) {
	content, found := root.contents[name]
	if !found {
		return nil, errCensusUnavailable
	}
	read, err := io.ReadAll(io.LimitReader(bytes.NewReader(content), int64(maxBytes)+1))
	if err != nil {
		return nil, err
	}
	root.readBytes = append(root.readBytes, len(read))
	if uint64(len(read)) > maxBytes {
		return nil, errCensusBound
	}
	return read, nil
}

func (root *censusBoundedReaderRoot) Close() error {
	root.closed = true
	return nil
}

func TestCensusGoASTPinsReadBudgetToRemainingAggregate(t *testing.T) {
	first, second := []byte("package p\n"), []byte("package q\n")
	for name, test := range map[string]struct {
		totalBytes uint64
		secondRead int
	}{
		"remaining bytes": {totalBytes: uint64(len(first) + 2), secondRead: 3},
		"zero remaining":  {totalBytes: uint64(len(first)), secondRead: 1},
	} {
		t.Run(name, func(t *testing.T) {
			root := &censusBoundedReaderRoot{contents: map[string][]byte{"first.go": first, "second.go": second}}
			var stdout, stderr bytes.Buffer
			code := runCensusWith(censusArgs("/source-root", "os", "Open", "first.go", "second.go"), &stdout, &stderr, census.Limits{
				MaxFiles: 2, MaxFileBytes: 20, MaxTotalBytes: test.totalBytes, MaxMatches: 1,
			}, func(string) (censusSourceRoot, error) {
				return root, nil
			})
			if code != exitBound || stdout.Len() != 0 || stderr.String() != "resource bound exceeded\n" {
				t.Fatalf("exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
			}
			// The bounded reader receives the remaining aggregate cap and reads
			// only its one-byte overflow sentinel beyond that cap.
			if len(root.readBytes) != 2 || root.readBytes[0] != len(first) || root.readBytes[1] != test.secondRead || !root.closed {
				t.Fatalf("actual reads, closed = %v, %v; want [%d %d], true", root.readBytes, root.closed, len(first), test.secondRead)
			}
		})
	}
}

func TestRunCommandCensusGoASTDispatchesBeforeCWD(t *testing.T) {
	root := t.TempDir()
	writeCensusFile(t, root, "source.go", "package fixture\nfunc f() { os.Open(\"x\") }\n")
	var stdout, stderr bytes.Buffer
	cwdCalls := 0
	code := runCommand(context.Background(), censusArgs(root, "os", "Open", "source.go"), nil, &stdout, &stderr, func() (string, error) {
		cwdCalls++
		return "", errors.New("cwd must not be used")
	}, nil, limits{}, nil)
	if code != 0 || stderr.Len() != 0 || cwdCalls != 0 {
		t.Fatalf("exit, stderr, cwd calls = %d, %q, %d", code, stderr.String(), cwdCalls)
	}
	document := decodeCensusDocument(t, stdout.Bytes())
	if document.Schema != censusCLISchemaV1 || document.Extractor != string(census.ExtractorVersion) || document.Query.Version != string(census.SelectorCallQueryV1) || document.Query.Receiver != "os" || document.Query.Selector != "Open" || len(document.SourceScope) != 1 || len(document.Matches) != 1 {
		t.Fatalf("document = %s", stdout.Bytes())
	}
	limits := census.DefaultLimits()
	if document.Limits.MaxFiles != limits.MaxFiles || document.Limits.MaxFileBytes != limits.MaxFileBytes || document.Limits.MaxTotalBytes != limits.MaxTotalBytes || document.Limits.MaxMatches != limits.MaxMatches {
		t.Fatalf("effective limits = %#v", document.Limits)
	}
	content := []byte("package fixture\nfunc f() { os.Open(\"x\") }\n")
	hash := sha256.Sum256(content)
	if document.SourceScope[0].PathB64 != base64.StdEncoding.EncodeToString([]byte("source.go")) || document.SourceScope[0].ContentSHA256 != hex.EncodeToString(hash[:]) || document.SourceScope[0].ByteLength != uint64(len(content)) {
		t.Fatalf("source scope = %#v", document.SourceScope[0])
	}
}

func TestRunCLICensusGoASTMalformedDoesNotFallThrough(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runCLI(context.Background(), "repository", []string{"census-go-ast", "unexpected"}, nil, &stdout, &stderr, nil, limits{}, nil); code != exitInvalid || stdout.Len() != 0 || stderr.String() != "invalid input\n" {
		t.Fatalf("exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
	}
}

func TestCensusGoASTDeterministicAcrossPathOrder(t *testing.T) {
	root := t.TempDir()
	writeCensusFile(t, root, "z.go", "package fixture\nfunc z() { os.Open(\"z\") }\n")
	writeCensusFile(t, root, "a.go", "package fixture\nfunc a() { os.Open(\"a\") }\n")
	firstCode, firstOut, firstErr := runCensusCommand(t, censusArgs(root, "os", "Open", "z.go", "a.go"))
	secondCode, secondOut, secondErr := runCensusCommand(t, censusArgs(root, "os", "Open", "a.go", "z.go"))
	if firstCode != 0 || secondCode != 0 || firstErr != "" || secondErr != "" || !bytes.Equal(firstOut, secondOut) {
		t.Fatalf("first = (%d, %s, %q), second = (%d, %s, %q)", firstCode, firstOut, firstErr, secondCode, secondOut, secondErr)
	}
}

func TestCensusGoASTRejectsMalformedArgumentsBeforeFilesystemAccess(t *testing.T) {
	root := t.TempDir()
	writeCensusFile(t, root, "source.go", "package fixture\n")
	for name, args := range map[string][]string{
		"empty":              {"census-go-ast"},
		"unknown flag":       {"census-go-ast", "--unknown", "x", "--"},
		"missing flag value": {"census-go-ast", "--source-root", root, "--receiver", "os", "--selector", "--"},
		"missing separator":  {"census-go-ast", "--source-root", root, "--receiver", "os", "--selector", "Open", "source.go"},
		"noncanonical root":  {"census-go-ast", "--source-root", root + "/", "--receiver", "os", "--selector", "Open", "--", "source.go"},
		"no paths":           {"census-go-ast", "--source-root", root, "--receiver", "os", "--selector", "Open", "--"},
		"bad receiver":       censusArgs(root, "not-valid", "Open", "source.go"),
		"bad selector":       censusArgs(root, "os", "not-valid", "source.go"),
		"traversal":          censusArgs(root, "os", "Open", "../outside.go"),
		"absolute":           censusArgs(root, "os", "Open", "/outside.go"),
		"noncanonical":       censusArgs(root, "os", "Open", "dir/../source.go"),
		"duplicate":          censusArgs(root, "os", "Open", "source.go", "source.go"),
		"non Go":             censusArgs(root, "os", "Open", "source.txt"),
	} {
		t.Run(name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			cwdCalls := 0
			code := runCommand(context.Background(), args, nil, &stdout, &stderr, func() (string, error) {
				cwdCalls++
				return "", errors.New("cwd must not be used")
			}, nil, limits{}, nil)
			if code != exitInvalid || stdout.Len() != 0 || stderr.String() != "invalid input\n" || cwdCalls != 0 {
				t.Fatalf("exit, stdout, stderr, cwd calls = %d, %q, %q, %d", code, stdout.String(), stderr.String(), cwdCalls)
			}
		})
	}
}

func TestCensusGoASTRejectsUnavailableAndParserInputs(t *testing.T) {
	root := t.TempDir()
	writeCensusFile(t, root, "bad.go", "package {\n")
	for name, test := range map[string]struct {
		paths      []string
		code       int
		stderrText string
	}{
		"missing": {paths: []string{"missing.go"}, code: exitAbsent, stderrText: "content unavailable\n"},
		"parser":  {paths: []string{"bad.go"}, code: exitInvalid, stderrText: "invalid input\n"},
	} {
		t.Run(name, func(t *testing.T) {
			code, stdout, stderr := runCensusCommand(t, censusArgs(root, "os", "Open", test.paths...))
			if code != test.code || len(stdout) != 0 || stderr != test.stderrText {
				t.Fatalf("exit, stdout, stderr = %d, %q, %q", code, stdout, stderr)
			}
		})
	}
}

func TestCensusGoASTDefaultLimitsRejectBoundedInputs(t *testing.T) {
	root := t.TempDir()
	writeCensusFile(t, root, "one.go", "package fixture\nfunc f() { os.Open(\"one\") }\n")
	writeCensusFile(t, root, "two.go", "package fixture\nfunc f() { os.Open(\"two\") }\n")
	for name, test := range map[string]struct {
		bounds census.Limits
		paths  []string
	}{
		"files":       {census.Limits{MaxFiles: 1, MaxFileBytes: 1024, MaxTotalBytes: 1024, MaxMatches: 8}, []string{"one.go", "two.go"}},
		"file bytes":  {census.Limits{MaxFiles: 8, MaxFileBytes: 4, MaxTotalBytes: 1024, MaxMatches: 8}, []string{"one.go"}},
		"total bytes": {census.Limits{MaxFiles: 8, MaxFileBytes: 1024, MaxTotalBytes: 4, MaxMatches: 8}, []string{"one.go"}},
		"matches":     {census.Limits{MaxFiles: 8, MaxFileBytes: 1024, MaxTotalBytes: 1024, MaxMatches: 1}, []string{"one.go", "two.go"}},
	} {
		t.Run(name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runCensusWith(censusArgs(root, "os", "Open", test.paths...), &stdout, &stderr, test.bounds, openCensusSourceRoot)
			if code != exitBound || stdout.Len() != 0 || stderr.String() != "resource bound exceeded\n" {
				t.Fatalf("exit, stdout, stderr = %d, %q, %q", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestCensusGoASTReturnsZeroMatchesAndIgnoresDecoys(t *testing.T) {
	root := t.TempDir()
	writeCensusFile(t, root, "zero.go", "package fixture\n// os.Open(\"comment\")\nvar _ = \"os.Open(\\\"string\\\")\"\nfunc f() {}\n")
	code, stdout, stderr := runCensusCommand(t, censusArgs(root, "os", "Open", "zero.go"))
	if code != 0 || stderr != "" || len(decodeCensusDocument(t, stdout).Matches) != 0 {
		t.Fatalf("exit, stdout, stderr = %d, %s, %q", code, stdout, stderr)
	}
	writeCensusFile(t, root, "decoys.go", "package fixture\n// os.Open(\"comment\")\nvar _ = \"os.Open(\\\"string\\\")\"\nfunc f() { os.Open(\"call\") }\n")
	code, stdout, stderr = runCensusCommand(t, censusArgs(root, "os", "Open", "decoys.go"))
	if code != 0 || stderr != "" || len(decodeCensusDocument(t, stdout).Matches) != 1 {
		t.Fatalf("exit, stdout, stderr = %d, %s, %q", code, stdout, stderr)
	}
}

func TestCensusGoASTPreservesNonUTF8PathsAndUnicodeIdentifiers(t *testing.T) {
	root := t.TempDir()
	rawPath := string([]byte{0xff, 'x', '.', 'g', 'o'})
	writeCensusFile(t, root, rawPath, "package fixture\nfunc f() { os.Open(\"x\") }\n")
	code, stdout, stderr := runCensusCommand(t, censusArgs(root, "os", "Open", rawPath))
	if code != 0 || stderr != "" {
		t.Fatalf("exit, stdout, stderr = %d, %s, %q", code, stdout, stderr)
	}
	document := decodeCensusDocument(t, stdout)
	decoded, err := base64.StdEncoding.DecodeString(document.SourceScope[0].PathB64)
	if err != nil || !bytes.Equal(decoded, []byte(rawPath)) {
		t.Fatalf("path bytes = %x, error = %v", decoded, err)
	}
	writeCensusFile(t, root, "unicode.go", "package fixture\nfunc f() { π.Écho() }\n")
	code, stdout, stderr = runCensusCommand(t, censusArgs(root, "π", "Écho", "unicode.go"))
	if code != 0 || stderr != "" || len(decodeCensusDocument(t, stdout).Matches) != 1 {
		t.Fatalf("exit, stdout, stderr = %d, %s, %q", code, stdout, stderr)
	}
}

func TestCensusGoASTOutputWriteFailureHasNoSuccessDocument(t *testing.T) {
	root := t.TempDir()
	writeCensusFile(t, root, "source.go", "package fixture\nfunc f() { os.Open(\"x\") }\n")
	var stderr bytes.Buffer
	code := runCommand(context.Background(), censusArgs(root, "os", "Open", "source.go"), nil, censusErrorWriter{}, &stderr, func() (string, error) {
		return "", errors.New("cwd must not be used")
	}, nil, limits{}, nil)
	if code != exitAbsent || stderr.String() != "content unavailable\n" {
		t.Fatalf("exit, stderr = %d, %q", code, stderr.String())
	}
}

func censusArgs(root, receiver, selector string, paths ...string) []string {
	args := []string{"census-go-ast", "--source-root", root, "--receiver", receiver, "--selector", selector, "--"}
	return append(args, paths...)
}

func runCensusCommand(t *testing.T, args []string) (int, []byte, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCommand(context.Background(), args, nil, &stdout, &stderr, func() (string, error) {
		return "", errors.New("cwd must not be used")
	}, nil, limits{}, nil)
	return code, stdout.Bytes(), stderr.String()
}

func decodeCensusDocument(t *testing.T, raw []byte) censusWireDocument {
	t.Helper()
	var document censusWireDocument
	if err := json.Unmarshal(raw, &document); err != nil {
		t.Fatalf("decode census JSON: %v\n%s", err, raw)
	}
	return document
}

func writeCensusFile(t *testing.T, root, name, content string) {
	t.Helper()
	fullPath := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestCanonicalAndLegacyCensusParity(t *testing.T) {
	root := t.TempDir()
	for name, content := range map[string]string{
		"match.go":   "package fixture\nfunc f() { os.Open(\"x\") }\n",
		"zero.go":    "package fixture\nfunc f() {}\n",
		"a.go":       "package fixture\nfunc a() { os.Open(\"a\") }\n",
		"z.go":       "package fixture\nfunc z() { os.Open(\"z\") }\n",
		"unicode.go": "package fixture\nfunc f() { π.Écho(\"x\") }\n",
		"bad.go":     "package {\n",
	} {
		writeCensusFile(t, root, name, content)
	}
	rawPath := string([]byte{0xff, 'x', '.', 'g', 'o'})
	writeCensusFile(t, root, rawPath, "package fixture\nfunc f() { os.Open(\"x\") }\n")
	for name, test := range map[string]struct {
		receiver, selector string
		paths              []string
	}{
		"match":          {"os", "Open", []string{"match.go"}},
		"zero":           {"os", "Open", []string{"zero.go"}},
		"reversed paths": {"os", "Open", []string{"z.go", "a.go"}},
		"unicode":        {"π", "Écho", []string{"unicode.go"}},
		"non-UTF-8":      {"os", "Open", []string{rawPath}},
	} {
		t.Run(name, func(t *testing.T) {
			legacy := censusArgs(root, test.receiver, test.selector, test.paths...)
			canonical := append([]string{"census", "go-ast"}, legacy[1:]...)
			legacyCode, legacyOut, legacyErr := runCensusCommand(t, legacy)
			canonicalCode, canonicalOut, canonicalErr := runCensusCommand(t, canonical)
			if legacyCode != 0 || canonicalCode != legacyCode || !bytes.Equal(canonicalOut, legacyOut) || canonicalErr != legacyErr {
				t.Fatalf("canonical = (%d, %q, %q), legacy = (%d, %q, %q)", canonicalCode, canonicalOut, canonicalErr, legacyCode, legacyOut, legacyErr)
			}
		})
	}
	duplicate := []string{
		"census-go-ast",
		"--source-root", root,
		"--receiver", "os",
		"--receiver", "os",
		"--selector", "Open",
		"--", "match.go",
	}
	for name, legacy := range map[string][]string{
		"invalid flag":       {"census-go-ast", "--invalid"},
		"missing options":    {"census-go-ast", "--source-root", root, "--"},
		"missing separator":  {"census-go-ast", "--source-root", root, "--receiver", "os", "--selector", "Open", "match.go"},
		"duplicate":          duplicate,
		"invalid path":       censusArgs(root, "os", "Open", "../outside.go"),
		"invalid identifier": censusArgs(root, "not-valid", "Open", "match.go"),
		"unavailable":        censusArgs(root, "os", "Open", "missing.go"),
		"parser":             censusArgs(root, "os", "Open", "bad.go"),
	} {
		t.Run(name, func(t *testing.T) {
			if name == "duplicate" {
				if _, err := parseCensusConfig(legacy); err == nil {
					t.Fatal("duplicate option unexpectedly parsed")
				}
				withoutDuplicate := append([]string(nil), legacy[:5]...)
				withoutDuplicate = append(withoutDuplicate, legacy[7:]...)
				if _, err := parseCensusConfig(withoutDuplicate); err != nil {
					t.Fatalf("duplicate removal did not isolate the option error: %v", err)
				}
			}
			canonical := append([]string{"census", "go-ast"}, legacy[1:]...)
			legacyCode, legacyOut, legacyErr := runCensusCommand(t, legacy)
			canonicalCode, canonicalOut, canonicalErr := runCensusCommand(t, canonical)
			if canonicalCode != legacyCode || !bytes.Equal(canonicalOut, legacyOut) || canonicalErr != legacyErr || len(canonicalOut) != 0 {
				t.Fatalf("canonical = (%d, %q, %q), legacy = (%d, %q, %q)", canonicalCode, canonicalOut, canonicalErr, legacyCode, legacyOut, legacyErr)
			}
		})
	}
	legacy := censusArgs(root, "os", "Open", "match.go")
	canonical := append([]string{"census", "go-ast"}, legacy[1:]...)
	for _, args := range [][]string{legacy, canonical} {
		var stderr bytes.Buffer
		if code := runCommand(context.Background(), args, nil, censusErrorWriter{}, &stderr, func() (string, error) { t.Fatal("cwd used"); return "", nil }, nil, limits{}, nil); code != exitAbsent || stderr.String() != "content unavailable\n" {
			t.Fatalf("writer failure = (%d, %q)", code, stderr.String())
		}
	}
	boundedLegacy := censusArgs(root, "os", "Open", "match.go", "zero.go")
	boundedCanonical := append([]string{"census", "go-ast"}, boundedLegacy[1:]...)
	for _, args := range [][]string{boundedLegacy, normalizeCensusArgs(boundedCanonical)} {
		var stdout, stderr bytes.Buffer
		if code := runCensusWith(args, &stdout, &stderr, census.Limits{MaxFiles: 1, MaxFileBytes: 1, MaxTotalBytes: 1, MaxMatches: 1}, func(string) (censusSourceRoot, error) { t.Fatal("source root opened"); return nil, nil }); code != exitBound || stdout.Len() != 0 || stderr.String() != "resource bound exceeded\n" {
			t.Fatalf("resource bound = (%d, %q, %q)", code, stdout.String(), stderr.String())
		}
	}
}

var _ io.Writer = censusErrorWriter{}

func TestNormalizeCensusArgs(t *testing.T) {
	valid := censusArgs("/source-root", "os", "Open", "source.go")
	for _, test := range []struct {
		name       string
		args, want []string
	}{
		{"empty tail", []string{"census", "go-ast"}, []string{"census-go-ast"}},
		{"malformed tail", []string{"census", "go-ast", "--invalid"}, []string{"census-go-ast", "--invalid"}},
		{"three tokens", []string{"census", "go-ast", "new-file"}, []string{"census-go-ast", "new-file"}},
		{"valid tail", append([]string{"census", "go-ast"}, valid[1:]...), valid},
		{"separator and raw bytes", []string{"census", "go-ast", "--", "\xffraw"}, []string{"census-go-ast", "--", "\xffraw"}},
		{"legacy unchanged", []string{"census-go-ast", "--invalid"}, []string{"census-go-ast", "--invalid"}},
		{"unrelated unchanged", []string{"census", "other", "new-file"}, []string{"census", "other", "new-file"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			original := append([]string(nil), test.args...)
			got := normalizeCensusArgs(test.args)
			if !slices.Equal(got, test.want) || !slices.Equal(test.args, original) || !slices.Equal(normalizeCensusArgs(got), got) {
				t.Fatalf("normalized, original, idempotent = %q, %q, %q", got, test.args, normalizeCensusArgs(got))
			}
			if len(test.args) > 1 && test.args[0] == "census" && test.args[1] == "go-ast" && &got[0] == &test.args[0] {
				t.Fatal("canonical normalization reused caller backing array")
			}
		})
	}
}
