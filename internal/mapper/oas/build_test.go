// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/oas"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

func TestBuildSchemaFromRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		op             *high.Operation
		expectedSchema *oas.OASSchema
	}{
		"default to application/json": {
			op: &high.Operation{
				RequestBody: &high.RequestBody{
					Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
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
					}),
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
					Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
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
					}),
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
					Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
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
					}),
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

			got, err := oas.BuildSchemaFromRequest(testCase.op, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

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
					Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
						"application/json": {
							Schema: nil,
						},
						"application/xml": {
							Schema: nil,
						},
					}),
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

			_, err := oas.BuildSchemaFromRequest(testCase.op, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})

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
					Codes: orderedmap.ToOrderedMap(map[string]*high.Response{
						"201": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the wrong one!",
										Type:        []string{"boolean"},
									}),
								},
							}),
						},
						"200": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
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
							}),
						},
					}),
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
					Codes: orderedmap.ToOrderedMap(map[string]*high.Response{
						"204": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the wrong one!",
										Type:        []string{"boolean"},
									}),
								},
							}),
						},
						"201": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
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
							}),
						},
					}),
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
					Codes: orderedmap.ToOrderedMap(map[string]*high.Response{
						"304": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this is the wrong one!",
										Type:        []string{"boolean"},
									}),
								},
							}),
						},
						"204": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
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
							}),
						},
					}),
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

			got, err := oas.BuildSchemaFromResponse(testCase.op, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

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
					Codes: orderedmap.ToOrderedMap(map[string]*high.Response{
						"300": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this won't be used!",
										Type:        []string{"string"},
									}),
								},
							}),
						},
						"skip-me!": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this won't be used!",
										Type:        []string{"string"},
									}),
								},
							}),
						},
						"199": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
								"application/json": {
									Schema: base.CreateSchemaProxy(&base.Schema{
										Description: "this won't be used!",
										Type:        []string{"string"},
									}),
								},
							}),
						},
					}),
				},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"200 response code with no valid schema": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: orderedmap.ToOrderedMap(map[string]*high.Response{
						"200": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
								"application/json": {
									Schema: nil,
								},
							}),
						},
					}),
				},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"201 response code with no valid schema": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: orderedmap.ToOrderedMap(map[string]*high.Response{
						"201": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
								"application/json": {
									Schema: nil,
								},
							}),
						},
					}),
				},
			},
			expectedErrRegex: oas.ErrSchemaNotFound.Error(),
		},
		"success response code with no valid schema": {
			op: &high.Operation{
				Responses: &high.Responses{
					Codes: orderedmap.ToOrderedMap(map[string]*high.Response{
						"204": {
							Content: orderedmap.ToOrderedMap(map[string]*high.MediaType{
								"application/json": {
									Schema: nil,
								},
							}),
						},
					}),
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

			_, err := oas.BuildSchemaFromResponse(testCase.op, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})

			if err == nil {
				t.Fatalf("Expected err to match %q, got nil", testCase.expectedErrRegex)
			}
			if !errRegex.Match([]byte(err.Error())) {
				t.Errorf("Expected error to match %q, got %q", testCase.expectedErrRegex, err.Error())
			}
		})
	}
}

