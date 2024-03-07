package authority

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ConnectorAttributesAPIController binds http requests to an api service and writes the service results to the http response
type ConnectorAttributesAPIController struct {
	service      ConnectorAttributesAPIServicer
	errorHandler model.ErrorHandler
}

// ConnectorAttributesAPIOption for how the controller is set up.
type ConnectorAttributesAPIOption func(*ConnectorAttributesAPIController)

// WithConnectorAttributesAPIErrorHandler inject model.ErrorHandler into controller
func WithConnectorAttributesAPIErrorHandler(h model.ErrorHandler) ConnectorAttributesAPIOption {
	return func(c *ConnectorAttributesAPIController) {
		c.errorHandler = h
	}
}

// NewConnectorAttributesAPIController creates a default api controller
func NewConnectorAttributesAPIController(s ConnectorAttributesAPIServicer, opts ...ConnectorAttributesAPIOption) model.Router {
	controller := &ConnectorAttributesAPIController{
		service:      s,
		errorHandler: model.DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ConnectorAttributesAPIController
func (c *ConnectorAttributesAPIController) Routes() model.Routes {
	return model.Routes{
		"ListAttributeDefinitions": model.Route{
			strings.ToUpper("Get"),
			"/v1/authorityProvider/{kind}/attributes",
			c.ListAttributeDefinitions,
		},
		"ValidateAttributes": model.Route{
			strings.ToUpper("Post"),
			"/v1/authorityProvider/{kind}/attributes/validate",
			c.ValidateAttributes,
		},
	}
}

// ListAttributeDefinitions - List available Attributes
func (c *ConnectorAttributesAPIController) ListAttributeDefinitions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	kindParam := params["kind"]
	if kindParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"kind"}, nil)
		return
	}
	result, err := c.service.ListAttributeDefinitions(r.Context(), kindParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}

// ValidateAttributes - Validate Attributes
func (c *ConnectorAttributesAPIController) ValidateAttributes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	kindParam := params["kind"]
	if kindParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"kind"}, nil)
		return
	}
	requestAttributeDtoParam := []model.RequestAttributeDto{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestAttributeDtoParam); err != nil {
		c.errorHandler(w, r, &model.ParsingError{Err: err}, nil)
		return
	}
	for _, el := range requestAttributeDtoParam {
		if err := model.AssertRequestAttributeDtoRequired(el); err != nil {
			c.errorHandler(w, r, err, nil)
			return
		}
	}
	result, err := c.service.ValidateAttributes(r.Context(), kindParam, requestAttributeDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}
