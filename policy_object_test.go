package changeevidence

import "testing"

func TestPolicyObject(t *testing.T) {
	tests := []struct {
		name, input, field, value, errorField, errorCode string
	}{
		{"allowed fields", `{"name":"source","maximum":"0.5"}`, "policy", `"source"`, "", ""},
		{"unknown field", `{"name":"source","unexpected":true}`, "policy", "", "policy.unexpected", "unknown_field"},
		{"authority field", `{"name":"source","release":true}`, "policy", "", "policy.release", "authority_field"},
		{"duplicate field", `{"name":"source","name":"other"}`, "policy", "", "policy.name", "duplicate_field"},
		{"not object", `[]`, "policy", "", "policy", "invalid_object"},
		{"escaped duplicate", `{"na\u006de":"source","name":"other"}`, "policy", "", "policy.name", "duplicate_field"},
		{"normalized authority", `{"delivery_authority":true}`, "policy", "", "policy.delivery_authority", "authority_field"},
		{"malformed", `{"name":`, "policy", "", "policy", "invalid_object"},
		{"trailing document", `{"name":"source"}{}`, "policy", "", "policy", "invalid_object"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			object, err := decodePolicyObject([]byte(tt.input), tt.field, "name", "maximum")
			if tt.errorCode == "" {
				if err != nil {
					t.Fatal(err)
				}
				if got := string(object["name"]); got != tt.value {
					t.Fatalf("object name = %q, want %q", got, tt.value)
				}
				return
			}
			contract, ok := err.(*ContractError)
			if !ok {
				t.Fatalf("error = %T %v, want *ContractError", err, err)
			}
			if contract.Field != tt.errorField || contract.Code != tt.errorCode {
				t.Fatalf("error = %#v, want field %q code %q", contract, tt.errorField, tt.errorCode)
			}
		})
	}
}
