package main

import (
	"os"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/cmd"

	"github.com/mattn/go-colorable"
)

func main() {
	// TODO: Temporary name for CLI :)
	name := "tfopenapigen"
	version := name + " Version " + version
	if commit != "" {
		version += " from commit " + commit
	}

	os.Exit(cmd.Run(
		name,
		version,
		os.Args[1:],
		os.Stdin,
		colorable.NewColorableStdout(),
		colorable.NewColorableStderr(),
	))
}
