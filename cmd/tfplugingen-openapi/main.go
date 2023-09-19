// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/cmd"
	"github.com/mitchellh/cli"

	"github.com/mattn/go-colorable"
)

func main() {
	name := "tfplugingen-openapi"
	version := name + " Version " + version
	if commit != "" {
		version += " from commit " + commit
	}

	os.Exit(runCLI(
		name,
		version,
		os.Args[1:],
		os.Stdin,
		colorable.NewColorableStdout(),
		colorable.NewColorableStderr(),
	))
}

func initCommands(ui cli.Ui) map[string]cli.CommandFactory {

	generateFactory := func() (cli.Command, error) {
		return &cmd.GenerateCommand{
			UI: ui,
		}, nil
	}

	return map[string]cli.CommandFactory{
		"generate": generateFactory,
	}
}

func runCLI(name, version string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,

		Ui: &cli.BasicUi{
			Reader:      stdin,
			Writer:      stdout,
			ErrorWriter: stderr,
		},
	}

	commands := initCommands(ui)
	openAPIGen := cli.CLI{
		Name:       name,
		Args:       args,
		Commands:   commands,
		HelpFunc:   cli.BasicHelpFunc(name),
		HelpWriter: stderr,
		Version:    version,
	}
	exitCode, err := openAPIGen.Run()
	if err != nil {
		return 1
	}

	return exitCode
}
