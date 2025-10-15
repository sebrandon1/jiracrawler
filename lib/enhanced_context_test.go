package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
)

// TestFetchIssueComments tests fetching comments for an issue
func TestFetchIssueComments(t *testing.T) {
	// Create mock Jira server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/rest/api/2/issue/TEST-123/comment", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		// Mock response
		response := map[string]interface{}{
			"comments": []map[string]interface{}{
				{
					"id":   "12345",
					"body": "Test comment",
					"author": map[string]string{
						"displayName": "Test User",
					},
					"created": "2025-01-01T12:00:00.000-0700",
					"updated": "2025-01-01T13:00:00.000-0700",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	tokenAuth := jira.BearerAuthTransport{
		Token: "test-token",
	}
	client, err := jira.NewClient(tokenAuth.Client(), server.URL)
	assert.NoError(t, err)

	// Test function
	comments, err := FetchIssueComments(client, server.URL, "TEST-123", "test-token")
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, "12345", comments[0].ID)
	assert.Equal(t, "Test comment", comments[0].Body)
	assert.Equal(t, "Test User", comments[0].Author)
}

// TestFetchIssueCommentsError tests error handling for comments
func TestFetchIssueCommentsError(t *testing.T) {
	// Create mock Jira server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Issue not found"))
	}))
	defer server.Close()

	// Create client
	tokenAuth := jira.BearerAuthTransport{
		Token: "test-token",
	}
	client, err := jira.NewClient(tokenAuth.Client(), server.URL)
	assert.NoError(t, err)

	// Test function
	comments, err := FetchIssueComments(client, server.URL, "TEST-999", "test-token")
	assert.Error(t, err)
	assert.Nil(t, comments)
	assert.Contains(t, err.Error(), "404")
}

