package changeevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestInventoryV1OwnershipEmptyAndLegacyCoexistence(t *testing.T) {
	policy := newPolicy(t, policyInput())
	empty, err := NewInventoryV1(policy, []InventoryEntryInput{})
	if err != nil || !bytes.Contains(empty.CanonicalBytes(), []byte(`"entries":[]`)) || empty.Digest() == "" || empty.AccountingPolicyDigest() == "" || empty.Entries() == nil || len(empty.Entries()) != 0 {
		t.Fatalf("verified empty inventory = (%#v, %v)", empty, err)
	}
	var zero PolicyDocumentV1
	failed, err := NewInventoryV1(zero, nil)
	if err == nil || failed.CanonicalBytes() != nil || failed.Digest() != "" || failed.AccountingPolicyDigest() != "" || failed.Entries() != nil {
		t.Fatalf("zero or failed inventory exposes evidence = (%#v, %v)", failed, err)
	}
	path, content := []byte("keep"), []byte("content")
	input := make([]InventoryEntryInput, 1, 2)
	input[0] = InventoryEntryInput{Path: path, Content: content}
	got, err := NewInventoryV1(policy, input)
	if err != nil {
		t.Fatal(err)
	}
	canonical, digest, link := got.CanonicalBytes(), got.Digest(), got.AccountingPolicyDigest()
	path[0], content[0], input[0] = 'x', 'x', InventoryEntryInput{Path: []byte("other"), Content: []byte("other")}
	input = append(input, InventoryEntryInput{Path: []byte("later"), Content: []byte("later")})
	returned, returnedBytes := got.Entries(), got.CanonicalBytes()
	returned[0], returned[0].Path[0], returnedBytes[0] = InventoryEntryV1{}, 'x', 'x'
	returned = append(returned, InventoryEntryV1{})
	later := got.Entries()
	contentSum := sha256.Sum256([]byte("content"))
	if !bytes.Equal(got.CanonicalBytes(), canonical) || got.Digest() != digest || got.AccountingPolicyDigest() != link || len(later) != 1 || !bytes.Equal(later[0].Path, []byte("keep")) || later[0].ContentSHA256 != Digest(hex.EncodeToString(contentSum[:])) || later[0].ByteLength != 7 {
		t.Fatal("mutable inventory state escaped")
	}
	legacy := NewUntrackedInventory([]UntrackedRecord{{Path: "legacy"}})
	records := legacy.Records()
	records[0].Path = "changed"
	if legacy.Records()[0].Path != "legacy" || (&InventoryUnavailableError{Code: InventoryMissing}).Error() != "inventory unavailable: missing_entry" {
		t.Fatal("legacy inventory behavior changed")
	}
}
