package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperBoolAttribute struct {
	resource.BoolAttribute

	Name string
}

func (a *MapperBoolAttribute) GetName() string {
	return a.Name
}

func (a *MapperBoolAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	boolAttribute, ok := mergeAttribute.(*MapperBoolAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = boolAttribute.Description
	}

	return a
}

func (a *MapperBoolAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: a.Name,
		Bool: &a.BoolAttribute,
	}
}
