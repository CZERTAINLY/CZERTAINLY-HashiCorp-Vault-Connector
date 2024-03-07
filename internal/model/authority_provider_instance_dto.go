package model

type AuthorityProviderInstanceDto struct {

	// Object identifier
	Uuid string `json:"uuid"`

	// Object Name
	Name string `json:"name"`

	// List of Authority instance Attributes
	Attributes []BaseAttributeDto `json:"attributes"`
}

// AssertAuthorityProviderInstanceDtoRequired checks if the required fields are not zero-ed
func AssertAuthorityProviderInstanceDtoRequired(obj AuthorityProviderInstanceDto) error {
	elements := map[string]interface{}{
		"uuid":       obj.Uuid,
		"name":       obj.Name,
		"attributes": obj.Attributes,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	// for _, el := range obj.Attributes {
	// 	if err := AssertBaseAttributeDtoRequired(el); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// AssertAuthorityProviderInstanceDtoConstraints checks if the values respects the defined constraints
func AssertAuthorityProviderInstanceDtoConstraints(obj AuthorityProviderInstanceDto) error {
	return nil
}
