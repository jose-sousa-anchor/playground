package main

import (
	"context"
	"fmt"
	"log"

	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var secretManagerClient *secretmanager.Client

func main() {
	var err error
	secretManagerClient, err = secretmanager.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create secret manager client: %v", err)
	}
	defer secretManagerClient.Close()

	fmt.Println("Secret Manager client created successfully.")

	// To get project run `gcloud config get-value project`
	// To get secret id run `gcloud secrets list`
	secretValue, err := GetSecret("development-204920", "default-galaxy_digital_rsa_private_key", "latest")
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}
	fmt.Printf("Retrieved secret: %s\n", secretValue)

	rsaKey, err := parseRSAPrivateKeyFromString(secretValue)
	if err != nil {
		log.Fatalf("Failed to parse RSA private key: %v", err)
	}
	fmt.Printf("Parsed RSA Private Key: %v\n", rsaKey)
}

// GetSecret retrieves the secret value from Google Cloud Secret Manager.
func GetSecret(projectID, secretID, version string) (string, error) {
	ctx := context.Background()

	// Build the request
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, secretID, version),
	}

	// Call the API
	result, err := secretManagerClient.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	// Return the secret payload
	return string(result.Payload.Data), nil
}

// ParseRSAPrivateKeyFromString parses a PEM-encoded RSA private key (PKCS#1 or PKCS#8).
func parseRSAPrivateKeyFromString(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("no PEM data found")
	}

	// Try PKCS#1 first
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// Fallback: PKCS#8
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}

	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed key is not RSA")
	}
	return rsaKey, nil
}
