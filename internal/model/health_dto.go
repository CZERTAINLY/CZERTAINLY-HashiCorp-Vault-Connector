package model

type HealthDto struct {
	Status HealthStatus `json:"status"`

	// Detailed status description
	Description string `json:"description,omitempty"`

	// Nested status of services
	Parts map[string]HealthDto `json:"parts,omitempty"`
}

// AssertHealthDtoRequired checks if the required fields are not zero-ed
func AssertHealthDtoRequired(obj HealthDto) error {
	elements := map[string]interface{}{
		"status": obj.Status,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertHealthDtoConstraints checks if the values respects the defined constraints
func AssertHealthDtoConstraints(obj HealthDto) error {
	return nil
}
