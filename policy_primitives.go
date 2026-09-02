package changeevidence

import (
	"bytes"
	"encoding/json"
	"io"
	"path"
	"strings"
)

type policyDecimal struct {
	digits string
	scale  int
}

func parsePolicyDecimal(value string) (policyDecimal, error) {
	whole, fraction, fractional := strings.Cut(value, ".")
	if whole == "" || (len(whole) > 1 && whole[0] == '0') || !policyDigits(whole) ||
		fractional && (fraction == "" || fraction[len(fraction)-1] == '0' || !policyDigits(fraction)) {
		return policyDecimal{}, &ContractError{"policy.decimal", "noncanonical"}
	}
	return policyDecimal{digits: whole + fraction, scale: len(fraction)}, nil
}

func (value policyDecimal) string() string {
	if value.scale == 0 {
		return value.digits
	}
	point := len(value.digits) - value.scale
	return value.digits[:point] + "." + value.digits[point:]
}

func comparePolicyDecimals(left, right policyDecimal) int {
	leftWhole, rightWhole := len(left.digits)-left.scale, len(right.digits)-right.scale
	if leftWhole < rightWhole {
		return -1
	}
	if leftWhole > rightWhole {
		return 1
	}
	for index := 0; index < leftWhole; index++ {
		if left.digits[index] < right.digits[index] {
			return -1
		}
		if left.digits[index] > right.digits[index] {
			return 1
		}
	}
	for index := 0; index < max(left.scale, right.scale); index++ {
		leftDigit, rightDigit := byte('0'), byte('0')
		if index < left.scale {
			leftDigit = left.digits[leftWhole+index]
		}
		if index < right.scale {
			rightDigit = right.digits[rightWhole+index]
		}
		if leftDigit < rightDigit {
			return -1
		}
		if leftDigit > rightDigit {
			return 1
		}
	}
	return 0
}

func policyDigits(value string) bool {
	for index := range value {
		if value[index] < '0' || value[index] > '9' {
			return false
		}
	}
	return true
}

func validPolicyGlob(value string) bool {
	return value != "" && !strings.ContainsAny(value, "\\\x00") && !path.IsAbs(value) && value != "." && value != ".." &&
		!strings.HasPrefix(value, "../") && path.Clean(value) == value && validPolicyGlobClasses(value)
}

func validPolicyGlobClasses(value string) bool {
	for start := 0; start < len(value); start++ {
		if value[start] == ']' {
			return false
		}
		if value[start] != '[' {
			continue
		}
		end := start + 1
		if end < len(value) && (value[end] == '!' || value[end] == '^') {
			end++
		}
		first := end
		for end < len(value) && value[end] != ']' {
			if value[end] == '/' {
				return false
			}
			end++
		}
		if first == end || end == len(value) {
			return false
		}
		for index := first; index+2 < end; index++ {
			if value[index+1] == '-' {
				if value[index] > value[index+2] {
					return false
				}
				index += 2
			}
		}
		start = end
	}
	return true
}

func decodePolicyObject(value []byte, field string, allowed ...string) (map[string]json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(value))
	start, err := decoder.Token()
	if err != nil || start != json.Delim('{') {
		return nil, invalidPolicyObject(field)
	}
	allowedFields, object := make(map[string]bool, len(allowed)), make(map[string]json.RawMessage, len(allowed))
	for _, name := range allowed {
		allowedFields[name] = true
	}
	for decoder.More() {
		token, err := decoder.Token()
		name, ok := token.(string)
		if err != nil || !ok {
			return nil, invalidPolicyObject(field)
		}
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			return nil, invalidPolicyObject(field)
		}
		if _, duplicate := object[name]; duplicate {
			return nil, &ContractError{field + "." + name, "duplicate_field"}
		}
		if !allowedFields[name] {
			code := "unknown_field"
			if isAuthorityField(name) {
				code = "authority_field"
			}
			return nil, &ContractError{field + "." + name, code}
		}
		object[name] = raw
	}
	end, err := decoder.Token()
	if err != nil || end != json.Delim('}') {
		return nil, invalidPolicyObject(field)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return nil, invalidPolicyObject(field)
	}
	return object, nil
}

func invalidPolicyObject(field string) error { return &ContractError{field, "invalid_object"} }
