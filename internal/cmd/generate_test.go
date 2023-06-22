// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/cmd"
	"github.com/hashicorp/terraform-plugin-codegen-spec/spec"
	"github.com/mitchellh/cli"
)

func TestGenerate_WithConfig(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		oasSpecPath    string
		configPath     string
		goldenFilePath string
	}{
		"GitHub v3 REST API": {
			oasSpecPath:    "testdata/github/openapi_spec.json",
			configPath:     "testdata/github/tfopenapigen_config.yml",
			goldenFilePath: "testdata/github/generated_framework_ir.json",
		},
		"Swagger Petstore - OpenAPI 3.0": {
			oasSpecPath:    "testdata/petstore3/openapi_spec.json",
			configPath:     "testdata/petstore3/tfopenapigen_config.yml",
			goldenFilePath: "testdata/petstore3/generated_framework_ir.json",
		},
		"Scaleway - Instance API": {
			oasSpecPath:    "testdata/scaleway/openapi_spec.yml",
			configPath:     "testdata/scaleway/tfopenapigen_config.yml",
			goldenFilePath: "testdata/scaleway/generated_framework_ir.json",
		},
		"Edgecase API": {
			oasSpecPath:    "testdata/edgecase/openapi_spec.yml",
			configPath:     "testdata/edgecase/tfopenapigen_config.yml",
			goldenFilePath: "testdata/edgecase/generated_framework_ir.json",
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockUi := cli.NewMockUi()
			c := cmd.GenerateCommand{UI: mockUi}
			args := []string{
				"--config", testCase.configPath,
				testCase.oasSpecPath,
			}

			exitCode := c.Run(args)
			if exitCode != 0 {
				t.Fatalf("unexpected error running generate cmd: %s", mockUi.ErrorWriter.String())
			}

			goldenFileBytes, err := os.ReadFile(testCase.goldenFilePath)
			if err != nil {
				t.Fatal(err)
			}

			// TODO: write logic to allow updating of golden files?
			if diff := cmp.Diff(mockUi.OutputWriter.Bytes(), goldenFileBytes, getFrameworkIRCmpOption()); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

// getFrameworkIRCmpOption returns a go-cmp transformer that will unmarshal JSON into Framework IR structs and compare.
// This produces a nicer diff for comparing golden files and will error with any unknown fields. Based off: https://github.com/google/go-cmp/issues/224#issuecomment-650429859
func getFrameworkIRCmpOption() cmp.Option {
	return cmp.FilterValues(func(x, y []byte) bool {
		return json.Valid(x) && json.Valid(y)
	}, cmp.Transformer("ParseIRJSON", func(in []byte) (out interface{}) {
		var irStruct spec.Specification
		decoder := json.NewDecoder(strings.NewReader(string(in)))
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&irStruct); err != nil {
			panic(fmt.Errorf("error parsing Framework IR JSON bytes: %w", err))
		}
		return irStruct
	}))
}
