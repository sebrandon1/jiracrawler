package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestValidateCmdStructure(t *testing.T) {
	assert.Equal(t, "validate", validateCmd.Use)
	assert.Equal(t, "Validate your Jira credentials", validateCmd.Short)
}

func TestValidateCmd_ValidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/myself", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	viper.Reset()
	viper.Set("jira_url", server.URL)
	viper.Set("apikey", "test-token")

	var buf bytes.Buffer
	validateCmd.SetOut(&buf)
	validateCmd.SetArgs([]string{})
	validateCmd.Run(validateCmd, []string{})

	// Output goes to stdout via fmt.Println, not to the command's writer,
	// so we just verify it doesn't panic or error
}

func TestValidateCmd_InvalidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	viper.Reset()
	viper.Set("jira_url", server.URL)
	viper.Set("apikey", "bad-token")

	validateCmd.SetArgs([]string{})
	validateCmd.Run(validateCmd, []string{})

	// Verify it doesn't panic - the command prints a message about invalid credentials
}

func TestValidateCmd_MissingAPIKey(t *testing.T) {
	viper.Reset()
	viper.Set("jira_url", "https://example.com")
	// Don't set apikey

	validateCmd.SetArgs([]string{})
	validateCmd.Run(validateCmd, []string{})

	// Verify it doesn't panic - the command prints a message about missing apikey
}
