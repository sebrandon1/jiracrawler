package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/jiracrawler/lib"
	"github.com/spf13/cobra"
)

// validateConfig checks that required Jira configuration values are set and returns them.
func validateConfig() (apikey, jiraURL, jiraUser string) {
	apikey = GetConfigValue("apikey")
	jiraURL = GetConfigValue("jira_url")
	jiraUser = GetConfigValue("jira_user")

	if apikey == "" || jiraURL == "" || jiraUser == "" {
		fmt.Println("Jira API key, URL, and user must be set in the config.")
		os.Exit(1)
	}
	return
}

// printOutput formats and prints data in the specified output format.
func printOutput(format string, data interface{}) error {
	switch format {
	case "yaml":
		return lib.PrintYAML(data)
	case "table":
		return lib.PrintTable(data)
	default:
		return lib.PrintJSON(data)
	}
}

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
		apikey, jiraURL, jiraUser := validateConfig()
		users := args
		results, err := lib.FetchAssignedIssuesWithProject(jiraURL, jiraUser, apikey, projectID, users)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := printOutput(output, results); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
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

		apikey, jiraURL, jiraUser := validateConfig()
		_ = jiraUser

		assignee := args[0]
		startDate := args[1]
		endDate := args[2]

		result, err := lib.FetchUserIssuesInDateRange(jiraURL, jiraUser, apikey, assignee, startDate, endDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := printOutput(output, result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(assignedIssuesCmd)
	getCmd.AddCommand(userUpdatesCmd)

	assignedIssuesCmd.Flags().StringP("output", "o", "json", "Output format: json|yaml|table")
	assignedIssuesCmd.PersistentFlags().StringP("projectID", "p", "CNF", "Jira project key (e.g., CNF)")

	userUpdatesCmd.Flags().StringP("output", "o", "json", "Output format: json|yaml|table")

	// Ensure getCmd and assignedIssuesCmd are initialized for root.go
}
