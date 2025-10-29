// internal/router/router.go (GÜNCELLENMİŞ)
package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/pkg/animation"
	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/pkg/sound"
	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/roblox"
	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/types" // YENİ IMPORT
)

// RawRequest yapısı silindi. Artık types.RawRequest kullanılıyor.

type Router struct {
	Client *roblox.Client
}

func NewRouter(c *roblox.Client) *Router {
	return &Router{Client: c}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	var rawReq types.RawRequest // types.RawRequest kullanıldı
	if err := json.NewDecoder(req.Body).Decode(&rawReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Invalid JSON request"})
		return
	}

	var err error
	switch rawReq.AssetType {
	case "Sound":
		err = sound.Reupload(r.Client, &rawReq)
	case "Animation":
		err = animation.Reupload(r.Client, &rawReq)
	default:
		err = fmt.Errorf("desteklenmeyen varlık türü: %s", rawReq.AssetType)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": fmt.Sprintf("%s ID %d başarıyla yüklendi.", rawReq.AssetType, rawReq.AssetID)})
}

