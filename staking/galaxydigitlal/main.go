package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	devApiKey         = "58ff731f9a3d49d689188fa333cdb8b8"
	devPrivateKeyFile = "keys/private_key.pem"
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
	data, err := os.ReadFile(devPrivateKeyFile)
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
	if _, err := os.Stat(devPrivateKeyFile); os.IsNotExist(err) {
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
		"sub":      devApiKey,
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

type AuthRoundTripper struct {
	base http.RoundTripper
	opts AuthConfig
}

type AuthConfig struct {
	PrivateKey    *rsa.PrivateKey
	SigningMethod jwt.SigningMethod
	TokenDuration time.Duration
	ApiKey        string
}

type AuthClaims struct {
	URI      string `json:"uri"`
	Nonce    string `json:"nonce"`
	Sub      string `json:"sub"`
	BodyHash string `json:"bodyHash,omitempty"`
	jwt.RegisteredClaims
}

func NewAuthRoundTripper(opts AuthConfig, baseRoundTripper http.RoundTripper) (*AuthRoundTripper, error) {
	if opts.PrivateKey == nil {
		return nil, errors.New("RSA private key is not set")
	}
	if opts.SigningMethod == nil {
		return nil, errors.New("signing method is not set")
	}
	if opts.TokenDuration > 30*time.Second {
		return nil, errors.New("token duration cannot be longer than 30 seconds")
	}
	if opts.TokenDuration <= 0 {
		return nil, errors.New("token duration must be positive")
	}
	if opts.ApiKey == "" {
		return nil, errors.New("API key is not set")
	}
	if baseRoundTripper == nil {
		baseRoundTripper = http.DefaultTransport
	}
	return &AuthRoundTripper{
		base: baseRoundTripper,
		opts: opts,
	}, nil
}

// https://docs.staking.galaxy.com/authentication-and-authorization
func (b *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	claims := AuthClaims{
		URI:      req.URL.Path,
		Nonce:    uuid.NewString(),
		Sub:      b.opts.ApiKey,
		BodyHash: hashBody(req),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(b.opts.TokenDuration)),
		},
	}

	token := jwt.NewWithClaims(b.opts.SigningMethod, claims)
	bearer, err := token.SignedString(b.opts.PrivateKey)
	if err != nil {
		return nil, errors.New("failed to sign Auth token")
	}

	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("X-API-KEY", b.opts.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	return b.base.RoundTrip(req)
}

func hashBody(req *http.Request) string {
	if req.Body == nil {
		return ""
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println(req.Context(), "failed to read request body", "error", err)
		return ""
	}

	// Restore the body
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if len(bodyBytes) == 0 {
		return ""
	}

	// Hash the raw body bytes as sent
	hash := sha256.Sum256(bodyBytes)
	hashStr := hex.EncodeToString(hash[:])
	fmt.Println(req.Context(), "hashing request body", "bodyLength", len(bodyBytes), "hash", hashStr)
	return hashStr
}
