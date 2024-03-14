package model

type FileAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	Data FileAttributeContentData `json:"data"`
}

// AssertFileAttributeContentRequired checks if the required fields are not zero-ed
func AssertFileAttributeContentRequired(obj FileAttributeContent) error {
	elements := map[string]interface{}{
		"data": obj.Data,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertFileAttributeContentDataRequired(obj.Data); err != nil {
		return err
	}
	return nil
}

// AssertFileAttributeContentConstraints checks if the values respects the defined constraints
func AssertFileAttributeContentConstraints(obj FileAttributeContent) error {
	return nil
}
