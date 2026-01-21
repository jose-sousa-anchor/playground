package vemflow

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	p2pclient "p2p/client"
	"time"
)

type VemStatusResponse struct {
	Status    string `json:"status"`
	VemResult string `json:"vemResult"`
}

func PollVemResult(
	ctx context.Context,
	client *http.Client,
	token string,
	requestID string,
) (string, error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			req, _ := p2pclient.NewVemStatusRequest(ctx, requestID, token)
			body, _, err := p2pclient.DoRequest(client, req)
			if err != nil {
				return "", err
			}

			var resp VemStatusResponse
			_ = json.Unmarshal(body, &resp)

			switch resp.Status {
			case "success":
				return resp.VemResult, nil
			case "error", "fault":
				return "", errors.New("vem request failed")
			}
		}
	}
}
