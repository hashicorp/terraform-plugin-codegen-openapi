// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper

import (
	"fmt"
	"log/slog"

	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/config"
	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/log"
	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var _ ProviderMapper = providerMapper{}

type ProviderMapper interface {
	MapToIR(*slog.Logger) (*provider.Provider, error)
}

type providerMapper struct {
	provider explorer.Provider
	document *high.Document
	//nolint:unused // Might be useful later!
	cfg config.Config
}

func NewProviderMapper(exploredProvider explorer.Provider, document *high.Document, cfg config.Config) ProviderMapper {
	return providerMapper{
		provider: exploredProvider,
		document: document,
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

	providerSchema, err := generateProviderSchema(pLogger, m.provider, m.document)
	if err != nil {
		return nil, err
	}

	providerIR.Schema = providerSchema
	return &providerIR, nil
}

func generateProviderSchema(logger *slog.Logger, exploredProvider explorer.Provider, document *high.Document) (*provider.Schema, error) {
	providerSchema := &provider.Schema{}

	schemaOpts := oas.SchemaOpts{
		Ignores: exploredProvider.Ignores,
	}
	s, err := oas.BuildSchema(exploredProvider.SchemaProxy, schemaOpts, oas.GlobalSchemaOpts{
		Document: document,
	})
	if err != nil {
		return nil, err
	}

	attributes, err := s.BuildProviderAttributes()
	if err != nil {
		log.WarnLogOnError(logger, err, "error mapping provider schema")

		return nil, fmt.Errorf("error mapping provider schema: %w", err)
	}

	providerSchema.Attributes = attributes.ToSpec()

	return providerSchema, nil
}
