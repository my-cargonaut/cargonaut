# TOOLCHAIN
GO				:= CGO_ENABLED=0 GOFLAGS=-mod=vendor GOBIN=$(CURDIR)/bin go
GO_BIN_IN_PATH  := CGO_ENABLED=0 GOFLAGS=-mod=vendor go
GO_NO_VENDOR    := CGO_ENABLED=0 GOFLAGS=-mod=readonly GOBIN=$(CURDIR)/bin go
GOFMT			:= $(GO)fmt

# ENVIRONMENT
VERBOSE 	=
GOPATH		:= $(GOPATH)
GOOS		?= $(shell echo $(shell uname -s) | tr A-Z a-z)
GOARCH		?= amd64
MOD_NAME	:= github.com/my-cargonaut/cargonaut

# APPLICATION INFORMATION
BUILD_DATE		:= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
REVISION        := $(shell git rev-parse --short HEAD)
RELEASE         := $(shell cat RELEASE)
USER            := $(shell whoami)

# TOOLS
GOLANGCI_LINT	:= bin/golangci-lint
GOTESTSUM		:= bin/gotestsum
STATIK			:= bin/statik

# MISC
COVERPROFILE	:= coverage.out

# FLAGS
GOFLAGS			:= -buildmode=exe -tags='osusergo netgo static_build' \
					-installsuffix=cgo -trimpath \
					-ldflags='-s -w -extldflags "-fno-PIC -static" \
					-X $(MOD_NAME)/version.release=$(RELEASE) \
					-X $(MOD_NAME)/version.revision=$(REVISION) \
					-X $(MOD_NAME)/version.buildDate=$(BUILD_DATE) \
					-X $(MOD_NAME)/version.buildUser=$(USER)'

GOTESTSUM_FLAGS := --jsonfile=tests.json --junitfile=junit.xml
GO_TEST_FLAGS   := -coverprofile=$(COVERPROFILE)

# DEPENDENCIES
GOMODDEPS = go.mod go.sum

# Enable verbose test output if explicitly set.
ifdef VERBOSE
	GOTESTSUM_FLAGS	+= --format=standard-verbose
endif

# FUNCS
# func go-list-pkg-sources(package)
go-list-pkg-sources = $(GO) list $(GOFLAGS) -f '{{ range $$index, $$filename := .GoFiles }}{{ $$.Dir }}/{{ $$filename }} {{end}}' $(1)
# func go-pkg-sourcefiles(package)
go-pkg-sourcefiles = $(shell $(call go-list-pkg-sources,$(strip $1)))

.PHONY: all
all: dep assets fmt lint test build ## Run dep, assets, fmt, lint, test and build.

.PHONY: assets
assets: internal/sql/migrations/migrations.go internal/ui/ui.go ## Build all assets.
	@echo ">> formatting embeeded code"
	@$(GOFMT) -s -w internal/sql/migrations/statik.go
	@$(GOFMT) -s -w internal/ui/statik.go

.PHONY: build
build: .build/cargonaut-$(GOOS)-$(GOARCH) ## Build all binaries.

.PHONY: clean
clean: ## Remove build and test artifacts.
	@echo ">> cleaning up artifacts"
	@rm -rf .build $(COVERPROFILE) tests.json junit.xml

.PHONY: cover
cover: $(COVERPROFILE) ## Calculate the code coverage score.
	@echo ">> calculating code coverage"
	@$(GO) tool cover -func=$(COVERPROFILE)

.PHONY: dep-clean
dep-clean: ## Remove obsolete dependencies.
	@echo ">> cleaning dependencies"
	@$(GO_NO_VENDOR) mod tidy

