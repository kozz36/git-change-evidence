package changeevidence

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

const (
	wirePolicyDigest    = DocumentDigest("1111111111111111111111111111111111111111111111111111111111111111")
	wireInventoryDigest = DocumentDigest("2222222222222222222222222222222222222222222222222222222222222222")
)

func TestResultV1CanonicalWireHostilePathAndIdentity(t *testing.T) {
	document, err := wireResultDocument(
		[]ResultEntryV1{{
			Path:     []byte{0xff, 0x0a, 'f', 'i', 'l', 'e'},
			Category: "source",
			Measurement: ResultEntryMeasurementV1{
				Kind: ResultEntryCountableV1, Additions: 2, Deletions: 1,
			},
		}},
		[]ResultCategoryTotalV1{{Category: "source", Additions: 2, Deletions: 1}},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(`{"schema":"git-change-evidence.accounting-result/v1","accounting_policy_sha256":"` + string(wirePolicyDigest) + `","inventory_sha256":"` + string(wireInventoryDigest) + `","revisions":{"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},"entries":[{"path_b64":"/wpmaWxl","category":"source","measurement":{"kind":"countable","additions":2,"deletions":1}}],"totals":[{"category":"source","additions":2,"deletions":1,"non_countable":0}],"observations":[]}` + "\n")
	if got := document.CanonicalBytes(); string(got) != string(want) {
		t.Fatalf("canonical bytes = %q, want %q", got, want)
	}

	sum := sha256.Sum256(want)
	wantDigest := DocumentDigest(hex.EncodeToString(sum[:]))
	if got := document.Digest(); got != wantDigest {
		t.Fatalf("digest = %q, want independently computed %q", got, wantDigest)
	}
}

func TestResultV1CanonicalWireUsesEmptyArrays(t *testing.T) {
	document, err := wireResultDocument(
		nil,
		[]ResultCategoryTotalV1{{Category: "source"}},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(`{"schema":"git-change-evidence.accounting-result/v1","accounting_policy_sha256":"` + string(wirePolicyDigest) + `","inventory_sha256":"` + string(wireInventoryDigest) + `","revisions":{"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},"entries":[],"totals":[{"category":"source","additions":0,"deletions":0,"non_countable":0}],"observations":[]}` + "\n")
	if got := document.CanonicalBytes(); string(got) != string(want) {
		t.Fatalf("canonical bytes = %q, want %q", got, want)
	}
}

func wireResultDocument(entries []ResultEntryV1, totals []ResultCategoryTotalV1, observations []ResultObservationV1) (ResultDocumentV1, error) {
	return resultDocument(
		wirePolicyDigest,
		wireInventoryDigest,
		RevisionIdentity{Base: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Head: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
		entries,
		totals,
		observations,
	)
}
