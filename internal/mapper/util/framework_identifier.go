// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// lowerToUpper will match a sequence of a lowercase letter followed by an uppercase letter
	lowerToUpper = regexp.MustCompile(`([a-z])[A-Z]`)

	// unsupportedCharacters matches any characters that are NOT alphanumeric or underscores
	unsupportedCharacters = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

	// leadingNumbers matches all numbers that are at the beginning of a string
	leadingNumbers = regexp.MustCompile(`^(\d+)`)
)

// TerraformIdentifier attempts to convert the given string to a valid Terraform identifier for usage in a Provider Code Specification.
func TerraformIdentifier(original string) string {
	if len(original) == 0 {
		return original
	}

	// Remove any characters that are either not supported in a Terraform indentifier, or can't be automatically converted
	removedUnsupported := unsupportedCharacters.ReplaceAllString(original, "")

	// Remove leading numbers
	noLeadingNumbers := leadingNumbers.ReplaceAllString(removedUnsupported, "")

	// Insert an underscore between lowercase letter followed by uppercase letter
	insertedUnderscores := lowerToUpper.ReplaceAllStringFunc(noLeadingNumbers, func(s string) string {
		firstRune, size := utf8.DecodeRuneInString(s)
		if firstRune == utf8.RuneError && size <= 1 {
			// The string is empty, return it
			return s
		}

		return fmt.Sprintf("%s_%s", string(firstRune), strings.ToLower(s[size:]))
	})

	// Lowercase the final string
	return strings.ToLower(insertedUnderscores)
}
