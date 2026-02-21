package indexing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client sends change events to the indexing service.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// ChangeRequest is the payload for the indexing service.
type ChangeRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Environment string `json:"environment"`
	Action      string `json:"action"`
	Content     string `json:"content"`
}

// NewClient creates a new indexing client.
func NewClient(baseURL string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendChange sends a change request to the indexing service.
func (c *Client) SendChange(ctx context.Context, configID string, req ChangeRequest) error {
	url := fmt.Sprintf("%s/configs/%s/changes", c.baseURL, configID)

	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("indexing service returned status %d", resp.StatusCode)
	}

	return nil
}

// DeleteConfig deletes all records for a config from the indexing service.
func (c *Client) DeleteConfig(ctx context.Context, configID string) error {
	url := fmt.Sprintf("%s/configs/%s", c.baseURL, configID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("indexing service returned status %d", resp.StatusCode)
	}
	return nil
}
