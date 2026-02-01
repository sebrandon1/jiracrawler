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

# View configuration
./jiracrawler config view

# Validate credentials
./jiracrawler validate

# Output as JSON/YAML
./jiracrawler get assignedissues <user> -o json
```

### Test
```bash
make test              # Run unit tests
make integration-test  # Run integration tests against real Jira
```

### Lint and Format
```bash
make lint
make vet
make fmt
```

### Clean
```bash
make clean
```

## Architecture

- **`cmd/`** - CLI command implementations using Cobra
- **`lib/`** - Jira API client library with rate limiting
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
- JSON/YAML output formats
- Credential validation

## Requirements

- Go 1.25+
- Jira personal access token (PAT)

## Code Style

- Follow standard Go conventions
- Use `go fmt` before committing
- Run `golangci-lint` for linting
