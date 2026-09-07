package changeevidence

import "testing"

func TestResultV1TotalShapesRejectClosedShapesDuplicatesAndAuthorityKeys(t *testing.T) {
	valid := `{"category":"release","additions":1,"deletions":2,"non_countable":3}`
	cases := []resultShapeCase{
		{"accepts empty totals", `[]`, "", ""},
		{"accepts structural authority-like category", `[{"category":"merge","additions":0,"deletions":0,"non_countable":0}]`, "", ""},
		{"rejects malformed totals", `[`, "totals", "invalid_array"},
		{"rejects trailing totals", `[] null`, "totals", "invalid_array"},
		{"rejects null totals", `null`, "totals", "invalid_array"},
		{"rejects non-array totals", `{}`, "totals", "invalid_array"},
		{"rejects null total", `[null]`, "totals[0]", "invalid_object"},
		{"rejects array total", `[[]]`, "totals[0]", "invalid_object"},
		{"rejects unknown total field", `[` + resultShapeAdd(valid, `"metadata":null`) + `]`, "totals[0].metadata", "unknown_field"},
		{"rejects authority total field", `[` + resultShapeAdd(valid, `"DeLiVeRy-Authority":null`) + `]`, "totals[0].DeLiVeRy-Authority", "authority_field"},
	}
	members := []resultShapeMember{{"category", `"release"`, "invalid"}, {"additions", "1", "invalid_uint"}, {"deletions", "2", "invalid_uint"}, {"non_countable", "3", "invalid_uint"}}
	for _, member := range members {
		field := "totals[0]." + member.name
		cases = append(cases,
			resultShapeCase{"missing " + member.name, `[` + resultShapeDelete(valid, member) + `]`, field, "required"},
			resultShapeCase{"null " + member.name, `[` + resultShapeSet(valid, member, "null") + `]`, field, member.code},
			resultShapeCase{"wrong primitive " + member.name, `[` + resultShapeSet(valid, member, "true") + `]`, field, member.code},
			resultShapeCase{"container " + member.name, `[` + resultShapeSet(valid, member, `{}`) + `]`, field, member.code},
			resultShapeCase{"duplicate " + member.name, `[` + resultShapeAdd(valid, `"`+member.name+`":`+member.value) + `]`, field, "duplicate_field"},
			resultShapeCase{"escaped duplicate " + member.name, `[` + resultShapeAdd(valid, resultShapeEscapedKey(member.name)+`:`+member.value) + `]`, field, "duplicate_field"},
		)
		if member.code == "invalid_uint" {
			for _, value := range []string{"0", "18446744073709551615"} {
				cases = append(cases, resultShapeCase{"accepts " + member.name + "=" + value, `[` + resultShapeSet(valid, member, value) + `]`, "", ""})
			}
			for _, value := range []string{"-1", "1.5", "1e1", "18446744073709551616", `"1"`, `[]`} {
				cases = append(cases, resultShapeCase{"rejects " + member.name + "=" + value, `[` + resultShapeSet(valid, member, value) + `]`, field, member.code})
			}
		}
	}
	second := resultShapeSet(valid, members[2], "null")
	cases = append(cases, resultShapeCase{"indexes divergent second total", `[` + valid + `,` + second + `]`, "totals[1].deletions", "invalid_uint"})
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			assertResultShape(t, validateResultTotalsShapes([]byte(test.raw)), test.field, test.code)
		})
	}
}
