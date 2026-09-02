package changeevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const AccountingPolicyContractV1 ContractVersion = "git-change-evidence.accounting-policy/v1"

type DocumentDigest string

func (d DocumentDigest) String() string { return string(d) }

type PolicyInput struct {
	Categories []CategoryInput
	Default    string
	Ratios     []RatioInput
}
type CategoryInput struct {
	Name          string
	PathGlobs     [][]byte
	LineThreshold uint64
}
type RatioInput struct{ Category, Reference, Maximum string }
type PolicyView PolicyInput

type PolicyDocumentV1 struct {
	canonical []byte
	digest    DocumentDigest
	view      PolicyInput
}
type policyWire struct {
	Schema     ContractVersion  `json:"schema"`
	Categories []policyCategory `json:"categories"`
	Default    string           `json:"default"`
	Ratios     []policyRatio    `json:"ratios"`
}
type policyCategory struct {
	Name          string   `json:"name"`
	PathGlobs     []string `json:"path_globs_b64"`
	LineThreshold uint64   `json:"line_threshold"`
}
type policyRatio struct {
	Category  string `json:"category"`
	Reference string `json:"reference"`
	Maximum   string `json:"maximum"`
}

func NewAccountingPolicyV1(input PolicyInput) (PolicyDocumentV1, error) {
	return makePolicy(input, true)
}
func (p PolicyDocumentV1) CanonicalBytes() []byte { return append([]byte(nil), p.canonical...) }
func (p PolicyDocumentV1) Digest() DocumentDigest { return p.digest }
func (p PolicyDocumentV1) View() PolicyView {
	if p.canonical == nil {
		return PolicyView{}
	}
	return PolicyView(clonePolicyInput(p.view))
}

func makePolicy(input PolicyInput, normalize bool) (PolicyDocumentV1, error) {
	input = clonePolicyInput(input)
	if len(input.Categories) == 0 {
		return PolicyDocumentV1{}, policyError("categories", "required")
	}
	names := make(map[string]bool, len(input.Categories))
	for i := range input.Categories {
		category := &input.Categories[i]
		if !validPolicyText(category.Name) || names[category.Name] {
			return PolicyDocumentV1{}, policyError("categories", "invalid")
		}
		names[category.Name] = true
		if normalize {
			sort.Slice(category.PathGlobs, func(i, j int) bool { return bytes.Compare(category.PathGlobs[i], category.PathGlobs[j]) < 0 })
		}
		for j, glob := range category.PathGlobs {
			if !validPolicyGlob(string(glob)) || j > 0 && bytes.Compare(category.PathGlobs[j-1], glob) >= 0 && !normalize {
				return PolicyDocumentV1{}, policyError("categories", "invalid_globs")
			}
		}
		if normalize {
			category.PathGlobs = slices.CompactFunc(category.PathGlobs, bytes.Equal)
		}
	}
	if !names[input.Default] {
		return PolicyDocumentV1{}, policyError("default", "invalid")
	}
	pairs := make(map[string]bool, len(input.Ratios))
	for _, ratio := range input.Ratios {
		_, err := parsePolicyDecimal(ratio.Maximum)
		maximum, exact := new(big.Rat).SetString(ratio.Maximum)
		key := ratio.Category + "\x00" + ratio.Reference
		if err != nil || !exact || !names[ratio.Category] || !names[ratio.Reference] || ratio.Category == ratio.Reference || pairs[key] || maximum.Sign() <= 0 || maximum.Cmp(big.NewRat(1, 1)) > 0 {
			return PolicyDocumentV1{}, policyError("ratios", "invalid")
		}
		pairs[key] = true
	}
	wire := policyWire{AccountingPolicyContractV1, make([]policyCategory, len(input.Categories)), input.Default, make([]policyRatio, len(input.Ratios))}
	for i, category := range input.Categories {
		globs := make([]string, len(category.PathGlobs))
		for j, glob := range category.PathGlobs {
			globs[j] = base64.StdEncoding.EncodeToString(glob)
		}
		wire.Categories[i] = policyCategory{category.Name, globs, category.LineThreshold}
	}
	for i, ratio := range input.Ratios {
		wire.Ratios[i] = policyRatio{ratio.Category, ratio.Reference, ratio.Maximum}
	}
	canonical, err := json.Marshal(wire)
	if err != nil {
		return PolicyDocumentV1{}, err
	}
	canonical = append(canonical, '\n')
	sum := sha256.Sum256(canonical)
	return PolicyDocumentV1{append([]byte(nil), canonical...), DocumentDigest(hex.EncodeToString(sum[:])), clonePolicyInput(input)}, nil
}

func clonePolicyInput(input PolicyInput) PolicyInput {
	result := PolicyInput{Categories: make([]CategoryInput, len(input.Categories)), Default: input.Default, Ratios: append([]RatioInput(nil), input.Ratios...)}
	for i, category := range input.Categories {
		result.Categories[i] = CategoryInput{Name: category.Name, PathGlobs: make([][]byte, len(category.PathGlobs)), LineThreshold: category.LineThreshold}
		for j, glob := range category.PathGlobs {
			result.Categories[i].PathGlobs[j] = append([]byte(nil), glob...)
		}
	}
	return result
}
func validPolicyText(value string) bool {
	return value != "" && utf8.ValidString(value) && strings.IndexFunc(value, unicode.IsControl) < 0
}
func policyError(field, code string) error { return &ContractError{field, code} }
