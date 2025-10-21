package twigga

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetTokenData(ctx context.Context, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/user/token/%s", c.accountBaseURL, token)

	body, statusCode, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	fmt.Println("statusCode: ", statusCode)

	return resp, nil
}
