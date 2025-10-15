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

// User represents a JIRA user (assignee, reporter, creator)
type User struct {
	Key          string `json:"key" yaml:"key"`
	Name         string `json:"name" yaml:"name"`
	EmailAddress string `json:"emailAddress" yaml:"emailAddress"`
	DisplayName  string `json:"displayName" yaml:"displayName"`
	Active       bool   `json:"active" yaml:"active"`
}

// Status represents a JIRA issue status
type Status struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// Priority represents a JIRA issue priority
type Priority struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// IssueType represents a JIRA issue type
type IssueType struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// Project represents a JIRA project
type Project struct {
	ID   string `json:"id" yaml:"id"`
	Key  string `json:"key" yaml:"key"`
	Name string `json:"name" yaml:"name"`
}

// Comment represents a Jira issue comment
type Comment struct {
	ID      string    `json:"id" yaml:"id"`
	Author  string    `json:"author" yaml:"author"`
	Body    string    `json:"body" yaml:"body"`
	Created time.Time `json:"created" yaml:"created"`
	Updated time.Time `json:"updated" yaml:"updated"`
}

// HistoryItem represents a Jira issue history entry
type HistoryItem struct {
	ID      string          `json:"id" yaml:"id"`
	Author  string          `json:"author" yaml:"author"`
	Created time.Time       `json:"created" yaml:"created"`
	Items   []HistoryChange `json:"items" yaml:"items"`
}

// HistoryChange represents a specific field change in history
type HistoryChange struct {
	Field      string `json:"field" yaml:"field"`
	FieldType  string `json:"fieldType" yaml:"fieldType"`
	FromString string `json:"fromString" yaml:"fromString"`
	ToString   string `json:"toString" yaml:"toString"`
}

// TimeTracking represents time tracking information
type TimeTracking struct {
	OriginalEstimate  string `json:"originalEstimate,omitempty" yaml:"originalEstimate,omitempty"`
	RemainingEstimate string `json:"remainingEstimate,omitempty" yaml:"remainingEstimate,omitempty"`
	TimeSpent         string `json:"timeSpent,omitempty" yaml:"timeSpent,omitempty"`
}

// IssuePermissions represents what the authenticated user can access
type IssuePermissions struct {
	CanViewComments bool
	CanViewHistory  bool
	CanViewIssue    bool
}

// EnhancedFields represents additional field data for an issue
type EnhancedFields struct {
	Labels       []string
	Components   []string
	Priority     string
	IssueType    string
	TimeTracking *TimeTracking
}

// Issue represents a JIRA issue with key fields
type Issue struct {
	Key         string    `json:"key" yaml:"key"`
	Summary     string    `json:"summary" yaml:"summary"`
	Description string    `json:"description" yaml:"description"`
	Status      Status    `json:"status" yaml:"status"`
	Priority    Priority  `json:"priority" yaml:"priority"`
	IssueType   IssueType `json:"issueType" yaml:"issueType"`
	Project     Project   `json:"project" yaml:"project"`
	Assignee    *User     `json:"assignee" yaml:"assignee"`
	Reporter    *User     `json:"reporter" yaml:"reporter"`
	Creator     *User     `json:"creator" yaml:"creator"`
	Created     string    `json:"created" yaml:"created"`
	Updated     string    `json:"updated" yaml:"updated"`
	Resolved    string    `json:"resolved" yaml:"resolved"`

	// Enhanced context fields (populated on demand)
	Comments     []Comment     `json:"comments,omitempty" yaml:"comments,omitempty"`
	History      []HistoryItem `json:"history,omitempty" yaml:"history,omitempty"`
	Labels       []string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Components   []string      `json:"components,omitempty" yaml:"components,omitempty"`
	TimeTracking *TimeTracking `json:"timeTracking,omitempty" yaml:"timeTracking,omitempty"`
}

// AssignedIssuesResult represents the result of fetching assigned issues
type AssignedIssuesResult struct {
	User   string  `json:"user" yaml:"user"`
	Issues []Issue `json:"issues" yaml:"issues"`
}

// UserUpdatesResult represents the result of fetching user updates in a date range
type UserUpdatesResult struct {
	User       string  `json:"user" yaml:"user"`
	DateRange  string  `json:"dateRange" yaml:"dateRange"`
	TotalCount int     `json:"totalCount" yaml:"totalCount"`
	Issues     []Issue `json:"issues" yaml:"issues"`
}

// convertJiraUser converts a JIRA user to our User struct
func convertJiraUser(jiraUser *jira.User) *User {
	if jiraUser == nil {
		return nil
	}
	return &User{
		Key:          jiraUser.Key,
		Name:         jiraUser.Name,
		EmailAddress: jiraUser.EmailAddress,
		DisplayName:  jiraUser.DisplayName,
		Active:       jiraUser.Active,
	}
}

