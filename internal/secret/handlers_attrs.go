package secret

import (
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
)

func getSecretAttributes(w http.ResponseWriter, r *http.Request) {
	resp := []sm.BaseAttributeDtoV3{}

	toJson(r.Context(), w, http.StatusOK, resp)
}

func (s *Server) listVaultAttributes(w http.ResponseWriter, r *http.Request) {
	var resp []sm.BaseAttributeDtoV3

	var vaultURI sm.BaseAttributeDtoV3
	if err := vaultURI.FromDataAttributeV3(vaultManagementURI); err != nil {
		slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
		internal(w, "Marshaling data structure failed.")
		return
	}
	resp = append(resp, vaultURI)

	var vaultRequestTimeout sm.BaseAttributeDtoV3
	if err := vaultRequestTimeout.FromDataAttributeV3(vaultManagementRequestTmout); err != nil {
		slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
		internal(w, "Marshaling data structure failed.")
		return
	}
	resp = append(resp, vaultRequestTimeout)

	var vaultMount sm.BaseAttributeDtoV3
	if err := vaultMount.FromDataAttributeV3(vaultManagementMount); err != nil {
		slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
		internal(w, "Marshaling data structure failed.")
		return
	}
	resp = append(resp, vaultMount)

	// auth methods

	credentialTypeContent := []sm.BaseAttributeContentDtoV3{}

	var appRole sm.BaseAttributeContentDtoV3
	if err := appRole.FromStringAttributeContentV3(credentialTypeAppRole); err != nil {
		slog.Error("Error marshaling StringAttributeContentV3 into BaseAttributeContentDtoV3", slog.String("error", err.Error()))
		internal(w, "Marshaling data structure failed.")
		return
	}
	credentialTypeContent = append(credentialTypeContent, appRole)

	var jwtToken sm.BaseAttributeContentDtoV3
	if err := jwtToken.FromStringAttributeContentV3(credentialTypeJwt); err != nil {
		slog.Error("Error marshaling StringAttributeContentV3 into BaseAttributeContentDtoV3", slog.String("error", err.Error()))
		internal(w, "Marshaling data structure failed.")
		return
	}
	credentialTypeContent = append(credentialTypeContent, jwtToken)

	if s.k8sToken != nil {
		var kubernetes sm.BaseAttributeContentDtoV3
		if err := kubernetes.FromStringAttributeContentV3(credentialTypeK8s); err != nil {
			slog.Error("Error marshaling StringAttributeContentV3 into BaseAttributeContentDtoV3", slog.String("error", err.Error()))
			internal(w, "Marshaling data structure failed.")
			return
		}
		credentialTypeContent = append(credentialTypeContent, kubernetes)
	}

	var credentialType sm.BaseAttributeDtoV3
	ctcpy := vaultManagementCredentialType
	ctcpy.Content = &credentialTypeContent
	if err := credentialType.FromDataAttributeV3(ctcpy); err != nil {
		slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
		internal(w, "Marshaling data structure failed.")
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
			CallbackMethod:  ptrStr("GET"),
			Mappings: []sm.AttributeCallbackMapping{
				{
					From:                 ptrStr(fmt.Sprintf("%s.data", vaultManagementCredentialType.Name)),
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
		slog.Error("Error marshaling GroupAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
		internal(w, "Marshaling data structure failed.")
		return
	}
	resp = append(resp, credentialGroup)

	var vaultPath sm.BaseAttributeDtoV3
	if err := vaultPath.FromDataAttributeV3(vaultManagementPath); err != nil {
		slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
		internal(w, "Marshaling data structure failed.")
		return
	}
	resp = append(resp, vaultPath)

	toJson(r.Context(), w, http.StatusOK, resp)
}

func (s *Server) credentialsType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	credType := vars["credentialsType"]

	slog.Debug(credType)

	var resp []sm.BaseAttributeDtoV3

	switch credType {
	case credentialTypeAppRole.Data:
		var roleID sm.BaseAttributeDtoV3
		if err := roleID.FromDataAttributeV3(vaultManagementRoleID); err != nil {
			slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
			internal(w, "Marshaling data structure failed.")
			return
		}
		resp = append(resp, roleID)

		var roleSecret sm.BaseAttributeDtoV3
		if err := roleSecret.FromDataAttributeV3(vaultManagementRoleSecret); err != nil {
			slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
			internal(w, "Marshaling data structure failed.")
			return
		}
		resp = append(resp, roleSecret)

	case credentialTypeJwt.Data:
		var role sm.BaseAttributeDtoV3
		if err := role.FromDataAttributeV3(vaultManagementRole); err != nil {
			slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
			internal(w, "Marshaling data structure failed.")
			return
		}
		resp = append(resp, role)

		var jwt sm.BaseAttributeDtoV3
		if err := jwt.FromDataAttributeV3(vaultManagementJwt); err != nil {
			slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
			internal(w, "Marshaling data structure failed.")
			return
		}
		resp = append(resp, jwt)

	case credentialTypeK8s.Data:
		if s.k8sToken == nil {
			badrequest(w, fmt.Sprintf("Credential type unknown: %s.", credType), sm.VALIDATIONFAILED)
			return
		}
		var role sm.BaseAttributeDtoV3
		if err := role.FromDataAttributeV3(vaultManagementRole); err != nil {
			slog.Error("Error marshaling DataAttributeV3 into BaseAttributeDtoV3", slog.String("error", err.Error()))
			internal(w, "Marshaling data structure failed.")
			return
		}
		resp = append(resp, role)

	default:
		badrequest(w, fmt.Sprintf("Credential type unknown: %s.", credType), sm.VALIDATIONFAILED)
		return
	}

	toJson(r.Context(), w, http.StatusOK, resp)
}
