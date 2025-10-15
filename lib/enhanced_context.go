package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	jira "github.com/andygrunwald/go-jira"
)

// getHTTPClient extracts the underlying HTTP client from jira.Client
// The jira.Client uses an internal HTTP client with the configured transport
func getHTTPClient(apikey string) *http.Client {
	tokenAuth := jira.BearerAuthTransport{
		Token: apikey,
	}
	return tokenAuth.Client()
}

// FetchIssueComments retrieves all comments for a specific issue
// Uses Bearer token authentication from existing jira.Client
func FetchIssueComments(client *jira.Client, baseURL, issueKey, apikey string) ([]Comment, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s/comment", baseURL, issueKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Bearer auth is handled by the client's transport
	req.Header.Set("Accept", "application/json")

	// Use the authenticated HTTP client
	httpClient := getHTTPClient(apikey)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jira API returned status %d: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var response struct {
		Comments []struct {
			ID     string `json:"id"`
			Body   string `json:"body"`
			Author struct {
				DisplayName string `json:"displayName"`
			} `json:"author"`
			Created string `json:"created"`
			Updated string `json:"updated"`
		} `json:"comments"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode comments response: %w", err)
	}

	var comments []Comment
	for _, c := range response.Comments {
		created, _ := time.Parse("2006-01-02T15:04:05.000-0700", c.Created)
		updated, _ := time.Parse("2006-01-02T15:04:05.000-0700", c.Updated)

		comments = append(comments, Comment{
			ID:      c.ID,
			Body:    c.Body,
			Author:  c.Author.DisplayName,
			Created: created,
			Updated: updated,
		})
	}

	return comments, nil
}

// FetchIssueHistory retrieves the change history for a specific issue
func FetchIssueHistory(client *jira.Client, baseURL, issueKey, apikey string) ([]HistoryItem, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s?expand=changelog", baseURL, issueKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	httpClient := getHTTPClient(apikey)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch history: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jira API returned status %d: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var response struct {
		Changelog struct {
			Histories []struct {
				ID     string `json:"id"`
				Author struct {
					DisplayName string `json:"displayName"`
				} `json:"author"`
				Created string `json:"created"`
				Items   []struct {
					Field      string `json:"field"`
					FieldType  string `json:"fieldtype"`
					FromString string `json:"fromString"`
					ToString   string `json:"toString"`
				} `json:"items"`
			} `json:"histories"`
		} `json:"changelog"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode history response: %w", err)
	}

	var history []HistoryItem
	for _, h := range response.Changelog.Histories {
		created, _ := time.Parse("2006-01-02T15:04:05.000-0700", h.Created)

		var items []HistoryChange
		for _, item := range h.Items {
			items = append(items, HistoryChange{
				Field:      item.Field,
				FieldType:  item.FieldType,
				FromString: item.FromString,
				ToString:   item.ToString,
			})
		}

		history = append(history, HistoryItem{
			ID:      h.ID,
			Author:  h.Author.DisplayName,
			Created: created,
			Items:   items,
		})
	}

	return history, nil
}

