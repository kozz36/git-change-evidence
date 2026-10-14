package cnsicprofile_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	"github.com/kozz36/git-change-evidence/internal/cnsicprofile"
)

func control() evidence.PolicyInput {
	return evidence.PolicyInput{Categories: []evidence.CategoryInput{
		{Name: "artifacts", PathGlobs: [][]byte{[]byte("openspec/changes/**")}},
		{Name: "tests", PathGlobs: [][]byte{[]byte("tests/**")}},
		{Name: "mechanical", PathGlobs: [][]byte{[]byte("nested/output.json")}},
		{Name: "production", LineThreshold: 400},
	}, Default: "production"}
}

func TestCNSICProfileAdmissionAndOwnership(t *testing.T) {
	want, err := evidence.NewAccountingPolicyV1(control())
	if err != nil {
		t.Fatal(err)
	}
	raw := want.CanonicalBytes()
	got, err := cnsicprofile.DecodePolicyV1(raw)
	if err != nil {
		t.Fatal(err)
	}
	view := got.View()
	view.Categories[0].Name = "changed"
	view.Categories[2].PathGlobs[0][0] = 'X'
	raw[0] = 'X'
	copy := got.CanonicalBytes()
	copy[0] = 'X'
	if !bytes.Equal(got.CanonicalBytes(), want.CanonicalBytes()) || got.Digest() != want.Digest() || !reflect.DeepEqual(got.View(), want.View()) {
		t.Fatal("admitted document changed or borrowed mutable storage")
	}
}

func TestCNSICProfileRejectsNeutralOutsideW2(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*evidence.PolicyInput)
	}{
		{"category", func(p *evidence.PolicyInput) { p.Categories[0].Name = "note" }},
		{"default", func(p *evidence.PolicyInput) { p.Default = "tests" }},
		{"glob", func(p *evidence.PolicyInput) { p.Categories[0].PathGlobs[0] = []byte("notes/**") }},
		{"threshold", func(p *evidence.PolicyInput) { p.Categories[3].LineThreshold = 7 }},
		{"other threshold", func(p *evidence.PolicyInput) { p.Categories[0].LineThreshold = 1 }},
		{"ratio", func(p *evidence.PolicyInput) {
			p.Ratios = []evidence.RatioInput{{Category: "tests", Reference: "production", Maximum: "0.5"}}
		}},
		{"order", func(p *evidence.PolicyInput) { p.Categories[0], p.Categories[1] = p.Categories[1], p.Categories[0] }},
		{"production glob", func(p *evidence.PolicyInput) { p.Categories[3].PathGlobs = [][]byte{[]byte("src/**")} }},
		{"mechanical glob", func(p *evidence.PolicyInput) { p.Categories[2].PathGlobs[0] = []byte("generated/*") }},
		{"mechanical question", func(p *evidence.PolicyInput) { p.Categories[2].PathGlobs[0] = []byte("output?") }},
		{"mechanical class", func(p *evidence.PolicyInput) { p.Categories[2].PathGlobs[0] = []byte("output[ab]") }},
		{"mechanical multiple", func(p *evidence.PolicyInput) { p.Categories[2].PathGlobs = [][]byte{[]byte("a"), []byte("b")} }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := control()
			tc.mutate(&input)
			neutral, err := evidence.NewAccountingPolicyV1(input)
			if err != nil {
				t.Fatal("neutral control:", err)
			}
			got, err := cnsicprofile.DecodePolicyV1(neutral.CanonicalBytes())
			assertError(t, got, err, "profile", "unsupported")
		})
	}
}

func TestCNSICProfileDecoderErrorsAndAuthority(t *testing.T) {
	valid, err := evidence.NewAccountingPolicyV1(control())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cnsicprofile.DecodePolicyV1(valid.CanonicalBytes()); err != nil {
		t.Fatal("positive control:", err)
	}
	for _, raw := range [][]byte{nil, []byte("{"), bytes.Replace(valid.CanonicalBytes(), []byte("{"), []byte("{\"approval\":true,"), 1)} {
		_, want := evidence.DecodeAccountingPolicyV1(raw)
		var contract *evidence.ContractError
		if !errors.As(want, &contract) {
			t.Fatal("missing neutral error")
		}
		got, err := cnsicprofile.DecodePolicyV1(raw)
		assertError(t, got, err, contract.Field, contract.Code)
		if bytes.Contains(raw, []byte("approval")) && (contract.Field != "document.approval" || contract.Code != "authority_field") {
			t.Fatal(contract)
		}
	}
}

func assertError(t *testing.T, got evidence.PolicyDocumentV1, err error, field, code string) {
	t.Helper()
	var contract *evidence.ContractError
	if !errors.As(err, &contract) || contract.Field != field || contract.Code != code || !reflect.DeepEqual(got, evidence.PolicyDocumentV1{}) {
		t.Fatalf("document/error = %#v/%v, want zero/%s:%s", got, err, field, code)
	}
}
