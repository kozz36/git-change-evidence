package changeevidence

import (
	"strings"
	"testing"
)

type resultRootShapeMember struct {
	name, value, primitive, container        string
	missing, invalid, invalidCode, duplicate string
}

func TestResultV1RootEnvelopeShapeRejectsClosedRoot(t *testing.T) {
	valid := `{"schema":"arbitrary schema","accounting_policy_sha256":"release","inventory_sha256":"merge","revisions":{"base":"arbitrary base","head":"arbitrary head"},"entries":[],"totals":[],"observations":[]}`
	for _, test := range []struct{ name, raw, field, code string }{
		{"accepts structural strings and empty arrays", valid, "", ""},
		{"accepts legal root and child-array whitespace", " \n" + strings.ReplaceAll(valid, `":[]`, `" : []`) + "\t ", "", ""},
		{"rejects malformed root", `{`, "document", "invalid_object"},
		{"rejects non-object root", `[]`, "document", "invalid_object"},
		{"rejects null root", `null`, "document", "invalid_object"},
		{"rejects trailing root value", valid + ` null`, "document", "invalid_object"},
	} {
		t.Run(test.name, func(t *testing.T) { assertResultRootShape(t, test.raw, test.field, test.code) })
	}
	members := []resultRootShapeMember{
		{"schema", `"arbitrary schema"`, "true", `[]`, "document.schema", "schema", "invalid", "document.schema"},
		{"accounting_policy_sha256", `"release"`, "true", `[]`, "document.accounting_policy_sha256", "accounting_policy_sha256", "invalid", "document.accounting_policy_sha256"},
		{"inventory_sha256", `"merge"`, "true", `[]`, "document.inventory_sha256", "inventory_sha256", "invalid", "document.inventory_sha256"},
		{"revisions", `{"base":"arbitrary base","head":"arbitrary head"}`, "true", `[]`, "document.revisions", "revisions", "invalid_object", "document.revisions"},
		{"entries", `[]`, "true", `{}`, "document.entries", "entries", "invalid_array", "document.entries"},
		{"totals", `[]`, "true", `{}`, "document.totals", "totals", "invalid_array", "document.totals"},
		{"observations", `[]`, "true", `{}`, "document.observations", "observations", "invalid_array", "document.observations"},
	}
	for _, member := range members {
		for _, test := range []struct{ name, raw, field, code string }{
			{"missing " + member.name, resultRootDelete(t, valid, member), member.missing, "required"},
			{"null " + member.name, resultRootSet(t, valid, member, "null"), member.invalid, member.invalidCode},
			{"wrong primitive " + member.name, resultRootSet(t, valid, member, member.primitive), member.invalid, member.invalidCode},
			{"wrong container " + member.name, resultRootSet(t, valid, member, member.container), member.invalid, member.invalidCode},
			{"duplicate " + member.name, resultRootAdd(t, valid, `"`+member.name+`":`+member.value), member.duplicate, "duplicate_field"},
			{"escaped duplicate " + member.name, resultRootAdd(t, valid, resultShapeEscapedKey(member.name)+`:`+member.value), member.duplicate, "duplicate_field"},
		} {
			t.Run(test.name, func(t *testing.T) { assertResultRootShape(t, test.raw, test.field, test.code) })
		}
	}
	for _, test := range []struct{ name, member, field, code string }{
		{"rejects unknown root member", `"metadata":null`, "document.metadata", "unknown_field"},
		{"rejects case authority root member", `"DeLiVeRy":null`, "document.DeLiVeRy", "authority_field"},
		{"rejects escaped authority root member", `"\u0044eliveryAuthority":null`, "document.DeliveryAuthority", "authority_field"},
		{"rejects separator authority root member", `"merge.gate":null`, "document.merge.gate", "authority_field"},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertResultRootShape(t, resultRootAdd(t, valid, test.member), test.field, test.code)
		})
	}
}

func assertResultRootShape(t *testing.T, raw, field, code string) {
	t.Helper()
	err := validateResultRootShape([]byte(raw))
	if field == "" {
		if err != nil {
			t.Fatalf("validateResultRootShape() error = %v", err)
		}
		return
	}
	contract, ok := err.(*ContractError)
	if !ok || contract.Field != field || contract.Code != code {
		t.Fatalf("validateResultRootShape() error = %#v, want ContractError{%q, %q}", err, field, code)
	}
}

func resultRootSet(t *testing.T, raw string, member resultRootShapeMember, value string) string {
	t.Helper()
	return resultRootReplace(t, raw, `"`+member.name+`":`+member.value, `"`+member.name+`":`+value)
}

func resultRootDelete(t *testing.T, raw string, member resultRootShapeMember) string {
	t.Helper()
	part := `"` + member.name + `":` + member.value
	if strings.HasPrefix(raw, "{"+part+",") {
		return resultRootReplace(t, raw, part+",", "")
	}
	return resultRootReplace(t, raw, ","+part, "")
}

func resultRootAdd(t *testing.T, raw, member string) string {
	t.Helper()
	if !strings.HasSuffix(raw, "}") {
		t.Fatalf("root fixture does not end in an object: %q", raw)
	}
	return raw[:len(raw)-1] + "," + member + "}"
}
