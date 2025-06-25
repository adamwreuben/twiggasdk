package twigga

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetTokenData(ctx context.Context, token string) (*TokenData, error) {
	url := fmt.Sprintf("%s/user/token/%s", c.accountBaseURL, token)

	fmt.Println("url: ", url)

	body, statusCode, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp TokenData
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	fmt.Println("statusCode: ", statusCode)

	return &resp, nil
}
