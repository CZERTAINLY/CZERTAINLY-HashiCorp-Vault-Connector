package model

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

var (
	URL_ATTR                   string = "URL_ATR_UUID"
	CREDENTIAL_TYPE_ATTR       string = "CREDENTIAL_TYPE_ATR_UUID"
	GROUP_CREDENTIAL_TYPE_ATTR string = "GROUP_CREDENTIAL_TYPE_ATR_UUID"
	JWT_TOKEN_ATTR             string = "JWT_TOKEN_ATR_UUID"
	ROLE_ID_ATTR               string = "ROLE_ID_ATR_UUID"
	ROLE_SECRET_ATTR           string = "ROLE_SECRET_ATR_UUID"
	AUTHORITY_ATTR             string = "AUTHORITY_ATTR_UUID"
)

func GetCredentialTypes() []AttributeContent {
	return []AttributeContent{

		StringAttributeContent{
			Reference: "Kubernetes token",
			Data:      "kubernetes",
		}, StringAttributeContent{
			Reference: "Role ID",
			Data:      "role",
		}, StringAttributeContent{
			Reference: "JWT",
			Data:      "jwt",
		},
	}
}

func GetAttributeDefByUUID(uuid string) Attribute {
	for _, attr := range GetAttributeList() {
		if attr.GetUuid() == uuid {
			return attr
		}
	}
	return nil
}

func GetAttributeList() []Attribute {
	attributeList := []Attribute{
		DataAttribute{
			Uuid:        URL_ATTR,
			Name:        "AuthorityUrl",
			Description: "Authority definition for discovery",
			Type:        DATA,
			Content:     nil,
			ContentType: STRING,
			Properties: &DataAttributeProperties{
				Label:       "Authority to discover",
				Visible:     true,
				Group:       "",
				Required:    true,
				ReadOnly:    false,
				List:        true,
				MultiSelect: false,
			},
		},
		DataAttribute{
			Uuid:        CREDENTIAL_TYPE_ATTR,
			Name:        "CredentialsType",
			Description: "Authority definition for discovery",
			Type:        DATA,
			Content:     nil,
			ContentType: STRING,
			Properties: &DataAttributeProperties{
				Label:       "Authority to discover",
				Visible:     true,
				Group:       "",
				Required:    true,
				ReadOnly:    false,
				List:        true,
				MultiSelect: false,
			},
		},
		GroupAttribute{
			Uuid:        GROUP_CREDENTIAL_TYPE_ATTR,
			Name:        "CredentialsGroup",
			Description: "Authority definition for discovery",
			Type:        GROUP,
			AttributeCallback: AttributeCallback{
				CallbackContext: "v1/authorityProvider/{credentialsType}/callback",
				CallbackMethod:  "GET",
				Mappings: []AttributeCallbackMapping{
					{
						From:                 "CredentialsType",
						AttributeType:        DATA,
						AttributeContentType: STRING,
						To:                   "credentialsType",
						Targets: []AttributeValueTarget{
							PATH_VARIABLE,
						},
					},
				},
			},
		},
		DataAttribute{
			Uuid:        ROLE_ID_ATTR,
			Name:        "RoleIdType",
			Description: "RoleId for connection to vault",
			Type:        DATA,
			Content:     nil,
			ContentType: SECRET,
			Properties: &DataAttributeProperties{
				Label:       "Role ID",
				Visible:     true,
				Group:       "",
				Required:    true,
				ReadOnly:    false,
				List:        true,
				MultiSelect: false,
			},
		},
		DataAttribute{
			Uuid:        ROLE_SECRET_ATTR,
			Name:        "Role Secret",
			Description: "RoleSecret for connection to vault",
			Type:        DATA,
			Content:     nil,
			ContentType: SECRET,
			Properties: &DataAttributeProperties{
				Label:       "Role Secret",
				Visible:     true,
				Group:       "",
				Required:    true,
				ReadOnly:    false,
				List:        true,
				MultiSelect: false,
			},
		},
		DataAttribute{
			Uuid:        JWT_TOKEN_ATTR,
			Name:        "JWT Token",
			Description: "JWT Token for connection to vault",
			Type:        DATA,
			Content:     nil,
			ContentType: SECRET,
			Properties: &DataAttributeProperties{
				Label:       "JWT Token",
				Visible:     true,
				Group:       "",
				Required:    true,
				ReadOnly:    false,
				List:        true,
				MultiSelect: false,
			},
		},
		DataAttribute{
			Uuid:        AUTHORITY_ATTR,
			Name:        "AuthorityAttributeIdentifier",
			Description: "Authority definition for discovery",
			Type:        DATA,
			Content:     nil,
			ContentType: OBJECT,
			Properties: &DataAttributeProperties{
				Label:       "Authority to discover",
				Visible:     true,
				Group:       "",
				Required:    true,
				ReadOnly:    false,
				List:        false,
				MultiSelect: false,
			},
		},
	}
	return attributeList
}

