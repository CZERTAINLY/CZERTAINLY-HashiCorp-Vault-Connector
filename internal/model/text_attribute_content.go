package model

type TextAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	// Text attribute value
	Data string `json:"data"`
}

func (a TextAttributeContent) GetData() interface{} {
	return a.Data
}

func (a TextAttributeContent) GetReference() string {
	return a.Reference
}

// AssertTextAttributeContentRequired checks if the required fields are not zero-ed
func AssertTextAttributeContentRequired(obj TextAttributeContent) error {
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

// AssertTextAttributeContentConstraints checks if the values respects the defined constraints
func AssertTextAttributeContentConstraints(obj TextAttributeContent) error {
	return nil
}
