package twigga

import "net/http"

// imp
func NewTwiggaClient(confPath string) (*Client, error) {

	bongoClient, err := LoadConfig(confPath)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL:        BaseURL,
		wSBaseURL:      WSBaseURL,
		accountBaseURL: AccountBaseURL,
		client:         *bongoClient,
		http:           &http.Client{},
	}, nil
}
