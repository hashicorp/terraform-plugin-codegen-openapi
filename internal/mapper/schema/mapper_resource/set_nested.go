package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperSetNestedAttribute struct {
	resource.SetNestedAttribute

	Name         string
	NestedObject MapperNestedAttributeObject
}

func (a *MapperSetNestedAttribute) GetName() string {
	return a.Name
}

func (a *MapperSetNestedAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	setNestedAttribute, ok := mergeAttribute.(*MapperSetNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(setNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *MapperSetNestedAttribute) ToSpec() resource.Attribute {
	a.SetNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:      a.Name,
		SetNested: &a.SetNestedAttribute,
	}
}
