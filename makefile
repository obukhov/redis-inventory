#!make
.PHONY: help docs build test
.DEFAULT_GOAL := help

help:## Show this help.
	@grep -E '^[a-zA-Z_-]+:.*?##\s*.*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?##\\s*"}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test:## Run tests
	go test ./...

docs:## Generate docs for cobra commands
	go run cmd/docgen/main.go

build:## Build docker services
	go build -o build/redis-inventory main.go
