package health

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

type Service interface {
	// GetHealth returns the health of the connector. Since the compliance provider connector
	// does not have any real checks, it returns a healthy status.
	GetHealth() (GetHealthResponse, error)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// get health request
////////////////////////////////////////////////////////////////////////////////////////////////////////

type (
	GetHealthResponse struct {
		Status      string `json:"status"`
		Description string `json:"description"`
	}
)

func MakeGetHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		healthStatus, err := s.GetHealth()
		if err != nil {
			return nil, fmt.Errorf("MakeGetHealthEndpoint: %w", err)
		}

		return healthStatus, nil
	}
}
