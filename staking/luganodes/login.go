package luganodes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AuthClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAuthClient returns an authentication helper with a base URL
func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SignupRequest is the body for the /api/signup call
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	OrgName  string `json:"orgName"`
}

// SignupResponse is the expected response from /api/signup
type SignupResponse struct {
	Result struct {
		User struct {
			APIKey string `json:"apiKey"`
		} `json:"user"`
	} `json:"result"`
}

// Signup registers a new organization user with Luganodes
func (a *AuthClient) Signup(
	ctx context.Context,
	email, password, orgName string,
) (*SignupResponse, error) {
	url := fmt.Sprintf("%s/api/signup", a.BaseURL)

	body := fmt.Sprintf(
		`{"email":"%s","password":"%s","orgName":"%s"}`,
		email, password, orgName,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SignupResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LoginRequest is the request body for the /api/login call
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the expected response from /api/login
type LoginResponse struct {
	Result struct {
		User struct {
			APIKey string `json:"apiKey"`
		} `json:"user"`
	} `json:"result"`
}

// Login logs into an existing Luganodes user and returns an API key
func (a *AuthClient) Login(
	ctx context.Context,
	email, password string,
) (*LoginResponse, error) {
	url := fmt.Sprintf("%s/api/login", a.BaseURL)

	body := fmt.Sprintf(
		`{"email":"%s","password":"%s"}`,
		email, password,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
