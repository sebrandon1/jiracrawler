package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryCmdStructure(t *testing.T) {
	assert.Equal(t, "query [jql]", queryCmd.Use)
	assert.Equal(t, "Run a custom JQL query", queryCmd.Short)

	// Verify flags
	outputFlag := queryCmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "o", outputFlag.Shorthand)
	assert.Equal(t, "json", outputFlag.DefValue)

	maxResultsFlag := queryCmd.Flags().Lookup("max-results")
	assert.NotNil(t, maxResultsFlag)
	assert.Equal(t, "m", maxResultsFlag.Shorthand)
	assert.Equal(t, "50", maxResultsFlag.DefValue)
}

func TestQueryCmdRegistered(t *testing.T) {
	subcommands := getCmd.Commands()
	names := make([]string, 0, len(subcommands))
	for _, cmd := range subcommands {
		names = append(names, cmd.Name())
	}
	assert.Contains(t, names, "query")
}
