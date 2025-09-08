# Use the .envrc file
include .envrc

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api --port=$(PORT) --env=$(ENV) --db-dsn=$(DB_DSN)

## run/tests: run the tests
.PHONY: run/tests
run/tests:
	go test ./...

## run/curl: run a curl command to test the healthcheck endpoint
.PHONY: run/curl
run/curl:
	curl -i localhost:$(PORT)/v$(VERSION)/healthcheck
	
# simple command to update version
.PHONY: bump/version
bump/version:
	@if [ ! -f .envrc ]; then echo "Error: .envrc file not found"; exit 1; fi; \
	current_version=$$(grep "VERSION=" .envrc | sed 's/.*VERSION=\([0-9]*\).*/\1/'); \
	if [ -z "$$current_version" ]; then echo "Error: Could not find VERSION in .envrc"; exit 1; fi; \
	new_version=$$((current_version + 1)); \
	sed -i.bak "s/VERSION=[0-9]*/VERSION=$$new_version/" .envrc; \
	echo "Version bumped: v$$current_version → v$$new_version"

# make command to post a comment
.PHONY: run/quote
run/quote:
	BODY='{"type":"Inspirational", "quote":"I am fond of pigs. Dogs look up to us. Cats look down on us. Pigs treat us as equals.", "author":"Winston S. Churchill"}'; \
	curl -i -H "Content-Type: application/json" -d "$$BODY" localhost:$(PORT)/v$(VERSION)/quote