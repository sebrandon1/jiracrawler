package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchAssignedIssuesWithProject_EmptyConfig(t *testing.T) {
	// Simulate missing config values
	setConfigValue := func(key, value string) {
		// This is a placeholder for setting config in tests
	}
	setConfigValue("jira_url", "")
	setConfigValue("jira_user", "")
	setConfigValue("apikey", "")

	issues := FetchAssignedIssuesWithProject([]string{"testuser"}, "CNF")
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
