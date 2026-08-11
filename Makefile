.PHONY: build test lint install clean build-mcp install-mcp build-media-mcp install-media-mcp build-all

BIN_EXT := $(if $(filter windows,$(shell go env GOOS)),.exe,)

build:
	go build -o bin/kie-pp-cli$(BIN_EXT) ./cmd/kie-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/kie-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/kie-pp-mcp$(BIN_EXT) ./cmd/kie-pp-mcp

install-mcp:
	go install ./cmd/kie-pp-mcp

build-media-mcp:
	go build -o bin/kie-media-mcp$(BIN_EXT) ./cmd/kie-media-mcp

install-media-mcp:
	go install ./cmd/kie-media-mcp

build-all: build build-mcp build-media-mcp
