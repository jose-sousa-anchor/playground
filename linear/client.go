package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const LinearAPIEndpoint = "https://api.linear.app/graphql"

// LinearClient represents a client for interacting with the Linear GraphQL API
type LinearClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewLinearClient creates a new Linear API client
func NewLinearClient(apiKey string) *LinearClient {
	return &LinearClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// GraphQLRequest represents a GraphQL request payload
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   json.RawMessage          `json:"data,omitempty"`
	Errors []GraphQLError           `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message string `json:"message"`
	Path    []string `json:"path,omitempty"`
}

// Query executes a GraphQL query and returns the response
func (c *LinearClient) Query(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", LinearAPIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return &gqlResp, fmt.Errorf("GraphQL errors: %+v", gqlResp.Errors)
	}

	return &gqlResp, nil
}

// GetViewer fetches the current authenticated user
func (c *LinearClient) GetViewer() (map[string]interface{}, error) {
	query := `
		query Me {
			viewer {
				id
				name
				email
			}
		}
	`

	resp, err := c.Query(query, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal viewer data: %w", err)
	}

	return result, nil
}

// GetTeams fetches all teams
func (c *LinearClient) GetTeams() (map[string]interface{}, error) {
	query := `
		query Teams {
			teams {
				nodes {
					id
					name
					key
				}
			}
		}
	`

	resp, err := c.Query(query, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal teams data: %w", err)
	}

	return result, nil
}

// GetTeamIssues fetches issues for a specific team
func (c *LinearClient) GetTeamIssues(teamID string) (map[string]interface{}, error) {
	query := `
		query TeamIssues($teamId: String!) {
			team(id: $teamId) {
				id
				name
				issues {
					nodes {
						id
						title
						identifier
						state {
							name
						}
						assignee {
							id
							name
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"teamId": teamID,
	}

	resp, err := c.Query(query, variables)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team issues data: %w", err)
	}

	return result, nil
}

// CreateIssue creates a new issue
func (c *LinearClient) CreateIssue(teamID, title, description string) (map[string]interface{}, error) {
	query := `
		mutation IssueCreate($teamId: String!, $title: String!, $description: String) {
			issueCreate(
				input: {
					teamId: $teamId
					title: $title
					description: $description
				}
			) {
				success
				issue {
					id
					title
					identifier
					url
				}
			}
		}
	`

	variables := map[string]interface{}{
		"teamId":      teamID,
		"title":       title,
		"description": description,
	}

	resp, err := c.Query(query, variables)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal issue create data: %w", err)
	}

	return result, nil
}

// GetProjectByName searches for a project by name and returns it
func (c *LinearClient) GetProjectByName(projectName string) (map[string]interface{}, error) {
	query := `
		query Projects($filter: ProjectFilter!) {
			projects(filter: $filter) {
				nodes {
					id
					name
					description
					state
					url
				}
			}
		}
	`

	variables := map[string]interface{}{
		"filter": map[string]interface{}{
			"name": map[string]interface{}{
				"eq": projectName,
			},
		},
	}

	resp, err := c.Query(query, variables)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project data: %w", err)
	}

	return result, nil
}

// GetProjectIssues fetches all issues for a specific project
func (c *LinearClient) GetProjectIssues(projectID string) (map[string]interface{}, error) {
	query := `
		query ProjectIssues($projectId: String!) {
			project(id: $projectId) {
				id
				name
				issues {
					nodes {
						id
						title
						identifier
						description
						priority
						state {
							name
							type
						}
						assignee {
							id
							name
							email
						}
						createdAt
						updatedAt
						url
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"projectId": projectID,
	}

	resp, err := c.Query(query, variables)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project issues data: %w", err)
	}

	return result, nil
}

// LoadAPIKeyFromEnv loads the Linear API key from environment variable
func LoadAPIKeyFromEnv() (string, error) {
	apiKey := os.Getenv("LINEAR_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("LINEAR_API_KEY environment variable is not set")
	}
	return apiKey, nil
}
