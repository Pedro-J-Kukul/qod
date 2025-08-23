## run/api: run the cmd/api application

.PHONY: run/api
run/api:
	go run ./cmd/api

## run/tests: run the tests
.PHONY: run/tests
run/tests:
	go test ./...