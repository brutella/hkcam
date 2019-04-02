package hkcam

import (
	"github.com/brutella/hc/characteristic"
)

const TypeTakeSnapshot = "E8AEE54F-6E4B-46D8-85B2-FECE188FDB08" // "7F8BA9BC-C3E0-4A75-8B24-81EB483F9C84"

type TakeSnapshot struct {
	*characteristic.Bool
}

func NewTakeSnapshot() *TakeSnapshot {
	b := characteristic.NewBool(TypeTakeSnapshot)
	b.Description = "Take Snapshot"
	b.Perms = []string{characteristic.PermWrite}

	return &TakeSnapshot{b}
}
