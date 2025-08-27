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

## run/curl: run a curl command to test the healthcheck endpoint
.PHONY: run/curl
run/curl:
	curl -i localhost:$(PORT)/$(VERSION)/healthcheck
	
# simple command to update version
.PHONY: bump/version
bump/version:
	@if [ ! -f .envrc ]; then echo "Error: .envrc file not found"; exit 1; fi; \
	current_version=$$(grep "VERSION=v" t.envrc | sed 's/.*VERSION=v\([0-9]*\).*/\1/'); \
	if [ -z "$$current_version" ]; then echo "Error: Could not find VERSION in .envrc"; exit 1; fi; \
	new_version=$$((curren_version + 1)); \
	sed -i.bak "s/VERSION=v[0-9]*/VERSION=v$$new_version/" .envrc; \
	echo "Version bumped: v$$current_version → v$$new_version"