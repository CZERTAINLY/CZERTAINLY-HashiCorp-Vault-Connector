//go:build integration

package integration

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// TestSecretProvider_CrossTypeRead verifies the behavior when reading a secret
// with a different type than it was created with.
//
// Cross-type compatibility:
//   - generic ↔ apiKey: COMPATIBLE (both stored as {"content": value})
//   - basicAuth ↔ generic: INCOMPATIBLE (different key structure)
//   - basicAuth ↔ apiKey: INCOMPATIBLE (different key structure)
func TestSecretProvider_CrossTypeRead(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	h := NewSecretTestHarness(t)
	defer h.Cleanup()

	tests := []struct {
		name           string
		createType     string
		createContent  map[string]any
		readType       string
		expectedStatus int
	}{
		// Compatible: generic and apiKey both use {"content": value}
		{
			name:           "create generic read as apiKey",
			createType:     "generic",
			createContent:  GenericSecret("cross-type-value"),
			readType:       "apiKey",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "create apiKey read as generic",
			createType:     "apiKey",
			createContent:  ApiKeySecret("cross-type-api-key"),
			readType:       "generic",
			expectedStatus: http.StatusOK,
		},
		// Incompatible: basicAuth uses {"username","password"}, others use {"content"}
		{
			name:           "create basicAuth read as generic",
			createType:     "basicAuth",
			createContent:  BasicAuthSecret("user1", "pass1"),
			readType:       "generic",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "create basicAuth read as apiKey",
			createType:     "basicAuth",
			createContent:  BasicAuthSecret("user2", "pass2"),
			readType:       "apiKey",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "create generic read as basicAuth",
			createType:     "generic",
			createContent:  GenericSecret("some-value"),
			readType:       "basicAuth",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "create apiKey read as basicAuth",
			createType:     "apiKey",
			createContent:  ApiKeySecret("some-api-key"),
			readType:       "basicAuth",
			expectedStatus: http.StatusInternalServerError,
		},
		// jwtToken/secretKey/apiKey/generic are all content-based — cross-compatible
		{
			name:           "create jwtToken read as apiKey",
			createType:     "jwtToken",
			createContent:  JwtTokenSecret("test.jwt.value"),
			readType:       "apiKey",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "create secretKey read as jwtToken",
			createType:     "secretKey",
			createContent:  SecretKeySecret("c2VjcmV0a2V5dmFsdWU="),
			readType:       "jwtToken",
			expectedStatus: http.StatusOK,
		},
		// privateKey: content-based, but reading non-PEM content as privateKey fails PEM validation
		{
			name:           "create apiKey read as privateKey",
			createType:     "apiKey",
			createContent:  ApiKeySecret("not-a-pem-value"),
			readType:       "privateKey",
			expectedStatus: http.StatusInternalServerError, // not PEM format → ErrNotDeclaredType → 500
		},
		{
			name:           "create privateKey read as apiKey",
			createType:     "privateKey",
			// base64 encoding of a structurally valid PEM block accepted by pem.Decode
			createContent:  PrivateKeySecret("LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBMlozVlM1SkpjZHMzCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg=="),
			readType:       "apiKey",
			expectedStatus: http.StatusOK, // privateKey stored as {"content": pem}, apiKey reads "content" key → 200
		},
		// keyValue: no "content" key — incompatible with all content-based types
		{
			name:           "create keyValue read as generic",
			createType:     "keyValue",
			createContent:  KeyValueSecret(map[string]any{"k1": "v1", "k2": "v2"}),
			readType:       "generic",
			expectedStatus: http.StatusInternalServerError, // fromGenericPayload expects exactly 1 key; keyValue has 2 → 500
		},
		{
			name:           "create keyValue read as apiKey",
			createType:     "keyValue",
			createContent:  KeyValueSecret(map[string]any{"somekey": "someval"}),
			readType:       "apiKey",
			expectedStatus: http.StatusInternalServerError, // no "content" key → ErrNotDeclaredType → 500
		},
		// keyStore: 3 keys — incompatible with single-key types
		{
			name:           "create keyStore read as generic",
			createType:     "keyStore",
			createContent:  KeyStoreSecret("ZmFrZWtleXN0b3Jl", "pass", "JKS"),
			readType:       "generic",
			expectedStatus: http.StatusInternalServerError, // fromGenericPayload expects 1 key; keyStore has 3 → 500
		},
		{
			name:           "create apiKey read as keyStore",
			createType:     "apiKey",
			createContent:  ApiKeySecret("some-api-key"),
			readType:       "keyStore",
			expectedStatus: http.StatusInternalServerError, // missing password and key-store-type → ErrNotDeclaredType → 500
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secretName := fmt.Sprintf("cross-%s", uuid.New().String())
			builder := h.Builder()

			// Create the secret with the specified type
			resp := doRequest(t, h.Client, http.MethodPost,
				h.URL()+"/v1/secretProvider/secrets",
				builder.BuildCreateSecretBody(secretName, tt.createContent))
			assertStatus(t, resp, http.StatusCreated)
			resp.Body.Close()

			// Read it with a different type
			resp = doRequest(t, h.Client, http.MethodPost,
				h.URL()+"/v1/secretProvider/secrets/content",
				builder.BuildReadSecretBody(secretName, tt.readType))
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Fatalf("cross-type read: create(%s) read(%s): expected HTTP %d, got %d\nBody: %s",
					tt.createType, tt.readType, tt.expectedStatus, resp.StatusCode, body)
			}
			t.Logf("create(%s) read(%s): got %d (expected)", tt.createType, tt.readType, resp.StatusCode)
		})
	}
}

