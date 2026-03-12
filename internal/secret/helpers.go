package secret

import (
	"context"
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

func ptrStr(v string) *string {
	return &v
}

func ptrAttributeResource(r sm.AttributeResource) *sm.AttributeResource {
	return &r
}

func vaultPath(path, name string) string {
	return fmt.Sprintf("%s/%s", path, name)
}

func toJson(_ context.Context, w http.ResponseWriter, resp any) {
	b, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal structure to json.",
			slog.String("error", err.Error()),
			slog.Any("structure", resp))
		internal(w, "Failed to marshal structure to json.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func obtainVClient(ctx context.Context, w http.ResponseWriter, r *http.Request, n Needs, body []byte) *vcg.Client {
	c, err := n.Client(ctx)
	switch {
	case vcg.IsErrorStatus(err, http.StatusUnauthorized):
		unauthorized(w, fmt.Sprintf("Authentication failed: %s.", err))
		return nil
	case err != nil:
		slog.Debug("Could not connect to Vault.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(body)))
		badrequest(w, fmt.Sprintf("Could not connect to Vault: %s", err), sm.ATTRIBUTESERROR)
		return nil
	}
	return c
}

func readRBody(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Calling `io.ReadAll()` failed.", slog.String("error", err.Error()))
		internal(w, "Reading request body failed.")
		return b, false
	}
	return b, true
}

func unmrshl(w http.ResponseWriter, body []byte, req any) bool {
	if err := json.Unmarshal(body, &req); err != nil {
		slog.Debug("Calling `json.Unmarshal()` failed.", slog.String("error", err.Error()))
		badrequest(w, "Failed to unmarshal request.", sm.ATTRIBUTESERROR)
		return false
	}
	return true
}

func obtainNeeds(ctx context.Context, w http.ResponseWriter, r *http.Request, k8sToken *string, vaultAttrs, secretAttrs *[]sm.RequestAttribute, b []byte) *Needs {
	n := NewNeeds(k8sToken)
	if err := n.Process(r.Context(), vaultAttrs, secretAttrs); err != nil {
		slog.Debug("Processing request attributes failed.",
			slog.String("error", err.Error()),
			slog.String("http-path", r.URL.Path),
			slog.String("request-body", string(b)))
		badrequest(w, fmt.Sprintf("Processing request attributes failed: %s.", err), sm.ATTRIBUTESERROR)
		return nil
	}
	return &n
}

func handleOpError(w http.ResponseWriter, r *http.Request, err error) (doReturn bool) {
	switch {
	case errors.Is(err, internalVault.ErrForbidden):
		slog.Debug("Authorization failed.", slog.String("error", err.Error()))
		forbidden(w, "Authorization failed.")
		return true

	case errors.Is(err, internalVault.ErrNotFound):
		slog.Debug("Not found.", slog.String("error", err.Error()))
		notfound(w, "Not found.")
		return true

	case errors.Is(err, internalVault.ErrAlreadyExists):
		precondition(w, "Secret already exists.", sm.RESOURCEALREADYEXISTS)

	case err != nil:
		slog.Error("Operation failed.", slog.String("error", err.Error()), slog.String("http-path", r.URL.Path))
		internal(w, "Operation failed.")
		return true
	}

	return false
}
