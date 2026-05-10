package twigga

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// doRequest handles standard JSON REST calls
func (c *httpClient) doRequest(ctx context.Context, method, url string, body any) ([]byte, int, error) {
	// ensure we have a valid token before proceeding!
	if err := c.ensureToken(ctx); err != nil {
		return nil, 0, fmt.Errorf("SDK authentication error: %v", err)
	}

	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("BONGO-TOKEN", c.token) // Updated to use the new struct property
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return bodyBytes, resp.StatusCode, nil
}

// doMultipartRequest handles file uploads (like zip files for Functions)
func (c *httpClient) doMultipartRequest(ctx context.Context, method, url string, body io.Reader, contentType string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", contentType)
	if c.token != "" {
		req.Header.Set("BONGO-TOKEN", c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	return bodyBytes, resp.StatusCode, err
}

func (c *httpClient) openWS(rawurl string) (*websocket.Conn, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	if c.token != "" {
		header.Set("BONGO-TOKEN", c.token)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	return conn, err
}

// ensureToken checks if the SDK's internal token is missing or expired.
// If it is, it silently calls BongoCloud to get a new one before continuing.
func (c *httpClient) ensureToken(ctx context.Context) error {
	if c.appId == "" || c.appSecret == "" {
		return nil // No credentials provided, operating unauthenticated
	}

	c.tokenMutex.RLock()
	if c.token != "" && time.Now().Add(5*time.Minute).Before(c.expiresAt) {
		c.tokenMutex.RUnlock()
		return nil
	}
	c.tokenMutex.RUnlock()

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if c.token != "" && time.Now().Add(5*time.Minute).Before(c.expiresAt) {
		return nil
	}

	url := fmt.Sprintf("%s/oauth/token", c.accountURL)
	reqBody := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.appId,
		"client_secret": c.appSecret,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("network error fetching SDK token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth server rejected SDK credentials (status %d)", resp.StatusCode)
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"` // 86400 (24h)
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &res)

	c.token = res.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)

	return nil
}