// TestSecretProvider_DuplicateCreate verifies that creating a secret with
// a name that already exists returns 412 Precondition Failed.
func TestSecretProvider_DuplicateCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	h := NewSecretTestHarness(t)
	defer h.Cleanup()

	secretName := fmt.Sprintf("dup-%s", uuid.New().String())
	builder := h.Builder()

	// First create: should succeed
	resp := doRequest(t, h.Client, http.MethodPost,
		h.URL()+"/v1/secretProvider/secrets",
		builder.BuildCreateSecretBody(secretName, GenericSecret("original")))
	assertStatus(t, resp, http.StatusCreated)
	resp.Body.Close()
	t.Log("First create: 201")

	// Second create with same name: should fail with 412
	resp = doRequest(t, h.Client, http.MethodPost,
		h.URL()+"/v1/secretProvider/secrets",
		builder.BuildCreateSecretBody(secretName, GenericSecret("duplicate")))
	assertStatus(t, resp, http.StatusPreconditionFailed)
	resp.Body.Close()
	t.Log("Duplicate create: 412")
}

// TestSecretProvider_NotFoundOperations verifies behavior when operating on
// secrets that do not exist.
//
//   - Read non-existent: 404
//   - Update non-existent: 404 (update reads before writing; fails with not-found)
//   - Delete non-existent: 204 (Vault KV delete is idempotent)
func TestSecretProvider_NotFoundOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	h := NewSecretTestHarness(t)
	defer h.Cleanup()

	builder := h.Builder()

	t.Run("read non-existent", func(t *testing.T) {
		name := fmt.Sprintf("nonexistent-%s", uuid.New().String())
		resp := doRequest(t, h.Client, http.MethodPost,
			h.URL()+"/v1/secretProvider/secrets/content",
			builder.BuildReadSecretBody(name, "generic"))
		assertStatus(t, resp, http.StatusNotFound)
		resp.Body.Close()
		t.Log("Read non-existent: 404")
	})

	t.Run("update non-existent", func(t *testing.T) {
		name := fmt.Sprintf("nonexistent-%s", uuid.New().String())
		resp := doRequest(t, h.Client, http.MethodPut,
			h.URL()+"/v1/secretProvider/secrets",
			builder.BuildUpdateSecretBody(name, GenericSecret("value")))
		assertStatus(t, resp, http.StatusNotFound)
		resp.Body.Close()
		t.Log("Update non-existent: 404")
	})

	t.Run("delete non-existent", func(t *testing.T) {
		name := fmt.Sprintf("nonexistent-%s", uuid.New().String())
		resp := doRequest(t, h.Client, http.MethodDelete,
			h.URL()+"/v1/secretProvider/secrets",
			builder.BuildDeleteSecretBody(name))
		assertStatus(t, resp, http.StatusNoContent)
		resp.Body.Close()
		t.Log("Delete non-existent: 204 (idempotent)")
	})
}

