package secret

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"CZERTAINLY-HashiCorp-Vault-Connector/internal/metrics"
	sv "CZERTAINLY-HashiCorp-Vault-Connector/internal/secret/vault"

	"github.com/gorilla/mux"
)

const DEFAULT_K8S_TOKEN_PATH = "/var/run/secrets/kubernetes.io/serviceaccount/token"

type Server struct {
	m        *sv.Manager
	k8sToken *string
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
			str := string(b)
			s.k8sToken = &str
		}
	}

	return s
}

func (s *Server) MuxRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Methods(http.MethodPost).Path("/v1/secretProvider/secrets").Handler(metrics.Middleware()(http.HandlerFunc(s.createSecret)))
	r.Methods(http.MethodPut).Path("/v1/secretProvider/secrets").Handler(metrics.Middleware()(http.HandlerFunc(s.updateSecret)))
	r.Methods(http.MethodPost).Path("/v1/secretProvider/secrets/content").Handler(metrics.Middleware()(http.HandlerFunc(s.getSecretValue)))
	r.Methods(http.MethodDelete).Path("/v1/secretProvider/secrets").Handler(metrics.Middleware()(http.HandlerFunc(s.deleteSecret)))

	r.Methods(http.MethodGet).Path("/v1/secretProvider/vaults/attributes").Handler(metrics.Middleware()(http.HandlerFunc(s.listVaultAttributes)))
	r.Methods(http.MethodGet).Path("/v1/secretProvider/secrets/{secretType}/attributes").Handler(metrics.Middleware()(http.HandlerFunc(getSecretAttributes)))
	r.Methods(http.MethodGet).Path("/v1/secretProvider/credentialType/{credentialsType}/callback").Handler(metrics.Middleware()(http.HandlerFunc(s.credentialsType)))

	r.Methods(http.MethodPost).Path("/v1/secretProvider/vaults").Handler(metrics.Middleware()(http.HandlerFunc(s.checkVaultConnection)))

	// Not implemented for now - returns HTTP 501 Not Implemented
	r.Methods(http.MethodPost).Path("/v1/secretProvider/secrets/rotate").Handler(metrics.Middleware()(http.HandlerFunc(s.rotateSecretValue)))
	r.Methods(http.MethodGet).Path("/v1/secretProvider/secrets/rotate/attributes").Handler(metrics.Middleware()(http.HandlerFunc(s.getRotateAttributes)))

	return r
}
