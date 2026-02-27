package secret

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	sv "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/vault"

	"github.com/gorilla/mux"
)

const DEFAULT_K8S_TOKEN_PATH = "/var/run/secrets/kubernetes.io/serviceaccount/token"

type Server struct {
	m *sv.Manager

	k8sTokenExists bool
	k8sToken       string
}

func New() Server {
	s := Server{
		m: sv.New(),
	}

	_, err := os.Stat(DEFAULT_K8S_TOKEN_PATH)
	switch {
	case errors.Is(err, os.ErrNotExist):
		slog.Debug("Kubernetes service account JWT file is not present.",
			slog.String("file", DEFAULT_K8S_TOKEN_PATH))
	case err != nil:
		slog.Error("Error executing `os.Stat()` for kubernetes service account JWT file.",
			slog.String("file", DEFAULT_K8S_TOKEN_PATH),
			slog.String("error", err.Error()))
	default:
		slog.Debug("Kubernetes service account JWT file exists.")
		b, rerr := os.ReadFile(DEFAULT_K8S_TOKEN_PATH)
		switch {
		case rerr != nil:
			slog.Error("Error reading kubernetes service account JWT file.",
				slog.String("file", DEFAULT_K8S_TOKEN_PATH),
				slog.String("error", rerr.Error()))
		default:
			s.k8sTokenExists = true
			s.k8sToken = string(b)
		}
	}

	return s
}

func (s *Server) MuxRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// WIP
	//r.Use(LimitBodySize())

	r.Methods(http.MethodPost).Path("/v1/secretProvider/secrets").Handler(http.HandlerFunc(s.createSecret))
	r.Methods(http.MethodPut).Path("/v1/secretProvider/secrets").Handler(http.HandlerFunc(s.updateSecret))
	r.Methods(http.MethodPost).Path("/v1/secretProvider/secrets/content").Handler(http.HandlerFunc(s.getSecretValue))
	r.Methods(http.MethodDelete).Path("/v1/secretProvider/secrets").Handler(http.HandlerFunc(s.deleteSecret))

	r.Methods(http.MethodGet).Path("/v1/secretProvider/vaults/attributes").Handler(http.HandlerFunc(s.listVaultAttributes))
	r.Methods(http.MethodGet).Path("/v1/secretProvider/secrets/{secretType}/attributes").Handler(http.HandlerFunc(getSecretAttributes))
	// TODO: v1/authorityProvider/credentialType/{credentialsType}/callback

	r.Methods(http.MethodPost).Path("/v1/secretProvider/vaults").Handler(http.HandlerFunc(s.checkVaultConnection))

	// Not implemented for now - returns HTTP 501 Not Implemented
	r.Methods(http.MethodPost).Path("/v1/secretProvider/secrets/rotate").Handler(http.HandlerFunc(s.rotateSecretValue))
	r.Methods(http.MethodGet).Path("/v1/secretProvider/secrets/rotate/attributes").Handler(http.HandlerFunc(s.getRotateAttributes))

	return r
}

// WIP: limiting maximum request body size
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
