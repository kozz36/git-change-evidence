package changeevidence

import "math"

type AccountingCategory struct {
	Name          string
	PathGlobs     []string
	LineThreshold int64
}
type CategoryRatio struct {
	Category, Reference string
	Maximum             float64
}
type AccountingPolicy struct {
	Categories []AccountingCategory
	Default    string
	Ratios     []CategoryRatio
}
type CategoryTotal struct {
	Category             string
	Additions, Deletions uint64
	NonCountable         uint64
}
type AccountingObservation struct {
	Category, Reference string
	Limit, Actual       float64
	Available, Exceeded bool
}
type AccountingResult struct {
	Totals       []CategoryTotal
	Observations []AccountingObservation
}
type AccountingError struct{ Field, Code string }

const (
	AccountingInvalidPolicy = "invalid_policy"
	AccountingOverflow      = "overflow"
)

func (e *AccountingError) Error() string { return "accounting " + e.Field + ": " + e.Code }

func Account(snapshot CommittedSnapshot, policy AccountingPolicy) (AccountingResult, error) {
	index, err := validateAccountingPolicy(policy)
	if err != nil {
		return AccountingResult{}, err
	}
	result := AccountingResult{Totals: make([]CategoryTotal, len(policy.Categories))}
	for i, category := range policy.Categories {
		result.Totals[i].Category = category.Name
	}
	for _, change := range snapshot.Entries() {
		total := &result.Totals[matchingCategory(policy, change.Path, index)]
		if !change.Lines.Countable {
			total.NonCountable++
			continue
		}
		var ok bool
		if total.Additions, ok = addLines(total.Additions, change.Lines.Additions); !ok {
			return AccountingResult{}, accountingFailure("totals", AccountingOverflow)
		}
		if total.Deletions, ok = addLines(total.Deletions, change.Lines.Deletions); !ok {
			return AccountingResult{}, accountingFailure("totals", AccountingOverflow)
		}
	}
	for i, category := range policy.Categories {
		if category.LineThreshold == 0 {
			continue
		}
		observation := AccountingObservation{Category: category.Name, Limit: float64(category.LineThreshold)}
		if total, ok := lineTotal(result.Totals[i]); ok {
			observation.Available, observation.Actual = true, float64(total)
			observation.Exceeded = total > uint64(category.LineThreshold)
		}
		result.Observations = append(result.Observations, observation)
	}
	for _, ratio := range policy.Ratios {
		observation := AccountingObservation{Category: ratio.Category, Reference: ratio.Reference, Limit: ratio.Maximum}
		numerator, numeratorOK := lineTotal(result.Totals[index[ratio.Category]])
		denominator, denominatorOK := lineTotal(result.Totals[index[ratio.Reference]])
		if numeratorOK && denominatorOK && denominator != 0 {
			observation.Available, observation.Actual = true, float64(numerator)/float64(denominator)
			observation.Exceeded = observation.Actual > ratio.Maximum
		}
		result.Observations = append(result.Observations, observation)
	}
	return result, nil
}

func validateAccountingPolicy(policy AccountingPolicy) (map[string]int, error) {
	if len(policy.Categories) == 0 {
		return nil, accountingFailure("categories", AccountingInvalidPolicy)
	}
	index := make(map[string]int, len(policy.Categories))
	for i, category := range policy.Categories {
		if category.Name == "" || category.LineThreshold < 0 {
			return nil, accountingFailure("categories", AccountingInvalidPolicy)
		}
		if _, exists := index[category.Name]; exists {
			return nil, accountingFailure("categories", AccountingInvalidPolicy)
		}
		for _, glob := range category.PathGlobs {
			if !validPathGlob(glob) {
				return nil, accountingFailure("patterns", AccountingInvalidPolicy)
			}
		}
		index[category.Name] = i
	}
	if _, exists := index[policy.Default]; policy.Default == "" || !exists {
		return nil, accountingFailure("default", AccountingInvalidPolicy)
	}
	seen := map[string]bool{}
	for _, ratio := range policy.Ratios {
		key := ratio.Category + "\x00" + ratio.Reference
		_, reference := index[ratio.Reference]
		if _, category := index[ratio.Category]; !category || ratio.Category == ratio.Reference || !reference || seen[key] || math.IsNaN(ratio.Maximum) || math.IsInf(ratio.Maximum, 0) || ratio.Maximum <= 0 || ratio.Maximum > 1 {
			return nil, accountingFailure("ratios", AccountingInvalidPolicy)
		}
		seen[key] = true
	}
	return index, nil
}

func accountingFailure(field, code string) error { return &AccountingError{Field: field, Code: code} }
func addLines(left, right uint64) (uint64, bool) {
	if ^uint64(0)-left < right {
		return 0, false
	}
	return left + right, true
}
func lineTotal(total CategoryTotal) (uint64, bool) {
	if total.NonCountable != 0 {
		return 0, false
	}
	return addLines(total.Additions, total.Deletions)
}
func matchingCategory(policy AccountingPolicy, path string, index map[string]int) int {
	for i, category := range policy.Categories {
		for _, glob := range category.PathGlobs {
			if matchPathGlob([]byte(glob), []byte(path)) {
				return i
			}
		}
	}
	return index[policy.Default]
}

func validPathGlob(glob string) bool {
	bytes := []byte(glob)
	if len(bytes) == 0 {
		return false
	}
	for i := 0; i < len(bytes); i++ {
		if bytes[i] != '[' {
			continue
		}
		end := i + 1
		if end < len(bytes) && (bytes[end] == '!' || bytes[end] == '^') {
			end++
		}
		start := end
		for end < len(bytes) && bytes[end] != ']' {
			if bytes[end] == '/' {
				return false
			}
			end++
		}
		if start == end || end == len(bytes) {
			return false
		}
		i = end
	}
	return true
}

func matchPathGlob(pattern, path []byte) bool {
	for len(pattern) != 0 {
		switch pattern[0] {
		case '*':
			globstar, skip := len(pattern) > 1 && pattern[1] == '*', 1
			if globstar {
				skip++
				if len(pattern) > skip && pattern[skip] == '/' && matchPathGlob(pattern[skip+1:], path) {
					return true
				}
			}
			for i := 0; ; i++ {
				if matchPathGlob(pattern[skip:], path[i:]) {
					return true
				}
				if i == len(path) || (!globstar && path[i] == '/') {
					return false
				}
			}
		case '?':
			if len(path) == 0 || path[0] == '/' {
				return false
			}
			pattern, path = pattern[1:], path[1:]
		case '[':
			end := 1
			for end < len(pattern) && pattern[end] != ']' {
				end++
			}
			if len(path) == 0 || path[0] == '/' || end == len(pattern) || !matchesClass(pattern[1:end], path[0]) {
				return false
			}
			pattern, path = pattern[end+1:], path[1:]
		default:
			if len(path) == 0 || pattern[0] != path[0] {
				return false
			}
			pattern, path = pattern[1:], path[1:]
		}
	}
	return len(path) == 0
}

func matchesClass(class []byte, value byte) bool {
	negated, match := len(class) > 0 && (class[0] == '!' || class[0] == '^'), false
	if negated {
		class = class[1:]
	}
	for i := 0; i < len(class); i++ {
		if i+2 < len(class) && class[i+1] == '-' {
			match = match || class[i] <= value && value <= class[i+2]
			i += 2
		} else {
			match = match || class[i] == value
		}
	}
	return match != negated
}
