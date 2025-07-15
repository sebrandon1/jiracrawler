# lib package

This package provides reusable functions for interacting with Jira and formatting output. These functions are intended to be called from the CLI layer or other Go code.

## Public Functions

### FetchAssignedIssuesWithProject
```go
func FetchAssignedIssuesWithProject(jiraURL, jiraUser, apikey string, projectID string, users []string) []map[string]interface{}
```
Fetches assigned issues for the given users in the specified Jira project. Returns a slice of maps containing user and issue data.

- `jiraURL`: The Jira instance URL (e.g., https://issues.redhat.com)
- `jiraUser`: The Jira username (usually an email address)
- `apikey`: The Jira personal access token
- `projectID`: The Jira project key (e.g., CNF)
- `users`: Slice of usernames to query

### FetchUserIssuesInDateRange
```go
func FetchUserIssuesInDateRange(jiraURL, jiraUser, apikey string, assignee string, startDate, endDate string) []map[string]interface{}
```
Fetches issues assigned to a specific user that were updated within a given date range. Returns a slice containing the user, date range, total count, and issues data.

- `jiraURL`: The Jira instance URL (e.g., https://issues.redhat.com)
- `jiraUser`: The Jira username (usually an email address)
- `apikey`: The Jira personal access token
- `assignee`: The username whose assigned issues to query
- `startDate`: Start date in YYYY-MM-DD format
- `endDate`: End date in YYYY-MM-DD format

### PrintJSON
```go
func PrintJSON(data interface{})
```
Prints the provided data as pretty-printed JSON to stdout.

### PrintYAML
```go
func PrintYAML(data interface{})
```
Prints the provided data as YAML to stdout.

---

These functions are used by the CLI in `cmd/` but can also be imported and used in other Go programs for Jira automation and reporting.
