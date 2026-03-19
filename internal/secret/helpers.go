package secret

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"CZERTAINLY-HashiCorp-Vault-Connector/internal/logger"
	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
	internalVault "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/vault"

	vcg "github.com/hashicorp/vault-client-go"
	"go.uber.org/zap"
)

func ptr[T any](v T) *T {
	return &v
}

func vaultPath(pathPrefix, secretPath, name string) string {
	var res string
	if pathPrefix != "" {
		res = fmt.Sprintf("%s/", pathPrefix)
	}
	if secretPath != "" {
		res = fmt.Sprintf("%s%s/", res, secretPath)
	}

	return fmt.Sprintf("%s%s", res, name)
}

func toJson(_ context.Context, w http.ResponseWriter, statusCode int, resp any) {
	b, err := json.Marshal(resp)
	if err != nil {
		logger.Get().Error("Failed to marshal structure to json.",
			zap.Error(err),
			zap.Any("structure", resp))
		internal(w, "Failed to marshal structure to json.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
}

func obtainVClient(ctx context.Context, w http.ResponseWriter, r *http.Request, n Needs, body []byte) *vcg.Client {
	c, err := n.Client(ctx)
	switch {
	case vcg.IsErrorStatus(err, http.StatusUnauthorized):
		unauthorized(w, fmt.Sprintf("Authentication failed: %s.", err))
		return nil

	case vcg.IsErrorStatus(err, http.StatusForbidden):
		forbidden(w, fmt.Sprintf("Authorization failed: %s.", err))
		return nil

	case err != nil:
		logger.Get().Debug("Could not connect to Vault.",
			zap.Error(err),
			zap.String("http-path", r.URL.Path),
			zap.String("request-body", string(body)))
		badrequest(w, fmt.Sprintf("Could not connect to Vault: %s", err), sm.ATTRIBUTESERROR)
		return nil
	}
	return c
}

func readRBody(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Get().Error("Calling `io.ReadAll()` failed.", zap.Error(err))
		internal(w, "Reading request body failed.")
		return b, false
	}
	return b, true
}

func unmrshl(w http.ResponseWriter, body []byte, req any) bool {
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Get().Debug("Calling `json.Unmarshal()` failed.", zap.Error(err))
		badrequest(w, "Failed to unmarshal request.", sm.ATTRIBUTESERROR)
		return false
	}
	return true
}

func obtainNeeds(_ context.Context, w http.ResponseWriter, r *http.Request, k8sToken *string, vaultAttrs, vaultProfileAttrs, secretAttrs *[]sm.RequestAttribute, b []byte) *Needs {
	n := NewNeeds(k8sToken)
	if err := n.Process(r.Context(), vaultAttrs, vaultProfileAttrs, secretAttrs); err != nil {
		logger.Get().Debug("Processing request attributes failed.",
			zap.Error(err),
			zap.String("http-path", r.URL.Path),
			zap.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Processing request attributes failed: %s.", err), sm.ATTRIBUTESERROR)
		return nil
	}
	return &n
}

func handleOpError(w http.ResponseWriter, r *http.Request, statusCode int, err error, reqName, reqType string) (doReturn bool) {
	log := logger.Get()
	switch {
	case errors.Is(err, internalVault.ErrForbidden):
		log.Debug("Authorization failed.", zap.Error(err))
		forbidden(w, "Authorization failed.")
		return true

	case errors.Is(err, internalVault.ErrNotFound):
		log.Debug("Not found.", zap.Error(err))
		notfound(w, err.Error())
		return true

	case errors.Is(err, internalVault.ErrAlreadyExists):
		precondition(w, "Secret already exists.", sm.RESOURCEALREADYEXISTS)
		return true

	case err != nil:
		log.Error("Operation failed.", zap.Error(err), zap.String("http-path", r.URL.Path))
		internal(w, "Operation failed.")
		return true
	}

	if reqName != "" && reqType != "" {
		toJson(r.Context(), w, statusCode, sm.SecretResponseDto{
			Name: reqName,
			Type: sm.SecretType(reqType),
		})
	} else if statusCode > 0 {
		w.WriteHeader(statusCode)
	}

	return false
}
