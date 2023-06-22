// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/cmd"
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

			if diff := cmp.Diff(mockUi.OutputWriter.Bytes(), goldenFileBytes); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
