# Integration Tests — Secret Provider

End-to-end integration tests for the HashiCorp Vault Connector's Secret Provider (`/v1/secretProvider/*`).

## Overview

Tests make real HTTP requests to the connector running **in-process** via `httptest.Server`,
backed by a real HashiCorp Vault instance in a **testcontainer**.

```
Test code
   │  HTTP (net/http)
   ▼
httptest.Server   ← secret.New().MuxRouter() — connector in-process
   │
   │  Vault SDK (AppRole auth + KV v2)
   ▼
Vault container   ← hashicorp/vault:latest, dev mode, random port
```

No database is required — the secret provider has no DB dependency.

## Infrastructure

| Service | Image | Purpose |
|---------|-------|---------|
| Vault | `hashicorp/vault:latest` | Real KV v2 + AppRole auth for CRUD tests |

Vault configuration applied automatically by the test harness:
- KV v2 enabled at mount `kv-test/`
- AppRole auth enabled at `approle/`
- Policy `czertainly-test`: read/write/delete on `kv-test/data/*`, `kv-test/metadata/*`, and `sys/internal/ui/mounts`
- AppRole role `czertainly-test` bound to the policy

## Running Tests

**Prerequisites:** Docker installed and running, Go 1.26+.

```bash
# All integration tests
go test -v -tags=integration ./tests/integration/...

# With race detector (used in CI)
go test -v -race -tags=integration ./tests/integration/...

# Single test
go test -v -tags=integration ./tests/integration/... -run TestSecretProvider_SecretCRUD

# Skip integration tests (unit-test-only run)
go test ./...
```

## Test Coverage

| Test | Endpoint | What it validates |
|------|----------|-------------------|
| `TestSecretProvider_ListVaultAttributes` | `GET /v1/secretProvider/vaults/attributes` | Returns 200 with all expected attribute UUIDs (static — no Vault needed) |
| `TestSecretProvider_CheckVaultConnection` | `POST /v1/secretProvider/vaults` | Valid AppRole → 204; invalid secret-id → 400 |
| `TestSecretProvider_SecretCRUD` | `POST`, `PUT`, `DELETE /v1/secretProvider/secrets*` | Full lifecycle (create → read → update → read → delete → 404) for all 8 secret types: `generic`, `apiKey`, `basicAuth`, `jwtToken`, `secretKey`, `privateKey`, `keyValue`, `keyStore` |
| `TestSecretProvider_CrossTypeRead` | `POST /v1/secretProvider/secrets/content` | Reading a secret with a different type than it was created with — compatible pairs return 200, incompatible pairs return 500 |
| `TestSecretProvider_DuplicateCreate` | `POST /v1/secretProvider/secrets` | Creating the same secret name twice → 412 Precondition Failed |
| `TestSecretProvider_NotFoundOperations` | `POST`, `PUT`, `DELETE /v1/secretProvider/secrets*` | Read/Update non-existent → 404; Delete non-existent → 204 (idempotent) |
| `TestSecretProvider_UpdateChangesType` | `PUT /v1/secretProvider/secrets` | Updating to a different type replaces the stored structure; reading with old type fails |
| `TestSecretProvider_SecretPathPrefix` | `POST`, `DELETE /v1/secretProvider/secrets*` | Secrets stored under a sub-path are isolated from the root path |

## File Structure

```
tests/integration/
├── README.md                    — this file
├── constants.go                 — timeouts, attribute UUIDs, credential type constants
├── setup.go                     — shared HTTP helpers (newHTTPClient, doRequest, bytesReader)
├── vault.go                     — VaultContainer: starts Vault testcontainer, configures KV+AppRole
├── builders.go                  — RequestBuilder + secret content helpers for all 8 types
├── assertions.go                — assertStatus, assertBodyContains, assertNotFound
├── fixtures.go                  — SecretTestHarness: Vault container + httptest.Server
├── secrets_test.go              — CRUD lifecycle + connection check + attribute listing
└── secrets_extended_test.go     — cross-type reads, duplicate create, not-found ops, path prefix
```

All files carry `//go:build integration` — excluded from `go test ./...`, only included with `-tags=integration`.

## Adding New Tests

1. Add a new `TestSecretProvider_<Feature>(t *testing.T)` function to `secrets_test.go` (or a new `*_test.go` file).
2. Start with `if testing.Short() { t.Skip(...) }`.
3. Create a harness: `h := NewSecretTestHarness(t); defer h.Cleanup()`.
4. Use `h.Builder()` to get a pre-configured `RequestBuilder`.
5. Use `doRequest()`, `assertStatus()`, `assertBodyContains()` for HTTP interactions.

### Example

```go
//go:build integration

package integration

import (
    "net/http"
    "testing"
)

func TestSecretProvider_MyFeature(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    h := NewSecretTestHarness(t)
    defer h.Cleanup()

    resp := doRequest(t, h.Client, http.MethodPost,
        h.URL()+"/v1/secretProvider/secrets",
        h.Builder().BuildCreateSecretBody("my-secret", GenericSecret("my-value")))
    defer resp.Body.Close()

    assertStatus(t, resp, http.StatusCreated)
}
```

## CI/CD

Integration tests run automatically on every push to `main` and every PR via `.github/workflows/integration.yml`.

The `hashicorp/vault:latest` image is pulled explicitly before running tests to reduce startup latency.

## Troubleshooting

**Vault container fails to start**
- Check Docker is running: `docker ps`
- Pull the image manually: `docker pull hashicorp/vault:latest`

**Tests hang / timeout**
- Run a single test with `-v` to see where it stalls

**`permission denied` from Vault on CRUD**
- The policy must include `sys/internal/ui/mounts` — this is already covered, but verify if you change the Vault setup in `vault.go`

**Build errors without the integration tag**
- All files in `tests/integration/` require `-tags=integration`. Run `go test ./...` (no tag) for unit tests only.
