# jiracrawler

[![PR Checks](https://github.com/sebrandon1/jiracrawler/actions/workflows/pre-main.yml/badge.svg)](https://github.com/sebrandon1/jiracrawler/actions/workflows/pre-main.yml)
[![Nightly Integration Tests](https://github.com/sebrandon1/jiracrawler/actions/workflows/integration-nightly.yml/badge.svg)](https://github.com/sebrandon1/jiracrawler/actions/workflows/integration-nightly.yml)
[![Release binaries](https://github.com/sebrandon1/jiracrawler/actions/workflows/release-binaries.yaml/badge.svg)](https://github.com/sebrandon1/jiracrawler/actions/workflows/release-binaries.yaml)

A CLI tool for querying, filtering, and reporting on Jira issues for Red Hat projects (or any Jira instance). It is designed to help teams and individuals track their assigned issues, filter by project, and validate Jira API credentials using a personal access token (PAT).

## Features

- **Fetch assigned issues**: Query Jira for issues assigned to one or more users, filtered by project and status (e.g., only active issues).
- **User activity tracking**: Find all issues assigned to a user that were updated within a specific date range.
- **Configurable project**: Use a flag to specify which Jira project to query (default: CNF).
- **Flexible output**: Output results in JSON or YAML for easy integration with other tools or reporting.
- **Credential validation**: Quickly check if your Jira API token and user are valid with a single command.
- **Config management**: Store and update your Jira URL, user, and API token securely in a local config file.

## Usage

### Setup

1. Build the binary:
   ```bash
   make build
   ```

2. Configure your Jira credentials:
   ```bash
   ./jiracrawler config set --user <your-email> --token <your-api-token> --url https://issues.redhat.com
   ```
   - `user`: Your Jira username (often your email address)
   - `token`: Your Jira personal access token
   - `url`: (Optional) Jira instance URL (defaults to https://issues.redhat.com)

### Validate Credentials

```bash
./jiracrawler validate
```

### Get Assigned Issues

```bash
./jiracrawler get assignedissues <user1> <user2> --projectID CNF --output json
```
- `projectID`: Jira project key (default: CNF)
- `output`: Output format (`json` or `yaml`)

### Get User Updates in Date Range

```bash
./jiracrawler get userupdates <user@example.com> <start-date> <end-date> --output json
```
- `user@example.com`: The user whose assigned issues you want to query
- `start-date`: Start date in YYYY-MM-DD format
- `end-date`: End date in YYYY-MM-DD format
- `output`: Output format (`json` or `yaml`)

Example:
```bash
./jiracrawler get userupdates user@redhat.com 2024-01-01 2024-01-31 --output json
```

### Example jq Usage

To list active (not closed/completed) issues per user from the JSON output:
```bash
jq -r '
  .[] |
  "\(.user): " + ([.issues[] | select(.fields.status.name | test("(?i)closed|done|completed") | not) | .key] | join(", "))
' output.json
```

## Testing

### Unit Tests
```bash
make test
```

### Integration Tests
```bash
make integration-test
```

Integration tests run against a real JIRA instance using your configured credentials. For full testing, set the `JIRACRAWLER_TEST_USER` environment variable:

```bash
export JIRACRAWLER_TEST_USER=your-email@example.com
make integration-test
```

### Nightly Integration Tests (CI/CD)

The repository includes a GitHub Actions workflow that runs integration tests nightly against a real JIRA instance. This workflow:

- Runs every night at 2 AM UTC
- Only executes on the upstream repository (forks are blocked for security)
- Can be manually triggered via workflow dispatch
- Requires the following secrets to be configured in the repository settings:

#### Required Secrets
- `JIRA_URL`: Your JIRA instance URL (e.g., `https://issues.redhat.com`)
- `JIRA_USER`: JIRA username/email for authentication and testing
- `JIRA_API_TOKEN`: Personal access token for JIRA API
- `JIRA_TEST_PROJECT`: (Optional) Project to test against (defaults to `CNF`)

The workflow automatically builds the binary, configures credentials, runs all integration tests, and cleans up sensitive data afterward.

## Requirements
- Go 1.18+
- Jira account with API token

## Why?
This tool is for Red Hat and open source teams who want a simple, scriptable way to:
- Track their Jira workload
- Integrate Jira data into dashboards or reports
- Validate and manage Jira API credentials

## License
Apache 2.0
