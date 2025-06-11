package twigga

import "net/http"

type Client struct {
	BaseURL        string // Document API
	WSBaseURL      string
	AccountBaseURL string // Account API
	Token          string
	http           *http.Client
}

type AppTokenRequest struct {
	AppID     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

type AppTokenResponse struct {
	AccessToken string `json:"accessToken"`
	Exp         int64  `json:"exp"`
}

type AuthenticateRequest struct {
	RedirectTo string `json:"redirectTo"`
}

type AuthenticateResponse struct {
	AuthURL string `json:"authUrl"`
}

type TokenData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}
