package twigga

import (
	"net/http"
)

type BongoCloudClient struct {
	Token string `json:"token"`
	Auth  struct {
		AppId     string `json:"appId"`
		AppSecret string `json:"appSecret"`
	} `json:"auth"`
	Twigga struct {
		DefaultDatabase string `json:"databaseId"`
	} `json:"twigga"`
}

type Client struct {
	baseURL        string // Document API
	wSBaseURL      string
	accountBaseURL string // Account API
	client         BongoCloudClient
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
	ID     string        `json:"id"`
	Email  string        `json:"email"`
	Events []interface{} `json:"events"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ReadAllDocumentsResult struct {
	Documents  []map[string]any `json:"documents"`
	Total      int              `json:"total"`
	NextCursor any              `json:"nextCursor"`
}

type AuthorizationTuple struct {
	ID          string `json:"id,omitempty"`
	ObjectType  string `json:"objectType"`
	ObjectID    string `json:"objectId"`
	Relation    string `json:"relation"`
	SubjectType string `json:"subjectType"` // "user" or "group"
	SubjectID   string `json:"subjectId"`
}

type Query struct {
	Where   []Condition `json:"where"`
	OrderBy string      `json:"orderBy"` // Field name
	Sort    string      `json:"sort"`    // "asc" or "desc"
	Limit   int         `json:"limit"`
	After   any         `json:"after"` // The Cursor
}

type Condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"` // "==", ">", "<", "array-contains", "in"
	Value    any    `json:"value"`
}
