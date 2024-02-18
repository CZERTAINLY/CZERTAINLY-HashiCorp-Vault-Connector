package info

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

type Service interface {
	// GetInfo returns the info for the given request. This information includes the list of end points
	// and the list of attributes.
	GetInfo(router *mux.Router) ([]GetInfoResponse, error)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// get info request
////////////////////////////////////////////////////////////////////////////////////////////////////////

type (
	EndPointsInfo struct {
		Name    string `json:"name"`
		Context string `json:"context"`
		Method  string `json:"method"`
	}

	GetInfoResponse struct {
		FunctionGroupCode string          `json:"functionGroupCode"`
		Kinds             []string        `json:"kinds"`
		EndPoints         []EndPointsInfo `json:"endPoints"`
	}
)

func MakeGetInfoEndpoint(s Service, router *mux.Router) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		a, err := s.GetInfo(router)
		if err != nil {
			return nil, fmt.Errorf("MakeGetInfoEndpoint: %w", err)
		}

		return a, nil
	}
}
