package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/ir"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/resource"

	"github.com/mitchellh/cli"
	"github.com/pb33f/libopenapi"
)

type generateCmd struct {
	ui             cli.Ui
	oasInputPath   string
	flagConfigPath string
	flagOutputPath string
}

func (cmd *generateCmd) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.StringVar(&cmd.flagConfigPath, "config", "./tfopenapigen_config.yml", "path to config file (YAML)")
	fs.StringVar(&cmd.flagOutputPath, "output", "", "path to output generated Framework IR file (JSON)")
	return fs
}

func (cmd *generateCmd) Help() string {
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

	strBuilder.WriteString("\nUsage: tfopenapigen generate [<args>] </path/to/oas_file.yml>\n\n")
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

func (cmd *generateCmd) Synopsis() string {
	return "Generates Framework Intermediate Representation (IR) JSON for an OpenAPI spec (JSON or YAML format)"
}

func (cmd *generateCmd) Run(args []string) int {
	fs := cmd.Flags()
	err := fs.Parse(args)
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("unable to parse flags: %s", err))
		return 1
	}

	cmd.oasInputPath = fs.Arg(0)
	if cmd.oasInputPath == "" {
		cmd.ui.Error("Error executing command: OpenAPI specification file is required as last argument")
		return 1
	}

	err = cmd.runInternal()
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("Error executing command: %s\n", err))
		return 1
	}

	return 0
}

func (cmd *generateCmd) runInternal() error {
	// 1. Read and parse generator config file
	configBytes, err := os.ReadFile(cmd.flagConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read generator config file: %w", err)
	}
	config, err := config.ParseConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to parse generator config file: %w", err)
	}

	// 2. Read and parse OpenAPI spec file
	oasBytes, err := os.ReadFile(cmd.oasInputPath)
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI spec file: %w", err)
	}
	doc, err := libopenapi.NewDocument(oasBytes)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec file: %w", err)
	}

	// 3. Build out the OpenAPI model, this will recursively load all local + remote references into one cohesive model
	model, errs := doc.BuildV3Model()
	// TODO: Determine how to handle circular ref errors - https://pb33f.io/libopenapi/circular-references/
	if len(errs) > 0 {
		var errResult error
		for _, err := range errs {
			errResult = errors.Join(errResult, err)
		}
		log.Printf("[WARN] Potential issues in model spec: %s", errResult)
	}

	// 4. Generate framework IR w/ config
	oasExplorer := explorer.NewConfigExplorer(model.Model, *config)
	frameworkIr, err := generateFrameworkIr(oasExplorer, *config)
	if err != nil {
		return err
	}

	// 5. Use framework IR to create JSON
	bytes, err := json.MarshalIndent(frameworkIr, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshalling Framework IR to JSON: %w", err)
	}

	// 6. Output to STDOUT or file
	output := os.Stdout
	if cmd.flagOutputPath != "" {
		output, err = os.Create(cmd.flagOutputPath)
		if err != nil {
			return fmt.Errorf("error creating output file for Framework IR: %w", err)
		}
	}

	_, err = output.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing framework IR to output: %w", err)
	}

	return nil
}

func generateFrameworkIr(dora explorer.Explorer, cfg config.Config) (*ir.IntermediateRepresentation, error) {
	// 1. Find TF resources
	resources, err := dora.FindResources()
	if err != nil {
		return nil, fmt.Errorf("error finding resources: %w", err)
	}

	// 2. Find TF data sources
	dataSources, err := dora.FindDataSources()
	if err != nil {
		return nil, fmt.Errorf("error finding data sources: %w", err)
	}

	// 3. Find TF provider
	provider, err := dora.FindProvider()
	if err != nil {
		return nil, fmt.Errorf("error finding provider: %w", err)
	}

	// 4. Use TF info to generate framework IR for resources
	resourceMapper := resource.NewResourceMapper(resources, cfg)
	resourcesIR, err := resourceMapper.MapToIR()
	if err != nil {
		return nil, fmt.Errorf("error generating Framework IR for resources: %w", err)
	}

	// 5. Use TF info to generate framework IR for data sources
	dataSourceMapper := datasource.NewDataSourceMapper(dataSources, cfg)
	dataSourcesIR, err := dataSourceMapper.MapToIR()
	if err != nil {
		return nil, fmt.Errorf("error generating Framework IR for data sources: %w", err)
	}

	return &ir.IntermediateRepresentation{
		Provider: ir.Provider{
			Name: provider.Name,
		},
		Resources:   resourcesIR,
		DataSources: dataSourcesIR,
	}, nil
}
