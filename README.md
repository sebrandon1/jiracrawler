# jiracrawler

[![PR Checks](https://github.com/sebrandon1/jiracrawler/actions/workflows/pre-main.yml/badge.svg)](https://github.com/sebrandon1/jiracrawler/actions/workflows/pre-main.yml)
[![Nightly Integration Tests](https://github.com/sebrandon1/jiracrawler/actions/workflows/integration-nightly.yml/badge.svg)](https://github.com/sebrandon1/jiracrawler/actions/workflows/integration-nightly.yml)
[![Release binaries](https://github.com/sebrandon1/jiracrawler/actions/workflows/release-binaries.yaml/badge.svg)](https://github.com/sebrandon1/jiracrawler/actions/workflows/release-binaries.yaml)

A CLI tool for querying, filtering, and reporting on Jira issues. Designed for Red Hat projects but works with any Jira instance.

## Key Features

- Fetch issues assigned to one or more users, filtered by project and status
- Track user activity within a specific date range
- JSON and YAML output for easy integration with other tools
- Credential validation and config management

## Quick Start

```bash
make build
./jiracrawler config set --user <email> --token <api-token> --url https://issues.redhat.com
./jiracrawler validate
./jiracrawler get assignedissues <user> --projectID CNF --output json
```

## Guides

| Document | Description |
|----------|-------------|
| [CLI Reference](docs/cli-reference.md) | All commands, flags, and configuration options |
| [Testing](docs/testing.md) | Unit tests, integration tests, and nightly CI setup |
| [Examples](docs/examples.md) | jq recipes, common workflows, and output tips |
| [Contributing](CONTRIBUTING.md) | How to contribute to the project |

## Development

```bash
make build              # Build the binary
make test               # Run unit tests
make integration-test   # Run integration tests (requires Jira credentials)
make lint               # Run linter
make fmt                # Format code
make clean              # Remove build artifacts
```

## Requirements

- Go 1.26+
- Jira account with API token

## License

Apache 2.0
