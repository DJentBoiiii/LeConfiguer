package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"indexing/internal/models"
	storageclient "indexing/internal/storage"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestHandler(t *testing.T, storageURL string) *Handler {
	t.Helper()

	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.ConfigChange{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	var c *storageclient.Client
	if storageURL != "" {
		c = storageclient.NewClient(storageURL)
	}

	return New(db, c)
}

func TestCreateChangeSuccess(t *testing.T) {
	h := newTestHandler(t, "")

	body := bytes.NewBufferString(`{"name":"nginx","type":"json","environment":"dev","action":"create","content":"{\"x\":1}"}`)
	req := httptest.NewRequest(http.MethodPost, "/configs/c1/changes", body)
	req = mux.SetURLVars(req, map[string]string{"id": "c1"})
	rr := httptest.NewRecorder()

	h.CreateChange(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestCreateChangeRejectsInvalidAction(t *testing.T) {
	h := newTestHandler(t, "")

	body := bytes.NewBufferString(`{"action":"invalid","content":"x"}`)
	req := httptest.NewRequest(http.MethodPost, "/configs/c1/changes", body)
	req = mux.SetURLVars(req, map[string]string{"id": "c1"})
	rr := httptest.NewRecorder()

	h.CreateChange(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestListVersionsRedactsContent(t *testing.T) {
	h := newTestHandler(t, "")

	if err := h.db.Create(&models.ConfigChange{ConfigID: "c1", Action: "create", Content: "secret"}).Error; err != nil {
		t.Fatalf("seed change: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/configs/c1/versions", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "c1"})
	rr := httptest.NewRecorder()

	h.ListVersions(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var versions []models.ConfigChange
	if err := json.Unmarshal(rr.Body.Bytes(), &versions); err != nil {
		t.Fatalf("decode versions: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(versions))
	}
	if versions[0].Content != "" {
		t.Fatalf("expected redacted content, got %q", versions[0].Content)
	}
}

func TestDiffReturnsLatestContent(t *testing.T) {
	h := newTestHandler(t, "")
	if err := h.db.Create(&models.ConfigChange{ConfigID: "c1", Action: "update", Content: "latest-content"}).Error; err != nil {
		t.Fatalf("seed change: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/diff/c1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "c1"})
	rr := httptest.NewRecorder()

	h.Diff(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Body.String() != "latest-content" {
		t.Fatalf("unexpected body: %q", rr.Body.String())
	}
}

func TestRollbackDeletesLatestAndUpdatesStorage(t *testing.T) {
	updated := 0
	storageSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/configs/c1" {
			t.Fatalf("unexpected storage call %s %s", r.Method, r.URL.Path)
		}
		updated++
		w.WriteHeader(http.StatusOK)
	}))
	defer storageSrv.Close()

	h := newTestHandler(t, storageSrv.URL)
	if err := h.db.Create(&models.ConfigChange{ConfigID: "c1", Name: "old", Type: "json", Environment: "dev", Action: "update", Content: "old-content"}).Error; err != nil {
		t.Fatalf("seed old change: %v", err)
	}
	if err := h.db.Create(&models.ConfigChange{ConfigID: "c1", Name: "new", Type: "json", Environment: "dev", Action: "update", Content: "new-content"}).Error; err != nil {
		t.Fatalf("seed new change: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/configs/c1/rollback", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "c1"})
	rr := httptest.NewRecorder()

	h.Rollback(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if updated != 1 {
		t.Fatalf("expected storage update call once, got %d", updated)
	}

	var count int64
	if err := h.db.Model(&models.ConfigChange{}).Where("config_id = ?", "c1").Count(&count).Error; err != nil {
		t.Fatalf("count changes: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 version remaining after rollback, got %d", count)
	}
}
