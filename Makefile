.PHONY: build build-debug run run-debug clean

GO_FILES = $(shell find . -name '*.go')
BIN_NAME := $(shell basename $(CURDIR))
BIN= ./bin/$(BIN_NAME)
DOCKER_NAME = $(BIN_NAME):latest
MAIN = ./main.go

build: $(GO_FILES)
	go build -ldflags "-s -w" -o $(BIN) $(MAIN)

build-debug: $(GO_FILES)
	go build -gcflags "all=-N -l" -o $(BIN) $(MAIN)

build-docker:
	docker build -t $(DOCKER_NAME) .

run: build
	@./$(BIN)

run-debug: build-debug
	@./$(BIN)

