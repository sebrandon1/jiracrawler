# Examples

## JSON Output with jq

List active (not closed/completed) issues per user from the JSON output:

```bash
jq -r '
  .[] |
  "\(.user): " + ([.issues[] | select(.fields.status.name | test("(?i)closed|done|completed") | not) | .key] | join(", "))
' output.json
```

## Common Workflows

### Weekly status report

Fetch all issues updated in the past week for a user:

```bash
./jiracrawler get userupdates user@redhat.com 2024-01-08 2024-01-15 --output json
```

### Multi-user assigned issues

Query assigned issues for multiple team members at once:

```bash
./jiracrawler get assignedissues user1@redhat.com user2@redhat.com --projectID CNF --output yaml
```

### Pipe to file

```bash
./jiracrawler get assignedissues user@redhat.com --output json > team-issues.json
```
