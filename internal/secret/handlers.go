package secret

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
	internalVault "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/vault"
)

func (s *Server) createSecret(w http.ResponseWriter, r *http.Request) {
	b, ok := readRBody(w, r)
	if !ok {
		return
	}

	var req sm.CreateSecretRequestDto
	if ok := unmrshl(w, b, &req); !ok {
		return
	}

	n := obtainNeeds(r.Context(), w, r, s.k8sToken, req.VaultAttributes, req.SecretAttributes, b)
	if n == nil {
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c := obtainVClient(r.Context(), w, r, *n, b)
	if c == nil {
		return
	}

	err := s.m.Create(r.Context(), c, n.mount, vaultPath(n.path, req.Name), req.Secret)
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
	b, ok := readRBody(w, r)
	if !ok {
		return
	}

	var req sm.UpdateSecretRequestDto
	if ok := unmrshl(w, b, &req); !ok {
		return
	}

	n := obtainNeeds(r.Context(), w, r, s.k8sToken, req.VaultAttributes, req.SecretAttributes, b)
	if n == nil {
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c := obtainVClient(r.Context(), w, r, *n, b)
	if c == nil {
		return
	}

	err := s.m.Update(r.Context(), c, n.mount, vaultPath(n.path, req.Name), req.Secret)
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
	b, ok := readRBody(w, r)
	if !ok {
		return
	}

	var req sm.SecretRequestDto
	if ok := unmrshl(w, b, &req); !ok {
		return
	}

	n := obtainNeeds(r.Context(), w, r, s.k8sToken, req.VaultAttributes, req.SecretAttributes, b)
	if n == nil {
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c := obtainVClient(r.Context(), w, r, *n, b)
	if c == nil {
		return
	}

	sc, err := s.m.Read(r.Context(), c, n.mount, vaultPath(n.path, req.Name), req.Type)
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

	toJson(r.Context(), w, sc)
}

func (s *Server) deleteSecret(w http.ResponseWriter, r *http.Request) {
	b, ok := readRBody(w, r)
	if !ok {
		return
	}

	var req sm.SecretRequestDto
	if ok := unmrshl(w, b, &req); !ok {
		return
	}

	n := obtainNeeds(r.Context(), w, r, s.k8sToken, req.VaultAttributes, req.SecretAttributes, b)
	if n == nil {
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c := obtainVClient(r.Context(), w, r, *n, b)
	if c == nil {
		return
	}

	err := s.m.Delete(r.Context(), c, n.mount, vaultPath(n.path, req.Name))
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
