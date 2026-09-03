package changeevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"reflect"
	"testing"
)

func TestNewInventoryV1ExactContentAndPaths(t *testing.T) {
	policy := newPolicy(t, policyInput())
	entries := []InventoryEntryInput{
		{Path: []byte("bin/data"), Content: []byte{0, 0xff, '\r', '\n'}},
		{Path: []byte("empty"), Content: []byte{}},
		{Path: []byte("line-crlf"), Content: []byte("line\r\n")},
		{Path: []byte("line-lf"), Content: []byte("line\n")},
	}
	got, err := NewInventoryV1(policy, entries)
	if err != nil {
		t.Fatal(err)
	}
	want := []struct {
		path    []byte
		content []byte
	}{
		{[]byte("bin/data"), []byte{0, 0xff, '\r', '\n'}},
		{[]byte("empty"), nil},
		{[]byte("line-crlf"), []byte("line\r\n")},
		{[]byte("line-lf"), []byte("line\n")},
	}
	for n, entry := range got.Entries() {
		sum := sha256.Sum256(want[n].content)
		if !bytes.Equal(entry.Path, want[n].path) || entry.ContentSHA256 != Digest(hex.EncodeToString(sum[:])) || entry.ByteLength != uint64(len(want[n].content)) {
			t.Fatalf("entry %d = %#v, want identity for %q", n, entry, want[n].path)
		}
	}
	if got.Entries()[2].ContentSHA256 == got.Entries()[3].ContentSHA256 || got.Entries()[2].ByteLength == got.Entries()[3].ByteLength {
		t.Fatal("CRLF and LF content identities collapsed")
	}
	ordered := []InventoryEntryInput{{Path: []byte("a"), Content: []byte("1")}, {Path: []byte{'a', '/', 0xff}, Content: []byte("2")}, {Path: []byte("a0"), Content: []byte("3")}, {Path: []byte{0xff}, Content: []byte("4")}}
	reversed := append([]InventoryEntryInput(nil), ordered...)
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}
	ascending, err := NewInventoryV1(policy, ordered)
	if err != nil {
		t.Fatal(err)
	}
	descending, err := NewInventoryV1(policy, reversed)
	if err != nil || !bytes.Equal(ascending.CanonicalBytes(), descending.CanonicalBytes()) {
		t.Fatal("equivalent raw-path sets have different canonical bytes")
	}
	for n, entry := range ascending.Entries() {
		if !bytes.Equal(entry.Path, ordered[n].Path) {
			t.Fatalf("raw-byte ordering at %d = %q, want %q", n, entry.Path, ordered[n].Path)
		}
	}
	for _, duplicate := range [][]InventoryEntryInput{{{Path: []byte("same"), Content: []byte("x")}, {Path: []byte("same"), Content: []byte("x")}}, {{Path: []byte("same"), Content: []byte("x")}, {Path: []byte("same"), Content: []byte("y")}}} {
		_, err := NewInventoryV1(policy, duplicate)
		var contract *ContractError
		if !errors.As(err, &contract) || contract.Field != "entries.path" || contract.Code != "duplicate_path" {
			t.Fatalf("duplicate error = %#v", err)
		}
	}
	for _, path := range [][]byte{nil, []byte("/absolute"), []byte("a//b"), []byte("./a"), []byte("a/./b"), []byte("a/../b"), []byte("../a"), []byte{'a', 0, 'b'}} {
		_, err := NewInventoryV1(policy, []InventoryEntryInput{{Path: path, Content: []byte("x")}})
		var contract *ContractError
		if !errors.As(err, &contract) || contract.Field != "entries.path" || contract.Code != "invalid_path" {
			t.Fatalf("invalid path %q error = %#v", path, err)
		}
	}
	validOpaque := []InventoryEntryInput{{Path: []byte{0xff}, Content: []byte("1")}, {Path: []byte("Case"), Content: []byte("2")}, {Path: []byte(`a\\b`), Content: []byte("3")}, {Path: []byte{'c', '\n'}, Content: []byte("4")}}
	if got, err := NewInventoryV1(policy, validOpaque); err != nil || len(got.Entries()) != len(validOpaque) {
		t.Fatalf("valid opaque paths = (%#v, %v)", got, err)
	}
}

func canonicalHighBitPathFixture(t *testing.T) (PolicyDocumentV1, InventoryEntryInput) {
	t.Helper()
	return newPolicy(t, policyInput()), InventoryEntryInput{Path: []byte{0xff, '\n', 'a'}, Content: []byte("x")}
}

func TestNewInventoryV1CanonicalContract(t *testing.T) {
	policy, entry := canonicalHighBitPathFixture(t)
	got, err := NewInventoryV1(policy, []InventoryEntryInput{entry})
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(`{"schema":"git-change-evidence.inventory/v1","accounting_policy_sha256":"` + policy.Digest().String() + `","entries":[{"path_b64":"/wph","content_sha256":"2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881","byte_length":1}]}` + "\n")
	sum := sha256.Sum256(want)
	if !bytes.Equal(got.CanonicalBytes(), want) || got.Digest().String() != hex.EncodeToString(sum[:]) || got.AccountingPolicyDigest() != policy.Digest() {
		t.Fatalf("canonical/link/digest = %q/%q/%q", got.CanonicalBytes(), got.AccountingPolicyDigest(), got.Digest())
	}
	equivalent, err := NewInventoryV1(newPolicy(t, policyInput()), []InventoryEntryInput{entry})
	if err != nil || !bytes.Equal(got.CanonicalBytes(), equivalent.CanonicalBytes()) || got.Digest() != equivalent.Digest() {
		t.Fatal("byte-identical policies do not reproduce inventory identity")
	}
	otherInput := policyInput()
	otherInput.Ratios[0].Maximum = "0.1"
	other, err := NewInventoryV1(newPolicy(t, otherInput), []InventoryEntryInput{entry})
	if err != nil || other.AccountingPolicyDigest() == got.AccountingPolicyDigest() || bytes.Equal(got.CanonicalBytes(), other.CanonicalBytes()) || got.Digest() == other.Digest() {
		t.Fatal("byte-distinct policy did not produce distinct inventory provenance")
	}
	var zero PolicyDocumentV1
	_, err = NewInventoryV1(zero, nil)
	var contract *ContractError
	if !errors.As(err, &contract) || contract.Field != "accounting_policy" || contract.Code != "invalid_document" {
		t.Fatalf("zero-policy error = %#v", err)
	}
	input := reflect.TypeFor[InventoryEntryInput]()
	if _, ok := input.FieldByName("ContentSHA256"); ok {
		t.Fatal("input accepts a content digest")
	}
	if _, ok := input.FieldByName("ByteLength"); ok {
		t.Fatal("input accepts a byte length")
	}
}
