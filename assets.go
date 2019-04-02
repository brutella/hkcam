package hkcam

import (
	"github.com/brutella/hc/characteristic"
)

const TypeAssets = "ACD9DFE7-948D-43D0-A205-D2F6F368541D" // "40F19475-404C-418F-B20C-0E3CF0389882"

type Assets struct {
	*characteristic.Bytes
}

func NewAssets() *Assets {
	b := characteristic.NewBytes(TypeAssets)
	b.Perms = characteristic.PermsAll()
	b.Value = []byte{}

	return &Assets{b}
}
