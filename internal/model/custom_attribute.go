package model

// CustomAttribute - Custom attribute allows to store and transfer dynamic data. Its content can be edited and send in requests to store.
type CustomAttribute struct {

	// UUID of the Attribute for unique identification
	Uuid string `json:"uuid"`

	// Name of the Attribute that is used for identification
	Name string `json:"name"`

	// Optional description of the Attribute, should contain helper text on what is expected
	Description string `json:"description,omitempty"`

	// Content of the Attribute
	Content []BaseAttributeContentDto `json:"content,omitempty"`

	Type AttributeType `json:"type"`

	ContentType AttributeContentType `json:"contentType"`

	Properties CustomAttributeProperties `json:"properties"`
}

// AssertCustomAttributeRequired checks if the required fields are not zero-ed
func AssertCustomAttributeRequired(obj CustomAttribute) error {
	elements := map[string]interface{}{
		"uuid":        obj.Uuid,
		"name":        obj.Name,
		"type":        obj.Type,
		"contentType": obj.ContentType,
		"properties":  obj.Properties,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Content {
		if err := AssertBaseAttributeContentDtoRequired(el); err != nil {
			return err
		}
	}
	if err := AssertCustomAttributePropertiesRequired(obj.Properties); err != nil {
		return err
	}
	return nil
}

// AssertCustomAttributeConstraints checks if the values respects the defined constraints
func AssertCustomAttributeConstraints(obj CustomAttribute) error {
	return nil
}
