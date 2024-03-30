package model

import (
	"fmt"
)

// AttributeValueTarget : Set of targets for propagating value.
type AttributeValueTarget string

// List of AttributeValueTarget
const (
	PATH_VARIABLE     AttributeValueTarget = "pathVariable"
	REQUEST_PARAMETER AttributeValueTarget = "requestParameter"
	BODY              AttributeValueTarget = "body"
)

// AllowedAttributeValueTargetEnumValues is all the allowed values of AttributeValueTarget enum
var AllowedAttributeValueTargetEnumValues = []AttributeValueTarget{
	"pathVariable",
	"requestParameter",
	"body",
}

// validAttributeValueTargetEnumValue provides a map of AttributeValueTargets for fast verification of use input
var validAttributeValueTargetEnumValues = map[AttributeValueTarget]struct{}{
	"pathVariable":     {},
	"requestParameter": {},
	"body":             {},
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v AttributeValueTarget) IsValid() bool {
	_, ok := validAttributeValueTargetEnumValues[v]
	return ok
}

// NewAttributeValueTargetFromValue returns a pointer to a valid AttributeValueTarget
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewAttributeValueTargetFromValue(v string) (AttributeValueTarget, error) {
	ev := AttributeValueTarget(v)
	if ev.IsValid() {
		return ev, nil
	} else {
		return "", fmt.Errorf("invalid value '%v' for AttributeValueTarget: valid values are %v", v, AllowedAttributeValueTargetEnumValues)
	}
}

// AssertAttributeValueTargetRequired checks if the required fields are not zero-ed
func AssertAttributeValueTargetRequired(obj AttributeValueTarget) error {
	return nil
}

// AssertAttributeValueTargetConstraints checks if the values respects the defined constraints
func AssertAttributeValueTargetConstraints(obj AttributeValueTarget) error {
	return nil
}
