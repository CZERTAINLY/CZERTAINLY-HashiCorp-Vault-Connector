package discovery

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"context"
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

	authorities, _ := s.authorityRepo.ListAuthorityInstances()

	objectContents := make([]model.AttributeContent, 0)

	for _, authority := range authorities {
		authorityData := map[string]interface{}{
			"uuid": authority.UUID,
		}
		m := make(map[string]interface{})
		m["uuid"] = authority.UUID
		objectContents = append(objectContents,
			model.ObjectAttributeContent{
				Reference: authority.Name,
				Data:      authorityData,
			})
	}
	attribute := model.GetAttributeDefByUUID(model.AUTHORITY_ATTR).(model.DataAttribute)
	attribute.Content = objectContents
	return model.Response(http.StatusOK, []model.Attribute{attribute}), nil

}

// ValidateAttributes - Validate Attributes
func (s *ConnectorAttributesAPIService) ValidateAttributes(ctx context.Context, kind string, requestAttributeDto []model.Attribute) (model.ImplResponse, error) {
	authorityAttribute := model.GetAttributeFromArrayByUUID(model.AUTHORITY_ATTR, requestAttributeDto).(model.DataAttribute)
	content := authorityAttribute.GetContent()[0]
	authUuid := content.(model.ObjectAttributeContent).GetData().(map[string]any)["uuid"].(string)
	auth, _ := s.authorityRepo.FindAuthorityInstanceByUUID(authUuid)
	if auth == nil {
		return model.Response(422, []string{"Authority not found"}), nil
	}
	return model.Response(http.StatusOK, nil), nil

}