// TestFetchIssueHistory tests fetching history for an issue
func TestFetchIssueHistory(t *testing.T) {
	// Create mock Jira server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Contains(t, r.URL.Path, "/rest/api/2/issue/TEST-123")
		assert.Contains(t, r.URL.RawQuery, "expand=changelog")

		// Mock response
		response := map[string]interface{}{
			"changelog": map[string]interface{}{
				"histories": []map[string]interface{}{
					{
						"id": "67890",
						"author": map[string]string{
							"displayName": "Test User",
						},
						"created": "2025-01-01T12:00:00.000-0700",
						"items": []map[string]interface{}{
							{
								"field":      "status",
								"fieldtype":  "jira",
								"fromString": "Open",
								"toString":   "In Progress",
							},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	tokenAuth := jira.BearerAuthTransport{
		Token: "test-token",
	}
	client, err := jira.NewClient(tokenAuth.Client(), server.URL)
	assert.NoError(t, err)

	// Test function
	history, err := FetchIssueHistory(client, server.URL, "TEST-123", "test-token")
	assert.NoError(t, err)
	assert.Len(t, history, 1)
	assert.Equal(t, "67890", history[0].ID)
	assert.Equal(t, "Test User", history[0].Author)
	assert.Len(t, history[0].Items, 1)
	assert.Equal(t, "status", history[0].Items[0].Field)
	assert.Equal(t, "Open", history[0].Items[0].FromString)
	assert.Equal(t, "In Progress", history[0].Items[0].ToString)
}

// TestFetchEnhancedFields tests fetching enhanced fields for an issue
func TestFetchEnhancedFields(t *testing.T) {
	// Create mock Jira server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Contains(t, r.URL.Path, "/rest/api/2/issue/TEST-123")
		assert.Contains(t, r.URL.RawQuery, "fields=labels,components,priority,issuetype,timetracking")

		// Mock response
		response := map[string]interface{}{
			"fields": map[string]interface{}{
				"labels": []string{"bug", "urgent"},
				"components": []map[string]string{
					{"name": "Frontend"},
					{"name": "Backend"},
				},
				"priority": map[string]string{
					"name": "High",
				},
				"issuetype": map[string]string{
					"name": "Bug",
				},
				"timetracking": map[string]string{
					"originalEstimate":  "8h",
					"remainingEstimate": "4h",
					"timeSpent":         "4h",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	tokenAuth := jira.BearerAuthTransport{
		Token: "test-token",
	}
	client, err := jira.NewClient(tokenAuth.Client(), server.URL)
	assert.NoError(t, err)

	// Test function
	fields, err := FetchEnhancedFields(client, server.URL, "TEST-123", "test-token")
	assert.NoError(t, err)
	assert.NotNil(t, fields)
	assert.Equal(t, []string{"bug", "urgent"}, fields.Labels)
	assert.Equal(t, []string{"Frontend", "Backend"}, fields.Components)
	assert.Equal(t, "High", fields.Priority)
	assert.Equal(t, "Bug", fields.IssueType)
	assert.NotNil(t, fields.TimeTracking)
	assert.Equal(t, "8h", fields.TimeTracking.OriginalEstimate)
	assert.Equal(t, "4h", fields.TimeTracking.RemainingEstimate)
	assert.Equal(t, "4h", fields.TimeTracking.TimeSpent)
}

// TestCheckIssuePermissions tests checking permissions for an issue
func TestCheckIssuePermissions(t *testing.T) {
	// Test successful permission check
	t.Run("CanView", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"123"}`))
		}))
		defer server.Close()

		tokenAuth := jira.BearerAuthTransport{
			Token: "test-token",
		}
		client, err := jira.NewClient(tokenAuth.Client(), server.URL)
		assert.NoError(t, err)

		perms, err := CheckIssuePermissions(client, server.URL, "TEST-123", "test-token")
		assert.NoError(t, err)
		assert.True(t, perms.CanViewIssue)
		assert.True(t, perms.CanViewComments)
		assert.True(t, perms.CanViewHistory)
	})

	// Test forbidden permission check
	t.Run("Forbidden", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer server.Close()

		tokenAuth := jira.BearerAuthTransport{
			Token: "test-token",
		}
		client, err := jira.NewClient(tokenAuth.Client(), server.URL)
		assert.NoError(t, err)

		perms, err := CheckIssuePermissions(client, server.URL, "TEST-123", "test-token")
		assert.NoError(t, err)
		assert.False(t, perms.CanViewIssue)
		assert.False(t, perms.CanViewComments)
		assert.False(t, perms.CanViewHistory)
	})
}

// TestFetchIssueWithEnhancedContext tests the main enhanced context function
func TestFetchIssueWithEnhancedContext(t *testing.T) {
	// Create mock Jira server
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")

		// Handle different endpoints
		if r.URL.Path == "/rest/api/2/issue/TEST-123" {
			if r.URL.RawQuery == "fields=id" {
				// Permission check
				json.NewEncoder(w).Encode(map[string]string{"id": "123"})
			} else if r.URL.RawQuery == "expand=changelog" {
				// History
				json.NewEncoder(w).Encode(map[string]interface{}{
					"changelog": map[string]interface{}{
						"histories": []interface{}{},
					},
				})
			} else if r.URL.RawQuery != "" {
				// Enhanced fields
				json.NewEncoder(w).Encode(map[string]interface{}{
					"fields": map[string]interface{}{
						"labels":     []string{"test"},
						"components": []interface{}{},
						"priority":   map[string]string{"name": "Medium"},
						"issuetype":  map[string]string{"name": "Task"},
						"timetracking": map[string]string{
							"originalEstimate": "1h",
						},
					},
				})
			} else {
				// Basic issue fetch (for client.Issue.Get)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"key": "TEST-123",
					"fields": map[string]interface{}{
						"summary":     "Test Issue",
						"description": "Test Description",
						"status": map[string]interface{}{
							"id":   "1",
							"name": "Open",
						},
						"priority": map[string]interface{}{
							"id":   "3",
							"name": "Medium",
						},
						"issuetype": map[string]interface{}{
							"id":   "1",
							"name": "Task",
						},
						"project": map[string]interface{}{
							"id":   "10000",
							"key":  "TEST",
							"name": "Test Project",
						},
						"created": "2025-01-01T12:00:00.000-0700",
						"updated": "2025-01-01T13:00:00.000-0700",
					},
				})
			}
		} else if r.URL.Path == "/rest/api/2/issue/TEST-123/comment" {
			// Comments
			json.NewEncoder(w).Encode(map[string]interface{}{
				"comments": []map[string]interface{}{
					{
						"id":   "1",
						"body": "Test comment",
						"author": map[string]string{
							"displayName": "Test User",
						},
						"created": "2025-01-01T12:00:00.000-0700",
						"updated": "2025-01-01T12:00:00.000-0700",
					},
				},
			})
		}
	}))
	defer server.Close()

	// Create client
	tokenAuth := jira.BearerAuthTransport{
		Token: "test-token",
	}
	client, err := jira.NewClient(tokenAuth.Client(), server.URL)
	assert.NoError(t, err)

	// Test function (non-verbose)
	issue, err := FetchIssueWithEnhancedContext(client, server.URL, "TEST-123", "test-token", false)
	assert.NoError(t, err)
	assert.NotNil(t, issue)
	assert.Equal(t, "TEST-123", issue.Key)
	assert.Equal(t, "Test Issue", issue.Summary)
	assert.Len(t, issue.Comments, 1)
	assert.NotNil(t, issue.Labels)
}

// TestGetHTTPClient tests the HTTP client helper
func TestGetHTTPClient(t *testing.T) {
	client := getHTTPClient("test-token")
	assert.NotNil(t, client)
	assert.NotNil(t, client.Transport)
}

// TestFetchUserIssuesInDateRangeWithContext tests the batch function
func TestFetchUserIssuesInDateRangeWithContext(t *testing.T) {
	// This is a more complex integration test that would require mocking
	// the entire Jira search API. For now, we test that it handles
	// the enhancedContext flag correctly.
	
	t.Run("WithoutEnhancedContext", func(t *testing.T) {
		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/rest/api/2/search" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"issues": []interface{}{},
					"total":  0,
				})
			}
		}))
		defer server.Close()

		result := FetchUserIssuesInDateRangeWithContext(
			server.URL,
			"test@example.com",
			"test-token",
			"test@example.com",
			"2025-01-01",
			"2025-01-31",
			false, // no enhanced context
			false,
		)
		assert.NotNil(t, result)
	})
}

// TestCommentStruct tests the Comment struct
func TestCommentStruct(t *testing.T) {
	now := time.Now()
	comment := Comment{
		ID:      "123",
		Author:  "Test User",
		Body:    "Test comment body",
		Created: now,
		Updated: now,
	}

	assert.Equal(t, "123", comment.ID)
	assert.Equal(t, "Test User", comment.Author)
	assert.Equal(t, "Test comment body", comment.Body)
	assert.Equal(t, now, comment.Created)
	assert.Equal(t, now, comment.Updated)
}

// TestHistoryItemStruct tests the HistoryItem struct
func TestHistoryItemStruct(t *testing.T) {
	now := time.Now()
	historyItem := HistoryItem{
		ID:      "456",
		Author:  "Test User",
		Created: now,
		Items: []HistoryChange{
			{
				Field:      "status",
				FieldType:  "jira",
				FromString: "Open",
				ToString:   "Closed",
			},
		},
	}

	assert.Equal(t, "456", historyItem.ID)
	assert.Equal(t, "Test User", historyItem.Author)
	assert.Equal(t, now, historyItem.Created)
	assert.Len(t, historyItem.Items, 1)
	assert.Equal(t, "status", historyItem.Items[0].Field)
}

// TestTimeTrackingStruct tests the TimeTracking struct
func TestTimeTrackingStruct(t *testing.T) {
	timeTracking := TimeTracking{
		OriginalEstimate:  "8h",
		RemainingEstimate: "4h",
		TimeSpent:         "4h",
	}

	assert.Equal(t, "8h", timeTracking.OriginalEstimate)
	assert.Equal(t, "4h", timeTracking.RemainingEstimate)
	assert.Equal(t, "4h", timeTracking.TimeSpent)
}

// TestIssuePermissionsStruct tests the IssuePermissions struct
func TestIssuePermissionsStruct(t *testing.T) {
	perms := IssuePermissions{
		CanViewComments: true,
		CanViewHistory:  true,
		CanViewIssue:    true,
	}

	assert.True(t, perms.CanViewComments)
	assert.True(t, perms.CanViewHistory)
	assert.True(t, perms.CanViewIssue)
}

// TestEnhancedFieldsStruct tests the EnhancedFields struct
func TestEnhancedFieldsStruct(t *testing.T) {
	fields := EnhancedFields{
		Labels:     []string{"bug", "urgent"},
		Components: []string{"Frontend"},
		Priority:   "High",
		IssueType:  "Bug",
		TimeTracking: &TimeTracking{
			OriginalEstimate: "8h",
		},
	}

	assert.Equal(t, []string{"bug", "urgent"}, fields.Labels)
	assert.Equal(t, []string{"Frontend"}, fields.Components)
	assert.Equal(t, "High", fields.Priority)
	assert.Equal(t, "Bug", fields.IssueType)
	assert.NotNil(t, fields.TimeTracking)
	assert.Equal(t, "8h", fields.TimeTracking.OriginalEstimate)
}