// TestSecretProvider_UpdateChangesType verifies that updating a secret
// with a different type replaces the stored content entirely.
// After update, the secret should be readable with the new type.
func TestSecretProvider_UpdateChangesType(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	h := NewSecretTestHarness(t)
	defer h.Cleanup()

	secretName := fmt.Sprintf("type-change-%s", uuid.New().String())
	builder := h.Builder()

	// Create as generic
	resp := doRequest(t, h.Client, http.MethodPost,
		h.URL()+"/v1/secretProvider/secrets",
		builder.BuildCreateSecretBody(secretName, GenericSecret("initial-value")))
	assertStatus(t, resp, http.StatusCreated)
	resp.Body.Close()
	t.Log("Created as generic: 201")

	// Read as generic: should succeed
	resp = doRequest(t, h.Client, http.MethodPost,
		h.URL()+"/v1/secretProvider/secrets/content",
		builder.BuildReadSecretBody(secretName, "generic"))
	assertStatus(t, resp, http.StatusOK)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	assertBodyContains(t, body, "initial-value")
	t.Log("Read as generic: 200")

	// Update to basicAuth
	resp = doRequest(t, h.Client, http.MethodPut,
		h.URL()+"/v1/secretProvider/secrets",
		builder.BuildUpdateSecretBody(secretName, BasicAuthSecret("newuser", "newpass")))
	assertStatus(t, resp, http.StatusOK)
	resp.Body.Close()
	t.Log("Updated to basicAuth: 200")

	// Read as basicAuth: should succeed with new credentials
	resp = doRequest(t, h.Client, http.MethodPost,
		h.URL()+"/v1/secretProvider/secrets/content",
		builder.BuildReadSecretBody(secretName, "basicAuth"))
	assertStatus(t, resp, http.StatusOK)
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	assertBodyContains(t, body, "newuser", "newpass")
	t.Log("Read as basicAuth after update: 200")

	// Read as old type (generic) should now fail — storage changed to {username, password}
	// fromGenericPayload expects exactly 1 key; basicAuth has 2 → ErrNotDeclaredType → 500
	resp = doRequest(t, h.Client, http.MethodPost,
		h.URL()+"/v1/secretProvider/secrets/content",
		builder.BuildReadSecretBody(secretName, "generic"))
	if resp.StatusCode == http.StatusOK {
		body, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Fatalf("expected non-200 when reading basicAuth as generic, got 200\nBody: %s", body)
	}
	resp.Body.Close()
	t.Logf("Read as old type (generic) after basicAuth update: %d (expected non-200)", resp.StatusCode)
}

// TestSecretProvider_SecretPathPrefix verifies that secrets can be created
// and retrieved when a relative secret path is specified.
func TestSecretProvider_SecretPathPrefix(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	h := NewSecretTestHarness(t)
	defer h.Cleanup()

	secretName := fmt.Sprintf("pfx-%s", uuid.New().String())

	// Build with a sub-path: secret will be stored at kv-test/subdir/secretName
	builder := h.Builder().WithSecretPath("subdir")

	// Create with path prefix
	resp := doRequest(t, h.Client, http.MethodPost,
		h.URL()+"/v1/secretProvider/secrets",
		builder.BuildCreateSecretBody(secretName, GenericSecret("prefixed-value")))
	assertStatus(t, resp, http.StatusCreated)
	resp.Body.Close()
	t.Log("Create with secret path prefix: 201")

	// Read with same path prefix: should find it
	resp = doRequest(t, h.Client, http.MethodPost,
		h.URL()+"/v1/secretProvider/secrets/content",
		builder.BuildReadSecretBody(secretName, "generic"))
	assertStatus(t, resp, http.StatusOK)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	assertBodyContains(t, body, "prefixed-value")
	t.Log("Read with secret path prefix: 200")

	// Read without path prefix: should NOT find it (different Vault path)
	builderNoPrefix := h.Builder() // no WithSecretPath
	resp = doRequest(t, h.Client, http.MethodPost,
		h.URL()+"/v1/secretProvider/secrets/content",
		builderNoPrefix.BuildReadSecretBody(secretName, "generic"))
	assertStatus(t, resp, http.StatusNotFound)
	resp.Body.Close()
	t.Log("Read without path prefix: 404 (correct — different Vault path)")

	// Delete with same path prefix
	resp = doRequest(t, h.Client, http.MethodDelete,
		h.URL()+"/v1/secretProvider/secrets",
		builder.BuildDeleteSecretBody(secretName))
	assertStatus(t, resp, http.StatusNoContent)
	resp.Body.Close()
	t.Log("Delete with secret path prefix: 204")
}
