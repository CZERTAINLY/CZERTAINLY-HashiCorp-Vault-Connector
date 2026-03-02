package secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
	internalVault "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/vault"

	vcg "github.com/hashicorp/vault-client-go"
)

func (s *Server) createSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		internal(w, "Reading request body failed.")
		return
	}

	var req sm.CreateSecretRequestDto
	if err := json.Unmarshal(b, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		badrequest(w, "Failed to unmarshal request.", sm.ATTRIBUTESERROR)
		return
	}

	n := NewNeeds(s.k8sToken)
	if err := n.Process(ctx, req.VaultAttributes, req.SecretAttributes); err != nil {
		slog.Debug("Processing request attributes failed.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Processing request attributes failed: %s.", err), sm.ATTRIBUTESERROR)
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c, err := n.Client(ctx)
	switch {
	case vcg.IsErrorStatus(err, http.StatusUnauthorized):
		unauthorized(w, fmt.Sprintf("Authentication failed: %s.", err))
		return
	case err != nil:
		slog.Debug("Could not connect to Vault.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Could not connect to Vault: %s", err), sm.ATTRIBUTESERROR)
		return
	}

	err = s.m.Create(ctx, c, n.mount, vaultPath(n.path, req.Name), req.Secret)
	switch {
	case errors.Is(err, internalVault.ErrForbidden):
		forbidden(w, fmt.Sprintf("Authorization failed: %s.", err))
		return
	case errors.Is(err, internalVault.ErrAlreadyExists):
		precondition(w, "Secret already exists.", sm.RESOURCEALREADYEXISTS)
		return

	case err != nil:
		slog.Error("Failed to create secret.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		internal(w, fmt.Sprintf("Failed to create secret: %s", err))
		return
	}
}

func (s *Server) updateSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		internal(w, "Reading request body failed.")
		return
	}

	var req sm.UpdateSecretRequestDto
	if err := json.Unmarshal(b, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		badrequest(w, "Failed to unmarshal request.", sm.ATTRIBUTESERROR)
		return
	}

	n := NewNeeds(s.k8sToken)
	if err := n.Process(ctx, req.VaultAttributes, req.SecretAttributes); err != nil {
		slog.Debug("Processing request attributes failed.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Processing request attributes failed: %s.", err), sm.ATTRIBUTESERROR)
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c, err := n.Client(ctx)
	switch {
	case vcg.IsErrorStatus(err, http.StatusUnauthorized):
		unauthorized(w, fmt.Sprintf("Authentication failed: %s.", err))
		return
	case err != nil:
		slog.Debug("Could not connect to Vault.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Could not connect to Vault: %s", err), sm.ATTRIBUTESERROR)
		return
	}

	err = s.m.Update(ctx, c, n.mount, vaultPath(n.path, req.Name), req.Secret)
	switch {
	case errors.Is(err, internalVault.ErrForbidden):
		forbidden(w, fmt.Sprintf("Authorization failed: %s.", err))
		return
	case errors.Is(err, internalVault.ErrNotFound):
		notfound(w, "Secret not found.")
		return

	case err != nil:
		slog.Error("Failed to update secret.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		internal(w, fmt.Sprintf("Failed to update secret: %s", err))
		return
	}

}

func (s *Server) getSecretValue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var b []byte

	b, err = io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		internal(w, "Reading request body failed.")
		return
	}

	var req sm.SecretRequestDto
	if err := json.Unmarshal(b, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		badrequest(w, "Failed to unmarshal request.", sm.ATTRIBUTESERROR)
		return
	}

	n := NewNeeds(s.k8sToken)
	if err := n.Process(ctx, req.VaultAttributes, req.SecretAttributes); err != nil {
		slog.Debug("Processing request attributes failed.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Processing request attributes failed: %s.", err), sm.ATTRIBUTESERROR)
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c, err := n.Client(ctx)
	switch {
	case vcg.IsErrorStatus(err, http.StatusUnauthorized):
		unauthorized(w, fmt.Sprintf("Authentication failed: %s.", err))
		return
	case err != nil:
		slog.Debug("Could not connect to Vault.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Could not connect to Vault: %s", err), sm.ATTRIBUTESERROR)
		return
	}

	sc, err := s.m.Read(ctx, c, n.mount, vaultPath(n.path, req.Name), req.Type)
	switch {
	case errors.Is(err, internalVault.ErrForbidden):
		forbidden(w, fmt.Sprintf("Authorization failed: %s.", err))
		return

	case errors.Is(err, internalVault.ErrNotFound):
		notfound(w, "Secret not found.")
		return

	case err != nil:
		slog.Error("Failed to read secret.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		internal(w, fmt.Sprintf("Failed to read secret: %s", err))
		return
	}

	toJson(ctx, w, sc)
}

func (s *Server) deleteSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error
	var b []byte

	b, err = io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		internal(w, "Reading request body failed.")
		return
	}

	var req sm.SecretRequestDto
	if err := json.Unmarshal(b, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		badrequest(w, "Failed to unmarshal request.", sm.ATTRIBUTESERROR)
		return
	}

	n := NewNeeds(s.k8sToken)
	if err := n.Process(ctx, req.VaultAttributes, req.SecretAttributes); err != nil {
		slog.Debug("Processing request attributes failed.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Processing request attributes failed: %s.", err), sm.ATTRIBUTESERROR)
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}
	c, err := n.Client(ctx)
	switch {
	case vcg.IsErrorStatus(err, http.StatusUnauthorized):
		unauthorized(w, fmt.Sprintf("Authentication failed: %s.", err))
		return
	case err != nil:
		slog.Debug("Could not connect to Vault.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Could not connect to Vault: %s", err), sm.ATTRIBUTESERROR)
		return
	}

	err = s.m.Delete(ctx, c, n.mount, vaultPath(n.path, req.Name))
	switch {
	case errors.Is(err, internalVault.ErrForbidden):
		forbidden(w, fmt.Sprintf("Authorization failed: %s.", err))
		return

	case errors.Is(err, internalVault.ErrNotFound):
		notfound(w, "Secret not found.")
		return

	case err != nil:
		slog.Error("Failed to delete secret.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		internal(w, fmt.Sprintf("Failed to delete secret: %s", err))
		return
	}

}

func (s *Server) rotateSecretValue(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not implemented yet.", http.StatusNotImplemented)
}
func (s *Server) getRotateAttributes(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not implemented yet.", http.StatusNotImplemented)
}
