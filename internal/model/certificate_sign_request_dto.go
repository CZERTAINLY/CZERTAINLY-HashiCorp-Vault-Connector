package model

import "github.com/tidwall/gjson"

type CertificateSignRequestDto struct {

	// Certificate sign request (PKCS#10) encoded as Base64 string
	Pkcs10 string `json:"pkcs10"`

	// List of RA Profiles attributes
	RaProfileAttributes []Attribute `json:"raProfileAttributes"`

	// List of Attributes to issue Certificate
	Attributes []Attribute `json:"attributes"`
}

func (a *CertificateSignRequestDto) Unmarshal(json []byte) {
	a.Pkcs10 = gjson.GetBytes(json, "pkcs10").String()
	a.RaProfileAttributes = UnmarshalAttributesValues([]byte(gjson.GetBytes(json, "raProfileAttributes").Raw))
}

// AssertCertificateSignRequestDtoRequired checks if the required fields are not zero-ed
func AssertCertificateSignRequestDtoRequired(obj CertificateSignRequestDto) error {
	elements := map[string]interface{}{
		"pkcs10":              obj.Pkcs10,
		"raProfileAttributes": obj.RaProfileAttributes,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.RaProfileAttributes {
		if err := AssertRequestAttributeDtoRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Attributes {
		if err := AssertRequestAttributeDtoRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertCertificateSignRequestDtoConstraints checks if the values respects the defined constraints
func AssertCertificateSignRequestDtoConstraints(obj CertificateSignRequestDto) error {
	return nil
}
