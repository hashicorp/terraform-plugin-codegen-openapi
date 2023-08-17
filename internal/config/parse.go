// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"errors"
	"fmt"

	"github.com/pb33f/libopenapi/index"
	"gopkg.in/yaml.v3"
)

// Config is tagged with `yaml` struct tags, however YAML is a superset of JSON, so it can also be parsed from JSON
type Config struct {
	Provider    Provider              `yaml:"provider"`
	Resources   map[string]Resource   `yaml:"resources"`
	DataSources map[string]DataSource `yaml:"data_sources"`
}

type Provider struct {
	Name      string `yaml:"name"`
	SchemaRef string `yaml:"schema_ref"`
}

type Resource struct {
	Create        *OpenApiSpecLocation `yaml:"create"`
	Read          *OpenApiSpecLocation `yaml:"read"`
	Update        *OpenApiSpecLocation `yaml:"update"`
	Delete        *OpenApiSpecLocation `yaml:"delete"`
	SchemaOptions SchemaOptions        `yaml:"schema"`
}
type DataSource struct {
	Read          *OpenApiSpecLocation `yaml:"read"`
	SchemaOptions SchemaOptions        `yaml:"schema"`
}

type OpenApiSpecLocation struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

type SchemaOptions struct {
	AttributeOptions AttributeOptions `yaml:"attributes"`
}
type AttributeOptions struct {
	Aliases map[string]string `yaml:"aliases"`
}

// ParseConfig takes in a byte array (of YAML), unmarshal into a Config struct, and validates the result
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
	// TODO: Add regex/validation of provider name
	if p.Name == "" {
		return errors.New("must have a 'name' property")
	}

	// All schema refs must be a local, file, or http resolve type
	if p.SchemaRef != "" && index.DetermineReferenceResolveType(p.SchemaRef) < 0 {
		return errors.New("'schema_ref' must be a valid JSON schema reference")
	}

	return nil
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
