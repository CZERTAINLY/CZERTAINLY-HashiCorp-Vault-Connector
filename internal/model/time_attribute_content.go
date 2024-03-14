package model

type TimeAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	// Time attribute value in format HH:mm:ss
	Data string `json:"data"`
}

// AssertTimeAttributeContentRequired checks if the required fields are not zero-ed
func AssertTimeAttributeContentRequired(obj TimeAttributeContent) error {
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

// AssertTimeAttributeContentConstraints checks if the values respects the defined constraints
func AssertTimeAttributeContentConstraints(obj TimeAttributeContent) error {
	return nil
}
