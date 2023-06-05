// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import "sort"

// Generics? ☜(ಠ_ಠ☜)
func SortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0)

	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
