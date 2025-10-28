// internal/roblox/client.go
package roblox

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Plugin'in ihtiyaç duyduğu temel yapılar (develop paketinden temizlendi)
type Creator struct {
	Type string 
	TargetID int64
}

type AssetInfo struct {
	Name string
	ID   int64
	Creator Creator
}

type Client struct {
	Cookie string
	httpClient *http.Client
	token string // CSRF Token
	tokenMutex sync.RWMutex
}

func NewClient(cookie string) (*Client, error) {
	c := &Client{
		Cookie: strings.TrimSpace(cookie),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	
	if err := c.fetchCSRFToken(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) fetchCSRFToken() error {
	// Options isteği ile CSRF token'ı çekmek
	req, _ := http.NewRequest("OPTIONS", "https://auth.roblox.com/v1/logout", nil)
	req.Header.Set("Cookie", ".ROBLOSECURITY=" + c.Cookie)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	token := resp.Header.Get("x-csrf-token")
	if token == "" || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CSRF token alınamadı. Durum: %d", resp.StatusCode)
	}
	c.SetToken(token)
	return nil
}

func (c *Client) GetToken() string {
	c.tokenMutex.RLock()
	defer c.tokenMutex.RUnlock()
	return c.token
}

func (c *Client) SetToken(s string) {
	c.tokenMutex.Lock()
	c.token = s
	c.tokenMutex.Unlock()
}