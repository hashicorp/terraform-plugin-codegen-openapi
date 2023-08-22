package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperListAttribute struct {
	resource.ListAttribute

	Name string
}

func (a *MapperListAttribute) GetName() string {
	return a.Name
}

func (a *MapperListAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	listAttribute, ok := mergeAttribute.(*MapperListAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = listAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, listAttribute.ElementType)

	return a
}

func (a *MapperListAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: a.Name,
		List: &a.ListAttribute,
	}
}
