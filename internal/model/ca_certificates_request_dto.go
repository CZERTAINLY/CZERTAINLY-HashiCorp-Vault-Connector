package model

type CaCertificatesRequestDto struct {

	// List of RA Profiles attributes
	RaProfileAttributes []RequestAttributeDto `json:"raProfileAttributes"`
}

// AssertCaCertificatesRequestDtoRequired checks if the required fields are not zero-ed
func AssertCaCertificatesRequestDtoRequired(obj CaCertificatesRequestDto) error {
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

// AssertCaCertificatesRequestDtoConstraints checks if the values respects the defined constraints
func AssertCaCertificatesRequestDtoConstraints(obj CaCertificatesRequestDto) error {
	return nil
}
