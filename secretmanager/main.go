package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"secretmanager/secrets"
)

func main() {
	// Check Luganodes API key
	secretValue, err := secrets.GetSecret(
		"staging-191601",
		"staging-luganodes-eth-staking-api-key",
		"latest",
	)
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}
	fmt.Printf("Retrieved Luganodes API key (first 20 chars): %s...\n", truncate(secretValue, 20))
	fmt.Printf("Full secret length: %d bytes\n", len(secretValue))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
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
