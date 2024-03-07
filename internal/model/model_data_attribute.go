package model

// DataAttribute - Data attribute allows to store and transfer dynamic data. Its content can be edited and send in requests to store.
type DataAttribute struct {
	Attribute `json:"-"`
	// UUID of the Attribute for unique identification
	Uuid string `json:"uuid"`

	// Name of the Attribute that is used for identification
	Name string `json:"name"`

	// Optional description of the Attribute, should contain helper text on what is expected
	Description string `json:"description,omitempty"`

	// Content of the Attribute
	Content []AttributeContent `json:"content,omitempty"`

	Type AttributeType `json:"type"`

	ContentType AttributeContentType `json:"contentType"`

	Properties *DataAttributeProperties `json:"properties"`

	// Optional regular expressions and constraints used for validating the Attribute content
	Constraints []BaseAttributeConstraint `json:"constraints,omitempty"`

	AttributeCallback *AttributeCallback `json:"attributeCallback,omitempty"`
}

func (obj DataAttribute) GetUuid() string {
	return obj.Uuid
}

func (obj DataAttribute) GetName() string {
	return obj.Name
}

func (obj DataAttribute) GetAttributeType() AttributeType {
	return obj.Type
}

func (obj DataAttribute) GetAttributeContentType() AttributeContentType {
	return obj.ContentType
}

// AssertDataAttributeRequired checks if the required fields are not zero-ed
func AssertDataAttributeRequired(obj DataAttribute) error {
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

	//for _, el := range obj.Content {
	//	if err := AssertBaseAttributeContentDtoRequired(el); err != nil {
	//		return err
	//	}
	//}
	//if err := AssertDataAttributePropertiesRequired(obj.Properties); err != nil {
	//	return err
	//}
	//for _, el := range obj.Constraints {
	//	if err := AssertBaseAttributeConstraintRequired(el); err != nil {
	//		return err
	//	}
	//}
	//if err := AssertAttributeCallbackRequired(obj.AttributeCallback); err != nil {
	//	return err
	//}
	return nil
}

// AssertDataAttributeConstraints checks if the values respects the defined constraints
func AssertDataAttributeConstraints(obj DataAttribute) error {
	return nil
}
