schema_version = 1

project {
  license        = "MPL-2.0"
  copyright_year = 2023

  header_ignore = [
    # internal catalog metadata (prose)
    "META.d/**/*.yaml",

    # changie tooling configuration and CHANGELOG entries (prose)
    ".changes/unreleased/*.yaml",
    ".changie.yaml",

    # GitHub issue template configuration
    ".github/ISSUE_TEMPLATE/*.yml",

    # GitHub Actions workflow-specific configurations
    ".github/labeler-*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",

    # GoReleaser tooling configuration
    ".goreleaser.yml",

    # Release Engineering tooling configuration
    ".release/*.hcl",

    # Unit test data files
    "internal/cmd/testdata/**",
  ]
}
