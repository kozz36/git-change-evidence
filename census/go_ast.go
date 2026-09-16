package census

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"

	ce "github.com/kozz36/git-change-evidence"
)

type sourceFile struct {
	path    []byte
	content []byte
	digest  ce.Digest
}

// GoASTV1 returns every syntax-only textual selector-call candidate in the
// complete, validated supplied inventory scope. It performs no semantic lookup.
func GoASTV1(inventory ce.InventoryDocumentV1, files []File, query QueryV1, limits Limits) (ResultV1, error) {
	if err := validateQuery(query); err != nil {
		return ResultV1{}, err
	}
	if err := validateLimits(limits); err != nil {
		return ResultV1{}, err
	}
	entries, err := validateInventory(inventory)
	if err != nil {
		return ResultV1{}, err
	}
	sources, err := validateSources(entries, files, limits)
	if err != nil {
		return ResultV1{}, err
	}

	matches := make([]MatchV1, 0)
	for _, source := range sources {
		set := token.NewFileSet()
		parsed, parseErr := parser.ParseFile(set, "supplied.go", source.content, parser.SkipObjectResolution)
		if parseErr != nil || parsed == nil {
			return ResultV1{}, contract("source", "parse_error")
		}
		file := set.File(parsed.Pos())
		if file == nil || file.Size() != len(source.content) {
			return ResultV1{}, contract("span", "invalid_span")
		}
		var walkErr error
		ast.Inspect(parsed, func(node ast.Node) bool {
			if walkErr != nil {
				return false
			}
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector := normalizedSelector(call.Fun)
			if selector == nil || !matchesQuery(selector, query) {
				return true
			}
			callSpan, err := physicalSpan(set, file, call.Pos(), call.End(), source.content)
			if err != nil {
				walkErr = err
				return false
			}
			selectorSpan, err := physicalSpan(set, file, selector.Pos(), selector.End(), source.content)
			if err != nil || selectorSpan.Start.Offset < callSpan.Start.Offset || selectorSpan.End.Offset > callSpan.End.Offset {
				walkErr = contract("span", "invalid_span")
				return false
			}
			if uint64(len(matches)) == limits.MaxMatches {
				walkErr = contract("limits.max_matches", "limit_exceeded")
				return false
			}
			matches = append(matches, MatchV1{
				Path:           append([]byte(nil), source.path...),
				FileSHA256:     source.digest,
				FragmentSHA256: digest(source.content[callSpan.Start.Offset:callSpan.End.Offset]),
				Call:           callSpan,
				Selector:       selectorSpan,
			})
			return true
		})
		if walkErr != nil {
			return ResultV1{}, walkErr
		}
	}
	sortMatches(matches)
	return ResultV1{inventory: inventory.Digest(), query: query, extractor: ExtractorVersion, matches: matches}, nil
}

func validateQuery(query QueryV1) error {
	if query.Version != SelectorCallQueryV1 {
		return contract("query.version", "unsupported_version")
	}
	if !token.IsIdentifier(query.Receiver) {
		return contract("query.receiver", "invalid_identifier")
	}
	if !token.IsIdentifier(query.Selector) {
		return contract("query.selector", "invalid_identifier")
	}
	return nil
}

func validateLimits(limits Limits) error {
	for _, limit := range []struct {
		field string
		value uint64
	}{
		{"limits.max_files", limits.MaxFiles},
		{"limits.max_file_bytes", limits.MaxFileBytes},
		{"limits.max_total_bytes", limits.MaxTotalBytes},
		{"limits.max_matches", limits.MaxMatches},
	} {
		if limit.value == 0 {
			return contract(limit.field, "invalid_limit")
		}
	}
	return nil
}

func validateInventory(inventory ce.InventoryDocumentV1) ([]ce.InventoryEntryV1, error) {
	canonical := inventory.CanonicalBytes()
	entries := inventory.Entries()
	if len(canonical) == 0 || inventory.AccountingPolicyDigest() == "" || entries == nil || inventory.Digest() != ce.DocumentDigest(digest(canonical)) {
		return nil, contract("inventory", "invalid_document")
	}
	return entries, nil
}

