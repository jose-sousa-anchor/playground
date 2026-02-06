# Linear API Integration in Go

A Go client for interacting with the Linear GraphQL API.

## Setup

1. Make sure you have your Linear API key in the `.env` file:
   ```
   LINEAR_API_KEY=*******
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

Run the example:
```bash
go run .
```

## Features

The client includes the following methods:

- `GetViewer()` - Get the current authenticated user
- `GetTeams()` - Fetch all teams
- `GetTeamIssues(teamID)` - Get issues for a specific team
- `CreateIssue(teamID, title, description)` - Create a new issue
- `Query(query, variables)` - Execute custom GraphQL queries

## Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    godotenv.Load()

    // Get API key from environment
    apiKey, _ := LoadAPIKeyFromEnv()

    // Create client
    client := NewLinearClient(apiKey)

    // Get current user
    viewer, err := client.GetViewer()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User: %+v\n", viewer)

    // Get all teams
    teams, err := client.GetTeams()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Teams: %+v\n", teams)
}
```

## Custom Queries

You can execute any custom GraphQL query:

```go
query := `
    query {
        issues(first: 10) {
            nodes {
                id
                title
                state {
                    name
                }
            }
        }
    }
`

resp, err := client.Query(query, nil)
// Handle response...
```
