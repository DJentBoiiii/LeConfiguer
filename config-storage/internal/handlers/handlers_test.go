package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"config-storage/internal/models"
	"config-storage/internal/storage"

	indexingclient "config-storage/internal/indexing"

	"github.com/gorilla/mux"
)

type mockStorage struct {
	createFn   func(config *models.Config) error
	getFn      func(id string) (*models.Config, error)
	updateFn   func(id string, config *models.Config) error
	deleteFn   func(id string) error
	listFn     func() ([]*models.Config, error)
	downloadFn func(id string) (*models.Config, io.ReadCloser, error)
}

func (m *mockStorage) Create(config *models.Config) error {
	if m.createFn == nil {
		return nil
	}
	return m.createFn(config)
}

func (m *mockStorage) Get(id string) (*models.Config, error) {
	if m.getFn == nil {
		return nil, nil
	}
	return m.getFn(id)
}

func (m *mockStorage) Update(id string, config *models.Config) error {
	if m.updateFn == nil {
		return nil
	}
	return m.updateFn(id, config)
}

func (m *mockStorage) Delete(id string) error {
	if m.deleteFn == nil {
		return nil
	}
	return m.deleteFn(id)
}

func (m *mockStorage) List() ([]*models.Config, error) {
	if m.listFn == nil {
		return nil, nil
	}
	return m.listFn()
}

func (m *mockStorage) Download(id string) (*models.Config, io.ReadCloser, error) {
	if m.downloadFn == nil {
		return nil, nil, nil
	}
	return m.downloadFn(id)
}

func multipartBody(t *testing.T, fields map[string]string, filename string, content []byte) (string, *bytes.Buffer) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatalf("write field: %v", err)
		}
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write file part: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	return writer.FormDataContentType(), body
}

func TestUploadConfigSuccess(t *testing.T) {
	storageMock := &mockStorage{}

	h := NewHandler(storageMock, nil)
	contentType, body := multipartBody(t, map[string]string{
		"name":        "nginx",
		"type":        "json",
		"environment": "dev",
	}, "config.json", []byte(`{"x":1}`))

	req := httptest.NewRequest(http.MethodPost, "/configs", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	h.UploadConfig(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var got models.Config
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID == "" {
		t.Fatalf("expected generated id")
	}
	if got.Name != "nginx.json" {
		t.Fatalf("unexpected name: %s", got.Name)
	}
}

func TestGetConfigNotFound(t *testing.T) {
	h := NewHandler(&mockStorage{getFn: func(id string) (*models.Config, error) {
		return nil, storage.ErrNotFound
	}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/configs/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	h.GetConfig(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestDeleteConfigCallsIndexing(t *testing.T) {
	indexChanges := 0
	deleteConfig := 0
	indexServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/configs/abc/changes":
			indexChanges++
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodDelete && r.URL.Path == "/configs/abc":
			deleteConfig++
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected indexing call %s %s", r.Method, r.URL.Path)
		}
	}))
	defer indexServer.Close()

	h := NewHandler(&mockStorage{
		getFn: func(id string) (*models.Config, error) {
			return &models.Config{ID: id, Name: "nginx", Type: "json", Environment: "dev"}, nil
		},
	}, indexingclient.NewClient(indexServer.URL))

	req := httptest.NewRequest(http.MethodDelete, "/configs/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	h.DeleteConfig(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if indexChanges != 1 || deleteConfig != 1 {
		t.Fatalf("expected indexing calls (changes=1, delete=1), got (%d, %d)", indexChanges, deleteConfig)
	}
}

func TestDownloadConfigReturnsFile(t *testing.T) {
	h := NewHandler(&mockStorage{downloadFn: func(id string) (*models.Config, io.ReadCloser, error) {
		return &models.Config{Name: "nginx", Type: "json"}, io.NopCloser(bytes.NewReader([]byte("hello"))), nil
	}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/configs/download/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	h.DownloadConfig(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if got := rr.Header().Get("Content-Disposition"); got == "" {
		t.Fatalf("expected content disposition header")
	}
	if rr.Body.String() != "hello" {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}
