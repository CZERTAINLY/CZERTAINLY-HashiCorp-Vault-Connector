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

// InfoAttribute - Info attribute contains content that is for information purpose or represents additional information for object (metadata). Its content can not be edited and is not send in requests to store.
type InfoAttribute struct {

	// UUID of the Attribute for unique identification
	Uuid string `json:"uuid"`

	// Name of the Attribute that is used for identification
	Name string `json:"name"`

	// Optional description of the Attribute, should contain helper text on what is expected
	Description string `json:"description,omitempty"`

	// Content of the Attribute
	Content []BaseAttributeContentDto `json:"content"`

	Type AttributeType `json:"type"`

	ContentType AttributeContentType `json:"contentType"`

	Properties InfoAttributeProperties `json:"properties"`
}

// AssertInfoAttributeRequired checks if the required fields are not zero-ed
func AssertInfoAttributeRequired(obj InfoAttribute) error {
	elements := map[string]interface{}{
		"uuid":        obj.Uuid,
		"name":        obj.Name,
		"content":     obj.Content,
		"type":        obj.Type,
		"contentType": obj.ContentType,
		"properties":  obj.Properties,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Content {
		if err := AssertBaseAttributeContentDtoRequired(el); err != nil {
			return err
		}
	}
	if err := AssertInfoAttributePropertiesRequired(obj.Properties); err != nil {
		return err
	}
	return nil
}

// AssertInfoAttributeConstraints checks if the values respects the defined constraints
func AssertInfoAttributeConstraints(obj InfoAttribute) error {
	return nil
}
