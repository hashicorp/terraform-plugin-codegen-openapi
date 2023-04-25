## Running the example spec
```sh
# Build the binary
go build .

# Pass OpenAPI spec (JSON or YML) to --input argument
./terraform-providers-code-generator-openapi generate --input ./examples/scaleway.instance.v1.api.yml

# Outputs to `generated_framework_ir.json`
cat generated_framework_ir.json
```