package changeevidence

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func TestDecodeInventoryV1ValueValidation(t *testing.T) {
	_, policy := inventoryDecodeFixture(t)
	policyDigest, content := policy.Digest().String(), strings.Repeat("b", 64)
	entry := inventoryDecodeEntry("YQ==", content, "1")
	raw := string(inventoryDecodeRaw(policy, "["+entry+"]"))
	for _, test := range []struct{ name, value string }{
		{"invalid", "!"}, {"URL", "YQ--"}, {"unpadded", "YQ"}, {"whitespace", "YQ== "},
	} {
		t.Run("base64 "+test.name, func(t *testing.T) {
			assertInventoryDecodeError(t, []byte(strings.Replace(raw, `"YQ=="`, `"`+test.value+`"`, 1)), policy, "entries.path_b64", "invalid")
		})
	}
	for _, test := range []struct {
		name string
		path []byte
	}{
		{"empty", nil}, {"absolute", []byte("/a")}, {"NUL", []byte("a\x00b")}, {"empty component", []byte("a//b")},
		{"dot component", []byte("a/./b")}, {"dotdot component", []byte("a/../b")}, {"trailing separator", []byte("a/")},
	} {
		t.Run("invalid path "+test.name, func(t *testing.T) {
			encoded := base64.StdEncoding.EncodeToString(test.path)
			assertInventoryDecodeError(t, []byte(strings.Replace(raw, "YQ==", encoded, 1)), policy, "entries.path_b64", "invalid_path")
		})
	}
	for _, test := range []struct{ name, raw, field string }{
		{"40-character policy digest", strings.Replace(raw, policyDigest, strings.Repeat("a", 40), 1), "accounting_policy_sha256"},
		{"uppercase policy digest", strings.Replace(raw, policyDigest, strings.ToUpper(policyDigest), 1), "accounting_policy_sha256"},
		{"40-character content digest", strings.Replace(raw, content, strings.Repeat("b", 40), 1), "entries.content_sha256"},
		{"uppercase content digest", strings.Replace(raw, content, strings.ToUpper(content), 1), "entries.content_sha256"},
		{"nonstring policy digest", strings.Replace(raw, `"`+policyDigest+`"`, `1`, 1), "accounting_policy_sha256"},
		{"nonstring content digest", strings.Replace(raw, `"`+content+`"`, `1`, 1), "entries.content_sha256"},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertInventoryDecodeError(t, []byte(test.raw), policy, test.field, "invalid_digest")
		})
	}
	for _, test := range []struct{ name, value, field, code string }{
		{"signed", "-1", "entries.byte_length", "invalid_uint"}, {"fraction", "1.5", "entries.byte_length", "invalid_uint"},
		{"exponent", "1e0", "entries.byte_length", "invalid_uint"}, {"string", `"1"`, "entries.byte_length", "invalid_uint"},
		{"leading zero", "01", "document", "invalid_object"}, {"overflow", "18446744073709551616", "entries.byte_length", "invalid_uint"},
	} {
		t.Run("byte length "+test.name, func(t *testing.T) {
			invalid := strings.Replace(raw, `"byte_length":1`, `"byte_length":`+test.value, 1)
			assertInventoryDecodeError(t, []byte(invalid), policy, test.field, test.code)
		})
	}
	maximum := []byte(strings.Replace(raw, `"byte_length":1`, `"byte_length":18446744073709551615`, 1))
	if document, err := DecodeInventoryV1(maximum, policy); err != nil || document.Entries()[0].ByteLength != ^uint64(0) {
		t.Fatalf("maximum byte length = (%#v, %v)", document, err)
	}
	zero := []byte(strings.Replace(raw, `"byte_length":1`, `"byte_length":0`, 1))
	if document, err := DecodeInventoryV1(zero, policy); err != nil || len(document.Entries()) != 1 || document.Entries()[0].ByteLength != 0 {
		t.Fatalf("zero byte length = (%#v, %v)", document, err)
	}
}

func TestDecodeInventoryV1OrderLinkAndCanonicality(t *testing.T) {
	raw, policy := inventoryDecodeFixture(t)
	assertInventoryDecodeError(t, raw, PolicyDocumentV1{}, "accounting_policy", "invalid_document")
	otherInput := policyInput()
	otherInput.Categories[0].LineThreshold = 1
	otherPolicy := newPolicy(t, otherInput)
	if policy.Digest() == otherPolicy.Digest() {
		t.Fatal("expected distinct valid policy documents")
	}
	assertInventoryDecodeError(t, raw, otherPolicy, "accounting_policy_sha256", "mismatch")
	empty, err := NewInventoryV1(policy, nil)
	if err != nil {
		t.Fatal(err)
	}
	if decoded, err := DecodeInventoryV1(empty.CanonicalBytes(), policy); err != nil || decoded.Entries() == nil || len(decoded.Entries()) != 0 {
		t.Fatalf("empty round trip = (%#v, %v)", decoded, err)
	}

	content := strings.Repeat("b", 64)
	entry := inventoryDecodeEntry("YQ==", content, "1")
	for _, test := range []struct{ name, entries, field, code string }{
		{"duplicate same content", "[" + entry + "," + entry + "]", "entries.path_b64", "duplicate_path"},
		{"duplicate different content", "[" + entry + "," + inventoryDecodeEntry("YQ==", strings.Repeat("c", 64), "1") + "]", "entries.path_b64", "duplicate_path"},
		{"noncanonical order", "[" + inventoryDecodeEntry("Yg==", content, "1") + "," + entry + "]", "entries", "noncanonical_order"},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertInventoryDecodeError(t, inventoryDecodeRaw(policy, test.entries), policy, test.field, test.code)
		})
	}
	source, err := NewInventoryV1(policy, []InventoryEntryInput{
		{Path: []byte("opaque"), Content: []byte("opaque")}, {Path: []byte("a\x01b"), Content: []byte("control")}, {Path: []byte{0xff}, Content: []byte("high-bit")},
	})
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeInventoryV1(source.CanonicalBytes(), policy)
	if err != nil || !bytes.Equal(decoded.CanonicalBytes(), source.CanonicalBytes()) || decoded.Digest() != source.Digest() {
		t.Fatalf("opaque/control/high-bit round trip = (%#v, %v)", decoded, err)
	}
	for _, want := range [][]byte{[]byte("opaque"), []byte("a\x01b"), {0xff}} {
		found := false
		for _, got := range decoded.Entries() {
			found = found || bytes.Equal(got.Path, want)
		}
		if !found {
			t.Fatalf("missing round-trip path %x", want)
		}
	}
}
