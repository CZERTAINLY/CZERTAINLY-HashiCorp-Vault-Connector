//go:build integration

package integration

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// TestSecretProvider_ListVaultAttributes verifies that GET /v1/secretProvider/vaults/attributes
// returns 200 with the expected attribute UUIDs in the response body.
// No Vault interaction required — this is a static response from the connector.
func TestSecretProvider_ListVaultAttributes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	h := NewSecretTestHarness(t)
	defer h.Cleanup()

	resp, err := h.Client.Get(h.URL() + "/v1/secretProvider/vaults/attributes")
	if err != nil {
		t.Fatalf("GET /v1/secretProvider/vaults/attributes failed: %v", err)
	}
	defer resp.Body.Close()

	assertStatus(t, resp, http.StatusOK)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	assertBodyContainsUUIDs(t, body,
		AttrUUIDVaultURI,
		AttrUUIDNamespace,
		AttrUUIDCredentialType,
		GroupUUIDVaultCredentials,
	)

	t.Log("ListVaultAttributes returned expected attribute UUIDs")
}

// TestSecretProvider_CheckVaultConnection verifies POST /v1/secretProvider/vaults
// with valid and invalid AppRole credentials.
//
// Note: the body is a bare []RequestAttribute JSON array (not a DTO struct wrapper).
// Valid credentials → 204 No Content.
// Invalid secret-id → 400 Bad Request.
func TestSecretProvider_CheckVaultConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	h := NewSecretTestHarness(t)
	defer h.Cleanup()

	tests := []struct {
		name           string
		secretID       string
		expectedStatus int
	}{
		{
			name:           "valid AppRole credentials",
			secretID:       h.Vault.SecretID(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid secret-id",
			secretID:       "00000000-0000-0000-0000-000000000000",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewRequestBuilder(
				h.Vault.VaultURL(),
				h.Vault.RoleID(),
				tt.secretID,
				h.Vault.Mount(),
			)

			body := builder.BuildCheckConnectionBody()

			req, err := http.NewRequest(http.MethodPost, h.URL()+"/v1/secretProvider/vaults", bytesReader(body))
			if err != nil {
				t.Fatalf("failed to build request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := h.Client.Do(req)
			if err != nil {
				t.Fatalf("POST /v1/secretProvider/vaults failed: %v", err)
			}
			defer resp.Body.Close()

			assertStatus(t, resp, tt.expectedStatus)
			t.Logf("%s: got %d", tt.name, resp.StatusCode)
		})
	}
}

// TestSecretProvider_SecretCRUD runs a full lifecycle test for each secret type:
// Create → Read (verify) → Update → Read (verify update) → Delete → Read (expect 404).
//
// Each run uses a UUID-suffixed secret name to avoid collisions across parallel runs.
func TestSecretProvider_SecretCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	h := NewSecretTestHarness(t)
	defer h.Cleanup()

	type secretCase struct {
		name          string
		secretType    string
		createContent map[string]any
		updateContent map[string]any
		// createValues are the strings expected in the read-after-create response body
		createValues []string
		// updateValues are the strings expected in the read-after-update response body
		updateValues []string
	}

	tests := []secretCase{
		{
			name:          "generic secret",
			secretType:    "generic",
			createContent: GenericSecret("initial-generic-value"),
			updateContent: GenericSecret("updated-generic-value"),
			createValues:  []string{"initial-generic-value"},
			updateValues:  []string{"updated-generic-value"},
		},
		{
			name:          "basicAuth secret",
			secretType:    "basicAuth",
			createContent: BasicAuthSecret("alice", "password123"),
			updateContent: BasicAuthSecret("alice", "newpassword456"),
			createValues:  []string{"alice", "password123"},
			updateValues:  []string{"alice", "newpassword456"},
		},
		{
			name:          "apiKey secret",
			secretType:    "apiKey",
			createContent: ApiKeySecret("initial-api-key-abc"),
			updateContent: ApiKeySecret("updated-api-key-xyz"),
			createValues:  []string{"initial-api-key-abc"},
			updateValues:  []string{"updated-api-key-xyz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unique name per test run to prevent collisions
			secretName := fmt.Sprintf("test-%s", uuid.New().String())
			builder := h.Builder()

			// --- 1. Create ---
			t.Logf("Creating secret %q (%s)", secretName, tt.secretType)
			resp := doRequest(t, h.Client, http.MethodPost,
				h.URL()+"/v1/secretProvider/secrets",
				builder.BuildCreateSecretBody(secretName, tt.createContent))
			assertStatus(t, resp, http.StatusCreated)
			resp.Body.Close()
			t.Log("Create: 201")

			// --- 2. Read after create ---
			t.Log("Reading secret after create")
			resp = doRequest(t, h.Client, http.MethodPost,
				h.URL()+"/v1/secretProvider/secrets/content",
				builder.BuildReadSecretBody(secretName, tt.secretType))
			assertStatus(t, resp, http.StatusOK)
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			assertBodyContains(t, body, tt.createValues...)
			t.Logf("Read after create: 200, values present")

			// --- 3. Update ---
			t.Log("Updating secret")
			resp = doRequest(t, h.Client, http.MethodPut,
				h.URL()+"/v1/secretProvider/secrets",
				builder.BuildUpdateSecretBody(secretName, tt.updateContent))
			assertStatus(t, resp, http.StatusOK)
			resp.Body.Close()
			t.Log("Update: 200")

			// --- 4. Read after update ---
			t.Log("Reading secret after update")
			resp = doRequest(t, h.Client, http.MethodPost,
				h.URL()+"/v1/secretProvider/secrets/content",
				builder.BuildReadSecretBody(secretName, tt.secretType))
			assertStatus(t, resp, http.StatusOK)
			body, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			assertBodyContains(t, body, tt.updateValues...)
			t.Logf("Read after update: 200, updated values present")

			// --- 5. Delete ---
			t.Log("Deleting secret")
			resp = doRequest(t, h.Client, http.MethodDelete,
				h.URL()+"/v1/secretProvider/secrets",
				builder.BuildDeleteSecretBody(secretName))
			resp.Body.Close()
			assertStatus(t, resp, http.StatusNoContent)
			t.Log("Delete: 204")

			// --- 6. Read after delete — expect 404 ---
			// Vault returns 404 → ErrNotFound → notfound() → HTTP 404
			t.Log("Reading secret after delete — expecting 404")
			resp = doRequest(t, h.Client, http.MethodPost,
				h.URL()+"/v1/secretProvider/secrets/content",
				builder.BuildReadSecretBody(secretName, tt.secretType))
			resp.Body.Close()
			assertNotFound(t, resp)
			t.Log("Read after delete: 404")
		})
	}
}
