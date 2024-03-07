package discovery

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"context"
	"net/http"
)

// ConnectorAttributesAPIRouter defines the required methods for binding the api requests to a responses for the ConnectorAttributesAPI
// The ConnectorAttributesAPIRouter implementation should parse necessary information from the http request,
// pass the data to a ConnectorAttributesAPIServicer to perform the required actions, then write the service results to the http response.
type ConnectorAttributesAPIRouter interface {
	ListAttributeDefinitions(http.ResponseWriter, *http.Request)
	ValidateAttributes(http.ResponseWriter, *http.Request)
}

// ConnectorInfoAPIRouter defines the required methods for binding the api requests to a responses for the ConnectorInfoAPI
// The ConnectorInfoAPIRouter implementation should parse necessary information from the http request,
// pass the data to a ConnectorInfoAPIServicer to perform the required actions, then write the service results to the http response.
type ConnectorInfoAPIRouter interface {
	ListSupportedFunctions(http.ResponseWriter, *http.Request)
}

// DiscoveryAPIRouter defines the required methods for binding the api requests to a responses for the DiscoveryAPI
// The DiscoveryAPIRouter implementation should parse necessary information from the http request,
// pass the data to a DiscoveryAPIServicer to perform the required actions, then write the service results to the http response.
type DiscoveryAPIRouter interface {
	DeleteDiscovery(http.ResponseWriter, *http.Request)
	DiscoverCertificate(http.ResponseWriter, *http.Request)
	GetDiscovery(http.ResponseWriter, *http.Request)
}

// HealthCheckAPIRouter defines the required methods for binding the api requests to a responses for the HealthCheckAPI
// The HealthCheckAPIRouter implementation should parse necessary information from the http request,
// pass the data to a HealthCheckAPIServicer to perform the required actions, then write the service results to the http response.
type HealthCheckAPIRouter interface {
	CheckHealth(http.ResponseWriter, *http.Request)
}

// ConnectorAttributesAPIServicer defines the api actions for the ConnectorAttributesAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ConnectorAttributesAPIServicer interface {
	ListAttributeDefinitions(context.Context, string) (model.ImplResponse, error)
	ValidateAttributes(context.Context, string, []model.RequestAttributeDto) (model.ImplResponse, error)
}

// ConnectorInfoAPIServicer defines the api actions for the ConnectorInfoAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ConnectorInfoAPIServicer interface {
	ListSupportedFunctions(context.Context) (model.ImplResponse, error)
}

// DiscoveryAPIServicer defines the api actions for the DiscoveryAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type DiscoveryAPIServicer interface {
	DeleteDiscovery(context.Context, string) (model.ImplResponse, error)
	DiscoverCertificate(context.Context, model.DiscoveryRequestDto) (model.ImplResponse, error)
	GetDiscovery(context.Context, string, model.DiscoveryDataRequestDto) (model.ImplResponse, error)
}

// HealthCheckAPIServicer defines the api actions for the HealthCheckAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type HealthCheckAPIServicer interface {
	CheckHealth(context.Context) (model.ImplResponse, error)
}
