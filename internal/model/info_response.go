package model

type InfoResponse struct {
	FunctionGroupCode FunctionGroupCode `json:"functionGroupCode"`

	// List of supported functional group kinds
	Kinds []string `json:"kinds"`

	// List of end points related to functional group
	EndPoints []EndpointDto `json:"endPoints"`
}

// AssertInfoResponseRequired checks if the required fields are not zero-ed
func AssertInfoResponseRequired(obj InfoResponse) error {
	elements := map[string]interface{}{
		"functionGroupCode": obj.FunctionGroupCode,
		"kinds":             obj.Kinds,
		"endPoints":         obj.EndPoints,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.EndPoints {
		if err := AssertEndpointDtoRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertInfoResponseConstraints checks if the values respects the defined constraints
func AssertInfoResponseConstraints(obj InfoResponse) error {
	return nil
}
