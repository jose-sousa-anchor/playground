package secrets

import (
	"context"
	"fmt"

	"cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// GetSecret retrieves the secret value from Google Cloud Secret Manager.
func GetSecret(projectID, secretID, version string) (string, error) {
	ctx := context.Background()

	secretManagerClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secret manager client: %w", err)
	}
	defer secretManagerClient.Close()

	// Build the full secret name path
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, secretID, version)

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API
	result, err := secretManagerClient.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	// Return the secret payload
	return string(result.Payload.Data), nil
}
