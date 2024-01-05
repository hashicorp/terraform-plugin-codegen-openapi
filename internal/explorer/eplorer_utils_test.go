package explorer

import (
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func TestMergeParameters(t *testing.T) {
	testCases := map[string]struct {
		noReadOp     bool
		readParams   []*high.Parameter
		commonParams []*high.Parameter
		want         []*high.Parameter
	}{
		"merge common and operation": {
			noReadOp: false,
			readParams: []*high.Parameter{
				{
					Name:        "string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
				{
					Name:     "bool_prop",
					Required: pointer(true),
					In:       "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey this is a bool, required!",
					}),
				},
				{
					Name: "float64_prop",
					In:   "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "float",
						Description: "hey this is a float64!",
					}),
				},
			},
			commonParams: []*high.Parameter{
				{
					Name:        "common_string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
			},
			want: []*high.Parameter{
				{
					Name:        "common_string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
				{
					Name:        "string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
				{
					Name:     "bool_prop",
					Required: pointer(true),
					In:       "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey this is a bool, required!",
					}),
				},
				{
					Name: "float64_prop",
					In:   "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "float",
						Description: "hey this is a float64!",
					}),
				},
			},
		},
		"nil common parameters": {
			noReadOp: false,
			readParams: []*high.Parameter{
				{
					Name:        "string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
				{
					Name:     "bool_prop",
					Required: pointer(true),
					In:       "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey this is a bool, required!",
					}),
				},
				{
					Name: "float64_prop",
					In:   "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "float",
						Description: "hey this is a float64!",
					}),
				},
			},
			commonParams: nil,
			want: []*high.Parameter{
				{
					Name:        "string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
				{
					Name:     "bool_prop",
					Required: pointer(true),
					In:       "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey this is a bool, required!",
					}),
				},
				{
					Name: "float64_prop",
					In:   "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "float",
						Description: "hey this is a float64!",
					}),
				},
			},
		},
		"nil read parameters": {
			noReadOp:   false,
			readParams: nil,
			commonParams: []*high.Parameter{
				{
					Name:        "common_string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
			},
			want: []*high.Parameter{
				{
					Name:        "common_string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
			},
		},
		"all nil": {
			noReadOp:     false,
			readParams:   nil,
			commonParams: nil,
			want:         nil,
		},
		"read overrides common": {
			noReadOp: false,
			readParams: []*high.Parameter{
				{
					Name:        "string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
				{
					Name:     "bool_prop",
					Required: pointer(true),
					In:       "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey this is a bool, required!",
					}),
				},
				{
					Name: "float64_prop",
					In:   "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "float",
						Description: "hey this is a float64!",
					}),
				},
			},
			commonParams: []*high.Parameter{
				{
					Name:        "string_prop",
					Required:    pointer(false),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
			},
			want: []*high.Parameter{
				{
					Name:        "string_prop",
					Required:    pointer(true),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
				{
					Name:     "bool_prop",
					Required: pointer(true),
					In:       "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey this is a bool, required!",
					}),
				},
				{
					Name: "float64_prop",
					In:   "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "float",
						Description: "hey this is a float64!",
					}),
				},
			},
		},
		"no read op": {
			noReadOp:   true,
			readParams: nil,
			commonParams: []*high.Parameter{
				{
					Name:        "string_prop",
					Required:    pointer(false),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
			},
			want: []*high.Parameter{
				{
					Name:        "string_prop",
					Required:    pointer(false),
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var readOp *high.Operation = nil
			if !testCase.noReadOp {
				readOp = &high.Operation{
					Parameters: testCase.readParams,
				}
			}

			resource := Resource{
				ReadOp:     readOp,
				Parameters: testCase.commonParams,
			}

			mergedParameters := resource.ReadOpParameters()

			if diff := cmp.Diff(mergedParameters, testCase.want, cmp.AllowUnexported(base.Schema{}, base.SchemaProxy{}, sync.Mutex{}, high.Parameter{})); diff != "" {
				t.Errorf("unexpected difference for resource: %s", diff)
			}

			dataSource := DataSource{
				ReadOp: &high.Operation{
					Parameters: testCase.readParams,
				},
				Parameters: testCase.commonParams,
			}

			mergedParameters = dataSource.ReadOpParameters()

			if diff := cmp.Diff(mergedParameters, testCase.want, cmp.AllowUnexported(base.Schema{}, base.SchemaProxy{}, sync.Mutex{}, high.Parameter{})); diff != "" {
				t.Errorf("unexpected difference for data source: %s", diff)
			}
		})
	}
}

func pointer[T any](value T) *T {
	return &value
}
