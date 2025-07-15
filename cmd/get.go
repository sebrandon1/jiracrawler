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
		results := lib.FetchAssignedIssuesWithProject(jiraURL, jiraUser, apikey, projectID, users)
		if output == "yaml" {
			lib.PrintYAML(results)
		} else {
			lib.PrintJSON(results)
		}
	},
}

var userUpdatesCmd = &cobra.Command{
	Use:   "userupdates [user] [start-date] [end-date]",
	Short: "Get issues assigned to a user that were updated within a date range",
	Long: `Get all issues assigned to a specific user that were updated within the specified date range.

Dates should be in YYYY-MM-DD format.

Example:
  jiracrawler get userupdates user@example.com 2024-01-01 2024-01-31`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")

		apikey := GetConfigValue("apikey")
		jiraURL := GetConfigValue("jira_url")
		jiraUser := GetConfigValue("jira_user")

		if apikey == "" || jiraURL == "" || jiraUser == "" {
			fmt.Println("Jira API key, URL, and user must be set in the config.")
			os.Exit(1)
		}

				assignee := args[0]
		startDate := args[1]
		endDate := args[2]

		result := lib.FetchUserIssuesInDateRange(jiraURL, jiraUser, apikey, assignee, startDate, endDate)
		
		if result == nil {
			fmt.Println("Failed to fetch issues. Please check your parameters and try again.")
			os.Exit(1)
		}

		if output == "yaml" {
			lib.PrintYAML(result)
		} else {
			lib.PrintJSON(result)
		}
	},
}

func init() {
	getCmd.AddCommand(assignedIssuesCmd)
	getCmd.AddCommand(userUpdatesCmd)

	assignedIssuesCmd.Flags().StringP("output", "o", "json", "Output format: json|yaml")
	assignedIssuesCmd.PersistentFlags().StringP("projectID", "p", "CNF", "Jira project key (e.g., CNF)")

	userUpdatesCmd.Flags().StringP("output", "o", "json", "Output format: json|yaml")

	// Ensure getCmd and assignedIssuesCmd are initialized for root.go
}
