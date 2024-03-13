package authority

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// CertificateManagementAPIController binds http requests to an api service and writes the service results to the http response
type CertificateManagementAPIController struct {
	service      CertificateManagementAPIServicer
	errorHandler model.ErrorHandler
}

// CertificateManagementAPIOption for how the controller is set up.
type CertificateManagementAPIOption func(*CertificateManagementAPIController)

// WithCertificateManagementAPIErrorHandler inject model.ErrorHandler into controller
func WithCertificateManagementAPIErrorHandler(h model.ErrorHandler) CertificateManagementAPIOption {
	return func(c *CertificateManagementAPIController) {
		c.errorHandler = h
	}
}

// NewCertificateManagementAPIController creates a default api controller
func NewCertificateManagementAPIController(s CertificateManagementAPIServicer, opts ...CertificateManagementAPIOption) model.Router {
	controller := &CertificateManagementAPIController{
		service:      s,
		errorHandler: model.DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// model.Routes returns all the api routes for the CertificateManagementAPIController
func (c *CertificateManagementAPIController) Routes() model.Routes {
	return model.Routes{
		"IdentifyCertificate": model.Route{
			strings.ToUpper("Post"),
			"/v2/authorityProvider/authorities/{uuid}/certificates/identify",
			c.IdentifyCertificate,
		},
		"IssueCertificate": model.Route{
			strings.ToUpper("Post"),
			"/v2/authorityProvider/authorities/{uuid}/certificates/issue",
			c.IssueCertificate,
		},
		"ListIssueCertificateAttributes": model.Route{
			strings.ToUpper("Get"),
			"/v2/authorityProvider/authorities/{uuid}/certificates/issue/attributes",
			c.ListIssueCertificateAttributes,
		},
		"ListRevokeCertificateAttributes": model.Route{
			strings.ToUpper("Get"),
			"/v2/authorityProvider/authorities/{uuid}/certificates/revoke/attributes",
			c.ListRevokeCertificateAttributes,
		},
		"RenewCertificate": model.Route{
			strings.ToUpper("Post"),
			"/v2/authorityProvider/authorities/{uuid}/certificates/renew",
			c.RenewCertificate,
		},
		"RevokeCertificate": model.Route{
			strings.ToUpper("Post"),
			"/v2/authorityProvider/authorities/{uuid}/certificates/revoke",
			c.RevokeCertificate,
		},
		"ValidateIssueCertificateAttributes": model.Route{
			strings.ToUpper("Post"),
			"/v2/authorityProvider/authorities/{uuid}/certificates/issue/attributes/validate",
			c.ValidateIssueCertificateAttributes,
		},
		"ValidateRevokeCertificateAttributes": model.Route{
			strings.ToUpper("Post"),
			"/v2/authorityProvider/authorities/{uuid}/certificates/revoke/attributes/validate",
			c.ValidateRevokeCertificateAttributes,
		},
	}
}

// IdentifyCertificate - Identify Certificate
func (c *CertificateManagementAPIController) IdentifyCertificate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"uuid"}, nil)
		return
	}
	certificateIdentificationRequestDtoParam := model.CertificateIdentificationRequestDto{}
	json, err := io.ReadAll(r.Body)
	if err != nil {
		c.errorHandler(w, r, &model.ParsingError{Err: err}, nil)
		return
	}

	certificateIdentificationRequestDtoParam.Unmarshal(json)
	if err := model.AssertCertificateIdentificationRequestDtoRequired(certificateIdentificationRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := model.AssertCertificateIdentificationRequestDtoConstraints(certificateIdentificationRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.IdentifyCertificate(r.Context(), uuidParam, certificateIdentificationRequestDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}

// IssueCertificate - Issue Certificate
func (c *CertificateManagementAPIController) IssueCertificate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"uuid"}, nil)
		return
	}
	certificateSignRequestDtoParam := model.CertificateSignRequestDto{}
	json, err := io.ReadAll(r.Body)
	if err != nil {
		c.errorHandler(w, r, &model.ParsingError{Err: err}, nil)
		return
	}

	certificateSignRequestDtoParam.Unmarshal(json)
	if err := model.AssertCertificateSignRequestDtoRequired(certificateSignRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := model.AssertCertificateSignRequestDtoConstraints(certificateSignRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.IssueCertificate(r.Context(), uuidParam, certificateSignRequestDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}

// ListIssueCertificateAttributes - List of Attributes to issue Certificate
func (c *CertificateManagementAPIController) ListIssueCertificateAttributes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"uuid"}, nil)
		return
	}
	result, err := c.service.ListIssueCertificateAttributes(r.Context(), uuidParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}

// ListRevokeCertificateAttributes - List of Attributes to revoke Certificate
func (c *CertificateManagementAPIController) ListRevokeCertificateAttributes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"uuid"}, nil)
		return
	}
	result, err := c.service.ListRevokeCertificateAttributes(r.Context(), uuidParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}

// RenewCertificate - Renew Certificate
func (c *CertificateManagementAPIController) RenewCertificate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"uuid"}, nil)
		return
	}
	certificateRenewRequestDtoParam := model.CertificateRenewRequestDto{}
	json, err := io.ReadAll(r.Body)
	if err != nil {
		c.errorHandler(w, r, &model.ParsingError{Err: err}, nil)
		return
	}

	certificateRenewRequestDtoParam.Unmarshal(json)
	if err := model.AssertCertificateRenewRequestDtoRequired(certificateRenewRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := model.AssertCertificateRenewRequestDtoConstraints(certificateRenewRequestDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.RenewCertificate(r.Context(), uuidParam, certificateRenewRequestDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}

// RevokeCertificate - Revoke Certificate
func (c *CertificateManagementAPIController) RevokeCertificate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"uuid"}, nil)
		return
	}
	certRevocationDtoParam := model.CertRevocationDto{}
	json, err := io.ReadAll(r.Body)
	if err != nil {
		c.errorHandler(w, r, &model.ParsingError{Err: err}, nil)
		return
	}

	certRevocationDtoParam.Unmarshal(json)
	if err := model.AssertCertRevocationDtoRequired(certRevocationDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	if err := model.AssertCertRevocationDtoConstraints(certRevocationDtoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.RevokeCertificate(r.Context(), uuidParam, certRevocationDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}

// ValidateIssueCertificateAttributes - Validate list of Attributes to issue Certificate
func (c *CertificateManagementAPIController) ValidateIssueCertificateAttributes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"uuid"}, nil)
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
	result, err := c.service.ValidateIssueCertificateAttributes(r.Context(), uuidParam, requestAttributeDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}

// ValidateRevokeCertificateAttributes - Validate list of Attributes to revoke certificate
func (c *CertificateManagementAPIController) ValidateRevokeCertificateAttributes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := params["uuid"]
	if uuidParam == "" {
		c.errorHandler(w, r, &model.RequiredError{"uuid"}, nil)
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
	result, err := c.service.ValidateRevokeCertificateAttributes(r.Context(), uuidParam, requestAttributeDtoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	model.EncodeJSONResponse(result.Body, &result.Code, w)
}
