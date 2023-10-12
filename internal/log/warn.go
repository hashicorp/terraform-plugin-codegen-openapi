// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package log

import (
	"errors"
	"log/slog"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
)

// WarnLogOnError inspects the error type and extracts additional information for structured logging if possible
func WarnLogOnError(logger *slog.Logger, err error, message string) {
	if err == nil {
		return
	}

	var propErr *oas.PropertyError

	if errors.As(err, &propErr) {
		logger.Warn(
			message,
			"err", err,
			"oas_property", propErr.Path(),
			"oas_line_number", propErr.LineNumber(),
		)

		return
	}

	logger.Warn(message, "err", err)
}