func GetAtributeByUUID(uuid string) AttributeDefinition {
	for _, attr := range GetAttributeList() {
		if attr.GetUuid() == uuid {
			return AttributeDefinition{
				Uuid:                 attr.GetUuid(),
				AttributeType:        attr.GetAttributeType(),
				AttributeContentType: attr.GetAttributeContentType(),
			}
		}
	}
	return AttributeDefinition{}
}

func unmarshalAttributeContent(content []byte, contentType AttributeContentType) AttributeContent {
	var result AttributeContent
	switch contentType {
	case STRING:
		stringContent := StringAttributeContent{}
		err := json.Unmarshal(content, &stringContent)
		result = stringContent
		if err != nil {
			panic(err)
		}
	case OBJECT:
		objectContent := ObjectAttributeContent{}
		err := json.Unmarshal(content, &objectContent)
		result = objectContent
		if err != nil {
			panic(err)
		}
	case BOOLEAN:
		booleanContent := BooleanAttributeContent{}
		err := json.Unmarshal(content, &booleanContent)
		result = booleanContent
		if err != nil {
			panic(err)
		}
	}

	return result
}

func unmarshalAttribute(content []byte, attrDef AttributeDefinition) Attribute {
	var result Attribute
	switch attrDef.AttributeType {
	case DATA:

		data := DataAttribute{}
		contents := gjson.GetBytes(content, "content")
		for _, content := range contents.Array() {
			data.Content = append(data.Content, unmarshalAttributeContent([]byte(content.Raw), attrDef.AttributeContentType))
		}
		data.Uuid = gjson.GetBytes(content, "uuid").String()
		data.Name = gjson.GetBytes(content, "name").String()
		data.Description = gjson.GetBytes(content, "description").String()
		data.Type = attrDef.AttributeType
		data.ContentType = attrDef.AttributeContentType
		properties := gjson.GetBytes(content, "properties").Raw
		if properties != "" {
			err := json.Unmarshal([]byte(properties), &data.Properties)
			if err != nil {
				panic(err)
			}
		}
		constraints := gjson.GetBytes(content, "constrains").Raw
		if constraints != "" {
			err := json.Unmarshal([]byte(constraints), &data.Constraints)
			if err != nil {
				panic(err)
			}
		}
		callbacks := gjson.GetBytes(content, "attributeCallback").Raw
		if callbacks != "" {
			err := json.Unmarshal([]byte(callbacks), &data.AttributeCallback)
			if err != nil {
				panic(err)
			}
		}
		result = data
	}

	return result
}

func unmarshalAttributeValue(content []byte, attrDef AttributeDefinition) Attribute {
	var result Attribute
	switch attrDef.AttributeType {
	case DATA:
		data := GetAttributeDefByUUID(attrDef.Uuid).(DataAttribute)
		contents := gjson.GetBytes(content, "content")
		data.Content = []AttributeContent{}
		for _, content := range contents.Array() {
			data.Content = append(data.Content, unmarshalAttributeContent([]byte(content.Raw), attrDef.AttributeContentType))
		}
		result = data
	}

	return result
}

func UnmarshalAttributesValues(content []byte) []Attribute {
	attributes := gjson.GetBytes(content, "@values")
	var result []Attribute
	for _, attribute := range attributes.Array() {
		def := GetAtributeByUUID(gjson.Get(attribute.Raw, "uuid").String())
		attributeObject := unmarshalAttributeValue([]byte(attribute.Raw), def)
		result = append(result, attributeObject)
	}
	return result
}

func UnmarshalAttributes(content []byte) []Attribute {
	attributes := gjson.GetBytes(content, "@values")
	var result []Attribute
	for _, attribute := range attributes.Array() {
		definition := AttributeDefinition{}
		json.Unmarshal([]byte(attribute.Raw), &definition)
		if definition.AttributeType == "" || definition.AttributeContentType == "" {
			def := GetAtributeByUUID(definition.Uuid)
			definition.AttributeType = def.AttributeType
			definition.AttributeContentType = def.AttributeContentType
		}
		attributeObject := unmarshalAttribute([]byte(attribute.Raw), definition)
		result = append(result, attributeObject)
	}
	return result
}

func GetAttributeFromArrayByUUID(uuid string, attributes []Attribute) Attribute {
	for _, attr := range attributes {
		if attr.GetUuid() == uuid {
			return attr
		}
	}
	return nil
}