func TestBuildSchema_MultiTypes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schemaProxy        *base.SchemaProxy
		expectedAttributes attrmapper.ResourceAttributes
	}{
		"nullable type - Type array": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"nullable_string_two"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nullable_string_one": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"null", "string"},
						Description: "hey there! I'm a nullable string type.",
					}),
					"nullable_string_two": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string", "null"},
						Description: "hey there! I'm a nullable string type, required.",
					}),
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "nullable_string_one",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a nullable string type."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "nullable_string_two",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a nullable string type, required."),
					},
				},
			},
		},
		"stringable types - Type array": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"stringable_number", "stringable_bool"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"stringable_bool": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string", "boolean"},
						Description: "hey there! I'm a stringable bool type, required.",
					}),
					"stringable_integer": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"integer", "string"},
						Description: "hey there! I'm a stringable integer type.",
					}),
					"stringable_number": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string", "number"},
						Description: "hey there! I'm a stringable number type, required.",
					}),
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "stringable_bool",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a stringable bool type, required."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "stringable_integer",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a stringable integer type."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "stringable_number",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a stringable number type, required."),
					},
				},
			},
		},
		"nullable type - anyOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"nullable_string_two"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "nullable_string_one",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a string type."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "nullable_string_two",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a string type, required."),
					},
				},
			},
		},
		"stringable types - anyOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"stringable_number", "stringable_bool"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"stringable_bool": base.CreateSchemaProxy(&base.Schema{
						AnyOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a stringable bool type, required.",
							}),
						},
					}),
					"stringable_integer": base.CreateSchemaProxy(&base.Schema{
						AnyOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a stringable integer type.",
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"integer"},
							}),
						},
					}),
					"stringable_number": base.CreateSchemaProxy(&base.Schema{
						AnyOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a stringable number type, required.",
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"number"},
							}),
						},
					}),
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "stringable_bool",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a stringable bool type, required."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "stringable_integer",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a stringable integer type."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "stringable_number",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a stringable number type, required."),
					},
				},
			},
		},
		"nullable type - oneOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"nullable_string_two"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "nullable_string_one",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a string type."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "nullable_string_two",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a string type, required."),
					},
				},
			},
		},
		"stringable types - oneOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"stringable_number", "stringable_bool"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"stringable_bool": base.CreateSchemaProxy(&base.Schema{
						OneOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a stringable bool type, required.",
							}),
						},
					}),
					"stringable_integer": base.CreateSchemaProxy(&base.Schema{
						OneOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a stringable integer type.",
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"integer"},
							}),
						},
					}),
					"stringable_number": base.CreateSchemaProxy(&base.Schema{
						OneOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a stringable number type, required.",
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"number"},
							}),
						},
					}),
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "stringable_bool",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a stringable bool type, required."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "stringable_integer",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a stringable integer type."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "stringable_number",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a stringable number type, required."),
					},
				},
			},
		},
		"random types - oneOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				OneOf: []*base.SchemaProxy{
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"object"},
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"field": base.CreateSchemaProxy(&base.Schema{
								Type: []string{"string"},
							}),
						}),
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type:  []string{"string"},
						Title: "Cool Type!",
					}),
				},
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					StringAttribute: resource.StringAttribute{ComputedOptionalRequired: "computed_optional"},
					Name:            "cool_type",
				},
				&attrmapper.ResourceSingleNestedAttribute{
					SingleNestedAttribute: resource.SingleNestedAttribute{ComputedOptionalRequired: "computed_optional"},
					Name:                  "field_0",
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
							},
							Name: "field",
						},
					},
				},
			},
		},
		"list attributes with nullable element type - Type array": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListAttribute{
					Name: "string_list_prop",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of nullable strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				&attrmapper.ResourceListAttribute{
					Name: "string_list_prop_required",
					ListAttribute: resource.ListAttribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListAttribute{
					Name: "string_list_prop",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of nullable strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				&attrmapper.ResourceListAttribute{
					Name: "string_list_prop_required",
					ListAttribute: resource.ListAttribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListAttribute{
					Name: "string_list_prop",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of nullable strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				&attrmapper.ResourceListAttribute{
					Name: "string_list_prop_required",
					ListAttribute: resource.ListAttribute{
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

			schema, err := oas.BuildSchema(testCase.schemaProxy, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})
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

func TestBuildSchema_AllOfSchemaComposition(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schemaProxy        *base.SchemaProxy
		expectedAttributes attrmapper.ResourceAttributes
	}{
		"allOf with one element - use subschema": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				AllOf: []*base.SchemaProxy{
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"object"},
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"nested_object": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"string"},
								Description: "hey there! I'm an object type.",
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"bool": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"boolean"},
										Description: "hey there! I'm a bool type.",
									}),
									"string": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"string"},
										Description: "hey there! I'm a string type, required.",
									}),
								}),
							}),
						}),
					}),
				},
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSingleNestedAttribute{
					Name: "nested_object",
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceBoolAttribute{
							Name: "bool",
							BoolAttribute: resource.BoolAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("hey there! I'm a bool type."),
							},
						},
						&attrmapper.ResourceStringAttribute{
							Name: "string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("hey there! I'm a string type, required."),
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm an object type."),
					},
				},
			},
		},
		"allOf with one element - use subschema and override description": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_allof_override"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"bool_allof_override": base.CreateSchemaProxy(&base.Schema{
						Description: "Override the bool's description",
						AllOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"boolean"},
								Description: "hey there! I'm a bool type.",
							}),
						},
					}),
					"string_allof_override": base.CreateSchemaProxy(&base.Schema{
						Description: "Override the string's description",
						AllOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a string type.",
							}),
						},
					}),
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceBoolAttribute{
					Name: "bool_allof_override",
					BoolAttribute: resource.BoolAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("Override the bool's description"),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "string_allof_override",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("Override the string's description"),
					},
				},
			},
		},
		"allOf with object composition": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				AllOf: []*base.SchemaProxy{
					base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"object"},
						Required: []string{"string_allof_override"},
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"field1": base.CreateSchemaProxy(&base.Schema{
								Description: "I m field1",
								Type:        []string{"string"},
							}),
						}),
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"object"},
						Required: []string{"string_allof_override"},
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"field2": base.CreateSchemaProxy(&base.Schema{
								Description: "I m field2",
								Type:        []string{"string"},
							}),
						}),
					}),
				},
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "field1",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("I m field1"),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "field2",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("I m field2"),
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema, err := oas.BuildSchema(testCase.schemaProxy, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})
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

