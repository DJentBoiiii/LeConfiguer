package proxy

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request) {
	baseURL := p.StorageURL

	// Route to indexing service for specific endpoints
	if strings.HasPrefix(r.URL.Path, "/diff/") ||
		strings.HasSuffix(r.URL.Path, "/versions") ||
		strings.Contains(r.URL.Path, "/versions/") ||
		strings.HasSuffix(r.URL.Path, "/rollback") {
		baseURL = p.IndexingURL
	}

	upstreamURL := baseURL + r.URL.Path
	if r.URL.RawQuery != "" {
		upstreamURL += "?" + r.URL.RawQuery
	}

	log.Printf("Proxying %s %s to %s", r.Method, r.URL.Path, upstreamURL)

	req, err := http.NewRequestWithContext(r.Context(), r.Method, upstreamURL, r.Body)
	if err != nil {
		log.Printf("Error creating upstream request: %v", err)
		http.Error(w, "failed to create upstream request", http.StatusBadGateway)
		return
	}

	copyHeaders(req.Header, r.Header)

	resp, err := p.Client.Do(req)
	if err != nil {
		log.Printf("Error reaching upstream at %s: %v", upstreamURL, err)
		http.Error(w, "failed to reach upstream service", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	log.Printf("Received response: status=%d", resp.StatusCode)

	copyHeaders(w.Header(), resp.Header)
	removeHopByHopHeaders(w.Header())
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
