package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) UpdateConfig(ctx context.Context, id, name, configType, environment, content string) error {
	url := fmt.Sprintf("%s/configs/%s", c.baseURL, id)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("name", name)
	_ = writer.WriteField("type", configType)
	_ = writer.WriteField("environment", environment)

	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, bytes.NewReader([]byte(content)))
	if err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("config-storage returned status %d", resp.StatusCode)
	}

	return nil
}
