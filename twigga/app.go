package twigga

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type App struct {
	baseClient *httpClient
	projectID  string
	databaseID string
}

type ServiceAccount struct {
	Type         string `json:"type"`
	ProjectID    string `json:"project_id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	DatabaseID   string `json:"database_id"`
	AuthURI      string `json:"auth_uri"`
}

type appConfig struct {
	credsFile string
	buildType BuildType
}

type httpClient struct {
	baseURL    string
	wsBaseURL  string
	accountURL string
	token      string
	http       *http.Client
}

type ClientOption func(*appConfig)

func WithCredentialsFile(filename string) ClientOption {
	return func(c *appConfig) {
		c.credsFile = filename
	}
}

func WithEnvironment(env BuildType) ClientOption {
	return func(c *appConfig) {
		c.buildType = env
	}
}

func InitializeApp(ctx context.Context, opts ...ClientOption) (*App, error) {
	config := &appConfig{buildType: PROD}
	for _, opt := range opts {
		opt(config)
	}

	var sa ServiceAccount

	if config.credsFile != "" {
		fileBytes, err := os.ReadFile(config.credsFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials file: %v", err)
		}
		if err := json.Unmarshal(fileBytes, &sa); err != nil {
			return nil, fmt.Errorf("invalid credentials JSON format: %v", err)
		}
	} else if os.Getenv("TWIGGA_APPLICATION_CREDENTIALS") != "" {
		fileBytes, err := os.ReadFile(os.Getenv("TWIGGA_APPLICATION_CREDENTIALS"))
		if err == nil {
			json.Unmarshal(fileBytes, &sa)
		}
	} else {
		return nil, errors.New("no credentials provided. Use WithCredentialsFile or set TWIGGA_APPLICATION_CREDENTIALS")
	}

	baseURL := PROD_BaseURL
	wsBaseURL := PROD_WSBaseURL
	accountURL := PROD_AccountBaseURL

	if config.buildType == TEST {
		baseURL = TEST_BaseURL
		wsBaseURL = TEST_WSBaseURL
		accountURL = TEST_AccountBaseURL
	}

	client := &httpClient{
		baseURL:    baseURL,
		accountURL: accountURL,
		wsBaseURL:  wsBaseURL,
		http:       &http.Client{},
	}

	if sa.AuthURI != "" && sa.ClientID != "" {
		token, err := fetchAccessToken(ctx, sa.AuthURI, sa.ClientID, sa.ClientSecret)
		if err != nil {
			return nil, fmt.Errorf("SDK authentication failed: %v", err)
		}
		client.token = token
	}

	return &App{
		baseClient: client,
		projectID:  sa.ProjectID,
		databaseID: sa.DatabaseID,
	}, nil
}

// Internal helper to fetch the OAuth token using the new Auth endpoint
func fetchAccessToken(ctx context.Context, authUri, clientId, clientSecret string) (string, error) {
	reqBody := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     clientId,
		"client_secret": clientSecret,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", authUri, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth server returned status: %d", resp.StatusCode)
	}

	var res struct {
		AccessToken string `json:"access_token"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &res)

	return res.AccessToken, nil
}

// NAMESPACE ACCESSORS
func (a *App) Auth() *AuthService {
	return &AuthService{client: a.baseClient, projectID: a.projectID}
}
func (a *App) Functions() *FunctionsService {
	return &FunctionsService{client: a.baseClient, projectID: a.projectID}
}
func (a *App) Rules() *RulesService {
	return &RulesService{client: a.baseClient, projectID: a.projectID}
}
func (a *App) Backups() *BackupService {
	return &BackupService{client: a.baseClient, projectID: a.projectID}
}
func (a *App) Storage() *StorageService {
	return &StorageService{client: a.baseClient, projectID: a.projectID}
}
func (a *App) Database() *DatabaseService {
	return &DatabaseService{client: a.baseClient, dbId: a.databaseID}
}
