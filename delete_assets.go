package hkcam

import (
	"github.com/brutella/hc/characteristic"
)

const TypeDeleteAssets = "3982EB69-1ECE-463E-96C6-E5A7DF2FA1CD" // "DDC309E5-A1B1-470F-AB0C-235202C531BA"

type DeleteAssets struct {
	*characteristic.Bytes
}

func NewDeleteAssets() *DeleteAssets {
	b := characteristic.NewBytes(TypeDeleteAssets)
	b.Perms = []string{characteristic.PermRead, characteristic.PermWrite}
	b.Value = []byte{}

	return &DeleteAssets{b}
}
