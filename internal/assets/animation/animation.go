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
	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/router"
)

// Reupload, Animation varlığını yeniden yükler
func Reupload(c *roblox.Client, r *router.RawRequest) error {
	uploadURL := fmt.Sprintf("https://data.roblox.com/ide/publish/uploadanimation?assetId=%d", r.AssetID)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(r.ExportPath) 
	if err != nil {
		return fmt.Errorf("animasyon dosyası açılamadı: %w", err)
	}
	defer file.Close()

	part, _ := writer.CreateFormFile("asset", r.ExportPath) // Animasyon için form alanı "asset"
	io.Copy(part, file)

	writer.Close()

	// İsteği oluştur
	req, _ := http.NewRequest("POST", uploadURL, body)
	req.Header.Set("Cookie", ".ROBLOSECURITY=" + c.Cookie)
	req.Header.Set("X-CSRF-TOKEN", c.GetToken()) // KRİTİK: CSRF Token
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// İsteği gönder
	resp, err := c.httpClient.Do(req)
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