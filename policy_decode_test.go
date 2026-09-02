package changeevidence

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func policyRaw(categories, ratios string) string {
	return `{"schema":"git-change-evidence.accounting-policy/v1","categories":` + categories + `,"default":"later","ratios":` + ratios + `}`
}
func TestDecodeAccountingPolicyV1FailsClosed(t *testing.T) {
	raw := string(newPolicy(t, policyInput()).CanonicalBytes())
	got, err := DecodeAccountingPolicyV1([]byte(raw))
	if err != nil || !bytes.Equal(got.CanonicalBytes(), []byte(raw)) {
		t.Fatalf("DecodeAccountingPolicyV1() = (%q, %v)", got.CanonicalBytes(), err)
	}
	view := got.View()
	view.Categories[0].PathGlobs[0][0] = 'x'
	if got.View().Categories[0].PathGlobs[0][0] != 'a' {
		t.Fatal("decoded view aliases nested bytes")
	}
	duplicateRatio := strings.Replace(raw, "]}\n", ",{\"category\":\"later\",\"reference\":\"first\",\"maximum\":\"0.10000000000000001\"}]}\n", 1)
	category := `[{"name":"later","path_globs_b64":[],"line_threshold":0}]`
	for _, value := range []string{
		"{", "[]", "null", raw + "{}", " " + raw, strings.Replace(raw, `{"schema"`, `{"unknown":true,"schema"`, 1), strings.Replace(raw, `,"default":"later"`, "", 1),
		policyRaw("null", "[]"), policyRaw("{}", "[]"), policyRaw("[null]", "[]"), policyRaw(`[{"name":0,"path_globs_b64":[],"line_threshold":0}]`, "[]"),
		policyRaw(`[{"name":"later","path_globs_b64":null,"line_threshold":0}]`, "[]"), policyRaw(`[{"name":"later","path_globs_b64":{},"line_threshold":0}]`, "[]"), policyRaw(`[{"name":"later","path_globs_b64":[null],"line_threshold":0}]`, "[]"),
		policyRaw(`[{"name":"later","path_globs_b64":[],"line_threshold":18446744073709551616}]`, "[]"), policyRaw(category, "null"), policyRaw(category, "{}"), policyRaw(category, "[null]"),
		policyRaw(`[{"name":"later","name":"other","path_globs_b64":[],"line_threshold":0}]`, "[]"), policyRaw(`[{"name":"later","unknown":true,"path_globs_b64":[],"line_threshold":0}]`, "[]"), policyRaw(`[{"path_globs_b64":[],"line_threshold":0}]`, "[]"),
		policyRaw(category, `[{"category":"later","category":"other","reference":"later","maximum":"1"}]`), policyRaw(category, `[{"category":"later","unknown":true,"reference":"later","maximum":"1"}]`), policyRaw(category, `[{"category":"later","reference":"later"}]`),
		strings.Replace(raw, `,"categories"`, `,"sche\u006da":"git-change-evidence.accounting-policy/v1","categories"`, 1), strings.Replace(raw, `"YQ=="`, `"YQ==\\n"`, 1), strings.Replace(raw, `"YQ=="`, `"YQ"`, 1), strings.Replace(raw, `"YQ==","/w=="`, `"/w==","YQ=="`, 1), strings.Replace(raw, `"YQ==","/w=="`, `"YQ==","YQ=="`, 1),
		strings.Replace(strings.Replace(raw, `{"schema"`, `{"default":"later","schema"`, 1), `,"default":"later","ratios"`, `,"ratios"`, 1), strings.Replace(raw, `,"categories"`, `, "categories"`, 1), strings.ReplaceAll(raw, "later", `\u006cater`), duplicateRatio,
	} {
		if _, err := DecodeAccountingPolicyV1([]byte(value)); err == nil {
			t.Errorf("noncanonical or invalid document accepted: %q", value)
		}
	}
	assertAuthority := func(value string) {
		_, err := DecodeAccountingPolicyV1([]byte(value))
		var contract *ContractError
		if !errors.As(err, &contract) || contract.Code != "authority_field" {
			t.Errorf("authority error = %#v", err)
		}
	}
	for _, field := range []string{"approval", "rejection", "blocking", "gating", "merge", "deploy", "release", "delivery"} {
		assertAuthority(strings.Replace(raw, `{"schema"`, `{"`+field+`":true,"schema"`, 1))
	}
	assertAuthority(strings.Replace(raw, `{"name":"later"`, `{"rejection":true,"name":"later"`, 1))
	assertAuthority(strings.Replace(raw, `{"category":"later"`, `{"delivery":true,"category":"later"`, 1))
}
