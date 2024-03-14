package model

type CodeBlockAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	Data CodeBlockAttributeContentData `json:"data"`
}

// AssertCodeBlockAttributeContentRequired checks if the required fields are not zero-ed
func AssertCodeBlockAttributeContentRequired(obj CodeBlockAttributeContent) error {
	elements := map[string]interface{}{
		"data": obj.Data,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertCodeBlockAttributeContentDataRequired(obj.Data); err != nil {
		return err
	}
	return nil
}

// AssertCodeBlockAttributeContentConstraints checks if the values respects the defined constraints
func AssertCodeBlockAttributeContentConstraints(obj CodeBlockAttributeContent) error {
	return nil
}
