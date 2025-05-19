# Makefile for Go project

APP_NAME = jiracrawler
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build test lint clean vet fmt run help

all: build

build:
	go build -o $(APP_NAME)

test:
	go test -v ./...

lint:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run ./...

clean:
	rm -rf bin/
	rm -f *.out *.test $(APP_NAME)

vet:
	go vet ./...

fmt:
	gofmt -s -w $(GO_FILES)

run: build
	./$(APP_NAME)

help:
	@echo "Common make targets:"
	@echo "  build   - Build the binary ($(APP_NAME))"
	@echo "  test    - Run all tests"
	@echo "  lint    - Run golangci-lint on the codebase"
	@echo "  vet     - Run go vet on the codebase"
	@echo "  fmt     - Format code with gofmt"
	@echo "  clean   - Remove build/test artifacts"
	@echo "  run     - Build and run the app"
