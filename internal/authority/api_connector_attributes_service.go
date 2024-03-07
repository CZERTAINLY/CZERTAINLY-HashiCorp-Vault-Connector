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
	attributes := make([]model.Attribute, 0)
	attributes = append(attributes, model.GetAttributeDefByUUID(model.URL_ATTR))
	attributes = append(attributes, model.GetAttributeDefByUUID(model.GROUP_CREDENTIAL_TYPE_ATTR))
	credentialTypeAttribute := model.GetAttributeDefByUUID(model.CREDENTIAL_TYPE_ATTR).(model.DataAttribute)
	credentialTypeAttribute.Content = model.GetCredentialTypes()
	attributes = append(attributes, credentialTypeAttribute)

	return model.Response(200, attributes), nil
}

// ValidateAttributes - Validate Attributes
func (s *ConnectorAttributesAPIService) ValidateAttributes(ctx context.Context, kind string, requestAttributeDto []model.Attribute) (model.ImplResponse, error) {

	return model.Response(http.StatusNotImplemented, nil), errors.New("ValidateAttributes method not implemented")
}
