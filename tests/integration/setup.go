//go:build integration

package integration

import (
	"bytes"
	"net/http"
	"testing"
	"time"
)

// newHTTPClient returns an HTTP client suitable for integration tests.
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

// bytesReader wraps a []byte as an io.Reader for use in http.NewRequest.
func bytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}

// doRequest performs an HTTP request with a JSON body and Content-Type header.
// It fails the test if the request itself errors (not on non-2xx status).
func doRequest(t *testing.T, client *http.Client, method, url string, body []byte) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to build %s %s: %v", method, url, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s %s failed: %v", method, url, err)
	}
	return resp
}
