package storage

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateConfigSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/configs/cfg-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		mediaType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(mediaType, "multipart/form-data;") {
			t.Fatalf("expected multipart content-type, got %q", mediaType)
		}

		if err := r.ParseMultipartForm(2 << 20); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		if got := r.FormValue("name"); got != "nginx" {
			t.Fatalf("unexpected name field: %q", got)
		}
		if got := r.FormValue("type"); got != "json" {
			t.Fatalf("unexpected type field: %q", got)
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		_ = file.Close()

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	err := c.UpdateConfig(context.Background(), "cfg-1", "nginx", "json", "dev", `{"x":1}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateConfigNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	err := c.UpdateConfig(context.Background(), "cfg-1", "nginx", "json", "dev", `{"x":1}`)
	if err == nil {
		t.Fatalf("expected error on non-2xx response")
	}
}

func TestMultipartWriterBuildsValidBoundary(t *testing.T) {
	body := &strings.Builder{}
	w := multipart.NewWriter(body)
	_ = w.WriteField("k", "v")
	if err := w.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	if !strings.Contains(w.FormDataContentType(), "boundary=") {
		t.Fatalf("expected boundary in content type")
	}
}
