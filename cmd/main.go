package main

import (
	authority "CZERTAINLY-HashiCorp-Vault-Connector/generated/authority"
	discovery "CZERTAINLY-HashiCorp-Vault-Connector/generated/discovery"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/config"
	db "CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/logger"
	"net/http"

	"go.uber.org/zap"
)

var version = "0.0.1"

func main() {
	log := logger.Get()
	c := config.Get()

	db.MigrateDB(c)
	conn,_:= db.ConnectDB(c)
	discoveryRepo,_ := db.NewDiscoveryRepository(conn)
	//authorityRepo,_ := db.NewAuthorityRepository(conn)
	
	DiscoveryAPIService := discovery.NewDiscoveryAPIService(discoveryRepo, log)
	DiscoveryAPIController := discovery.NewDiscoveryAPIController(DiscoveryAPIService)

	AuthorityManagementAPIService := authority.NewAuthorityManagementAPIService()
	AuthorityManagementAPIController := authority.NewAuthorityManagementAPIController(AuthorityManagementAPIService)

	CertificateManagementAPIService := authority.NewCertificateManagementAPIService()
	CertificateManagementAPIController := authority.NewCertificateManagementAPIController(CertificateManagementAPIService)

	
	DiscoveryConnectorAttributesAPIService := discovery.NewConnectorAttributesAPIService()
	DiscoveryConnectorAttributesAPIController := discovery.NewConnectorAttributesAPIController(DiscoveryConnectorAttributesAPIService)
	
	AuthorityConnectorAttributesAPIService := authority.NewConnectorAttributesAPIService()
	AuthorityConnectorAttributesAPIController := authority.NewConnectorAttributesAPIController(AuthorityConnectorAttributesAPIService)
	
	// HealthCheckAPIService := discovery.NewHealthCheckAPIService()
	// HealthCheckAPIController := discovery.NewHealthCheckAPIController(HealthCheckAPIService)
	
	// ConnectorInfoAPIService := discovery.NewConnectorInfoAPIService()
	// ConnectorInfoAPIController := discovery.NewConnectorInfoAPIController(ConnectorInfoAPIService)
	
	topMux := http.NewServeMux()

	

	topMux.Handle("/v1", logMiddleware(discovery.NewRouter(ConnectorInfoAPIController)))
	topMux.Handle("/v1/", logMiddleware(discovery.NewRouter(HealthCheckAPIController)))

	topMux.Handle("/v1/authorityProvider/", logMiddleware(authority.NewRouter(AuthorityConnectorAttributesAPIController, AuthorityManagementAPIController, CertificateManagementAPIController)))
	topMux.Handle("/v1/discoveryProvider/", logMiddleware(discovery.NewRouter(DiscoveryConnectorAttributesAPIController, DiscoveryAPIController)))

	err := http.ListenAndServe(":"+c.Server.Port, topMux)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func logMiddleware(next http.Handler) http.Handler {
	l := logger.Get()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("Request received", zap.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
