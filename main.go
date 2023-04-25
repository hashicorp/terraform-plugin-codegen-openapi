package main

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/cmd"
	"os"

	"github.com/mattn/go-colorable"
)

func main() {
	name := "openapi"
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
