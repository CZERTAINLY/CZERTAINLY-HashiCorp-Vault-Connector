package secret

import (
	"fmt"
	"net/http"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
)

func (s *Server) createSecret(w http.ResponseWriter, r *http.Request) {
	crtBody, ok := readRBody(w, r)
	if !ok {
		return
	}

	var req sm.CreateSecretRequestDto
	if ok := unmrshl(w, crtBody, &req); !ok {
		return
	}

	n := obtainNeeds(r.Context(), w, r, s.k8sToken, req.VaultAttributes, req.VaultProfileAttributes, req.SecretAttributes, crtBody)
	if n == nil {
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c := obtainVClient(r.Context(), w, r, *n, crtBody)
	if c == nil {
		return
	}

	scrtType, err := s.m.Create(r.Context(), c, n.mount, vaultPath(n.pathPrefix, n.secretPath, req.Name), req.Secret)
	_ = handleOpError(w, r, http.StatusCreated, err, req.Name, string(scrtType))
}

func (s *Server) updateSecret(w http.ResponseWriter, r *http.Request) {
	uptdBody, ok := readRBody(w, r)
	if !ok {
		return
	}

	var req sm.UpdateSecretRequestDto
	if ok := unmrshl(w, uptdBody, &req); !ok {
		return
	}

	n := obtainNeeds(r.Context(), w, r, s.k8sToken, req.VaultAttributes, req.VaultProfileAttributes, req.SecretAttributes, uptdBody)
	if n == nil {
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c := obtainVClient(r.Context(), w, r, *n, uptdBody)
	if c == nil {
		return
	}

	secretType, err := s.m.Update(r.Context(), c, n.mount, vaultPath(n.pathPrefix, n.secretPath, req.Name), req.Secret)
	_ = handleOpError(w, r, http.StatusOK, err, req.Name, string(secretType))
}

func (s *Server) getSecretValue(w http.ResponseWriter, r *http.Request) {
	getBody, ok := readRBody(w, r)
	if !ok {
		return
	}

	var req sm.SecretRequestDto
	if ok := unmrshl(w, getBody, &req); !ok {
		return
	}

	n := obtainNeeds(r.Context(), w, r, s.k8sToken, req.VaultAttributes, req.VaultProfileAttributes, req.SecretAttributes, getBody)
	if n == nil {
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c := obtainVClient(r.Context(), w, r, *n, getBody)
	if c == nil {
		return
	}

	sc, err := s.m.Read(r.Context(), c, n.mount, vaultPath(n.pathPrefix, n.secretPath, req.Name), req.Type)
	if handleOpError(w, r, 0, err, "", "") {
		return
	}

	toJson(r.Context(), w, http.StatusOK, sm.SecretContentResponseDto{
		Content: sc,
	})
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

	n := obtainNeeds(r.Context(), w, r, s.k8sToken, req.VaultAttributes, req.VaultProfileAttributes, req.SecretAttributes, b)
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

	err := s.m.Delete(r.Context(), c, n.mount, vaultPath(n.pathPrefix, n.secretPath, req.Name))
	_ = handleOpError(w, r, http.StatusNoContent, err, "", "")
}

func (s *Server) rotateSecretValue(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not implemented yet.", http.StatusNotImplemented)
}
func (s *Server) getRotateAttributes(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not implemented yet.", http.StatusNotImplemented)
}
