// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper"
	"github.com/hashicorp/terraform-plugin-codegen-spec/spec"

	"github.com/hashicorp/cli"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/index"
)

type GenerateCommand struct {
	UI             cli.Ui
	oasInputPath   string
	flagConfigPath string
	flagOutputPath string
}

func (cmd *GenerateCommand) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.StringVar(&cmd.flagConfigPath, "config", "./generator_config.yml", "path to generator config file (YAML)")
	fs.StringVar(&cmd.flagOutputPath, "output", "./provider_code_spec.json", "destination file path for generated provider code spec (JSON)")
	return fs
}

func (cmd *GenerateCommand) Help() string {
	strBuilder := &strings.Builder{}

	longestName := 0
	longestUsage := 0
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		if len(f.Name) > longestName {
			longestName = len(f.Name)
		}
		if len(f.Usage) > longestUsage {
			longestUsage = len(f.Usage)
		}
	})

	strBuilder.WriteString("\nUsage: tfplugingen-openapi generate [<args>] </path/to/oas_file.yml>\n\n")
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		if f.DefValue != "" {
			strBuilder.WriteString(fmt.Sprintf("    --%s <ARG> %s%s%s  (default: %q)\n",
				f.Name,
				strings.Repeat(" ", longestName-len(f.Name)+2),
				f.Usage,
				strings.Repeat(" ", longestUsage-len(f.Usage)+2),
				f.DefValue,
			))
		} else {
			strBuilder.WriteString(fmt.Sprintf("    --%s <ARG> %s%s%s\n",
				f.Name,
				strings.Repeat(" ", longestName-len(f.Name)+2),
				f.Usage,
				strings.Repeat(" ", longestUsage-len(f.Usage)+2),
			))
		}
	})
	strBuilder.WriteString("\n")

	return strBuilder.String()
}

func (cmd *GenerateCommand) Synopsis() string {
	return "Generates Provider Code Specification from an OpenAPI 3.x Specification"
}

func (cmd *GenerateCommand) Run(args []string) int {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	fs := cmd.Flags()
	err := fs.Parse(args)
	if err != nil {
		logger.Error("error parsing flags", "err", err)
		return 1
	}

	cmd.oasInputPath = fs.Arg(0)
	if cmd.oasInputPath == "" {
		logger.Error("error executing command", "err", "OpenAPI specification file is required as last argument")
		return 1
	}

	err = cmd.runInternal(logger)
	if err != nil {
		logger.Error("error executing command", "err", err)
		return 1
	}

	return 0
}

func (cmd *GenerateCommand) runInternal(logger *slog.Logger) error {
	// 1. Read and parse generator config file
	configBytes, err := os.ReadFile(cmd.flagConfigPath)
	if err != nil {
		return fmt.Errorf("error reading generator config file: %w", err)
	}
	config, err := config.ParseConfig(configBytes)
	if err != nil {
		return fmt.Errorf("error parsing generator config file: %w", err)
	}

	// 2. Read and parse OpenAPI spec file
	oasBytes, err := os.ReadFile(cmd.oasInputPath)
	if err != nil {
		return fmt.Errorf("error reading OpenAPI spec file: %w", err)
	}
	doc, err := libopenapi.NewDocument(oasBytes)
	if err != nil {
		return fmt.Errorf("error parsing OpenAPI spec file: %w", err)
	}

	// 3. Build out the OpenAPI model, this will recursively load all local + remote references into one cohesive model
	model, errs := doc.BuildV3Model()

	// 4. Log circular references as warnings and fail on any other model building errors
	var errResult error
	for _, err := range errs {
		if rslvErr, ok := err.(*index.ResolvingError); ok {
			logger.Warn(
				"circular reference found in OpenAPI spec",
				"circular_ref", rslvErr.CircularReference.GenerateJourneyPath())
			continue
		}

		errResult = errors.Join(errResult, err)
	}
	if errResult != nil {
		return fmt.Errorf("error building OpenAPI 3.x model: %w", errResult)
	}

	// 5. Generate provider code spec w/ config
	oasExplorer := explorer.NewConfigExplorer(model.Model, *config)
	providerCodeSpec, err := generateProviderCodeSpec(logger, oasExplorer, *config)
	if err != nil {
		return err
	}

	// 6. Use provider code spec to create JSON
	bytes, err := json.MarshalIndent(providerCodeSpec, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshalling provider code spec to JSON: %w", err)
	}

	// 7. Log a warning if the provider code spec is not valid based on the JSON schema
	err = spec.Validate(context.TODO(), bytes)
	if err != nil {
		logger.Warn(
			"generated provider code spec failed validation",
			"validation_msg", err)
	}

	// 8. Output to file
	output, err := os.Create(cmd.flagOutputPath)
	if err != nil {
		return fmt.Errorf("error creating output file for provider code spec: %w", err)
	}

	_, err = output.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing provider code spec to output: %w", err)
	}

	return nil
}

func generateProviderCodeSpec(logger *slog.Logger, dora explorer.Explorer, cfg config.Config) (*spec.Specification, error) {
	// 1. Find TF resources in OAS
	explorerResources, err := dora.FindResources()
	if err != nil {
		return nil, fmt.Errorf("error finding resource(s): %w", err)
	}

	// 2. Find TF data sources in OAS
	explorerDataSources, err := dora.FindDataSources()
	if err != nil {
		return nil, fmt.Errorf("error finding data source(s): %w", err)
	}

	// 3. Find TF provider in OAS
	explorerProvider, err := dora.FindProvider()
	if err != nil {
		return nil, fmt.Errorf("error finding provider: %w", err)
	}

	// 4. Use TF info to generate provider code spec for resources
	resourceMapper := mapper.NewResourceMapper(explorerResources, cfg)
	resourcesIR, err := resourceMapper.MapToIR(logger)
	if err != nil {
		return nil, fmt.Errorf("error generating provider code spec for resources: %w", err)
	}

	// 5. Use TF info to generate provider code spec for data sources
	dataSourceMapper := mapper.NewDataSourceMapper(explorerDataSources, cfg)
	dataSourcesIR, err := dataSourceMapper.MapToIR(logger)
	if err != nil {
		return nil, fmt.Errorf("error generating provider code spec for data sources: %w", err)
	}

	// 6. Use TF info to generate provider code spec for provider
	providerMapper := mapper.NewProviderMapper(explorerProvider, cfg)
	providerIR, err := providerMapper.MapToIR(logger)
	if err != nil {
		return nil, fmt.Errorf("error generating provider code spec for provider: %w", err)
	}

	requestMapper := mapper.NewRequestMapper(explorerResources, cfg)
	requestsIR, err := requestMapper.MapToIR(logger)
	if err != nil {
		return nil, fmt.Errorf("error generating provider code spec for request: %w", err)
	}

	return &spec.Specification{
		Version:     spec.Version0_1,
		Provider:    providerIR,
		Resources:   resourcesIR,
		DataSources: dataSourcesIR,
		Requests:    requestsIR,
	}, nil
}
