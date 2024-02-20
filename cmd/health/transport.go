package health

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
		MakeGetHealthEndpoint(s), decodeGetInfoRequest, encodeGetInfoResponse)

	router.Methods(http.MethodGet).Path("/v1/health").Handler(getInfoHandler).Name("getHealth")
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
