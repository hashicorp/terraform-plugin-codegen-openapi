## Running the example spec
```sh
# Build the binary
go build -o tfopenapigen .

# Pass generator config to --config and OpenAPI spec (JSON or YML) as positional argument
# Outputs to STDOUT
./tfopenapigen generate --config ./examples/petstore3/tfopenapigen_config.yml ./examples/petstore3/openapi_spec.json
```