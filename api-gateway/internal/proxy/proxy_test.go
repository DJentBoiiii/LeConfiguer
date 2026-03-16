//go:build ignore

package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewProxyTrimsTrailingSlash(t *testing.T) {
	p := NewProxy("http://storage.local///", "http://index.local//")

	if p.StorageURL != "http://storage.local" {
		t.Fatalf("unexpected storage URL: %s", p.StorageURL)
	}
	if p.IndexingURL != "http://index.local" {
		t.Fatalf("unexpected indexing URL: %s", p.IndexingURL)
	}
}

func TestCopyHeadersRemovesHopByHop(t *testing.T) {
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

func TestForwardRoutesToIndexingForDiff(t *testing.T) {
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
	if got := rr.Header().Get("X-Upstream"); got != "indexing" {
		t.Fatalf("unexpected upstream header: %q", got)
	}
	if body := rr.Body.String(); body != "indexed" {
		t.Fatalf("unexpected response body: %q", body)
	}
}

func TestForwardRoutesToStorageByDefault(t *testing.T) {
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
	if body := rr.Body.String(); body != "stored" {
		t.Fatalf("unexpected response body: %q", body)
	}
}

func TestForwardReturnsBadGatewayOnUpstreamError(t *testing.T) {
	p := NewProxy("http://127.0.0.1:1", "http://127.0.0.1:1")

	req := httptest.NewRequest(http.MethodGet, "/configs", nil)
	rr := httptest.NewRecorder()

	p.Forward(rr, req)

	if rr.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, rr.Code)
	}
}






















































































































}	}		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, rr.Code)	if rr.Code != http.StatusBadGateway {	p.Forward(rr, req)	rr := httptest.NewRecorder()	req := httptest.NewRequest(http.MethodGet, "/configs", nil)	p := NewProxy("http://127.0.0.1:1", "http://127.0.0.1:1")func TestForwardReturnsBadGatewayOnUpstreamError(t *testing.T) {}	}		t.Fatalf("unexpected response body: %q", body)	if body := rr.Body.String(); body != "stored" {	}		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)	if rr.Code != http.StatusOK {	p.Forward(rr, req)	rr := httptest.NewRecorder()	req.Header.Set("X-Test", "hello")	req := httptest.NewRequest(http.MethodGet, "/configs", nil)	p := NewProxy(storage.URL, indexing.URL)	defer indexing.Close()	}))		t.Fatalf("request must not be routed to indexing for regular endpoint")	indexing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {	defer storage.Close()	}))		_, _ = io.WriteString(w, "stored")		w.WriteHeader(http.StatusOK)		}			t.Fatalf("expected header to be forwarded, got %q", got)		if got := r.Header.Get("X-Test"); got != "hello" {		}			t.Fatalf("unexpected path forwarded to storage: %s", r.URL.Path)		if r.URL.Path != "/configs" {	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {func TestForwardRoutesToStorageByDefault(t *testing.T) {}	}		t.Fatalf("unexpected response body: %q", body)	if body := rr.Body.String(); body != "indexed" {	}		t.Fatalf("unexpected upstream header: %q", got)	if got := rr.Header().Get("X-Upstream"); got != "indexing" {	}		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rr.Code)	if rr.Code != http.StatusAccepted {	p.Forward(rr, req)	rr := httptest.NewRecorder()	req := httptest.NewRequest(http.MethodGet, "/diff/abc?v=1", nil)	p := NewProxy(storage.URL, indexing.URL)	defer storage.Close()	}))		t.Fatalf("request must not be routed to storage for diff endpoint")	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {	defer indexing.Close()	}))		_, _ = w.Write([]byte("indexed"))		w.WriteHeader(http.StatusAccepted)		w.Header().Set("X-Upstream", "indexing")		}			t.Fatalf("unexpected query forwarded to indexing: %s", r.URL.RawQuery)		if r.URL.RawQuery != "v=1" {		}			t.Fatalf("unexpected path forwarded to indexing: %s", r.URL.Path)		if r.URL.Path != "/diff/abc" {	indexing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {func TestForwardRoutesToIndexingForDiff(t *testing.T) {}	}		t.Fatalf("expected Transfer-Encoding to be removed, got %q", got)	if got := dst.Get("Transfer-Encoding"); got != "" {	}		t.Fatalf("expected Connection to be removed, got %q", got)	if got := dst.Get("Connection"); got != "" {	}		t.Fatalf("expected X-Test header to be copied, got %q", got)	if got := dst.Get("X-Test"); got != "ok" {	copyHeaders(dst, src)	dst := http.Header{}	src.Set("Transfer-Encoding", "chunked")	src.Set("Connection", "keep-alive")	src.Set("X-Test", "ok")	src := http.Header{}func TestCopyHeadersRemovesHopByHop(t *testing.T) {}	}		t.Fatalf("unexpected indexing URL: %s", p.IndexingURL)	if p.IndexingURL != "http://index.local" {	}		t.Fatalf("unexpected storage URL: %s", p.StorageURL)	if p.StorageURL != "http://storage.local" {	p := NewProxy("http://storage.local///", "http://index.local//")func TestNewProxyTrimsTrailingSlash(t *testing.T) {)	"testing"	"net/http/httptest"	"net/http"
*/
