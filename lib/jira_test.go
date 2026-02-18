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

func TestFetchIssuesWithJQL_EmptyURL(t *testing.T) {
	result, err := FetchIssuesWithJQL("", "token", "project = CNF", 50)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "jiraURL and apikey must be provided")
}

func TestFetchIssuesWithJQL_EmptyJQL(t *testing.T) {
	result, err := FetchIssuesWithJQL("https://example.com", "token", "", 50)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "JQL query must not be empty")
}

func TestPrintTable_QueryResult(t *testing.T) {
	data := &QueryResult{
		JQL:        "project = CNF",
		TotalCount: 1,
		Issues: []Issue{
			{
				Key:      "CNF-1234",
				Summary:  "Test issue",
				Status:   Status{Name: "Open"},
				Priority: Priority{Name: "Major"},
			},
		},
	}
	err := PrintTable(data)
	assert.NoError(t, err)
}

func TestPrintTable_AssignedIssuesResult(t *testing.T) {
	data := []AssignedIssuesResult{
		{
			User: "testuser",
			Issues: []Issue{
				{
					Key:      "CNF-1234",
					Summary:  "Fix network policy",
					Status:   Status{Name: "Open"},
					Priority: Priority{Name: "Major"},
				},
				{
					Key:      "CNF-1235",
					Summary:  "Update SDK version",
					Status:   Status{Name: "In Progress"},
					Priority: Priority{Name: "Critical"},
				},
			},
		},
	}
	err := PrintTable(data)
	assert.NoError(t, err)
}

func TestPrintTable_UserUpdatesResult(t *testing.T) {
	data := &UserUpdatesResult{
		User:       "testuser",
		DateRange:  "2024-01-01 to 2024-01-31",
		TotalCount: 1,
		Issues: []Issue{
			{
				Key:      "CNF-1234",
				Summary:  "Fix network policy",
				Status:   Status{Name: "Closed"},
				Priority: Priority{Name: "Minor"},
			},
		},
	}
	err := PrintTable(data)
	assert.NoError(t, err)
}

func TestPrintTable_NilUserUpdatesResult(t *testing.T) {
	var data *UserUpdatesResult
	err := PrintTable(data)
	assert.NoError(t, err)
}

func TestPrintTable_UnsupportedType(t *testing.T) {
	err := PrintTable("unsupported")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported data type")
}

func TestPrintTable_LongSummaryTruncation(t *testing.T) {
	data := []AssignedIssuesResult{
		{
			User: "testuser",
			Issues: []Issue{
				{
					Key:      "CNF-1234",
					Summary:  "This is a very long summary that exceeds sixty characters and should be truncated",
					Status:   Status{Name: "Open"},
					Priority: Priority{Name: "Major"},
				},
			},
		},
	}
	err := PrintTable(data)
	assert.NoError(t, err)
}
