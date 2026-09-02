package changeevidence

import "testing"

func TestPolicyDecimal(t *testing.T) {
	tests := []struct {
		name, input, want string
		valid             bool
	}{
		{"zero", "0", "0", true},
		{"fraction", "0.5", "0.5", true},
		{"whole and fraction", "12.34", "12.34", true},
		{"leading zero", "01", "", false},
		{"trailing fractional zero", "1.20", "", false},
		{"exponent", "1e2", "", false},
		{"signed", "+1", "", false},
		{"missing whole", ".5", "", false},
		{"missing fraction", "1.", "", false},
		{"zero fraction", "0.0", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePolicyDecimal(tt.input)
			if (err == nil) != tt.valid {
				t.Fatalf("parsePolicyDecimal(%q) error = %v, valid = %v", tt.input, err, tt.valid)
			}
			if err == nil && got.string() != tt.want {
				t.Fatalf("parsePolicyDecimal(%q).string() = %q, want %q", tt.input, got.string(), tt.want)
			}
		})
	}
}

func TestPolicyDecimalCompare(t *testing.T) {
	tests := []struct {
		name, left, right string
		want              int
	}{
		{"fractional precision", "0.10000000000000001", "0.1", 1},
		{"different scales", "0.9", "0.01", 1},
		{"equal wholes", "12", "12", 0},
		{"large whole exactness", "9007199254740993", "9007199254740992", 1},
		{"fractional ordering", "0.001", "0.01", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, err := parsePolicyDecimal(tt.left)
			if err != nil {
				t.Fatal(err)
			}
			right, err := parsePolicyDecimal(tt.right)
			if err != nil {
				t.Fatal(err)
			}
			if got := comparePolicyDecimals(left, right); got != tt.want {
				t.Fatalf("comparePolicyDecimals(%q, %q) = %d, want %d", tt.left, tt.right, got, tt.want)
			}
		})
	}
}
