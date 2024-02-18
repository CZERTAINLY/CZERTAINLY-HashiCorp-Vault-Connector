package info

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/cmd/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var endpoints []EndPointsInfo

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) GetInfo(router *mux.Router) ([]GetInfoResponse, error) {
	l := logger.Get()
	l.Info("Entering GetInfo method")

	endpoints = endpoints[:0]
	err := router.Walk(s.gorillaWalkFn)
	if err != nil {
		return nil, err
	}

	var infoResponse []GetInfoResponse
	infoResponse = append(infoResponse, GetInfoResponse{
		FunctionGroupCode: "authorityProvider",
		Kinds:             []string{"x509"},
		EndPoints:         endpoints,
	})

	l.Info("List of endpoints for the connector is ", zap.Any("endpoints", infoResponse))

	return infoResponse, nil
}

func (s *service) gorillaWalkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	l := logger.Get()

	l.Info("Entering gorillaWalkFn to calculate the list of end points")

	path, _ := route.GetPathTemplate()
	l.Debug("Path: ", zap.String("path", path))

	method, _ := route.GetMethods()
	l.Debug("Method: ", zap.String("method", method[0]))

	name := route.GetName()
	l.Debug("Name: ", zap.String("name", name))

	endpoint := EndPointsInfo{
		Name:    name,
		Context: path,
		Method:  method[0],
	}
	l.Debug("Endpoint created: ", zap.Any("endpoint", endpoint))

	endpoints = append(endpoints, endpoint)
	return nil
}
