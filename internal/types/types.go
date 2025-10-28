// internal/types/types.go
package types

// RawRequest: Plugin'den gelen ham JSON verisi yapısı (Artık Router paketinden bağımsız)
type RawRequest struct {
	UniverseID    int64  `json:"universeId"`
	PlaceID       int64  `json:"placeId"`
	AssetID       int64  `json:"assetId"`
	AssetType     string `json:"assetType"` // "Sound", "Animation", etc.
	ExportPath    string `json:"exportPath"` // Plugin'in gönderdiği dosya yolu
}
