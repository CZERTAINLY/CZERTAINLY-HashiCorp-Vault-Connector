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
	"bytes"
	"github.com/lib/pq"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var version = "0.0.1"

var routes map[string][]model.EndpointDto

var log = logger.Get()

func main() {
	routes = make(map[string][]model.EndpointDto)
	c := config.Get()
	log.Info("Starting CZERTAINLY-HashiCorp-Vault-Connector", zap.String("version", version))
	conn, _ := db.ConnectDB(c)
	schema := config.Get().Database.Schema
	err := conn.Exec("CREATE SCHEMA IF NOT EXISTS " + pq.QuoteIdentifier(schema)).Error
	if err != nil {
		log.Error("Error creating schema", zap.Error(err))
	}
	db.MigrateDB(c)
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
			Kinds:             []string{model.CONNECTOR_KIND},
			EndPoints:         routes["discoveryProvider"],
		},
		{
			FunctionGroupCode: "authorityProvider",
			Kinds:             []string{model.CONNECTOR_KIND},
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

	err = http.ListenAndServe(":"+c.Server.Port, topMux)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func logMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: remove body logging
		buf, _ := io.ReadAll(r.Body)
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		rdr2 := io.NopCloser(bytes.NewBuffer(buf))
		body, _ := io.ReadAll(rdr1)
		log.Debug("Request received", zap.String("path", r.URL.Path), zap.String("body", string(body)))
		r.Body = rdr2
		next.ServeHTTP(w, r)
	})
}

func populateRoutes(router *mux.Router, routeKey string) {
	routes[routeKey] = make([]model.EndpointDto, 0)
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		name := route.GetName()
		endpoint := model.EndpointDto{
			Method:   met[0],
			Name:     strings.ToLower(string(name[0])) + name[1:],
			Uuid:     utils.DeterministicGUID(met[0] + tpl),
			Context:  tpl,
			Required: true,
		}
		log.Info(strings.Join(met, ", ") + " " + tpl)
		routes[routeKey] = append(routes[routeKey], endpoint)
		return nil
	})
	if err != nil {
		log.Error("Unable to walk routers:" + err.Error())
	}
}
