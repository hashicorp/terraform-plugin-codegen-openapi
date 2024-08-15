// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

// This regex matches attribute locations, dot-separated, as represented as {attribute_name}.{nested_attribute_name}
//   - category = MATCH
//   - category.id = MATCH
//   - category.tags.name = MATCH
//   - category. = NO MATCH
//   - .category = NO MATCH
var attributeLocationRegex = regexp.MustCompile(`^[\w]+(?:\.[\w]+)*$`)

// Config represents a YAML generator config.
type Config struct {
	Provider    Provider              `yaml:"provider"`
	Resources   map[string]Resource   `yaml:"resources"`
	DataSources map[string]DataSource `yaml:"data_sources"`
}

// Provider generator config section.
type Provider struct {
	Name      string `yaml:"name"`
	SchemaRef string `yaml:"schema_ref"`

	// TODO: At some point, this should probably be refactored to work with the SchemaOptions struct
	// Ignores are a slice of strings, representing an attribute location to ignore during mapping (dot-separated for nested attributes).
	Ignores []string `yaml:"ignores"`
}

// Resource generator config section.
type Resource struct {
	Create        *OpenApiSpecLocation `yaml:"create"`
	Read          *OpenApiSpecLocation `yaml:"read"`
	Update        *OpenApiSpecLocation `yaml:"update"`
	Delete        *OpenApiSpecLocation `yaml:"delete"`
	SchemaOptions SchemaOptions        `yaml:"schema"`
}

// DataSource generator config section.
type DataSource struct {
	Read          *OpenApiSpecLocation `yaml:"read"`
	SchemaOptions SchemaOptions        `yaml:"schema"`
}

// OpenApiSpecLocation defines a location in an OpenAPI spec for an API operation.
type OpenApiSpecLocation struct {
	// Matches the path key for a path item (refer to [OAS Paths Object]).
	//
	// [OAS Paths Object]: https://spec.openapis.org/oas/v3.1.0#paths-object
	Path string `yaml:"path"`
	// Matches the operation method in a path item: GET, POST, etc (refer to [OAS Path Item Object]).
	//
	// [OAS Path Item Object]: https://spec.openapis.org/oas/v3.1.0#pathItemObject
	Method string `yaml:"method"`
}

// SchemaOptions generator config section. This section contains options for modifying the output of the generator.
type SchemaOptions struct {
	// Ignores are a slice of strings, representing an attribute location to ignore during mapping (dot-separated for nested attributes).
	Ignores          []string         `yaml:"ignores"`
	AttributeOptions AttributeOptions `yaml:"attributes"`
}

// AttributeOptions generator config section. This section is used to modify the output of specific attributes.
type AttributeOptions struct {
	// Aliases are a map, with the key being a parameter name in an OpenAPI operation and the value being the new name (alias).
	Aliases map[string]string `yaml:"aliases"`
	// Overrides are a map, with the key being an attribute location (dot-separated for nested attributes) and the value being overrides to apply to the attribute.
	Overrides map[string]Override `yaml:"overrides"`
}

// Override generator config section.
type Override struct {
	// Description overrides the description that was mapped/merged from the OpenAPI specification.
	Description string `yaml:"description"`
	// ComputedOptionalRequired overrides the inferred value from the OpenAPI specification.
	ComputedOptionalRequired string `yaml:"computed_optional_required"`
}

// ParseConfig takes in a byte array (of YAML), unmarshals into a Config struct, and validates the result
func ParseConfig(bytes []byte) (*Config, error) {
	var result Config
	err := yaml.Unmarshal(bytes, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err = result.Validate(); err != nil {
		return nil, fmt.Errorf("config validation error(s):\n%w", err)
	}

	return &result, nil
}

func (c Config) Validate() error {
	var result error

	if len(c.DataSources) == 0 && len(c.Resources) == 0 {
		result = errors.Join(result, errors.New("\tat least one object is required in either 'resources' or 'data_sources'"))
	}

	// Validate Provider
	err := c.Provider.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("\tprovider %w", err))
	}

	// Validate all Resources
	for name, resource := range c.Resources {
		err := resource.Validate()
		if err != nil {
			result = errors.Join(result, fmt.Errorf("\tresource '%s' %w", name, err))
		}
	}

	// Validate all Data Sources
	for name, dataSource := range c.DataSources {
		err := dataSource.Validate()
		if err != nil {
			result = errors.Join(result, fmt.Errorf("\tdata_source '%s' %w", name, err))
		}
	}

	return result
}

func (p Provider) Validate() error {
	var result error

	if p.Name == "" {
		result = errors.Join(result, errors.New("must have a 'name' property"))
	}

	for _, ignore := range p.Ignores {
		if !attributeLocationRegex.MatchString(ignore) {
			result = errors.Join(result, fmt.Errorf("invalid item for ignores: %q - must be dot-separated string", ignore))
		}
	}

	return result
}

func (r Resource) Validate() error {
	var result error

	if r.Create == nil {
		result = errors.Join(result, errors.New("must have a create object"))
	}
	if r.Read == nil {
		result = errors.Join(result, errors.New("must have a read object"))
	}

	err := r.Create.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid create: %w", err))
	}

	err = r.Read.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid read: %w", err))
	}

	err = r.Update.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid update: %w", err))
	}

	err = r.Delete.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid delete: %w", err))
	}

	err = r.SchemaOptions.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid schema: %w", err))
	}

	return result
}

func (d DataSource) Validate() error {
	var result error

	if d.Read == nil {
		result = errors.Join(result, errors.New("must have a read object"))
	}

	err := d.Read.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid read: %w", err))
	}

	err = d.SchemaOptions.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid schema: %w", err))
	}

	return result
}

func (o *OpenApiSpecLocation) Validate() error {
	var result error
	if o == nil {
		return nil
	}

	if o.Path == "" {
		result = errors.Join(result, errors.New("'path' property is required"))
	}

	if o.Method == "" {
		result = errors.Join(result, errors.New("'method' property is required"))
	}

	return result
}

func (s *SchemaOptions) Validate() error {
	var result error

	err := s.AttributeOptions.Validate()
	if err != nil {
		result = errors.Join(result, fmt.Errorf("invalid attributes: %w", err))
	}

	for _, ignore := range s.Ignores {
		if !attributeLocationRegex.MatchString(ignore) {
			result = errors.Join(result, fmt.Errorf("invalid item for ignores: %q - must be dot-separated string", ignore))
		}
	}

	return result
}

func (s *AttributeOptions) Validate() error {
	var result error

	for path := range s.Overrides {
		if !attributeLocationRegex.MatchString(path) {
			result = errors.Join(result, fmt.Errorf("invalid key for override: %q - must be dot-separated string", path))
		}
	}

	return result
}
