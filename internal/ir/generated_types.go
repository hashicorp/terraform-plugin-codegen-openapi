// ir package contains Go bindings for the Framework IR JSON schema. This should eventually be deleted in favor of official bindings
package ir

type IntermediateRepresentation struct {
	Provider  Provider   `json:"provider"`
	Resources []Resource `json:"resources,omitempty"`
}

type Provider struct {
	Name string `json:"name"`
}

type Resource struct {
	Name   string         `json:"name"`
	Schema ResourceSchema `json:"schema"`
}

type ResourceSchema struct {
	Attributes []ResourceAttribute `json:"attributes,omitempty"`
}

type ResourceAttribute struct {
	Name         string                         `json:"name"`
	Bool         *ResourceBoolAttribute         `json:"bool,omitempty"`
	String       *ResourceStringAttribute       `json:"string,omitempty"`
	Int64        *ResourceInt64Attribute        `json:"int64,omitempty"`
	Number       *ResourceNumberAttribute       `json:"number,omitempty"`
	Float64      *ResourceFloat64Attribute      `json:"float64,omitempty"`
	List         *ResourceListAttribute         `json:"list,omitempty"`
	ListNested   *ResourceListNestedAttribute   `json:"list_nested,omitempty"`
	SingleNested *ResourceSingleNestedAttribute `json:"single_nested,omitempty"`
}

type ResourceBoolAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type ResourceStringAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
	Sensitive                *bool                    `json:"sensitive,omitempty"`
}

type ResourceInt64Attribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type ResourceNumberAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type ResourceFloat64Attribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type ResourceListAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
	ElementType              ElementType              `json:"element_type"`
}

type ResourceListNestedAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired      `json:"computed_optional_required"`
	Description              *string                       `json:"description,omitempty"`
	NestedObject             ResourceAttributeNestedObject `json:"nested_object"`
}

type ResourceSingleNestedAttribute struct {
	Attributes               []ResourceAttribute      `json:"attributes"`
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type ResourceAttributeNestedObject struct {
	Attributes []ResourceAttribute `json:"attributes"`
}

type ElementType struct {
	Bool    *BoolElement    `json:"bool,omitempty"`
	Int64   *Int64Element   `json:"int64,omitempty"`
	Float64 *Float64Element `json:"float64,omitempty"`
	Number  *NumberElement  `json:"number,omitempty"`
	String  *StringElement  `json:"string,omitempty"`
	List    *ListElement    `json:"list,omitempty"`
	Object  []ObjectElement `json:"object,omitempty"`
}

type BoolElement struct {
}

type Int64Element struct {
}

type Float64Element struct {
}

type NumberElement struct {
}

type StringElement struct {
}

type ListElement struct {
	*ElementType
}

type ObjectElement struct {
	Name string `json:"name"`
	*ElementType
}

type ComputedOptionalRequired string

const (
	Computed         ComputedOptionalRequired = "computed"
	ComputedOptional ComputedOptionalRequired = "computed_optional"
	Optional         ComputedOptionalRequired = "optional"
	Required         ComputedOptionalRequired = "required"
)
