package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperListNestedAttribute struct {
	resource.ListNestedAttribute

	Name         string
	NestedObject MapperNestedAttributeObject
}

func (a *MapperListNestedAttribute) GetName() string {
	return a.Name
}

func (a *MapperListNestedAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	listNestedAttribute, ok := mergeAttribute.(*MapperListNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = listNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(listNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *MapperListNestedAttribute) ToSpec() resource.Attribute {
	a.ListNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:       a.Name,
		ListNested: &a.ListNestedAttribute,
	}
}
