.PHONY: build build-debug run run-debug clean

GO_FILES = $(shell find . -name '*.go')
BIN = ./bin/$(shell basename $(CURDIR))
MAIN = ./main.go

build: $(GO_FILES)
	go build -ldflags "-s -w" -o $(BIN) $(MAIN)

build-debug: $(GO_FILES)
	go build -gcflags "all=-N -l" -o $(BIN) $(MAIN)

run: build
	@./$(BIN)

run-debug: build-debug
	@./$(BIN)

