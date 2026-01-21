package luganodes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ExitChallengeRequest struct {
	Challenge string `json:"challenge"`
	Signature string `json:"signature"`
}

type ExitResponse struct {
	Message string `json:"message"`
}

func (c *Client) SubmitExit(
	ctx context.Context,
	keyAddr string,
	challenge, signature string,
) (*ExitResponse, error) {
	url := fmt.Sprintf("%s/api/exit?key=%s", c.BaseURL, keyAddr)

	bodyReq := ExitChallengeRequest{Challenge: challenge, Signature: signature}
	b, _ := json.Marshal(bodyReq)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")

	body, _, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var resp ExitResponse
	json.Unmarshal(body, &resp)
	return &resp, nil
}

func (c *Client) GenerateExitMessage(
	ctx context.Context,
	keyAddr, challenge, signature string,
) (*ExitResponse, error) {
	url := fmt.Sprintf("%s/api/exit/message?key=%s", c.BaseURL, keyAddr)

	bodyReq := ExitChallengeRequest{Challenge: challenge, Signature: signature}
	b, _ := json.Marshal(bodyReq)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")

	body, _, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var resp ExitResponse
	json.Unmarshal(body, &resp)
	return &resp, nil
}
