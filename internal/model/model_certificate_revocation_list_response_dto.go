package model

type CertificateRevocationListResponseDto struct {

	// Base64 encoded CRL data
	CrlData []string `json:"crlData"`
}

// AssertCertificateRevocationListResponseDtoRequired checks if the required fields are not zero-ed
func AssertCertificateRevocationListResponseDtoRequired(obj CertificateRevocationListResponseDto) error {
	elements := map[string]interface{}{
		"crlData": obj.CrlData,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertCertificateRevocationListResponseDtoConstraints checks if the values respects the defined constraints
func AssertCertificateRevocationListResponseDtoConstraints(obj CertificateRevocationListResponseDto) error {
	return nil
}
