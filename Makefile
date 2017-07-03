NAME := prompt
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

.DEFAULT_GOAL := help

setup:  ## Setup for required tools.
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	glide install

fmt: ## Formatting source codes.
	@goimports -w $$(glide nv -x)

lint: ## Run golint and go vet.
	@golint $$(glide novendor)
	@go vet $$(glide novendor)

test:  ## Run the tests.
	@go test $$(glide novendor)

build: main.go  ## Build a binary.
	go build -ldflags "$(LDFLAGS)"

help: ## Show help text
	@echo "Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: setup fmt lint test help build

