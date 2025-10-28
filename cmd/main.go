// cmd/main.go
package main

import (
	"bufio" // Kullanıcıdan girdi almak için
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/roblox"
	"github.com/111222Bomba/Asset-Reuploader-Clean/internal/router"
)

// Plugin'in dinlemesi gereken Port ve Cookie dosyası
const (
	Port = "38073" 
	CookieFile = "cookie.txt"
)

// getCookieFromUser, kullanıcıdan .ROBLOSECURITY çerezini alır.
func getCookieFromUser() string {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("Lütfen .ROBLOSECURITY çerezinizi yapıştırın (UYARI: Bu çerezi kimseyle paylaşmayın!): ")
		cookie, _ := reader.ReadString('\n')
		cookie = strings.TrimSpace(cookie)
		
		if cookie == "" {
			fmt.Println("Hata: Çerez boş olamaz. Tekrar deneyin.")
			continue
		}

		return cookie
	}
}


func main() {
	log.Println("Asset-Reuploader-Clean başlatılıyor...")
	
	var cookieStr string
	
	// 1. Cookie'yi Dosyadan Oku
	cookie, readErr := os.ReadFile(CookieFile)
	
	if readErr == nil {
		// Dosya mevcutsa ve okunduysa, kullan
		cookieStr = strings.TrimSpace(string(cookie))
	} else if os.IsNotExist(readErr) || readErr != nil { 
        // Dosya mevcut değilse veya başka bir okuma hatası varsa, kullanıcıdan girdi al
		fmt.Println("--- Cookie dosyası (cookie.txt) bulunamadı veya okunamadı. ---")
		cookieStr = getCookieFromUser()
		
		// Yeni Cookie'yi dosyaya yaz
		if err := os.WriteFile(CookieFile, []byte(cookieStr), 0644); err != nil {
			log.Printf("UYARI: Cookie dosyaya yazılamadı: %v", err)
		} else {
            fmt.Printf("Cookie başarıyla %s dosyasına kaydedildi. Sonraki çalıştırmalarda tekrar girmeyeceksiniz.\n", CookieFile)
        }
	}

	// 2. Roblox İstemcisini Oluştur (CSRF Token'ı Çeker)
	c, clientErr := roblox.NewClient(cookieStr)
	if clientErr != nil {
		// Eğer Cookie hatalıysa (CSRF çekilemiyorsa), programı sonlandır
		log.Fatalf("HATA: Roblox istemcisi oluşturulamadı veya CSRF token çekilemedi. Çereziniz geçersiz olabilir: %v", clientErr)
	}

	// 3. Plugin İsteklerini Dinlemeye Başla
	fmt.Printf("Plugin isteklerini http://localhost:%s adresinde dinliyor...\n", Port)
	
	r := router.NewRouter(c)
	http.Handle("/", r)

	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}