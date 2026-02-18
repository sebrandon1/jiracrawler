package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmdStructure(t *testing.T) {
	assert.Equal(t, "jiracrawler", rootCmd.Use)
	assert.Equal(t, "Jira issue crawler CLI", rootCmd.Short)
}

func TestExecuteReturnsNoError(t *testing.T) {
	// Execute with --help to avoid needing a real config/server
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}
