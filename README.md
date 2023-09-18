# terraform-plugin-codegen-openapi

_Experimental: This code is under active development and not intended for production usage._

## Running an example spec
```sh
make build

# Pass generator config to --config and OpenAPI spec (JSON or YML) as positional argument
# Outputs to STDOUT
./tfplugingen-openapi  generate --config ./internal/cmd/testdata/petstore3/tfopenapigen_config.yml ./internal/cmd/testdata/petstore3/openapi_spec.json
```
