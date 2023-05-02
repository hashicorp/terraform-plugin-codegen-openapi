package cmd

import (
	"io"

	"github.com/mitchellh/cli"
)

func initCommands(ui cli.Ui) map[string]cli.CommandFactory {

	generateFactory := func() (cli.Command, error) {
		return &generateCmd{
			ui: ui,
		}, nil
	}

	return map[string]cli.CommandFactory{
		"generate": generateFactory,
	}
}

func Run(name, version string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
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
