# Variables
PKG          := ./...
BIN_DIR      := bin
LINT_CONTEXT := ./...

# Standard Go commands
GOCMD        := go
GOTEST       := $(GOCMD) test
GOMOD        := $(GOCMD) mod
GOLINT       := golangci-lint

.PHONY: all help tidy modernize fmt lint test test-race cover clean

all: help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## tidy: Format code and sync dependencies
tidy:
	$(GOCMD) fmt $(PKG)
	$(GOMOD) tidy

## modernize: Apply modern Go suggestions (requires modernize tool)
modernize:
	go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest -fix ./...

## fmt: Run gofumpt or gofmt if you prefer stricter formatting
fmt:
	$(GOCMD) fmt $(PKG)

## fix: Run all automated fixers (fmt, modernize, and lint --fix)
fix: tidy modernize
	$(GOLINT) run --fix ./...

## lint: Run golangci-lint (requires installation)
lint:
	$(GOLINT) run $(LINT_CONTEXT)

## test: Run all unit tests
test:
	$(GOTEST) -v -short $(PKG)

## test-race: Run tests with the data race detector
test-race:
	$(GOTEST) -v -race $(PKG)

## cover: Run tests and generate coverage report
cover:
	$(GOTEST) -count=1 -covermode=atomic -coverprofile=coverage.out $(PKG)
	$(GOCMD) tool cover -html=coverage.out

## clean: Remove temporary files and test cache
clean:
	rm -f coverage.out
	$(GOCMD) clean -testcache