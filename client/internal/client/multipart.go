package client

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func (c *Client) multipart(path string, fields map[string]string) (*bytes.Buffer, string, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	for k, v := range fields {
		w.WriteField(k, v)
	}
	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			return nil, "", err
		}
		defer f.Close()
		p, _ := w.CreateFormFile("file", filepath.Base(path))
		io.Copy(p, f)
	}
	w.Close()
	return buf, w.FormDataContentType(), nil
}
