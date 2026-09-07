package changeevidence

import "testing"

func testResultV1ObservationShapeClosure(t *testing.T) {
	fixtures := resultV1ObservationShapeFixtures()
	for _, fixture := range fixtures {
		for _, field := range fixture.fields {
			if field.code != "invalid_uint" {
				continue
			}
			for _, test := range []struct{ name, value, code string }{
				{"zero", "0", ""}, {"max uint64", "18446744073709551615", ""},
				{"negative", "-1", "invalid_uint"}, {"fraction", "1.5", "invalid_uint"}, {"exponent", "1e0", "invalid_uint"},
				{"string", `"1"`, "invalid_uint"}, {"boolean", "true", "invalid_uint"}, {"null", "null", "invalid_uint"},
				{"array", "[]", "invalid_uint"}, {"object", "{}", "invalid_uint"}, {"overflow", "18446744073709551616", "invalid_uint"},
			} {
				t.Run(fixture.name+" "+field.name+" "+test.name, func(t *testing.T) {
					raw := "[" + observationSet(fixture.raw, field, test.value) + "]"
					if test.code == "" {
						assertObservationShape(t, raw, "", "")
					} else {
						assertObservationShape(t, raw, "observations[0]."+field.name, test.code)
					}
				})
			}
		}
		for _, field := range fixture.fields {
			if field.name != "available" && field.name != "exceeded" {
				continue
			}
			for _, test := range []struct{ name, value string }{{"numeric", "1"}, {"string", `"true"`}, {"null", "null"}, {"array", "[]"}, {"object", "{}"}} {
				t.Run(fixture.name+" "+field.name+" rejects "+test.name, func(t *testing.T) {
					assertObservationShape(t, "["+observationSet(fixture.raw, field, test.value)+"]", "observations[0]."+field.name, "invalid")
				})
			}
			if field.name == "exceeded" {
				for _, value := range []string{"true", "false"} {
					t.Run(fixture.name+" exceeded "+value, func(t *testing.T) {
						assertObservationShape(t, "["+observationSet(fixture.raw, field, value)+"]", "", "")
					})
				}
			}
		}
	}
	assertObservationShape(t, `[{"kind":"line_threshold","category":"release","limit":0,"available":false,"exceeded":true},{"kind":"ratio","category":"merge","reference":"release","limit":"anything","available":false,"exceeded":false,"denominator":null}]`, "observations[1].denominator", "unknown_field")
	for _, test := range []struct{ name, raw string }{
		{"threshold available true", `[{"kind":"line_threshold","category":"release","limit":0,"available":true,"actual":0,"exceeded":true}]`},
		{"threshold unavailable false", `[{"kind":"line_threshold","category":"release","limit":0,"available":false,"exceeded":true}]`},
		{"ratio available true", `[{"kind":"ratio","category":"release","reference":"merge","limit":"anything","available":true,"numerator":0,"denominator":0,"exceeded":true}]`},
		{"ratio unavailable false", `[{"kind":"ratio","category":"release","reference":"merge","limit":"anything","available":false,"exceeded":true}]`},
	} {
		t.Run("availability selects "+test.name, func(t *testing.T) { assertObservationShape(t, test.raw, "", "") })
	}
}
