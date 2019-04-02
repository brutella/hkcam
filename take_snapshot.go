package hkcam

import (
	"github.com/brutella/hc/characteristic"
)

const TypeTakeSnapshot = "E8AEE54F-6E4B-46D8-85B2-FECE188FDB08"

type TakeSnapshot struct {
	*characteristic.Bool
}

func NewTakeSnapshot() *TakeSnapshot {
	b := characteristic.NewBool(TypeTakeSnapshot)
	b.Description = "Take Snapshot"
	b.Perms = []string{characteristic.PermWrite}

	return &TakeSnapshot{b}
}
