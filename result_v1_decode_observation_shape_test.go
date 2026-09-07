package changeevidence

import (
	"fmt"
	"strings"
	"testing"
)

type observationShapeField struct{ name, value, code string }
type observationShapeFixture struct {
	name, raw         string
	fields, forbidden []observationShapeField
}

func TestResultV1ObservationShapesRejectClosedShapesDuplicatesAndAuthorityKeys(t *testing.T) {
	testResultV1ObservationShapeCore(t)
	testResultV1ObservationShapeClosure(t)
}

func testResultV1ObservationShapeCore(t *testing.T) {
	for _, test := range []struct{ name, raw, field, code string }{
		{"accepts empty observations", `[]`, "", ""}, {"rejects leading whitespace", ` []`, "observations", "invalid_array"},
		{"rejects malformed observations", `[`, "observations", "invalid_array"}, {"rejects trailing observations", `[] null`, "observations", "invalid_array"},
		{"rejects null observations", `null`, "observations", "invalid_array"}, {"rejects object observations", `{}`, "observations", "invalid_array"},
		{"rejects null observation", `[null]`, "observations[0]", "invalid_object"}, {"rejects array observation", `[[]]`, "observations[0]", "invalid_object"},
	} {
		t.Run(test.name, func(t *testing.T) { assertObservationShape(t, test.raw, test.field, test.code) })
	}
	fixtures := resultV1ObservationShapeFixtures()
	for _, fixture := range fixtures {
		t.Run("accepts "+fixture.name, func(t *testing.T) { assertObservationShape(t, "["+fixture.raw+"]", "", "") })
		for _, name := range []string{"metadata", "approval", "DeLiVeRy-Authority"} {
			code := "unknown_field"
			if name != "metadata" {
				code = "authority_field"
			}
			t.Run(fixture.name+" rejects "+name, func(t *testing.T) {
				assertObservationShape(t, "["+observationAdd(fixture.raw, `"`+name+`":null`)+"]", "observations[0]."+name, code)
			})
		}
		for _, field := range fixture.fields {
			tests := []struct{ name, raw, code string }{
				{"missing", observationDelete(fixture.raw, field), "required"},
				{"null", observationSet(fixture.raw, field, "null"), field.code},
				{"wrong primitive", observationSet(fixture.raw, field, observationWrongPrimitive(field)), field.code},
				{"container", observationSet(fixture.raw, field, "{}"), field.code},
				{"duplicate", observationAdd(fixture.raw, `"`+field.name+`":`+field.value), "duplicate_field"},
				{"escaped duplicate", observationAdd(fixture.raw, observationEscapedKey(field.name)+":"+field.value), "duplicate_field"},
			}
			for _, test := range tests {
				t.Run(fixture.name+" "+field.name+" "+test.name, func(t *testing.T) {
					assertObservationShape(t, "["+test.raw+"]", "observations[0]."+field.name, test.code)
				})
			}
		}
		for _, field := range fixture.forbidden {
			for _, value := range []string{field.value, "null"} {
				t.Run(fixture.name+" closes "+field.name+"="+value, func(t *testing.T) {
					assertObservationShape(t, "["+observationAdd(fixture.raw, `"`+field.name+`":`+value)+"]", "observations[0]."+field.name, "unknown_field")
				})
			}
		}
		for _, value := range []struct{ name, value, code string }{{"unsupported", `"unsupported"`, "unsupported_kind"}, {"authority", `"merge-gate"`, "authority_field"}} {
			t.Run(fixture.name+" kind "+value.name, func(t *testing.T) {
				assertObservationShape(t, "["+observationSet(fixture.raw, fixture.fields[0], value.value)+"]", "observations[0].kind", value.code)
			})
		}
	}
	second := observationSet(fixtures[0].raw, fixtures[0].fields[1], "null")
	assertObservationShape(t, "["+fixtures[0].raw+","+second+"]", "observations[1].category", "invalid")
}

func resultV1ObservationShapeFixtures() []observationShapeFixture {
	return []observationShapeFixture{
		{"available threshold", `{"kind":"line_threshold","category":"release","limit":1,"available":true,"actual":2,"exceeded":false}`, []observationShapeField{{"kind", `"line_threshold"`, "invalid"}, {"category", `"release"`, "invalid"}, {"limit", "1", "invalid_uint"}, {"available", "true", "invalid"}, {"actual", "2", "invalid_uint"}, {"exceeded", "false", "invalid"}}, []observationShapeField{{"reference", `"merge"`, ""}, {"numerator", "1", ""}, {"denominator", "1", ""}}},
		{"unavailable threshold", `{"kind":"line_threshold","category":"merge","limit":0,"available":false,"exceeded":false}`, []observationShapeField{{"kind", `"line_threshold"`, "invalid"}, {"category", `"merge"`, "invalid"}, {"limit", "0", "invalid_uint"}, {"available", "false", "invalid"}, {"exceeded", "false", "invalid"}}, []observationShapeField{{"actual", "0", ""}, {"reference", `"release"`, ""}, {"numerator", "0", ""}, {"denominator", "0", ""}}},
		{"available ratio", `{"kind":"ratio","category":"release","reference":"merge","limit":"not-a-decimal","available":true,"numerator":2,"denominator":0,"exceeded":false}`, []observationShapeField{{"kind", `"ratio"`, "invalid"}, {"category", `"release"`, "invalid"}, {"reference", `"merge"`, "invalid"}, {"limit", `"not-a-decimal"`, "invalid"}, {"available", "true", "invalid"}, {"numerator", "2", "invalid_uint"}, {"denominator", "0", "invalid_uint"}, {"exceeded", "false", "invalid"}}, []observationShapeField{{"actual", "0", ""}}},
		{"unavailable ratio", `{"kind":"ratio","category":"merge","reference":"release","limit":"release-authority","available":false,"exceeded":false}`, []observationShapeField{{"kind", `"ratio"`, "invalid"}, {"category", `"merge"`, "invalid"}, {"reference", `"release"`, "invalid"}, {"limit", `"release-authority"`, "invalid"}, {"available", "false", "invalid"}, {"exceeded", "false", "invalid"}}, []observationShapeField{{"actual", "0", ""}, {"numerator", "0", ""}, {"denominator", "0", ""}}},
	}
}

func observationSet(raw string, field observationShapeField, value string) string {
	return strings.Replace(raw, `"`+field.name+`":`+field.value, `"`+field.name+`":`+value, 1)
}
func observationDelete(raw string, field observationShapeField) string {
	part := `"` + field.name + `":` + field.value
	if strings.HasPrefix(raw, "{"+part+",") {
		return strings.Replace(raw, part+",", "", 1)
	}
	return strings.Replace(raw, ","+part, "", 1)
}
func observationAdd(raw, member string) string { return raw[:len(raw)-1] + "," + member + "}" }
func observationEscapedKey(name string) string { return fmt.Sprintf(`"\u%04x%s"`, name[0], name[1:]) }
func observationWrongPrimitive(field observationShapeField) string {
	if field.code == "invalid_uint" {
		return "true"
	}
	return "1"
}
func assertObservationShape(t *testing.T, raw, field, code string) {
	t.Helper()
	err := validateResultObservationShapes([]byte(raw))
	if field == "" {
		if err != nil {
			t.Fatalf("validateResultObservationShapes() error = %v", err)
		}
		return
	}
	contract, ok := err.(*ContractError)
	if !ok || contract.Field != field || contract.Code != code {
		t.Fatalf("validateResultObservationShapes() error = %#v, want ContractError{%q, %q}", err, field, code)
	}
}
