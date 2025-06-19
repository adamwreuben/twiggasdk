package twigga

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// imp
func (c *Client) doRequest(ctx context.Context, method, url string, body any) ([]byte, int, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.client.Token != "" {
		req.Header.Set("BONGO-TOKEN", c.client.Token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println("doReq: ", resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return bodyBytes, resp.StatusCode, nil
}

// CreateDocumentAuto creates a new document with auto-generated ID
func (c *Client) CreateDocumentAuto(ctx context.Context, collection string, doc any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s", c.baseURL, c.client.Twigga.DefaultDatabase, collection)
	res, _, err := c.doRequest(ctx, http.MethodPost, url, doc)
	return res, err
}

// CreateDocumentWithID creates a document with a specified ID
func (c *Client) CreateDocumentWithID(ctx context.Context, collection, id string, doc any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", c.baseURL, c.client.Twigga.DefaultDatabase, collection, id)
	res, _, err := c.doRequest(ctx, http.MethodPost, url, doc)
	return res, err
}

// GetDocument fetches a document by ID
func (c *Client) GetDocument(ctx context.Context, collection, id string) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", c.baseURL, c.client.Twigga.DefaultDatabase, collection, id)

	res, _, err := c.doRequest(ctx, http.MethodGet, url, nil)
	return res, err
}

// return list of filetered documents
func (c *Client) QueryDocuments(ctx context.Context, collection string, filter map[string]any) (map[string]interface{}, error) {
	fmt.Println("WEWE herrer*")
	url := fmt.Sprintf("%s/document/%s/%s/filter", c.baseURL, c.client.Twigga.DefaultDatabase, collection)

	body, statusCode, err := c.doRequest(ctx, http.MethodPost, url, filter)
	fmt.Println("statusCode for querying: ", statusCode)

	if err != nil {
		fmt.Println("do req error: ", err.Error())
		return nil, err
	}

	if statusCode == http.StatusOK {
		var doc map[string]interface{}
		if err := json.Unmarshal(body, &doc); err != nil {
			return nil, err
		}
	}

	if statusCode == 429 {
		return nil, errors.New("too many request per IP, please try again later")
	}

	return nil, errors.New("Unknown error!")
}

func (c *Client) CollectionExists(ctx context.Context, collection string) (bool, error) {
	url := fmt.Sprintf("%s/collection/%s/%s/exists", c.baseURL, c.client.Twigga.DefaultDatabase, collection)

	body, _, err := c.doRequest(ctx, http.MethodGet, url, nil)

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
	url := fmt.Sprintf("%s/document/%s/%s/exists", c.baseURL, c.client.Twigga.DefaultDatabase, collection)

	body, statusCode, err := c.doRequest(ctx, http.MethodPost, url, filter)

	if err != nil {
		return false, err
	}

	var result struct {
		Exists bool `json:"exists"`
	}

	if statusCode == http.StatusOK {

		if err := json.Unmarshal(body, &result); err != nil {
			return false, err
		}

		return result.Exists, nil
	}

	if statusCode == 429 { // DNS too many request
		return false, errors.New("too many request per IP, please try again later")
	}

	return false, nil

}

// GetCollection fetches all documents from a table
func (c *Client) GetCollection(ctx context.Context, collection string) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s", c.baseURL, c.client.Twigga.DefaultDatabase, collection)

	res, _, err := c.doRequest(ctx, http.MethodGet, url, nil)
	return res, err
}

// UpdateDocument updates a document by ID
func (c *Client) UpdateDocument(ctx context.Context, collection, id string, doc any) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", c.baseURL, c.client.Twigga.DefaultDatabase, collection, id)

	res, _, err := c.doRequest(ctx, http.MethodPut, url, doc)
	return res, err
}

// DeleteDocument deletes a document by ID
func (c *Client) DeleteDocument(ctx context.Context, collection, id string) ([]byte, error) {
	url := fmt.Sprintf("%s/document/%s/%s/%s", c.baseURL, c.client.Twigga.DefaultDatabase, collection, id)

	res, _, err := c.doRequest(ctx, http.MethodDelete, url, nil)
	return res, err
}

// CreateDatabase creates a new database
func (c *Client) CreateDatabase(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/database/%s", c.baseURL, c.client.Twigga.DefaultDatabase)

	res, _, err := c.doRequest(ctx, http.MethodPost, url, nil)
	return res, err
}

// DeleteDatabase deletes a database
func (c *Client) DeleteDatabase(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/database/%s", c.baseURL, c.client.Twigga.DefaultDatabase)

	res, _, err := c.doRequest(ctx, http.MethodDelete, url, nil)
	return res, err
}

// ListAllCollections lists collections in a database
func (c *Client) ListAllCollections(ctx context.Context) ([]byte, error) {
	url := fmt.Sprintf("%s/database/%s", c.baseURL, c.client.Twigga.DefaultDatabase)

	res, _, err := c.doRequest(ctx, http.MethodGet, url, nil)
	return res, err
}

// DeleteCollection deletes a collection in a database
func (c *Client) DeleteCollection(ctx context.Context, collection string) ([]byte, error) {
	url := fmt.Sprintf("%s/collection/%s/%s", c.baseURL, c.client.Twigga.DefaultDatabase, collection)

	res, _, err := c.doRequest(ctx, http.MethodDelete, url, nil)
	return res, err
}
