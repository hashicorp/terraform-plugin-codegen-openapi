package mapper

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
)

var _ ProviderMapper = providerMapper{}

type ProviderMapper interface {
	MapToIR() (*provider.Provider, error)
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

func (m providerMapper) MapToIR() (*provider.Provider, error) {
	providerIR := provider.Provider{
		Name: m.provider.Name,
	}

	if m.provider.SchemaProxy == nil {
		return &providerIR, nil
	}

	providerSchema, err := generateProviderSchema(m.provider)
	if err != nil {
		return nil, fmt.Errorf("error mapping provider schema: %w", err)
	}

	providerIR.Schema = providerSchema
	return &providerIR, nil
}

func generateProviderSchema(exploredProvider explorer.Provider) (*provider.Schema, error) {
	providerSchema := &provider.Schema{}

	s, err := oas.BuildSchema(exploredProvider.SchemaProxy, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})
	if err != nil {
		return nil, err
	}

	attributes, err := s.BuildProviderAttributes()
	if err != nil {
		return nil, err
	}

	providerSchema.Attributes = *attributes

	return providerSchema, nil
}
