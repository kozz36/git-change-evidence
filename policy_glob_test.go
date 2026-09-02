package changeevidence

import (
	"bytes"
	"testing"
)

func TestPolicyGlob(t *testing.T) {
	tests := []struct {
		name, glob string
		valid      bool
	}{
		{"single segment wildcard", "*.go", true},
		{"normalized recursive path", "docs/**/*.md", true},
		{"normalized class", "src/[a-z]?.go", true},
		{"empty", "", false},
		{"absolute", "/etc/*", false},
		{"current directory", "./docs/*.md", false},
		{"parent directory", "docs/../*.md", false},
		{"invalid class", "src/[a-z", false},
		{"trailing separator", "docs/", false},
		{"duplicate separator", "docs//*.md", false},
		{"leading parent", "../docs/*.md", false},
		{"backslash separator", `docs\\*.md`, false},
		{"NUL", "docs/\x00*.md", false},
		{"unmatched closing class", "src/].go", false},
		{"descending class range", "src/[z-a].go", false},
		{"ascending class range", "src/[a-z].go", true},
		{"bang-negated class range", "src/[!a-z].go", true},
		{"caret-negated class range", "src/[^0-9].go", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validPolicyGlob(tt.glob); got != tt.valid {
				t.Fatalf("validPolicyGlob(%q) = %v, want %v", tt.glob, got, tt.valid)
			}
		})
	}
}

func TestPolicyGlobByteIdentity(t *testing.T) {
	const composed = "caf\u00e9/*.go"
	const decomposed = "cafe\u0301/*.go"
	tests := []struct {
		name, glob string
	}{
		{"composed", composed},
		{"decomposed", decomposed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !validPolicyGlob(tt.glob) {
				t.Fatalf("validPolicyGlob(%q) = false", tt.glob)
			}
		})
	}
	if bytes.Equal([]byte(composed), []byte(decomposed)) {
		t.Fatal("Unicode-equivalent glob spellings have equal bytes")
	}
	if matchPathGlob([]byte(composed), []byte(decomposed)) {
		t.Fatal("Unicode-equivalent glob spellings matched after validation")
	}
}
