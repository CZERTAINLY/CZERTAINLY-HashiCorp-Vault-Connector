//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	vaultImage     = "hashicorp/vault:latest"
	vaultPort      = "8200/tcp"
	vaultRootToken = "root"
	vaultMount     = "kv-test"
	vaultRole      = "czertainly-test"
	vaultPolicy    = "czertainly-test"
)

// VaultContainer wraps a testcontainers Vault instance.
type VaultContainer struct {
	container testcontainers.Container
	host      string
	port      string
	roleID    string
	secretID  string
}

// SetupVault starts a Vault dev-mode container and configures it for integration tests:
//  1. Enables KV v2 at mount "kv-test/"
//  2. Enables AppRole auth at "approle/"
//  3. Creates policy allowing kv-test + sys/internal/ui/mounts
//  4. Creates an AppRole role bound to the policy
//  5. Fetches and stores role-id and secret-id
func SetupVault(ctx context.Context, t *testing.T) *VaultContainer {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        vaultImage,
		ExposedPorts: []string{vaultPort},
		Env: map[string]string{
			"VAULT_DEV_ROOT_TOKEN_ID":  vaultRootToken,
			"VAULT_DEV_LISTEN_ADDRESS": "0.0.0.0:8200",
		},
		WaitingFor: wait.ForHTTP("/v1/sys/health").
			WithPort("8200/tcp").
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start Vault container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get Vault container host: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, "8200")
	if err != nil {
		t.Fatalf("failed to get Vault container port: %v", err)
	}

	vc := &VaultContainer{
		container: container,
		host:      host,
		port:      mappedPort.Port(),
	}

	t.Logf("Vault started at %s", vc.VaultURL())

	vc.configure(ctx, t)

	return vc
}

// VaultURL returns the base URL for the Vault container.
func (vc *VaultContainer) VaultURL() string {
	return fmt.Sprintf("http://%s:%s", vc.host, vc.port)
}

// RoleID returns the AppRole role-id configured during setup.
func (vc *VaultContainer) RoleID() string { return vc.roleID }

// SecretID returns the AppRole secret-id configured during setup.
func (vc *VaultContainer) SecretID() string { return vc.secretID }

// Mount returns the KV v2 mount name used in tests.
func (vc *VaultContainer) Mount() string { return vaultMount + "/" }

// Terminate stops and removes the Vault container.
func (vc *VaultContainer) Terminate(ctx context.Context) {
	if vc.container != nil {
		_ = vc.container.Terminate(ctx)
	}
}

// configure runs the Vault setup sequence using the HTTP API with the root token.
func (vc *VaultContainer) configure(ctx context.Context, t *testing.T) {
	t.Helper()

	// 1. Enable KV v2 at kv-test/
	// Vault dev mode pre-mounts "secret/" — using a different name avoids 400 Bad Request.
	vc.vaultAPI(t, http.MethodPost, "/v1/sys/mounts/"+vaultMount, map[string]any{
		"type":    "kv",
		"options": map[string]any{"version": "2"},
	})
	t.Log("KV v2 enabled at kv-test/")

	// 2. Enable AppRole auth
	vc.vaultAPI(t, http.MethodPost, "/v1/sys/auth/approle", map[string]any{
		"type": "approle",
	})
	t.Log("AppRole auth enabled")

	// 3. Create policy
	// kv-test/data/* and kv-test/metadata/* for CRUD
	// sys/internal/ui/mounts for DetectKVVersion + ListVisibleMounts (called on every operation)
	policy := `
path "kv-test/data/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
path "kv-test/metadata/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
path "sys/internal/ui/mounts" {
  capabilities = ["read"]
}
`
	vc.vaultAPI(t, http.MethodPost, "/v1/sys/policies/acl/"+vaultPolicy, map[string]any{
		"policy": policy,
	})
	t.Log("Policy czertainly-test created")

	// 4. Create AppRole role bound to the policy
	vc.vaultAPI(t, http.MethodPost, "/v1/auth/approle/role/"+vaultRole, map[string]any{
		"policies": []string{vaultPolicy},
	})
	t.Log("AppRole role czertainly-test created")

	// 5. Fetch role-id
	resp := vc.vaultAPIGet(t, "/v1/auth/approle/role/"+vaultRole+"/role-id")
	vc.roleID = resp["data"].(map[string]any)["role_id"].(string)
	t.Logf("Role ID: %s", vc.roleID)

	// 6. Generate secret-id
	resp = vc.vaultAPI(t, http.MethodPost, "/v1/auth/approle/role/"+vaultRole+"/secret-id", map[string]any{})
	vc.secretID = resp["data"].(map[string]any)["secret_id"].(string)
	t.Log("Secret ID generated")
}

// vaultAPI makes an authenticated Vault API call and returns the parsed response body.
// It fails the test on any HTTP error or non-2xx status.
func (vc *VaultContainer) vaultAPI(t *testing.T, method, path string, body map[string]any) map[string]any {
	t.Helper()

	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("vault API marshal error on %s %s: %v", method, path, err)
	}

	req, err := http.NewRequest(method, vc.VaultURL()+path, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("vault API request build error on %s %s: %v", method, path, err)
	}
	req.Header.Set("X-Vault-Token", vaultRootToken)
	req.Header.Set("Content-Type", "application/json")

	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("vault API call failed on %s %s: %v", method, path, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("vault API %s %s returned %d: %s", method, path, resp.StatusCode, respBody)
	}

	if len(respBody) == 0 {
		return nil
	}
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		t.Fatalf("vault API response parse failed on %s %s: %v (body: %s)", method, path, err, respBody)
	}
	return result
}

// vaultAPIGet is a convenience wrapper for GET requests.
func (vc *VaultContainer) vaultAPIGet(t *testing.T, path string) map[string]any {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, vc.VaultURL()+path, nil)
	if err != nil {
		t.Fatalf("vault GET request build error on %s: %v", path, err)
	}
	req.Header.Set("X-Vault-Token", vaultRootToken)

	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("vault GET failed on %s: %v", path, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("vault GET %s returned %d: %s", path, resp.StatusCode, respBody)
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		t.Fatalf("vault GET response parse failed on %s: %v", path, err)
	}
	return result
}
