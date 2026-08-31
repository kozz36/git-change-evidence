package changeevidence

type UntrackedRecord struct{ Path string }

type UntrackedInventory struct{ records []UntrackedRecord }

func NewUntrackedInventory(records []UntrackedRecord) UntrackedInventory {
	return UntrackedInventory{records: append([]UntrackedRecord(nil), records...)}
}

func (i UntrackedInventory) Records() []UntrackedRecord {
	return append([]UntrackedRecord(nil), i.records...)
}

type InventoryUnavailableCode string

const (
	InventoryUnsupportedPlatform InventoryUnavailableCode = "unsupported_platform"
	InventoryInvalidPath         InventoryUnavailableCode = "invalid_path"
	InventoryMissing             InventoryUnavailableCode = "missing_entry"
	InventorySymlink             InventoryUnavailableCode = "symlink_entry"
	InventorySpecial             InventoryUnavailableCode = "special_entry"
	InventoryRacing              InventoryUnavailableCode = "racing_entry"
	InventoryUnverifiable        InventoryUnavailableCode = "unverifiable_identity"
)

type InventoryUnavailableError struct{ Code InventoryUnavailableCode }

func (e *InventoryUnavailableError) Error() string {
	return "inventory unavailable: " + string(e.Code)
}
