package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/jiracrawler/lib"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query [jql]",
	Short: "Run a custom JQL query",
	Long: `Run an arbitrary JQL query against the configured Jira instance.

The JQL string should be quoted. Results are returned in the specified output format.

Examples:
  jiracrawler get query "project = CNF AND status = Open ORDER BY priority DESC"
  jiracrawler get query "assignee = currentUser() AND resolution = Unresolved" -o table
  jiracrawler get query "project = CNF AND created >= -7d" --max-results 10 -o yaml`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		maxResults, _ := cmd.Flags().GetInt("max-results")

		apikey := GetConfigValue("apikey")
		jiraURL := GetConfigValue("jira_url")

		if apikey == "" || jiraURL == "" {
			fmt.Fprintln(os.Stderr, "Jira API key and URL must be set in the config.")
			os.Exit(1)
		}

		jql := args[0]
		result, err := lib.FetchIssuesWithJQL(jiraURL, apikey, jql, maxResults)
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
	getCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringP("output", "o", "json", "Output format: json|yaml|table")
	queryCmd.Flags().IntP("max-results", "m", 50, "Maximum number of results to return")
}
