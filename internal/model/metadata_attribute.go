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
	Content []AttributeContent `json:"content"`

	Type AttributeType `json:"type"`

	ContentType AttributeContentType `json:"contentType"`

	Properties MetadataAttributeProperties `json:"properties"`
}

func (d MetadataAttribute) GetContent() []AttributeContent {
	return d.Content
}
func (d MetadataAttribute) GetUuid() string {
	return d.Uuid
}

func (d MetadataAttribute) GetName() string {
	return d.Name
}

func (d MetadataAttribute) GetAttributeType() AttributeType {
	return d.Type
}

func (d MetadataAttribute) GetAttributeContentType() AttributeContentType {
	return d.ContentType
}

// AssertMetadataAttributeRequired checks if the required fields are not zero-ed
func AssertMetadataAttributeRequired(obj Attribute) error {
	objType := obj.(MetadataAttribute)
	elements := map[string]interface{}{
		"uuid":        objType.Uuid,
		"name":        objType.Name,
		"content":     objType.Content,
		"type":        objType.Type,
		"contentType": objType.ContentType,
		"properties":  objType.Properties,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range objType.Content {
		if err := AssertBaseAttributeContentDtoRequired(el); err != nil {
			return err
		}
	}
	if err := AssertMetadataAttributePropertiesRequired(objType.Properties); err != nil {
		return err
	}
	return nil
}

// AssertMetadataAttributeConstraints checks if the values respects the defined constraints
func AssertMetadataAttributeConstraints(obj MetadataAttribute) error {
	return nil
}
