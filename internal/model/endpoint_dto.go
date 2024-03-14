package model

// EndpointDto - List of end points related to functional group
type EndpointDto struct {

	// Object identifier
	Uuid string `json:"uuid"`

	// Object Name
	Name string `json:"name"`

	// Context of the Endpoint
	Context string `json:"context"`

	// Method to be used for the Endpoint
	Method string `json:"method"`

	// True if the Endpoint is required for implementation
	Required bool `json:"required"`
}

// AssertEndpointDtoRequired checks if the required fields are not zero-ed
func AssertEndpointDtoRequired(obj EndpointDto) error {
	elements := map[string]interface{}{
		"uuid":     obj.Uuid,
		"name":     obj.Name,
		"context":  obj.Context,
		"method":   obj.Method,
		"required": obj.Required,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertEndpointDtoConstraints checks if the values respects the defined constraints
func AssertEndpointDtoConstraints(obj EndpointDto) error {
	return nil
}
