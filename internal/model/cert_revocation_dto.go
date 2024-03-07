package model

type CertRevocationDto struct {
	Reason CertificateRevocationReason `json:"reason"`

	// List of RA Profiles attributes
	RaProfileAttributes []RequestAttributeDto `json:"raProfileAttributes"`

	// List of Attributes to revoke Certificate
	Attributes []RequestAttributeDto `json:"attributes"`

	// Base64 Certificate content. (Certificate to be revoked)
	Certificate string `json:"certificate"`
}

// AssertCertRevocationDtoRequired checks if the required fields are not zero-ed
func AssertCertRevocationDtoRequired(obj CertRevocationDto) error {
	elements := map[string]interface{}{
		"reason":              obj.Reason,
		"raProfileAttributes": obj.RaProfileAttributes,
		"attributes":          obj.Attributes,
		"certificate":         obj.Certificate,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.RaProfileAttributes {
		if err := AssertRequestAttributeDtoRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Attributes {
		if err := AssertRequestAttributeDtoRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertCertRevocationDtoConstraints checks if the values respects the defined constraints
func AssertCertRevocationDtoConstraints(obj CertRevocationDto) error {
	return nil
}
