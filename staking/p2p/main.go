package main

import (
	"context"
	"fmt"
	p2pclient "p2p/client"
	"p2p/vemcrypto"
	"p2p/vemflow"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	token := "<P2P_BEARER_TOKEN>"
	client := p2pclient.NewHTTPClient()

	/*
	   ----------------------------------------------------------------
	   1. PROVISION: CREATE NODE REQUEST
	   ----------------------------------------------------------------
	*/

	provisionID := "3611b95c-e1b3-40c0-9086-3de0a4379943"

	createPayload := p2pclient.CreateNodeRequestPayload{
		ID:                        provisionID,
		Type:                      "REGULAR",
		ValidatorsCount:           2,
		AmountPerValidator:        "32000000000",
		WithdrawalCredentialsType: "0x01",
		WithdrawalAddress:         "0x39D02C253dA1d9F85ddbEB3B6Dc30bc1EcBbFA17",
		EigenPodOwnerAddress:      "",
		ControllerAddress:         "0x39D02C253dA1d9F85ddbEB3B6Dc30bc1EcBbFA17",
		FeeRecipientAddress:       "0x53da3c92fCCEb0CFE1764f65DDfF1564A2b15585",
		NodesOptions: p2pclient.NodesOptionsInput{
			Location:  "any",
			RelaysSet: "",
		},
	}

	createReq, err := p2pclient.NewCreateNodeRequest(ctx, createPayload, token)
	if err != nil {
		panic(err)
	}

	body, status, err := p2pclient.DoRequest(client, createReq)
	if err != nil {
		panic(err)
	}

	fmt.Println("Provision create status:", status)
	fmt.Println("Provision create response:", string(body))

	/*
	   ----------------------------------------------------------------
	   2. PROVISION: CHECK STATUS
	   ----------------------------------------------------------------
	*/

	statusReq, err := p2pclient.NewGetNodeRequestStatusRequest(
		ctx,
		provisionID,
		token,
	)
	if err != nil {
		panic(err)
	}

	statusBody, statusCode, err := p2pclient.DoRequest(client, statusReq)
	if err != nil {
		panic(err)
	}

	fmt.Println("Provision status HTTP:", statusCode)
	fmt.Println("Provision status response:", string(statusBody))

	/*
	   ----------------------------------------------------------------
	   3. EXIT: GENERATE ECDH KEYS
	   ----------------------------------------------------------------
	*/

	ecdhPrivKey, ecdhPubKeyBase64, err := vemcrypto.GenerateECDHKeypair()
	if err != nil {
		panic(err)
	}

	/*
	   ----------------------------------------------------------------
	   4. EXIT: BUILD & SIGN VEM REQUEST
	   ----------------------------------------------------------------
	*/

	// Inner VEM request MUST be a JSON STRING
	vemRequest := fmt.Sprintf(
		`{"action":"vem_request","pubkeys":["0xVALIDATOR_BLS_PUBKEY"],"ecdh_client_pubkey":"%s"}`,
		ecdhPubKeyBase64,
	)

	// Ethereum signature of vemRequest (EIP-191)
	// This MUST be produced by your custody / signer
	vemSignature := "<ETH_SIGNATURE_HEX>"
	vemSignedBy := "0x39D02C253dA1d9F85ddbEB3B6Dc30bc1EcBbFA17"

	vemID := "uuid-vem-request-id"

	vemPayload := p2pclient.VemCreatePayload{
		ID:                  vemID,
		Type:                "off_chain",
		VemRequest:          vemRequest,
		VemRequestSignature: vemSignature,
		VemRequestSignedBy:  vemSignedBy,
	}

	/*
	   ----------------------------------------------------------------
	   5. EXIT: CREATE VEM
	   ----------------------------------------------------------------
	*/

	vemCreateReq, err := p2pclient.NewVemCreateRequest(ctx, vemPayload, token)
	if err != nil {
		panic(err)
	}

	vemCreateBody, vemCreateStatus, err := p2pclient.DoRequest(client, vemCreateReq)
	if err != nil {
		panic(err)
	}

	fmt.Println("VEM create HTTP:", vemCreateStatus)
	fmt.Println("VEM create response:", string(vemCreateBody))

	/*
	   ----------------------------------------------------------------
	   6. EXIT: POLL FOR SIGNED EXIT MESSAGE
	   ----------------------------------------------------------------
	*/

	encryptedVemResult, err := vemflow.PollVemResult(
		ctx,
		client,
		token,
		vemID,
	)
	if err != nil {
		panic(err)
	}

	/*
	   ----------------------------------------------------------------
	   7. EXIT: DECRYPT SIGNED VALIDATOR EXIT MESSAGE
	   ----------------------------------------------------------------
	*/

	signedExitMessage, err := vemcrypto.DecryptVemResult(
		ecdhPrivKey,
		encryptedVemResult,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("SIGNED VALIDATOR EXIT MESSAGE:")
	fmt.Println(string(signedExitMessage))
}
