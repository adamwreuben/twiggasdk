package twigga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

type DatabaseService struct {
	client *httpClient
	dbId   string
}

type CollectionReference struct {
	client *httpClient
	dbId   string
	name   string
}

type getAllOptions struct {
	limit int
	sort  string
	field string
	after string
}

type DocumentReference struct {
	client *httpClient
	dbId   string
	coll   string
	id     string
}

type GetAllOption func(*getAllOptions)

type GetAllDocumentsResult struct {
	Documents  []map[string]any `json:"documents"`
	Total      int              `json:"total"`
	NextCursor any              `json:"nextCursor"`
}

func (d *DatabaseService) Collection(name string) *CollectionReference {
	return &CollectionReference{
		client: d.client,
		dbId:   d.dbId,
		name:   name,
	}
}

func (d *DatabaseService) CreateDatabase(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/database/%s", d.client.baseURL, d.dbId)
	res, _, err := d.client.doRequest(ctx, http.MethodPost, url, nil)
	return res, err
}

func (d *DatabaseService) DeleteDatabase(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/database/%s", d.client.baseURL, d.dbId)
	res, _, err := d.client.doRequest(ctx, http.MethodDelete, url, nil)
	return res, err
}

func (d *DatabaseService) ListCollections(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/database/%s", d.client.baseURL, d.dbId)
	res, _, err := d.client.doRequest(ctx, http.MethodGet, url, nil)
	return res, err
}

func (d *DocumentReference) Get(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", d.client.baseURL, d.dbId, d.coll, d.id)
	res, _, err := d.client.doRequest(ctx, http.MethodGet, url, nil)
	return res, err
}

func (d *DocumentReference) Exists(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/document/%s/%s/exists", d.client.baseURL, d.dbId, d.coll)

	filter := map[string]any{"id": d.id}

	body, status, err := d.client.doRequest(ctx, http.MethodPost, url, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check document existence: %w", err)
	}

	if status == http.StatusNotFound {
		return false, nil
	}

	if status != http.StatusOK {
		return false, fmt.Errorf("unexpected status %d: %s", status, string(body))
	}

	var result struct {
		Exists bool `json:"exists"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("failed to parse exists response: %v", err)
	}

	return result.Exists, nil
}

func (d *DocumentReference) Set(ctx context.Context, data any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", d.client.baseURL, d.dbId, d.coll, d.id)
	res, _, err := d.client.doRequest(ctx, http.MethodPost, url, data)
	return res, err
}

func (d *DocumentReference) Update(ctx context.Context, data any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", d.client.baseURL, d.dbId, d.coll, d.id)
	res, _, err := d.client.doRequest(ctx, http.MethodPut, url, data)
	return res, err
}

func (d *DocumentReference) Delete(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", d.client.baseURL, d.dbId, d.coll, d.id)
	res, _, err := d.client.doRequest(ctx, http.MethodDelete, url, nil)
	return res, err
}

func (d *DocumentReference) Listen() (*websocket.Conn, error) {
	endpoint := fmt.Sprintf("%s/document/%s/%s/%s/changes", d.client.wsBaseURL, d.dbId, d.coll, d.id)
	return d.client.openWS(endpoint)
}

func (c *CollectionReference) Doc(id string) *DocumentReference {
	return &DocumentReference{
		client: c.client,
		dbId:   c.dbId,
		coll:   c.name,
		id:     id,
	}
}

func (c *CollectionReference) Add(ctx context.Context, data any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s", c.client.baseURL, c.dbId, c.name)
	res, _, err := c.client.doRequest(ctx, http.MethodPost, url, data)
	return res, err
}

func (c *CollectionReference) BulkAdd(ctx context.Context, docs []any, groupFields []string) (string, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/document/%s/%s/bulk", c.client.baseURL, c.dbId, c.name))

	q := u.Query()
	for _, f := range groupFields {
		q.Add("group", f)
	}
	u.RawQuery = q.Encode()

	body, statusCode, err := c.client.doRequest(ctx, http.MethodPost, u.String(), docs)
	if err != nil {
		return "", err
	}
	if statusCode != http.StatusOK {
		return "", fmt.Errorf("bulk insert failed: %d, %s", statusCode, string(body))
	}
	return string(body), nil
}

func (c *CollectionReference) CollectionExists(ctx context.Context) (bool, error) {
	url := fmt.Sprintf("%s/collection/%s/%s/exists", c.client.baseURL, c.dbId, c.name)

	body, status, err := c.client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to check collection existence: %w", err)
	}

	if status != http.StatusOK {
		return false, fmt.Errorf("unexpected status %d: %s", status, string(body))
	}

	var result struct {
		Exists bool `json:"exists"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("failed to parse exists response: %v", err)
	}

	return result.Exists, nil
}

func (c *CollectionReference) Filter(ctx context.Context, filter map[string]any, options ...map[string]string) (*ReadAllDocumentsResult, error) {
	basePath := fmt.Sprintf("%s/document/%s/%s/filter", c.client.baseURL, c.dbId, c.name)

	u, _ := url.Parse(basePath)
	if len(options) > 0 {
		q := u.Query()
		for k, v := range options[0] {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	body, status, err := c.client.doRequest(ctx, http.MethodPost, u.String(), filter)
	if err != nil {
		return nil, err
	}
	if status == http.StatusTooManyRequests {
		return nil, errors.New("too many requests per IP")
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", status, string(body))
	}

	var res ReadAllDocumentsResult
	json.Unmarshal(body, &res)
	return &res, nil
}

func (c *CollectionReference) GetAll(ctx context.Context, opts ...GetAllOption) (*ReadAllDocumentsResult, error) {
	basePath := fmt.Sprintf("%s/document/%s/%s", c.client.baseURL, c.dbId, c.name)

	u, err := url.Parse(basePath)
	if err != nil {
		return nil, err
	}

	// Default options
	options := &getAllOptions{
		limit: 300,
		sort:  "asc",
		field: "id",
	}

	// Apply user options
	for _, opt := range opts {
		opt(options)
	}

	// Build query parameters
	q := u.Query()
	q.Set("limit", strconv.Itoa(options.limit))
	q.Set("sort", options.sort)
	q.Set("field", options.field)

	if options.after != "" {
		q.Set("after", options.after)
	}

	u.RawQuery = q.Encode()

	body, status, err := c.client.doRequest(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if status == http.StatusTooManyRequests {
		return nil, errors.New("too many requests per IP")
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", status, string(body))
	}

	var res ReadAllDocumentsResult
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &res, nil
}

func (c *CollectionReference) Listen() (*websocket.Conn, error) {
	endpoint := fmt.Sprintf("%s/document/%s/%s/changes", c.client.wsBaseURL, c.dbId, c.name)
	return c.client.openWS(endpoint)
}

func WithLimit(limit int) GetAllOption {
	return func(o *getAllOptions) {
		if limit > 0 && limit <= 1000 {
			o.limit = limit
		}
	}
}

func WithSort(sort string) GetAllOption {
	return func(o *getAllOptions) {
		if sort == "asc" || sort == "desc" {
			o.sort = sort
		}
	}
}

func WithField(field string) GetAllOption {
	return func(o *getAllOptions) {
		if field != "" {
			o.field = field
		}
	}
}

func WithAfter(cursor string) GetAllOption {
	return func(o *getAllOptions) {
		o.after = cursor
	}
}
