package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	// The project ID from the URL: https://linear.app/anchorlabs/project/asset-registry-1c62eb6f727d/overview
	projectID := "1c62eb6f727d"

	// Query for comprehensive project information
	query := `
		query Project($projectId: String!) {
			project(id: $projectId) {
				id
				name
				description
				state
				startDate
				targetDate
				completedAt
				url
				progress
				lead {
					id
					name
					email
				}
				members {
					nodes {
						id
						name
						email
					}
				}
				teams {
					nodes {
						id
						name
						key
					}
				}
				issues {
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
							id
							name
							email
						}
						creator {
							name
							email
						}
						parent {
							id
							title
							identifier
						}
						children {
							nodes {
								id
								title
								identifier
							}
						}
						labels {
							nodes {
								id
								name
								color
							}
						}
						createdAt
						updatedAt
						completedAt
						url
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"projectId": projectID,
	}

	resp, err := client.Query(query, variables)
	if err != nil {
		log.Fatalf("Failed to query project: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Pretty print the entire response
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Save to file
	outputFile := "project_asset_registry.json"
	if err := os.WriteFile(outputFile, jsonBytes, 0644); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Printf("Project information saved to %s\n", outputFile)
	fmt.Println(string(jsonBytes))
}
