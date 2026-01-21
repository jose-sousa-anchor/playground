package luganodes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ProvisionRequest struct {
	WithdrawalAddress  string  `json:"withdrawalAddress"`
	ControllerAddress  string  `json:"controllerAddress,omitempty"`
	FeeRecipient       string  `json:"feeRecipient,omitempty"`
	ValidatorsCount    int     `json:"validatorsCount"`
	Batch              bool    `json:"batch"`
	Compounding        bool    `json:"compounding"`
	AmountPerValidator float64 `json:"amountPerValidator"`
}

type ProvisionResponse struct {
	ProvisionId       string `json:"provisionId"`
	WithdrawalAddress string `json:"withdrawalAddress"`
	Status            string `json:"status"`
	Created           string `json:"created"`
	ValidatorsCount   int    `json:"validatorsCount"`
	ControllerAddress string `json:"controllerAddress"`
	FeeRecipient      string `json:"feeRecipient"`
}

func (c *Client) CreateProvision(
	ctx context.Context,
	reqBody ProvisionRequest,
) (*ProvisionResponse, error) {
	url := fmt.Sprintf("%s/api/provision", c.BaseURL)

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, _, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result ProvisionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type ValidatorObjectsResponse struct {
	Result []struct {
		Amount           float64 `json:"amount"`
		ValidatorIndex   int     `json:"validatorIndex"`
		Status           string  `json:"status"`
		ValidatorAddress string  `json:"validatorAddress"`
		DepositInput     string  `json:"depositInput"`
	} `json:"result"`
}

func (c *Client) GetValidatorObjects(
	ctx context.Context,
	provisionId string,
	page, perPage int,
) (*ValidatorObjectsResponse, error) {
	url := fmt.Sprintf("%s/api/validators?provisionId=%s&page=%d&per_page=%d",
		c.BaseURL, provisionId, page, perPage)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	body, _, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result ValidatorObjectsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
