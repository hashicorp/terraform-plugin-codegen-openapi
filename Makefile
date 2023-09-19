build:
	go build ./cmd/tfplugingen-openapi

lint:
	golangci-lint run

fmt:
	gofmt -s -w -e .

test:
	go test $$(go list ./... | grep -v /output) -v -cover -timeout=120s -parallel=4

# Generate copywrite headers
generate:
	cd tools; go generate ./...

# Regenerate testdata folder
testdata:
	go run ./cmd/tfplugingen-openapi generate \
		--config ./internal/cmd/testdata/petstore3/tfopenapigen_config.yml \
		./internal/cmd/testdata/petstore3/openapi_spec.json \
		> ./internal/cmd/testdata/petstore3/generated_framework_ir.json

	go run ./cmd/tfplugingen-openapi generate \
		--config ./internal/cmd/testdata/github/tfopenapigen_config.yml \
		./internal/cmd/testdata/github/openapi_spec.json \
		> ./internal/cmd/testdata/github/generated_framework_ir.json

	go run ./cmd/tfplugingen-openapi generate \
		--config ./internal/cmd/testdata/scaleway/tfopenapigen_config.yml \
		./internal/cmd/testdata/scaleway/openapi_spec.yml \
		> ./internal/cmd/testdata/scaleway/generated_framework_ir.json

	go run ./cmd/tfplugingen-openapi generate \
		--config ./internal/cmd/testdata/edgecase/tfopenapigen_config.yml \
		./internal/cmd/testdata/edgecase/openapi_spec.yml \
		> ./internal/cmd/testdata/edgecase/generated_framework_ir.json

.PHONY: lint fmt test