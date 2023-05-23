# terraform-plugin-codegen-openapi

_Experimental: This code is under active development and not intended for production usage._

## Running an example spec
```sh
# Build the binary
go build -o tfopenapigen .

# Pass generator config to --config and OpenAPI spec (JSON or YML) as positional argument
# Outputs to STDOUT
./tfopenapigen generate --config ./examples/petstore3/tfopenapigen_config.yml ./examples/petstore3/openapi_spec.json
```
