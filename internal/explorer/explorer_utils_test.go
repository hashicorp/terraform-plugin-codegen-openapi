// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package explorer

import (
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
)

func TestReadOpParameters_Resource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		readOp       *high.Operation
		commonParams []*high.Parameter
		want         []*high.Parameter
	}{
		"merge common and operation": {
			readOp: &high.Operation{
				Parameters: []*high.Parameter{
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
			readOp: &high.Operation{
				Parameters: []*high.Parameter{
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
			readOp: &high.Operation{
				Parameters: nil,
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
			},
		},
		"all nil": {
			readOp:       nil,
			commonParams: nil,
			want:         []*high.Parameter{},
		},
		"read overrides common": {
			readOp: &high.Operation{
				Parameters: []*high.Parameter{
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
			readOp: nil,
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resource := Resource{
				ReadOp:           testCase.readOp,
				CommonParameters: testCase.commonParams,
			}

			mergedParameters := resource.ReadOpParameters()

			if diff := cmp.Diff(mergedParameters, testCase.want, cmpopts.IgnoreUnexported(sync.Mutex{}), cmp.AllowUnexported(base.Schema{}, base.SchemaProxy{}, high.Parameter{})); diff != "" {
				t.Errorf("unexpected difference for resource: %s", diff)
			}
		})
	}
}

func TestReadOpParameters_DataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		readOp       *high.Operation
		commonParams []*high.Parameter
		want         []*high.Parameter
	}{
		"merge common and operation": {
			readOp: &high.Operation{
				Parameters: []*high.Parameter{
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
			readOp: &high.Operation{
				Parameters: []*high.Parameter{
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
			readOp: &high.Operation{
				Parameters: nil,
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
			},
		},
		"all nil": {
			readOp:       nil,
			commonParams: nil,
			want:         []*high.Parameter{},
		},
		"read overrides common": {
			readOp: &high.Operation{
				Parameters: []*high.Parameter{
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
			readOp: nil,
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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dataSource := DataSource{
				ReadOp:           testCase.readOp,
				CommonParameters: testCase.commonParams,
			}

			mergedParameters := dataSource.ReadOpParameters()

			if diff := cmp.Diff(mergedParameters, testCase.want, cmpopts.IgnoreUnexported(sync.Mutex{}), cmp.AllowUnexported(base.Schema{}, base.SchemaProxy{}, high.Parameter{})); diff != "" {
				t.Errorf("unexpected difference for data source: %s", diff)
			}
		})
	}
}

func pointer[T any](value T) *T {
	return &value
}
