// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
)

func TestParseConfig_Valid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input string
	}{
		"valid single resource": {
			input: `
provider:
  name: example

resources:
  thing:
    create:
      path: /example/path/to/things
      method: POST
    read:
      path: /example/path/to/thing/{id}
      method: GET`,
		},
		"valid resource with parameter matches": {
			input: `
provider:
  name: example

resources:
  thing:
    create:
      path: /example/path/to/things
      method: POST
    read:
      path: /example/path/to/thing/{id}
      method: GET
    merge_options:
      param_matches:
        otherId: id`,
		},
		"valid single data source": {
			input: `
provider:
  name: example

data_sources:
  thing:
    read:
      path: /example/path/to/thing/{id}
      method: GET`,
		},
		"valid data source with parameter matches": {
			input: `
provider:
  name: example

data_sources:
  thing:
    read:
      path: /example/path/to/thing/{id}
      method: GET
    merge_options:
      param_matches:
        otherId: id`,
		},
		"valid combo of resources and data sources": {
			input: `
provider:
  name: example

resources:
  thing_one:
    create:
      path: /example/path/to/things
      method: POST
    read:
      path: /example/path/to/thing/{id}
      method: GET
  thing_two:
    create:
      path: /example/path/to/things
      method: POST
    read:
      path: /example/path/to/thing/{id}
      method: GET
    update:
      path: /example/path/to/thing/{id}
      method: PATCH
    delete:
      path: /example/path/to/thing/{id}
      method: DELETE
data_sources:
  thing_one:
    read:
      path: /example/path/to/thing/{id}
      method: GET
  thing_two:
    read:
      path: /example/path/to/thing/{id}
      method: GET`,
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := config.ParseConfig([]byte(testCase.input))
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			}
		})
	}
}

func TestParseConfig_Invalid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input            string
		expectedErrRegex string
	}{
		"invalid YAML": {
			input:            `^&*`,
			expectedErrRegex: `error unmarshaling config`,
		},
		"provider - name required": {
			input:            ``,
			expectedErrRegex: `provider must have a 'name' property`,
		},
		"provider - invalid schema_ref - not resolvable": {
			input: `
provider:
  name: example
  schema_ref: thisaintvalid

data_sources:
  thing_one:
    read:
      path: /example/path/to/thing/{id}
      method: GET`,
			expectedErrRegex: `provider 'schema_ref' must be a valid JSON schema reference`,
		},
		"at least one resource or data source required": {
			input: `
provider:
  name: example`,
			expectedErrRegex: `at least one object is required in either 'resources' or 'data_sources'`,
		},
		"resource - create required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    read:
      path: /example/path/to/thing/{id}
      method: GET`,
			expectedErrRegex: `resource 'thing_one' must have a create object`,
		},
		"resource - read required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    create:
      path: /example/path/to/things
      method: POST`,
			expectedErrRegex: `resource 'thing_one' must have a read object`,
		},
		"resource - invalid create - path required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    create:
      method: POST`,
			expectedErrRegex: `invalid create: 'path' property is required`,
		},
		"resource - invalid create - method required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    create:
      path: /example/path/to/things`,
			expectedErrRegex: `invalid create: 'method' property is required`,
		},
		"resource - invalid read - path required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    read:
      method: POST`,
			expectedErrRegex: `invalid read: 'path' property is required`,
		},
		"resource - invalid read - method required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    read:
      path: /example/path/to/things`,
			expectedErrRegex: `invalid read: 'method' property is required`,
		},
		"resource - invalid update - path required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    update:
      method: POST`,
			expectedErrRegex: `invalid update: 'path' property is required`,
		},
		"resource - invalid update - method required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    update:
      path: /example/path/to/things`,
			expectedErrRegex: `invalid update: 'method' property is required`,
		},
		"resource - invalid delete - path required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    delete:
      method: POST`,
			expectedErrRegex: `invalid delete: 'path' property is required`,
		},
		"resource - invalid delete - method required": {
			input: `
provider:
  name: example

resources:
  thing_one:
    delete:
      path: /example/path/to/things`,
			expectedErrRegex: `invalid delete: 'method' property is required`,
		},
		"data source - read required": {
			input: `
provider:
  name: example

data_sources:
  thing_one:`,
			expectedErrRegex: `data_source 'thing_one' must have a read object`,
		},
		"data source - invalid read - path required": {
			input: `
provider:
  name: example

data_sources:
  thing_one:
    read:
      method: GET`,
			expectedErrRegex: `invalid read: 'path' property is required`,
		},
		"data source - invalid read - method required": {
			input: `
provider:
  name: example

data_sources:
  thing_one:
    read:
      path: /example/path/to/thing/{id}`,
			expectedErrRegex: `invalid read: 'method' property is required`,
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		errRegex := regexp.MustCompile(testCase.expectedErrRegex)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := config.ParseConfig([]byte(testCase.input))
			if err == nil {
				t.Fatalf("Expected err to match %q, got nil", testCase.expectedErrRegex)
			}
			if !errRegex.Match([]byte(err.Error())) {
				t.Errorf("Expected error to match %q, got %q", testCase.expectedErrRegex, err.Error())
			}
		})
	}
}
