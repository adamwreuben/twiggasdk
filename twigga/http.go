package twigga

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// imp
func (c *Client) doRequest(ctx context.Context, method, url string, body any) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("BONGO-TOKEN", c.Token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// CreateDocumentAuto creates a new document with auto-generated ID
func (c *Client) CreateDocumentAuto(ctx context.Context, collection string, doc any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s", c.BaseURL, c.DefaultUserDatabase, collection)
	return c.doRequest(ctx, http.MethodPost, url, doc)
}

// CreateDocumentWithID creates a document with a specified ID
func (c *Client) CreateDocumentWithID(ctx context.Context, collection, id string, doc any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", c.BaseURL, c.DefaultUserDatabase, collection, id)
	return c.doRequest(ctx, http.MethodPost, url, doc)
}

// GetDocument fetches a document by ID
func (c *Client) GetDocument(ctx context.Context, collection, id string) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", c.BaseURL, c.DefaultUserDatabase, collection, id)
	return c.doRequest(ctx, http.MethodGet, url, nil)
}

// return list of filetered documents
func (c *Client) QueryDocuments(ctx context.Context, collection string, filter map[string]any) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/document/%s/%s/filter", c.BaseURL, c.DefaultUserDatabase, collection)

	body, err := c.doRequest(ctx, http.MethodPost, url, filter)
	if err != nil {
		return nil, err
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (c *Client) CollectionExists(ctx context.Context, collection string) (bool, error) {
	url := fmt.Sprintf("%s/collection/%s/%s/exists", c.BaseURL, c.DefaultUserDatabase, collection)

	body, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	var result struct {
		Exists bool `json:"exists"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	return result.Exists, nil
}

func (c *Client) DocumentExists(ctx context.Context, collection string, filter map[string]any) (bool, error) {
	url := fmt.Sprintf("%s/document/%s/%s/exists", c.BaseURL, c.DefaultUserDatabase, collection)

	body, err := c.doRequest(ctx, http.MethodPost, url, filter)
	if err != nil {
		return false, err
	}

	var result struct {
		Exists bool `json:"exists"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	return result.Exists, nil
}

// GetCollection fetches all documents from a table
func (c *Client) GetCollection(ctx context.Context, collection string) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s", c.BaseURL, c.DefaultUserDatabase, collection)
	return c.doRequest(ctx, http.MethodGet, url, nil)
}

// UpdateDocument updates a document by ID
func (c *Client) UpdateDocument(ctx context.Context, collection, id string, doc any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", c.BaseURL, c.DefaultUserDatabase, collection, id)
	return c.doRequest(ctx, http.MethodPut, url, doc)
}

// DeleteDocument deletes a document by ID
func (c *Client) DeleteDocument(ctx context.Context, collection, id string) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", c.BaseURL, c.DefaultUserDatabase, collection, id)
	return c.doRequest(ctx, http.MethodDelete, url, nil)
}

// CreateDatabase creates a new database
func (c *Client) CreateDatabase(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/database/%s", c.BaseURL, c.DefaultUserDatabase)
	return c.doRequest(ctx, http.MethodPost, url, nil)
}

// DeleteDatabase deletes a database
func (c *Client) DeleteDatabase(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/database/%s", c.BaseURL, c.DefaultUserDatabase)
	return c.doRequest(ctx, http.MethodDelete, url, nil)
}

// ListAllCollections lists collections in a database
func (c *Client) ListAllCollections(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/database/%s", c.BaseURL, c.DefaultUserDatabase)
	return c.doRequest(ctx, http.MethodGet, url, nil)
}

// DeleteCollection deletes a collection in a database
func (c *Client) DeleteCollection(ctx context.Context, collection string) ([]byte, error) {
	url := fmt.Sprintf("%s/collection/%s/%s", c.BaseURL, c.DefaultUserDatabase, collection)
	return c.doRequest(ctx, http.MethodDelete, url, nil)
}
