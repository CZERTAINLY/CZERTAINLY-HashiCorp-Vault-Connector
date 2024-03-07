package model

type DiscoveryDataRequestDto struct {

	// Name of the Discovery
	Name string `json:"name"`

	// Discovery Kind
	Kind string `json:"kind"`

	// Page number for the retrieved certificates
	PageNumber int32 `json:"pageNumber"`

	// Number of certificates per page
	ItemsPerPage int32 `json:"itemsPerPage"`
}

// AssertDiscoveryDataRequestDtoRequired checks if the required fields are not zero-ed
func AssertDiscoveryDataRequestDtoRequired(obj DiscoveryDataRequestDto) error {
	elements := map[string]interface{}{
		"name":         obj.Name,
		"kind":         obj.Kind,
		"pageNumber":   obj.PageNumber,
		"itemsPerPage": obj.ItemsPerPage,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertDiscoveryDataRequestDtoConstraints checks if the values respects the defined constraints
func AssertDiscoveryDataRequestDtoConstraints(obj DiscoveryDataRequestDto) error {
	return nil
}
