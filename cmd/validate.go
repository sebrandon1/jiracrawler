package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate your Jira credentials",
	Run: func(cmd *cobra.Command, args []string) {
		jiraURL := GetConfigValue("jira_url")
		apikey := GetConfigValue("apikey")
		if jiraURL == "" {
			jiraURL = "https://issues.redhat.com"
		}
		if apikey == "" {
			fmt.Println("Missing apikey in config.")
			return
		}
		url := jiraURL + "/rest/api/2/myself"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Authorization", "Bearer "+apikey)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		if resp.StatusCode == 200 {
			fmt.Println("Jira credentials are valid!")
		} else {
			fmt.Printf("Jira credentials are invalid. Status: %s\n", resp.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
