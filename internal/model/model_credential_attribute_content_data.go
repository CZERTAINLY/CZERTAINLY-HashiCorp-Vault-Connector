package model

// CredentialAttributeContentData - Credential attribute content data
type CredentialAttributeContentData struct {

	// Object identifier
	Uuid string `json:"uuid"`

	// Object Name
	Name string `json:"name"`

	// Credential Kind
	Kind string `json:"kind"`

	// List of Credential Attributes
	Attributes []DataAttribute `json:"attributes"`
}

// // AssertCredentialAttributeContentDataRequired checks if the required fields are not zero-ed
// func AssertCredentialAttributeContentDataRequired(obj CredentialAttributeContentData) error {
// 	elements := map[string]interface{}{
// 		"uuid":       obj.Uuid,
// 		"name":       obj.Name,
// 		"kind":       obj.Kind,
// 		"attributes": obj.Attributes,
// 	}
// 	for name, el := range elements {
// 		if isZero := IsZeroValue(el); isZero {
// 			return &RequiredError{Field: name}
// 		}
// 	}

// 	for _, el := range obj.Attributes {
// 		if err := AssertDataAttributeRequired(el); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// AssertCredentialAttributeContentDataConstraints checks if the values respects the defined constraints
func AssertCredentialAttributeContentDataConstraints(obj CredentialAttributeContentData) error {
	return nil
}
