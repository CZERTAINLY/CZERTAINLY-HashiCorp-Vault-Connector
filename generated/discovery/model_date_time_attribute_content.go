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

import (
	"time"
)

type DateTimeAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	// DateTime attribute value in format yyyy-MM-ddTHH:mm:ss.SSSXXX
	Data time.Time `json:"data"`
}

// AssertDateTimeAttributeContentRequired checks if the required fields are not zero-ed
func AssertDateTimeAttributeContentRequired(obj DateTimeAttributeContent) error {
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

// AssertDateTimeAttributeContentConstraints checks if the values respects the defined constraints
func AssertDateTimeAttributeContentConstraints(obj DateTimeAttributeContent) error {
	return nil
}
