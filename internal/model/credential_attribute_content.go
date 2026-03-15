package model

type CredentialAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	Data CredentialAttributeContentData `json:"data"`
}

// // AssertCredentialAttributeContentRequired checks if the required fields are not zero-ed
// func AssertCredentialAttributeContentRequired(obj CredentialAttributeContent) error {
// 	elements := map[string]interface{}{
// 		"data": obj.Data,
// 	}
// 	for name, el := range elements {
// 		if isZero := IsZeroValue(el); isZero {
// 			return &RequiredError{Field: name}
// 		}
// 	}

// 	if err := AssertCredentialAttributeContentDataRequired(obj.Data); err != nil {
// 		return err
// 	}
// 	return nil
// }

// AssertCredentialAttributeContentConstraints checks if the values respects the defined constraints
func AssertCredentialAttributeContentConstraints(obj CredentialAttributeContent) error {
	return nil
}
