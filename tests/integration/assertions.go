//go:build integration

package integration

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// assertStatus fails the test if the response status code does not match expected.
// It reads and logs the body on failure for easier debugging.
func assertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected HTTP %d, got %d\nBody: %s", expected, resp.StatusCode, body)
	}
}

// assertBodyContains fails the test if the response body does not contain all expected strings.
func assertBodyContains(t *testing.T, body []byte, expected ...string) {
	t.Helper()
	bodyStr := string(body)
	for _, e := range expected {
		if !strings.Contains(bodyStr, e) {
			t.Fatalf("expected body to contain %q\nBody: %s", e, bodyStr)
		}
	}
}

// assertBodyContainsUUIDs fails the test if the response body does not contain all given UUIDs.
func assertBodyContainsUUIDs(t *testing.T, body []byte, uuids ...string) {
	t.Helper()
	assertBodyContains(t, body, uuids...)
}

// assertNotFound fails the test if the response is not HTTP 404.
func assertNotFound(t *testing.T, resp *http.Response) {
	t.Helper()
	assertStatus(t, resp, http.StatusNotFound)
}
