/*
 * CZERTAINLY Discovery Provider API
 *
 * REST API for implementations of custom Discovery Provider
 *
 * API version: 2.11.0
 * Contact: getinfo@czertainly.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package discovery

// BaseAttributeContentDto - Base Attribute content definition
type BaseAttributeContentDto struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	// Content Data
	Data map[string]interface{} `json:"data"`
}

// AssertBaseAttributeContentDtoRequired checks if the required fields are not zero-ed
func AssertBaseAttributeContentDtoRequired(obj BaseAttributeContentDto) error {
	elements := map[string]interface{}{
		"data": obj.Data,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertBaseAttributeContentDtoConstraints checks if the values respects the defined constraints
func AssertBaseAttributeContentDtoConstraints(obj BaseAttributeContentDto) error {
	return nil
}