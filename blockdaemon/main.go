package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"secretmanager/secrets"
)

const (
	defaultTimeout = 30 * time.Second
	// This is customized. It's necessary to keep http connections open to reduce latency from re-connecting over and over again
	defaultMaxIdleConns = 1000
	// This is customized. It's necessary for making many short lived requests in parallel
	// If this number is exceeded, new requests will fail with "dial tcp 127.0.0.1:8332: connect: cannot assign requested address"
	defaultMaxConnsPerHost = 1000
)

// HTTPClientOptions allows for popular http client configuration tweaks
// zero values are anchorage-set defaults or go library defaults
type HTTPClientOptions struct {
	MaxConnsPerHost int
	Timeout         time.Duration
	MaxIdleConns    int
	RoundTripper    http.RoundTripper
	InsecureTLS     bool
}

// NetworkName represents the name of a network
type NetworkName string

const (
	// Mainnet is the Ethereum mainnet
	Mainnet NetworkName = "mainnet"
	// Holesky is the Holesky network
	Holesky NetworkName = "holesky"
	// Hoodi is the Hoodi network
	Hoodi NetworkName = "hoodi"
)

// PostStakeIntentRequest represents the request to PostStakeIntent
type PostStakeIntentRequest struct {
	NetworkName NetworkName `json:"-"`
	// PlanID is the ID of the plan defining the staking parameters (region, etc.)
	// This feature enables you to stake to validators from the specified plan(s). When no plan id is specified, validators across all plans that match the API route will be available for staking. If it is a shared plan, you can to stake to validators on a shared node from the specified plan.
	PlanID string `json:"plan_id"`
	// Stakes is a list of stake parameters to generate intents for.
	Stakes []StakeRequest `json:"stakes"`
}

// StakeRequest represents a single stake tx in the PostStakeIntentRequest
type StakeRequest struct {
	// Amount of ETH to be staked (denominated in Gwei).
	Amount string `json:"amount"`
	// WithdrawalAddress is an hex-encoded ethereum account or smart contract address. (e.g. 0x93247f2209abcacf57b75a51dafae777f9dd38bc)
	WithdrawalAddress string `json:"withdrawal_address"`
	// Quantity is the number of validators to create that will share the same withdrawal credentials.
	Quantity int `json:"quantity"`
	// FeeRecipient is an ethereum address to receive transaction fees from published blocks. 20-bytes, hex encoded with 0x prefix, case insensitive.
	// Defaults to the withdrawal address if not provided
	FeeRecipient string `json:"fee_recipient"`
}

// EthereumStakeIntent represents the ethereum's stake intents
type EthereumStakeIntent struct {
	// ContractAddress is an hex-encoded ethereum account or smart contract address.
	ContractAddress string `json:"contract_address"`
	// EstimatedGas is the estimated gas for the transaction
	EstimatedGas int `json:"estimated_gas"`
	// ExpirationTime is the transaction expiration time
	ExpirationTime int `json:"expiration_time"`
	// Stakes being made.
	Stakes []Stake
	// UnsignedTransaction is the unsigned transaction to be signed by the user.
	UnsignedTransaction string `json:"unsigned_transaction"`
}

// Stake represents a single stake tx in a stake intent
type Stake struct {
	// Amount of ETH (denominated in Gwei).
	Amount string `json:"amount"`
	// FeeRecipient is an ethereum address to receive transaction fees from published blocks. 20-bytes, hex encoded with 0x prefix, case insensitive.
	FeeRecipient string `json:"fee_recipient"`
	// StakeID is the unique identifier for the stake.
	StakeID string `json:"stake_id"`
	// ValidatorPublicKey is a BLS public Key.
	ValidatorPublicKey string `json:"validator_public_key"`
	// WithdrawalCredentials is an hexadecimal encoded withdrawal credentials which can be either a BLS public key or an ethereum account address.
	WithdrawalCredentials string `json:"withdrawal_credentials"`
}

// StakeIntent represents the response body for a stake intent
type StakeIntent struct {
	// CustomerID is the ID of the customer who made the request
	CustomerID string `json:"customer_id"`
	// Ethereum is the response for Ethereum staking
	Ethereum EthereumStakeIntent `json:"ethereum"`
	// Network is the network on which the staking is being done (e.g. mainnet, holesky)
	Network string `json:"network"`
	// Protocol is the protocol on which the staking is being done (e.g. ethereum)
	Protocol string `json:"protocol"`
	// StakeIntentID is an unique idenifier for a group of stakes.
	StakeIntentID string `json:"stake_intent_id"`
}

// PostStakeIntentResponse is the response to PostStakeIntent
type PostStakeIntentResponse StakeIntent

