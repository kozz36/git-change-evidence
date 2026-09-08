package changeevidence

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewAccountingResultV1ClassifiesRenameByPrimaryPath(t *testing.T) {
	policy, inventory := resultUnit5Antecedents(t, PolicyInput{
		Categories: []CategoryInput{
			{Name: "docs", PathGlobs: [][]byte{[]byte("docs/**")}},
			{Name: "source", PathGlobs: [][]byte{[]byte("src/**")}},
			{Name: "other"},
		},
		Default: "other",
	})
	domain := resultAccountingDomain(policy.View())
	if old, current := domain.match([]byte("docs/old.txt")), domain.match([]byte("src/new.go")); old != 0 || current != 1 {
		t.Fatalf("fixture policy categories = (%d, %d), want (0, 1)", old, current)
	}

	result, err := NewAccountingResultV1(policy, inventory, resultUnit5Snapshot([]CommittedChange{{
		Path:         "src/new.go",
		PreviousPath: "docs/old.txt",
		Lines:        CommittedLineCounts{Additions: 4, Deletions: 2, Countable: true},
	}}))
	if err != nil {
		t.Fatal(err)
	}
	wantEntries := []ResultEntryV1{{
		Path:     []byte("src/new.go"),
		Category: "source",
		Measurement: ResultEntryMeasurementV1{
			Kind: ResultEntryCountableV1, Additions: 4, Deletions: 2,
		},
	}}
	wantTotals := []ResultCategoryTotalV1{
		{Category: "docs"},
		{Category: "source", Additions: 4, Deletions: 2},
		{Category: "other"},
	}
	if !reflect.DeepEqual(result.Entries(), wantEntries) || !reflect.DeepEqual(result.Totals(), wantTotals) {
		t.Fatalf("rename result = (%#v, %#v)", result.Entries(), result.Totals())
	}
}

func TestNewAccountingResultV1ReproducesIdentityFromIndependentAntecedents(t *testing.T) {
	policy, inventory := resultUnit5Antecedents(t, resultUnit5Policy())
	policyBytes, inventoryBytes := policy.CanonicalBytes(), inventory.CanonicalBytes()
	leftPolicy, err := DecodeAccountingPolicyV1(append([]byte(nil), policyBytes...))
	if err != nil {
		t.Fatal(err)
	}
	rightPolicy, err := DecodeAccountingPolicyV1(append([]byte(nil), policyBytes...))
	if err != nil {
		t.Fatal(err)
	}
	leftInventory, err := DecodeInventoryV1(append([]byte(nil), inventoryBytes...), leftPolicy)
	if err != nil {
		t.Fatal(err)
	}
	rightInventory, err := DecodeInventoryV1(append([]byte(nil), inventoryBytes...), rightPolicy)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(leftPolicy.CanonicalBytes(), policyBytes) || !bytes.Equal(leftPolicy.CanonicalBytes(), rightPolicy.CanonicalBytes()) || leftPolicy.Digest() != rightPolicy.Digest() || !bytes.Equal(leftInventory.CanonicalBytes(), inventoryBytes) || !bytes.Equal(leftInventory.CanonicalBytes(), rightInventory.CanonicalBytes()) || leftInventory.Digest() != rightInventory.Digest() || leftInventory.AccountingPolicyDigest() != leftPolicy.Digest() || rightInventory.AccountingPolicyDigest() != rightPolicy.Digest() {
		t.Fatal("independent antecedents differ or carry the wrong policy link")
	}

	left, err := NewAccountingResultV1(leftPolicy, leftInventory, resultUnit5Snapshot([]CommittedChange{{
		Path: "a/reproduced", Lines: CommittedLineCounts{Additions: 4, Deletions: 2, Countable: true},
	}}))
	if err != nil {
		t.Fatal(err)
	}
	right, err := NewAccountingResultV1(rightPolicy, rightInventory, resultUnit5Snapshot([]CommittedChange{{
		Path: "a/reproduced", Lines: CommittedLineCounts{Additions: 4, Deletions: 2, Countable: true},
	}}))
	if err != nil {
		t.Fatal(err)
	}
	if left.AccountingPolicyDigest() != leftPolicy.Digest() || left.InventoryDigest() != leftInventory.Digest() || right.AccountingPolicyDigest() != rightPolicy.Digest() || right.InventoryDigest() != rightInventory.Digest() || !bytes.Equal(left.CanonicalBytes(), right.CanonicalBytes()) || left.Digest() != right.Digest() {
		t.Fatal("equivalent constructions did not reproduce Result provenance identity")
	}
}
