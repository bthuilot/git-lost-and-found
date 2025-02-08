GO_FILES = $(shell find . -name '*.go')
BIN_NAME := "git-lost-and-found"
BIN_DIR := ./bin
BIN= $(BIN_DIR)/$(BIN_NAME)
DOCKER_REPO=$(BIN_NAME)
DOCKER_TAG=dev
DOCKER_REGISTRY=
MAIN=./main.go
ARGS=

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
CGO_ENABLED := $(shell go env CGO_ENABLED)

############
# Building #
############
.PHONY: build build-debug

$(BIN): $(GO_FILES)
	@go build -ldflags "-s -w" -o $(BIN) $(MAIN)

build: $(BIN)

build-debug: $(GO_FILES)
	@go build -gcflags "all=-N -l" -o $(BIN) $(MAIN)

build-docker: $(BIN)
	@docker build -t $(DOCKER_REGISTRY)$(DOCKER_REPO):$(DOCKER_TAG) .

production-build: $(GO_FILES)
	@go build -ldflags "-s -w" -o $(BIN)-$(GOOS)-$(GOARCH) $(MAIN)

###########
# Running #
###########

run: build
	@./$(BIN) $(ARGS)

run-debug: build-debug
	@./$(BIN) $(ARGS)

####################
# Install & Config #
####################

install: build
	@cp ./$(BIN) /usr/local/bin/$(BIN_NAME)

clean:
	@rm -rf $(BIN_DIR)/* $(DIST)/*

###########
# Linting #
###########

fmt:
	@gofmt -l -s -w .

lint:
	@test -z $(gofmt -l .)
	@golangci-lint run

###########
# Testing #
###########

test:
	@go test -v ./...