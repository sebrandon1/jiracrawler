package lib

import (
	"encoding/json"
	"fmt"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"gopkg.in/yaml.v3"
)

// JiraConfig holds the Jira instance URL and user email
// These should be set in the config as 'jira_url' and 'jira_email'
type JiraConfig struct {
	URL    string
	APIKey string
}

// FetchAssignedIssuesWithProject fetches assigned issues for the given users and project from Jira
func FetchAssignedIssuesWithProject(jiraURL, jiraUser, apikey string, projectID string, users []string) []map[string]interface{} {
	if jiraURL == "" || jiraUser == "" || projectID == "" {
		fmt.Println("jiraURL, jiraUser, and projectID must be provided (no defaults in lib)")
		return nil
	}
	tokenAuth := jira.BearerAuthTransport{
		Token: apikey,
	}
	client, err := jira.NewClient(tokenAuth.Client(), jiraURL)
	if err != nil {
		fmt.Printf("Error creating Jira client: %v\n", err)
		return nil
	}
	var allIssues []map[string]interface{}
	for _, user := range users {
		jql := fmt.Sprintf("project=%s AND assignee=\"%s\" ORDER BY created DESC", projectID, user)
		issues, resp, err := client.Issue.Search(jql, nil)
		if err != nil {
			fmt.Printf("Error fetching issues for %s: %v\n", user, err)
			continue
		}
		if resp != nil && resp.StatusCode != 200 {
			fmt.Printf("Jira API error for %s: %s\n", user, resp.Status)
			continue
		}
		var issuesList []map[string]interface{}
		for _, issue := range issues {
			issueMap := map[string]interface{}{
				"key":    issue.Key,
				"fields": issue.Fields,
			}
			issuesList = append(issuesList, issueMap)
		}
		allIssues = append(allIssues, map[string]interface{}{
			"user":   user,
			"issues": issuesList,
		})
	}
	return allIssues
}

// FetchUserIssuesInDateRange fetches issues assigned to a user that were updated within a specific date range
// startDate and endDate should be in YYYY-MM-DD format
func FetchUserIssuesInDateRange(jiraURL, jiraUser, apikey string, assignee string, startDate, endDate string) []map[string]interface{} {
	if jiraURL == "" || jiraUser == "" || assignee == "" || startDate == "" || endDate == "" {
		fmt.Println("jiraURL, jiraUser, assignee, startDate, and endDate must be provided")
		return nil
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		fmt.Printf("Invalid start date format. Use YYYY-MM-DD: %v\n", err)
		return nil
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		fmt.Printf("Invalid end date format. Use YYYY-MM-DD: %v\n", err)
		return nil
	}

	tokenAuth := jira.BearerAuthTransport{
		Token: apikey,
	}
	client, err := jira.NewClient(tokenAuth.Client(), jiraURL)
	if err != nil {
		fmt.Printf("Error creating Jira client: %v\n", err)
		return nil
	}

	// JQL query to find issues assigned to the user that were updated in the date range
	jql := fmt.Sprintf("assignee=\"%s\" AND updated >= \"%s\" AND updated <= \"%s\" ORDER BY updated DESC", assignee, startDate, endDate)

	issues, resp, err := client.Issue.Search(jql, nil)
	if err != nil {
		fmt.Printf("Error fetching issues for %s: %v\n", assignee, err)
		return nil
	}
	if resp != nil && resp.StatusCode != 200 {
		fmt.Printf("Jira API error for %s: %s\n", assignee, resp.Status)
		return nil
	}

	var issuesList []map[string]interface{}
	for _, issue := range issues {
		issueMap := map[string]interface{}{
			"key":    issue.Key,
			"fields": issue.Fields,
		}
		issuesList = append(issuesList, issueMap)
	}

	return []map[string]interface{}{
		{
			"user":       assignee,
			"dateRange":  fmt.Sprintf("%s to %s", startDate, endDate),
			"totalCount": len(issuesList),
			"issues":     issuesList,
		},
	}
}

func PrintJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(b))
}

func PrintYAML(data interface{}) {
	b, err := yaml.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling YAML:", err)
		return
	}
	fmt.Println(string(b))
}
