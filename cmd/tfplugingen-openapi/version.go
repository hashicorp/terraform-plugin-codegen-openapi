// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"runtime/debug"
)

var (
	// These vars will be set by goreleaser.
	version string
	commit  string
)

func getVersion() string {
	// Prefer global version as it's set by goreleaser via ldflags
	// https://goreleaser.com/cookbooks/using-main.version/
	if version != "" {
		if commit != "" {
			version = fmt.Sprintf("%s from commit: %s", version, commit)
		}
		return version
	}

	// If not built with goreleaser, check the binary for VCS revision/module version info
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return fmt.Sprintf("commit: %s", setting.Value)
			}
		}

		return fmt.Sprintf("module: %s", info.Main.Version)
	}

	return "local"
}
