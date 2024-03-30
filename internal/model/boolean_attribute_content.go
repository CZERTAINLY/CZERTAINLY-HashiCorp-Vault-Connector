package model

type BooleanAttributeContent struct {
	AttributeContent `json:"-"`
	// Content Reference
	Reference string `json:"reference,omitempty"`

	// Boolean attribute value
	Data bool `json:"data"`
}

// AssertBooleanAttributeContentRequired checks if the required fields are not zero-ed
func AssertBooleanAttributeContentRequired(obj BooleanAttributeContent) error {
	elements := map[string]interface{}{
		"data": obj.Data,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertBooleanAttributeContentConstraints checks if the values respects the defined constraints
func AssertBooleanAttributeContentConstraints(obj BooleanAttributeContent) error {
	return nil
}
