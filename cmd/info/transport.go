package info

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type Router interface {
	Handle(method, path string, handler http.Handler)
}

func RegisterRoutes(router *mux.Router, s Service) {
	getInfoHandler := httptransport.NewServer(
		MakeGetInfoEndpoint(s, router), decodeGetInfoRequest, encodeGetInfoResponse)

	router.Methods(http.MethodGet).Path("/v1").Handler(getInfoHandler).Name("listSupportedFunctions")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// get info request
////////////////////////////////////////////////////////////////////////////////////////////////////////

func encodeGetInfoResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(writer).Encode(response)
}

func decodeGetInfoRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return request, nil
}
