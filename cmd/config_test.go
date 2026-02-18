package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigValue(t *testing.T) {
	viper.Reset()
	viper.Set("test_key", "test_value")
	assert.Equal(t, "test_value", GetConfigValue("test_key"))
}

func TestGetConfigValue_Empty(t *testing.T) {
	viper.Reset()
	assert.Equal(t, "", GetConfigValue("nonexistent_key"))
}

func TestDefaultJiraURL(t *testing.T) {
	assert.Equal(t, "https://issues.redhat.com", DefaultJiraURL)
}

func TestConfigCmdStructure(t *testing.T) {
	assert.Equal(t, "config", configCmd.Use)
	assert.Equal(t, "Set or get configuration values", configCmd.Short)
}

func TestSetCmdStructure(t *testing.T) {
	assert.Equal(t, "set", setCmd.Use)
	assert.Equal(t, "Set a key value pair to the configuration", setCmd.Short)

	// Verify flags are registered
	tokenFlag := setCmd.PersistentFlags().Lookup("token")
	assert.NotNil(t, tokenFlag)
	assert.Equal(t, "t", tokenFlag.Shorthand)

	urlFlag := setCmd.PersistentFlags().Lookup("url")
	assert.NotNil(t, urlFlag)
	assert.Equal(t, "u", urlFlag.Shorthand)

	userFlag := setCmd.PersistentFlags().Lookup("user")
	assert.NotNil(t, userFlag)
	assert.Equal(t, "s", userFlag.Shorthand)
}

func TestViewCmdStructure(t *testing.T) {
	assert.Equal(t, "view", viewCmd.Use)
	assert.Equal(t, "View the configuration", viewCmd.Short)
}
