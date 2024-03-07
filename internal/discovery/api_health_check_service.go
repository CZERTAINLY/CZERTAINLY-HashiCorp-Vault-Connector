package discovery

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"context"
	"errors"
	"net/http"
)

// HealthCheckAPIService is a service that implements the logic for the HealthCheckAPIServicer
// This service should implement the business logic for every endpoint for the HealthCheckAPI API.
// Include any external packages or services that will be required by this service.
type HealthCheckAPIService struct {
}

// NewHealthCheckAPIService creates a default api service
func NewHealthCheckAPIService() HealthCheckAPIServicer {
	return &HealthCheckAPIService{}
}

// CheckHealth - Health check
func (s *HealthCheckAPIService) CheckHealth(ctx context.Context) (model.ImplResponse, error) {
	// TODO - update CheckHealth with the required logic for this service method.
	// Add api_health_check_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response model.Response(400, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(400, ErrorMessageDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(200, HealthDto{}) or use other options such as http.Ok ...
	// return model.Response(200, HealthDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(500, {}) or use other options such as http.Ok ...
	// return model.Response(500, nil),nil

	// TODO: Uncomment the next line to return response model.Response(404, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(404, ErrorMessageDto{}), nil

	return model.Response(http.StatusNotImplemented, nil), errors.New("CheckHealth method not implemented")
}
