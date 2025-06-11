package twigga

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GenerateAppToken(ctx context.Context, appID, appSecret string) (*AppTokenResponse, error) {
	url := fmt.Sprintf("%s/application/token", c.AccountBaseURL)

	reqBody := AppTokenRequest{
		AppID:     appID,
		AppSecret: appSecret,
	}

	body, err := c.doRequest(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	var resp AppTokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Authenticate(ctx context.Context, redirectTo string) (*AuthenticateResponse, error) {
	url := fmt.Sprintf("%s/application/authenticate", c.AccountBaseURL)

	reqBody := AuthenticateRequest{RedirectTo: redirectTo}

	body, err := c.doRequest(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	var resp AuthenticateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateAccount(ctx context.Context, req CreateAccountRequest) (*MessageResponse, error) {
	url := fmt.Sprintf("%s/user/create", c.AccountBaseURL)

	body, err := c.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}

	var resp MessageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	url := fmt.Sprintf("%s/user/login", c.AccountBaseURL)

	reqBody := LoginRequest{Email: email, Password: password}

	body, err := c.doRequest(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}

	var resp LoginResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Logout(ctx context.Context, userID string) (*MessageResponse, error) {
	url := fmt.Sprintf("%s/user/logout/%s", c.AccountBaseURL, userID)

	body, err := c.doRequest(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	var resp MessageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
