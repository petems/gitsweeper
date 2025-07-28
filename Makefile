# Setup name variables for the package/tool
NAME := gitsweeper
PKG := github.com/petems/$(NAME)
LINT := golangci-lint
GIT_COMMIT := $(shell git log -1 --pretty=format:"%h" .)
VERSION := $(shell grep "const Version " main.go | sed -E 's/.*"(.+)"$$/\1/')

.PHONY: all
all: clean build fmt lint test install

.PHONY: clean build size-comparison
build:
	@echo "building ${NAME} ${VERSION} (ultra-optimized)"
	@echo "GOPATH=${GOPATH}"
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.gitCommit=${GIT_COMMIT}" -trimpath -a -installsuffix cgo -o bin/${NAME}

# Legacy compatibility builds (deprecated)
build-legacy:
	@echo "building legacy ${NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X main.gitCommit=${GIT_COMMIT}" -o bin/${NAME}-legacy

size-comparison: build build-legacy
	@echo "Binary size comparison:"
	@ls -lh bin/${NAME}* | awk '{print $$9 " - " $$5}' | sort

.PHONY: fmt
fmt: ## Verifies all files have men `gofmt`ed
	@echo "+ $@"
	@gofmt -s -l . | grep -v '.pb.go:' | grep -v vendor | tee /dev/stderr

## Run linter
lint:
	@echo "Checking golangci-lint version..."
	@$(LINT) version | grep -q "golangci-lint has version" || (echo "golangci-lint not found. Please install it first." && exit 1)
	@$(LINT) version | grep -oE "version [0-9]+\.[0-9]+\.[0-9]+" | cut -d' ' -f2 | awk -F. '{if ($$1 > 2 || ($$1 == 2 && $$2 >= 3)) exit 0; else exit 1}' || (echo "golangci-lint version 2.3.0 or higher required. Current version:" && $(LINT) version && exit 1)
	$(LINT) run

lint-fix:
	$(LINT) run --fix

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

.PHONY: cucumber
cucumber: ## Runs the cucumber integration tests
	@echo "+ $@"
	bundle exec cucumber

.PHONY: install
install: ## Installs the executable or package
	@echo "+ $@"
	go install -a .

.PHONY: build-all
build-all: ## Build binaries for all platforms
	@echo "+ $@"
	@echo "building ${NAME} ${VERSION} for multiple platforms"
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.gitCommit=${GIT_COMMIT}" -o bin/${NAME}-linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.gitCommit=${GIT_COMMIT}" -o bin/${NAME}-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.gitCommit=${GIT_COMMIT}" -o bin/${NAME}-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.gitCommit=${GIT_COMMIT}" -o bin/${NAME}-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.gitCommit=${GIT_COMMIT}" -o bin/${NAME}-windows-amd64.exe

.PHONY: release-archives
release-archives: build-all ## Create release archives for all platforms
	@echo "+ $@"
	@mkdir -p dist
	tar -czf dist/${NAME}-${VERSION}-linux-amd64.tar.gz -C bin ${NAME}-linux-amd64 -C .. README.md
	tar -czf dist/${NAME}-${VERSION}-linux-arm64.tar.gz -C bin ${NAME}-linux-arm64 -C .. README.md
	tar -czf dist/${NAME}-${VERSION}-darwin-amd64.tar.gz -C bin ${NAME}-darwin-amd64 -C .. README.md
	tar -czf dist/${NAME}-${VERSION}-darwin-arm64.tar.gz -C bin ${NAME}-darwin-arm64 -C .. README.md
	cd bin && zip ../dist/${NAME}-${VERSION}-windows-amd64.zip ${NAME}-windows-amd64.exe ../README.md
