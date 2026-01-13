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
./jiracrawler issues --user <email> --project CNF

# Validate credentials
./jiracrawler config validate

# Output as JSON/YAML
./jiracrawler issues --format json
```

### Test
```bash
go test ./...
```

## Architecture

- **`cmd/`** - CLI command implementations using Cobra
- **`lib/`** - Jira API client library and utilities
- **`scripts/`** - Helper scripts
- **`main.go`** - Application entry point

## Configuration

Config file location: `.jiracrawler-config.yaml`
```yaml
jira_url: "https://issues.redhat.com"
jira_user: "your-email@example.com"
jira_token: "your-api-token"
```

## Features

- Fetch issues assigned to users
- Filter by project and status
- User activity tracking by date range
- JSON/YAML output formats
- Credential validation

## Requirements

- Go 1.21+
- Jira personal access token (PAT)

## Code Style

- Follow standard Go conventions
- Use `go fmt` before committing
- Run `golangci-lint` for linting
