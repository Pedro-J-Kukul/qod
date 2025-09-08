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
	@echo "Posting a quote..."
	BODY='{"type":"Inspirational", "quote":"I am fond of pigs. Dogs look up to us. Cats look down on us. Pigs treat us as equals.", "author":"Winston S. Churchill"}'; \
	curl -i -H "Content-Type: application/json" -d "$$BODY" localhost:$(PORT)/v$(VERSION)/quote

# Create a new migration file
.PHONY: migration/create
migration/create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a name for the migration using 'make migration/create name=your_migration_name'"; \
		exit 1; \
	fi
	@if [ ! -d "./migrations" ]; then mkdir ./migrations; fi
	migrate create -seq -ext=.sql -dir=./migrations $(name)

# Apply all up migrations
.PHONY: migration/up
migration/up:
	migrate -path ./migrations -database "$(DB_DSN)" up

# Apply all down migrations
.PHONY: migration/down
migration/down:
	migrate -path ./migrations -database "$(DB_DSN)" down

# Login to psql
.PHONY: psql/login
psql/login:
	psql "$(DB_DSN)"

# login to postgresql as sudo
.PHONY: psql/sudo
psql/sudo:
	sudo -u postgres psql
