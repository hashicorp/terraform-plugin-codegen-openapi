package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/explorer"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper"
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/pb33f/libopenapi"
)

type generateCmd struct {
	ui            cli.Ui
	flagInputPath string
}

func (cmd *generateCmd) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.StringVar(&cmd.flagInputPath, "input", "", "path to OpenAPI 3.1 spec file (.json or .yml)")
	return fs
}

func (cmd *generateCmd) Help() string {
	return "help text placeholder"
}

func (cmd *generateCmd) Synopsis() string {
	return "help text placeholder"
}

func (cmd *generateCmd) Run(args []string) int {
	fs := cmd.Flags()
	err := fs.Parse(args)
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("unable to parse flags: %s", err))
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
	// TODO: validate input path/default

	// 1. Read open API file
	example, err := os.ReadFile(cmd.flagInputPath)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	// 2. Parse basics of OpenAPI file
	doc, err := libopenapi.NewDocument(example)
	if err != nil {
		return fmt.Errorf("cannot create new doc: %w", err)
	}

	// 3. Build out the rest of the model, this will recursively load all local + remote references into one cohesive model
	model, errs := doc.BuildV3Model()
	// 4. Circular ref error handling - https://pb33f.io/libopenapi/circular-references/
	// - Not sure the implications of this yet, need to review
	if len(errs) > 0 {
		var errResult error
		for _, err := range errs {
			errResult = errors.Join(errResult, err)
		}
		log.Printf("[WARN] Potential issues in model spec: %s", errResult)
	}

	// TODO: Initial discovery takes too many guesses, will work on config-driven explorer next
	oasExplorer := explorer.NewGuesstimatorExplorer(model.Model)

	// 5. Generate framework IR
	frameworkIr, err := generateFrameworkIr(oasExplorer)
	if err != nil {
		return err
	}

	// 6. Use framework IR to create JSON file
	// TODO: accept output location as a CLI arg
	bytes, err := json.MarshalIndent(frameworkIr, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshalling Framework IR to JSON: %w", err)
	}

	_, err = os.Stdout.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing framework IR to standard output: %w", err)
	}

	return nil
}

func generateFrameworkIr(dora explorer.Explorer) (*ir.FrameworkIR, error) {
	// 1. Find TF resources
	resources, err := dora.FindResources()
	if err != nil {
		return nil, fmt.Errorf("error finding resources: %w", err)
	}

	// 2. Find TF data sources
	_, err = dora.FindDataSources()
	if err != nil {
		return nil, fmt.Errorf("error finding data sources: %w", err)
	}

	// 3. Find TF provider
	// TBD

	// 4. Use TF info to generate framework IR for resources
	resourceMapper := mapper.NewResourceMapper()
	resourcesIR, err := resourceMapper.MapToIR(resources)
	if err != nil {
		return nil, fmt.Errorf("error generating Framework IR for resources: %w", err)
	}

	// 5. Use TF info to generate framework IR for data sources
	// TBD

	return &ir.FrameworkIR{
		// TODO: eventually provider should move into explorer/mapper logic
		Provider: ir.Provider{
			Name: "placeholder",
		},
		Resources: *resourcesIR,
	}, nil
}
