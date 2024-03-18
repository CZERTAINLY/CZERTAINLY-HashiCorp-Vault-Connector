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
	Name                 string
	Uuid                 string
	AttributeType        AttributeType
	AttributeContentType AttributeContentType
}

type AttributeConstraint interface {
	GetConstraintType() AttributeConstraintType
}

type Unmarshalable interface {
	Unmarshal(json []byte)
}
