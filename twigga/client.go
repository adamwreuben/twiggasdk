package twigga

import "net/http"

type BuildType string

const (
	TEST BuildType = "TEST"
	PROD BuildType = "PROD"
)

func NewTwiggaClient(confPath string, buildType BuildType) (*Client, error) {

	bongoClient, err := LoadConfig(confPath)
	if err != nil {
		return nil, err
	}

	baseUrl := PROD_BaseURL
	wsBaseUrl := PROD_BaseURL
	accountBaseUrl := PROD_AccountBaseURL

	if buildType == TEST {
		baseUrl = TEST_BaseURL
		wsBaseUrl = TEST_BaseURL
		accountBaseUrl = TEST_AccountBaseURL
	}

	return &Client{
		baseURL:        baseUrl,
		wSBaseURL:      wsBaseUrl,
		accountBaseURL: accountBaseUrl,
		client:         *bongoClient,
		http:           &http.Client{},
	}, nil
}
