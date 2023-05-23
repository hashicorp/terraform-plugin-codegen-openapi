package util

import "github.com/hashicorp/terraform-plugin-codegen-spec/schema"

// TODO: hackish conversion ¯\_(ツ)_/¯
func ConvertToAttributeType(name string, elemType schema.ElementType) schema.ObjectAttributeType {
	attrType := schema.ObjectAttributeType{Name: name}

	switch {
	case elemType.Bool != nil:
		attrType.Bool = elemType.Bool
	case elemType.Float64 != nil:
		attrType.Float64 = elemType.Float64
	case elemType.Int64 != nil:
		attrType.Int64 = elemType.Int64
	case elemType.List != nil:
		attrType.List = elemType.List
	case elemType.Map != nil:
		attrType.Map = elemType.Map
	case elemType.Number != nil:
		attrType.Number = elemType.Number
	case elemType.Object != nil:
		attrType.Object = elemType.Object
	case elemType.Set != nil:
		attrType.Set = elemType.Set
	case elemType.String != nil:
		attrType.String = elemType.String
	}

	return attrType
}

// TODO: hackish conversion ¯\_(ツ)_/¯
func ConvertToElementType(attrType schema.ObjectAttributeType) schema.ElementType {
	elemType := schema.ElementType{}

	switch {
	case attrType.Bool != nil:
		elemType.Bool = attrType.Bool
	case attrType.Float64 != nil:
		elemType.Float64 = attrType.Float64
	case attrType.Int64 != nil:
		elemType.Int64 = attrType.Int64
	case attrType.List != nil:
		elemType.List = attrType.List
	case attrType.Map != nil:
		elemType.Map = attrType.Map
	case attrType.Number != nil:
		elemType.Number = attrType.Number
	case attrType.Object != nil:
		elemType.Object = attrType.Object
	case attrType.Set != nil:
		elemType.Set = attrType.Set
	case attrType.String != nil:
		elemType.String = attrType.String
	}

	return elemType
}
