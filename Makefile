PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: go-jsonschema
go-jsonschema:
	GOBIN=$(PROJECT_DIR)/bin go install github.com/atombender/go-jsonschema@v0.17.0

.PHONY: generate
generate: go-jsonschema
	$(PROJECT_DIR)/bin/go-jsonschema -p apis $(PROJECT_DIR)/specification/schema/schema.json -o apis/apis.go

