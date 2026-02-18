# Contributing to jiracrawler

Thank you for your interest in contributing to jiracrawler!

## Prerequisites

- Go 1.26+
- A Jira instance with a personal access token (for integration tests)

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch from `main`

## Building

```bash
make build
```

## Testing

### Unit Tests

```bash
make test
```

### Integration Tests

Integration tests run against a real Jira instance:

```bash
export JIRACRAWLER_TEST_USER=your-email@example.com
make integration-test
```

## Linting and Formatting

```bash
make lint    # Run golangci-lint
make vet     # Run go vet
make fmt     # Format code with gofmt
```

Please ensure all lint and vet checks pass before submitting a PR.

## Submitting Changes

1. Create a feature branch from `main`
2. Make your changes
3. Run `make fmt`, `make vet`, and `make test`
4. Commit with a clear, descriptive message
5. Push your branch and open a pull request against `main`

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Keep changes focused and minimal
