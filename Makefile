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
# Commands to interact with the API
########################################################################################################

# get healthcheck endpoint
.PHONY: api/healthcheck
api/healthcheck:
	@echo "Testing healthcheck endpoint..."
	@curl -i localhost:$(PORT)/v5/healthcheck

# --------------------------------------------------------------------------------------------------------	
# Quotes API commands
# --------------------------------------------------------------------------------------------------------
# POST API commands using type, quote and author variables from .envrc
.PHONY: api/quotes/post
api/quotes/post:
	@echo "Posting a quote..."
	BODY='{"type":$(TYPE), "quote":$(QUOTE), "author":$(AUTHOR)}'; \
	curl -i -H "Content-Type: application/json" -d "$$BODY" localhost:$(PORT)/v2/quotes
	

# GET a quote with id input
.PHONY: api/quotes/get
api/quotes/get:
	@echo "Getting a quote..."
	curl -i localhost:$(PORT)/v1/quotes/$(id)

# UPDATE a quote  with id and UPDATEBODY from .envrc input
.PHONY: api/quotes/update
api/quotes/update:
	@echo "Updating a quote..."
	curl -i -X PATCH -H "Content-Type: application/json" -d $(QUOTESBODY) localhost:$(PORT)/v1/quotes/$(id)

# DELETE a quote with id input
.PHONY: api/quotes/delete
api/quotes/delete:
	@echo "Deleting a quote..."
	curl -i -X DELETE localhost:$(PORT)/v1/quotes/$(id)

# GETALL quotes, list quotes
.PHONY: api/quotes/list
api/quotes/list:
	@echo "Listing quotes..."
	curl -i localhost:$(PORT)/v1/quotes

# make command to list quotes with query parameters
.PHONY: api/quotes/list/query
api/quotes/list/query:
	@echo "Listing quotes with filters..."
	curl -i "localhost:$(PORT)/v1/quotes?$(QUERY)"
# --------------------------------------------------------------------------------------------------------
# Users API commands
# --------------------------------------------------------------------------------------------------------
# POST API commands using name, email and password variables from .envrc
.PHONY: api/users/post
api/users/post:
	@echo "Registering a user..."
	BODY='{"username":$(USERNAME), "email":$(USEREMAIL), "password":$(USERPASSWORD)}'; \
	curl -i -d "$$BODY" localhost:$(PORT)/v1/users

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

