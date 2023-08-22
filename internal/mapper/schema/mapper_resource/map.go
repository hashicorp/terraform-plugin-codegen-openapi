package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperMapAttribute struct {
	resource.MapAttribute

	Name string
}

func (a *MapperMapAttribute) GetName() string {
	return a.Name
}

func (a *MapperMapAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	mapAttribute, ok := mergeAttribute.(*MapperMapAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, mapAttribute.ElementType)

	return a
}

func (a *MapperMapAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: a.Name,
		Map:  &a.MapAttribute,
	}
}
