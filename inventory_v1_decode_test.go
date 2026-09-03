package changeevidence

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestDecodeInventoryV1CanonicalAndClosedShape(t *testing.T) {
	raw, policy := inventoryDecodeFixture(t)
	document, err := DecodeInventoryV1(raw, policy)
	if err != nil || !bytes.Equal(document.CanonicalBytes(), raw) || len(document.Entries()) != 1 {
		t.Fatalf("decoded document = (%#v, %v)", document, err)
	}
	canonical, entries := document.CanonicalBytes(), document.Entries()
	canonical[0], entries[0].Path[0] = '!', '!'
	if bytes.Equal(canonical, document.CanonicalBytes()) || bytes.Equal(entries[0].Path, document.Entries()[0].Path) {
		t.Fatal("decoder egress aliases document state")
	}

	policyDigest, content := policy.Digest().String(), strings.Repeat("b", 64)
	entry := inventoryDecodeEntry("YQ==", content, "1")
	base := string(inventoryDecodeRaw(policy, "["+entry+"]"))
	cases := []struct{ name, raw, field, code string }{
		{"missing schema", strings.Replace(base, `"schema":"`+string(InventoryContractV1)+`",`, "", 1), "document.schema", "required"},
		{"missing policy digest", strings.Replace(base, `"accounting_policy_sha256":"`+policyDigest+`",`, "", 1), "document.accounting_policy_sha256", "required"},
		{"missing entries", strings.Replace(base, `,"entries":[`+entry+`]`, "", 1), "document.entries", "required"},
		{"missing path", strings.Replace(base, `"path_b64":"YQ==",`, "", 1), "entries.path_b64", "required"},
		{"missing content digest", strings.Replace(base, `"content_sha256":"`+content+`",`, "", 1), "entries.content_sha256", "required"},
		{"missing byte length", strings.Replace(base, `,"byte_length":1`, "", 1), "entries.byte_length", "required"},
		{"unknown root field", strings.Replace(base, `,"entries":`, `,"unknown":true,"entries":`, 1), "document.unknown", "unknown_field"},
		{"unknown entry field", strings.Replace(base, `,"byte_length":1`, `,"unknown":true,"byte_length":1`, 1), "entries.unknown", "unknown_field"},
		{"duplicate root field", strings.Replace(base, `,"entries":`, `,"schema":"other","entries":`, 1), "document.schema", "duplicate_field"},
		{"escaped duplicate root field", strings.Replace(base, `,"entries":`, `,"\u0073chema":"other","entries":`, 1), "document.schema", "duplicate_field"},
		{"duplicate entry field", strings.Replace(base, `,"byte_length":1`, `,"path_b64":"YQ==","byte_length":1`, 1), "entries.path_b64", "duplicate_field"},
		{"escaped duplicate entry field", strings.Replace(base, `,"byte_length":1`, `,"\u0070ath_b64":"YQ==","byte_length":1`, 1), "entries.path_b64", "duplicate_field"},
		{"reordered root fields", `{"accounting_policy_sha256":"` + policyDigest + `","schema":"` + string(InventoryContractV1) + `","entries":[` + entry + `]}` + "\n", "document", "noncanonical"},
		{"reordered entry fields", string(inventoryDecodeRaw(policy, `[{"content_sha256":"`+content+`","path_b64":"YQ==","byte_length":1}]`)), "document", "noncanonical"},
		{"absent LF", base[:len(base)-1], "document", "noncanonical"},
		{"escaped root equivalent", strings.Replace(base, `inventory/v1`, `inventory\/v1`, 1), "document", "noncanonical"},
		{"escaped entry equivalent", strings.Replace(base, `"YQ=="`, `"\u0059Q=="`, 1), "document", "noncanonical"},
		{"leading whitespace", " " + base, "document", "noncanonical"},
		{"trailing whitespace", base + " ", "document", "noncanonical"},
		{"internal whitespace", strings.Replace(base, `,"entries":`, `, "entries":`, 1), "document", "noncanonical"},
		{"truncation", base[:len(base)-2], "document", "invalid_object"},
		{"root is null", `null`, "document", "invalid_object"},
		{"root is boolean", `true`, "document", "invalid_object"},
		{"root is string", `"document"`, "document", "invalid_object"},
		{"root is array", `[]`, "document", "invalid_object"},
		{"trailing root value", base + `null`, "document", "invalid_object"},
		{"invalid UTF-8", string(append([]byte(base[:len(base)-1]), 0xff)), "document", "invalid"},
		{"entries is null", strings.Replace(base, `[`+entry+`]`, `null`, 1), "entries", "invalid_array"},
		{"entries is boolean", strings.Replace(base, `[`+entry+`]`, `true`, 1), "entries", "invalid_array"},
		{"entries is string", strings.Replace(base, `[`+entry+`]`, `"entries"`, 1), "entries", "invalid_array"},
		{"entries is object", strings.Replace(base, `[`+entry+`]`, `{}`, 1), "entries", "invalid_array"},
		{"entry is not object", string(inventoryDecodeRaw(policy, `[null]`)), "entries", "invalid_object"},
		{"unsupported schema", strings.Replace(base, string(InventoryContractV1), `git-change-evidence.inventory/v2`, 1), "schema", "unsupported_version"},
		{"nonstring schema", strings.Replace(base, `"schema":"`+string(InventoryContractV1)+`"`, `"schema":1`, 1), "document", "invalid"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			assertInventoryDecodeError(t, []byte(test.raw), policy, test.field, test.code)
		})
	}
	for _, field := range []string{"approval", "approve", "rejection", "reject", "denial", "deny", "blocking", "block", "gating", "gate", "merge", "deploy", "release", "delivery", "deliver", "merge_gate", "delivery_authority", "delivery-authority", "delivery.authority", "deliveryAuthority"} {
		t.Run("root authority "+field, func(t *testing.T) {
			assertInventoryDecodeError(t, []byte(`{"`+field+`":true,`+base[1:]), policy, "document."+field, "authority_field")
		})
		t.Run("entry authority "+field, func(t *testing.T) {
			raw := inventoryDecodeRaw(policy, "["+strings.TrimSuffix(entry, "}")+`,"`+field+`":true}]`)
			assertInventoryDecodeError(t, raw, policy, "entries."+field, "authority_field")
		})
	}
}

