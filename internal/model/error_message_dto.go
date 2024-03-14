package model

type ErrorMessageDto struct {

	// Error message detail
	Message string `json:"message"`
}

// AssertErrorMessageDtoRequired checks if the required fields are not zero-ed
func AssertErrorMessageDtoRequired(obj ErrorMessageDto) error {
	elements := map[string]interface{}{
		"message": obj.Message,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertErrorMessageDtoConstraints checks if the values respects the defined constraints
func AssertErrorMessageDtoConstraints(obj ErrorMessageDto) error {
	return nil
}
