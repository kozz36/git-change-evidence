package changeevidence

const InventoryContractV1 ContractVersion = "git-change-evidence.inventory/v1"

type InventoryEntryInput struct {
	Path    []byte
	Content []byte
}

type InventoryEntryV1 struct {
	Path          []byte
	ContentSHA256 Digest
	ByteLength    uint64
}

type InventoryDocumentV1 struct {
	canonical              []byte
	digest                 DocumentDigest
	accountingPolicyDigest DocumentDigest
	entries                []InventoryEntryV1
}

func NewInventoryV1(policy PolicyDocumentV1, entries []InventoryEntryInput) (InventoryDocumentV1, error) {
	return buildInventoryV1(policy, entries)
}

func (i InventoryDocumentV1) CanonicalBytes() []byte { return append([]byte(nil), i.canonical...) }
func (i InventoryDocumentV1) Digest() DocumentDigest { return i.digest }
func (i InventoryDocumentV1) AccountingPolicyDigest() DocumentDigest {
	return i.accountingPolicyDigest
}
func (i InventoryDocumentV1) Entries() []InventoryEntryV1 {
	if i.canonical == nil {
		return nil
	}
	return cloneInventoryEntries(i.entries)
}
