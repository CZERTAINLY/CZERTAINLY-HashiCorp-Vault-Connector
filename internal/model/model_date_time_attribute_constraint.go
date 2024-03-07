package model

type DateTimeAttributeConstraint struct {

	// Description of the constraint
	Description string `json:"description,omitempty"`

	// Error message to be displayed for wrong data
	ErrorMessage string `json:"errorMessage,omitempty"`

	Type AttributeConstraintType `json:"type"`

	Data DateTimeAttributeConstraintData `json:"data,omitempty"`
}

// AssertDateTimeAttributeConstraintRequired checks if the required fields are not zero-ed
func AssertDateTimeAttributeConstraintRequired(obj DateTimeAttributeConstraint) error {
	elements := map[string]interface{}{
		"type": obj.Type,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertDateTimeAttributeConstraintDataRequired(obj.Data); err != nil {
		return err
	}
	return nil
}

// AssertDateTimeAttributeConstraintConstraints checks if the values respects the defined constraints
func AssertDateTimeAttributeConstraintConstraints(obj DateTimeAttributeConstraint) error {
	return nil
}
