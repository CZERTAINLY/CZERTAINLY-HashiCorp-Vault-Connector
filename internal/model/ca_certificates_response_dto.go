package model

type CaCertificatesResponseDto struct {

	// List of Certificates
	Certificates []CertificateDataResponseDto `json:"certificates"`
}

// AssertCaCertificatesResponseDtoRequired checks if the required fields are not zero-ed
func AssertCaCertificatesResponseDtoRequired(obj CaCertificatesResponseDto) error {
	elements := map[string]interface{}{
		"certificates": obj.Certificates,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Certificates {
		if err := AssertCertificateDataResponseDtoRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertCaCertificatesResponseDtoConstraints checks if the values respects the defined constraints
func AssertCaCertificatesResponseDtoConstraints(obj CaCertificatesResponseDto) error {
	return nil
}
