package twigga

import (
	"context"
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
