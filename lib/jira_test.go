package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchAssignedIssuesWithProject_EmptyConfig(t *testing.T) {
	results, err := FetchAssignedIssuesWithProject("", "", "", "CNF", []string{"testuser"})
	assert.Error(t, err)
	assert.Nil(t, results)
}

func TestFetchUserIssuesInDateRange_EmptyConfig(t *testing.T) {
	result, err := FetchUserIssuesInDateRange("", "", "", "testuser", "2024-01-01", "2024-01-31")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFetchUserIssuesInDateRange_EmptyUser(t *testing.T) {
	result, err := FetchUserIssuesInDateRange("https://example.com", "user", "token", "", "2024-01-01", "2024-01-31")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFetchUserIssuesInDateRange_EmptyDates(t *testing.T) {
	result, err := FetchUserIssuesInDateRange("https://example.com", "user", "token", "testuser", "", "2024-01-31")
	assert.Error(t, err)
	assert.Nil(t, result)

	result, err = FetchUserIssuesInDateRange("https://example.com", "user", "token", "testuser", "2024-01-01", "")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFetchUserIssuesInDateRange_InvalidDateFormat(t *testing.T) {
	result, err := FetchUserIssuesInDateRange("https://example.com", "user", "token", "testuser", "invalid-date", "2024-01-31")
	assert.Error(t, err)
	assert.Nil(t, result)

	result, err = FetchUserIssuesInDateRange("https://example.com", "user", "token", "testuser", "2024-01-01", "invalid-date")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPrintJSON(t *testing.T) {
	data := map[string]interface{}{"foo": "bar"}
	err := PrintJSON(data)
	assert.NoError(t, err)
}

func TestPrintYAML(t *testing.T) {
	data := map[string]interface{}{"foo": "bar"}
	err := PrintYAML(data)
	assert.NoError(t, err)
}
