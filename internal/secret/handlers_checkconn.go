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

func (s *Server) checkVaultConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		internal(w, "Reading request body failed.")
		return
	}

	req := []sm.RequestAttribute{}
	if err := json.Unmarshal(b, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		badrequest(w, "Failed to unmarshal request.", sm.ATTRIBUTESERROR)
		return
	}

	n := NewNeeds(s.k8sToken)
	if err = n.Process(ctx, &req, nil); err != nil {
		slog.Debug("Processing request attributes failed.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Processing request attributes failed: %s.", err), sm.ATTRIBUTESERROR)
		return
	}

	if err = n.ConnectionCheck(); err != nil {
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

	err = s.m.ConnCheck(ctx, c)
	switch {
	case errors.Is(err, internalVault.ErrForbidden):
		forbidden(w, fmt.Sprintf("Authorization failed: %s.", err))
		return

	case err != nil:
		slog.Error("Could not connect to Vault.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		internal(w, fmt.Sprintf("Could not connect to Vault: %s", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
