package hkcam

import (
	"github.com/brutella/hc/characteristic"
)

const TypeDeleteAssets = "3982EB69-1ECE-463E-96C6-E5A7DF2FA1CD"

type DeleteAssets struct {
	*characteristic.Bytes
}

func NewDeleteAssets() *DeleteAssets {
	b := characteristic.NewBytes(TypeDeleteAssets)
	b.Perms = []string{characteristic.PermRead, characteristic.PermWrite}
	b.Value = []byte{}

	return &DeleteAssets{b}
}
