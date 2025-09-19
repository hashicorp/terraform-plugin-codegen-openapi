// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package log

import (
	"errors"
	"log/slog"

	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/oas"
)

// WarnLogOnError inspects the error type and extracts additional information for structured logging if possible
func WarnLogOnError(logger *slog.Logger, err error, message string) {
	if err == nil {
		return
	}

	var schemaErr *oas.SchemaError
	if errors.As(err, &schemaErr) {
		if schemaErr.Path() != "" {
			logger = logger.With("oas_path", schemaErr.Path())
		}
		if schemaErr.LineNumber() != 0 {
			logger = logger.With("oas_line_number", schemaErr.LineNumber())
		}
	}

	logger.Warn(message, "err", err)
}
