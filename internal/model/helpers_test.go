package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestUnmarshalAttributeValue(t *testing.T) {
	result := UnmarshalAttributesValues([]byte(JSON_STRING_VALUE_ARR))
	content := GetAttributeFromArrayByUUID(AUTHORITY_URL_ATTR, result).GetContent()[0]
	URL := content.GetData().(string)
	fmt.Println(URL)
}

func TestUnmarshalAttribute(t *testing.T) {
	result := UnmarshalAttributes([]byte(JSON_STRING_ARR))
	fmt.Println(result)
	resultString, _ := json.Marshal(result)
	var unmarshaled []interface{}
	err := json.Unmarshal([]byte(JSON_STRING_ARR), &unmarshaled)
	if err != nil {
		return
	}
	marshaled, _ := json.Marshal(unmarshaled)
	if equal, err := compareJSON(string(marshaled), string(resultString)); err != nil {
		t.Fatalf("Error comparing JSON:")
	} else if equal {
		fmt.Println("JSON strings are equal")
	} else {
		t.Fatalf("JSON strings are not equal")
	}

}

func compareJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(o1, o2), nil
}

const (
	JSON_STRING_ARR = `
[
{
"uuid": "8a68156a-d1f5-4322-b2a5-26e872a6fc0e",
"name": "Attribute",
"description": "string",
"type": "data",
"content": [
{
"reference": "string",
"data": "bola"
}
],
"contentType": "string",
"properties": {
"label": "Attribute Name",
"visible": true,
"group": "requiredAttributes",
"required": false,
"readOnly": false,
"list": false,
"multiSelect": false
}
}
]	
`

	JSON_STRING_VALUE_ARR = `
[
    {
        "uuid": "8a68156a-d1f5-4322-b2a5-26e872a6fc0e",
        "name": "authority_url",
        "content": [
            {
                "reference": "string",
                "data": "bola"
            }
        ]
    },
    {
        "uuid": "85197836-2ceb-4e77-b14e-53d2e9761cfc",
        "name": "CredentialsType",
        "description": "Authority definition for discovery",
        "content": [
            {
                "reference": "Kubernetes token",
                "data": "kubernetes"
            }
        ]
    }
]
`
)
