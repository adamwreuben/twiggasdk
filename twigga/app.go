package twigga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type App struct {
	baseClient *httpClient
	projectID  string
}

type ServiceAccount struct {
	ProjectID string `json:"projectId"`
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

type appConfig struct {
	credsFile string
	buildType BuildType
}

type httpClient struct {
	baseURL    string
	wsBaseURL  string
	accountURL string
	http       *http.Client

	appId      string
	appSecret  string
	token      string
	tokenMutex sync.RWMutex
	expiresAt  time.Time
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
		json.Unmarshal(fileBytes, &sa)
	} else if os.Getenv("TWIGGA_APPLICATION_CREDENTIALS") != "" {
		fileBytes, _ := os.ReadFile(os.Getenv("TWIGGA_APPLICATION_CREDENTIALS"))
		json.Unmarshal(fileBytes, &sa)
	} else {
		return nil, errors.New("no credentials provided")
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
		appId:      sa.AppId,
		appSecret:  sa.AppSecret,
	}

	return &App{
		baseClient: client,
		projectID:  sa.ProjectID,
	}, nil
}

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
	return &DatabaseService{client: a.baseClient, dbId: a.projectID}
}
