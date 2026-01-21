package p2pclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://api-test.p2p.org"
)

/*
   ---------- HTTP INFRA ----------
*/

// NewHTTPClient returns a configured HTTP client
func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
	}
}

// DoRequest executes an HTTP request and returns raw body + status
func DoRequest(client *http.Client, req *http.Request) ([]byte, int, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

/*
   ---------- CREATE NODE REQUEST (POST) ----------
*/

// CreateNodeRequestPayload mirrors the API payload
type CreateNodeRequestPayload struct {
	ID                        string            `json:"id"`
	Type                      string            `json:"type"`
	ValidatorsCount           int               `json:"validatorsCount"`
	AmountPerValidator        string            `json:"amountPerValidator"`
	WithdrawalCredentialsType string            `json:"withdrawalCredentialsType"`
	WithdrawalAddress         string            `json:"withdrawalAddress"`
	EigenPodOwnerAddress      string            `json:"eigenPodOwnerAddress"`
	ControllerAddress         string            `json:"controllerAddress"`
	FeeRecipientAddress       string            `json:"feeRecipientAddress"`
	NodesOptions              NodesOptionsInput `json:"nodesOptions"`
}

type NodesOptionsInput struct {
	Location  string `json:"location"`
	RelaysSet string `json:"relaysSet"`
}

// NewCreateNodeRequest builds the POST request
func NewCreateNodeRequest(
	ctx context.Context,
	payload CreateNodeRequestPayload,
	bearerToken string,
) (*http.Request, error) {
	url := baseURL + "/api/v1/eth/staking/direct/nodes-request/create"

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	return req, nil
}

/*
   ---------- GET NODE REQUEST STATUS (GET) ----------
*/

// NewGetNodeRequestStatusRequest builds the GET request
func NewGetNodeRequestStatusRequest(
	ctx context.Context,
	nodeRequestID string,
	bearerToken string,
) (*http.Request, error) {
	url := fmt.Sprintf(
		"%s/api/v1/eth/staking/direct/nodes-request/status/%s",
		baseURL,
		nodeRequestID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	return req, nil
}

// doRequest executes and reads the HTTP response
func doRequest(client *http.Client, req *http.Request) ([]byte, int, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

const vemBaseURL = baseURL + "/api/v1/eth/staking/direct/vem"

type VemCreatePayload struct {
	ID                  string `json:"id"`
	Type                string `json:"type"`
	VemRequest          string `json:"vemRequest"`
	VemRequestSignature string `json:"vemRequestSignature"`
	VemRequestSignedBy  string `json:"vemRequestSignedBy"`
	VemRequestTxId      string `json:"vemRequestTxId"`
	VemRequestProof     string `json:"vemRequestProof"`
}

func NewVemCreateRequest(
	ctx context.Context,
	payload VemCreatePayload,
	token string,
) (*http.Request, error) {
	b, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		vemBaseURL+"/create",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func NewVemStatusRequest(
	ctx context.Context,
	id, token string,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/status/%s", vemBaseURL, id),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	return req, nil
}
