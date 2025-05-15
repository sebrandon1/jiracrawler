package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		apikey := viper.GetString("apikey")
		jiraURL := viper.GetString("jira_url")
		jiraEmail := viper.GetString("jira_email")
		if apikey == "" || jiraURL == "" || jiraEmail == "" {
			fmt.Println("Jira API key, URL, and email must be set. Use 'jiracrawler config set apikey <key>', 'jiracrawler config set jira_url <url>', and 'jiracrawler config set jira_email <email>' to set them.")
			os.Exit(1)
		}
		users := args
		issues := FetchAssignedIssues(apikey, users)
		if output == "yaml" {
			PrintYAML(issues)
		} else {
			PrintJSON(issues)
		}
	},
}

func init() {
	getCmd.AddCommand(assignedIssuesCmd)
	assignedIssuesCmd.Flags().StringP("output", "o", "json", "Output format: json|yaml")
	// Ensure getCmd and assignedIssuesCmd are initialized for root.go
}
