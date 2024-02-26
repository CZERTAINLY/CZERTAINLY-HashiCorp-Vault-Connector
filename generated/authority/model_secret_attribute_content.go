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

type SecretAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	Data SecretAttributeContentData `json:"data"`
}

// AssertSecretAttributeContentRequired checks if the required fields are not zero-ed
func AssertSecretAttributeContentRequired(obj SecretAttributeContent) error {
	elements := map[string]interface{}{
		"data": obj.Data,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertSecretAttributeContentDataRequired(obj.Data); err != nil {
		return err
	}
	return nil
}

// AssertSecretAttributeContentConstraints checks if the values respects the defined constraints
func AssertSecretAttributeContentConstraints(obj SecretAttributeContent) error {
	return nil
}