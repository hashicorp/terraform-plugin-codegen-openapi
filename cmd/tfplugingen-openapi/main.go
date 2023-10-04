// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"io"
	"os"
	"runtime/debug"

	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/cmd"
)

// version will be set by goreleaser via ldflags
// https://goreleaser.com/cookbooks/using-main.version/
func main() {
	name := "tfplugingen-openapi"
	version := name + " commit: " + func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					return setting.Value
				}
			}

			return info.Main.Version
		}

		return "local"
	}()

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
