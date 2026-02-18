package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompletionCmdStructure(t *testing.T) {
	assert.Equal(t, "completion [bash|zsh|fish|powershell]", completionCmd.Use)
	assert.Equal(t, "Generate shell completion scripts", completionCmd.Short)
	assert.Equal(t, []string{"bash", "zsh", "fish", "powershell"}, completionCmd.ValidArgs)
	assert.True(t, completionCmd.DisableFlagsInUseLine)
}
