package schema

import (
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/util"
)

type ElementTypeBuilder func() (*ir.ElementType, error)

func (s *OASSchema) BuildElementType() (*ir.ElementType, error) {
	switch s.Type {
	case util.OAS_type_string:
		return s.BuildStringElementType()
	case util.OAS_type_integer:
		return s.BuildIntegerElementType()
	case util.OAS_type_number:
		return s.BuildNumberElementType()
	case util.OAS_type_boolean:
		return s.BuildBoolElementType()
	case util.OAS_type_array:
		return s.BuildListElementType()
	case util.OAS_type_object:
		return s.BuildObjectElementType()

	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}
