package api

import (
	"encoding/json"
	"testing"
)

func TestFetchEstimate_ObjectFormat(t *testing.T) {
	mockResponse := `{
		"min": 0.04444,
		"rates": {
			"Results": [
				{
					"Exchange": "PegasusSwap",
					"Amount": 0.18522283,
					"MinAmount": 0.001
				},
				{
					"Exchange": "SimpleSwap",
					"Amount": 0.1845,
					"MinAmount": 0.001
				},
				{
					"Exchange": "ChangeNow",
					"Amount": 0.1860,
					"MinAmount": 0.001
				}
			],
			"Min": 0,
			"TradeValue_fiat": 100.0,
			"TradeValue_btc": 0.001
		}
	}`

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
	err := json.Unmarshal([]byte(mockResponse), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal object format: %v", err)
	}

	if len(result.Rates.Results) != 3 {
		t.Errorf("Expected 3 rates, got %d", len(result.Rates.Results))
	}

	estimates := populateEstimates(result.Rates.Results, "btc", "eth", 1.0, "btc", "eth", result.Rates.TradeValue_fiat)
	if len(estimates) != 3 {
		t.Errorf("Expected 3 estimates after population, got %d", len(estimates))
	}

	if estimates[0].ReceiveAmount != 0.1860 {
		t.Errorf("Expected first estimate to have highest ReceiveAmount (0.1860), got %f", estimates[0].ReceiveAmount)
	}

	for _, est := range estimates {
		if est.Coin1 != "btc" {
			t.Errorf("Expected Coin1 to be 'btc', got '%s'", est.Coin1)
		}
		if est.Coin2 != "eth" {
			t.Errorf("Expected Coin2 to be 'eth', got '%s'", est.Coin2)
		}
		if est.SendAmount != 1.0 {
			t.Errorf("Expected SendAmount to be 1.0, got %f", est.SendAmount)
		}
	}
}

func TestFetchEstimate_Sorting(t *testing.T) {
	mockResponse := `{
		"min": 0.001,
		"rates": {
			"Results": [
				{"Exchange": "low", "Amount": 0.0100, "MinAmount": 0.001},
				{"Exchange": "high", "Amount": 0.0200, "MinAmount": 0.001},
				{"Exchange": "medium", "Amount": 0.0150, "MinAmount": 0.001}
			],
			"Min": 0,
			"TradeValue_fiat": 50.0,
			"TradeValue_btc": 0.0005
		}
	}`

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
	if err := json.Unmarshal([]byte(mockResponse), &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	estimates := populateEstimates(result.Rates.Results, "btc", "eth", 1.0, "btc", "eth", result.Rates.TradeValue_fiat)

	if estimates[0].ReceiveAmount != 0.0200 {
		t.Errorf("Expected highest amount first, got %f", estimates[0].ReceiveAmount)
	}
	if estimates[1].ReceiveAmount != 0.0150 {
		t.Errorf("Expected second highest amount second, got %f", estimates[1].ReceiveAmount)
	}
	if estimates[2].ReceiveAmount != 0.0100 {
		t.Errorf("Expected lowest amount last, got %f", estimates[2].ReceiveAmount)
	}
}

func TestFetchEstimate_SingleEstimate(t *testing.T) {
	mockResponse := `{
		"min": 0.001,
		"rates": {
			"Results": [
				{"Exchange": "only", "Amount": 0.0150, "MinAmount": 0.001}
			],
			"Min": 0,
			"TradeValue_fiat": 50.0,
			"TradeValue_btc": 0.0005
		}
	}`

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
	err := json.Unmarshal([]byte(mockResponse), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	estimates := populateEstimates(result.Rates.Results, "btc", "eth", 1.0, "btc", "eth", result.Rates.TradeValue_fiat)
	if len(estimates) != 1 {
		t.Errorf("Expected 1 estimate, got %d", len(estimates))
	}

	if estimates[0].ExchangeName != "only" {
		t.Errorf("Expected exchange name 'only', got '%s'", estimates[0].ExchangeName)
	}
}

func TestFetchEstimate_EmptyResponse(t *testing.T) {
	mockResponse := `{
		"min": 0,
		"rates": {
			"Results": [],
			"Min": 0,
			"TradeValue_fiat": 0,
			"TradeValue_btc": 0
		}
	}`

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
	err := json.Unmarshal([]byte(mockResponse), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	estimates := populateEstimates(result.Rates.Results, "btc", "eth", 1.0, "btc", "eth", result.Rates.TradeValue_fiat)
	if len(estimates) != 0 {
		t.Errorf("Expected 0 estimates, got %d", len(estimates))
	}
}
