# Setup name variables for the package/tool
NAME := gitsweeper
PKG := github.com/petems/$(NAME)
GIT_COMMIT := $(shell git log -1 --pretty=format:"%h" .)
VERSION := $(shell grep "const Version " main.go | sed -E 's/.*"(.+)"$$/\1/')

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: all
all: clean build fmt lint test install

.PHONY: clean build
build:
	@echo "building ${NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X main.gitCommit=${GIT_COMMIT}" -o bin/${NAME}

.PHONY: fmt
fmt: ## Verifies all files have men `gofmt`ed
	@echo "+ $@"
	@gofmt -s -l . | grep -v '.pb.go:' | grep -v vendor | tee /dev/stderr

.PHONY: lint
lint: ## Verifies `golint` passes
	@echo "+ $@"
	@golangci-lint run ./... | tee /dev/stderr

.PHONY: cover
cover: ## Runs go test with coverage
	@for d in $(shell go list ./... | grep -v vendor); do \
		go test -race -coverprofile=profile.out -covermode=atomic "$$d"; \
	done;

.PHONY: cover_html
cover_html: ## Runs go test with coverage
	@go tool cover -html=profile.out

.PHONY: clean
clean: clean-coverage ## Cleanup any build binaries or packages
	@echo "+ $@"
	$(RM) $(NAME)
	$(RM) -r $(BUILDDIR)

.PHONY: test
test: ## Runs the go tests
	@echo "+ $@"
	@go test ./...

.PHONY: acceptance-test
acceptance-test: ## Runs the acceptance tests
	@echo "+ $@"
	@go test -v -run "Test.*Command" .

.PHONY: test-all
test-all: test acceptance-test ## Runs all tests (unit and acceptance)

.PHONY: acceptance-cover
acceptance-cover: ## Runs acceptance tests with coverage
	@echo "+ $@"
	@go test -v -coverprofile=acceptance-coverage.out -run "Test.*Command" .

.PHONY: acceptance-cover-html
acceptance-cover-html: acceptance-cover ## Generates HTML coverage report for acceptance tests
	@echo "+ $@"
	@go tool cover -html=acceptance-coverage.out -o acceptance-coverage.html
	@echo "Coverage report generated: acceptance-coverage.html"

.PHONY: acceptance-cover-func
acceptance-cover-func: acceptance-cover ## Shows function-level coverage for acceptance tests
	@echo "+ $@"
	@go tool cover -func=acceptance-coverage.out

.PHONY: all-cover
all-cover: ## Runs all tests with coverage and generates combined report
	@echo "+ $@"
	@go test -v -coverprofile=all-coverage.out ./...
	@echo "Combined coverage report generated: all-coverage.out"

.PHONY: all-cover-html
all-cover-html: all-cover ## Generates HTML coverage report for all tests
	@echo "+ $@"
	@go tool cover -html=all-coverage.out -o all-coverage.html
	@echo "Combined coverage report generated: all-coverage.html"

.PHONY: integration-cover
integration-cover: ## Runs acceptance tests with integration coverage (Go 1.20+)
	@echo "+ $@"
	@rm -rf covdatafiles
	@mkdir -p covdatafiles
	@echo "Building coverage-instrumented binary..."
	@go build -cover -o bin/$(NAME)-integration-test main.go
	@echo "Running integration tests with coverage..."
	@GOCOVERDIR=$(PWD)/covdatafiles GITSWEEPER_TEST_BINARY=$(PWD)/bin/$(NAME)-integration-test go test -v -run "Test.*Command" .
	@echo "Processing coverage data..."
	@go tool covdata percent -i=$(PWD)/covdatafiles

.PHONY: integration-cover-html
integration-cover-html: integration-cover ## Generates HTML coverage report for integration tests
	@echo "+ $@"
	@go tool covdata textfmt -i=$(PWD)/covdatafiles -o=integration-coverage.txt
	@go tool cover -html=integration-coverage.txt -o=integration-coverage.html
	@echo "Integration coverage report generated: integration-coverage.html"

.PHONY: integration-cover-func
integration-cover-func: integration-cover ## Shows function-level coverage for integration tests
	@echo "+ $@"
	@go tool covdata textfmt -i=$(PWD)/covdatafiles -o=integration-coverage.txt
	@go tool cover -func=integration-coverage.txt

.PHONY: integration-cover-merge
integration-cover-merge: integration-cover ## Merges integration coverage data files
	@echo "+ $@"
	@rm -rf merged-covdata
	@mkdir -p merged-covdata
	@go tool covdata merge -i=$(PWD)/covdatafiles -o=merged-covdata
	@echo "Merged coverage data available in merged-covdata/"
	@go tool covdata percent -i=merged-covdata

.PHONY: clean-coverage
clean-coverage: ## Removes all coverage files
	@echo "+ $@"
	$(RM) *.out *.html *.txt profile.out acceptance-coverage.out all-coverage.out integration-coverage.txt
	$(RM) -r covdatafiles merged-covdata

.PHONY: install
install: ## Installs the executable or package
	@echo "+ $@"
	go install -a .
