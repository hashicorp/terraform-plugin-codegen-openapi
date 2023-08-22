package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperNumberAttribute struct {
	resource.NumberAttribute

	Name string
}

func (a *MapperNumberAttribute) GetName() string {
	return a.Name
}

func (a *MapperNumberAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	numberAttribute, ok := mergeAttribute.(*MapperNumberAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = numberAttribute.Description
	}

	return a
}

func (a *MapperNumberAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:   a.Name,
		Number: &a.NumberAttribute,
	}
}
