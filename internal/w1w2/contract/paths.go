package contract

import "bytes"

// ValidatePath checks raw Linux names; it grants no filesystem or confidentiality authority.
func ValidatePath(path []byte) error {
	if len(path) > MaxPathBytes {
		return &Failure{Status: 4, Token: "resource_limit"}
	}
	if len(path) == 0 || path[0] == '/' || bytes.IndexByte(path, 0) >= 0 {
		return &Failure{Status: 2, Token: "unsafe_path"}
	}
	for _, component := range bytes.Split(path, []byte{'/'}) {
		if len(component) > MaxComponentBytes {
			return &Failure{Status: 4, Token: "resource_limit"}
		}
		if len(component) == 0 || bytes.Equal(component, []byte(".")) || bytes.Equal(component, []byte("..")) || gitComponent(component) {
			return &Failure{Status: 2, Token: "unsafe_path"}
		}
	}
	return nil
}

// ValidatePaths checks raw identities and charges every path without reordering or normalization.
func ValidatePaths(paths [][]byte) error {
	if len(paths) > MaxPaths {
		return &Failure{Status: 4, Token: "resource_limit"}
	}
	remaining := MaxTotalPathBytes
	seen := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		if err := ValidatePath(path); err != nil {
			return err
		}
		if len(path) > remaining {
			return &Failure{Status: 4, Token: "resource_limit"}
		}
		key := string(path)
		if _, duplicate := seen[key]; duplicate {
			return &Failure{Status: 2, Token: "invalid_document"}
		}
		seen[key] = struct{}{}
		remaining -= len(path)
	}
	return nil
}

// Metadata exclusion folds ASCII only; it is not a physical hardlink/alias certificate.
func gitComponent(component []byte) bool {
	if len(component) != 4 || component[0] != '.' {
		return false
	}
	for i, want := range []byte("git") {
		value := component[i+1]
		if value >= 'A' && value <= 'Z' {
			value += 'a' - 'A'
		}
		if value != want {
			return false
		}
	}
	return true
}
