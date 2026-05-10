package twigga

import (
	"context"
	"fmt"
	"net/http"
)

type CreateUserReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ClientId  string `json:"client_id"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// AUTH SERVICE
type AuthService struct {
	client    *httpClient
	projectID string
}

// Signup creates a new user and returns their initial OAuth tokens.
func (a *AuthService) Signup(ctx context.Context, req CreateUserReq) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/auth/signup", a.client.accountURL)

	body, status, err := a.client.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusCreated && status != http.StatusOK {
		return nil, fmt.Errorf("signup failed with status %d: %s", status, string(body))
	}

	var res TokenResponse
	Unmarshal(body, &res)
	return &res, nil
}

// Login authenticates an existing user via Email & Password.
func (a *AuthService) Login(ctx context.Context, email, password, clientId string) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/auth/login", a.client.accountURL)
	req := map[string]string{
		"email":     email,
		"password":  password,
		"client_id": clientId,
	}

	body, status, err := a.client.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("login failed with status %d: %s", status, string(body))
	}

	var res TokenResponse
	Unmarshal(body, &res)
	return &res, nil
}

// LoginWithProvider handles Identity (e.g., passing a Google Token).
func (a *AuthService) LoginWithProvider(ctx context.Context, providerName, credential, clientId string) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/auth/provider/%s", a.client.accountURL, providerName)
	req := map[string]string{
		"credential": credential,
		"client_id":  clientId,
	}

	body, status, err := a.client.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("provider login failed with status %d: %s", status, string(body))
	}

	var res TokenResponse
	Unmarshal(body, &res)
	return &res, nil
}

// RefreshToken trades a 30-day Refresh Token for a new 15-minute Access Token.
func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/oauth/token", a.client.accountURL)
	req := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	body, status, err := a.client.doRequest(ctx, http.MethodPost, url, req)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed with status %d", status)
	}

	var res TokenResponse
	Unmarshal(body, &res)
	return &res, nil
}

// Logout kills a specific session device by revoking its Refresh Token.
func (a *AuthService) Logout(ctx context.Context, refreshToken string) error {
	url := fmt.Sprintf("%s/oauth/revoke", a.client.accountURL)
	req := map[string]string{"token": refreshToken}
	_, _, err := a.client.doRequest(ctx, http.MethodPost, url, req)
	return err
}

// ForgotPassword triggers an email to the user with a reset link.
func (a *AuthService) ForgotPassword(ctx context.Context, email string) error {
	url := fmt.Sprintf("%s/auth/forgot-password", a.client.accountURL)
	req := map[string]string{"email": email}
	_, _, err := a.client.doRequest(ctx, http.MethodPost, url, req)
	return err
}

// ResetPassword takes the token from the email and applies the new password.
func (a *AuthService) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
	url := fmt.Sprintf("%s/auth/reset-password", a.client.accountURL)
	req := map[string]string{
		"token":        resetToken,
		"new_password": newPassword,
	}
	_, _, err := a.client.doRequest(ctx, http.MethodPost, url, req)
	return err
}

// ZANZIBAR NAMESPACE
type AuthorizationTuple struct {
	ObjectType  string `json:"objectType"`
	ObjectID    string `json:"objectId"`
	Relation    string `json:"relation"`
	SubjectType string `json:"subjectType"`
	SubjectID   string `json:"subjectId"`
}

func (a *AuthService) Zanzibar() *ZanzibarService {
	return &ZanzibarService{client: a.client}
}

type ZanzibarService struct {
	client *httpClient
}

func (z *ZanzibarService) Write(ctx context.Context, tuple AuthorizationTuple) error {
	url := fmt.Sprintf("%s/authorize/write", z.client.baseURL)

	_, status, err := z.client.doRequest(ctx, http.MethodPost, url, tuple)
	if err != nil {
		return err
	}

	if status != http.StatusOK && status != http.StatusCreated {
		return fmt.Errorf("failed to assign authorization, status: %d", status)
	}
	return nil
}

func (z *ZanzibarService) Check(ctx context.Context, subID, subType, relation, objType, objID string) (bool, error) {
	url := fmt.Sprintf("%s/authorize/check?subjectId=%s&subjectType=%s&relation=%s&objectType=%s&objectId=%s",
		z.client.baseURL, subID, subType, relation, objType, objID)

	body, _, err := z.client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	var res struct {
		Allowed bool `json:"allowed"`
	}
	Unmarshal(body, &res)
	return res.Allowed, nil
}
