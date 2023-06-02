// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// TODO: write tests for error paths (nullable types + build schema functions)

func TestBuildSchemaFromRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		op             *high.Operation
		expectedSchema *oas.OASSchema
	}{
		"default to application/json": {
			op: &high.Operation{
				RequestBody: &high.RequestBody{
					Content: map[string]*high.MediaType{
						"application/xml": {
							Schema: base.CreateSchemaProxy(&base.Schema{
								Description: "this is the wrong one!",
								Type:        []string{"boolean"},
							}),
						},
						"application/json": {
							Schema: base.CreateSchemaProxy(&base.Schema{
								Description: "this is the correct one!",
								Type:        []string{"string"},
							}),
						},
					},
				},
			},
			expectedSchema: &oas.OASSchema{
				Type: "string",
				Schema: &base.Schema{
					Description: "this is the correct one!",
					Type:        []string{"string"},
				},
			},
		},
		"utilizes other media types in sorted order": {
			op: &high.Operation{
				RequestBody: &high.RequestBody{
					Content: map[string]*high.MediaType{
						"application/xml": {
							Schema: base.CreateSchemaProxy(&base.Schema{
								Description: "this won't be used because of sorting!",
								Type:        []string{"boolean"},
							}),
						},
						"application/jay-son": {
							Schema: base.CreateSchemaProxy(&base.Schema{
								Description: "this is will get used because of sorting!",
								Type:        []string{"string"},
							}),
						},
					},
				},
			},
			expectedSchema: &oas.OASSchema{
				Type: "string",
				Schema: &base.Schema{
					Description: "this is will get used because of sorting!",
					Type:        []string{"string"},
				},
			},
		},
		"utilizes other media types when nil schemas in priority media types": {
			op: &high.Operation{
				RequestBody: &high.RequestBody{
					Content: map[string]*high.MediaType{
						"application/json": {
							Schema: nil,
						},
						"application/starts-with-ayyy-son": {
							Schema: nil,
						},
						"application/xml": {
							Schema: base.CreateSchemaProxy(&base.Schema{
								Description: "this will get used!",
								Type:        []string{"string"},
							}),
						},
					},
				},
			},
			expectedSchema: &oas.OASSchema{
				Type: "string",
				Schema: &base.Schema{
					Description: "this will get used!",
					Type:        []string{"string"},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := oas.BuildSchemaFromRequest(testCase.op)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			// TODO: this is hacky + not recommended, should see if there is a better comparison method long-term
			if diff := cmp.Diff(got, testCase.expectedSchema, cmpopts.IgnoreUnexported(base.Schema{}, oas.OASSchema{})); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}

}

func TestBuildSchemaFromRequest_Errors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		op               *high.Operation
		expectedErrRegex string
	}{
		"nil op": {
			op:               nil,
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"nil request body": {
			op: &high.Operation{
				RequestBody: nil,
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"empty request body content": {
			op: &high.Operation{
				RequestBody: &high.RequestBody{},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"no media type schemas": {
			op: &high.Operation{
				RequestBody: &high.RequestBody{
					Content: map[string]*high.MediaType{
						"application/json": {
							Schema: nil,
						},
						"application/xml": {
							Schema: nil,
						},
					},
				},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		errRegex := regexp.MustCompile(testCase.expectedErrRegex)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := oas.BuildSchemaFromRequest(testCase.op)

			if err == nil {
				t.Errorf("Expected err to match %q, got nil", testCase.expectedErrRegex)
				return
			}
			if !errRegex.Match([]byte(err.Error())) {
				t.Errorf("Expected error to match %q, got %q", testCase.expectedErrRegex, err.Error())
			}
		})
	}
}

func TestBuildSchemaFromResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		op             *high.Operation
		expectedSchema *oas.OASSchema
	}{
		"default to 200 and application/json": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: map[string]*high.Response{
						"201": {
							Content: map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the wrong one!",
										Type:        []string{"boolean"},
									}),
								},
							},
						},
						"200": {
							Content: map[string]*high.MediaType{
								"application/xml": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the wrong one!",
										Type:        []string{"boolean"},
									}),
								},
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the correct one!",
										Type:        []string{"string"},
									}),
								},
							},
						},
					},
				},
			},
			expectedSchema: &oas.OASSchema{
				Type: "string",
				Schema: &base.Schema{
					Description: "this is the correct one!",
					Type:        []string{"string"},
				},
			},
		},
		"fallback to 201 and application/json": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: map[string]*high.Response{
						"204": {
							Content: map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the wrong one!",
										Type:        []string{"boolean"},
									}),
								},
							},
						},
						"201": {
							Content: map[string]*high.MediaType{
								"application/xml": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the wrong one!",
										Type:        []string{"boolean"},
									}),
								},
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the correct one!",
										Type:        []string{"string"},
									}),
								},
							},
						},
					},
				},
			},
			expectedSchema: &oas.OASSchema{
				Type: "string",
				Schema: &base.Schema{
					Description: "this is the correct one!",
					Type:        []string{"string"},
				},
			},
		},
		"fallback to success code and any media type in sorted order": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: map[string]*high.Response{
						"304": {
							Content: map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the wrong one!",
										Type:        []string{"boolean"},
									}),
								},
							},
						},
						"204": {
							Content: map[string]*high.MediaType{
								"application/xml": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the wrong one!",
										Type:        []string{"boolean"},
									}),
								},
								"application/jay-son": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the correct one!",
										Type:        []string{"string"},
									}),
								},
							},
						},
					},
				},
			},
			expectedSchema: &oas.OASSchema{
				Type: "string",
				Schema: &base.Schema{
					Description: "this is the correct one!",
					Type:        []string{"string"},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := oas.BuildSchemaFromResponse(testCase.op)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			// TODO: this is hacky + not recommended, should see if there is a better comparison method long-term
			if diff := cmp.Diff(got, testCase.expectedSchema, cmpopts.IgnoreUnexported(base.Schema{}, oas.OASSchema{})); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}

}

