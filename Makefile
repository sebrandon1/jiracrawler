# Makefile for Go project

.PHONY: all build test lint clean

all: build

build:
	go build -o jiracrawler

test:
	go test -v ./...

lint:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run ./...

clean:
	rm -rf bin/
	rm -f *.out *.test

vet:
	go vet ./...
