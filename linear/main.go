package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	apiKey, err := LoadAPIKeyFromEnv()
	if err != nil {
		log.Fatalf("Failed to load API key: %v", err)
	}

	client := NewLinearClient(apiKey)

	// First get the viewer to find Jose's user ID
	viewer, err := client.GetViewer()
	if err != nil {
		log.Fatalf("Failed to get viewer: %v", err)
	}
	
	viewerData := viewer["viewer"].(map[string]interface{})
	userID := viewerData["id"].(string)
	userName := viewerData["name"].(string)
	userEmail := viewerData["email"].(string)
	
	fmt.Printf("User: %s <%s> (ID: %s)\n\n", userName, userEmail, userID)

	// Query for issues assigned to Jose updated in 2025
	query := `
		query UserIssues($userId: ID!, $after: String) {
			issues(
				first: 100
				after: $after
				filter: {
					assignee: { id: { eq: $userId } }
					updatedAt: { gte: "2025-01-01T00:00:00Z" }
				}
				orderBy: updatedAt
			) {
				pageInfo {
					hasNextPage
					endCursor
				}
				nodes {
					id
					title
					identifier
					description
					priority
					estimate
					state {
						name
						type
					}
					assignee {
						name
						email
					}
					creator {
						name
					}
					parent {
						id
						title
						identifier
					}
					labels {
						nodes {
							name
						}
					}
					createdAt
					updatedAt
					completedAt
					url
					project {
						name
					}
					team {
						name
						key
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"userId": userID,
		"after":  nil,
	}

	resp, err := client.Query(query, variables)
	if err != nil {
		log.Fatalf("Failed to query issues: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	issuesData := result["issues"].(map[string]interface{})
	nodes := issuesData["nodes"].([]interface{})
	
	fmt.Printf("Found %d issues assigned to %s in 2025\n\n", len(nodes), userName)
	
	// Print summary by state
	stateCount := make(map[string]int)
	projectCount := make(map[string]int)
	
	for _, node := range nodes {
		issue := node.(map[string]interface{})
		state := issue["state"].(map[string]interface{})
		stateName := state["name"].(string)
		stateCount[stateName]++
		
		if proj, ok := issue["project"].(map[string]interface{}); ok && proj != nil {
			projName := proj["name"].(string)
			projectCount[projName]++
		}
	}
	
	fmt.Println("=== Issues by State ===")
	for state, count := range stateCount {
		fmt.Printf("  %s: %d\n", state, count)
	}
	
	fmt.Println("\n=== Issues by Project ===")
	for proj, count := range projectCount {
		fmt.Printf("  %s: %d\n", proj, count)
	}
	
	fmt.Println("\n=== All Issues ===")
	jsonBytes, _ := json.MarshalIndent(nodes, "", "  ")
	fmt.Println(string(jsonBytes))
}
