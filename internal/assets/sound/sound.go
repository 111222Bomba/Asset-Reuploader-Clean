// internal/assets/sound/sound.go
package sound

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/roblox"
	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/types" // Döngüyü kırmak için kullanılan yeni paket
)

// Reupload, Sound varlığını yeniden yükler
func Reupload(c *roblox.Client, r *types.RawRequest) error {
	// Sound yükleme API URL'i
	uploadURL := fmt.Sprintf("https://assetgame.roblox.com/asset/?id=%d", r.AssetID)
	
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// 1. Dosya yolunu Plugin'den alıp açıyoruz
	file, err := os.Open(r.ExportPath) 
	if err != nil {
		return fmt.Errorf("ses dosyası açılamadı: %w", err)
	}
	defer file.Close()
	
	// 2. Dosyayı forma ekliyoruz
	part, err := writer.CreateFormFile("fileForUpload", r.ExportPath) // Form alan adı
	if err != nil {
		return err
	}
	io.Copy(part, file)
	
	writer.Close() // Multipart formu kapat
	
	// 3. POST İsteğini oluştur
	req, _ := http.NewRequest("POST", uploadURL, body)
	req.Header.Set("Cookie", ".ROBLOSECURITY=" + c.Cookie)
	req.Header.Set("X-CSRF-TOKEN", c.GetToken()) // CSRF Token ekleniyor
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// 4. İsteği gönder
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 5. Yanıtı kontrol et
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Sound API hata kodu: %d. Yanıt: %s", resp.StatusCode, string(responseBody))
	}
	
	return nil
}
