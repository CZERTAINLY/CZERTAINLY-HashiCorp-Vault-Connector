package model

type ObjectAttributeContent struct {
	// Content Reference
	Reference string `json:"reference,omitempty"`

	// Object attribute content data
	Data map[string]interface{} `json:"data"`
}

func (o ObjectAttributeContent) GetData() interface{} {
	return o.Data
}

func (o ObjectAttributeContent) GetReference() string {
	return o.Reference
}

// AssertObjectAttributeContentRequired checks if the required fields are not zero-ed
func AssertObjectAttributeContentRequired(obj ObjectAttributeContent) error {
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

// AssertObjectAttributeContentConstraints checks if the values respects the defined constraints
func AssertObjectAttributeContentConstraints(obj ObjectAttributeContent) error {
	return nil
}
