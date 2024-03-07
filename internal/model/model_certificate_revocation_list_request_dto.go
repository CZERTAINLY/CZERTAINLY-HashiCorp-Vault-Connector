package model

type CertificateRevocationListRequestDto struct {

	// If true, the delta CRL is returned, otherwise the full CRL is returned
	Delta bool `json:"delta,omitempty"`

	// List of RA Profiles attributes
	RaProfileAttributes []RequestAttributeDto `json:"raProfileAttributes"`
}

// AssertCertificateRevocationListRequestDtoRequired checks if the required fields are not zero-ed
func AssertCertificateRevocationListRequestDtoRequired(obj CertificateRevocationListRequestDto) error {
	elements := map[string]interface{}{
		"raProfileAttributes": obj.RaProfileAttributes,
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
	return nil
}

// AssertCertificateRevocationListRequestDtoConstraints checks if the values respects the defined constraints
func AssertCertificateRevocationListRequestDtoConstraints(obj CertificateRevocationListRequestDto) error {
	return nil
}
