package model

// GroupAttribute - Group attribute and its content represents dynamic list of additional attributes retrieved by callback. Its content can not be edited and is not send in requests to store.
type GroupAttribute struct {
	// UUID of the Attribute for unique identification
	Uuid string `json:"uuid"`

	// Name of the Attribute that is used for identification
	Name string `json:"name"`

	// Optional description of the Attribute, should contain helper text on what is expected
	Description string `json:"description,omitempty"`

	// List of all different types of attributes
	Content []AttributeContent `json:"content,omitempty"`

	Type              AttributeType      `json:"type"`
	AttributeCallback *AttributeCallback `json:"attributeCallback,omitempty"`
}

func (obj GroupAttribute) GetUuid() string {
	return obj.Uuid
}

func (obj GroupAttribute) GetName() string {
	return obj.Name
}

func (obj GroupAttribute) GetAttributeType() AttributeType {
	return obj.Type
}

func (obj GroupAttribute) GetAttributeContentType() AttributeContentType {
	return ""
}

func (obj GroupAttribute) GetContent() []AttributeContent {
	return []AttributeContent{}
}

// AssertGroupAttributeRequired checks if the required fields are not zero-ed
func AssertGroupAttributeRequired(obj GroupAttribute) error {
	elements := map[string]interface{}{
		"uuid": obj.Uuid,
		"name": obj.Name,
		"type": obj.Type,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	//for _, el := range obj.Content {
	//	if err := AssertBaseAttributeDtoRequired(el); err != nil {
	//		return err
	//	}
	//}
	if err := AssertAttributeCallbackRequired(*obj.AttributeCallback); err != nil {
		return err
	}
	return nil
}

// AssertGroupAttributeConstraints checks if the values respects the defined constraints
func AssertGroupAttributeConstraints(obj GroupAttribute) error {
	return nil
}
