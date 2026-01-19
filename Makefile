# Copyright (C) 2024-2026 Bryce Thuilot
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the FSF, either version 3 of the License, or (at your option) any later version.
# See the LICENSE file in the root of this repository for full license text or
# visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

GO_FILES = $(shell find . -name '*.go')

BIN_NAME := "git-lost-and-found"

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

BIN = $(LOCALBIN)/$(BIN_NAME)
BIN_DEBUG = $(LOCALBIN)/$(BIN_NAME)-debug
ARGS ?= 

MAIN=./main.go

VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || echo "dev")
DOCKER_TAG ?= $(VERSION)
DOCKER_IMAGE ?= "ghcr.io/bthuilot/git-lost-and-found:$(VERSION)"

COMMIT_SHA := $(shell git rev-parse HEAD)
DIRTY := $(shell git diff --quiet || echo "-dirty")
BUILD_TIME := $(shell TZ="UTC" date -Iseconds)

METADATA_LD_FLAGS = -X 'github.com/bthuilot/git-lost-and-found/v2/cmd.version=$(VERSION)' -X 'github.com/bthuilot/git-lost-and-found/v2/cmd.gitCommit=$(COMMIT_SHA)$(DIRTY)' -X 'github.com/bthuilot/git-lost-and-found/v2/cmd.buildTime=$(BUILD_TIME)'
LD_FLAGS ?= -s -w

############
# Building #
############
.PHONY: build build-debug build-docker production-build

$(BIN): $(GO_FILES)
	go build -ldflags "$(LD_FLAGS) $(METADATA_LD_FLAGS)" -trimpath -o $(BIN) $(MAIN)


.PHONY: build
build: $(BIN)

.PHONY: build-debug
build-debug: $(GO_FILES)
	@go build -gcflags "all=-N -l" -o $(BIN)-debug $(MAIN)

docker-build:
	@docker build -t $(DOCKER_IMAGE) --build-arg LD_FLAGS="$(LD_FLAGS)" .

###########
# Running #
###########

.PHONY: run run-debug install clean fmt lint test
run: build
	@./$(BIN) $(ARGS)

run-debug: build-debug
	@./$(BIN) $(ARGS)

####################
# Install & Config #
####################

install: build
	@cp $(BIN) /usr/local/bin/$(BIN_NAME)

clean:
	@rm -rf $(BIN)

###########
# Linting #
###########

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter
	$(GOLANGCI_LINT) run
	license-eye header check


.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix
	license-eye header fix

###########
# Testing #
###########

test:
	@go test -v ./...

####################
# Tools & Binaries #
####################


GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
LICENSE_EYE ?= $(LOCALBIN)/license-eye

GOLANGCI_LINT_VERSION ?= v2.5.0
LICENSE_EYE_VERSION ?= v0.7.0


.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))


.PHONY: license-eye
license-eye: $(LICENSE_EYE) ## Download skywalking-eyes locally if necessary.
$(LICENSE_EYE): $(LOCALBIN)
	$(call go-install-tool,$(LICENSE_EYE),github.com/apache/skywalking-eyes/cmd/license-eye,$(LICENSE_EYE_VERSION))


# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef