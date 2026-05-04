package twigga

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// doRequest handles standard JSON REST calls
func (c *httpClient) doRequest(ctx context.Context, method, url string, body any) ([]byte, int, error) {
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
