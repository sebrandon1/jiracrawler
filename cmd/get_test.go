package cmd

import (
	"testing"

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

func TestGetCmdHasSubcommands(t *testing.T) {
	subcommands := getCmd.Commands()
	names := make([]string, 0, len(subcommands))
	for _, cmd := range subcommands {
		names = append(names, cmd.Name())
	}
	assert.Contains(t, names, "assignedissues")
	assert.Contains(t, names, "userupdates")
}
