package client

func (c *Client) Upload(endpoint, path string, fields map[string]string) ([]byte, error) {
	return c.withFile("POST", endpoint, path, fields)
}

func (c *Client) Update(endpoint, path string, fields map[string]string) ([]byte, error) {
	return c.withFile("PUT", endpoint, path, fields)
}

func (c *Client) request(method, endpoint string) ([]byte, error) {
	return c.do(method, endpoint, nil, "")
}

func (c *Client) Get(endpoint string) ([]byte, error) {
	return c.request("GET", endpoint)
}

func (c *Client) Delete(endpoint string) ([]byte, error) {
	return c.request("DELETE", endpoint)
}

func (c *Client) Post(endpoint string) ([]byte, error) {
	return c.request("POST", endpoint)
}
