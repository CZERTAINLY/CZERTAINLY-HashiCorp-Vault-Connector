package model

// RegexpAttributeConstraint - Regular Expression Attribute Constraint implements the BaseAttributeConstraint
type RegexpAttributeConstraint struct {

	// Description of the constraint
	Description string `json:"description,omitempty"`

	// Error message to be displayed for wrong data
	ErrorMessage string `json:"errorMessage,omitempty"`

	Type AttributeConstraintType `json:"type"`

	// Regular Expression Attribute Constraint Data
	Data string `json:"data,omitempty"`
}

func (a RegexpAttributeConstraint) GetConstraintType() AttributeConstraintType {
	return a.Type
}

// AssertRegexpAttributeConstraintRequired checks if the required fields are not zero-ed
func AssertRegexpAttributeConstraintRequired(obj RegexpAttributeConstraint) error {
	elements := map[string]interface{}{
		"type": obj.Type,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRegexpAttributeConstraintConstraints checks if the values respects the defined constraints
func AssertRegexpAttributeConstraintConstraints(obj RegexpAttributeConstraint) error {
	return nil
}
