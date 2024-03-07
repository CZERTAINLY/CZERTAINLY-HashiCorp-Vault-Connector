package model

// CodeBlockAttributeContentData - CodeBlock attribute content data
type CodeBlockAttributeContentData struct {
	Language ProgrammingLanguageEnum `json:"language"`

	// Block of the code in Base64. Formatting of the code is specified by variable language
	Code string `json:"code"`
}

// AssertCodeBlockAttributeContentDataRequired checks if the required fields are not zero-ed
func AssertCodeBlockAttributeContentDataRequired(obj CodeBlockAttributeContentData) error {
	elements := map[string]interface{}{
		"language": obj.Language,
		"code":     obj.Code,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertCodeBlockAttributeContentDataConstraints checks if the values respects the defined constraints
func AssertCodeBlockAttributeContentDataConstraints(obj CodeBlockAttributeContentData) error {
	return nil
}
