package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
)

const URL = "api.cyphergoat.com"

var API_KEY string

type Estimate struct {
	ExchangeName  string  `json:"Exchange"`
	ReceiveAmount float64 `json:"Amount"`
	MinAmount     float64 `json:"MinAmount"`
	KYCScore      int     `json:"KYCScore"`
	Network1      string
	Network2      string
	Coin1         string
	Coin2         string
	SendAmount    float64
	Address       string
	ImageURL      string
	TradeValueUSD float64
}

type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
}
type Transaction struct {
	Coin1          string    `json:"Coin1,omitempty"`
	Coin2          string    `json:"Coin2,omitempty"`
	Network1       string    `json:"Network1,omitempty"`
	Network2       string    `json:"Network2,omitempty"`
	Address        string    `json:"Address,omitempty"`
	EstimateAmount float64   `json:"EstimateAmount,omitempty"`
	Provider       string    `json:"Provider,omitempty"`
	Id             string    `json:"Id,omitempty"`
	SendAmount     float64   `json:"SendAmount,omitempty"`
	Track          string    `json:"Track,omitempty"`
	Status         string    `json:"Status,omitempty"`
	KYC            string    `json:"KYC,omitempty"`
	Token          string    `json:"Token,omitempty"`
	Done           bool      `json:"Done,omitempty"`
	CGID           string    `json:"CGID,omitempty"`
	CreatedAt      time.Time `json:"CreatedAt,omitempty"`
}

func init() {
	API_KEY = GetAPIKeyFromEnv()
}

func GetAPIKeyFromEnv() string {
	if key := os.Getenv("CYPHERGOAT_API_KEY"); key != "" {
		return key
	}
	return os.Getenv("API_KEY")
}

func GetAPIKey() string {
	return API_KEY
}

func FetchEstimateFromAPI(ctx context.Context, coin1, coin2 string, amount float64, best bool, network1, network2 string) ([]Estimate, error) {
	params := url.Values{}
	params.Set("coin1", coin1)
	params.Set("coin2", coin2)
	params.Set("amount", fmt.Sprintf("%f", amount))
	params.Set("network1", network1)
	params.Set("network2", network2)
	if best {
		params.Set("best", "true")
	}

	requestURL := fmt.Sprintf("https://%s/estimate?%s", URL, params.Encode())

	data, err := SendRequestWithContext(ctx, requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch estimate: %w", err)
	}

	type RatesWrapper struct {
		Results         []Estimate `json:"Results"`
		Min             float64    `json:"Min"`
		TradeValue_fiat float64    `json:"TradeValue_fiat"`
		TradeValue_btc  float64    `json:"TradeValue_btc"`
	}

	type ApiResponse struct {
		Min   float64      `json:"min"`
		Rates RatesWrapper `json:"rates"`
	}

	var result ApiResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal estimate response: %w", err)
	}

	priceService := NewPriceService()
	coin2USDPrice, err := priceService.GetPrice(ctx, coin2)
	if err != nil {
		coin2USDPrice = 0
	}

	estimates := populateEstimates(result.Rates.Results, coin1, coin2, amount, network1, network2, coin2USDPrice)
	return estimates, nil
}

func populateEstimates(estimates []Estimate, coin1, coin2 string, amount float64, network1, network2 string, coin2USDPrice float64) []Estimate {
	for i := range estimates {
		estimates[i].Coin1 = coin1
		estimates[i].Coin2 = coin2
		estimates[i].SendAmount = amount
		estimates[i].Network1 = network1
		estimates[i].Network2 = network2
		estimates[i].TradeValueUSD = estimates[i].ReceiveAmount * coin2USDPrice
	}
	slices.SortFunc(estimates, func(a, b Estimate) int {
		if a.ReceiveAmount > b.ReceiveAmount {
			return -1
		}
		if a.ReceiveAmount < b.ReceiveAmount {
			return 1
		}
		return 0
	})
	return estimates
}

func CreateTradeFromAPI(ctx context.Context, coin1, coin2 string, amount float64, address, partner string, network1, network2 string) (Transaction, error) {
	params := url.Values{}
	params.Set("coin1", coin1)
	params.Set("coin2", coin2)
	params.Set("amount", fmt.Sprintf("%f", amount))
	params.Set("partner", partner)
	params.Set("address", address)
	params.Set("network1", network1)
	params.Set("network2", network2)

	requestURL := fmt.Sprintf("https://%s/swap?%s", URL, params.Encode())

	data, err := SendRequestWithContext(ctx, requestURL)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to create trade: %w", err)
	}

	var result TransactionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return Transaction{}, fmt.Errorf("failed to unmarshal trade response: %w", err)
	}

	transaction := result.Transaction
	return transaction, nil
}

func TrackTxFromAPI(ctx context.Context, t Transaction) (Transaction, error) {
	requestURL := fmt.Sprintf("https://%s/transaction?id=%s", URL, strings.ToLower(t.Provider))

	data, err := SendRequestWithContext(ctx, requestURL)
	if err != nil {
		return t, fmt.Errorf("failed to track transaction: %w", err)
	}

	var responseMap map[string]any
	if err := json.Unmarshal(data, &responseMap); err != nil {
		return t, fmt.Errorf("failed to unmarshal track response: %w", err)
	}

	status, ok := responseMap["status"].(string)
	if !ok {
		return t, fmt.Errorf("status field is missing or not a string")
	}
	t.Status = status

	return t, nil
}

func GetTransactionFromAPI(ctx context.Context, id string) (Transaction, error) {
	requestURL := fmt.Sprintf("https://%s/transaction?id=%s", URL, id)

	data, err := SendRequestWithContext(ctx, requestURL)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to get transaction: %w", err)
	}

	var result map[string]Transaction
	if err := json.Unmarshal(data, &result); err != nil {
		return Transaction{}, fmt.Errorf("failed to unmarshal transaction response: %w", err)
	}

	transaction := result["transaction"]
	return transaction, nil
}
