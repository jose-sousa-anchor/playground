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

	// Find the "Asset Registry Tagging" project
	projectName := "Asset Registry Tagging"
	fmt.Printf("=== Searching for Project: %s ===\n", projectName)
	projectData, err := client.GetProjectByName(projectName)
	if err != nil {
		log.Fatalf("Failed to find project: %v", err)
	}

	// Extract project ID from the response
	projects, ok := projectData["projects"].(map[string]interface{})
	if !ok {
		log.Fatalf("Unexpected project data format")
	}

	nodes, ok := projects["nodes"].([]interface{})
	if !ok || len(nodes) == 0 {
		log.Fatalf("No project found with name: %s", projectName)
	}

	project := nodes[0].(map[string]interface{})
	projectID := project["id"].(string)

	fmt.Printf("Found project: %s (ID: %s)\n\n", project["name"], projectID)

	// Get all issues for the project
	fmt.Println("=== Getting Issues for Asset Registry Tagging ===")
	issues, err := client.GetProjectIssues(projectID)
	if err != nil {
		log.Fatalf("Failed to get project issues: %v", err)
	}

	printJSON(issues)
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
