package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchAssignedIssuesWithProject_EmptyConfig(t *testing.T) {
	issues := FetchAssignedIssuesWithProject("", "", "", "CNF", []string{"testuser"})
	assert.Nil(t, issues)
}

func TestFetchUserIssuesInDateRange_EmptyConfig(t *testing.T) {
	issues := FetchUserIssuesInDateRange("", "", "", "testuser", "2024-01-01", "2024-01-31")
	assert.Nil(t, issues)
}

func TestFetchUserIssuesInDateRange_EmptyUser(t *testing.T) {
	issues := FetchUserIssuesInDateRange("https://example.com", "user", "token", "", "2024-01-01", "2024-01-31")
	assert.Nil(t, issues)
}

func TestFetchUserIssuesInDateRange_EmptyDates(t *testing.T) {
	issues := FetchUserIssuesInDateRange("https://example.com", "user", "token", "testuser", "", "2024-01-31")
	assert.Nil(t, issues)

	issues = FetchUserIssuesInDateRange("https://example.com", "user", "token", "testuser", "2024-01-01", "")
	assert.Nil(t, issues)
}

func TestFetchUserIssuesInDateRange_InvalidDateFormat(t *testing.T) {
	issues := FetchUserIssuesInDateRange("https://example.com", "user", "token", "testuser", "invalid-date", "2024-01-31")
	assert.Nil(t, issues)

	issues = FetchUserIssuesInDateRange("https://example.com", "user", "token", "testuser", "2024-01-01", "invalid-date")
	assert.Nil(t, issues)
}

func TestPrintJSON(t *testing.T) {
	data := map[string]interface{}{"foo": "bar"}
	PrintJSON(data)
	// No assertion, just ensure no panic
}

func TestPrintYAML(t *testing.T) {
	data := map[string]interface{}{"foo": "bar"}
	PrintYAML(data)
	// No assertion, just ensure no panic
}
