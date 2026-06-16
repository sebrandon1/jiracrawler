# CLI Reference

## Configuration

Store your Jira credentials in a local config file:

```bash
./jiracrawler config set --user <your-email> --token <your-api-token> --url https://issues.redhat.com
```

| Flag    | Description                                          | Default                      |
|---------|------------------------------------------------------|------------------------------|
| `user`  | Your Jira username (often your email address)        | —                            |
| `token` | Your Jira personal access token                      | —                            |
| `url`   | Jira instance URL                                    | `https://issues.redhat.com`  |

View current configuration:

```bash
./jiracrawler config view
```

## Validate Credentials

Quickly check whether your Jira API token and user are valid:

```bash
./jiracrawler validate
```

## Get Assigned Issues

Query Jira for issues assigned to one or more users, filtered by project:

```bash
./jiracrawler get assignedissues <user1> <user2> --projectID CNF --output json
```

| Flag        | Description                        | Default |
|-------------|------------------------------------|---------|
| `projectID` | Jira project key                   | `CNF`   |
| `output`    | Output format (`json` or `yaml`)   | —       |

## Get User Updates in Date Range

Find all issues assigned to a user that were updated within a specific date range:

```bash
./jiracrawler get userupdates <user@example.com> <start-date> <end-date> --output json
```

| Argument       | Description                         |
|----------------|-------------------------------------|
| `user`         | Email of the user to query          |
| `start-date`   | Start date in `YYYY-MM-DD` format   |
| `end-date`     | End date in `YYYY-MM-DD` format     |

| Flag     | Description                        | Default |
|----------|------------------------------------|---------|
| `output` | Output format (`json` or `yaml`)   | —       |

### Example

```bash
./jiracrawler get userupdates user@redhat.com 2024-01-01 2024-01-31 --output json
```
