package secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"CZERTAINLY-HashiCorp-Vault-Connector/internal/logger"
	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
	internalVault "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/vault"

	"go.uber.org/zap"
)

func (s *Server) checkVaultConnection(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Calling `io.ReadAll()` failed.", zap.Error(err))
		internal(w, "Reading request body failed.")
		return
	}

	req := []sm.RequestAttribute{}
	if err := json.Unmarshal(b, &req); err != nil {
		log.Debug("Calling `json.Unmarshal()` failed.", zap.Error(err))
		badrequest(w, "Failed to unmarshal request.", sm.ATTRIBUTESERROR)
		return
	}

	n := NewNeeds(s.k8sToken)
	if err = n.Process(ctx, &req, nil); err != nil {
		log.Debug("Processing request attributes failed.",
			zap.Error(err),
			zap.String("http-path", r.URL.Path),
			zap.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Processing request attributes failed: %s.", err), sm.ATTRIBUTESERROR)
		return
	}

	if err = n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c := obtainVClient(ctx, w, r, n, b)
	if c == nil {
		return
	}

	err = s.m.ConnCheck(ctx, c)
	switch {
	case errors.Is(err, internalVault.ErrForbidden):
		forbidden(w, fmt.Sprintf("Authorization failed: %s.", err))
		return

	case err != nil:
		log.Error("Could not connect to Vault.",
			zap.Error(err),
			zap.String("http-path", r.URL.Path),
			zap.String("request-body", string(b)))
		internal(w, fmt.Sprintf("Could not connect to Vault: %s", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
