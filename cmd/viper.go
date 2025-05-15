package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// JiraConfig holds the Jira instance URL and user email
// These should be set in the config as 'jira_url' and 'jira_email'
type JiraConfig struct {
	URL    string
	Email  string
	APIKey string
}

// FetchAssignedIssues fetches assigned issues for the given users from Jira
func FetchAssignedIssues(apikey string, users []string) []map[string]interface{} {
	jiraURL := viper.GetString("jira_url")
	jiraEmail := viper.GetString("jira_email")
	if jiraURL == "" || jiraEmail == "" {
		fmt.Println("Jira URL and email must be set in config.")
		return nil
	}
	var allIssues []map[string]interface{}
	for _, user := range users {
		issues, err := fetchUserIssues(jiraURL, jiraEmail, apikey, user)
		if err != nil {
			fmt.Printf("Error fetching issues for %s: %v\n", user, err)
			continue
		}
		allIssues = append(allIssues, map[string]interface{}{
			"user":   user,
			"issues": issues,
		})
	}
	return allIssues
}

func fetchUserIssues(jiraURL, jiraEmail, apikey, user string) ([]map[string]interface{}, error) {
	apiEndpoint := fmt.Sprintf("%s/rest/api/3/search", strings.TrimRight(jiraURL, "/"))
	jql := url.QueryEscape(fmt.Sprintf("assignee=\"%s\" AND resolution=Unresolved ORDER BY updated DESC", user))
	fullURL := fmt.Sprintf("%s?jql=%s", apiEndpoint, jql)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(jiraEmail, apikey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Jira API error: %s", string(body))
	}

	var result struct {
		Issues []map[string]interface{} `json:"issues"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Issues, nil
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
