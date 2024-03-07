package model

// RequestAttributeDto - Request attribute to send attribute content for object
type RequestAttributeDto struct {
	Attribute
	// // UUID of the Attribute
	// Uuid string `json:"uuid,omitempty"`

	// // Name of the Attribute
	// Name string `json:"name"`

	// // Content of the Attribute
	// Content []BaseAttributeContentDto `json:"content"`
}

// AssertRequestAttributeDtoRequired checks if the required fields are not zero-ed
func AssertRequestAttributeDtoRequired(obj RequestAttributeDto) error {
	elements := map[string]interface{}{
		"name":    obj.GetName(),
		"content": obj.GetContent(),
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	// for _, el := range obj.GetContent() {
	// 	if err := AssertBaseAttributeContentDtoRequired(el); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// AssertRequestAttributeDtoConstraints checks if the values respects the defined constraints
func AssertRequestAttributeDtoConstraints(obj RequestAttributeDto) error {
	return nil
}
