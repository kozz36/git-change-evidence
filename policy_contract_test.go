package changeevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func policyInput() PolicyInput {
	return PolicyInput{
		Categories: []CategoryInput{{Name: "later", PathGlobs: [][]byte{{0xff}, []byte("a"), []byte("a")}}, {Name: "first", PathGlobs: [][]byte{[]byte("src/**")}, LineThreshold: 7}},
		Default:    "later", Ratios: []RatioInput{{Category: "later", Reference: "first", Maximum: "0.10000000000000001"}, {Category: "first", Reference: "later", Maximum: "1"}},
	}
}
func newPolicy(t *testing.T, input PolicyInput) PolicyDocumentV1 {
	got, err := NewAccountingPolicyV1(input)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestAccountingPolicyV1Contract(t *testing.T) {
	input := policyInput()
	got := newPolicy(t, input)
	want := []byte("{\"schema\":\"git-change-evidence.accounting-policy/v1\",\"categories\":[{\"name\":\"later\",\"path_globs_b64\":[\"YQ==\",\"/w==\"],\"line_threshold\":0},{\"name\":\"first\",\"path_globs_b64\":[\"c3JjLyoq\"],\"line_threshold\":7}],\"default\":\"later\",\"ratios\":[{\"category\":\"later\",\"reference\":\"first\",\"maximum\":\"0.10000000000000001\"},{\"category\":\"first\",\"reference\":\"later\",\"maximum\":\"1\"}]}\n")
	sum := sha256.Sum256(want)
	if !bytes.Equal(got.CanonicalBytes(), want) || got.Digest().String() != hex.EncodeToString(sum[:]) {
		t.Fatalf("canonical = %q, digest = %q", got.CanonicalBytes(), got.Digest())
	}
	equivalent := policyInput()
	equivalent.Categories[0].PathGlobs = [][]byte{[]byte("a"), {0xff}}
	if !bytes.Equal(got.CanonicalBytes(), newPolicy(t, equivalent).CanonicalBytes()) || got.Digest() != newPolicy(t, equivalent).Digest() {
		t.Fatal("equivalent typed inputs differ")
	}
	other := policyInput()
	other.Ratios[0].Maximum = "0.1"
	if got.Digest() == newPolicy(t, other).Digest() {
		t.Fatal("exact decimal values collapsed")
	}
	input.Categories[0].PathGlobs[0][0] = 'x'
	view := got.View()
	view.Categories[0].PathGlobs[0][0], view.Ratios[0].Maximum = 'x', "1"
	got.CanonicalBytes()[0] = 'x'
	if got.View().Categories[0].PathGlobs[0][0] != 'a' || got.View().Ratios[0].Maximum != "0.10000000000000001" || got.CanonicalBytes()[0] != '{' {
		t.Fatal("mutable state escaped")
	}
	empty := policyInput()
	empty.Categories[0].PathGlobs, empty.Ratios = nil, nil
	if bytes.Contains(newPolicy(t, empty).CanonicalBytes(), []byte(":null")) {
		t.Fatal("canonical arrays must not be null")
	}
	maxThreshold := policyInput()
	maxThreshold.Categories[0].LineThreshold = ^uint64(0)
	if _, err := DecodeAccountingPolicyV1(newPolicy(t, maxThreshold).CanonicalBytes()); err != nil {
		t.Fatal(err)
	}
	var zero PolicyDocumentV1
	if zero.CanonicalBytes() != nil || zero.Digest() != "" || zero.View().Categories != nil {
		t.Fatal("zero value presents a document")
	}
	valid := policyInput()
	valid.Categories[0].PathGlobs = [][]byte{[]byte("src/[a-z]?.go")}
	newPolicy(t, valid)
	tiny := policyInput()
	tiny.Ratios[0].Maximum = "0.00000000000000001"
	newPolicy(t, tiny)
	for name, mutate := range map[string]func(*PolicyInput){
		"empty categories": func(p *PolicyInput) { p.Categories = nil }, "empty name": func(p *PolicyInput) { p.Categories[0].Name = "" }, "control name": func(p *PolicyInput) { p.Categories[0].Name = "\x01" },
		"duplicate name": func(p *PolicyInput) { p.Categories = append(p.Categories, p.Categories[0]) }, "absolute glob": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{[]byte("/x")} }, "empty component": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{[]byte("a//b")} },
		"dot glob": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{[]byte("./a")} }, "dotdot glob": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{[]byte("a/../b")} }, "backslash glob": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{[]byte(`a\\b`)} },
		"NUL glob": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{[]byte("a\x00b")} }, "empty glob": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{{}} }, "empty class": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{[]byte("[]")} }, "malformed class": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{[]byte("[")} },
		"cross separator class": func(p *PolicyInput) { p.Categories[0].PathGlobs = [][]byte{[]byte("a/[a/b]")} }, "missing default": func(p *PolicyInput) { p.Default = "missing" }, "missing ratio category": func(p *PolicyInput) { p.Ratios[0].Category = "missing" },
		"missing ratio reference": func(p *PolicyInput) { p.Ratios[0].Reference = "missing" }, "self ratio": func(p *PolicyInput) { p.Ratios[0].Reference = "later" }, "duplicate pair": func(p *PolicyInput) { p.Ratios = append(p.Ratios, p.Ratios[0]) },
		"maximum above one": func(p *PolicyInput) { p.Ratios[0].Maximum = "1.0000000000000001" }, "maximum signed": func(p *PolicyInput) { p.Ratios[0].Maximum = "+1" }, "maximum negative": func(p *PolicyInput) { p.Ratios[0].Maximum = "-1" },
		"maximum zero": func(p *PolicyInput) { p.Ratios[0].Maximum = "0" }, "maximum noncanonical": func(p *PolicyInput) { p.Ratios[0].Maximum = "1.0" },
	} {
		t.Run(name, func(t *testing.T) {
			candidate := policyInput()
			mutate(&candidate)
			if _, err := NewAccountingPolicyV1(candidate); err == nil {
				t.Fatal("invalid policy accepted")
			}
		})
	}
}
