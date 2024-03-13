package discovery

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/db"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/utils"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/vault"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// DiscoveryAPIService is a service that implements the logic for the DiscoveryAPIServicer
// This service should implement the business logic for every endpoint for the DiscoveryAPI API.
// Include any external packages or services that will be required by this service.
type DiscoveryAPIService struct {
	discoveryRepo *db.DiscoveryRepository
	authorityRepo *db.AuthorityRepository
	log           *zap.Logger
}

// NewDiscoveryAPIService creates a default api service
func NewDiscoveryAPIService(discoveryRepo *db.DiscoveryRepository, authorityRepo *db.AuthorityRepository, logger *zap.Logger) DiscoveryAPIServicer {
	return &DiscoveryAPIService{
		discoveryRepo: discoveryRepo,
		authorityRepo: authorityRepo,
		log:           logger,
	}
}

// DeleteDiscovery - Delete Discovery
func (s *DiscoveryAPIService) DeleteDiscovery(ctx context.Context, uuid string) (model.ImplResponse, error) {
	discovery, err := s.discoveryRepo.FindDiscoveryByUUID(uuid)
	if err != nil {
		return model.Response(http.StatusNotFound, model.ErrorMessageDto{Message: "Discovery " + uuid + " not found."}), nil
	}
	s.discoveryRepo.DeleteDiscovery(discovery)

	return model.Response(204, nil), nil
}

// DiscoverCertificate - Initiate certificate Discovery
func (s *DiscoveryAPIService) DiscoverCertificate(ctx context.Context, discoveryRequestDto model.DiscoveryRequestDto) (model.ImplResponse, error) {
	id := uuid.New()
	fmt.Println(id)
	response := model.DiscoveryProviderDto{
		Uuid:                        utils.DeterministicGUID("name"),
		Name:                        discoveryRequestDto.Name,
		Status:                      model.IN_PROGRESS,
		TotalCertificatesDiscovered: 0,
		CertificateData:             nil,
		Meta:                        nil,
	}
	discovery := &db.Discovery{
		UUID:         response.Uuid,
		Name:         response.Name,
		Status:       string(response.Status),
		Meta:         nil,
		Certificates: nil,
	}
	s.discoveryRepo.CreateDiscovery(discovery)
	go s.DiscoveryCertificates(&db.AuthorityInstance{}, discovery)

	return model.Response(http.StatusOK, response), nil
}

// GetDiscovery - Get Discovery status and result
func (s *DiscoveryAPIService) GetDiscovery(ctx context.Context, uuid string, discoveryDataRequestDto model.DiscoveryDataRequestDto) (model.ImplResponse, error) {
	discovery, err := s.discoveryRepo.FindDiscoveryByUUID(uuid)
	if err != nil {
		return model.Response(http.StatusNotFound, model.ErrorMessageDto{Message: "Discovery " + uuid + " not found."}), nil
	}
	if discovery.Status == "IN_PROGRESS" {
		return model.Response(http.StatusOK, model.DiscoveryProviderDto{Uuid: discovery.UUID, Name: discovery.Name, Status: model.IN_PROGRESS, TotalCertificatesDiscovered: 0, CertificateData: nil, Meta: nil}), nil
	} else {
		pagination := db.Pagination{
			Page:  1,
			Limit: 10,
		}
		result, _ := s.discoveryRepo.List(pagination)
		var certificateDtos []model.DiscoveryProviderCertificateDataDto
		rows, _ := result.Rows.([]*db.Certificate) // Convert interface{} to []db.CertificateData
		for _, certificateData := range rows {
			discoveryProviderCertificateDataDto := model.DiscoveryProviderCertificateDataDto{
				Uuid:          certificateData.UUID,
				Base64Content: certificateData.Base64Content,
			}
			certificateDtos = append(certificateDtos, discoveryProviderCertificateDataDto)
		}

		return model.Response(http.StatusOK, model.DiscoveryProviderDto{Uuid: discovery.UUID, Name: discovery.Name, Status: model.COMPLETED, TotalCertificatesDiscovered: 0, CertificateData: certificateDtos, Meta: nil}), nil
	}

}

func (s *DiscoveryAPIService) DiscoveryCertificates(authority *db.AuthorityInstance, discovery *db.Discovery) {
	// get the vault client
	client, err := vault.GetClient(*authority)
	if err != nil {
		discovery.Status = "FAILED"
		s.discoveryRepo.UpdateDiscovery(discovery)
		s.log.Fatal(err.Error())
		return
	}
	// get the certificates
	ctx := context.Background()
	certificates, err := client.Secrets.PkiListCerts(ctx)
	if err != nil {
		discovery.Status = "FAILED"
		s.discoveryRepo.UpdateDiscovery(discovery)
		s.log.Fatal(err.Error())
		return
	}
	var certificateKeys []*db.Certificate
	for _, certificateKey := range certificates.Data.Keys {
		certificateData, err := client.Secrets.PkiReadCert(ctx, certificateKey)
		if err != nil {
			discovery.Status = "FAILED"
			s.discoveryRepo.UpdateDiscovery(discovery)
			s.log.Fatal(err.Error())
			return
		}
		certificate := db.Certificate{
			SerialNumber:  certificateKey,
			UUID:          utils.DeterministicGUID(certificateKey),
			Base64Content: certificateData.Data.Certificate,
		}
		certificateKeys = append(certificateKeys, &certificate)
	}
	s.discoveryRepo.AssociateCertificatesToDiscovery(discovery, certificateKeys...)

	// Update discovery status to "COMPLETED"
	discovery.Status = "COMPLETED"
	s.discoveryRepo.UpdateDiscovery(discovery)

}
