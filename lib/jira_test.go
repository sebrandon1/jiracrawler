package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchAssignedIssuesWithProject_EmptyConfig(t *testing.T) {
	issues := FetchAssignedIssuesWithProject("", "", "", "CNF", []string{"testuser"})
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
