package changeevidence

import "testing"

func TestResultV1EntryShapesRejectClosedShapesDuplicatesAndAuthorityKeys(t *testing.T) {
	entry := func(measurement string) string {
		return `{"path_b64":"not-base64","category":"release","measurement":` + measurement + `}`
	}
	countable := entry(`{"kind":"countable","additions":2,"deletions":1}`)
	cases := []struct {
		name, raw, field, code string
	}{
		{"accepts countable release category", "[" + countable + "]", "", ""},
		{"accepts countable zero and max uint64", `[{"path_b64":"max","category":"release","measurement":{"kind":"countable","additions":0,"deletions":18446744073709551615}}]`, "", ""},
		{"accepts non-countable merge category", "[" + `{"path_b64":"also-not-base64","category":"merge","measurement":{"kind":"non_countable"}}` + "]", "", ""},
		{"malformed entries", `[`, "entries", "invalid_array"},
		{"trailing entry value", `[] null`, "entries", "invalid_array"},
		{"null entries", `null`, "entries", "invalid_array"},
		{"non-array entries", `{}`, "entries", "invalid_array"},
		{"null entry", `[null]`, "entries[0]", "invalid_object"},
		{"missing path", `[{"category":"release","measurement":{"kind":"non_countable"}}]`, "entries[0].path_b64", "required"},
		{"missing category", `[{"path_b64":"x","measurement":{"kind":"non_countable"}}]`, "entries[0].category", "required"},
		{"missing measurement", `[{"path_b64":"x","category":"release"}]`, "entries[0].measurement", "required"},
		{"second entry path is indexed", "[" + countable + `,{"path_b64":"x","category":"merge","measurement":{"kind":"countable","additions":1}}]`, "entries[1].measurement.deletions", "required"},
		{"unknown entry field", "[" + countable[:len(countable)-1] + `,"unknown":true}]`, "entries[0].unknown", "unknown_field"},
		{"authority entry field", "[" + countable[:len(countable)-1] + `,"review-gate":true}]`, "entries[0].review-gate", "authority_field"},
		{"approval entry field", "[" + countable[:len(countable)-1] + `,"approval":true}]`, "entries[0].approval", "authority_field"},
		{"normalized denial entry field", "[" + countable[:len(countable)-1] + `,"DeNiAl-status":true}]`, "entries[0].DeNiAl-status", "authority_field"},
		{"duplicate path", `[{"path_b64":"x","path_b64":"y","category":"release","measurement":{"kind":"non_countable"}}]`, "entries[0].path_b64", "duplicate_field"},
		{"escaped duplicate path", `[{"path_b64":"x","\u0070ath_b64":"y","category":"release","measurement":{"kind":"non_countable"}}]`, "entries[0].path_b64", "duplicate_field"},
		{"non-string path", `[{"path_b64":1,"category":"release","measurement":{"kind":"non_countable"}}]`, "entries[0].path_b64", "invalid"},
		{"non-string category", `[{"path_b64":"x","category":true,"measurement":{"kind":"non_countable"}}]`, "entries[0].category", "invalid"},
		{"non-object measurement", `[{"path_b64":"x","category":"release","measurement":[]}]`, "entries[0].measurement", "invalid_object"},
		{"missing measurement kind", "[" + entry(`{}`) + "]", "entries[0].measurement.kind", "required"},
		{"countable missing additions", "[" + entry(`{"kind":"countable","deletions":1}`) + "]", "entries[0].measurement.additions", "required"},
		{"countable missing deletions", "[" + entry(`{"kind":"countable","additions":1}`) + "]", "entries[0].measurement.deletions", "required"},
		{"non-countable additions are closed", "[" + entry(`{"kind":"non_countable","additions":1}`) + "]", "entries[0].measurement.additions", "unknown_field"},
		{"non-countable deletions are closed", "[" + entry(`{"kind":"non_countable","deletions":1}`) + "]", "entries[0].measurement.deletions", "unknown_field"},
		{"unknown measurement field", "[" + entry(`{"kind":"non_countable","unknown":true}`) + "]", "entries[0].measurement.unknown", "unknown_field"},
		{"authority measurement field", "[" + entry(`{"kind":"non_countable","delivery_authority":true}`) + "]", "entries[0].measurement.delivery_authority", "authority_field"},
		{"normalized rejection measurement field", "[" + entry(`{"kind":"non_countable","ReJeCt_ion":true}`) + "]", "entries[0].measurement.ReJeCt_ion", "authority_field"},
		{"duplicate measurement kind", "[" + entry(`{"kind":"countable","kind":"non_countable"}`) + "]", "entries[0].measurement.kind", "duplicate_field"},
		{"escaped duplicate measurement kind", "[" + entry(`{"kind":"countable","\u006bind":"non_countable"}`) + "]", "entries[0].measurement.kind", "duplicate_field"},
		{"non-string measurement kind", "[" + entry(`{"kind":1}`) + "]", "entries[0].measurement.kind", "invalid"},
		{"unsupported measurement kind", "[" + entry(`{"kind":"unsupported"}`) + "]", "entries[0].measurement.kind", "unsupported_kind"},
		{"authority measurement discriminator", "[" + entry(`{"kind":"merge-gate"}`) + "]", "entries[0].measurement.kind", "authority_field"},
		{"non-uint additions", "[" + entry(`{"kind":"countable","additions":true,"deletions":1}`) + "]", "entries[0].measurement.additions", "invalid_uint"},
		{"non-uint deletions", "[" + entry(`{"kind":"countable","additions":1,"deletions":"1"}`) + "]", "entries[0].measurement.deletions", "invalid_uint"},
		{"negative additions", "[" + entry(`{"kind":"countable","additions":-1,"deletions":1}`) + "]", "entries[0].measurement.additions", "invalid_uint"},
		{"negative deletions", "[" + entry(`{"kind":"countable","additions":1,"deletions":-1}`) + "]", "entries[0].measurement.deletions", "invalid_uint"},
		{"fractional additions", "[" + entry(`{"kind":"countable","additions":1.5,"deletions":1}`) + "]", "entries[0].measurement.additions", "invalid_uint"},
		{"fractional deletions", "[" + entry(`{"kind":"countable","additions":1,"deletions":1.5}`) + "]", "entries[0].measurement.deletions", "invalid_uint"},
		{"additions above max uint64", "[" + entry(`{"kind":"countable","additions":18446744073709551616,"deletions":1}`) + "]", "entries[0].measurement.additions", "invalid_uint"},
		{"deletions above max uint64", "[" + entry(`{"kind":"countable","additions":1,"deletions":18446744073709551616}`) + "]", "entries[0].measurement.deletions", "invalid_uint"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			err := validateResultEntryShapes([]byte(test.raw))
			if test.field == "" {
				if err != nil {
					t.Fatalf("validateResultEntryShapes() error = %v", err)
				}
				return
			}
			contract, ok := err.(*ContractError)
			if !ok || contract.Field != test.field || contract.Code != test.code {
				t.Fatalf("validateResultEntryShapes() error = %#v, want ContractError{%q, %q}", err, test.field, test.code)
			}
		})
	}
}
