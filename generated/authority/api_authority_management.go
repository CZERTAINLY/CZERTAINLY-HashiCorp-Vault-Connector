/*
 * CZERTAINLY Authority Provider v2 API
 *
 * REST API for implementations of custom v2 Authority Provider
 *
 * API version: 2.11.0
 * Contact: getinfo@czertainly.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package authority

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// AuthorityManagementAPIController binds http requests to an api service and writes the service results to the http response
type AuthorityManagementAPIController struct {
	service      AuthorityManagementAPIServicer
	errorHandler ErrorHandler
}

// AuthorityManagementAPIOption for how the controller is set up.
type AuthorityManagementAPIOption func(*AuthorityManagementAPIController)

// WithAuthorityManagementAPIErrorHandler inject ErrorHandler into controller
func WithAuthorityManagementAPIErrorHandler(h ErrorHandler) AuthorityManagementAPIOption {
	return func(c *AuthorityManagementAPIController) {
		c.errorHandler = h
	}
}

// NewAuthorityManagementAPIController creates a default api controller
func NewAuthorityManagementAPIController(s AuthorityManagementAPIServicer, opts ...AuthorityManagementAPIOption) Router {
	controller := &AuthorityManagementAPIController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the AuthorityManagementAPIController
func (c *AuthorityManagementAPIController) Routes() Routes {
	return Routes{
		"CreateAuthorityInstance": Route{
			strings.ToUpper("Post"),
			"/v1/authorityProvider/authorities",
			c.CreateAuthorityInstance,
		},
		"GetAuthorityInstance": Route{
			strings.ToUpper("Get"),
			"/v1/authorityProvider/authorities/{uuid}",
			c.GetAuthorityInstance,
		},
		"GetCaCertificates": Route{
			strings.ToUpper("Post"),
			"/v1/authorityProvider/authorities/{uuid}/caCertificates",
			c.GetCaCertificates,
		},
		"GetConnection": Route{
			strings.ToUpper("Get"),
			"/v1/authorityProvider/authorities/{uuid}/connect",
			c.GetConnection,
		},
		"GetCrl": Route{
			strings.ToUpper("Post"),
			"/v1/authorityProvider/authorities/{uuid}/crl",
			c.GetCrl,
		},
		"ListAuthorityInstances": Route{
			strings.ToUpper("Get"),
			"/v1/authorityProvider/authorities",
			c.ListAuthorityInstances,
		},
		"ListRAProfileAttributes": Route{
			strings.ToUpper("Get"),
			"/v1/authorityProvider/authorities/{uuid}/raProfile/attributes",
			c.ListRAProfileAttributes,
		},
		"RemoveAuthorityInstance": Route{
			strings.ToUpper("Delete"),
			"/v1/authorityProvider/authorities/{uuid}",
			c.RemoveAuthorityInstance,
		},
		"UpdateAuthorityInstance": Route{
			strings.ToUpper("Post"),
			"/v1/authorityProvider/authorities/{uuid}",
			c.UpdateAuthorityInstance,
		},
		"ValidateRAProfileAttributes": Route{
			strings.ToUpper("Post"),
			"/v1/authorityProvider/authorities/{uuid}/raProfile/attributes/validate",
			c.ValidateRAProfileAttributes,
		},
	}
}

// CreateAuthorityInstance - Create Authority instance
func (c *AuthorityManagementAPIController) CreateAuthorityInstance(w http.ResponseWriter, r *http.Request) {
	authorityProviderInstanceRequestDtoParam := AuthorityProviderInstanceRequestDto{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&authorityProviderInstanceRequestDtoParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertAuthorityProviderInstanceRequestDtoRequired(authorityProviderInstanceRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertAuthorityProviderInstanceRequestDtoConstraints(authorityProviderInstanceRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateAuthorityInstance(r.Context(), authorityProviderInstanceRequestDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetAuthorityInstance - Get an Authority instance
func (c *AuthorityManagementAPIController) GetAuthorityInstance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &RequiredError{"uuid"}, nil)
		return
	}
	result, err := c.service.GetAuthorityInstance(r.Context(), uuidParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetCaCertificates - Get the Authority Instance's certificate chain
func (c *AuthorityManagementAPIController) GetCaCertificates(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &RequiredError{"uuid"}, nil)
		return
	}
	caCertificatesRequestDtoParam := CaCertificatesRequestDto{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&caCertificatesRequestDtoParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertCaCertificatesRequestDtoRequired(caCertificatesRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertCaCertificatesRequestDtoConstraints(caCertificatesRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.GetCaCertificates(r.Context(), uuidParam, caCertificatesRequestDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetConnection - Connect to Authority
func (c *AuthorityManagementAPIController) GetConnection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &RequiredError{"uuid"}, nil)
		return
	}
	result, err := c.service.GetConnection(r.Context(), uuidParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetCrl - Get the latest CRL for the Authority Instance
func (c *AuthorityManagementAPIController) GetCrl(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &RequiredError{"uuid"}, nil)
		return
	}
	certificateRevocationListRequestDtoParam := CertificateRevocationListRequestDto{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&certificateRevocationListRequestDtoParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertCertificateRevocationListRequestDtoRequired(certificateRevocationListRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertCertificateRevocationListRequestDtoConstraints(certificateRevocationListRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.GetCrl(r.Context(), uuidParam, certificateRevocationListRequestDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ListAuthorityInstances - List Authority instances
func (c *AuthorityManagementAPIController) ListAuthorityInstances(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.ListAuthorityInstances(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ListRAProfileAttributes - List RA Profile Attributes
func (c *AuthorityManagementAPIController) ListRAProfileAttributes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &RequiredError{"uuid"}, nil)
		return
	}
	result, err := c.service.ListRAProfileAttributes(r.Context(), uuidParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// RemoveAuthorityInstance - Remove Authority instance
func (c *AuthorityManagementAPIController) RemoveAuthorityInstance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &RequiredError{"uuid"}, nil)
		return
	}
	result, err := c.service.RemoveAuthorityInstance(r.Context(), uuidParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateAuthorityInstance - Update Authority instance
func (c *AuthorityManagementAPIController) UpdateAuthorityInstance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &RequiredError{"uuid"}, nil)
		return
	}
	authorityProviderInstanceRequestDtoParam := AuthorityProviderInstanceRequestDto{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&authorityProviderInstanceRequestDtoParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertAuthorityProviderInstanceRequestDtoRequired(authorityProviderInstanceRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := AssertAuthorityProviderInstanceRequestDtoConstraints(authorityProviderInstanceRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateAuthorityInstance(r.Context(), uuidParam, authorityProviderInstanceRequestDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// ValidateRAProfileAttributes - Validate RA Profile attributes
func (c *AuthorityManagementAPIController) ValidateRAProfileAttributes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &RequiredError{"uuid"}, nil)
		return
	}
	requestAttributeDtoParam := []RequestAttributeDto{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&requestAttributeDtoParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	for _, el := range requestAttributeDtoParam {
		if err := AssertRequestAttributeDtoRequired(el); err != nil {
			c.errorHandler(w, r, err, nil)
			return
		}
	}
	result, err := c.service.ValidateRAProfileAttributes(r.Context(), uuidParam, requestAttributeDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
