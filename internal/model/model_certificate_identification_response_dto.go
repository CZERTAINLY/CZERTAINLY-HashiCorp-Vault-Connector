package model

type CertificateIdentificationResponseDto struct {

	// Metadata for identified certificate
	Meta []MetadataAttribute `json:"meta"`
}

// AssertCertificateIdentificationResponseDtoRequired checks if the required fields are not zero-ed
func AssertCertificateIdentificationResponseDtoRequired(obj CertificateIdentificationResponseDto) error {
	elements := map[string]interface{}{
		"meta": obj.Meta,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Meta {
		if err := AssertMetadataAttributeRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertCertificateIdentificationResponseDtoConstraints checks if the values respects the defined constraints
func AssertCertificateIdentificationResponseDtoConstraints(obj CertificateIdentificationResponseDto) error {
	return nil
}
