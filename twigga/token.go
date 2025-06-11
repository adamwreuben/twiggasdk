package twigga

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetTokenData(ctx context.Context, token string) (*TokenData, error) {
	url := fmt.Sprintf("%s/data/%s", c.AccountBaseURL, token)

	body, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp TokenData
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
