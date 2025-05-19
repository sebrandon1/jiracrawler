# jiracrawler

A CLI tool for querying, filtering, and reporting on Jira issues for Red Hat projects (or any Jira instance). It is designed to help teams and individuals track their assigned issues, filter by project, and validate Jira API credentials using a personal access token (PAT).

## Features

- **Fetch assigned issues**: Query Jira for issues assigned to one or more users, filtered by project and status (e.g., only active issues).
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

### Example jq Usage

To list active (not closed/completed) issues per user from the JSON output:
```bash
jq -r '
  .[] |
  "\(.user): " + ([.issues[] | select(.fields.status.name | test("(?i)closed|done|completed") | not) | .key] | join(", "))
' output.json
```

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
