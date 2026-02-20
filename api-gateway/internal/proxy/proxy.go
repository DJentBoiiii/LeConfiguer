package proxy

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// Proxy forwards requests to the Config Storage service.
type Proxy struct {
	BaseURL string
	Client  *http.Client
}

// NewProxy creates a new Proxy instance.
func NewProxy(baseURL string) *Proxy {
	return &Proxy{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Client:  &http.Client{},
	}
}

func (p *Proxy) Forward(w http.ResponseWriter, r *http.Request) {
	upstreamURL := p.BaseURL + r.URL.Path
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
		log.Printf("Error reaching config storage at %s: %v", upstreamURL, err)
		http.Error(w, "failed to reach config storage", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	log.Printf("Received response from config storage: status=%d", resp.StatusCode)

	copyHeaders(w.Header(), resp.Header)
	removeHopByHopHeaders(w.Header())
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
	removeHopByHopHeaders(dst)
}

func removeHopByHopHeaders(header http.Header) {
	header.Del("Connection")
	header.Del("Proxy-Connection")
	header.Del("Keep-Alive")
	header.Del("Proxy-Authenticate")
	header.Del("Proxy-Authorization")
	header.Del("Te")
	header.Del("Trailer")
	header.Del("Transfer-Encoding")
	header.Del("Upgrade")
}