func TestBuildSchemaFromResponse_Errors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		op               *high.Operation
		expectedErrRegex string
	}{
		"nil op": {
			op:               nil,
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"nil responses": {
			op: &high.Operation{
				Responses: nil,
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"empty response codes": {
			op: &high.Operation{
				Responses: &high.Responses{},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"no success response code media type schemas": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: map[string]*high.Response{
						"300": {
							Content: map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this won't be used!",
										Type:        []string{"string"},
									}),
								},
							},
						},
						"skip-me!": {
							Content: map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this won't be used!",
										Type:        []string{"string"},
									}),
								},
							},
						},
						"199": {
							Content: map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this won't be used!",
										Type:        []string{"string"},
									}),
								},
							},
						},
					},
				},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"200 response code with no valid schema": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: map[string]*high.Response{
						"200": {
							Content: map[string]*high.MediaType{
								"application/json": {
									Schema: nil,
								},
							},
						},
					},
				},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"201 response code with no valid schema": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: map[string]*high.Response{
						"201": {
							Content: map[string]*high.MediaType{
								"application/json": {
									Schema: nil,
								},
							},
						},
					},
				},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"success response code with no valid schema": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: map[string]*high.Response{
						"204": {
							Content: map[string]*high.MediaType{
								"application/json": {
									Schema: nil,
								},
							},
						},
					},
				},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		errRegex := regexp.MustCompile(testCase.expectedErrRegex)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := oas.BuildSchemaFromResponse(testCase.op)

			if err == nil {
				t.Fatalf("Expected err to match %q, got nil", testCase.expectedErrRegex)
			}
			if !errRegex.Match([]byte(err.Error())) {
				t.Errorf("Expected error to match %q, got %q", testCase.expectedErrRegex, err.Error())
			}
		})
	}

}

func TestBuildSchema_NullableMultiTypes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schemaProxy        *base.SchemaProxy
		expectedAttributes *[]resource.Attribute
	}{
		"nullable type - Type array": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"nullable_string_two"},
				Properties: map[string]*base.SchemaProxy{
					"nullable_string_one": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"null", "string"},
						Description: "hey there! I'm a nullable string type.",
					}),
					"nullable_string_two": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string", "null"},
						Description: "hey there! I'm a nullable string type, required.",
					}),
				},
			}),
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nullable_string_one",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a nullable string type."),
						Sensitive:                pointer(false),
					},
				},
				{
					Name: "nullable_string_two",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a nullable string type, required."),
						Sensitive:                pointer(false),
					},
				},
			},
		},
		"nullable type - anyOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"nullable_string_two"},
				Properties: map[string]*base.SchemaProxy{
					"nullable_string_one": base.CreateSchemaProxy(&base.Schema{
						AnyOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null"},
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a string type.",
							}),
						},
					}),
					"nullable_string_two": base.CreateSchemaProxy(&base.Schema{
						AnyOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a string type, required.",
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null"},
							}),
						},
					}),
				},
			}),
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nullable_string_one",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a string type."),
						Sensitive:                pointer(false),
					},
				},
				{
					Name: "nullable_string_two",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a string type, required."),
						Sensitive:                pointer(false),
					},
				},
			},
		},
		"nullable type - oneOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"nullable_string_two"},
				Properties: map[string]*base.SchemaProxy{
					"nullable_string_one": base.CreateSchemaProxy(&base.Schema{
						OneOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null"},
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a string type.",
							}),
						},
					}),
					"nullable_string_two": base.CreateSchemaProxy(&base.Schema{
						OneOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a string type, required.",
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null"},
							}),
						},
					}),
				},
			}),
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nullable_string_one",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a string type."),
						Sensitive:                pointer(false),
					},
				},
				{
					Name: "nullable_string_two",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a string type, required."),
						Sensitive:                pointer(false),
					},
				},
			},
		},
		"list attributes with nullable element type - Type array": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"string_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null", "string"},
							}),
						},
					}),
					"string_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"string", "null"},
							}),
						},
					}),
				},
			}),
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "string_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of nullable strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "string_list_prop_required",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of nullable strings, required."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
			},
		},
		"list attributes with nullable element type - anyOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"string_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								AnyOf: []*base.SchemaProxy{
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"null"},
									}),
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"string"},
									}),
								},
							}),
						},
					}),
					"string_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								AnyOf: []*base.SchemaProxy{
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"string"},
									}),
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"null"},
									}),
								},
							}),
						},
					}),
				},
			}),
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "string_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of nullable strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "string_list_prop_required",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of nullable strings, required."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
			},
		},
		"list attributes with nullable element type - oneOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"string_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								OneOf: []*base.SchemaProxy{
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"null"},
									}),
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"string"},
									}),
								},
							}),
						},
					}),
					"string_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								OneOf: []*base.SchemaProxy{
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"string"},
									}),
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"null"},
									}),
								},
							}),
						},
					}),
				},
			}),
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "string_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of nullable strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "string_list_prop_required",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of nullable strings, required."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema, err := oas.BuildSchema(testCase.schemaProxy)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			attributes, err := schema.BuildResourceAttributes()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(attributes, testCase.expectedAttributes); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
