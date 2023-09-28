// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package attrmapper contains types and methods that provide an intermediate step between the OpenAPI schema
// types (libopenapi) and the Provider Code Specification types (terraform-plugin-codegen-spec). This intermediate
// step enables merging of attributes, overriding of specific properties, and converting into a provider code spec type
// to be marshalled to JSON.
package attrmapper
