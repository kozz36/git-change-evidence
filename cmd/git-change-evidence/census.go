package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go/token"
	"io"
	"path"
	"strings"

	changeevidence "github.com/kozz36/git-change-evidence"
	"github.com/kozz36/git-change-evidence/census"
)

const censusCLISchemaV1 = "git-change-evidence.census-go-ast-cli/v1"

var (
	errCensusBound       = errors.New("census resource bound exceeded")
	errCensusUnavailable = errors.New("census source unavailable")
)

type censusSourceRoot interface {
	ReadFile(name string, maxBytes uint64) ([]byte, error)
	Close() error
}

type censusConfig struct {
	sourceRoot string
	receiver   string
	selector   string
	paths      []string
}

type censusDocumentV1 struct {
	Schema           string                 `json:"schema"`
	ExtractorVersion string                 `json:"extractor_version"`
	Query            censusQueryV1          `json:"query"`
	Limits           censusLimitsV1         `json:"limits"`
	Inventory        censusInventoryV1      `json:"inventory"`
	SourceScope      []censusSourceRecordV1 `json:"source_scope"`
	Matches          []censusMatchV1        `json:"matches"`
}

type censusQueryV1 struct {
	Version  string `json:"version"`
	Receiver string `json:"receiver"`
	Selector string `json:"selector"`
}

type censusLimitsV1 struct {
	MaxFiles      uint64 `json:"max_files"`
	MaxFileBytes  uint64 `json:"max_file_bytes"`
	MaxTotalBytes uint64 `json:"max_total_bytes"`
	MaxMatches    uint64 `json:"max_matches"`
}

type censusInventoryV1 struct {
	Schema                 string `json:"schema"`
	SHA256                 string `json:"sha256"`
	AccountingPolicySHA256 string `json:"accounting_policy_sha256"`
}

type censusSourceRecordV1 struct {
	PathB64       string `json:"path_b64"`
	ContentSHA256 string `json:"content_sha256"`
	ByteLength    uint64 `json:"byte_length"`
}

type censusMatchV1 struct {
	PathB64        string       `json:"path_b64"`
	FileSHA256     string       `json:"file_sha256"`
	FragmentSHA256 string       `json:"fragment_sha256"`
	Call           censusSpanV1 `json:"call"`
	Selector       censusSpanV1 `json:"selector"`
}

type censusSpanV1 struct {
	Start censusPositionV1 `json:"start"`
	End   censusPositionV1 `json:"end"`
}

type censusPositionV1 struct {
	Offset int `json:"offset"`
	Line   int `json:"line"`
	Column int `json:"column"`
}

func runCensus(args []string, stdout, stderr io.Writer) int {
	return runCensusWith(args, stdout, stderr, census.DefaultLimits(), openCensusSourceRoot)
}

func runCensusWith(args []string, stdout, stderr io.Writer, bounds census.Limits, openRoot func(string) (censusSourceRoot, error)) int {
	config, err := parseCensusConfig(args)
	if err != nil || !validCensusLimits(bounds) || openRoot == nil {
		return censusExit(stderr, exitInvalid)
	}
	if uint64(len(config.paths)) > bounds.MaxFiles {
		return censusExit(stderr, exitBound)
	}

	root, err := openRoot(config.sourceRoot)
	if err != nil || root == nil {
		return censusExit(stderr, exitAbsent)
	}
	defer root.Close()

	entries := make([]changeevidence.InventoryEntryInput, 0, len(config.paths))
	files := make([]census.File, 0, len(config.paths))
	remaining := bounds.MaxTotalBytes
	for _, name := range config.paths {
		readLimit := bounds.MaxFileBytes
		if remaining < readLimit {
			readLimit = remaining
		}
		content, err := root.ReadFile(name, readLimit)
		if err != nil {
			if errors.Is(err, errCensusBound) {
				return censusExit(stderr, exitBound)
			}
			return censusExit(stderr, exitAbsent)
		}
		size := uint64(len(content))
		if size > bounds.MaxFileBytes || size > remaining {
			return censusExit(stderr, exitBound)
		}
		remaining -= size
		pathBytes := []byte(name)
		entries = append(entries, changeevidence.InventoryEntryInput{Path: pathBytes, Content: content})
		files = append(files, census.File{Path: pathBytes, Content: content})
	}

	policy, err := changeevidence.NewAccountingPolicyV1(changeevidence.PolicyInput{
		Categories: []changeevidence.CategoryInput{{Name: "census"}},
		Default:    "census",
	})
	if err != nil {
		return censusExit(stderr, exitInvalid)
	}
	inventory, err := changeevidence.NewInventoryV1(policy, entries)
	if err != nil {
		return censusExit(stderr, exitInvalid)
	}
	result, err := census.GoASTV1(inventory, files, census.QueryV1{
		Version:  census.SelectorCallQueryV1,
		Receiver: config.receiver,
		Selector: config.selector,
	}, bounds)
	if err != nil {
		if censusLimitError(err) {
			return censusExit(stderr, exitBound)
		}
		return censusExit(stderr, exitInvalid)
	}
	if err := writeCensusDocument(stdout, censusOutput(inventory, result, bounds)); err != nil {
		return censusExit(stderr, exitAbsent)
	}
	return 0
}

