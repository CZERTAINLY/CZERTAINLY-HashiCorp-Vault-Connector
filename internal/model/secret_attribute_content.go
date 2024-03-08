package model

type SecretAttributeContent struct {

	// Content Reference
	Reference string `json:"reference,omitempty"`

	Data SecretAttributeContentData `json:"data"`
}

func (a SecretAttributeContent) GetData() interface{} {
	return a.Data
}

func (a SecretAttributeContent) GetReference() string {
	return a.Reference
}

// AssertSecretAttributeContentRequired checks if the required fields are not zero-ed
func AssertSecretAttributeContentRequired(obj SecretAttributeContent) error {
	elements := map[string]interface{}{
		"data": obj.Data,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertSecretAttributeContentDataRequired(obj.Data); err != nil {
		return err
	}
	return nil
}

// AssertSecretAttributeContentConstraints checks if the values respects the defined constraints
func AssertSecretAttributeContentConstraints(obj SecretAttributeContent) error {
	return nil
}
