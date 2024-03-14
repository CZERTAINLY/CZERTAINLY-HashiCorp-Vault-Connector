package model

// RangeAttributeConstraintData - Integer Range Attribute Constraint Data
type RangeAttributeConstraintData struct {

	// Start of the range for validation
	From int32 `json:"from,omitempty"`

	// End of the range for validation
	To int32 `json:"to,omitempty"`
}

// AssertRangeAttributeConstraintDataRequired checks if the required fields are not zero-ed
func AssertRangeAttributeConstraintDataRequired(obj RangeAttributeConstraintData) error {
	return nil
}

// AssertRangeAttributeConstraintDataConstraints checks if the values respects the defined constraints
func AssertRangeAttributeConstraintDataConstraints(obj RangeAttributeConstraintData) error {
	return nil
}
