package mapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
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
	return &provider.Provider{
		Name: m.provider.Name,
	}, nil
}
