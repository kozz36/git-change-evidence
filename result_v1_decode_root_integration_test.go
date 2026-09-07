package changeevidence

import (
	"strings"
	"testing"
)

func TestResultV1RootEnvelopeShapeDelegatesClosedSubtrees(t *testing.T) {
	valid := resultRootEnvelopeFixture()
	assertResultRootShape(t, valid, "", "")
	for _, test := range []struct{ name, old, new, field, code string }{
		{"revisions duplicate", `{"base":"release","head":"merge"}`, `{"base":"release","head":"merge","base":"again"}`, "revisions.base", "duplicate_field"},
		{"revisions escaped duplicate", `{"base":"release","head":"merge"}`, `{"base":"release","head":"merge","\u0062ase":"again"}`, "revisions.base", "duplicate_field"},
		{"revisions authority", `{"base":"release","head":"merge"}`, `{"base":"release","head":"merge","delivery_authority":null}`, "revisions.delivery_authority", "authority_field"},
		{"entries duplicate", `"measurement":{"kind":"countable","additions":1,"deletions":2}}`, `"measurement":{"kind":"countable","additions":1,"deletions":2},"path_b64":"second"}`, "entries[0].path_b64", "duplicate_field"},
		{"entries escaped duplicate", `"measurement":{"kind":"countable","additions":1,"deletions":2}}`, `"measurement":{"kind":"countable","additions":1,"deletions":2},"\u0070ath_b64":"second"}`, "entries[0].path_b64", "duplicate_field"},
		{"entries authority", `"measurement":{"kind":"countable","additions":1,"deletions":2}}`, `"measurement":{"kind":"countable","additions":1,"deletions":2},"delivery_authority":null}`, "entries[0].delivery_authority", "authority_field"},
		{"totals duplicate", `{"category":"release","additions":1,"deletions":2,"non_countable":0}`, `{"category":"release","additions":1,"deletions":2,"non_countable":0,"category":"again"}`, "totals[0].category", "duplicate_field"},
		{"totals escaped duplicate", `{"category":"release","additions":1,"deletions":2,"non_countable":0}`, `{"category":"release","additions":1,"deletions":2,"non_countable":0,"\u0063ategory":"again"}`, "totals[0].category", "duplicate_field"},
		{"totals authority", `{"category":"release","additions":1,"deletions":2,"non_countable":0}`, `{"category":"release","additions":1,"deletions":2,"non_countable":0,"delivery_authority":null}`, "totals[0].delivery_authority", "authority_field"},
		{"observations duplicate", `{"kind":"ratio","category":"merge","reference":"release","limit":"other","available":false,"exceeded":false}`, `{"kind":"ratio","category":"merge","reference":"release","limit":"other","available":false,"exceeded":false,"kind":"ratio"}`, "observations[3].kind", "duplicate_field"},
		{"observations escaped duplicate", `{"kind":"ratio","category":"merge","reference":"release","limit":"other","available":false,"exceeded":false}`, `{"kind":"ratio","category":"merge","reference":"release","limit":"other","available":false,"exceeded":false,"\u006bind":"ratio"}`, "observations[3].kind", "duplicate_field"},
		{"observations authority", `{"kind":"ratio","category":"merge","reference":"release","limit":"other","available":false,"exceeded":false}`, `{"kind":"ratio","category":"merge","reference":"release","limit":"other","available":false,"exceeded":false,"delivery_authority":null}`, "observations[3].delivery_authority", "authority_field"},
		{"entries invalid nested uint", `"measurement":{"kind":"countable","additions":1,"deletions":2}`, `"measurement":{"kind":"countable","additions":true,"deletions":2}`, "entries[0].measurement.additions", "invalid_uint"},
		{"totals indexes divergent second item", `{"category":"merge","additions":3,"deletions":4,"non_countable":0}`, `{"category":"merge","additions":3,"deletions":true,"non_countable":0}`, "totals[1].deletions", "invalid_uint"},
		{"entries indexes divergent second item", `{"path_b64":"second","category":"merge","measurement":{"kind":"non_countable"}}`, `{"path_b64":"second","category":"merge","measurement":{"kind":"non_countable","additions":null}}`, "entries[1].measurement.additions", "unknown_field"},
		{"observations indexes divergent second item", `{"kind":"line_threshold","category":"merge","limit":0,"available":false,"exceeded":false}`, `{"kind":"line_threshold","category":"merge","limit":true,"available":false,"exceeded":false}`, "observations[1].limit", "invalid_uint"},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertResultRootShape(t, resultRootReplace(t, valid, test.old, test.new), test.field, test.code)
		})
	}
}

func TestResultV1RootEnvelopeShapeAcceptsStructuralSubtreeVariants(t *testing.T) {
	assertResultRootShape(t, resultRootEnvelopeFixture(), "", "")
}

func resultRootEnvelopeFixture() string {
	return `{"schema":"arbitrary schema","accounting_policy_sha256":"release","inventory_sha256":"merge","revisions":{"base":"release","head":"merge"},"entries":[{"path_b64":"first","category":"release","measurement":{"kind":"countable","additions":1,"deletions":2}},{"path_b64":"second","category":"merge","measurement":{"kind":"non_countable"}}],"totals":[{"category":"release","additions":1,"deletions":2,"non_countable":0},{"category":"merge","additions":3,"deletions":4,"non_countable":0}],"observations":[{"kind":"line_threshold","category":"release","limit":1,"available":true,"actual":2,"exceeded":false},{"kind":"line_threshold","category":"merge","limit":0,"available":false,"exceeded":false},{"kind":"ratio","category":"release","reference":"merge","limit":"arbitrary","available":true,"numerator":1,"denominator":0,"exceeded":true},{"kind":"ratio","category":"merge","reference":"release","limit":"other","available":false,"exceeded":false}]}`
}

func resultRootReplace(t *testing.T, raw, old, new string) string {
	t.Helper()
	updated := strings.Replace(raw, old, new, 1)
	if updated == raw {
		t.Fatalf("root fixture mutation %q was not applied", old)
	}
	return updated
}
