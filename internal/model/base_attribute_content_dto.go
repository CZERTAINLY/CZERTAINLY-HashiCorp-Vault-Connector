package model

// BaseAttributeContentDto - Base Attribute content definition
type BaseAttributeContentDto struct {
	AttributeContent
	// Content Reference
	Reference string `json:"reference,omitempty"`

	// Content Data
	Data map[string]interface{} `json:"data"`
}

// AssertBaseAttributeContentDtoRequired checks if the required fields are not zero-ed
func AssertBaseAttributeContentDtoRequired(obj AttributeContent) error {
	elements := map[string]interface{}{
		"data": obj.GetData(),
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertBaseAttributeContentDtoConstraints checks if the values respects the defined constraints
func AssertBaseAttributeContentDtoConstraints(obj BaseAttributeContentDto) error {
	return nil
}
