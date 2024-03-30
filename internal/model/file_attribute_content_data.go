package model

// FileAttributeContentData - File attribute content data
type FileAttributeContentData struct {

	// File content
	Content string `json:"content"`

	// Name of the file
	FileName string `json:"fileName"`

	// Type of the file uploaded
	MimeType string `json:"mimeType"`
}

// AssertFileAttributeContentDataRequired checks if the required fields are not zero-ed
func AssertFileAttributeContentDataRequired(obj FileAttributeContentData) error {
	elements := map[string]interface{}{
		"content":  obj.Content,
		"fileName": obj.FileName,
		"mimeType": obj.MimeType,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertFileAttributeContentDataConstraints checks if the values respects the defined constraints
func AssertFileAttributeContentDataConstraints(obj FileAttributeContentData) error {
	return nil
}
