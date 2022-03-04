package hkcam

import (
	"github.com/brutella/hap/characteristic"
)

// TypeAssets is the uuid of the Assets characteristic
const TypeAssets = "ACD9DFE7-948D-43D0-A205-D2F6F368541D"

// Assets contains a list of assets encoded as JSON.
// A valid JSON looks like this. `{"assets":[{"id":"1.jpg", "date":"2019-04-01T10:00:00+00:00"}]}`
// Writing to this characteristic is discouraged.
type Assets struct {
	*characteristic.Bytes
}

func NewAssets() *Assets {
	b := characteristic.NewBytes(TypeAssets)
	b.Permissions = []string{characteristic.PermissionRead, characteristic.PermissionEvents}

	b.SetValue([]byte{})

	return &Assets{b}
}

type AssetsMetadataResponse struct {
	Assets []CameraAssetMetadata `json:"assets"`
}

type CameraAssetMetadata struct {
	ID   string `json:"id"`
	Date string `json:"date"`
}
