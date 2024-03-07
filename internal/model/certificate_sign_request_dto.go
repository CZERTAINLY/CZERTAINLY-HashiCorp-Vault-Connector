package model

type CertificateSignRequestDto struct {

	// Certificate sign request (PKCS#10) encoded as Base64 string
	Pkcs10 string `json:"pkcs10"`

	// List of RA Profiles attributes
	RaProfileAttributes []RequestAttributeDto `json:"raProfileAttributes"`

	// List of Attributes to issue Certificate
	Attributes []RequestAttributeDto `json:"attributes"`
}

// AssertCertificateSignRequestDtoRequired checks if the required fields are not zero-ed
func AssertCertificateSignRequestDtoRequired(obj CertificateSignRequestDto) error {
	elements := map[string]interface{}{
		"pkcs10":              obj.Pkcs10,
		"raProfileAttributes": obj.RaProfileAttributes,
		"attributes":          obj.Attributes,
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

// AssertCertificateSignRequestDtoConstraints checks if the values respects the defined constraints
func AssertCertificateSignRequestDtoConstraints(obj CertificateSignRequestDto) error {
	return nil
}
