# OpenAPI Provider Spec Generator

> _[Terraform Provider Code Generation](https://developer.hashicorp.com/terraform/plugin/code-generation) is currently in tech preview._

## Overview

The OpenAPI Provider Spec Generator is a CLI tool that transforms an [OpenAPI 3.x](https://www.openapis.org/) Specification (OAS) into a [Provider Code Specification](https://developer.hashicorp.com/terraform/plugin/code-generation/specification). 

A Provider Code Specification can be used to generate [Terraform Provider](https://developer.hashicorp.com/terraform/plugin) code, for example, with the [Plugin Framework Code Generator](https://developer.hashicorp.com/terraform/plugin/code-generation/framework-generator).

## Documentation

Full usage info, examples, and config file documentation live on the HashiCorp developer site: https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator

For more in-depth details about the design and internals of the OpenAPI Provider Spec Generator, see [`DESIGN.md`](./DESIGN.md).

## Usage

### Installation

You install a copy of the binary manually from the [releases](https://github.com/hashicorp/terraform-plugin-codegen-openapi/releases) tab, or install via the Go toolchain:

```shell-session
go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest
```

### Generate

The primary `generate` command requires a [generator config](https://developer.hashicorp.com/terraform/plugin/code-generation/openapi-generator#generator-config) and an OpenAPI 3.x specification:

```shell-session
tfplugingen-openapi generate \
  --config <path/to/generator_config.yml> \
  --output <output/for/provider_code_spec.json> \
  <path/to/openapi_spec.json>
```

### Examples

Example generator configs, OpenAPI specifications, and Provider Code Specification output can be found in the [`./internal/cmd/testdata/`](./internal/cmd/testdata/) folder. Here is an example running `petstore3`, built from source:

```shell-session
go run ./cmd/tfplugingen-openapi generate \
	--config ./internal/cmd/testdata/petstore3/generator_config.yml \
	--output ./internal/cmd/testdata/petstore3/provider_code_spec.json \
	./internal/cmd/testdata/petstore3/openapi_spec.json
```

## License

Refer to [Mozilla Public License v2.0](./LICENSE).

## Experimental Status

By using the software in this repository (the "Software"), you acknowledge that: (1) the Software is still in development, may change, and has not been released as a commercial product by HashiCorp and is not currently supported in any way by HashiCorp; (2) the Software is provided on an "as-is" basis, and may include bugs, errors, or other issues; (3) the Software is NOT INTENDED FOR PRODUCTION USE, use of the Software may result in unexpected results, loss of data, or other unexpected results, and HashiCorp disclaims any and all liability resulting from use of the Software; and (4) HashiCorp reserves all rights to make all decisions about the features, functionality and commercial release (or non-release) of the Software, at any time and without any obligation or liability whatsoever.
