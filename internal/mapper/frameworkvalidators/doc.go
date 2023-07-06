// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package frameworkvalidators contains functionality for mapping validations
// onto specification that uses terraform-plugin-framework-validators.
//
// Currently, the specification requires all schema validations to be written
// as "custom" validations. Over time, the specification may begin to support
// "native" validations for very common use cases as specific properties, which
// would alleviate some of the need of this package.
package frameworkvalidators
