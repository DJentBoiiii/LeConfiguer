package client

import (
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	url string
	c   *http.Client
}

func NewClient(url string) *Client {
	return &Client{url, &http.Client{}}
}

func (c *Client) do(method, endpoint string, body io.Reader, ct string) ([]byte, error) {
	req, _ := http.NewRequest(method, c.url+endpoint, body)
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, data)
	}
	return data, nil
}

func (c *Client) withFile(method, endpoint, path string, fields map[string]string) ([]byte, error) {
	body, ct, err := c.multipart(path, fields)
	if err != nil {
		return nil, err
	}
	return c.do(method, endpoint, body, ct)
}
