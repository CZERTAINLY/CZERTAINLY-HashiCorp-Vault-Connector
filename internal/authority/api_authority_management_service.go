package authority

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// AuthorityManagementAPIService is a service that implements the logic for the AuthorityManagementAPIServicer
// This service should implement the business logic for every endpoint for the AuthorityManagementAPI API.
// Include any external packages or services that will be required by this service.
type AuthorityManagementAPIService struct {
	authorityRepo *db.AuthorityRepository
	log           *zap.Logger
}

// NewAuthorityManagementAPIService creates a default api service
func NewAuthorityManagementAPIService(authorityRepo *db.AuthorityRepository, logger *zap.Logger) AuthorityManagementAPIServicer {
	return &AuthorityManagementAPIService{
		authorityRepo: authorityRepo,
		log:           logger,
	}
}

// CreateAuthorityInstance - Create Authority instance
func (s *AuthorityManagementAPIService) CreateAuthorityInstance(ctx context.Context, request model.AuthorityProviderInstanceRequestDto) (model.ImplResponse, error) {
	attributes := request.Attributes
	URL := model.GetAttributeFromArrayByUUID(model.URL_ATTR, attributes).GetContent()[0].GetData().(string)
	credentialType := model.GetAttributeFromArrayByUUID(model.CREDENTIAL_TYPE_ATTR, attributes).GetContent()[0].GetData().(string)
	authorityName := request.Name
	marshaledAttrs, err := json.Marshal(attributes)
	if err != nil {
		return model.Response(500, model.ErrorMessageDto{
			Message: "Failed to marshal attributes",
		}), err
	}
	authority := db.AuthorityInstance{
		UUID:           utils.DeterministicGUID(authorityName),
		Name:           authorityName,
		URL:            URL,
		Attributes:     string(marshaledAttrs),
		CredentialType: credentialType,
	}
	err = s.authorityRepo.CreateAuthorityInstance(&authority)
	if err != nil {
		return model.Response(500, model.ErrorMessageDto{
			Message: "Failed to create authority",
		}), err
	}
	dto := model.AuthorityProviderInstanceDto{
		Uuid:       authority.UUID,
		Name:       authority.Name,
		Attributes: attributes,
	}
	return model.Response(200, dto), nil
}

// GetAuthorityInstance - Get an Authority instance
func (s *AuthorityManagementAPIService) GetAuthorityInstance(ctx context.Context, uuid string) (model.ImplResponse, error) {
	authority, err := s.authorityRepo.FindAuthorityInstanceByUUID(uuid)
	if err != nil {
		return model.Response(404, model.ErrorMessageDto{
			Message: "Authority not found",
		}), nil
	}
	attributes := model.UnmarshalAttributes([]byte(authority.Attributes))
	authorityDto := model.AuthorityProviderInstanceDto{
		Uuid:       authority.UUID,
		Name:       authority.Name,
		Attributes: attributes,
	}
	return model.Response(200, authorityDto), nil
}

// GetCaCertificates - Get the Authority Instance&#39;s certificate chain
func (s *AuthorityManagementAPIService) GetCaCertificates(ctx context.Context, uuid string, caCertificatesRequestDto model.CaCertificatesRequestDto) (model.ImplResponse, error) {
	// TODO - update GetCaCertificates with the required logic for this service method.
	// Add api_authority_management_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response model.Response(200, CaCertificatesResponseDto{}) or use other options such as http.Ok ...
	// return model.Response(200, CaCertificatesResponseDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(400, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(400, ErrorMessageDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(500, {}) or use other options such as http.Ok ...
	// return model.Response(500, nil),nil

	// TODO: Uncomment the next line to return response model.Response(404, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(404, ErrorMessageDto{}), nil

	return model.Response(http.StatusNotImplemented, nil), errors.New("GetCaCertificates method not implemented")
}

// GetConnection - Connect to Authority
func (s *AuthorityManagementAPIService) GetConnection(ctx context.Context, uuid string) (model.ImplResponse, error) {
	// TODO - update GetConnection with the required logic for this service method.
	// Add api_authority_management_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response model.Response(204, {}) or use other options such as http.Ok ...
	// return model.Response(204, nil),nil

	// TODO: Uncomment the next line to return response model.Response(400, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(400, ErrorMessageDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(500, {}) or use other options such as http.Ok ...
	// return model.Response(500, nil),nil

	// TODO: Uncomment the next line to return response model.Response(404, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(404, ErrorMessageDto{}), nil

	return model.Response(http.StatusNotImplemented, nil), errors.New("GetConnection method not implemented")
}