.PHONY: dep-upgrade
dep-upgrade: ## Upgrade all direct dependencies to their latest version.
	@echo ">> upgrading dependencies"
	@$(GO) get $(shell $(GO_NO_VENDOR) list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	@$(GO) mod vendor
	@make dep

.PHONY: dep
dep: dep-clean dep.stamp ## Install and verify dependencies and remove obsolete ones.

dep.stamp: $(GOMODDEPS)
	@echo ">> installing dependencies"
	@$(GO) mod download
	@$(GO) mod verify
	@$(GO) mod vendor
	@touch $@

.PHONY: fmt
fmt: ## Format and simplify the source code using `gofmt`.
	@echo ">> formatting code"
	@! $(GOFMT) -s -w $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

.PHONY: install
install: $(GOPATH)/bin/cargonaut ## Install all binaries into the $GOPATH/bin directory.

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint the source code.
	@echo ">> linting code"
	@GO111MODULE=on $(GOLANGCI_LINT) run

.PHONY: test
test: $(GOTESTSUM) ## Run all tests. Run with VERBOSE=1 to get verbose test output ('-v' flag).
	@echo ">> running tests"
	@$(GOTESTSUM) $(GOTESTSUM_FLAGS) -- $(GO_TEST_FLAGS) ./...

.PHONY: tools
tools: $(GOLANGCI_LINT) $(GOTESTSUM) $(STATIK)  ## Install all tools into the projects local $GOBIN directory.

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# BUILD TARGETS

.build/cargonaut-darwin-amd64: dep.stamp $(call go-pkg-sourcefiles, ./...)
	@echo ">> building cargonaut production binary for darwin/amd64"
	@GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o=.build/cargonaut-darwin-amd64 ./cmd/cargonaut

.build/cargonaut-linux-amd64: dep.stamp $(call go-pkg-sourcefiles, ./...)
	@echo ">> building cargonaut production binary for linux/amd64"
	@GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o=.build/cargonaut-linux-amd64 ./cmd/cargonaut

.build/cargonaut-windows-amd64: dep.stamp $(call go-pkg-sourcefiles, ./...)
	@echo ">> building cargonaut production binary for windows/amd64"
	@GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o=.build/cargonaut-windows-amd64 ./cmd/cargonaut

# INSTALL TARGETS

$(GOPATH)/bin/cargonaut: dep.stamp $(call go-pkg-sourcefiles, ./...)
	@echo ">> installing cargonaut binary"
	@$(GO_BIN_IN_PATH) install $(GOFLAGS) ./cmd/cargonaut

# STATIK TARGETS

internal/sql/migrations/migrations.go: migrations/*.sql $(STATIK)
	@echo ">> embedding database migrations"
	@$(STATIK) -src=migrations -dest=internal/sql -p=migrations -f -m -ns=migrations -include=*.sql

internal/ui/ui.go: web/dist $(STATIK)
	@echo ">> embedding web ui"
	@$(STATIK) -src=web/dist -dest=internal -p=ui -f -m -ns=ui

# WEB TARGETS

# web/dist: web/src/App.vue web/src/main.js web/src/api/*.js web/src/assets/* web/src/components/*.vue web/src/plugins/*.js web/src/router/*.js web/src/store/*.js web/src/store/modules/*.js web/src/views/*.vue
web/dist: web/src/App.vue web/src/main.js web/src/api/*.js web/src/components/*.vue web/src/plugins/*.js web/src/router/*.js web/src/store/*.js web/src/store/modules/*.js web/src/views/*.vue
	@echo ">> building web ui"
	@cd web && yarn install
	@cd web && yarn lint
	@cd web && yarn build

# MISC TARGETS

$(COVERPROFILE):
	@make test

# TOOLS

$(GOLANGCI_LINT): dep.stamp $(call go-pkg-sourcefiles, ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint)
	@echo ">> installing golangci-lint"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint

$(GOTESTSUM): dep.stamp $(call go-pkg-sourcefiles, ./vendor/gotest.tools/gotestsum)
	@echo ">> installing gotestsum"
	@$(GO) install gotest.tools/gotestsum

$(STATIK): dep.stamp $(call go-pkg-sourcefiles, ./vendor/github.com/rakyll/statik)
	@echo ">> installing statik"
	@$(GO) install github.com/rakyll/statik
