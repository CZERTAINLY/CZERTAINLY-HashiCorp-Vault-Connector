package model

import (
	"fmt"
)

// CertificateRequestFormat : Format of the Certificate Request
type CertificateRequestFormat string

// List of CertificateRequestFormat
const (
	CERTIFICATEREQUESTFORMAT_PKCS10 CertificateRequestFormat = "pkcs10"
	CERTIFICATEREQUESTFORMAT_CRMF   CertificateRequestFormat = "crmf"
)

// AllowedCertificateRequestFormatEnumValues is all the allowed values of CertificateRequestFormat enum
var AllowedCertificateRequestFormatEnumValues = []CertificateRequestFormat{
	"pkcs101",
	"crmf",
}

// validCertificateRequestFormatEnumValue provides a map of CertificateRequestFormat for fast verification of use input
var validCertificateRequestFormatEnumValues = map[CertificateRequestFormat]struct{}{
	"pkcs10": {},
	"crmf":   {},
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v CertificateRequestFormat) IsValid() bool {
	_, ok := validCertificateRequestFormatEnumValues[v]
	return ok
}

// NewCertificateRequestFormatFromValue returns a pointer to a valid CertificateRequestFormat
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewCertificateRequestFormatFromValue(v string) (CertificateRequestFormat, error) {
	ev := CertificateRequestFormat(v)
	if ev.IsValid() {
		return ev, nil
	} else {
		return "", fmt.Errorf("invalid value '%v' for CertificateRequestFormat: valid values are %v", v, AllowedCertificateRequestFormatEnumValues)
	}
}

// AssertCertificateRequestFormatRequired checks if the required fields are not zero-ed
func AssertCertificateRequestFormatRequired(obj CertificateRequestFormat) error {
	return nil
}

// AssertCertificateRequestFormatConstraints checks if the values respects the defined constraints
func AssertCertificateRequestFormatConstraints(obj CertificateType) error {
	return nil
}
