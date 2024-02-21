package main

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/cmd/config"
	"CZERTAINLY-HashiCorp-Vault-Connector/cmd/logger"
	authority "CZERTAINLY-HashiCorp-Vault-Connector/generated/authority"
	discovery "CZERTAINLY-HashiCorp-Vault-Connector/generated/discovery"
	"go.uber.org/zap"
	"log"
	"net/http"
)

var version = "0.0.1"

func main() {
	l := logger.Get()
	c := config.Get()

	l.Info("Starting the server version: " + version)

	DiscoveryAPIService := discovery.NewDiscoveryAPIService()
	DiscoveryAPIController := discovery.NewDiscoveryAPIController(DiscoveryAPIService)

	AuthorityManagementAPIService := authority.NewAuthorityManagementAPIService()
	AuthorityManagementAPIController := authority.NewAuthorityManagementAPIController(AuthorityManagementAPIService)

	CertificateManagementAPIService := authority.NewCertificateManagementAPIService()
	CertificateManagementAPIController := authority.NewCertificateManagementAPIController(CertificateManagementAPIService)

	HealthCheckAPIService := discovery.NewHealthCheckAPIService()
	HealthCheckAPIController := discovery.NewHealthCheckAPIController(HealthCheckAPIService)

	ConnectorAttributesAPIService := discovery.NewConnectorAttributesAPIService()
	ConnectorAttributesAPIController := discovery.NewConnectorAttributesAPIController(ConnectorAttributesAPIService)

	ConnectorInfoAPIService := discovery.NewConnectorInfoAPIService()
	ConnectorInfoAPIController := discovery.NewConnectorInfoAPIController(ConnectorInfoAPIService)


	topMux := http.NewServeMux()
	topMux.Handle("/v1", logMiddleware(discovery.NewRouter(ConnectorInfoAPIController)))
	topMux.Handle("/v1/", logMiddleware(discovery.NewRouter(ConnectorAttributesAPIController, ConnectorInfoAPIController, HealthCheckAPIController)))
	topMux.Handle("/v1/authorityProvider/", logMiddleware(authority.NewRouter(AuthorityManagementAPIController, CertificateManagementAPIController)))
	topMux.Handle("/v1/discoveryProvider/", logMiddleware(discovery.NewRouter(DiscoveryAPIController)))


	log.Fatal(http.ListenAndServe(":"+c.Server.Port, topMux))

}

func logMiddleware(next http.Handler) http.Handler {
	l := logger.Get()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("Request received", zap.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
