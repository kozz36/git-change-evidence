// Package cnsicprofile admits only the observed W2 policy into neutral contracts.
package cnsicprofile

import (
	"strings"

	evidence "github.com/kozz36/git-change-evidence"
)

// DecodePolicyV1 admits an existing canonical neutral document, not a new
// consumer configuration format. Unsupported profile values are not inferred.
func DecodePolicyV1(raw []byte) (evidence.PolicyDocumentV1, error) {
	policy, err := evidence.DecodeAccountingPolicyV1(raw)
	if err != nil {
		return evidence.PolicyDocumentV1{}, err
	}
	view := policy.View()
	unsupported := func() (evidence.PolicyDocumentV1, error) {
		return evidence.PolicyDocumentV1{}, &evidence.ContractError{Field: "profile", Code: "unsupported"}
	}
	if view.Default != "production" || len(view.Ratios) != 0 || len(view.Categories) != 4 {
		return unsupported()
	}
	for i, name := range []string{"artifacts", "tests", "mechanical", "production"} {
		category := view.Categories[i]
		if category.Name != name {
			return unsupported()
		}
		if i == 3 {
			if category.LineThreshold != 400 || len(category.PathGlobs) != 0 {
				return unsupported()
			}
			continue
		}
		if category.LineThreshold != 0 || len(category.PathGlobs) != 1 {
			return unsupported()
		}
		glob := string(category.PathGlobs[0])
		if i == 0 && glob != "openspec/changes/**" || i == 1 && glob != "tests/**" || i == 2 && !literalAnnotationPath(glob) {
			return unsupported()
		}
	}
	return policy, nil
}

func literalAnnotationPath(path string) bool {
	if path == "" || strings.HasPrefix(path, "/") || strings.ContainsAny(path, "\x00*?[]\\") {
		return false
	}
	for _, part := range strings.Split(path, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}
