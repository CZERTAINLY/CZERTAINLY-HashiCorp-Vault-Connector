package secret

import (
	"fmt"
	"net/http"

	sm "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/model"
)

func (s *Server) checkVaultConnection(w http.ResponseWriter, r *http.Request) {
	ckBody, ok := readRBody(w, r)
	if !ok {
		return
	}

	req := []sm.RequestAttribute{}
	if ok := unmrshl(w, ckBody, &req); !ok {
		return
	}

	n := obtainNeeds(r.Context(), w, r, s.k8sToken, &req, nil, ckBody)
	if n == nil {
		return
	}

	if err := n.CommonCheck(); err != nil {
		badrequest(w, fmt.Sprintf("Missing request attribute or validation failed: %s.", err), sm.VALIDATIONFAILED)
		return
	}

	c := obtainVClient(r.Context(), w, r, *n, ckBody)
	if c == nil {
		return
	}

	err := s.m.ConnCheck(r.Context(), c, n.mount)
	if handleOpError(w, r, 0, err, "", "") {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
