package health

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/cmd/logger"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) GetHealth() (GetHealthResponse, error) {
	l := logger.Get()
	l.Info("Entering GetHealth method")

	return GetHealthResponse{
		Status:      "ok",
		Description: "OK",
	}, nil
}
