package indexing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendChangeSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/configs/cfg-1/changes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	err := c.SendChange(context.Background(), "cfg-1", ChangeRequest{Action: "create", Content: "{}"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSendChangeNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	err := c.SendChange(context.Background(), "cfg-1", ChangeRequest{Action: "create", Content: "{}"})
	if err == nil {
		t.Fatalf("expected error on non-2xx response")
	}
}

func TestDeleteConfigSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/configs/cfg-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	if err := c.DeleteConfig(context.Background(), "cfg-1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
