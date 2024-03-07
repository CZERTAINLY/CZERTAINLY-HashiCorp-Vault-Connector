package connectorInfo

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"net/http"
	"strings"
)

// ConnectorInfoAPIController binds http requests to an api service and writes the service results to the http response
type ConnectorInfoAPIController struct {
	service      ConnectorInfoAPIServicer
	errorHandler model.ErrorHandler
}

// ConnectorInfoAPIOption for how the controller is set up.
type ConnectorInfoAPIOption func(*ConnectorInfoAPIController)

// WithConnectorInfoAPIErrorHandler inject model.ErrorHandler into controller
func WithConnectorInfoAPIErrorHandler(h model.ErrorHandler) ConnectorInfoAPIOption {
	return func(c *ConnectorInfoAPIController) {
		c.errorHandler = h
	}
}

// NewConnectorInfoAPIController creates a default api controller
func NewConnectorInfoAPIController(s ConnectorInfoAPIServicer, opts ...ConnectorInfoAPIOption) model.Router {
	controller := &ConnectorInfoAPIController{
		service:      s,
		errorHandler: model.DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// model.Routes returns all the api routes for the ConnectorInfoAPIController
func (c *ConnectorInfoAPIController) Routes() model.Routes {
	return model.Routes{
		"ListSupportedFunctions": model.Route{
			strings.ToUpper("Get"),
			"/v1",
			c.ListSupportedFunctions,
		},
	}
}

// ListSupportedFunctions - List supported functions of the connector
func (c *ConnectorInfoAPIController) ListSupportedFunctions(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.ListSupportedFunctions(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}
