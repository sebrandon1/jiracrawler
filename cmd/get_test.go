package cmd

import (
	"testing"

	"github.com/sebrandon1/jiracrawler/lib"
	"github.com/stretchr/testify/assert"
)

func TestGetCmdStructure(t *testing.T) {
	assert.Equal(t, "get", getCmd.Use)
	assert.Equal(t, "Get Jira data", getCmd.Short)
}

func TestAssignedIssuesCmdStructure(t *testing.T) {
	assert.Equal(t, "assignedissues [users...]", assignedIssuesCmd.Use)
	assert.Equal(t, "Get assigned issues for users", assignedIssuesCmd.Short)

	// Verify flags
	outputFlag := assignedIssuesCmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "o", outputFlag.Shorthand)
	assert.Equal(t, "json", outputFlag.DefValue)

	projectFlag := assignedIssuesCmd.PersistentFlags().Lookup("projectID")
	assert.NotNil(t, projectFlag)
	assert.Equal(t, "p", projectFlag.Shorthand)
	assert.Equal(t, "CNF", projectFlag.DefValue)
}

func TestUserUpdatesCmdStructure(t *testing.T) {
	assert.Equal(t, "userupdates [user] [start-date] [end-date]", userUpdatesCmd.Use)
	assert.Equal(t, "Get issues assigned to a user that were updated within a date range", userUpdatesCmd.Short)

	// Verify flags
	outputFlag := userUpdatesCmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "o", outputFlag.Shorthand)
	assert.Equal(t, "json", outputFlag.DefValue)
}

func TestPrintOutput_JSON(t *testing.T) {
	data := map[string]string{"key": "value"}
	err := printOutput("json", data)
	assert.NoError(t, err)
}

func TestPrintOutput_YAML(t *testing.T) {
	data := map[string]string{"key": "value"}
	err := printOutput("yaml", data)
	assert.NoError(t, err)
}

func TestPrintOutput_Table(t *testing.T) {
	data := &lib.UserUpdatesResult{
		User:       "test",
		TotalCount: 0,
		Issues:     []lib.Issue{},
	}
	err := printOutput("table", data)
	assert.NoError(t, err)
}

func TestPrintOutput_DefaultIsJSON(t *testing.T) {
	data := map[string]string{"key": "value"}
	err := printOutput("", data)
	assert.NoError(t, err)
}

func TestGetCmdHasSubcommands(t *testing.T) {
	subcommands := getCmd.Commands()
	names := make([]string, 0, len(subcommands))
	for _, cmd := range subcommands {
		names = append(names, cmd.Name())
	}
	assert.Contains(t, names, "assignedissues")
	assert.Contains(t, names, "userupdates")
}
