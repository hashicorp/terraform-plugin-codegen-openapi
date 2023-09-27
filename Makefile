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
		--config ./internal/cmd/testdata/petstore3/generator_config.yml \
		--output ./internal/cmd/testdata/petstore3/provider_code_spec.json \
		./internal/cmd/testdata/petstore3/openapi_spec.json

	go run ./cmd/tfplugingen-openapi generate \
		--config ./internal/cmd/testdata/github/generator_config.yml \
		--output ./internal/cmd/testdata/github/provider_code_spec.json \
		./internal/cmd/testdata/github/openapi_spec.json

	go run ./cmd/tfplugingen-openapi generate \
		--config ./internal/cmd/testdata/scaleway/generator_config.yml \
		--output ./internal/cmd/testdata/scaleway/provider_code_spec.json \
		./internal/cmd/testdata/scaleway/openapi_spec.yml

	go run ./cmd/tfplugingen-openapi generate \
		--config ./internal/cmd/testdata/edgecase/generator_config.yml \
		--output ./internal/cmd/testdata/edgecase/provider_code_spec.json \
		./internal/cmd/testdata/edgecase/openapi_spec.yml

.PHONY: lint fmt test