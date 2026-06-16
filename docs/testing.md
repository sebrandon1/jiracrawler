# Testing

## Unit Tests

```bash
make test
```

## Integration Tests

Integration tests run against a real Jira instance using your configured credentials:

```bash
make integration-test
```

For full testing, set the `JIRACRAWLER_TEST_USER` environment variable:

```bash
export JIRACRAWLER_TEST_USER=your-email@example.com
make integration-test
```

## Nightly Integration Tests (CI/CD)

The repository includes a GitHub Actions workflow that runs integration tests nightly against a real Jira instance. This workflow:

- Runs every night at 2 AM UTC
- Only executes on the upstream repository (forks are blocked for security)
- Can be manually triggered via workflow dispatch
- Requires the following secrets to be configured in the repository settings

### Required Secrets

| Secret             | Description                                                  |
|--------------------|--------------------------------------------------------------|
| `JIRA_URL`         | Your Jira instance URL (e.g., `https://issues.redhat.com`)  |
| `JIRA_USER`        | Jira username/email for authentication and testing           |
| `JIRA_API_TOKEN`   | Personal access token for Jira API                           |
| `JIRA_TEST_PROJECT`| (Optional) Project to test against (defaults to `CNF`)       |

The workflow automatically builds the binary, configures credentials, runs all integration tests, and cleans up sensitive data afterward.