// GetCrl - Get the latest CRL for the Authority Instance
func (s *AuthorityManagementAPIService) GetCrl(ctx context.Context, uuid string, certificateRevocationListRequestDto model.CertificateRevocationListRequestDto) (model.ImplResponse, error) {
	// TODO - update GetCrl with the required logic for this service method.
	// Add api_authority_management_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response model.Response(400, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(400, ErrorMessageDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(500, {}) or use other options such as http.Ok ...
	// return model.Response(500, nil),nil

	// TODO: Uncomment the next line to return response model.Response(404, ErrorMessageDto{}) or use other options such as http.Ok ...
	// return model.Response(404, ErrorMessageDto{}), nil

	// TODO: Uncomment the next line to return response model.Response(200, CertificateRevocationListResponseDto{}) or use other options such as http.Ok ...
	// return model.Response(200, CertificateRevocationListResponseDto{}), nil

	return model.Response(http.StatusNotImplemented, nil), errors.New("GetCrl method not implemented")
}

// ListAuthorityInstances - List Authority instances
func (s *AuthorityManagementAPIService) ListAuthorityInstances(ctx context.Context) (model.ImplResponse, error) {
	authorities, _ := s.authorityRepo.ListAuthorityInstances()
	var authoritiesDto []model.AuthorityProviderInstanceDto
	for _, authority := range authorities {
		attributes := model.UnmarshalAttributes([]byte(authority.Attributes))
		authoritiesDto = append(authoritiesDto, model.AuthorityProviderInstanceDto{
			Uuid:       authority.UUID,
			Name:       authority.Name,
			Attributes: attributes,
		})
	}

	return model.Response(200, authoritiesDto), nil
}

// ListRAProfileAttributes - List RA Profile Attributes
func (s *AuthorityManagementAPIService) ListRAProfileAttributes(ctx context.Context, uuid string) (model.ImplResponse, error) {
	return model.Response(200, []model.BaseAttributeDto{}), nil
}

// RemoveAuthorityInstance - Remove Authority instance
func (s *AuthorityManagementAPIService) RemoveAuthorityInstance(ctx context.Context, uuid string) (model.ImplResponse, error) {
	authority, err := s.authorityRepo.FindAuthorityInstanceByUUID(uuid)
	if err != nil {
		return model.Response(204, nil), nil
	}

	// Delete the authority if it has been found
	err = s.authorityRepo.DeleteAuthorityInstance(authority)
	if err != nil {
		// Handle error, failed to delete authority
		return model.Response(500, model.ErrorMessageDto{}), err
	}

	// Return success response
	return model.Response(200, nil), nil

}

// UpdateAuthorityInstance - Update Authority instance
func (s *AuthorityManagementAPIService) UpdateAuthorityInstance(ctx context.Context, uuid string, request model.AuthorityProviderInstanceRequestDto) (model.ImplResponse, error) {
	authority, err := s.authorityRepo.FindAuthorityInstanceByUUID(uuid)
	if err != nil {
		return model.Response(500, model.ErrorMessageDto{
			Message: "Failed to marshal attributes",
		}), err
	}
	attributes := request.Attributes
	URL := model.GetAttributeFromArrayByUUID(model.URL_ATTR, attributes).GetContent()[0].GetData().(string)
	credentialType := model.GetAttributeFromArrayByUUID(model.CREDENTIAL_TYPE_ATTR, attributes).GetContent()[0].GetData().(string)
	authorityName := request.Name

	marshaledAttrs, err := json.Marshal(attributes)
	if err != nil {
		return model.Response(500, model.ErrorMessageDto{
			Message: "Failed to marshal attributes",
		}), err
	}
	authority.Name = authorityName
	authority.URL = URL
	authority.CredentialType = credentialType
	authority.Attributes = string(marshaledAttrs)

	err = s.authorityRepo.UpdateAuthorityInstance(authority)
	if err != nil {
		// Handle error, failed to delete authority
		return model.Response(500, model.ErrorMessageDto{}), err
	}
	attributesEntity := model.UnmarshalAttributes([]byte(authority.Attributes))
	authorityDto := model.AuthorityProviderInstanceDto{
		Uuid:       authority.UUID,
		Name:       authority.Name,
		Attributes: attributesEntity,
	}
	return model.Response(200, authorityDto), nil

}

// ValidateRAProfileAttributes - Validate RA Profile attributes
func (s *AuthorityManagementAPIService) ValidateRAProfileAttributes(ctx context.Context, uuid string, requestAttributeDto []model.RequestAttributeDto) (model.ImplResponse, error) {
	return model.Response(200, nil), nil
}
