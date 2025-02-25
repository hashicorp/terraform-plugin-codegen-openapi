// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd_test

import (
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/cli"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/cmd"
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
			configPath:     "testdata/github/generator_config.yml",
			goldenFilePath: "testdata/github/provider_code_spec.json",
		},
		"Swagger Petstore - OpenAPI 3.0": {
			oasSpecPath:    "testdata/petstore3/openapi_spec.json",
			configPath:     "testdata/petstore3/generator_config.yml",
			goldenFilePath: "testdata/petstore3/provider_code_spec.json",
		},
		"Scaleway - Instance API": {
			oasSpecPath:    "testdata/scaleway/openapi_spec.yml",
			configPath:     "testdata/scaleway/generator_config.yml",
			goldenFilePath: "testdata/scaleway/provider_code_spec.json",
		},
		"EdgeCase API": {
			oasSpecPath:    "testdata/edgecase/openapi_spec.yml",
			configPath:     "testdata/edgecase/generator_config.yml",
			goldenFilePath: "testdata/edgecase/provider_code_spec.json",
		},
		"Kubernetes API": {
			oasSpecPath:    "testdata/kubernetes/openapi_spec.json",
			configPath:     "testdata/kubernetes/generator_config.yml",
			goldenFilePath: "testdata/kubernetes/provider_code_spec.json",
		},
	}
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempProviderSpecPath := path.Join(t.TempDir(), "provider_code_spec.json")

			mockUi := cli.NewMockUi()
			c := cmd.GenerateCommand{UI: mockUi}
			args := []string{
				"--config", testCase.configPath,
				"--output", tempProviderSpecPath,
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

			tempProviderSpecBytes, err := os.ReadFile(tempProviderSpecPath)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tempProviderSpecBytes, goldenFileBytes); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
