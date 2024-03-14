package model

import (
	"fmt"
)

// CertificateRevocationReason : Revocation reason
type CertificateRevocationReason string

// List of CertificateRevocationReason
const (
	UNSPECIFIED            CertificateRevocationReason = "unspecified"
	KEY_COMPROMISE         CertificateRevocationReason = "keyCompromise"
	C_A_COMPROMISE         CertificateRevocationReason = "cACompromise"
	AFFILIATION_CHANGED    CertificateRevocationReason = "affiliationChanged"
	SUPERSEDED             CertificateRevocationReason = "superseded"
	CESSATION_OF_OPERATION CertificateRevocationReason = "cessationOfOperation"
	CERTIFICATE_HOLD       CertificateRevocationReason = "certificateHold"
	PRIVILEGE_WITHDRAWN    CertificateRevocationReason = "privilegeWithdrawn"
	A_A_COMPROMISE         CertificateRevocationReason = "aACompromise"
)

// AllowedCertificateRevocationReasonEnumValues is all the allowed values of CertificateRevocationReason enum
var AllowedCertificateRevocationReasonEnumValues = []CertificateRevocationReason{
	"unspecified",
	"keyCompromise",
	"cACompromise",
	"affiliationChanged",
	"superseded",
	"cessationOfOperation",
	"certificateHold",
	"privilegeWithdrawn",
	"aACompromise",
}

// validCertificateRevocationReasonEnumValue provides a map of CertificateRevocationReasons for fast verification of use input
var validCertificateRevocationReasonEnumValues = map[CertificateRevocationReason]struct{}{
	"unspecified":          {},
	"keyCompromise":        {},
	"cACompromise":         {},
	"affiliationChanged":   {},
	"superseded":           {},
	"cessationOfOperation": {},
	"certificateHold":      {},
	"privilegeWithdrawn":   {},
	"aACompromise":         {},
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v CertificateRevocationReason) IsValid() bool {
	_, ok := validCertificateRevocationReasonEnumValues[v]
	return ok
}

// NewCertificateRevocationReasonFromValue returns a pointer to a valid CertificateRevocationReason
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewCertificateRevocationReasonFromValue(v string) (CertificateRevocationReason, error) {
	ev := CertificateRevocationReason(v)
	if ev.IsValid() {
		return ev, nil
	} else {
		return "", fmt.Errorf("invalid value '%v' for CertificateRevocationReason: valid values are %v", v, AllowedCertificateRevocationReasonEnumValues)
	}
}

// AssertCertificateRevocationReasonRequired checks if the required fields are not zero-ed
func AssertCertificateRevocationReasonRequired(obj CertificateRevocationReason) error {
	return nil
}

// AssertCertificateRevocationReasonConstraints checks if the values respects the defined constraints
func AssertCertificateRevocationReasonConstraints(obj CertificateRevocationReason) error {
	return nil
}
