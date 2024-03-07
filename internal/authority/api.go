package authority

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"context"
	"net/http"
)

// AuthorityManagementAPIRouter defines the required methods for binding the api requests to a responses for the AuthorityManagementAPI
// The AuthorityManagementAPIRouter implementation should parse necessary information from the http request,
// pass the data to a AuthorityManagementAPIServicer to perform the required actions, then write the service results to the http response.
type AuthorityManagementAPIRouter interface {
	CreateAuthorityInstance(http.ResponseWriter, *http.Request)
	GetAuthorityInstance(http.ResponseWriter, *http.Request)
	GetCaCertificates(http.ResponseWriter, *http.Request)
	GetConnection(http.ResponseWriter, *http.Request)
	GetCrl(http.ResponseWriter, *http.Request)
	ListAuthorityInstances(http.ResponseWriter, *http.Request)
	ListRAProfileAttributes(http.ResponseWriter, *http.Request)
	RemoveAuthorityInstance(http.ResponseWriter, *http.Request)
	UpdateAuthorityInstance(http.ResponseWriter, *http.Request)
	ValidateRAProfileAttributes(http.ResponseWriter, *http.Request)
}

// CertificateManagementAPIRouter defines the required methods for binding the api requests to a responses for the CertificateManagementAPI
// The CertificateManagementAPIRouter implementation should parse necessary information from the http request,
// pass the data to a CertificateManagementAPIServicer to perform the required actions, then write the service results to the http response.
type CertificateManagementAPIRouter interface {
	IdentifyCertificate(http.ResponseWriter, *http.Request)
	IssueCertificate(http.ResponseWriter, *http.Request)
	ListIssueCertificateAttributes(http.ResponseWriter, *http.Request)
	ListRevokeCertificateAttributes(http.ResponseWriter, *http.Request)
	RenewCertificate(http.ResponseWriter, *http.Request)
	RevokeCertificate(http.ResponseWriter, *http.Request)
	ValidateIssueCertificateAttributes(http.ResponseWriter, *http.Request)
	ValidateRevokeCertificateAttributes(http.ResponseWriter, *http.Request)
}

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

// HealthCheckAPIRouter defines the required methods for binding the api requests to a responses for the HealthCheckAPI
// The HealthCheckAPIRouter implementation should parse necessary information from the http request,
// pass the data to a HealthCheckAPIServicer to perform the required actions, then write the service results to the http response.
type HealthCheckAPIRouter interface {
	CheckHealth(http.ResponseWriter, *http.Request)
}

// AuthorityManagementAPIServicer defines the api actions for the AuthorityManagementAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AuthorityManagementAPIServicer interface {
	CreateAuthorityInstance(context.Context, model.AuthorityProviderInstanceRequestDto) (model.ImplResponse, error)
	GetAuthorityInstance(context.Context, string) (model.ImplResponse, error)
	GetCaCertificates(context.Context, string, model.CaCertificatesRequestDto) (model.ImplResponse, error)
	GetConnection(context.Context, string) (model.ImplResponse, error)
	GetCrl(context.Context, string, model.CertificateRevocationListRequestDto) (model.ImplResponse, error)
	ListAuthorityInstances(context.Context) (model.ImplResponse, error)
	ListRAProfileAttributes(context.Context, string) (model.ImplResponse, error)
	RemoveAuthorityInstance(context.Context, string) (model.ImplResponse, error)
	UpdateAuthorityInstance(context.Context, string, model.AuthorityProviderInstanceRequestDto) (model.ImplResponse, error)
	ValidateRAProfileAttributes(context.Context, string, []model.RequestAttributeDto) (model.ImplResponse, error)
}

// CertificateManagementAPIServicer defines the api actions for the CertificateManagementAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type CertificateManagementAPIServicer interface {
	IdentifyCertificate(context.Context, string, model.CertificateIdentificationRequestDto) (model.ImplResponse, error)
	IssueCertificate(context.Context, string, model.CertificateSignRequestDto) (model.ImplResponse, error)
	ListIssueCertificateAttributes(context.Context, string) (model.ImplResponse, error)
	ListRevokeCertificateAttributes(context.Context, string) (model.ImplResponse, error)
	RenewCertificate(context.Context, string, model.CertificateRenewRequestDto) (model.ImplResponse, error)
	RevokeCertificate(context.Context, string, model.CertRevocationDto) (model.ImplResponse, error)
	ValidateIssueCertificateAttributes(context.Context, string, []model.RequestAttributeDto) (model.ImplResponse, error)
	ValidateRevokeCertificateAttributes(context.Context, string, []model.RequestAttributeDto) (model.ImplResponse, error)
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

// HealthCheckAPIServicer defines the api actions for the HealthCheckAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type HealthCheckAPIServicer interface {
	CheckHealth(context.Context) (model.ImplResponse, error)
}
