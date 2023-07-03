// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frameworkvalidators

import "github.com/hashicorp/terraform-plugin-codegen-spec/code"

const (
	// CodeImportBasePath is the base code import path for framework validators.
	CodeImportBasePath = "github.com/hashicorp/terraform-plugin-framework-validators"
)

// CodeImport returns the framework validators code import for the given path.
func CodeImport(packagePath string) code.Import {
	return code.Import{
		Path: CodeImportBasePath + "/" + packagePath,
	}
}
