package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	// "github.com/gorilla/mux"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
	// sv "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/vault"
)

func (s *Server) checkVaultConnection(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented yet.", http.StatusNotImplemented)
}

func (s *Server) createSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
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
	ctx := r.Context()

	b, err := io.ReadAll(r.Body)
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
	ctx := r.Context()

	var err error
	var b []byte

	b, err = io.ReadAll(r.Body)
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

func (s *Server) listVaultAttributes(w http.ResponseWriter, r *http.Request) {
	var resp []sm.BaseAttributeDtoV3

	var vaultURI sm.BaseAttributeDtoV3
	if err := vaultURI.FromDataAttributeV3(sm.DataAttributeV3{
		Uuid:          VaultManagementUriUUID,
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          VaultManagementUriName,
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("Vault URI should be in the following format: `http(s)://<vault-url>:<port>`."),
		Properties: sm.DataAttributeProperties{
			Label:    "Vault URI",
			Visible:  true,
			Required: true,
		},
	}); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	resp = append(resp, vaultURI)

	var vaultRequestTimeout sm.BaseAttributeDtoV3
	if err := vaultRequestTimeout.FromDataAttributeV3(sm.DataAttributeV3{
		Uuid:          VaultManagementRequestTimeoutUUID,
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          VaultManagementRequestTimeoutName,
		ContentType:   sm.AttributeContentTypeInteger,
		Description:   ptrStr("Request timeout in seconds applied to each Vault request."),
		Properties: sm.DataAttributeProperties{
			Label:    "Individual Vault request timeout",
			Visible:  true,
			Required: true,
		},
	}); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	resp = append(resp, vaultRequestTimeout)

	var vaultMount sm.BaseAttributeDtoV3
	if err := vaultMount.FromDataAttributeV3(sm.DataAttributeV3{
		Uuid:          VaultManagementMountUUID,
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          VaultManagementMountName,
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("Vault mount."),
		Properties: sm.DataAttributeProperties{
			Label:    "Vault mount",
			Visible:  true,
			Required: true,
		},
	}); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	resp = append(resp, vaultMount)

	// auth methods

	credentialTypeContent := []sm.BaseAttributeContentDtoV3{}

	var appRole sm.BaseAttributeContentDtoV3
	if err := appRole.FromStringAttributeContentV3(sm.StringAttributeContentV3{
		Reference: ptrStr("AppRole"),
		Data:      "approle",
	}); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	credentialTypeContent = append(credentialTypeContent, appRole)

	var jwtToken sm.BaseAttributeContentDtoV3
	if err := jwtToken.FromStringAttributeContentV3(sm.StringAttributeContentV3{
		Reference: ptrStr("JWT/OIDC"),
		Data:      "jwt",
	}); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	credentialTypeContent = append(credentialTypeContent, jwtToken)

	if s.k8sTokenExists {
		var kubernetes sm.BaseAttributeContentDtoV3
		if err := jwtToken.FromStringAttributeContentV3(sm.StringAttributeContentV3{
			Reference: ptrStr("Kubernetes"),
			Data:      "kubernetes",
		}); err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		credentialTypeContent = append(credentialTypeContent, kubernetes)
	}

	var credentialType sm.BaseAttributeDtoV3
	if err := credentialType.FromDataAttributeV3(sm.DataAttributeV3{
		Uuid:          VaultManagementCredentialTypeUUID,
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          VaultManagementCredentialTypeName,
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("List of available Vault authentication methods."),
		Properties: sm.DataAttributeProperties{
			Label:       "Please select an authentication method",
			Visible:     true,
			Required:    true,
			ReadOnly:    false,
			List:        true,
			MultiSelect: false,
		},
		Content: &credentialTypeContent,
	}); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	resp = append(resp, credentialType)

	credentialGroupAttrType := sm.Data
	credentialGroupAttrContentType := sm.AttributeContentTypeString
	credentialGroupVersion := int32(3)

	var credentialGroup sm.BaseAttributeDtoV3
	if err := credentialGroup.FromGroupAttributeV3(sm.GroupAttributeV3{
		Uuid:          VaultManagementCredentialGroupUUID,
		Version:       &credentialGroupVersion,
		SchemaVersion: sm.V3,
		Name:          VaultManagementCredentialGroupName,
		AttributeCallback: &sm.AttributeCallback{
			CallbackContext: ptrStr("v1/secretProvider/credentialType/{credentialsType}/callback"),
			CallbackMethod:  "GET",
			Mappings: []sm.AttributeCallbackMapping{
				{
					From:                 ptrStr(fmt.Sprintf("%s.data", VaultManagementCredentialTypeName)),
					AttributeType:        &credentialGroupAttrType,
					AttributeContentType: &credentialGroupAttrContentType,
					To:                   "credentialsType",
					Targets: []sm.AttributeValueTarget{
						sm.PathVariable,
					},
				},
			},
		},
	}); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	resp = append(resp, credentialGroup)

	/* TODO: Remove
	var vaultToken sm.BaseAttributeDtoV3
	if err := vaultToken.FromDataAttributeV3(sm.DataAttributeV3{
		Uuid:          VaultManagementTokenUUID,
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          VaultManagementTokenName,
		Type:          sm.Data,
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("Vault authentication token."),
		Properties: sm.DataAttributeProperties{
			Label:    "Vault token",
			Visible:  true,
			Required: true,
		},
	}); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	resp = append(resp, vaultToken)
	*/

	toJson(r.Context(), w, resp)
}

func getSecretAttributes(w http.ResponseWriter, r *http.Request) {
	var resp []sm.BaseAttributeDtoV3

	var secretPath sm.BaseAttributeDtoV3
	if err := secretPath.FromDataAttributeV3(sm.DataAttributeV3{
		Uuid:          SecretManagementSecretPathUUID,
		Version:       3,
		SchemaVersion: sm.V3,
		Name:          SecretManagementSecretPathName,
		Type:          sm.Data,
		ContentType:   sm.AttributeContentTypeString,
		Description:   ptrStr("Path of secret in Vault without trailing slash."),
		Properties: sm.DataAttributeProperties{
			Label:    "Secret Path",
			Visible:  true,
			Required: true,
		},
	}); err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	resp = append(resp, secretPath)
	toJson(r.Context(), w, resp)
}

func (s *Server) rotateSecretValue(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not implemented yet.", http.StatusNotImplemented)
}
func (s *Server) getRotateAttributes(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not implemented yet.", http.StatusNotImplemented)
}

func vaultPath(path, name string) string {
	return fmt.Sprintf("%s/%s", path, name)
}

func ptrStr(v string) *string {
	return &v
}

func toJson(_ context.Context, w http.ResponseWriter, resp any) {
	b, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal structure to json.",
			slog.String("error", err.Error()),
			slog.Any("structure", resp))
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
