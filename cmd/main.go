package main

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/authority"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/config"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/connectorInfo"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/discovery"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/health"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/logger"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/utils"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var version = "0.0.1"

var routes map[string][]model.EndpointDto

func main() {
	routes = make(map[string][]model.EndpointDto)
	log := logger.Get()
	c := config.Get()

	db.MigrateDB(c)
	conn, _ := db.ConnectDB(c)
	discoveryRepo, _ := db.NewDiscoveryRepository(conn)
	authorityRepo, _ := db.NewAuthorityRepository(conn)

	DiscoveryAPIService := discovery.NewDiscoveryAPIService(discoveryRepo, authorityRepo, log)
	DiscoveryAPIController := discovery.NewDiscoveryAPIController(DiscoveryAPIService)

	AuthorityManagementAPIService := authority.NewAuthorityManagementAPIService(authorityRepo, log)
	AuthorityManagementAPIController := authority.NewAuthorityManagementAPIController(AuthorityManagementAPIService)

	CertificateManagementAPIService := authority.NewCertificateManagementAPIService(authorityRepo, log)
	CertificateManagementAPIController := authority.NewCertificateManagementAPIController(CertificateManagementAPIService)

	DiscoveryConnectorAttributesAPIService := discovery.NewConnectorAttributesAPIService(authorityRepo, log)
	DiscoveryConnectorAttributesAPIController := discovery.NewConnectorAttributesAPIController(DiscoveryConnectorAttributesAPIService)

	AuthorityConnectorAttributesAPIService := authority.NewConnectorAttributesAPIService(authorityRepo, log)
	AuthorityConnectorAttributesAPIController := authority.NewConnectorAttributesAPIController(AuthorityConnectorAttributesAPIService)

	HealthAPIService := health.NewHealthCheckAPIService()
	HealthAPIController := health.NewHealthCheckAPIController(HealthAPIService)

	topMux := http.NewServeMux()

	healthRouter := model.NewRouter(HealthAPIController)

	authorityRouter := model.NewRouter(AuthorityConnectorAttributesAPIController, AuthorityManagementAPIController, CertificateManagementAPIController)
	populateRoutes(authorityRouter, "authorityProvider")

	// needs to be separate as it uses v2 prefix!
	certificateRouter := model.NewRouter(CertificateManagementAPIController)
	populateRoutes(authorityRouter, "authorityProvider")

	discoveryRouter := model.NewRouter(DiscoveryConnectorAttributesAPIController, DiscoveryAPIController)
	populateRoutes(discoveryRouter, "discoveryProvider")

	info := []model.InfoResponse{
		{
			FunctionGroupCode: "discoveryProvider",
			Kinds:             []string{"Vault"},
			EndPoints:         routes["discoveryProvider"],
		},
		{
			FunctionGroupCode: "authorityProvider",
			Kinds:             []string{"Vault"},
			EndPoints:         routes["authorityProvider"],
		},
	}

	ConnectorInfoAPIService := connectorInfo.NewConnectorInfoAPIService(info)
	ConnectorInfoAPIController := connectorInfo.NewConnectorInfoAPIController(ConnectorInfoAPIService)
	connectorInfoRouter := model.NewRouter(ConnectorInfoAPIController)

	topMux.Handle("/v1", logMiddleware(connectorInfoRouter))
	topMux.Handle("/v1/", logMiddleware(healthRouter))
	topMux.Handle("/v1/authorityProvider/", logMiddleware(authorityRouter))
	topMux.Handle("/v2/authorityProvider/", logMiddleware(certificateRouter))
	topMux.Handle("/v1/discoveryProvider/", logMiddleware(discoveryRouter))

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

func populateRoutes(router *mux.Router, routeKey string) {
	log := logger.Get()
	routes[routeKey] = make([]model.EndpointDto, 0)
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		endpoint := model.EndpointDto{
			Method:   met[0],
			Uuid:     utils.DeterministicGUID(met[0] + tpl),
			Context:  tpl,
			Required: true,
		}
		log.Info(strings.Join(met, ", ") + " " + tpl)
		routes[routeKey] = append(routes[routeKey], endpoint)
		return nil
	})
}
