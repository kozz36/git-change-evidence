package changeevidence

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestValidateResultAntecedentsRequiresExactPolicyInventoryChain(t *testing.T) {
	policy, inventory := resultAntecedentFixture(t)
	view, policyDigest, inventoryDigest, err := validateResultAntecedents(policy, inventory)
	if err != nil || policyDigest != policy.Digest() || inventoryDigest != inventory.Digest() {
		t.Fatalf("validated antecedents = (%#v, %q, %q, %v)", view, policyDigest, inventoryDigest, err)
	}
	assertResultPolicyGlob(t, view, "src/**")
	callerCachedPolicy := policy.CanonicalBytes()
	callerCachedInventory := inventory.CanonicalBytes()
	otherPolicy := newPolicy(t, PolicyInput{Categories: []CategoryInput{{Name: "source", PathGlobs: [][]byte{[]byte("src/**")}, LineThreshold: 1}}, Default: "source"})
	otherInventory, err := NewInventoryV1(otherPolicy, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name      string
		policy    PolicyDocumentV1
		inventory InventoryDocumentV1
		field     string
	}{
		{"zero policy", PolicyDocumentV1{}, inventory, "accounting_policy"},
		{"policy digest only", PolicyDocumentV1{digest: policy.Digest()}, inventory, "accounting_policy"},
		{"policy stale digest", PolicyDocumentV1{canonical: policy.CanonicalBytes(), digest: DocumentDigest(strings.Repeat("0", 64)), view: PolicyInput(policy.View())}, inventory, "accounting_policy"},
		{"policy substituted bytes retained identity", PolicyDocumentV1{canonical: otherPolicy.CanonicalBytes(), digest: policy.Digest(), view: PolicyInput(policy.View())}, inventory, "accounting_policy"},
		{"policy noncanonical", PolicyDocumentV1{canonical: append(policy.CanonicalBytes(), ' '), digest: policy.Digest(), view: PolicyInput(policy.View())}, inventory, "accounting_policy"},
		{"zero inventory", policy, InventoryDocumentV1{}, "inventory"},
		{"inventory digest only", policy, InventoryDocumentV1{digest: inventory.Digest()}, "inventory"},
		{"valid policy then late inventory stale digest", policy, InventoryDocumentV1{canonical: inventory.CanonicalBytes(), digest: DocumentDigest(strings.Repeat("0", 64)), accountingPolicyDigest: policy.Digest()}, "inventory"},
		{"inventory noncanonical", policy, InventoryDocumentV1{canonical: append(inventory.CanonicalBytes(), ' '), digest: inventory.Digest(), accountingPolicyDigest: policy.Digest()}, "inventory"},
		{"inventory substituted bytes retained identity", policy, InventoryDocumentV1{canonical: otherInventory.CanonicalBytes(), digest: inventory.Digest(), accountingPolicyDigest: policy.Digest()}, "inventory"},
		{"supplied inventory policy link differs after strict decode", policy, InventoryDocumentV1{canonical: inventory.CanonicalBytes(), digest: inventory.Digest(), accountingPolicyDigest: otherPolicy.Digest()}, "inventory"},
		{"valid inventory for another policy", policy, otherInventory, "inventory"},
	} {
		t.Run(test.name, func(t *testing.T) {
			gotView, gotPolicy, gotInventory, err := validateResultAntecedents(test.policy, test.inventory)
			assertResultAntecedentValidationError(t, err, test.field, "invalid_document")
			if !reflect.DeepEqual(gotView, PolicyView{}) || gotPolicy != "" || gotInventory != "" {
				t.Fatalf("failure returned partial validation = (%#v, %q, %q)", gotView, gotPolicy, gotInventory)
			}
		})
	}

	callerPolicyView := policy.View()
	callerPolicyView.Categories[0].PathGlobs[0][0] = '!'
	view.Categories[0].PathGlobs[0][0] = '?'
	policyAccessorBuffer := policy.CanonicalBytes()
	inventoryAccessorBuffer := inventory.CanonicalBytes()
	policyAccessorBuffer[0], inventoryAccessorBuffer[0] = '#', '$'
	if !bytes.Equal(policy.CanonicalBytes(), callerCachedPolicy) || !bytes.Equal(inventory.CanonicalBytes(), callerCachedInventory) || policy.Digest() != policyDigest || inventory.Digest() != inventoryDigest {
		t.Fatalf("canonical accessor mutation altered antecedent identity or body")
	}
	assertResultPolicyGlob(t, policy.View(), "src/**")

	freshPolicy, err := DecodeAccountingPolicyV1(callerCachedPolicy)
	if err != nil {
		t.Fatal(err)
	}
	freshInventory, err := DecodeInventoryV1(callerCachedInventory, freshPolicy)
	if err != nil {
		t.Fatal(err)
	}
	freshView, freshPolicyDigest, freshInventoryDigest, err := validateResultAntecedents(freshPolicy, freshInventory)
	if err != nil || freshPolicyDigest != policyDigest || freshInventoryDigest != inventoryDigest {
		t.Fatalf("caller-cached fresh chain = (%#v, %q, %q, %v)", freshView, freshPolicyDigest, freshInventoryDigest, err)
	}
	assertResultPolicyGlob(t, freshView, "src/**")

	laterView, laterPolicyDigest, laterInventoryDigest, err := validateResultAntecedents(policy, inventory)
	if err != nil || laterPolicyDigest != policyDigest || laterInventoryDigest != inventoryDigest {
		t.Fatalf("same-chain later validation = (%#v, %q, %q, %v)", laterView, laterPolicyDigest, laterInventoryDigest, err)
	}
	assertResultPolicyGlob(t, laterView, "src/**")
}

func resultAntecedentFixture(t *testing.T) (PolicyDocumentV1, InventoryDocumentV1) {
	t.Helper()
	policy := newPolicy(t, PolicyInput{Categories: []CategoryInput{{Name: "source", PathGlobs: [][]byte{[]byte("src/**")}}}, Default: "source"})
	inventory, err := NewInventoryV1(policy, []InventoryEntryInput{{Path: []byte("untracked"), Content: []byte("content")}})
	if err != nil {
		t.Fatal(err)
	}
	return policy, inventory
}

func assertResultAntecedentValidationError(t *testing.T, err error, field, code string) {
	t.Helper()
	var contract *ContractError
	if !errors.As(err, &contract) || contract.Field != field || contract.Code != code {
		t.Fatalf("error = %#v, want %s/%s", err, field, code)
	}
}

func assertResultPolicyGlob(t *testing.T, view PolicyView, want string) {
	t.Helper()
	if len(view.Categories) != 1 || len(view.Categories[0].PathGlobs) != 1 || !bytes.Equal(view.Categories[0].PathGlobs[0], []byte(want)) {
		t.Fatalf("policy view glob = %#v, want %q", view, want)
	}
}
