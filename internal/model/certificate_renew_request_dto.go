package model

type CertificateRenewRequestDto struct {

	// Certificate sign request (PKCS#10) encoded as Base64 string
	Pkcs10 string `json:"pkcs10"`

	// List of RA Profiles attributes
	RaProfileAttributes []RequestAttributeDto `json:"raProfileAttributes"`

	// Base64 Certificate content. (Certificate to be renewed)
	Certificate string `json:"certificate"`

	// Metadata for the Certificate
	Meta []MetadataAttribute `json:"meta"`
}

// AssertCertificateRenewRequestDtoRequired checks if the required fields are not zero-ed
func AssertCertificateRenewRequestDtoRequired(obj CertificateRenewRequestDto) error {
	elements := map[string]interface{}{
		"pkcs10":              obj.Pkcs10,
		"raProfileAttributes": obj.RaProfileAttributes,
		"certificate":         obj.Certificate,
		"meta":                obj.Meta,
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
	for _, el := range obj.Meta {
		if err := AssertMetadataAttributeRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertCertificateRenewRequestDtoConstraints checks if the values respects the defined constraints
func AssertCertificateRenewRequestDtoConstraints(obj CertificateRenewRequestDto) error {
	return nil
}
