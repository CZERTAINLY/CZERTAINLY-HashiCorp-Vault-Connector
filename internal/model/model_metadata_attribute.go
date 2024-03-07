package model

// MetadataAttribute - Info attribute contains content that is for metadata. Its content can not be edited and is not send in requests to store.
type MetadataAttribute struct {

	// UUID of the Attribute for unique identification
	Uuid string `json:"uuid"`

	// Name of the Attribute that is used for identification
	Name string `json:"name"`

	// Optional description of the Attribute, should contain helper text on what is expected
	Description string `json:"description,omitempty"`

	// Content of the Attribute
	Content []BaseAttributeContentDto `json:"content"`

	Type AttributeType `json:"type"`

	ContentType AttributeContentType `json:"contentType"`

	Properties MetadataAttributeProperties `json:"properties"`
}

// AssertMetadataAttributeRequired checks if the required fields are not zero-ed
func AssertMetadataAttributeRequired(obj MetadataAttribute) error {
	elements := map[string]interface{}{
		"uuid":        obj.Uuid,
		"name":        obj.Name,
		"content":     obj.Content,
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
	if err := AssertMetadataAttributePropertiesRequired(obj.Properties); err != nil {
		return err
	}
	return nil
}

// AssertMetadataAttributeConstraints checks if the values respects the defined constraints
func AssertMetadataAttributeConstraints(obj MetadataAttribute) error {
	return nil
}
