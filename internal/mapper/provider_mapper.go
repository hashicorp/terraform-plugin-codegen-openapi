// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper

import (
	"fmt"
	"log/slog"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
)

var _ ProviderMapper = providerMapper{}

type ProviderMapper interface {
	MapToIR(*slog.Logger) (*provider.Provider, error)
}

type providerMapper struct {
	provider explorer.Provider
	//nolint:unused // Might be useful later!
	cfg config.Config
}

func NewProviderMapper(exploredProvider explorer.Provider, cfg config.Config) ProviderMapper {
	return providerMapper{
		provider: exploredProvider,
		cfg:      cfg,
	}
}

func (m providerMapper) MapToIR(logger *slog.Logger) (*provider.Provider, error) {
	providerIR := provider.Provider{
		Name: m.provider.Name,
	}

	if m.provider.SchemaProxy == nil {
		return &providerIR, nil
	}

	pLogger := logger.With("provider", providerIR.Name)

	providerSchema, err := generateProviderSchema(pLogger, m.provider)
	if err != nil {
		return nil, err
	}

	providerIR.Schema = providerSchema
	return &providerIR, nil
}

func generateProviderSchema(logger *slog.Logger, exploredProvider explorer.Provider) (*provider.Schema, error) {
	providerSchema := &provider.Schema{}

	s, err := oas.BuildSchema(exploredProvider.SchemaProxy, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})
	if err != nil {
		return nil, err
	}

	attributes, propErr := s.BuildProviderAttributes()
	if propErr != nil {
		logger.Error("error mapping provider schema", "err", propErr,
			"oas_property", propErr.Path(),
			"oas_line_number", propErr.LineNumber())

		return nil, fmt.Errorf("error mapping provider schema: %w", propErr)
	}

	providerSchema.Attributes = attributes.ToSpec()

	return providerSchema, nil
}
