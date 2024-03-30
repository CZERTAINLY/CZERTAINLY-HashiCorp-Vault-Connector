package model

type FloatAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	// Float attribute value
	Data float32 `json:"data"`
}

// AssertFloatAttributeContentRequired checks if the required fields are not zero-ed
func AssertFloatAttributeContentRequired(obj FloatAttributeContent) error {
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

// AssertFloatAttributeContentConstraints checks if the values respects the defined constraints
func AssertFloatAttributeContentConstraints(obj FloatAttributeContent) error {
	return nil
}
