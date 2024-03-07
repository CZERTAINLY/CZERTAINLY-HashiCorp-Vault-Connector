package discovery

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"context"
)

// ConnectorInfoAPIService is a service that implements the logic for the ConnectorInfoAPIServicer
// This service should implement the business logic for every endpoint for the ConnectorInfoAPI API.
// Include any external packages or services that will be required by this service.
type ConnectorInfoAPIService struct {
	endpoints []model.EndpointDto
}

// NewConnectorInfoAPIService creates a default api service
func NewConnectorInfoAPIService() ConnectorInfoAPIServicer {
	return &ConnectorInfoAPIService{}
}

// ListSupportedFunctions - List supported functions of the connector
func (s *ConnectorInfoAPIService) ListSupportedFunctions(ctx context.Context) (model.ImplResponse, error) {
	// TODO - update ListSupportedFunctions with the required logic for this service method.
	// Add api_connector_info_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response model.Response(400, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(400, ErrorMessageDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(500, {}) or use other options such as http.Ok ...
	// return model.Response(500, nil),nil

	infoResponses := []model.InfoResponse{}

	return model.Response(200, infoResponses), nil

	// TODO: Uncomment the next line to return response model.Response(404, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(404, ErrorMessageDto{}), nil
}
