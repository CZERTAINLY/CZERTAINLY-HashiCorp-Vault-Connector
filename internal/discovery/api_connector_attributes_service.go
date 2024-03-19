package discovery

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/vault"
	"context"
	"fmt"
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
	if kind != "HVault" {
		message := fmt.Sprintf("Unrecognized Authority Instance kind: %s", kind)
		return model.Response(http.StatusUnprocessableEntity, message), nil
	}

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
	var attributes []model.Attribute
	attribute := model.GetAttributeDefByUUID(model.DISCOVERY_AUTHORITY_ATTR).(model.DataAttribute)
	attribute.Content = objectContents
	attributes = append(attributes, attribute)
	attribute = model.GetAttributeDefByUUID(model.DISCOVERY_PKI_ENGINE_ATTR).(model.DataAttribute)
	attributes = append(attributes, attribute)
	return model.Response(http.StatusOK, attributes), nil

}

// ValidateAttributes - Validate Attributes
func (s *ConnectorAttributesAPIService) ValidateAttributes(ctx context.Context, kind string, requestAttributeDto []model.Attribute) (model.ImplResponse, error) {
	if kind != "HVault" {
		message := fmt.Sprintf("Unrecognized Authority Instance kind: %s", kind)
		return model.Response(http.StatusUnprocessableEntity, message), nil
	}

	authorityAttribute := model.GetAttributeFromArrayByUUID(model.DISCOVERY_AUTHORITY_ATTR, requestAttributeDto).(model.DataAttribute)
	content := authorityAttribute.GetContent()[0]
	authUuid := content.(model.ObjectAttributeContent).GetData().(map[string]any)["uuid"].(string)
	auth, _ := s.authorityRepo.FindAuthorityInstanceByUUID(authUuid)
	if auth == nil {
		return model.Response(422, []string{"Authority not found"}), nil
	}
	return model.Response(http.StatusOK, nil), nil

}

func (s *ConnectorAttributesAPIService) PkiEnginesCallback(ctx context.Context, authorityUuid string) (model.ImplResponse, error) {
	authority, err := s.authorityRepo.FindAuthorityInstanceByUUID(authorityUuid)
	if err != nil {
		return model.Response(http.StatusNotFound, model.ErrorMessageDto{
			Message: "Authority not found",
		}), nil
	}
	client, err := vault.GetClient(*authority)
	if err != nil {
		return model.Response(http.StatusInternalServerError, model.ErrorMessageDto{
			Message: "Failed to create vault client",
		}), err
	}
	mounts, _ := client.System.MountsListSecretsEngines(ctx)
	var engineList []model.AttributeContent
	for engineName, engineData := range mounts.Data {
		if engineData.(map[string]any)["type"] == "pki" {

			engineDataObject := make(map[string]interface{})
			engineDataObject["engineName"] = engineName
			engineDataObject["engineAccesor"] = engineData.(map[string]any)["accessor"]
			engineDataObject["runningPluginVersion"] = engineData.(map[string]any)["running_plugin_version"]

			engineList = append(engineList, model.ObjectAttributeContent{
				Reference: engineName,
				Data:      engineDataObject,
			})
		}
	}
	return model.Response(http.StatusOK, engineList), nil
}
