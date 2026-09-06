package changeevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

const (
	resultUnit5Base = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	resultUnit5Head = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

func TestNewAccountingResultV1ConstructsExclusiveCanonicalEvidence(t *testing.T) {
	policy, inventory := resultUnit5Antecedents(t, resultUnit5Policy())
	reversed := []CommittedChange{
		{Path: string([]byte{0xff}), Lines: CommittedLineCounts{Additions: 99}},
		{Path: "a0", Lines: CommittedLineCounts{Additions: 5, Countable: true}},
		{Path: string(append([]byte("a/"), 0xff)), Lines: CommittedLineCounts{Additions: 4, Countable: true}},
		{Path: "a", Lines: CommittedLineCounts{Additions: 2, Deletions: 1, Countable: true}},
	}
	got, err := NewAccountingResultV1(policy, inventory, resultUnit5Snapshot(reversed))
	if err != nil {
		t.Fatal(err)
	}
	wantEntries := []ResultEntryV1{
		{Path: []byte("a"), Category: "other", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: 2, Deletions: 1}},
		{Path: append([]byte("a/"), 0xff), Category: "raw", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: 4}},
		{Path: []byte("a0"), Category: "other", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1, Additions: 5}},
		{Path: []byte{0xff}, Category: "other", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryNonCountableV1}},
	}
	wantTotals := []ResultCategoryTotalV1{{Category: "raw", Additions: 4}, {Category: "source"}, {Category: "other", Additions: 7, Deletions: 1, NonCountable: 1}}
	wantObservations := []ResultObservationV1{
		{Kind: ResultLineThresholdObservationV1, Category: "raw", LineThresholdLimit: 3, Available: true, Actual: 4, Exceeded: true},
		{Kind: ResultLineThresholdObservationV1, Category: "other", LineThresholdLimit: 3},
		{Kind: ResultRatioObservationV1, Category: "raw", Reference: "source", RatioLimit: "0.5"},
	}
	if !reflect.DeepEqual(got.Entries(), wantEntries) || !reflect.DeepEqual(got.Totals(), wantTotals) || !reflect.DeepEqual(got.Observations(), wantObservations) {
		t.Fatalf("result = (%#v, %#v, %#v)", got.Entries(), got.Totals(), got.Observations())
	}
	wantCanonical := []byte(fmt.Sprintf(`{"schema":"git-change-evidence.accounting-result/v1","accounting_policy_sha256":"%s","inventory_sha256":"%s","revisions":{"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},"entries":[{"path_b64":"YQ==","category":"other","measurement":{"kind":"countable","additions":2,"deletions":1}},{"path_b64":"YS//","category":"raw","measurement":{"kind":"countable","additions":4,"deletions":0}},{"path_b64":"YTA=","category":"other","measurement":{"kind":"countable","additions":5,"deletions":0}},{"path_b64":"/w==","category":"other","measurement":{"kind":"non_countable"}}],"totals":[{"category":"raw","additions":4,"deletions":0,"non_countable":0},{"category":"source","additions":0,"deletions":0,"non_countable":0},{"category":"other","additions":7,"deletions":1,"non_countable":1}],"observations":[{"kind":"line_threshold","category":"raw","limit":3,"available":true,"actual":4,"exceeded":true},{"kind":"line_threshold","category":"other","limit":3,"available":false,"exceeded":false},{"kind":"ratio","category":"raw","reference":"source","limit":"0.5","available":false,"exceeded":false}]}`, policy.Digest(), inventory.Digest()) + "\n")
	wantHash := sha256.Sum256(wantCanonical)
	if got.AccountingPolicyDigest() != policy.Digest() || got.InventoryDigest() != inventory.Digest() || got.Revisions() != (RevisionIdentity{Base: resultUnit5Base, Head: resultUnit5Head}) {
		t.Fatalf("provenance = (%q, %q, %#v)", got.AccountingPolicyDigest(), got.InventoryDigest(), got.Revisions())
	}
	if !bytes.Equal(got.CanonicalBytes(), wantCanonical) || got.Digest() != DocumentDigest(hex.EncodeToString(wantHash[:])) {
		t.Fatalf("canonical identity = (%q, %q)", got.CanonicalBytes(), got.Digest())
	}
	ordered, err := NewAccountingResultV1(policy, inventory, resultUnit5Snapshot(wantCommittedChanges(wantEntries)))
	if err != nil || !bytes.Equal(got.CanonicalBytes(), ordered.CanonicalBytes()) || got.Digest() != ordered.Digest() {
		t.Fatalf("ordered result = (%#v, %v)", ordered, err)
	}
}

func resultUnit5Policy() PolicyInput {
	return PolicyInput{
		Categories: []CategoryInput{{Name: "raw", PathGlobs: [][]byte{[]byte("a/**")}, LineThreshold: 3}, {Name: "source", PathGlobs: [][]byte{[]byte("a/**")}}, {Name: "other", LineThreshold: 3}},
		Default:    "other",
		Ratios:     []RatioInput{{Category: "raw", Reference: "source", Maximum: "0.5"}},
	}
}

func resultUnit5Antecedents(t *testing.T, input PolicyInput) (PolicyDocumentV1, InventoryDocumentV1) {
	t.Helper()
	policy := newPolicy(t, input)
	inventory, err := NewInventoryV1(policy, nil)
	if err != nil {
		t.Fatal(err)
	}
	return policy, inventory
}

func resultUnit5Snapshot(entries []CommittedChange) CommittedSnapshot {
	return NewCommittedSnapshot(resultUnit5Base, resultUnit5Head, entries)
}

func wantCommittedChanges(entries []ResultEntryV1) []CommittedChange {
	changes := make([]CommittedChange, len(entries))
	for i, entry := range entries {
		changes[i] = CommittedChange{Path: string(entry.Path), Lines: CommittedLineCounts{Additions: entry.Measurement.Additions, Deletions: entry.Measurement.Deletions, Countable: entry.Measurement.Kind == ResultEntryCountableV1}}
	}
	return changes
}