// setDefaults creates a copy with overridden zero values to anchorage-set defaults
func (o HTTPClientOptions) withDefaults() HTTPClientOptions {
	optionsWithDefault := o

	if o.Timeout == 0 {
		optionsWithDefault.Timeout = defaultTimeout
	}
	if o.MaxConnsPerHost == 0 {
		optionsWithDefault.MaxConnsPerHost = defaultMaxConnsPerHost
	}
	if o.MaxIdleConns == 0 {
		optionsWithDefault.MaxIdleConns = defaultMaxIdleConns
	}

	return optionsWithDefault
}

type ValidatorType string

const (
	ValidatorType0x01 ValidatorType = "0x01"
	ValidatorType0x02 ValidatorType = "0x02"
)

func CreateStakeIntent(
	ctx context.Context,
	httpClient http.Client,
	apiKey string,
	req PostStakeIntentRequest,
	validatorType ValidatorType) (*PostStakeIntentResponse, error) {
	url := fmt.Sprintf("https://svc.blockdaemon.com/boss/v1/ethereum/hoodi/stake-intents?validator_type=%s", validatorType)

	// Marshal request body to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Unmarshal response
	var stakeResp PostStakeIntentResponse
	if err := json.Unmarshal(body, &stakeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &stakeResp, nil
}

func main() {
	ctx := context.Background()

	// Get the API key from secret manager
	apiKey, err := secrets.GetSecret(
		"staging-191601",
		"default-blockdaemon-staking-api-credentials",
		"latest",
	)
	if err != nil {
		panic(err)
	}

	// Create HTTP client
	httpClient, _ := NewHTTPClient(ctx, HTTPClientOptions{})

	// Create a stake intent request
	stakeReq := PostStakeIntentRequest{
		NetworkName: Hoodi,
		PlanID:      "8ecb1a4f-225d-491c-ad5a-c33fb1770f76", // Replace with your actual plan ID
		Stakes: []StakeRequest{
			{
				Amount:            "32000000000",                                // 32 ETH in wei
				WithdrawalAddress: "0x601cae643a1cbad5509762dca13d54498eab171b", // Replace with your withdrawal address
				FeeRecipient:      "0x601cae643a1cbad5509762dca13d54498eab171b", // Replace with your fee recipient address
				Quantity:          1,
			},
		},
	}

	// Create stake intent with timing
	startTime := time.Now()
	stakeResp, err := CreateStakeIntent(ctx, httpClient, apiKey, stakeReq, ValidatorType0x02)
	elapsed := time.Since(startTime)

	if err != nil {
		fmt.Printf("Error creating stake intent: %v\n", err)
		panic(err)
	}

	fmt.Printf("Time to create stake intent: %v\n", elapsed)

	// Print response
	fmt.Printf("Stake Intent Created Successfully!\n")
	fmt.Printf("Stake Intent ID: %s\n", stakeResp.StakeIntentID)
	for _, stake := range stakeResp.Ethereum.Stakes {
		fmt.Printf("  Stake ID: %s, Amount: %s, Validator Public Key: %s\n",
			stake.StakeID, stake.Amount, stake.ValidatorPublicKey)
	}
}

// NewRoundTripper exists to give us an easy way to create new http clients with good configs and proper connection cleanup
func NewRoundTripper(ctx context.Context, opts HTTPClientOptions) (http.RoundTripper, func()) {
	optionsWithDefault := opts.withDefaults()

	rawTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2: true,
		// We only connect to one host here, so it' set to the same number
		MaxIdleConnsPerHost: optionsWithDefault.MaxIdleConns,
		MaxIdleConns:        optionsWithDefault.MaxIdleConns,

		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		// This is customized. It's necessary for making many short lived requests in parallel
		// If this number is exceeded, new requests will fail with "dial tcp 127.0.0.1:8332: connect: cannot assign requested address"
		ExpectContinueTimeout: 1 * time.Second,
		MaxConnsPerHost:       optionsWithDefault.MaxConnsPerHost,
	}
	if optionsWithDefault.InsecureTLS {
		rawTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: optionsWithDefault.InsecureTLS}
	}
	roundTripper := http.RoundTripper(rawTransport)
	return roundTripper, rawTransport.CloseIdleConnections
}

// NewHTTPClient exists to give us an easy way to create new http clients with good configs and proper connection cleanup
func NewHTTPClient(ctx context.Context, opts HTTPClientOptions) (http.Client, func()) {
	optionsWithDefault := opts.withDefaults()

	rt := optionsWithDefault.RoundTripper
	if optionsWithDefault.RoundTripper == nil {
		rt, _ = NewRoundTripper(ctx, optionsWithDefault)
	}

	client := http.Client{
		Transport: rt,
		Timeout:   optionsWithDefault.Timeout,
	}
	return client, client.CloseIdleConnections // CloseIdleConnections will delegate to the underlying transport
}
