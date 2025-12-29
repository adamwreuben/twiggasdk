package twigga

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GenerateAppToken(ctx context.Context, appID, appSecret string) (*AppTokenResponse, error) {
	url := fmt.Sprintf("%s/application/token", c.accountBaseURL)

	reqBody := AppTokenRequest{
		AppID:     appID,
		AppSecret: appSecret,
	}

	body, statusCode, err := c.doRequest(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	var resp AppTokenResponse
	if err := Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	fmt.Println("statusCode: ", statusCode)

	return &resp, nil
}

func (c *Client) Authenticate(ctx context.Context, redirectTo string) (*AuthenticateResponse, error) {
	url := fmt.Sprintf("%s/application/authenticate", c.accountBaseURL)

	reqBody := AuthenticateRequest{RedirectTo: redirectTo}

	body, statusCode, err := c.doRequest(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	var resp AuthenticateResponse
	if err := Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	fmt.Println("statusCode: ", statusCode)

	return &resp, nil
}

func (c *Client) CreateAccount(ctx context.Context, req CreateAccountRequest) (*MessageResponse, error) {
	url := fmt.Sprintf("%s/user/create", c.accountBaseURL)

	body, statusCode, err := c.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}

	var resp MessageResponse
	if err := Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	fmt.Println("statusCode: ", statusCode)

	return &resp, nil
}

func (c *Client) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	url := fmt.Sprintf("%s/user/login", c.accountBaseURL)

	reqBody := LoginRequest{Email: email, Password: password}

	body, statusCode, err := c.doRequest(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	fmt.Println("LoginRequest body: ", string(body))

	var resp LoginResponse
	if err := Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	fmt.Println("statusCode: ", statusCode)

	return &resp, nil
}

func (c *Client) Logout(ctx context.Context, userID string) (*MessageResponse, error) {
	url := fmt.Sprintf("%s/user/logout/%s", c.accountBaseURL, userID)

	body, statusCode, err := c.doRequest(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	var resp MessageResponse
	if err := Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	fmt.Println("statusCode: ", statusCode)

	return &resp, nil
}

func (c *Client) CheckAuthorization(ctx context.Context, subID, subType, relation, objType, objID string) (bool, error) {
	url := fmt.Sprintf("%s/authorize/check?subjectId=%s&subjectType=%s&relation=%s&objectType=%s&objectId=%s",
		c.baseURL, subID, subType, relation, objType, objID)

	body, status, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	if status != http.StatusOK {
		return false, fmt.Errorf("auth check failed with status %d", status)
	}

	var result struct {
		Allowed bool `json:"allowed"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	return result.Allowed, nil
}

func (c *Client) AssignAuthorization(ctx context.Context, authTuple AuthorizationTuple) error {
	url := fmt.Sprintf("%s/authorize/write", c.baseURL)
	_, status, err := c.doRequest(ctx, http.MethodPost, url, authTuple)

	if status != http.StatusOK && status != http.StatusCreated {
		return fmt.Errorf("failed to assign authorization, status: %d", status)
	}
	return err
}
