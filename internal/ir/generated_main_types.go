// ir package contains Go bindings for the Framework IR JSON schema. This should eventually be deleted in favor of official bindings
package ir

type IntermediateRepresentation struct {
	Provider    Provider     `json:"provider"`
	Resources   []Resource   `json:"resources,omitempty"`
	DataSources []DataSource `json:"datasources,omitempty"`
}

type Provider struct {
	Name string `json:"name"`
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
