package model

import (
	"time"
)

// DateTimeAttributeConstraintData - DateTime Range Attribute Constraint Data
type DateTimeAttributeConstraintData struct {

	// Start of the datetime for validation
	From time.Time `json:"from,omitempty"`

	// End of the datetime for validation
	To time.Time `json:"to,omitempty"`
}

// AssertDateTimeAttributeConstraintDataRequired checks if the required fields are not zero-ed
func AssertDateTimeAttributeConstraintDataRequired(obj DateTimeAttributeConstraintData) error {
	return nil
}

// AssertDateTimeAttributeConstraintDataConstraints checks if the values respects the defined constraints
func AssertDateTimeAttributeConstraintDataConstraints(obj DateTimeAttributeConstraintData) error {
	return nil
}
