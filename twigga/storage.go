package twigga

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Bucket struct {
	Name string `json:"name"`
}

type FileObject struct {
	Name         string    `json:"name"`
	ContentType  string    `json:"contentType"`
	Size         string    `json:"size"`
	LastModified time.Time `json:"lastModified"`
	ETag         string    `json:"etag"`
	StorageClass string    `json:"storageClass"`
}

type FilesResponse struct {
	Files []FileObject `json:"files"`
}

type FileResponse struct {
	URL    string `json:"url"`
	Bucket string `json:"bucket"`
	File   string `json:"file"`
}

// -------------------- Buckets --------------------

// AddBucket creates a new bucket
func (c *Client) AddBucket(ctx context.Context, bucket string) error {
	url := fmt.Sprintf("%s/storage/buckets", c.baseURL)
	body := map[string]string{"name": bucket}

	_, status, err := c.doRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}
	if status != http.StatusOK && status != http.StatusCreated {
		return fmt.Errorf("failed to create bucket, status %d", status)
	}
	return nil
}

// GetBuckets lists all buckets
func (c *Client) GetBuckets(ctx context.Context) ([]Bucket, error) {
	url := fmt.Sprintf("%s/storage/buckets", c.baseURL)
	body, _, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var buckets []Bucket
	if err := json.Unmarshal(body, &buckets); err != nil {
		return nil, err
	}
	return buckets, nil
}

// GetBucket fetches details of a single bucket
func (c *Client) GetBucket(ctx context.Context, bucket string) (*Bucket, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s", c.baseURL, bucket)
	body, status, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("bucket %s not found", bucket)
	}
	var b Bucket
	if err := json.Unmarshal(body, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// DeleteBucket deletes a bucket
func (c *Client) DeleteBucket(ctx context.Context, bucket string) error {
	url := fmt.Sprintf("%s/storage/buckets/%s", c.baseURL, bucket)
	_, status, err := c.doRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("failed to delete bucket %s, status %d", bucket, status)
	}
	return nil
}

// -------------------- Objects --------------------

// UploadFile uploads a file to a bucket
func (c *Client) UploadFile(ctx context.Context, bucket, filePath string) (*FileObject, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/objects", c.baseURL, bucket)

	// open file
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// prepare multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.client.Token != "" {
		req.Header.Set("BONGO-TOKEN", c.client.Token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("upload failed: %s", string(bodyBytes))
	}

	var fo FileObject
	if err := json.Unmarshal(bodyBytes, &fo); err != nil {
		return nil, err
	}

	return &fo, nil
}

// GetFiles lists all files in a bucket
func (c *Client) GetFiles(ctx context.Context, bucket string) ([]FileObject, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/objects", c.baseURL, bucket)
	body, _, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var files []FileObject
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, err
	}
	return files, nil
}

// GetFile downloads a file from bucket
func (c *Client) GetFile(ctx context.Context, bucket, object string) ([]byte, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/objects/%s", c.baseURL, bucket, object)
	body, status, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("file %s not found", object)
	}
	return body, nil
}

// DeleteFile deletes a file from bucket
func (c *Client) DeleteFile(ctx context.Context, bucket, object string) error {
	url := fmt.Sprintf("%s/storage/buckets/%s/objects/%s", c.baseURL, bucket, object)
	_, status, err := c.doRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("failed to delete file %s, status %d", object, status)
	}
	return nil
}
