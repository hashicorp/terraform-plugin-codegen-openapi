package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperFloat64Attribute struct {
	resource.Float64Attribute

	Name string
}

func (a *MapperFloat64Attribute) GetName() string {
	return a.Name
}

func (a *MapperFloat64Attribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	float64Attribute, ok := mergeAttribute.(*MapperFloat64Attribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = float64Attribute.Description
	}

	return a
}

func (a *MapperFloat64Attribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:    a.Name,
		Float64: &a.Float64Attribute,
	}
}
