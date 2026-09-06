package changeevidence

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func TestNewAccountingResultV1OwnsIngressAndEgress(t *testing.T) {
	policy, _ := resultUnit5Antecedents(t, resultUnit5Policy())
	inventory, err := NewInventoryV1(policy, []InventoryEntryInput{{Path: []byte("untracked"), Content: []byte("content")}})
	if err != nil {
		t.Fatal(err)
	}
	source := []CommittedChange{{Path: string(append([]byte("a/"), 0xff)), Lines: CommittedLineCounts{Additions: 4, Countable: true}}}
	snapshot := resultUnit5Snapshot(source)
	document, err := NewAccountingResultV1(policy, inventory, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	canonical := append([]byte(nil), document.CanonicalBytes()...)
	sourceEntries := document.Entries()
	entries := make([]ResultEntryV1, len(sourceEntries))
	for i, entry := range sourceEntries {
		entries[i] = entry
		entries[i].Path = append([]byte(nil), entry.Path...)
	}
	totals := append([]ResultCategoryTotalV1(nil), document.Totals()...)
	observations := append([]ResultObservationV1(nil), document.Observations()...)
	digest, policyDigest, inventoryDigest, revisions := document.Digest(), document.AccountingPolicyDigest(), document.InventoryDigest(), document.Revisions()
	source[0] = CommittedChange{Path: "changed", Lines: CommittedLineCounts{Additions: 99, Countable: true}}
	snapshot.Base, snapshot.Head = "changed", "changed"
	snapshotEntries := snapshot.Entries()
	snapshotEntries[0].Path, snapshotEntries[0].Lines.Additions = "changed", 99
	policyView, inventoryEntries := policy.View(), inventory.Entries()
	policyView.Categories[0].PathGlobs[0][0] = '!'
	if len(inventoryEntries) != 0 {
		inventoryEntries[0].Path[0] = '!'
	}
	returnedCanonical, returnedEntries, returnedTotals, returnedObservations := document.CanonicalBytes(), document.Entries(), document.Totals(), document.Observations()
	returnedCanonical[0] = '!'
	returnedEntries[0].Path[0] = '!'
	returnedEntries[0] = ResultEntryV1{}
	returnedTotals[0].Category = "changed"
	returnedObservations[0].Category = "changed"
	if !bytes.Equal(document.CanonicalBytes(), canonical) || document.Digest() != digest || document.AccountingPolicyDigest() != policyDigest || document.InventoryDigest() != inventoryDigest || document.Revisions() != revisions || !reflect.DeepEqual(document.Entries(), entries) || !reflect.DeepEqual(document.Totals(), totals) || !reflect.DeepEqual(document.Observations(), observations) {
		t.Fatal("ingress or egress mutation changed result")
	}
}

func TestNewAccountingResultV1ReturnsZeroOnFatalFailure(t *testing.T) {
	policy, inventory := resultUnit5Antecedents(t, resultUnit5Policy())
	_, otherInventory := resultUnit5Antecedents(t, PolicyInput{Categories: []CategoryInput{{Name: "alternate"}}, Default: "alternate"})
	ratioPolicy, ratioInventory := resultUnit5Antecedents(t, PolicyInput{Categories: []CategoryInput{{Name: "raw", PathGlobs: [][]byte{[]byte("a/**")}}, {Name: "source"}}, Default: "source", Ratios: []RatioInput{{Category: "raw", Reference: "source", Maximum: "0.5"}}})
	for _, test := range []struct {
		name      string
		policy    PolicyDocumentV1
		inventory InventoryDocumentV1
		snapshot  CommittedSnapshot
	}{
		{"antecedent", PolicyDocumentV1{}, inventory, resultUnit5Snapshot(nil)},
		{"zero inventory", policy, InventoryDocumentV1{}, resultUnit5Snapshot(nil)},
		{"inventory for another policy", policy, otherInventory, resultUnit5Snapshot(nil)},
		{"snapshot", policy, inventory, NewCommittedSnapshot("invalid", resultUnit5Head, nil)},
		{"late entry", policy, inventory, resultUnit5Snapshot([]CommittedChange{{Path: "a/x"}, {Path: "a/../bad"}})},
		{"checked total", policy, inventory, resultUnit5Snapshot([]CommittedChange{{Path: "a/x", Lines: CommittedLineCounts{Additions: math.MaxUint64, Countable: true}}, {Path: "a/y", Lines: CommittedLineCounts{Additions: 1, Countable: true}}})},
		{"threshold overflow", policy, inventory, resultUnit5Snapshot([]CommittedChange{{Path: "a/x", Lines: CommittedLineCounts{Additions: math.MaxUint64, Deletions: 1, Countable: true}}})},
		{"ratio overflow", ratioPolicy, ratioInventory, resultUnit5Snapshot([]CommittedChange{{Path: "a/x", Lines: CommittedLineCounts{Additions: math.MaxUint64, Deletions: 1, Countable: true}}, {Path: "note", Lines: CommittedLineCounts{Additions: 1, Countable: true}}})},
		{"ratio reference overflow after unavailable numerator", ratioPolicy, ratioInventory, resultUnit5Snapshot([]CommittedChange{{Path: "a/x", Lines: CommittedLineCounts{Countable: false}}, {Path: "note", Lines: CommittedLineCounts{Additions: math.MaxUint64, Deletions: 1, Countable: true}}})},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewAccountingResultV1(test.policy, test.inventory, test.snapshot)
			if err == nil {
				t.Fatal("constructor succeeded")
			}
			assertZeroResultDocument(t, got)
		})
	}
}

func assertZeroResultDocument(t *testing.T, got ResultDocumentV1) {
	t.Helper()
	if !reflect.DeepEqual(got, ResultDocumentV1{}) || got.CanonicalBytes() != nil || got.Digest() != "" || got.AccountingPolicyDigest() != "" || got.InventoryDigest() != "" || got.Revisions() != (RevisionIdentity{}) || got.Entries() != nil || got.Totals() != nil || got.Observations() != nil {
		t.Fatalf("nonzero result = %#v", got)
	}
}
