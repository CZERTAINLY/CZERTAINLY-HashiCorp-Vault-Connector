/*
 * CZERTAINLY Authority Provider v2 API
 *
 * REST API for implementations of custom v2 Authority Provider
 *
 * API version: 2.11.0
 * Contact: getinfo@czertainly.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package authority

type CredentialAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	Data CredentialAttributeContentData `json:"data"`
}

// AssertCredentialAttributeContentRequired checks if the required fields are not zero-ed
func AssertCredentialAttributeContentRequired(obj CredentialAttributeContent) error {
	elements := map[string]interface{}{
		"data": obj.Data,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertCredentialAttributeContentDataRequired(obj.Data); err != nil {
		return err
	}
	return nil
}

// AssertCredentialAttributeContentConstraints checks if the values respects the defined constraints
func AssertCredentialAttributeContentConstraints(obj CredentialAttributeContent) error {
	return nil
}