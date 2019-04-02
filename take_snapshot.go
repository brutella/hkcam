package hkcam

import (
	"github.com/brutella/hc/characteristic"
)

const TypeTakeSnapshot = "E8AEE54F-6E4B-46D8-85B2-FECE188FDB08"

// TakeSnapshot is used to take a snapshot.
// After writing `true` to this characteristic,
// a snapshot is taked and persisted on disk.
type TakeSnapshot struct {
	*characteristic.Bool
}

func NewTakeSnapshot() *TakeSnapshot {
	b := characteristic.NewBool(TypeTakeSnapshot)
	b.Description = "Take Snapshot"
	b.Perms = []string{characteristic.PermWrite}

	return &TakeSnapshot{b}
}
