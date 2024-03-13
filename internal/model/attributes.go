package model

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

const (
	URL_ATTR                   string = "URL_ATR_UUID"
	CREDENTIAL_TYPE_ATTR       string = "CREDENTIAL_TYPE_ATR_UUID"
	GROUP_CREDENTIAL_TYPE_ATTR string = "GROUP_CREDENTIAL_TYPE_ATR_UUID"
	JWT_TOKEN_ATTR             string = "JWT_TOKEN_ATR_UUID"
	ROLE_ID_ATTR               string = "ROLE_ID_ATR_UUID"
	ROLE_SECRET_ATTR           string = "ROLE_SECRET_ATR_UUID"
	AUTHORITY_ATTR             string = "AUTHORITY_ATTR_UUID"
	RA_PROFILE_ENGINE_ATTR     string = "RA_PROFILE_ENGINE_ATTR_UUID"
	RA_PROFILE_ROLE_ATTR       string = "RA_PROFILE_ROLE_ATTR_UUID"
	RA_PROFILE_AUTHORITY_ATTR  string = "RA_PROFILE_AUTHORITY_ATTR_UUID"
)

type AttributeName string

const (
	KUBERNETES_CRED string = "kubernetes"
	ROLE_CRED       string = "role"
	TOKEN_CRED      string = "jwt"
)

type CredentialType string

func GetCredentialTypeByName(credentialType string) AttributeContent {
	for _, attribute := range GetCredentialTypes() {
		if attribute.GetData() == credentialType {
			return attribute
		}
	}
	return nil
}
func GetCredentialTypes() []AttributeContent {
	return []AttributeContent{

		StringAttributeContent{
			Reference: "Kubernetes token",
			Data:      KUBERNETES_CRED,
		}, StringAttributeContent{
			Reference: "Role ID",
			Data:      ROLE_CRED,
		}, StringAttributeContent{
			Reference: "JWT",
			Data:      TOKEN_CRED,
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

const (
	AuthorityManagementAttributes string = "AuthorityManagementAttributes"
	DisoveryAttributes            string = "DiscoveryAttributes"
	RAProfilesAttributes          string = "RAProfilesAttributes"
)

func GetAttributeListBySet(attributeSet string) []Attribute {
	switch attributeSet {
	case AuthorityManagementAttributes:
		return getAuthorityManagementAttributes()
	case DisoveryAttributes:
		return getDiscoveryAttributes()
	case RAProfilesAttributes:
		return getRAProfilesAttributes()
	}

	return nil
}

func GetAttributeList() []Attribute {
	attributeList := append(getAuthorityManagementAttributes(), getDiscoveryAttributes()...)
	attributeList = append(attributeList, getRAProfilesAttributes()...)
	attributeList = append(attributeList, getAuthorityManagementAttributes()...)
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

	case SECRET:
		secretAttributeContent := SecretAttributeContent{}
		err := json.Unmarshal(content, &secretAttributeContent)
		result = secretAttributeContent
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
		err := json.Unmarshal([]byte(attribute.Raw), &definition)
		if err != nil {
			return nil
		}
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

func getRAProfilesAttributes() []Attribute {
	return []Attribute{
		DataAttribute{
			Uuid:        RA_PROFILE_ENGINE_ATTR,
			Name:        "ra_profile_engine",
			Description: "Selection of RA Profile engine",
			Type:        DATA,
			Content:     nil,
			ContentType: OBJECT,
			Properties: &DataAttributeProperties{
				Label:       "RA Profile PKI engine selection",
				Visible:     true,
				Group:       "",
				Required:    true,
				ReadOnly:    false,
				List:        false,
				MultiSelect: false,
			},
		},
		DataAttribute{
			Uuid:        RA_PROFILE_AUTHORITY_ATTR,
			Name:        "ra_profile_authority",
			Description: "",
			Type:        DATA,
			Content:     nil,
			ContentType: OBJECT,
			Properties: &DataAttributeProperties{
				Label:       "Authority used for PKI engine query",
				Visible:     false,
				Group:       "",
				Required:    true,
				ReadOnly:    false,
				List:        false,
				MultiSelect: false,
			},
		},
		DataAttribute{
			Uuid:        RA_PROFILE_ROLE_ATTR,
			Name:        "ra_profile_role",
			Description: "Select RA Profile Role",
			Type:        DATA,
			Content:     nil,
			ContentType: STRING,
			Properties: &DataAttributeProperties{
				Label:       "RA profile role",
				Visible:     true,
				Group:       "",
				Required:    true,
				ReadOnly:    false,
				List:        true,
				MultiSelect: false,
			},
			AttributeCallback: &AttributeCallback{
				CallbackContext: "/v1/authorityProvider/authorities/{uuid}/raProfileRole/{engineName}/callback",
				CallbackMethod:  "GET",
				Mappings: []AttributeCallbackMapping{
					{
						From:                 "ra_profile_engine.data.engineName",
						AttributeType:        DATA,
						AttributeContentType: STRING,
						To:                   "engineName",
						Targets: []AttributeValueTarget{
							PATH_VARIABLE,
						},
					},
					{
						From:                 "ra_profile_authority.uuid",
						AttributeType:        DATA,
						AttributeContentType: STRING,
						To:                   "uuid",
						Targets: []AttributeValueTarget{
							PATH_VARIABLE,
						},
					},
				},
			},
		},
	}
}

func getDiscoveryAttributes() []Attribute {
	return []Attribute{
		DataAttribute{
			Uuid:        AUTHORITY_ATTR,
			Name:        "authority_to_discover",
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

}

func getAuthorityManagementAttributes() []Attribute {
	return []Attribute{
		DataAttribute{
			Uuid:        URL_ATTR,
			Name:        "authority_url",
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
			Name:        "credentials_type",
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
			Name:        "credential_group",
			Description: "Authority definition for discovery",
			Type:        GROUP,
			AttributeCallback: &AttributeCallback{
				CallbackContext: "v1/authorityProvider/credentialType/{credentialsType}/callback",
				CallbackMethod:  "GET",
				Mappings: []AttributeCallbackMapping{
					{
						From:                 "credentials_type.data",
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
			Name:        "role_id",
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
			Name:        "role_secret",
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
			Name:        "jwt_token",
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
	}
}
