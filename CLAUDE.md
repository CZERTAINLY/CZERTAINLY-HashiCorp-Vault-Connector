# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
# Build binary
go build -v ./...
go build -ldflags "-s -w -X main.version=<version>" ./cmd

# Run tests (with race detector)
go test -race -vet=off -coverprofile=coverage.out ./...

# Run a single test
go test -run TestName ./internal/model/

# Lint
golangci-lint run --timeout=5m

# Vet
go vet ./...

# Verify dependencies
go mod verify
```

## Code Generation

API server code is generated from OpenAPI specs. To regenerate:
```bash
./generate.sh
```
Do not manually edit files in `internal/model/` that are marked as generated.

## Architecture

This is a **CZERTAINLY framework connector** for HashiCorp Vault PKI. It exposes REST APIs that the CZERTAINLY platform calls to manage and discover certificates via Vault's PKI secrets engine.

### Providers

Two function groups are implemented (both labeled `HVault`):

- **Authority Provider** (`/v1/authorityProvider/*`, `/v2/authorityProvider/*`) — issue, renew, revoke, and identify certificates; download CA certs and CRLs
- **Discovery Provider** (`/v1/discoveryProvider/*`) — discover certificates stored in Vault instances

A **Secret Provider** (`/v1/secretProvider/*`) is also exposed as a v2 interface.

### Request flow

`cmd/main.go` → wires together all services and controllers → Gorilla Mux router with correlation-ID middleware → service layer → `internal/vault/` client → HashiCorp Vault API

### Key packages

| Package | Purpose |
|---|---|
| `internal/authority/` | Authority provider: certificate issuance, renewal, revocation |
| `internal/discovery/` | Discovery provider: scans Vault for certificates |
| `internal/vault/` | Vault HTTP client + auth methods (AppRole, Kubernetes, JWT/OIDC) |
| `internal/secret/` | Secret provider implementation |
| `internal/db/` | PostgreSQL via GORM; auto-migration on startup |
| `internal/model/` | Generated OpenAPI DTOs and shared types |
| `internal/config/` | Environment-variable-based configuration |
| `internal/logger/` | uber/zap structured logging with context propagation (zax) |

### Database

PostgreSQL (≥12) with schema `hvault` (configurable). Tables: `authority_instances`, `certificates`, `discoveries`, `discovery_certificates`. Migrations live in `migrations/` and run automatically at startup.

### Vault authentication

AppRole (RoleID + SecretID), Kubernetes service account token, and JWT/OIDC are all supported — selected per authority instance via connector attributes.

## Required Environment Variables

| Variable | Default | Notes |
|---|---|---|
| `DATABASE_NAME` | — | **Required** |
| `DATABASE_USER` | — | **Required** |
| `DATABASE_PASSWORD` | — | **Required** |
| `DATABASE_HOST` | `localhost` | |
| `DATABASE_PORT` | `5432` | |
| `DATABASE_SCHEMA` | `hvault` | |
| `DATABASE_SSL_MODE` | `require` | |
| `SERVER_PORT` | `8080` | |
| `LOG_LEVEL` | `INFO` | |
