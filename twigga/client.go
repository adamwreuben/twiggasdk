package twigga

import "net/http"

// imp
func NewTwiggaClient(token string) *Client {
	return &Client{
		BaseURL:        "https://twiga.bongocloud.co.tz",
		WSBaseURL:      "wss://twiga.bongocloud.co.tz",
		AccountBaseURL: "https://account.bongocloud.co.tz",
		Token:          token,
		http:           &http.Client{},
	}
}
