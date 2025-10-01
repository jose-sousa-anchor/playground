package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	stagingPrivateKeyFile = "keys/staging/private_key.pem"
	stagingPublicKeyFile  = "keys/staging/public_key.pem"
	prodPrivateKeyFile    = "keys/prod/private_key.pem"
	prodPublicKeyFile     = "keys/prod/public_key.pem"
)

type JWTSignablePayload struct {
	Url string `json:"url"`
	// Nonce is a unique request ID
	Nonce string `json:"nonce"`
	// Iat is the request generated timestamp in UTC
	Iat int64 `json:"iat"`
	// Exp is the expiration timestamp of request in UTC (must be less than iat+ 30sec)
	Exp int64 `json:"exp"`
	// Sub is the API key shared by the BlueShift
	Sub string `json:"sub"`
	// BodyHash is the SHA256 hash of the request body if applicable.
	BodyHash string `json:"bodyHash"`
}

// loadPrivateKey loads an RSA private key from PEM file
func loadPrivateKey() (any, error) {
	data, err := os.ReadFile(prodPrivateKeyFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS8PrivateKey(block.Bytes)
}

func main() {
	if _, err := os.Stat(prodPrivateKeyFile); os.IsNotExist(err) {
		panic(err)
	}

	// Step 2: load private key from file
	privateKey, err := loadPrivateKey()
	if err != nil {
		panic(err)
	}

	// Step 3: build JWT claims
	claims := jwt.MapClaims{
		"uri":      "/v1/resource",
		"nonce":    uuid.New().String(),
		"iat":      time.Now().UTC().Unix(),
		"exp":      time.Now().UTC().Add(30 * time.Second).Unix(),
		"sub":      "YOUR_API_KEY_HERE",
		"bodyHash": "BASE64_SHA256_OF_BODY",
	}

	// Step 4: sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("JWT token:")
	fmt.Println(tokenString)
}
