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

// InfoAttributeProperties - Properties of the Attributes
type InfoAttributeProperties struct {

	// Friendly name of the the Attribute
	Label string `json:"label"`

	// Boolean determining if the Attribute is visible and can be displayed, otherwise it should be hidden to the user.
	Visible bool `json:"visible"`

	// Group of the Attribute, used for the logical grouping of the Attribute
	Group string `json:"group,omitempty"`
}

// AssertInfoAttributePropertiesRequired checks if the required fields are not zero-ed
func AssertInfoAttributePropertiesRequired(obj InfoAttributeProperties) error {
	elements := map[string]interface{}{
		"label":   obj.Label,
		"visible": obj.Visible,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertInfoAttributePropertiesConstraints checks if the values respects the defined constraints
func AssertInfoAttributePropertiesConstraints(obj InfoAttributeProperties) error {
	return nil
}
