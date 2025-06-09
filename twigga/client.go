package twigga

import "net/http"

type Client struct {
	BaseURL   string
	WSBaseURL string
	Token     string
	http      *http.Client
}

// imp
func NewTwiggaClient(token string) *Client {
	return &Client{
		BaseURL:   "https://twiga.bongocloud.co.tz",
		WSBaseURL: "wss://twiga.bongocloud.co.tz",
		Token:     token,
		http:      &http.Client{},
	}
}
