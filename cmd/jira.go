package cmd

import (
	"encoding/json"
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"gopkg.in/yaml.v3"
)

// JiraConfig holds the Jira instance URL and user email
// These should be set in the config as 'jira_url' and 'jira_email'
type JiraConfig struct {
	URL    string
	APIKey string
}

// FetchAssignedIssues fetches assigned issues for the given users from Jira
func FetchAssignedIssues(users []string) []map[string]interface{} {
	jiraURL := GetConfigValue("jira_url")
	jiraUser := GetConfigValue("jira_user")
	apikey := GetConfigValue("apikey")
	if jiraURL == "" {
		jiraURL = "https://issues.redhat.com"
	}
	if jiraUser == "" {
		fmt.Println("Jira username must be set in config as 'jira_user'.")
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
		// JQL: Only issues in the CNF project assigned to the user, not Closed, ordered by created date
		jql := fmt.Sprintf("project=CNF AND assignee=\"%s\" ORDER BY created DESC", user)
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

// FetchAssignedIssuesWithProject fetches assigned issues for the given users and project from Jira
func FetchAssignedIssuesWithProject(users []string, projectID string) []map[string]interface{} {
	jiraURL := GetConfigValue("jira_url")
	jiraUser := GetConfigValue("jira_user")
	apikey := GetConfigValue("apikey")
	if jiraURL == "" {
		jiraURL = "https://issues.redhat.com"
	}
	if jiraUser == "" {
		fmt.Println("Jira username must be set in config as 'jira_user'.")
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
