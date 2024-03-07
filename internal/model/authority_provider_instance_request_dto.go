package model

type AuthorityProviderInstanceRequestDto struct {

	// Authority instance name
	Name string `json:"name"`

	// Kind of Authority instance
	Kind string `json:"kind"`

	// List of Authority instance Attributes
	Attributes []Attribute `json:"attributes"`
}

//func (a *AuthorityProviderInstanceRequestDto) Unmarshal(json []byte) {
//	a.Name = gjson.GetBytes(json, "name").String()
//	a.Kind = gjson.GetBytes(json, "kind").String()
//	a.Attributes = UnmarshalAttributesValues([]byte(gjson.GetBytes(json, "attributes").Raw))
//}

// AssertAuthorityProviderInstanceRequestDtoRequired checks if the required fields are not zero-ed
func AssertAuthorityProviderInstanceRequestDtoRequired(obj AuthorityProviderInstanceRequestDto) error {
	elements := map[string]interface{}{
		"name":       obj.Name,
		"kind":       obj.Kind,
		"attributes": obj.Attributes,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	//for _, el := range obj.Attributes {
	//	if err := AssertRequestAttributeDtoRequired(el); err != nil {
	//		return err
	//	}
	//}
	return nil
}

// AssertAuthorityProviderInstanceRequestDtoConstraints checks if the values respects the defined constraints
func AssertAuthorityProviderInstanceRequestDtoConstraints(obj AuthorityProviderInstanceRequestDto) error {
	return nil
}
