# Setup name variables for the package/tool
NAME := gitsweeper
PKG := github.com/petems/$(NAME)
GIT_COMMIT := $(shell git log -1 --pretty=format:"%h" .)
VERSION := $(shell grep "const Version " main.go | sed -E 's/.*"(.+)"$$/\1/')

.PHONY: all
all: clean build fmt lint test install

.PHONY: clean build build-optimized build-ultra-optimized size-comparison
build:
	@echo "building ${NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X main.gitCommit=${GIT_COMMIT}" -o bin/${NAME}

build-optimized:
	@echo "building optimized ${NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -tags optimized -ldflags="-s -w -X main.gitCommit=${GIT_COMMIT}" -trimpath -o bin/${NAME}-optimized

build-ultra-optimized:
	@echo "building ultra-optimized ${NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	CGO_ENABLED=0 go build -tags optimized -ldflags="-s -w -X main.gitCommit=${GIT_COMMIT}" -trimpath -a -installsuffix cgo -o bin/${NAME}-ultra

size-comparison: build build-optimized build-ultra-optimized
	@echo "Binary size comparison:"
	@ls -lh bin/${NAME}* | awk '{print $$9 " - " $$5}' | sort

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
clean: ## Cleanup any build binaries or packages
	@echo "+ $@"
	$(RM) $(NAME)
	$(RM) -r $(BUILDDIR)

.PHONY: test
test: ## Runs the go tests
	@echo "+ $@"
	@go test ./...

.PHONY: install
install: ## Installs the executable or package
	@echo "+ $@"
	go install -a .