func TestBuildSchema_Errors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schemaProxy      *base.SchemaProxy
		expectedErrRegex string
	}{
		"no type or schema composition": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type: []string{},
			}),
			expectedErrRegex: `no 'type' array or supported allOf, oneOf, anyOf constraint - attribute cannot be created`,
		},
		"unsupported multi-type array": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"string", "object"},
			}),
			expectedErrRegex: `\[string object\] - unsupported multi-type, attribute cannot be created`,
		},
		"unsupported multi-type oneOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				OneOf: []*base.SchemaProxy{
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"object"},
					}),
				},
			}),
			expectedErrRegex: `\[string object\] - unsupported multi-type, attribute cannot be created`,
		},
		"unsupported multi-type anyOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				OneOf: []*base.SchemaProxy{
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"object"},
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
			}),
			expectedErrRegex: `\[object string\] - unsupported multi-type, attribute cannot be created`,
		},
		"too many allOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				AllOf: []*base.SchemaProxy{
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"null"},
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
			}),
			expectedErrRegex: `found 2 allOf subschema\(s\), schema composition is currently not supported`,
		},
		"too many anyOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				AnyOf: []*base.SchemaProxy{
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"null"},
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"integer"},
					}),
				},
			}),
			expectedErrRegex: `found 3 anyOf subschema\(s\), schema composition is currently not supported`,
		},
		"too many oneOf": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				OneOf: []*base.SchemaProxy{
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"null"},
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"integer"},
					}),
					base.CreateSchemaProxy(&base.Schema{
						Type: []string{"number"},
					}),
				},
			}),
			expectedErrRegex: `found 4 oneOf subschema\(s\), schema composition is currently not supported`,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase
		errRegex := regexp.MustCompile(testCase.expectedErrRegex)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := oas.BuildSchema(testCase.schemaProxy, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})
			if err == nil {
				t.Fatalf("Expected err to match %q, got nil", testCase.expectedErrRegex)
			}

			if !errRegex.Match([]byte(err.Error())) {
				t.Errorf("Expected error to match %q, got %q", testCase.expectedErrRegex, err.Error())
			}
		})
	}
}

func TestBuildSchema_EdgeCases(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schemaProxy        *base.SchemaProxy
		expectedAttributes attrmapper.ResourceAttributes
	}{
		"no type with properties - defaults to object": {
			schemaProxy: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{},
				Required: []string{"string", "object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Description: "hey there! I'm a string type, required.",
					}),
					"object": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{},
						Required:    []string{"bool"},
						Description: "hey there! I'm an object type, required.",
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"bool": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"boolean"},
								Description: "hey there! I'm a bool type, required.",
							}),
						}),
					}),
				}),
			}),
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSingleNestedAttribute{
					Name: "object",
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceBoolAttribute{
							Name: "bool",
							BoolAttribute: resource.BoolAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("hey there! I'm a bool type, required."),
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm an object type, required."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "string",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a string type, required."),
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema, err := oas.BuildSchema(testCase.schemaProxy, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})
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