// convertJiraIssue converts a JIRA issue to our Issue struct
func convertJiraIssue(jiraIssue jira.Issue) Issue {
	issue := Issue{
		Key:         jiraIssue.Key,
		Summary:     jiraIssue.Fields.Summary,
		Description: jiraIssue.Fields.Description,
		Created:     time.Time(jiraIssue.Fields.Created).Format(time.RFC3339),
		Updated:     time.Time(jiraIssue.Fields.Updated).Format(time.RFC3339),
		Assignee:    convertJiraUser(jiraIssue.Fields.Assignee),
		Reporter:    convertJiraUser(jiraIssue.Fields.Reporter),
		Creator:     convertJiraUser(jiraIssue.Fields.Creator),
	}

	// Handle status
	if jiraIssue.Fields.Status != nil {
		issue.Status = Status{
			ID:   jiraIssue.Fields.Status.ID,
			Name: jiraIssue.Fields.Status.Name,
		}
	}

	// Handle priority
	if jiraIssue.Fields.Priority != nil {
		issue.Priority = Priority{
			ID:   jiraIssue.Fields.Priority.ID,
			Name: jiraIssue.Fields.Priority.Name,
		}
	}

	// Handle issue type
	if jiraIssue.Fields.Type.ID != "" {
		issue.IssueType = IssueType{
			ID:   jiraIssue.Fields.Type.ID,
			Name: jiraIssue.Fields.Type.Name,
		}
	}

	// Handle project
	if jiraIssue.Fields.Project.ID != "" {
		issue.Project = Project{
			ID:   jiraIssue.Fields.Project.ID,
			Key:  jiraIssue.Fields.Project.Key,
			Name: jiraIssue.Fields.Project.Name,
		}
	}

	// Handle resolution date
	if !time.Time(jiraIssue.Fields.Resolutiondate).IsZero() {
		issue.Resolved = time.Time(jiraIssue.Fields.Resolutiondate).Format(time.RFC3339)
	}

	return issue
}

// FetchAssignedIssuesWithProject fetches assigned issues for the given users and project from Jira
func FetchAssignedIssuesWithProject(jiraURL, jiraUser, apikey string, projectID string, users []string) []AssignedIssuesResult {
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
	var allResults []AssignedIssuesResult
	for _, user := range users {
		// Include all issues regardless of status (including resolved/closed)
		jql := fmt.Sprintf("project=%s AND assignee=\"%s\" AND (resolution is empty OR resolution is not empty) ORDER BY created DESC", projectID, user)
		issues, resp, err := client.Issue.Search(jql, nil)
		if err != nil {
			fmt.Printf("Error fetching issues for %s: %v\n", user, err)
			continue
		}
		if resp != nil && resp.StatusCode != 200 {
			fmt.Printf("Jira API error for %s: %s\n", user, resp.Status)
			continue
		}
		var convertedIssues []Issue
		for _, issue := range issues {
			convertedIssues = append(convertedIssues, convertJiraIssue(issue))
		}
		allResults = append(allResults, AssignedIssuesResult{
			User:   user,
			Issues: convertedIssues,
		})
	}
	return allResults
}

// FetchUserIssuesInDateRange fetches issues assigned to a user that were updated within a specific date range
// startDate and endDate should be in YYYY-MM-DD format
func FetchUserIssuesInDateRange(jiraURL, jiraUser, apikey string, assignee string, startDate, endDate string) *UserUpdatesResult {
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
	// Include all issues regardless of status (including resolved/closed)
	jql := fmt.Sprintf("assignee=\"%s\" AND updated >= \"%s\" AND updated <= \"%s\" AND (resolution is empty OR resolution is not empty) ORDER BY updated DESC", assignee, startDate, endDate)

	issues, resp, err := client.Issue.Search(jql, nil)
	if err != nil {
		fmt.Printf("Error fetching issues for %s: %v\n", assignee, err)
		return nil
	}
	if resp != nil && resp.StatusCode != 200 {
		fmt.Printf("Jira API error for %s: %s\n", assignee, resp.Status)
		return nil
	}

	var convertedIssues []Issue
	for _, issue := range issues {
		convertedIssues = append(convertedIssues, convertJiraIssue(issue))
	}

	return &UserUpdatesResult{
		User:       assignee,
		DateRange:  fmt.Sprintf("%s to %s", startDate, endDate),
		TotalCount: len(convertedIssues),
		Issues:     convertedIssues,
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

// FetchUserIssuesInDateRangeWithContext - Enhanced version of FetchUserIssuesInDateRange
// Adds optional parameter to fetch enhanced context for all issues
func FetchUserIssuesInDateRangeWithContext(
	jiraURL, jiraUser, apikey, userEmail, startDate, endDate string,
	enhancedContext bool,
	verbose bool,
) *UserUpdatesResult {
	// First, get basic issues using existing function
	result := FetchUserIssuesInDateRange(jiraURL, jiraUser, apikey, userEmail, startDate, endDate)
	if result == nil {
		return nil
	}

	// If enhanced context not requested, return as-is
	if !enhancedContext {
		return result
	}

	// Create client for enhanced fetching
	tokenAuth := jira.BearerAuthTransport{
		Token: apikey,
	}
	client, err := jira.NewClient(tokenAuth.Client(), jiraURL)
	if err != nil {
		fmt.Printf("Warning: failed to create client for enhanced context: %v\n", err)
		return result
	}

	// Enhance each issue with additional context
	for i, issue := range result.Issues {
		enhanced, err := FetchIssueWithEnhancedContext(client, jiraURL, issue.Key, apikey, verbose)
		if err != nil {
			fmt.Printf("Warning: failed to enhance issue %s: %v\n", issue.Key, err)
			continue
		}
		result.Issues[i] = *enhanced

		// Rate limiting: sleep between requests to avoid 429 errors
		if i < len(result.Issues)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return result
}
