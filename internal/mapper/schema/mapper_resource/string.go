package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperStringAttribute struct {
	resource.StringAttribute

	Name string
}

func (a *MapperStringAttribute) GetName() string {
	return a.Name
}

func (a *MapperStringAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	stringAttribute, ok := mergeAttribute.(*MapperStringAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = stringAttribute.Description
	}

	return a
}

func (a *MapperStringAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:   a.Name,
		String: &a.StringAttribute,
	}
}
