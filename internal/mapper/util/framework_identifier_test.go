// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
)

func TestFrameworkIdentifier(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		original string
		want     string
	}{
		"no change - empty string": {
			original: "",
			want:     "",
		},
		"no change - lowercase alphabet": {
			original: "thing",
			want:     "thing",
		},
		"no change - one letter lowercase": {
			original: "f",
			want:     "f",
		},
		"no change - leading underscore": {
			original: "_thing",
			want:     "_thing",
		},
		"no change - middle underscore": {
			original: "fake_thing",
			want:     "fake_thing",
		},
		"no change - alphanumeric": {
			original: "thing123",
			want:     "thing123",
		},
		"no change - alphanumeric with underscores": {
			original: "fake_thing_123",
			want:     "fake_thing_123",
		},
		"change - middle hyphen removed": {
			original: "fake-thing",
			want:     "fakething",
		},
		"change - middle hyphen removed and separated": {
			original: "Fake-Thing",
			want:     "fake_thing",
		},
		"change - special symbols": {
			original: "<fakeThing>",
			want:     "fake_thing",
		},
		"change - remove leading number": {
			original: "123fakeThing",
			want:     "fake_thing",
		},
		"change - one letter uppercase": {
			original: "F",
			want:     "f",
		},
		"change - capitalized with underscore": {
			original: "Fake_Thing",
			want:     "fake_thing",
		},
		"change - lower camelCase": {
			original: "fakeThing",
			want:     "fake_thing",
		},
		"change - PascalCase": {
			original: "FakeThing",
			want:     "fake_thing",
		},
		"change - lower camelCase with initialism": {
			original: "fakeID",
			want:     "fake_id",
		},
		"change - PascalCase with initialism": {
			original: "FakeURL",
			want:     "fake_url",
		},
		"change - all uppercase": {
			original: "FAKETHING",
			want:     "fakething",
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := util.TerraformIdentifier(testCase.original)
			if got != testCase.want {
				t.Fatalf("expected %s, got %s", testCase.want, got)
			}
		})
	}
}
