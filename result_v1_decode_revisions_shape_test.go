package changeevidence

import (
	"fmt"
	"strings"
	"testing"
)

func TestResultV1RevisionShapesRejectClosedShapesDuplicatesAndAuthorityKeys(t *testing.T) {
	valid := `{"base":"release","head":"merge"}`
	cases := []resultShapeCase{
		{"accepts arbitrary strings", valid, "", ""},
		{"accepts empty revision strings", `{"base":"","head":""}`, "", ""},
		{"rejects malformed revisions", `{`, "revisions", "invalid_object"},
		{"rejects trailing revisions", valid + ` null`, "revisions", "invalid_object"},
		{"rejects null revisions", `null`, "revisions", "invalid_object"},
		{"rejects array revisions", `[]`, "revisions", "invalid_object"},
		{"rejects unknown revision field", resultShapeAdd(valid, `"metadata":null`), "revisions.metadata", "unknown_field"},
		{"rejects authority revision field", resultShapeAdd(valid, `"DeLiVeRy-Authority":null`), "revisions.DeLiVeRy-Authority", "authority_field"},
		{"rejects escaped authority revision field", resultShapeAdd(valid, `"\u0044eliveryAuthority":null`), "revisions.DeliveryAuthority", "authority_field"},
	}
	for _, member := range []resultShapeMember{{"base", `"release"`, "invalid"}, {"head", `"merge"`, "invalid"}} {
		cases = append(cases,
			resultShapeCase{"missing " + member.name, resultShapeDelete(valid, member), "revisions." + member.name, "required"},
			resultShapeCase{"null " + member.name, resultShapeSet(valid, member, "null"), "revisions." + member.name, member.code},
			resultShapeCase{"wrong primitive " + member.name, resultShapeSet(valid, member, "true"), "revisions." + member.name, member.code},
			resultShapeCase{"container " + member.name, resultShapeSet(valid, member, `{}`), "revisions." + member.name, member.code},
			resultShapeCase{"duplicate " + member.name, resultShapeAdd(valid, `"`+member.name+`":`+member.value), "revisions." + member.name, "duplicate_field"},
			resultShapeCase{"escaped duplicate " + member.name, resultShapeAdd(valid, resultShapeEscapedKey(member.name)+`:`+member.value), "revisions." + member.name, "duplicate_field"},
		)
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			assertResultShape(t, validateResultRevisionsShape([]byte(test.raw), "revisions"), test.field, test.code)
		})
	}
}

type resultShapeCase struct{ name, raw, field, code string }
type resultShapeMember struct{ name, value, code string }

func resultShapeSet(raw string, member resultShapeMember, value string) string {
	return strings.Replace(raw, `"`+member.name+`":`+member.value, `"`+member.name+`":`+value, 1)
}
func resultShapeDelete(raw string, member resultShapeMember) string {
	part := `"` + member.name + `":` + member.value
	if strings.HasPrefix(raw, "{"+part+",") {
		return strings.Replace(raw, part+",", "", 1)
	}
	return strings.Replace(raw, ","+part, "", 1)
}
func resultShapeAdd(raw, member string) string { return raw[:len(raw)-1] + "," + member + "}" }
func resultShapeEscapedKey(name string) string { return fmt.Sprintf(`"\u%04x%s"`, name[0], name[1:]) }

func assertResultShape(t *testing.T, err error, field, code string) {
	t.Helper()
	if field == "" {
		if err != nil {
			t.Fatalf("shape error = %v", err)
		}
		return
	}
	contract, ok := err.(*ContractError)
	if !ok || contract.Field != field || contract.Code != code {
		t.Fatalf("shape error = %#v, want ContractError{%q, %q}", err, field, code)
	}
}
