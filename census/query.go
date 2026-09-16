// Package census produces exact syntax-only selector-call candidates for a validated inventory scope.
package census

import ce "github.com/kozz36/git-change-evidence"

const (
	// SelectorCallQueryV1 identifies the exact textual selector predicate.
	SelectorCallQueryV1 ce.ContractVersion = "git-change-evidence.selector-call-query/v1"
	// ExtractorVersion identifies this Go AST extraction behavior.
	ExtractorVersion ce.ContractVersion = "git-change-evidence.go-ast-census/v1"
)

// File supplies the complete bytes for one raw inventory path.
type File struct {
	Path    []byte
	Content []byte
}

// QueryV1 names one exact, versioned textual selector predicate.
type QueryV1 struct {
	Version            ce.ContractVersion
	Receiver, Selector string
}

// Limits bound admitted input and published matches, not parser resources.
type Limits struct {
	MaxFiles, MaxFileBytes, MaxTotalBytes, MaxMatches uint64
}

// DefaultLimits returns the exact bounds for one synchronous census operation.
func DefaultLimits() Limits {
	return Limits{MaxFiles: 256, MaxFileBytes: 1 << 20, MaxTotalBytes: 8 << 20, MaxMatches: 10000}
}

// Position identifies a physical source byte and its one-based line and byte column.
type Position struct {
	Offset, Line, Column int
}

// Span is a half-open physical source range.
type Span struct {
	Start, End Position
}

// MatchV1 binds one candidate to its path, complete file, full call, and selector evidence.
type MatchV1 struct {
	Path                       []byte
	FileSHA256, FragmentSHA256 ce.Digest
	Call, Selector             Span
}

// ResultV1 retains the bindings and ordered evidence from one successful census.
type ResultV1 struct {
	inventory ce.DocumentDigest
	query     QueryV1
	extractor ce.ContractVersion
	matches   []MatchV1
}

// InventoryDigest returns the supplied inventory's exact canonical digest.
func (r ResultV1) InventoryDigest() ce.DocumentDigest { return r.inventory }

// Query returns the verbatim versioned query.
func (r ResultV1) Query() QueryV1 { return r.query }

// ExtractorVersion returns the implementation-owned extraction version.
func (r ResultV1) ExtractorVersion() ce.ContractVersion { return r.extractor }

// Matches returns a defensive copy of the ordered candidate evidence.
func (r ResultV1) Matches() []MatchV1 {
	if r.matches == nil {
		return nil
	}
	matches := make([]MatchV1, len(r.matches))
	for index, match := range r.matches {
		matches[index] = match
		matches[index].Path = append([]byte(nil), match.Path...)
	}
	return matches
}
