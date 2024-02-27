package main

import (
	authority "CZERTAINLY-HashiCorp-Vault-Connector/generated/authority"
	discovery "CZERTAINLY-HashiCorp-Vault-Connector/generated/discovery"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/config"
	db "CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/logger"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var version = "0.0.1"

type InfoResponse struct {
	FunctionGroupCode string `json:"functionGroupCode"`

	// List of supported functional group kinds
	Kinds []string `json:"kinds"`

	// List of end points related to functional group
	EndPoints []EndpointDto `json:"endPoints"`
}

type EndpointDto struct {

	// Object identifier
	Uuid string `json:"uuid"`

	// Object Name
	Name string `json:"name"`

	// Context of the Endpoint
	Context string `json:"context"`

	// Method to be used for the Endpoint
	Method string `json:"method"`

	// True if the Endpoint is required for implementation
	Required bool `json:"required"`
}

// HealthStatus : Current connector operational status
type HealthStatus string

// List of HealthStatus
const (
	OK      HealthStatus = "ok"
	NOK     HealthStatus = "nok"
	UNKNOWN HealthStatus = "unknown"
)

type HealthDto struct {
	Status HealthStatus `json:"status"`

	// Detailed status description
	Description string `json:"description,omitempty"`

	// Nested status of services
	Parts map[string]HealthDto `json:"parts,omitempty"`
}

var routes map[string][]EndpointDto

func main() {
	routes = make(map[string][]EndpointDto)
	log := logger.Get()
	c := config.Get()

	db.MigrateDB(c)
	conn, _ := db.ConnectDB(c)
	discoveryRepo, _ := db.NewDiscoveryRepository(conn)
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

	topMux := http.NewServeMux()

	infoRouter := mux.NewRouter()
	infoRouter.HandleFunc("/v1", infoHandler).Methods("GET")
	populateRoutes(infoRouter, "info")
	topMux.Handle("/v1", logMiddleware(infoRouter))


	healthRouter := mux.NewRouter()
	healthRouter.HandleFunc("/v1/health", healthHandler).Methods("GET")
	populateRoutes(healthRouter, "health")
	topMux.Handle("/v1/health", logMiddleware(healthRouter))


	authorityRouter := authority.NewRouter(AuthorityConnectorAttributesAPIController, AuthorityManagementAPIController, CertificateManagementAPIController)
	populateRoutes(authorityRouter, "authorityProvider")

	discoveryRouter := discovery.NewRouter(DiscoveryConnectorAttributesAPIController, DiscoveryAPIController)
	populateRoutes(discoveryRouter, "discoveryProvider")

	topMux.Handle("/v1/authorityProvider/", logMiddleware(authorityRouter))
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HealthDto{
		Status:      OK,
		Description: "Service is running properly",
		Parts:       nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	info := []InfoResponse{
		{
			FunctionGroupCode: "authorityProvider",
			Kinds:             []string{"Vault"},
			EndPoints:         routes["discoveryProvider"],
		},
		{
			FunctionGroupCode: "authorityProvider",
			Kinds:             []string{"Vault"},
			EndPoints:         routes["authorityProvider"],
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func populateRoutes(router *mux.Router, routeKey string) {
	log := logger.Get()
	routes[routeKey] = make([]EndpointDto, 0)
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		endpoint := EndpointDto{
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
