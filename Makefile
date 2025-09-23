# Use the .envrc file
include .envrc

########################################################################################################
# Commands to run the application and tests
########################################################################################################
# run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo "Starting API server on port $(PORT) in $(ENV) mode..."
	@go run ./cmd/api \
	-port=$(PORT) \
	-env=$(ENV) \
	-db-dsn=$(DB_DSN) \
	-cors-trusted-origins="$(CORS_TRUSTED_ORIGINS)\
	-rate-limiter-enabled=$(RATE_LIMITER_ENABLED) \
	-rate-limiter-rps=$(RATE_LIMITER_RPS) \
	-rate-limiter-burst=$(RATE_LIMITER_BURST)"

## run/tests: run the tests
.PHONY: run/tests
run/tests:
	@echo "Running tests..."
	@go test ./...


########################################################################################################
# Commands to run example applications
########################################################################################################
# run/cors-basic
.PHONY: cors/basic
cors/basic:
	@echo "Basic CORS test"
	@go run ./cmd/examples/cors/basic

# run/cors-preflight
.PHONY: cors/preflight
cors/preflight:
	@echo "preflight CORS test"
	@go run ./cmd/examples/cors/preflight

########################################################################################################
# Commands to interact with the Quotes API
########################################################################################################

# simple curl command to test healthcheck endpoint
.PHONY: api/healthcheck
api/healthcheck:
	@echo "Testing healthcheck endpoint..."
	@curl -i localhost:$(PORT)/v5/healthcheck
	
# make command to post a comment using QOUTE, AUTHOR, TYPE from .envrc	
.PHONY: api/post/individual
api/post/individual:
	@echo "Posting a quote..."
	BODY='{"type":$(TYPE), "quote":$(QUOTE), "author":$(AUTHOR)}'; \
	curl -i -H "Content-Type: application/json" -d "$$BODY" localhost:$(PORT)/v2/quotes

# command to post a comment using body data from .envrc
.PHONY: api/post/body
api/post/body:
	@echo "Posting a quote..."
	BODY='$(UPDATEBODY)'; \
	curl -i -H "Content-Type: application/json" -d "$$BODY" localhost:$(PORT)/v2/quotes

# make command to get a comment with id input
.PHONY: api/get
api/get:
	@echo "Getting a quote..."
	curl -i localhost:$(PORT)/v1/quotes/$(id)

# make command to update a comment with body data from .envrc
.PHONY: api/update
api/update:
	@echo "Updating a quote..."
	BODY='$(UPDATEBODY)'; \
	curl -i -X PATCH -H "Content-Type: application/json" -d "$$BODY" localhost:$(PORT)/v1/quotes/$(id)

# make command to test update with no fields
.PHONY: api/update/empty
api/update/empty:
	@echo "Updating a quote with no fields..."
	BODY='{}'; \
	curl -i -X PATCH -H "Content-Type: application/json" -d "$$BODY" localhost:$(PORT)/v1/quotes/$(id)

# make command to delete a comment with id input
.PHONY: api/delete
api/delete:
	@echo "Deleting a quote..."
	curl -i -X DELETE localhost:$(PORT)/v1/quotes/$(id)

# make command to list quotes
.PHONY: api/list
api/list:
	@echo "Listing quotes..."
	curl -i localhost:$(PORT)/v1/quotes

# make command to list quotes with query parameters
.PHONY: api/list/query
api/list/query:
	@echo "Listing quotes with filters..."
	curl -i "localhost:$(PORT)/v1/quotes?$(QUERY)"

# make command to list quotes with page and page_size filter
.PHONY: api/list/pagination
api/list/pagination:
	@echo "Listing quotes with page and page_size filter..."
	curl -i "localhost:$(PORT)/v1/quotes?page=$(pg)&page_size=$(sz)"

# make command to list quotes with type filter
.PHONY: api/list/pagination-sort
api/list/pagination-sort:
	@echo "Listing quotes with sort filter..."
	curl -i "localhost:$(PORT)/v1/quotes?page=$(pg)&page_size=$(sz)&sort=$(sort)"

# command to list quotes with sort
.PHONY: api/list/sort
api/list/sort:
	@echo "Listing quotes with sort..."
	curl -i "localhost:$(PORT)/v1/quotes?sort=$(sort)"

# list with sort pagination and query
.PHONY: api/list/allfilters
api/list/allfilters:
	@echo "Listing quotes with all filters..."
	curl -i "localhost:$(PORT)/v1/quotes?$(QUERY)&page=$(pg)&page_size=$(sz)&sort=$(sort)"

########################################################################################################
# Database migration commands
########################################################################################################
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

########################################################################################################
# PostgreSQL commands
########################################################################################################
# Login to psql
.PHONY: psql/login
psql/login:
	psql "$(DB_DSN)"

# login to postgresql as sudo
.PHONY: psql/sudo
psql/sudo:
	sudo -u postgres psql

########################################################################################################
# Git commands
########################################################################################################
# Git Deleting all local Branches except main
.PHONY: git/cleanup
git/cleanup:
	@echo "Deleting all local branches except 'main'..."
	@git branch | grep -v "^main$" | xargs git branch -D