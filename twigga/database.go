package twigga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type DatabaseService struct {
	client *httpClient
	dbId   string
}

// Collection gives you access to a specific table
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

// Collection reference
type CollectionReference struct {
	client *httpClient
	dbId   string
	name   string
}

// Doc drills down into a specific document by ID
func (c *CollectionReference) Doc(id string) *DocumentReference {
	return &DocumentReference{
		client: c.client,
		dbId:   c.dbId,
		coll:   c.name,
		id:     id,
	}
}

// Add creates a document with an Auto-Generated ID
func (c *CollectionReference) Add(ctx context.Context, data any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s", c.client.baseURL, c.dbId, c.name)
	res, _, err := c.client.doRequest(ctx, http.MethodPost, url, data)
	return res, err
}

// BulkAdd inserts multiple documents at once
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

func (c *CollectionReference) Exists(ctx context.Context, filter map[string]any) (bool, error) {

	if filter == nil || len(filter) == 0 {
		url := fmt.Sprintf("%s/collection/%s/%s/exists", c.client.baseURL, c.dbId, c.name)
		body, _, err := c.client.doRequest(ctx, http.MethodGet, url, nil)
		if err != nil {
			return false, err
		}

		var result struct {
			Exists bool `json:"exists"`
		}
		json.Unmarshal(body, &result)
		return result.Exists, nil
	}

	options := map[string]string{
		"limit": "1",
	}

	res, err := c.Filter(ctx, filter, options)
	if err != nil {
		return false, err
	}

	return len(res.Documents) > 0, nil
}

// Filter allows querying the collection
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

// DOCUMENT REFERENCE
type DocumentReference struct {
	client *httpClient
	dbId   string
	coll   string
	id     string
}

// Get fetches the document
func (d *DocumentReference) Get(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", d.client.baseURL, d.dbId, d.coll, d.id)
	res, _, err := d.client.doRequest(ctx, http.MethodGet, url, nil)
	return res, err
}

// Set creates or overwrites a document with a specific ID
func (d *DocumentReference) Set(ctx context.Context, data any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", d.client.baseURL, d.dbId, d.coll, d.id)
	res, _, err := d.client.doRequest(ctx, http.MethodPost, url, data)
	return res, err
}

// Update modifies an existing document
func (d *DocumentReference) Update(ctx context.Context, data any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", d.client.baseURL, d.dbId, d.coll, d.id)
	res, _, err := d.client.doRequest(ctx, http.MethodPut, url, data)
	return res, err
}

// Delete removes the document
func (d *DocumentReference) Delete(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", d.client.baseURL, d.dbId, d.coll, d.id)
	res, _, err := d.client.doRequest(ctx, http.MethodDelete, url, nil)
	return res, err
}

func (c *CollectionReference) Listen() (*websocket.Conn, error) {
	endpoint := fmt.Sprintf("%s/document/%s/%s/changes", c.client.wsBaseURL, c.dbId, c.name)
	return c.client.openWS(endpoint)
}

func (d *DocumentReference) Listen() (*websocket.Conn, error) {
	endpoint := fmt.Sprintf("%s/document/%s/%s/%s/changes", d.client.wsBaseURL, d.dbId, d.coll, d.id)
	return d.client.openWS(endpoint)
}
