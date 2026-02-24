package secret

import (
	"net/http"

	sv "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/vault"

	"github.com/gorilla/mux"
)

type Server struct {
	m *sv.Manager
}

func New() Server {
	return Server{
		m: sv.New(),
	}
}

func (s *Server) MuxRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// TODO: Finish
	//r.Use(LimitBodySize())

	r.Methods(http.MethodPost).Path("/v1/secretProvider/vaults").Handler(http.HandlerFunc(checkVaultConnection))

	r.Methods(http.MethodGet).Path("/v1/secretProvider/secrets").Handler(http.HandlerFunc(listSecrets))
	r.Methods(http.MethodPost).Path("/v1/secretProvider/secrets").Handler(http.HandlerFunc(s.createSecret))
	r.Methods(http.MethodPut).Path("/v1/secretProvider/secrets").Handler(http.HandlerFunc(s.updateSecret))
	r.Methods(http.MethodPost).Path("/v1/secretProvider/secrets/value").Handler(http.HandlerFunc(s.getSecretValue))
	r.Methods(http.MethodPost).Path("/v1/secretProvider/secrets/{uuid}/delete").Handler(http.HandlerFunc(s.deleteSecret))

	r.Methods(http.MethodGet).Path("/v1/secretProvider/vaults/attributes").Handler(http.HandlerFunc(listVaultAttributes))
	r.Methods(http.MethodGet).Path("/v1/secretProvider/secrets/{secretType}/attributes").Handler(http.HandlerFunc(getSecretAttributes))

	// TODO: callback for a group attribute

	return r
}

// TODO: finish idea: limiting maximum request body size
/*
func LimitBodySize(next http.Handler, maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Replace the request body with a MaxBytesReader
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
*/
