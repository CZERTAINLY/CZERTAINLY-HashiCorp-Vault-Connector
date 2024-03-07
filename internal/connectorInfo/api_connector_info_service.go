package connectorInfo

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"context"
)

// ConnectorInfoAPIService is a service that implements the logic for the ConnectorInfoAPIServicer
// This service should implement the business logic for every endpoint for the ConnectorInfoAPI API.
// Include any external packages or services that will be required by this service.
type ConnectorInfoAPIService struct {
	routes []model.InfoResponse
}

// NewConnectorInfoAPIService creates a default api service
func NewConnectorInfoAPIService(routes []model.InfoResponse) ConnectorInfoAPIServicer {
	return &ConnectorInfoAPIService{
		routes: routes,
	}
}

// ListSupportedFunctions - List supported functions of the connector
func (s *ConnectorInfoAPIService) ListSupportedFunctions(ctx context.Context) (model.ImplResponse, error) {
	return model.Response(200, s.routes), nil
}
