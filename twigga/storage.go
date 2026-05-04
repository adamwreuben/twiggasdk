package twigga

import (
	"bytes"
	"context"
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

type StorageService struct {
	client    *httpClient
	projectID string
}

// Bucket gives you a reference to a specific bucket for chaining methods
func (s *StorageService) Bucket(name string) *BucketReference {
	return &BucketReference{
		client: s.client,
		name:   name,
	}
}

// CreateBucket creates a new storage bucket (automatically passes projectId for Rule evaluation)
func (s *StorageService) CreateBucket(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/storage/buckets", s.client.baseURL)
	body := map[string]string{
		"name":      name,
		"projectId": s.projectID,
	}

	_, _, err := s.client.doRequest(ctx, http.MethodPost, url, body)
	return err
}

// ListBuckets fetches all buckets
func (s *StorageService) ListBuckets(ctx context.Context) ([]Bucket, error) {
	url := fmt.Sprintf("%s/storage/buckets", s.client.baseURL)
	body, _, err := s.client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Buckets []Bucket `json:"buckets"`
	}
	Unmarshal(body, &res)
	return res.Buckets, nil
}

// Usage retrieves the total storage metrics for the project
func (s *StorageService) Usage(ctx context.Context) (map[string]any, error) {
	url := fmt.Sprintf("%s/storage/projects/%s/usage", s.client.baseURL, s.projectID)
	body, _, err := s.client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]any
	Unmarshal(body, &res)
	return res, nil
}

type BucketReference struct {
	client *httpClient
	name   string
}

func (b *BucketReference) File(objectName string) *FileReference {
	return &FileReference{
		client: b.client,
		bucket: b.name,
		name:   objectName,
	}
}

// Get fetches bucket metadata
func (b *BucketReference) Get(ctx context.Context) (*Bucket, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s", b.client.baseURL, b.name)
	body, _, err := b.client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var bucket Bucket
	Unmarshal(body, &bucket)
	return &bucket, nil
}

// Delete removes the bucket (fails if not empty)
func (b *BucketReference) Delete(ctx context.Context) error {
	url := fmt.Sprintf("%s/storage/buckets/%s", b.client.baseURL, b.name)
	_, _, err := b.client.doRequest(ctx, http.MethodDelete, url, nil)
	return err
}

// SetPolicy updates the bucket access level ("public" or "private")
func (b *BucketReference) SetPolicy(ctx context.Context, policy string) error {
	url := fmt.Sprintf("%s/storage/buckets/%s/policy", b.client.baseURL, b.name)
	body := map[string]string{"policy": policy}
	_, _, err := b.client.doRequest(ctx, http.MethodPost, url, body)
	return err
}

// ListFiles fetches all file metadata inside the bucket
func (b *BucketReference) ListFiles(ctx context.Context) ([]FileObject, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/objects", b.client.baseURL, b.name)
	body, _, err := b.client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Files []FileObject `json:"files"`
	}
	Unmarshal(body, &res)
	return res.Files, nil
}

// Upload safely streams a local file to the bucket.
func (b *BucketReference) Upload(ctx context.Context, filePath string) ([]string, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/objects", b.client.baseURL, b.name)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Server expects the form key to be "files"
	part, err := writer.CreateFormFile("files", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, f); err != nil {
		return nil, err
	}
	writer.Close()

	bodyBytes, _, err := b.client.doMultipartRequest(ctx, http.MethodPost, url, &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var res struct {
		Files []string `json:"files"`
	}
	Unmarshal(bodyBytes, &res)
	return res.Files, nil
}

type FileReference struct {
	client *httpClient
	bucket string
	name   string
}

// Get retrieves the pre-signed URL and file metadata
func (f *FileReference) Get(ctx context.Context) (map[string]any, error) {
	url := fmt.Sprintf("%s/storage/buckets/%s/objects/%s", f.client.baseURL, f.bucket, f.name)
	body, _, err := f.client.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]any
	Unmarshal(body, &res)
	return res, nil
}

// Delete permanently removes the file
func (f *FileReference) Delete(ctx context.Context) error {
	url := fmt.Sprintf("%s/storage/buckets/%s/objects/%s", f.client.baseURL, f.bucket, f.name)
	_, _, err := f.client.doRequest(ctx, http.MethodDelete, url, nil)
	return err
}
