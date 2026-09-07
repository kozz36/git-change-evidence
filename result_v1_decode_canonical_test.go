package changeevidence

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDecodeAccountingResultV1RejectsNoncanonicalRepresentations(t *testing.T) {
	policy, inventory, snapshot := resultDecodeFixture(t)
	want, err := NewAccountingResultV1(policy, inventory, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		mutate func(t *testing.T, raw []byte) []byte
	}{
		{"whitespace", func(_ *testing.T, raw []byte) []byte { return append([]byte(" "), raw...) }},
		{"key order", resultDecodeReorderedObject},
		{"escaped string", func(t *testing.T, raw []byte) []byte {
			return resultDecodeReplaceOne(t, raw, `"path_b64":"c3JjL2E=","category":"source"`, `"path_b64":"c3JjL2E=","category":"\u0073ource"`)
		}},
		{"missing final LF", func(_ *testing.T, raw []byte) []byte { return bytes.TrimSuffix(raw, []byte{'\n'}) }},
		{"trailing JSON", func(_ *testing.T, raw []byte) []byte { return append(raw, []byte(`{}`)...) }},
		{"numeric fraction", func(t *testing.T, raw []byte) []byte {
			return resultDecodeReplaceOne(t, raw, `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":2,"deletions":1`, `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":2.0,"deletions":1`)
		}},
		{"numeric exponent", func(t *testing.T, raw []byte) []byte {
			return resultDecodeReplaceOne(t, raw, `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":2,"deletions":1`, `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":2e0,"deletions":1`)
		}},
		{"numeric string", func(t *testing.T, raw []byte) []byte {
			return resultDecodeReplaceOne(t, raw, `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":2,"deletions":1`, `"path_b64":"c3JjL2E=","category":"source","measurement":{"kind":"countable","additions":"2","deletions":1`)
		}},
		{"null array", resultDecodeNullObservations},
		{"unknown root field", resultDecodeUnknownRootField},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := DecodeAccountingResultV1(test.mutate(t, want.CanonicalBytes()), policy, inventory)
			if err == nil {
				t.Fatal("noncanonical representation was accepted")
			}
			assertZeroResultDocument(t, got)
		})
	}

	hostile, err := NewAccountingResultV1(policy, inventory, resultUnit5Snapshot([]CommittedChange{{Path: string([]byte{0xff}), Lines: CommittedLineCounts{Countable: false}}}))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ name, replacement string }{
		{"URL Base64", "_w=="},
		{"unpadded Base64", "/w"},
		{"nonzero Base64 pad bits", "/x=="},
		{"Base64 CRLF", `/w\r\n==`},
	} {
		t.Run(test.name, func(t *testing.T) {
			raw := resultDecodeReplaceOne(t, hostile.CanonicalBytes(), "/w==", test.replacement)
			if test.name == "Base64 CRLF" {
				resultDecodeBase64CRLFPremise(t, raw)
			}
			got, err := DecodeAccountingResultV1(raw, policy, inventory)
			if err == nil {
				t.Fatal("Base64 alternate was accepted")
			}
			assertZeroResultDocument(t, got)
		})
	}
}

func TestDecodeAccountingResultV1RequiresExactAntecedentsAndAcceptsIndependentChain(t *testing.T) {
	policy, inventory, snapshot := resultDecodeFixture(t)
	original, err := NewAccountingResultV1(policy, inventory, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	policy2, inventory2 := resultUnit5Antecedents(t, PolicyInput{Categories: []CategoryInput{{Name: "release"}}, Default: "release"})
	independent, err := NewAccountingResultV1(policy2, inventory2, resultUnit5Snapshot([]CommittedChange{{Path: "src/a", Lines: CommittedLineCounts{Additions: 2, Countable: true}}}))
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeAccountingResultV1(independent.CanonicalBytes(), policy2, inventory2)
	if err != nil || !bytes.Equal(got.CanonicalBytes(), independent.CanonicalBytes()) || got.Entries()[0].Category != "release" {
		t.Fatalf("independent P2->I2 evidence = %#v, %v", got, err)
	}
	alternateInventory, err := NewInventoryV1(policy, []InventoryEntryInput{{Path: []byte("untracked"), Content: []byte("x")}})
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name string
		p    PolicyDocumentV1
		i    InventoryDocumentV1
	}{
		{"same policy alternate inventory", policy, alternateInventory},
		{"different policy and inventory", policy2, inventory2},
		{"cross-policy inventory", policy2, inventory},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := DecodeAccountingResultV1(original.CanonicalBytes(), test.p, test.i)
			if err == nil {
				t.Fatal("substituted antecedents were accepted")
			}
			assertZeroResultDocument(t, got)
		})
	}
}

func TestDecodeAccountingResultV1AcceptsEmptyEqualRevisionEvidence(t *testing.T) {
	policy, inventory, _ := resultDecodeFixture(t)
	empty, err := NewAccountingResultV1(policy, inventory, NewCommittedSnapshot(resultUnit5Base, resultUnit5Base, nil))
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeAccountingResultV1(empty.CanonicalBytes(), policy, inventory)
	if err != nil || !bytes.Equal(got.CanonicalBytes(), empty.CanonicalBytes()) || got.Revisions() != empty.Revisions() || got.Entries() == nil {
		t.Fatalf("empty equal-revision evidence = %#v, %v", got, err)
	}
}

func resultDecodeReorderedObject(t *testing.T, raw []byte) []byte {
	t.Helper()
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		t.Fatal(err)
	}
	reordered, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	reordered = append(reordered, '\n')
	if bytes.Equal(reordered, raw) {
		t.Fatal("key-order mutation did not change fixture")
	}
	return reordered
}

func resultDecodeUnknownRootField(t *testing.T, raw []byte) []byte {
	t.Helper()
	if len(raw) < 2 || raw[len(raw)-2] != '}' || raw[len(raw)-1] != '\n' {
		t.Fatal("canonical fixture has no terminal object and LF")
	}
	return append(append([]byte(nil), raw[:len(raw)-2]...), []byte(",\"unknown\":true}\n")...)
}

func resultDecodeNullObservations(t *testing.T, raw []byte) []byte {
	t.Helper()
	index := bytes.LastIndex(raw, []byte(`,"observations":`))
	if index < 0 {
		t.Fatal("canonical fixture has no observations array")
	}
	return append(append([]byte(nil), raw[:index]...), []byte(",\"observations\":null}\n")...)
}
