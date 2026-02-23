package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load API key from environment
	apiKey, err := LoadAPIKeyFromEnv()
	if err != nil {
		log.Fatalf("Failed to load API key: %v", err)
	}

	// Create Linear client
	client := NewLinearClient(apiKey)

	// Fetch Asset Registry Tagging project
	projectName := "Asset Registry Tagging"
	fmt.Printf("=== Fetching Project: %s ===\n", projectName)
	projectData, err := client.GetProjectByName(projectName)
	if err != nil {
		log.Fatalf("Failed to get project: %v", err)
	}

	// Extract project ID
	projects, ok := projectData["projects"].(map[string]interface{})
	if !ok {
		log.Fatalf("Invalid project data structure")
	}

	nodes, ok := projects["nodes"].([]interface{})
	if !ok || len(nodes) == 0 {
		log.Fatalf("No project found with name: %s", projectName)
	}

	project := nodes[0].(map[string]interface{})
	projectID := project["id"].(string)
	fmt.Printf("\nProject found: %s (ID: %s)\n\n", project["name"], projectID)

	// Fetch all issues for this project
	fmt.Printf("=== Fetching Issues for Project: %s ===\n", projectName)
	issuesData, err := client.GetProjectIssues(projectID)
	if err != nil {
		log.Fatalf("Failed to get project issues: %v", err)
	}

	printJSON(issuesData)
}

// printJSON pretty prints JSON data
func printJSON(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
