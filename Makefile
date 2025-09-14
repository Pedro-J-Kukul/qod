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
.PHONY: api/insert
api/insert:
	@echo "Posting a quote..."
	BODY='{"type":"funny", "quote":"I am fond of pigs. Dogs look up to us. Cats look down on us. Pigs treat us as equals.", "author":"Winston S. Churchill"}'; \
	curl -i -H "Content-Type: application/json" -d "$$BODY" localhost:$(PORT)/v6/quotes

.PHONY: api/get
api/get:
	@echo "Getting a quote..."
	curl -i localhost:$(PORT)/v1/quotes/1

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
	migrate -path ./migrations -database "$(DB_DSN)" up 1

# Apply all down 1 migrations
.PHONY: migration/down
migration/down:
	migrate -path ./migrations -database "$(DB_DSN)" down 1

# fix and reapply the last migration and fix dirty state
.PHONY: migration/fix
migration/fix:
	@echo 'Checking migration status...'
	@migrate -path ./migrations -database "${DB_DSN}" version > /tmp/migrate_version 2>&1
	@cat /tmp/migrate_version
	@if grep -q "dirty" /tmp/migrate_version; then \
		version=$$(grep -o '[0-9]\+' /tmp/migrate_version | head -1); \
		echo "Found dirty migration at version $$version"; \
		echo "Forcing version $$version..."; \
		migrate -path ./migrations -database "${DB_DSN}" force $$version; \
		echo "Running down migration..."; \
		migrate -path ./migrations -database "${DB_DSN}" down 1; \
		echo "Running up migration..."; \
		migrate -path ./migrations -database "${DB_DSN}" up; \
	else \
		echo "No dirty migration found"; \
	fi
	@rm -f /tmp/migrate_version

# Login to psql
.PHONY: psql/login
psql/login:
	psql "$(DB_DSN)"

# login to postgresql as sudo
.PHONY: psql/sudo
psql/sudo:
	sudo -u postgres psql
