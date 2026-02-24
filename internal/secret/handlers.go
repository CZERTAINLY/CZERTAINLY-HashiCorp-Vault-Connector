package secret

import (
	// "context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	// "github.com/gorilla/mux"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
	// sv "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/vault"
)

func checkVaultConnection(w http.ResponseWriter, r *http.Request) {
	// POST
	// request body ContentType "application/json"
	// array of RequestAttributeDto
}

func listSecrets(w http.ResponseWriter, r *http.Request) {
	// GET
	// Response
	//  ContentType "application/json"
	//  json array of SecretResponseDto
}

func (s *Server) createSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
	// TODO: rework after limiting request size
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	var req sm.CreateSecretRequestDto
	if err := json.Unmarshal(b, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to unmarshal request: %s.", err), http.StatusBadRequest)
		return
	}

	attrs, err := ProcessAttrs(ctx, *req.VaultAttributes, *req.SecretAttributes)
	if err != nil {
		slog.Debug("Processing request attributes failed.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Processing request attributes failed: %s", err), http.StatusBadRequest)
		return
	}

	c, err := attrs.Client()
	if err != nil {
		slog.Debug("Could not connect to Vault.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Could not connect to Vault: %s", err), http.StatusBadRequest)
		return
	}

	if err := s.m.Create(ctx, c, attrs.mount, vaultPath(attrs.path, req.Name), req.Secret); err != nil {
		slog.Error("Error creating secret.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Error creating secret: %s", err), http.StatusBadRequest)
		return
	}
}

func (s *Server) updateSecret(w http.ResponseWriter, r *http.Request) {
	slog.Debug("!!! KEBABI !!! Inside `updateSecret`")
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
	// TODO: rework after limiting request size
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	var req sm.UpdateSecretRequestDto
	if err := json.Unmarshal(b, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to unmarshal request: %s.", err), http.StatusBadRequest)
		return
	}

	attrs, err := ProcessAttrs(ctx, *req.VaultAttributes, *req.SecretAttributes)
	if err != nil {
		slog.Debug("Processing request attributes failed.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Processing request attributes failed: %s", err), http.StatusBadRequest)
		return
	}

	c, err := attrs.Client()
	if err != nil {
		slog.Debug("Could not connect to Vault.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Could not connect to Vault: %s", err), http.StatusBadRequest)
		return
	}

	if err := s.m.Update(ctx, c, attrs.mount, vaultPath(attrs.path, req.Name), req.Secret); err != nil {
		slog.Error("Error creating secret.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Error creating secret: %s", err), http.StatusBadRequest)
		return
	}
}

func (s *Server) getSecretValue(w http.ResponseWriter, r *http.Request) {
	slog.Debug("!!! KEBABI !!! Inside `getSecretValue`")
	ctx := r.Context()

	var err error
	var b []byte

	b, err = io.ReadAll(r.Body)
	// TODO: rework after limiting request size
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	var req sm.SecretRequestDto
	if err := json.Unmarshal(b, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to unmarshal request: %s.", err), http.StatusBadRequest)
		return
	}

	attrs, err := ProcessAttrs(ctx, *req.VaultAttributes, *req.SecretAttributes)
	if err != nil {
		slog.Debug("Processing request attributes failed.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Processing request attributes failed: %s", err), http.StatusBadRequest)
		return
	}

	c, err := attrs.Client()
	if err != nil {
		slog.Debug("Could not connect to Vault.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Could not connect to Vault: %s", err), http.StatusBadRequest)
		return
	}

	sc, err := s.m.Read(ctx, c, attrs.mount, vaultPath(attrs.path, req.Name), req.Type)
	if err != nil {
		slog.Error("Error reading secret.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Error reading secret: %s", err), http.StatusBadRequest)
		return
	}

	b, err = json.Marshal(sc)
	if err != nil {
		slog.Error("Failed to marshal structure to json.", slog.String("error", err.Error()))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal server error."))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (s *Server) deleteSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var b []byte

	b, err = io.ReadAll(r.Body)
	// TODO: rework after limiting request size
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	var req sm.SecretRequestDto
	if err := json.Unmarshal(b, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to unmarshal request: %s.", err), http.StatusBadRequest)
		return
	}

	attrs, err := ProcessAttrs(ctx, *req.VaultAttributes, *req.SecretAttributes)
	if err != nil {
		slog.Debug("Processing request attributes failed.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Processing request attributes failed: %s", err), http.StatusBadRequest)
		return
	}

	c, err := attrs.Client()
	if err != nil {
		slog.Debug("Could not connect to Vault.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Could not connect to Vault: %s", err), http.StatusBadRequest)
		return
	}

	if err := s.m.Delete(ctx, c, attrs.mount, vaultPath(attrs.path, req.Name)); err != nil {
		slog.Error("Error deleting secret.", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Error deleting secret: %s", err), http.StatusBadRequest)
		return
	}
}

func listVaultAttributes(w http.ResponseWriter, r *http.Request) {
}

func getSecretAttributes(w http.ResponseWriter, r *http.Request) {
}

func vaultPath(path, name string) string {
	return fmt.Sprintf("%s/%s", path, name)
}
