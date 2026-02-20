package proxy

import (
	"net/http"
	"strings"
)

type Proxy struct {
	StorageURL  string
	IndexingURL string
	Client      *http.Client
}

func NewProxy(storageURL, indexingURL string) *Proxy {
	return &Proxy{
		StorageURL:  strings.TrimRight(storageURL, "/"),
		IndexingURL: strings.TrimRight(indexingURL, "/"),
		Client:      &http.Client{},
	}
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
