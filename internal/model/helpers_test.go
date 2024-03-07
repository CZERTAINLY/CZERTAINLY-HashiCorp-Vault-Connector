package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestUnmarshalAttributeValue(t *testing.T) {
	result := UnmarshalAttributesValues([]byte(JSON_STRING_VALUE_ARR))
	content := GetAttributeFromArrayByUUID(URL_ATTR, result).GetContent()[0]
	URL := content.GetData().(string)
	fmt.Println(URL)
}

func TestUnmarshalAttribute(t *testing.T) {
	result := UnmarshalAttributes([]byte(JSON_STRING_ARR))
	fmt.Println(result)
	resultString, _ := json.Marshal(result)
	var unmarshaled []interface{}
	json.Unmarshal([]byte(JSON_STRING_ARR), &unmarshaled)
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

//var attributeList1 = []DataAttribute{{
//	Uuid:        "sfgsdf",
//	Name:        "AuthorityAttributeIdentifier",
//	Description: "Authority definition for discovery",
//	Type:        "data",
//	Content:     []AttributeContent{StringAttributeContent{Reference: "string", Data: "s"}},
//	ContentType: "object",
//	Properties: DataAttributeProperties{
//		Label:       "Authority to discover",
//		Visible:     true,
//		Group:       "",
//		Required:    true,
//		ReadOnly:    false,
//		List:        false,
//		MultiSelect: false,
//	},
//},
//}

const (
	JSON_STRING = `
	{
    "uuid": "166b5cf52-63f2-11ec-90d6-0242ac120003",
    "name": "Attribute",
    "description": "string",
    "type": "data",
    "content": [
      {
        "reference": "string",
        "data": "bla"
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
}`

	JSON_STRING_ARR = `
[
{
"uuid": "166b5cf52-63f2-11ec-90d6-0242ac120003",
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
        "uuid": "URL_ATR_UUID",
        "name": "AuthorityUrl",
        "content": [
            {
                "reference": "string",
                "data": "bola"
            }
        ]
    },
    {
        "uuid": "CREDENTIAL_TYPE_ATR_UUID",
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
