package main

import (
	"context"
	"fmt"
	"luganodes"
)

func main() {
	ctx := context.Background()
	baseURL := "https://testnet.eth-staking.lgns.net"

	authClient := luganodes.NewAuthClient(baseURL)

	email := "jose.sousa@anchorlabs.com"
	password := "SuperSecret123!"
	orgName := "Anchorage"
	signupResp, err := authClient.Signup(ctx, email, password, orgName)
	if err != nil {
		panic(err)
	}
	apiKey := signupResp.Result.User.APIKey
	fmt.Println("New API key:", apiKey)

	// OR Login to existing account
	loginResp, err := authClient.Login(ctx, email, password)
	apiKey = loginResp.Result.User.APIKey

	// 3) Create Luganodes REST client with API key
	client := luganodes.NewClient(apiKey, baseURL)

	// 1) Provision
	provReq := luganodes.ProvisionRequest{
		WithdrawalAddress:  "0xyourWithdrawAddr",
		ValidatorsCount:    2,
		Batch:              false,
		Compounding:        false,
		AmountPerValidator: 32,
	}

	provResp, _ := client.CreateProvision(ctx, provReq)
	fmt.Println("Provision ID:", provResp.ProvisionId)

	// 2) Fetch validator objects
	valObjs, _ := client.GetValidatorObjects(ctx, provResp.ProvisionId, 1, 20)
	fmt.Println("Validator Objects:", valObjs.Result)

	// 3) EXIT – sign the challenge locally and call
	challenge := "<packed challenge>" // derived off-chain
	sig, err := luganodes.SignMessage([]byte(challenge), "<private key hex>")
	if err != nil {
		panic(err)
	}

	exitResp, _ := client.GenerateExitMessage(ctx, provReq.WithdrawalAddress, challenge, sig)
	fmt.Println("Exit Message:", exitResp.Message)
}
