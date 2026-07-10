# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A CLI tool for querying, filtering, and reporting on Jira issues. Designed for Red Hat projects but works with any Jira instance. Supports fetching assigned issues, tracking user activity, and credential validation.

## Common Commands

### Build
```bash
make build
```

### Configure
```bash
./jiracrawler config set --user <email> --token <api-token> --url https://issues.redhat.com
```

### Run
```bash
# Get issues assigned to a user
./jiracrawler get assignedissues <user> --projectID CNF

# Get issues updated by user in date range
./jiracrawler get userupdates <user> <start-date> <end-date>

# Run a custom JQL query
./jiracrawler get query "project = CNF AND status = Open" --max-results 10 -o table

# View configuration
./jiracrawler config view

# Validate credentials
./jiracrawler validate

# Output as JSON/YAML
./jiracrawler get assignedissues <user> -o json

# Generate shell completion scripts (bash/zsh/fish/powershell)
./jiracrawler completion bash
```

### Test
```bash
make test              # Run unit tests
make coverage          # Run tests with coverage report
make integration-test  # Run integration tests against real Jira
```

### Lint and Format
```bash
make lint
make vet
make fmt
```

### Other
```bash
make clean             # Remove build/test artifacts
make run               # Build and run the app
make help              # Show all available make targets
```

## Architecture

- **`cmd/`** - CLI command implementations using Cobra
  - `root.go` - Root command and global flags
  - `config.go` - Config set/view subcommands
  - `get.go` - Parent get command with `assignedissues` and `userupdates` subcommands
  - `query.go` - `get query` subcommand for arbitrary JQL queries
  - `validate.go` - Credential validation command
  - `completion.go` - Shell completion script generation (bash/zsh/fish/powershell)
- **`lib/`** - Jira API client library
  - `jira.go` - Core Jira client, issue types, and API interaction
  - `enhanced_context.go` - Enhanced issue context fetching (comments, change history, additional fields, permissions checks)
  - `rate_limiter.go` - Thread-safe rate limiter with exponential backoff, Retry-After header support, and 429 retry logic
- **`scripts/`** - Helper scripts
- **`main.go`** - Application entry point

## Configuration

Config file location: `.jiracrawler-config.yaml`
```yaml
url: "https://issues.redhat.com"
username: "your-email@example.com"
apikey: "your-api-token"
```

Environment variables can also be used with prefix `JIRACRAWLER_` (e.g., `JIRACRAWLER_URL`, `JIRACRAWLER_APIKEY`).

## Features

- Fetch issues assigned to users
- Filter by project and status
- User activity tracking by date range
- Custom JQL query support (`get query` command with `--max-results` flag)
- Enhanced issue context fetching (comments, change history, labels, components, time tracking)
- Rate limiting with exponential backoff and Retry-After header support for Jira API calls
- JSON/YAML/table output formats
- Credential validation
- Shell completion generation (bash/zsh/fish/powershell)

## Requirements

- Go 1.26+
- Jira personal access token (PAT)

## Code Style

- Follow standard Go conventions
- Use `go fmt` before committing
- Run `golangci-lint` for linting