// FetchEnhancedFields retrieves additional field data for an issue
func FetchEnhancedFields(client *jira.Client, baseURL, issueKey, apikey string) (*EnhancedFields, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s?fields=labels,components,priority,issuetype,timetracking",
		baseURL, issueKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	httpClient := getHTTPClient(apikey)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch enhanced fields: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jira API returned status %d: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var response struct {
		Fields struct {
			Labels     []string `json:"labels"`
			Components []struct {
				Name string `json:"name"`
			} `json:"components"`
			Priority struct {
				Name string `json:"name"`
			} `json:"priority"`
			IssueType struct {
				Name string `json:"name"`
			} `json:"issuetype"`
			TimeTracking struct {
				OriginalEstimate  string `json:"originalEstimate"`
				RemainingEstimate string `json:"remainingEstimate"`
				TimeSpent         string `json:"timeSpent"`
			} `json:"timetracking"`
		} `json:"fields"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode enhanced fields response: %w", err)
	}

	var labels []string
	labels = response.Fields.Labels

	var components []string
	for _, component := range response.Fields.Components {
		components = append(components, component.Name)
	}

	var timeTracking *TimeTracking
	if response.Fields.TimeTracking.OriginalEstimate != "" ||
		response.Fields.TimeTracking.TimeSpent != "" {
		timeTracking = &TimeTracking{
			OriginalEstimate:  response.Fields.TimeTracking.OriginalEstimate,
			RemainingEstimate: response.Fields.TimeTracking.RemainingEstimate,
			TimeSpent:         response.Fields.TimeTracking.TimeSpent,
		}
	}

	return &EnhancedFields{
		Labels:       labels,
		Components:   components,
		Priority:     response.Fields.Priority.Name,
		IssueType:    response.Fields.IssueType.Name,
		TimeTracking: timeTracking,
	}, nil
}

// CheckIssuePermissions verifies what the user can access for an issue
func CheckIssuePermissions(client *jira.Client, baseURL, issueKey, apikey string) (IssuePermissions, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s?fields=id", baseURL, issueKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return IssuePermissions{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	httpClient := getHTTPClient(apikey)
	resp, err := httpClient.Do(req)
	if err != nil {
		return IssuePermissions{}, fmt.Errorf("failed to check permissions: %w", err)
	}
	defer resp.Body.Close()

	permissions := IssuePermissions{
		CanViewIssue:    resp.StatusCode == http.StatusOK,
		CanViewComments: resp.StatusCode == http.StatusOK,
		CanViewHistory:  resp.StatusCode == http.StatusOK,
	}

	// If we can't view the basic issue, we definitely can't view comments or history
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		permissions.CanViewComments = false
		permissions.CanViewHistory = false
	}

	return permissions, nil
}

// FetchIssueWithEnhancedContext retrieves an issue with all available enhanced context
// This is the main function consumers should use
func FetchIssueWithEnhancedContext(client *jira.Client, baseURL, issueKey, apikey string, verbose bool) (*Issue, error) {
	// Get basic issue first (use existing jiracrawler logic)
	jiraIssue, _, err := client.Issue.Get(issueKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue: %w", err)
	}

	if verbose {
		fmt.Printf("Fetching enhanced context for issue %s...\n", issueKey)
	}

	// Check permissions first
	permissions, err := CheckIssuePermissions(client, baseURL, issueKey, apikey)
	if err != nil {
		if verbose {
			fmt.Printf("Warning: failed to check permissions for %s: %v\n", issueKey, err)
		}
		// Continue anyway, we'll get errors on individual fetches if needed
	} else if !permissions.CanViewIssue {
		fmt.Printf("Warning: insufficient permissions to view issue %s - skipping enhanced context\n", issueKey)
		result := convertJiraIssue(*jiraIssue)
		return &result, nil
	}

	// Create result issue
	result := convertJiraIssue(*jiraIssue)

	// Fetch comments
	if permissions.CanViewComments {
		comments, err := FetchIssueComments(client, baseURL, issueKey, apikey)
		if err != nil {
			fmt.Printf("Warning: failed to fetch comments for %s: %v\n", issueKey, err)
		} else {
			result.Comments = comments
			if verbose {
				fmt.Printf("  ✓ Fetched %d comments for %s\n", len(comments), issueKey)
			}
		}
	} else if verbose {
		fmt.Printf("  ⊘ Skipping comments for %s (insufficient permissions)\n", issueKey)
	}

	// Fetch history
	if permissions.CanViewHistory {
		history, err := FetchIssueHistory(client, baseURL, issueKey, apikey)
		if err != nil {
			fmt.Printf("Warning: failed to fetch history for %s: %v\n", issueKey, err)
		} else {
			result.History = history
			if verbose {
				fmt.Printf("  ✓ Fetched %d history entries for %s\n", len(history), issueKey)
			}
		}
	} else if verbose {
		fmt.Printf("  ⊘ Skipping history for %s (insufficient permissions)\n", issueKey)
	}

	// Fetch additional fields
	enhancedFields, err := FetchEnhancedFields(client, baseURL, issueKey, apikey)
	if err != nil {
		fmt.Printf("Warning: failed to fetch enhanced fields for %s: %v\n", issueKey, err)
	} else {
		if verbose {
			fmt.Printf("  ✓ Fetched enhanced fields for %s\n", issueKey)
		}
		if enhancedFields.Labels != nil {
			result.Labels = enhancedFields.Labels
		}
		if enhancedFields.Components != nil {
			result.Components = enhancedFields.Components
		}
		if enhancedFields.TimeTracking != nil {
			result.TimeTracking = enhancedFields.TimeTracking
		}
	}

	return &result, nil
}

