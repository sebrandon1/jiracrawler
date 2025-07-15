# Integration Test CI/CD Setup

This document describes how to configure the nightly integration tests workflow for the jiracrawler repository.

## Overview

The `integration-nightly.yml` workflow automatically runs integration tests against a real JIRA instance every night at 2 AM UTC. It includes multiple security safeguards to prevent unauthorized access to secrets.

## Security Features

### Fork Protection
- **Repository Verification**: Only runs on `sebrandon1/jiracrawler` 
- **Secret Isolation**: Secrets are never exposed to forked repositories
- **Multi-step Validation**: Security check job runs first and gates the integration tests

### Credential Handling
- **Temporary Configuration**: JIRA credentials are only set during test execution
- **Automatic Cleanup**: Config files and binaries are removed after tests complete
- **No Persistence**: Credentials are never committed or stored in the repository

## Required Repository Secrets

Configure these secrets in **Settings > Secrets and variables > Actions**:

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `JIRA_URL` | JIRA instance URL | `https://issues.redhat.com` |
| `JIRA_USER` | JIRA username/email (used for auth and testing) | `user@redhat.com` |
| `JIRA_API_TOKEN` | Personal access token | `ATATT3xFfGF0...` |
| `JIRA_TEST_PROJECT` | (Optional) Project to test | `CNF` |

## How to Generate JIRA API Token

1. Go to your JIRA instance (e.g., https://issues.redhat.com)
2. Click your profile → **Account settings**
3. Navigate to **Security** → **Create and manage API tokens**
4. Click **Create API token**
5. Give it a descriptive name like "jiracrawler-ci"
6. Copy the generated token immediately (it won't be shown again)

## Workflow Triggers

### Automatic (Nightly)
- **Schedule**: Every day at 2:00 AM UTC
- **Cron**: `0 2 * * *`

### Manual
- Go to **Actions** tab in GitHub
- Select **Nightly Integration Tests**
- Click **Run workflow**
- Choose the branch (usually `main`)

## Test Coverage

The integration workflow tests:

1. **Credentials Validation** - Verifies JIRA API connectivity
2. **Assigned Issues Retrieval** - Tests existing functionality  
3. **User Updates in Date Range** - Tests new feature
4. **Error Handling** - Validates input validation
5. **Help Commands** - Ensures CLI help works
6. **Output Formats** - Tests both JSON and YAML output

## Monitoring

### Success Indicators
- ✅ All jobs complete successfully
- ✅ Integration tests find and process real JIRA data
- ✅ JSON output validation passes

### Failure Scenarios
- ❌ Repository is a fork (expected security behavior)
- ❌ Required secrets are missing or invalid
- ❌ JIRA API is unreachable
- ❌ Test user has no accessible issues
- ❌ Output format validation fails

## Troubleshooting

### Missing Secrets Error
```
❌ Required secrets not set. Please configure:
  - JIRA_URL (e.g., https://issues.redhat.com)
  - JIRA_USER (e.g., user@redhat.com)
  - JIRA_API_TOKEN
```
**Solution**: Add the missing secrets in repository settings.

### Fork Security Block
```
⚠️ Running on fork: username/jiracrawler, skipping integration tests
```
**Solution**: This is expected behavior. Only the upstream repository should run integration tests.

### JIRA Authentication Failed
```
❌ Credentials validation failed
```
**Solution**: Verify the JIRA_API_TOKEN is valid and the user has proper permissions.

### No Test Data Found
```
Found 0 issues updated in the date range
```
**Solution**: Verify JIRA_USER has assigned issues or recent activity.

## Maintenance

### Token Rotation
- JIRA API tokens should be rotated periodically
- Update the `JIRA_API_TOKEN` secret when rotating
- Test the workflow manually after rotation

### User Changes
- If the test user changes, update `JIRA_USER`
- Ensure the user has sufficient JIRA activity for meaningful tests

### Schedule Adjustments
- Modify the cron expression in `integration-nightly.yml`
- Consider timezone differences for your team 