func validateSources(entries []ce.InventoryEntryV1, files []File, limits Limits) ([]sourceFile, error) {
	if uint64(len(entries)) > limits.MaxFiles || uint64(len(files)) > limits.MaxFiles {
		return nil, contract("limits.max_files", "limit_exceeded")
	}
	sources := make([]sourceFile, len(files))
	for index, supplied := range files {
		sources[index] = sourceFile{path: append([]byte(nil), supplied.Path...), content: append([]byte(nil), supplied.Content...)}
	}
	sort.Slice(sources, func(left, right int) bool { return bytes.Compare(sources[left].path, sources[right].path) < 0 })
	for index := 1; index < len(sources); index++ {
		if bytes.Equal(sources[index-1].path, sources[index].path) {
			return nil, contract("files.path", "duplicate_path")
		}
	}
	for index := 0; index < len(entries) && index < len(sources); index++ {
		comparison := bytes.Compare(entries[index].Path, sources[index].path)
		if comparison < 0 {
			return nil, contract("files.path", "missing_path")
		}
		if comparison > 0 {
			return nil, contract("files.path", "extra_path")
		}
	}
	if len(entries) > len(sources) {
		return nil, contract("files.path", "missing_path")
	}
	if len(entries) < len(sources) {
		return nil, contract("files.path", "extra_path")
	}

	var total uint64
	for _, source := range sources {
		size := uint64(len(source.content))
		if size > limits.MaxFileBytes {
			return nil, contract("limits.max_file_bytes", "limit_exceeded")
		}
		if !withinTotal(limits.MaxTotalBytes, total, size) {
			return nil, contract("limits.max_total_bytes", "limit_exceeded")
		}
		total += size
	}
	for index := range sources {
		if uint64(len(sources[index].content)) != entries[index].ByteLength {
			return nil, contract("files.byte_length", "length_mismatch")
		}
		sources[index].digest = digest(sources[index].content)
		if sources[index].digest != entries[index].ContentSHA256 {
			return nil, contract("files.content_sha256", "digest_mismatch")
		}
	}
	return sources, nil
}

func withinTotal(limit, total, next uint64) bool {
	return total <= limit && next <= limit-total
}

func normalizedSelector(expression ast.Expr) *ast.SelectorExpr {
	for {
		switch node := expression.(type) {
		case *ast.ParenExpr:
			expression = node.X
		case *ast.IndexExpr:
			expression = node.X
		case *ast.IndexListExpr:
			expression = node.X
		default:
			selector, _ := expression.(*ast.SelectorExpr)
			return selector
		}
	}
}

func matchesQuery(selector *ast.SelectorExpr, query QueryV1) bool {
	receiver, ok := selector.X.(*ast.Ident)
	return ok && receiver.Name == query.Receiver && selector.Sel.Name == query.Selector
}

func physicalSpan(set *token.FileSet, file *token.File, start, end token.Pos, content []byte) (Span, error) {
	startPosition, err := physicalPosition(set, file, start, content)
	if err != nil {
		return Span{}, err
	}
	endPosition, err := physicalPosition(set, file, end, content)
	if err != nil || startPosition.Offset >= endPosition.Offset {
		return Span{}, contract("span", "invalid_span")
	}
	return Span{Start: startPosition, End: endPosition}, nil
}

func physicalPosition(set *token.FileSet, file *token.File, position token.Pos, content []byte) (Position, error) {
	if position == token.NoPos || set.File(position) != file || file.Size() != len(content) || int(position) < file.Base() || int(position) > file.Base()+file.Size() {
		return Position{}, contract("span", "invalid_span")
	}
	offset := file.Offset(position)
	physical := file.PositionFor(position, false)
	if offset < 0 || offset > len(content) || physical.Offset != offset || physical.Line < 1 || physical.Column < 1 {
		return Position{}, contract("span", "invalid_span")
	}
	return Position{Offset: offset, Line: physical.Line, Column: physical.Column}, nil
}

func sortMatches(matches []MatchV1) {
	sort.Slice(matches, func(left, right int) bool {
		if comparison := bytes.Compare(matches[left].Path, matches[right].Path); comparison != 0 {
			return comparison < 0
		}
		leftKey := []int{matches[left].Call.Start.Offset, matches[left].Call.End.Offset, matches[left].Selector.Start.Offset, matches[left].Selector.End.Offset}
		rightKey := []int{matches[right].Call.Start.Offset, matches[right].Call.End.Offset, matches[right].Selector.Start.Offset, matches[right].Selector.End.Offset}
		for index := range leftKey {
			if leftKey[index] != rightKey[index] {
				return leftKey[index] < rightKey[index]
			}
		}
		return false
	})
}

func digest(content []byte) ce.Digest {
	sum := sha256.Sum256(content)
	return ce.Digest(hex.EncodeToString(sum[:]))
}

func contract(field, code string) error { return &ce.ContractError{Field: field, Code: code} }