func parseCensusConfig(args []string) (censusConfig, error) {
	if len(args) == 0 || args[0] != "census-go-ast" {
		return censusConfig{}, errCensusUnavailable
	}
	var config censusConfig
	seen := make(map[string]bool, 3)
	separator := -1
	for index := 1; index < len(args); {
		if args[index] == "--" {
			separator = index
			break
		}
		if index+1 >= len(args) || seen[args[index]] {
			return censusConfig{}, errCensusUnavailable
		}
		seen[args[index]] = true
		value := args[index+1]
		switch args[index] {
		case "--source-root":
			config.sourceRoot = value
		case "--receiver":
			config.receiver = value
		case "--selector":
			config.selector = value
		default:
			return censusConfig{}, errCensusUnavailable
		}
		index += 2
	}
	if separator < 0 || !seen["--source-root"] || !seen["--receiver"] || !seen["--selector"] || separator == len(args)-1 || !token.IsIdentifier(config.receiver) || !token.IsIdentifier(config.selector) {
		return censusConfig{}, errCensusUnavailable
	}
	if _, err := censusAbsoluteRootParts(config.sourceRoot); err != nil {
		return censusConfig{}, err
	}
	paths, err := validateCensusPaths(args[separator+1:])
	if err != nil {
		return censusConfig{}, err
	}
	config.paths = paths
	return config, nil
}

func validateCensusPaths(values []string) ([]string, error) {
	paths := make([]string, len(values))
	seen := make(map[string]bool, len(values))
	for index, value := range values {
		if _, err := censusRelativePathParts(value); err != nil || seen[value] {
			return nil, errCensusUnavailable
		}
		seen[value] = true
		paths[index] = value
	}
	return paths, nil
}

func censusRelativePathParts(value string) ([]string, error) {
	if value == "" || strings.ContainsRune(value, 0) || strings.HasPrefix(value, "/") || path.Clean(value) != value || !strings.HasSuffix(value, ".go") {
		return nil, errCensusUnavailable
	}
	parts := strings.Split(value, "/")
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return nil, errCensusUnavailable
		}
	}
	return parts, nil
}

func censusAbsoluteRootParts(value string) ([]string, error) {
	if value == "" || strings.ContainsRune(value, 0) || !strings.HasPrefix(value, "/") || path.Clean(value) != value {
		return nil, errCensusUnavailable
	}
	if value == "/" {
		return nil, nil
	}
	parts := strings.Split(value[1:], "/")
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return nil, errCensusUnavailable
		}
	}
	return parts, nil
}

func validCensusLimits(bounds census.Limits) bool {
	return bounds.MaxFiles != 0 && bounds.MaxFileBytes != 0 && bounds.MaxTotalBytes != 0 && bounds.MaxMatches != 0
}

func censusLimitError(err error) bool {
	var contract *changeevidence.ContractError
	return errors.As(err, &contract) && strings.HasPrefix(contract.Field, "limits.")
}

func censusOutput(inventory changeevidence.InventoryDocumentV1, result census.ResultV1, bounds census.Limits) censusDocumentV1 {
	entries := inventory.Entries()
	scope := make([]censusSourceRecordV1, len(entries))
	for index, entry := range entries {
		scope[index] = censusSourceRecordV1{
			PathB64:       base64.StdEncoding.EncodeToString(entry.Path),
			ContentSHA256: string(entry.ContentSHA256),
			ByteLength:    entry.ByteLength,
		}
	}
	matches := result.Matches()
	outputMatches := make([]censusMatchV1, len(matches))
	for index, match := range matches {
		outputMatches[index] = censusMatchV1{
			PathB64:        base64.StdEncoding.EncodeToString(match.Path),
			FileSHA256:     string(match.FileSHA256),
			FragmentSHA256: string(match.FragmentSHA256),
			Call:           censusSpanOutput(match.Call),
			Selector:       censusSpanOutput(match.Selector),
		}
	}
	query := result.Query()
	return censusDocumentV1{
		Schema:           censusCLISchemaV1,
		ExtractorVersion: string(result.ExtractorVersion()),
		Query:            censusQueryV1{Version: string(query.Version), Receiver: query.Receiver, Selector: query.Selector},
		Limits: censusLimitsV1{
			MaxFiles: bounds.MaxFiles, MaxFileBytes: bounds.MaxFileBytes, MaxTotalBytes: bounds.MaxTotalBytes, MaxMatches: bounds.MaxMatches,
		},
		Inventory: censusInventoryV1{
			Schema:                 string(changeevidence.InventoryContractV1),
			SHA256:                 inventory.Digest().String(),
			AccountingPolicySHA256: inventory.AccountingPolicyDigest().String(),
		},
		SourceScope: scope,
		Matches:     outputMatches,
	}
}

func censusSpanOutput(span census.Span) censusSpanV1 {
	return censusSpanV1{
		Start: censusPositionV1{Offset: span.Start.Offset, Line: span.Start.Line, Column: span.Start.Column},
		End:   censusPositionV1{Offset: span.End.Offset, Line: span.End.Line, Column: span.End.Column},
	}
}

func writeCensusDocument(stdout io.Writer, document censusDocumentV1) error {
	if stdout == nil {
		return errCensusUnavailable
	}
	var encoded bytes.Buffer
	encoder := json.NewEncoder(&encoded)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(document); err != nil {
		return err
	}
	written, err := stdout.Write(encoded.Bytes())
	if err != nil {
		return err
	}
	if written != encoded.Len() {
		return io.ErrShortWrite
	}
	return nil
}

func censusExit(stderr io.Writer, code int) int {
	if stderr != nil {
		fmt.Fprint(stderr, diagnostic(code))
	}
	return code
}
