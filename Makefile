.DEFAULT_GOAL := help

.PHONY: setup
setup:  ## Setup for required tools.
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/tools/cmd/stringer
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: fmt
fmt: ## Formatting source codes.
	@goimports -w $(find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: lint
lint: ## Run golint and go vet.
	@golint .
	@golint -set_exit_status=1
	@go vet .

.PHONY: test
test:  ## Run the tests with race condition checking.
	@go test -race .

.PHONY: coverage
cover:  ## Run the tests.
	@go test -coverprofile=coverage.o
	@go tool cover -func=coverage.o

.PHONY: help
help: ## Show help text
	@echo "Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
