package model

type IntegerAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	// Integer attribute value
	Data int32 `json:"data"`
}

// AssertIntegerAttributeContentRequired checks if the required fields are not zero-ed
func AssertIntegerAttributeContentRequired(obj IntegerAttributeContent) error {
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

// AssertIntegerAttributeContentConstraints checks if the values respects the defined constraints
func AssertIntegerAttributeContentConstraints(obj IntegerAttributeContent) error {
	return nil
}
