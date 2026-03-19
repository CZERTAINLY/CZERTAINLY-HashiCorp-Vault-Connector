//go:build integration

package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"CZERTAINLY-HashiCorp-Vault-Connector/internal/secret"
)

// SecretTestHarness holds all infrastructure for secret provider integration tests.
type SecretTestHarness struct {
	Vault  *VaultContainer
	Server *httptest.Server
	Client *http.Client

	ctx context.Context
	t   *testing.T
}

// NewSecretTestHarness starts a Vault container, configures it (KV v2 + AppRole),
// and starts the connector's secret provider in-process via httptest.Server.
//
// Usage:
//
//	h := NewSecretTestHarness(t)
//	defer h.Cleanup()
func NewSecretTestHarness(t *testing.T) *SecretTestHarness {
	t.Helper()

	ctx := context.Background()

	// Start and configure Vault container
	vault := SetupVault(ctx, t)

	// Start the connector's secret provider in-process.
	// IMPORTANT: secret.New() returns a value, not a pointer.
	// MuxRouter() has a pointer receiver, so we must assign to a named variable first —
	// calling MuxRouter() on a temporary (non-addressable) value would panic.
	s := secret.New()
	server := httptest.NewServer(s.MuxRouter())
	t.Logf("Connector secret provider started at %s", server.URL)

	return &SecretTestHarness{
		Vault:  vault,
		Server: server,
		Client: newHTTPClient(),
		ctx:    ctx,
		t:      t,
	}
}

// Cleanup stops the httptest.Server and terminates the Vault container.
// Call with defer immediately after NewSecretTestHarness.
func (h *SecretTestHarness) Cleanup() {
	h.t.Helper()
	h.Server.Close()
	h.Vault.Terminate(h.ctx)
}

// Builder returns a RequestBuilder pre-configured with the harness's Vault credentials
// and the kv-test/ mount.
func (h *SecretTestHarness) Builder() *RequestBuilder {
	return NewRequestBuilder(
		h.Vault.VaultURL(),
		h.Vault.RoleID(),
		h.Vault.SecretID(),
		h.Vault.Mount(),
	)
}

// URL returns the base URL of the connector's httptest.Server.
func (h *SecretTestHarness) URL() string {
	return h.Server.URL
}
