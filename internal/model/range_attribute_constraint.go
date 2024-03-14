package model

type RangeAttributeConstraint struct {

	// Description of the constraint
	Description string `json:"description,omitempty"`

	// Error message to be displayed for wrong data
	ErrorMessage string `json:"errorMessage,omitempty"`

	Type AttributeConstraintType `json:"type"`

	Data RangeAttributeConstraintData `json:"data,omitempty"`
}

// AssertRangeAttributeConstraintRequired checks if the required fields are not zero-ed
func AssertRangeAttributeConstraintRequired(obj RangeAttributeConstraint) error {
	elements := map[string]interface{}{
		"type": obj.Type,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertRangeAttributeConstraintDataRequired(obj.Data); err != nil {
		return err
	}
	return nil
}

// AssertRangeAttributeConstraintConstraints checks if the values respects the defined constraints
func AssertRangeAttributeConstraintConstraints(obj RangeAttributeConstraint) error {
	return nil
}
