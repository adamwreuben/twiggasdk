package twigga

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
)

type FunctionsService struct {
	client    *httpClient
	projectID string
}

func (f *FunctionsService) Deploy(ctx context.Context, functionName string, zipReader io.Reader, runtime string) (*map[string]any, error) {
	url := fmt.Sprintf("%s/functions/%s/%s/deploy", f.client.baseURL, f.projectID, functionName)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add runtime field
	writer.WriteField("runtime", runtime)

	// Add file
	part, _ := writer.CreateFormFile("code", "function.zip")
	io.Copy(part, zipReader)
	writer.Close()

	body, _, err := f.client.doMultipartRequest(ctx, "POST", url, &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var res map[string]any
	Unmarshal(body, &res)
	return &res, nil
}

func (f *FunctionsService) Rollback(ctx context.Context, functionName, deploymentId string) error {
	url := fmt.Sprintf("%s/functions/%s/%s/rollback", f.client.baseURL, f.projectID, functionName)
	req := map[string]string{"deploymentId": deploymentId}
	_, _, err := f.client.doRequest(ctx, "POST", url, req)
	return err
}

func (f *FunctionsService) Invoke(ctx context.Context, functionName string, payload any) ([]byte, error) {
	url := fmt.Sprintf("%s/functions/invoke/%s/%s", f.client.baseURL, f.projectID, functionName)
	res, _, err := f.client.doRequest(ctx, "POST", url, payload)
	return res, err
}

type RulesService struct {
	client    *httpClient
	projectID string
}

type RulesResponse struct {
	CurrentRules string           `json:"currentRules"`
	History      []map[string]any `json:"history"`
}

func (r *RulesService) Get(ctx context.Context, serviceType string) (*RulesResponse, error) {
	url := fmt.Sprintf("%s/rules/projects/%s/rules?serviceType=%s", r.client.baseURL, r.projectID, serviceType)
	body, _, err := r.client.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var res RulesResponse
	Unmarshal(body, &res)
	return &res, nil
}

func (r *RulesService) Update(ctx context.Context, serviceType, rawRules string) error {
	url := fmt.Sprintf("%s/rules/projects/%s/rules", r.client.baseURL, r.projectID)
	req := map[string]string{
		"serviceType": serviceType,
		"rawRules":    rawRules,
	}
	_, _, err := r.client.doRequest(ctx, "PUT", url, req)
	return err
}

type BackupService struct {
	client    *httpClient
	projectID string
}

func (b *BackupService) Create(ctx context.Context, backupType string) (*map[string]any, error) {
	url := fmt.Sprintf("%s/backups/%s/create?type=%s", b.client.baseURL, b.projectID, backupType)
	body, _, err := b.client.doRequest(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	var res map[string]any
	Unmarshal(body, &res)
	return &res, nil
}

func (b *BackupService) Restore(ctx context.Context, backupId string) error {
	url := fmt.Sprintf("%s/backups/%s/restore/%s", b.client.baseURL, b.projectID, backupId)
	_, _, err := b.client.doRequest(ctx, "POST", url, nil)
	return err
}
