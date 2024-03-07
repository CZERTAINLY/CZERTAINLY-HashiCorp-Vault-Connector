package model

type StringAttributeContent struct {
	AttributeContent `json:"-"`
	// Content Reference
	Reference string `json:"reference,omitempty"`

	// String attribute value
	Data string `json:"data"`
}

// AssertStringAttributeContentRequired checks if the required fields are not zero-ed
func AssertStringAttributeContentRequired(obj StringAttributeContent) error {
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

// AssertStringAttributeContentConstraints checks if the values respects the defined constraints
func AssertStringAttributeContentConstraints(obj StringAttributeContent) error {
	return nil
}
