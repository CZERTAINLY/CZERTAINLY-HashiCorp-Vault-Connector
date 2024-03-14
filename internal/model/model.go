package model

type Attribute interface {
	GetContent() []AttributeContent
	GetName() string
	GetUuid() string
	GetAttributeType() AttributeType
	GetAttributeContentType() AttributeContentType
}

type AttributeContent interface {
	GetData() interface{}
	GetReference() string
}

type AttributeDefinition struct {
	Uuid                 string
	AttributeType        AttributeType
	AttributeContentType AttributeContentType
}

type Unmarshalable interface {
	Unmarshal(json []byte)
}
