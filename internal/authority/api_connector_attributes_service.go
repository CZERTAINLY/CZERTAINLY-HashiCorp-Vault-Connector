package authority

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"context"
	// "encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// ConnectorAttributesAPIService is a service that implements the logic for the ConnectorAttributesAPIServicer
// This service should implement the business logic for every endpoint for the ConnectorAttributesAPI API.
// Include any external packages or services that will be required by this service.
type ConnectorAttributesAPIService struct {
	authorityRepo *db.AuthorityRepository
	log           *zap.Logger
}

// NewConnectorAttributesAPIService creates a default api service
func NewConnectorAttributesAPIService(authorityRepo *db.AuthorityRepository, logger *zap.Logger) ConnectorAttributesAPIServicer {
	return &ConnectorAttributesAPIService{
		authorityRepo: authorityRepo,
		log:           logger,
	}
}

// ListAttributeDefinitions - List available Attributes
func (s *ConnectorAttributesAPIService) ListAttributeDefinitions(ctx context.Context, kind string) (model.ImplResponse, error) {

	// TODO - update ListAttributeDefinitions with the required logic for this service method.
	// Add api_connector_attributes_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response model.Response(400, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(400, ErrorMessageDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(500, {}) or use other options such as http.Ok ...
	// return model.Response(500, nil),nil

	// TODO: Uncomment the next line to return response model.Response(404, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(404, ErrorMessageDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(200, []BaseAttributeDto{}) or use other options such as http.Ok ...
	// return model.Response(200, []BaseAttributeDto{}), nil

	//return model.Response(200, model.GetAttributeList()), nil
	return model.Response(200, nil), nil
}

// ValidateAttributes - Validate Attributes
func (s *ConnectorAttributesAPIService) ValidateAttributes(ctx context.Context, kind string, requestAttributeDto []model.RequestAttributeDto) (model.ImplResponse, error) {
	// TODO - update ValidateAttributes with the required logic for this service method.
	// Add api_connector_attributes_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response model.Response(200, {}) or use other options such as http.Ok ...
	// return model.Response(200, nil),nil

	// TODO: Uncomment the next line to return response model.Response(422, []string{}) or use other options such as http.Ok ...
	// return model.Response(422, []string{}), nil

	// TODO: Uncomment the next line to return response model.Response(400, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(400, ErrorMessageDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(500, {}) or use other options such as http.Ok ...
	// return model.Response(500, nil),nil

	// TODO: Uncomment the next line to return response model.Response(404, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(404, ErrorMessageDto{}), nil

	return model.Response(http.StatusNotImplemented, nil), errors.New("ValidateAttributes method not implemented")
}
