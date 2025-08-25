# Use the .envrc file
include .envrc

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api --port=$(PORT) --env=$(ENV) --version=$(VERSION)

## run/tests: run the tests
.PHONY: run/tests
run/tests:
	go test ./...