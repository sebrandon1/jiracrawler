# Makefile for Go project

APP_NAME = jiracrawler
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build test integration-test lint clean vet fmt run help

all: build

build:
	go build -o $(APP_NAME)

test:
	go test -v ./...

integration-test: build
	@chmod +x scripts/integration-test.sh
	./scripts/integration-test.sh

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
	@echo "  build            - Build the binary ($(APP_NAME))"
	@echo "  test             - Run all unit tests"
	@echo "  integration-test - Run integration tests against real JIRA"
	@echo "  lint             - Run golangci-lint on the codebase"
	@echo "  vet              - Run go vet on the codebase"
	@echo "  fmt              - Format code with gofmt"
	@echo "  clean            - Remove build/test artifacts"
	@echo "  run              - Build and run the app"
