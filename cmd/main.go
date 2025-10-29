// cmd/main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
    "os" // Sadece örnek olması için
    
	// Proje paketlerinizi içe aktarın
	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/roblox"
	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/assets"
)

const Port = "38073"

// Payload yapısı (Lua script'inden gelen JSON)
type ReuploadPayload struct {
	UniverseId int64  `json:"universeId"`
	PlaceId    int64  `json:"placeId"`
	AssetId    int64  `json:"assetId"`
	AssetType  string `json:"assetType"`
	ExportPath string `json:"exportPath"`
}

// Global Roblox istemcisi
var robloxClient *roblox.Client

// --- HTTP İŞLEYİCİSİ (Handler) ---

func handleReupload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    
    // Yalnızca POST isteklerini kabul et
    if r.Method != http.MethodPost {
        http.Error(w, `{"success":false,"error":"Sadece POST istekleri kabul edilir."}`, http.StatusMethodNotAllowed)
        return
    }

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("KRİTİK SUNUCU HATASI: İstek gövdesi okunamadı: %s\n", err.Error()) 
		http.Error(w, `{"success":false,"error":"Invalid request body"}`, http.StatusInternalServerError)
		return
	}

	var payload ReuploadPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		fmt.Printf("KRİTİK SUNUCU HATASI: JSON çözümleme hatası: %s\n", err.Error()) 
		http.Error(w, `{"success":false,"error":"Invalid JSON request"}`, http.StatusBadRequest)
		return
	}
    
    // Payload verilerini logla
    log.Printf("Gelen istek: Varlık ID: %d, Tip: %s, Yol: %s", payload.AssetId, payload.AssetType, payload.ExportPath)

	// --- Varlık Yükleme İşlemi ---
	
	// Varsayım: robloxClient global olarak ayarlandı
	if robloxClient == nil {
		err = fmt.Errorf("Roblox istemcisi başlatılmamış.")
	} else {
        // Gerçek yükleme işlevini çağır
        // Bu fonksiyonun, varlık yükleme ve çerez yönetimini yaptığı varsayılır.
		err = asset.ReuploadAsset(robloxClient, payload.AssetId, payload.AssetType, payload.ExportPath, payload.UniverseId) 
	}
    
	if err != nil {
		// KRİTİK DÜZELTME: Hatanın detayını terminale basıyoruz
		fmt.Printf("KRİTİK SUNUCU HATASI (HTTP 500): Yükleme başarısız oldu: %s\n", err.Error()) 
		
        // Lua'ya hata mesajını gönder
        errorResponse := fmt.Sprintf(`{"success":false,"error":"Yükleme Başarısız: %s"}`, err.Error())
		http.Error(w, errorResponse, http.StatusInternalServerError)
		return
	}

	// Başarılı yanıt
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf(`{"success":true,"message":"%s başarıyla yeniden yüklendi (ID: %d)"}`, payload.AssetType, payload.AssetId)
	w.Write([]byte(response))
}

// --- ANA FONKSİYON ---

func main() {
    // 1. Cookie okuma ve istemci oluşturma (Bu kısım sizde zaten olmalı)
    // Varsayım: .ROBLOSECURITY çerezi "cookie.txt" dosyasından okunuyor.
    cookie, err := os.ReadFile("cookie.txt")
	if err != nil {
		log.Fatalf("HATA: Cookie dosyası okunamadı: %v", err)
	}
    
	robloxClient, err = roblox.NewClient(string(cookie))
	if err != nil {
		log.Fatalf("HATA: Roblox istemcisi oluşturulamadı veya CSRF token çekilemedi. Çereziniz geçersiz olabilir: %v", err)
	}

    // 2. HTTP sunucusunu başlatma
	http.HandleFunc("/", handleReupload)

	log.Printf("Plugin isteklerini http://localhost:%s adresinde dinliyor...", Port)
	log.Printf("CSRF Token: %s", robloxClient.GetToken())
    
	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatalf("HTTP sunucusu başlatılamadı: %v", err)
	}
}

