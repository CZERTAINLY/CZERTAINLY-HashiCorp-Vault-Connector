package model

import (
	"time"
)

type DateTimeAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitzero"`

	// DateTime attribute value in format yyyy-MM-ddTHH:mm:ss.SSSXXX
	Data time.Time `json:"data"`
}

func (d DateTimeAttributeContent) GetData() any {
	return d.Data
}

func (d DateTimeAttributeContent) GetReference() string {
	return d.Reference
}

// AssertDateTimeAttributeContentRequired checks if the required fields are not zero-ed
func AssertDateTimeAttributeContentRequired(obj DateTimeAttributeContent) error {
	elements := map[string]any{
		"data": obj.Data,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertDateTimeAttributeContentConstraints checks if the values respects the defined constraints
func AssertDateTimeAttributeContentConstraints(obj DateTimeAttributeContent) error {
	return nil
}
