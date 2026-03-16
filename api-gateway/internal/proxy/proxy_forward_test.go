package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewProxyTrimsTrailingSlash_Clean(t *testing.T) {
	p := NewProxy("http://storage.local///", "http://index.local//")
	if p.StorageURL != "http://storage.local" {
		t.Fatalf("unexpected storage URL: %s", p.StorageURL)
	}
	if p.IndexingURL != "http://index.local" {
		t.Fatalf("unexpected indexing URL: %s", p.IndexingURL)
	}
}

func TestCopyHeadersRemovesHopByHop_Clean(t *testing.T) {
	src := http.Header{}
	src.Set("X-Test", "ok")
	src.Set("Connection", "keep-alive")
	src.Set("Transfer-Encoding", "chunked")

	dst := http.Header{}
	copyHeaders(dst, src)

	if got := dst.Get("X-Test"); got != "ok" {
		t.Fatalf("expected X-Test header to be copied, got %q", got)
	}
	if got := dst.Get("Connection"); got != "" {
		t.Fatalf("expected Connection to be removed, got %q", got)
	}
	if got := dst.Get("Transfer-Encoding"); got != "" {
		t.Fatalf("expected Transfer-Encoding to be removed, got %q", got)
	}
}

func TestForwardRoutesToIndexingForDiff_Clean(t *testing.T) {
	indexing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/diff/abc" {
			t.Fatalf("unexpected path forwarded to indexing: %s", r.URL.Path)
		}
		if r.URL.RawQuery != "v=1" {
			t.Fatalf("unexpected query forwarded to indexing: %s", r.URL.RawQuery)
		}
		w.Header().Set("X-Upstream", "indexing")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("indexed"))
	}))
	defer indexing.Close()

	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("request must not be routed to storage for diff endpoint")
	}))
	defer storage.Close()

	p := NewProxy(storage.URL, indexing.URL)
	req := httptest.NewRequest(http.MethodGet, "/diff/abc?v=1", nil)
	rr := httptest.NewRecorder()

	p.Forward(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}
	if rr.Header().Get("X-Upstream") != "indexing" {
		t.Fatalf("expected response from indexing upstream")
	}
	if rr.Body.String() != "indexed" {
		t.Fatalf("unexpected response body: %q", rr.Body.String())
	}
}

func TestForwardRoutesToStorageByDefault_Clean(t *testing.T) {
	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/configs" {
			t.Fatalf("unexpected path forwarded to storage: %s", r.URL.Path)
		}
		if got := r.Header.Get("X-Test"); got != "hello" {
			t.Fatalf("expected header to be forwarded, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "stored")
	}))
	defer storage.Close()

	indexing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("request must not be routed to indexing for regular endpoint")
	}))
	defer indexing.Close()

	p := NewProxy(storage.URL, indexing.URL)
	req := httptest.NewRequest(http.MethodGet, "/configs", nil)
	req.Header.Set("X-Test", "hello")
	rr := httptest.NewRecorder()

	p.Forward(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Body.String() != "stored" {
		t.Fatalf("unexpected response body: %q", rr.Body.String())
	}
}

func TestForwardReturnsBadGatewayOnUpstreamError_Clean(t *testing.T) {
	p := NewProxy("http://127.0.0.1:1", "http://127.0.0.1:1")
	req := httptest.NewRequest(http.MethodGet, "/configs", nil)
	rr := httptest.NewRecorder()

	p.Forward(rr, req)

	if rr.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, rr.Code)
	}
}
