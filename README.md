# terraform-plugin-codegen-openapi

_Experimental: This code is under active development and not intended for production usage._

## Running an example spec
```sh
# Build the binary
go build ./cmd/terraform-plugin-codegen-openapi

# Pass generator config to --config and OpenAPI spec (JSON or YML) as positional argument
# Outputs to STDOUT
./terraform-plugin-codegen-openapi  generate --config ./internal/cmd/testdata/petstore3/tfopenapigen_config.yml ./internal/cmd/testdata/petstore3/openapi_spec.json
```
