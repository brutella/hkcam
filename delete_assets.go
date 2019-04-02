package hkcam

import (
	"github.com/brutella/hc/characteristic"
)

// TypeDeleteAssets is the uuid of the DeleteAssets characteristic
const TypeDeleteAssets = "3982EB69-1ECE-463E-96C6-E5A7DF2FA1CD"

// DeleteAssets is used to handle request to delete assets.
// A valid JSON looks like this. `{"ids":["1.jpg"]}`
// Reading the value of this characteristic is discouraged.
type DeleteAssets struct {
	*characteristic.Bytes
}

func NewDeleteAssets() *DeleteAssets {
	b := characteristic.NewBytes(TypeDeleteAssets)
	b.Perms = []string{characteristic.PermRead, characteristic.PermWrite}
	b.Value = []byte{}

	return &DeleteAssets{b}
}

type DeleteAssetsRequest struct {
	IDs []string `json:"ids"`
}
