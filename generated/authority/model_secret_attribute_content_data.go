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

// SecretAttributeContentData - Secret attribute content data
type SecretAttributeContentData struct {

	// Secret attribute data
	Secret string `json:"secret,omitempty"`

	// Level of protection of the data
	ProtectionLevel string `json:"protectionLevel,omitempty"`
}

// AssertSecretAttributeContentDataRequired checks if the required fields are not zero-ed
func AssertSecretAttributeContentDataRequired(obj SecretAttributeContentData) error {
	return nil
}

// AssertSecretAttributeContentDataConstraints checks if the values respects the defined constraints
func AssertSecretAttributeContentDataConstraints(obj SecretAttributeContentData) error {
	return nil
}