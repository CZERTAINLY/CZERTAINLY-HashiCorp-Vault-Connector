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

type DiscoveryRequestDto struct {

	// Name of the Discovery
	Name string `json:"name"`

	// Discovery Kind
	Kind string `json:"kind"`

	// Discovery Provider Attributes. Mandatory for creating new Discovery
	Attributes []RequestAttributeDto `json:"attributes,omitempty"`
}

// AssertDiscoveryRequestDtoRequired checks if the required fields are not zero-ed
func AssertDiscoveryRequestDtoRequired(obj DiscoveryRequestDto) error {
	elements := map[string]interface{}{
		"name": obj.Name,
		"kind": obj.Kind,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Attributes {
		if err := AssertRequestAttributeDtoRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertDiscoveryRequestDtoConstraints checks if the values respects the defined constraints
func AssertDiscoveryRequestDtoConstraints(obj DiscoveryRequestDto) error {
	return nil
}