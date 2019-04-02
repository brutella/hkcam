package hkcam

import (
	"github.com/brutella/hc/characteristic"
)

const TypeGetAsset = "6A6C39F5-67F0-4BE1-BA9D-E56BD27C9606"

type GetAsset struct {
	*characteristic.Bytes
}

func NewGetAsset() *GetAsset {
	b := characteristic.NewBytes(TypeGetAsset)
	b.Perms = []string{characteristic.PermRead, characteristic.PermWrite}
	b.Value = []byte{}

	return &GetAsset{b}
}
