// internal/assets/animation/animation.go
package animation

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/roblox"
	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/types" // Döngüyü kırmak için
)

// Reupload, Animation varlığını yeniden yükler
func Reupload(c *roblox.Client, r *types.RawRequest) error {
	// Animasyon yükleme API URL'i
	uploadURL := fmt.Sprintf("https://data.roblox.com/ide/publish/uploadanimation?assetId=%d", r.AssetID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Dosya yolu Plugin'den geliyor
	file, err := os.Open(r.ExportPath) 
	if err != nil {
		return fmt.Errorf("animasyon dosyası açılamadı: %w", err)
	}
	defer file.Close()

	// Dosyayı forma ekle: Form alanı "asset"
	part, err := writer.CreateFormFile("asset", r.ExportPath) 
	if err != nil {
		return err
	}
	io.Copy(part, file)

	writer.Close()

	// İsteği oluştur
	req, _ := http.NewRequest("POST", uploadURL, body)
	req.Header.Set("Cookie", ".ROBLOSECURITY=" + c.Cookie)
	req.Header.Set("X-CSRF-TOKEN", c.GetToken()) // CSRF Token
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// İsteği gönder (DÜZELTME UYGULANDI: c.HTTPClient)
	resp, err := c.HTTPClient.Do(req) 
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Animasyon API hata kodu: %d. Yanıt: %s", resp.StatusCode, string(responseBody))
	}
	
	return nil
}
