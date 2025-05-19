package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/jiracrawler/lib"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Jira data",
}

var assignedIssuesCmd = &cobra.Command{
	Use:   "assignedissues [users...]",
	Short: "Get assigned issues for users",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		projectID, _ := cmd.Flags().GetString("projectID")
		if projectID == "" {
			projectID = "CNF"
		}
		apikey := GetConfigValue("apikey")
		jiraURL := GetConfigValue("jira_url")
		jiraUser := GetConfigValue("jira_user")

		if apikey == "" || jiraURL == "" || jiraUser == "" {
			fmt.Println("Jira API key, URL, and user must be set in the config.")
			os.Exit(1)
		}
		users := args
		issues := lib.FetchAssignedIssuesWithProject(jiraURL, jiraUser, apikey, projectID, users)
		if output == "yaml" {
			lib.PrintYAML(issues)
		} else {
			lib.PrintJSON(issues)
		}
	},
}

func init() {
	getCmd.AddCommand(assignedIssuesCmd)
	assignedIssuesCmd.Flags().StringP("output", "o", "json", "Output format: json|yaml")
	assignedIssuesCmd.PersistentFlags().StringP("projectID", "p", "CNF", "Jira project key (e.g., CNF)")
	// Ensure getCmd and assignedIssuesCmd are initialized for root.go
}