func inventoryDecodeFixture(t *testing.T) ([]byte, PolicyDocumentV1) {
	t.Helper()
	policy := newPolicy(t, policyInput())
	document, err := NewInventoryV1(policy, []InventoryEntryInput{{Path: []byte("a"), Content: []byte("a")}})
	if err != nil {
		t.Fatal(err)
	}
	return document.CanonicalBytes(), policy
}

func inventoryDecodeRaw(policy PolicyDocumentV1, entries string) []byte {
	return []byte(`{"schema":"` + string(InventoryContractV1) + `","accounting_policy_sha256":"` + policy.Digest().String() + `","entries":` + entries + `}` + "\n")
}

func inventoryDecodeEntry(path, digest, length string) string {
	return `{"path_b64":"` + path + `","content_sha256":"` + digest + `","byte_length":` + length + `}`
}

func assertInventoryDecodeError(t *testing.T, raw []byte, policy PolicyDocumentV1, field, code string) {
	t.Helper()
	document, err := DecodeInventoryV1(raw, policy)
	var contract *ContractError
	if !errors.As(err, &contract) || contract.Field != field || contract.Code != code {
		t.Fatalf("error = %#v, want %s/%s", err, field, code)
	}
	if document.CanonicalBytes() != nil || document.Entries() != nil || document.Digest() != "" || document.AccountingPolicyDigest() != "" {
		t.Fatalf("rejected document = %#v, want exact zero", document)
	}
}
