package model

// AttributeCallback - Optional definition of callback for getting the content of the Attribute based on the action
type AttributeCallback struct {

	// Context part of callback URL
	CallbackContext string `json:"callbackContext"`

	// HTTP method of the callback
	CallbackMethod string `json:"callbackMethod"`

	// Mappings for the callback method
	Mappings []AttributeCallbackMapping `json:"mappings"`
}

// AssertAttributeCallbackRequired checks if the required fields are not zero-ed
func AssertAttributeCallbackRequired(obj AttributeCallback) error {
	elements := map[string]interface{}{
		"callbackContext": obj.CallbackContext,
		"callbackMethod":  obj.CallbackMethod,
		"mappings":        obj.Mappings,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Mappings {
		if err := AssertAttributeCallbackMappingRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertAttributeCallbackConstraints checks if the values respects the defined constraints
func AssertAttributeCallbackConstraints(obj AttributeCallback) error {
	return nil
}
