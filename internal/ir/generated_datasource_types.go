package ir

type DataSource struct {
	Name   string           `json:"name"`
	Schema DataSourceSchema `json:"schema"`
}

type DataSourceSchema struct {
	Attributes []DataSourceAttribute `json:"attributes,omitempty"`
}

type DataSourceAttribute struct {
	Name string `json:"name"`

	Bool         *DataSourceBoolAttribute         `json:"bool,omitempty"`
	String       *DataSourceStringAttribute       `json:"string,omitempty"`
	Int64        *DataSourceInt64Attribute        `json:"int64,omitempty"`
	Number       *DataSourceNumberAttribute       `json:"number,omitempty"`
	Float64      *DataSourceFloat64Attribute      `json:"float64,omitempty"`
	List         *DataSourceListAttribute         `json:"list,omitempty"`
	ListNested   *DataSourceListNestedAttribute   `json:"list_nested,omitempty"`
	SingleNested *DataSourceSingleNestedAttribute `json:"single_nested,omitempty"`
}

type DataSourceBoolAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type DataSourceStringAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
	Sensitive                *bool                    `json:"sensitive,omitempty"`
}

type DataSourceInt64Attribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type DataSourceNumberAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type DataSourceFloat64Attribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type DataSourceListAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
	ElementType              ElementType              `json:"element_type"`
}

type DataSourceListNestedAttribute struct {
	ComputedOptionalRequired ComputedOptionalRequired        `json:"computed_optional_required"`
	Description              *string                         `json:"description,omitempty"`
	NestedObject             DataSourceAttributeNestedObject `json:"nested_object"`
}

type DataSourceSingleNestedAttribute struct {
	Attributes               []DataSourceAttribute    `json:"attributes"`
	ComputedOptionalRequired ComputedOptionalRequired `json:"computed_optional_required"`
	Description              *string                  `json:"description,omitempty"`
}

type DataSourceAttributeNestedObject struct {
	Attributes []DataSourceAttribute `json:"attributes"`
}
