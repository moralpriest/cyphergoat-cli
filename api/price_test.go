package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestGetPrice_Success(t *testing.T) {
	expectedPrice := 45250.50
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ids") != "bitcoin" {
			t.Errorf("Expected coin ID 'bitcoin', got '%s'", r.URL.Query().Get("ids"))
		}
		if r.URL.Query().Get("vs_currencies") != "usd" {
			t.Errorf("Expected currency 'usd', got '%s'", r.URL.Query().Get("vs_currencies"))
		}
		response := map[string]map[string]float64{
			"bitcoin": {"usd": expectedPrice},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	testService := NewPriceServiceWithURL(server.URL)
	testService.lastCall = time.Now().Add(-1 * time.Second)

	price, err := testService.GetPrice(context.Background(), "btc")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if price != expectedPrice {
		t.Errorf("Expected price %f, got %f", expectedPrice, price)
	}
}

func TestGetPrice_Stablecoin(t *testing.T) {
	testCases := []struct {
		coin string
		want float64
	}{
		{"usdt", 1.0},
		{"usdc", 1.0},
		{"dai", 1.0},
		{"busd", 1.0},
		{"TUSd", 1.0},
		{"USDT", 1.0},
	}

	for _, tc := range testCases {
		if !IsStablecoin(tc.coin) {
			t.Errorf("Expected %s to be stablecoin", tc.coin)
		}
	}
}

func TestGetPrice_CacheHit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Server should not be called for cached price")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cachedPrice := 3250.00
	testService := NewPriceServiceWithURL(server.URL)
	testService.cache = map[string]PriceCache{
		"eth": {
			price:     cachedPrice,
			timestamp: time.Now(),
		},
	}

	price, err := testService.GetPrice(context.Background(), "eth")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if price != cachedPrice {
		t.Errorf("Expected cached price %f, got %f", cachedPrice, price)
	}
}

func TestGetPrice_CacheExpired(t *testing.T) {
	expectedPrice := 85.50
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("ids")
		response := make(map[string]map[string]float64)

		if ids == "litecoin" {
			response = map[string]map[string]float64{
				"litecoin": {"usd": expectedPrice},
			}
		} else if ids == "bitcoin" {
			response = map[string]map[string]float64{
				"bitcoin": {"usd": 45250.50},
			}
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	testService := NewPriceServiceWithURL(server.URL)
	testService.cache = map[string]PriceCache{
		"ltc": {
			price:     80.00,
			timestamp: time.Now().Add(-10 * time.Minute),
		},
	}
	testService.lastCall = time.Now().Add(-1 * time.Second)

	price, err := testService.GetPrice(context.Background(), "ltc")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if price != expectedPrice {
		t.Errorf("Expected fresh price %f, got %f", expectedPrice, price)
	}
}

func TestGetPrice_UnknownCoin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]map[string]float64{}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	testService := NewPriceServiceWithURL(server.URL)
	testService.lastCall = time.Now().Add(-1 * time.Second)

	_, err := testService.GetPrice(context.Background(), "unknowncoin123")
	if err == nil {
		t.Error("Expected error for unknown coin, got nil")
	}
}

func TestGetPrice_RateLimit(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		response := map[string]map[string]float64{
			"sol": {"usd": 100.00},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	testService := NewPriceServiceWithURL(server.URL)
	testService.lastCall = time.Now().Add(-1 * time.Second)

	ctx := context.Background()

	testService.GetPrice(ctx, "sol")
	testService.GetPrice(ctx, "sol")
	testService.GetPrice(ctx, "sol")

	if callCount != 3 {
		t.Errorf("Expected 3 API calls (no cache), got %d", callCount)
	}
}

func TestGetPrice_PrivacyCoins(t *testing.T) {
	testCoins := []string{"xmr", "arrr", "zec", "dero", "wow", "firo", "zano", "dash", "bdx", "ban"}

	for _, coin := range testCoins {
		coinID := getCoinGeckoID(coin)
		if coinID == "" {
			t.Errorf("Missing CoinGecko ID for privacy coin: %s", coin)
		}
		if IsStablecoin(coin) {
			t.Errorf("Privacy coin %s incorrectly identified as stablecoin", coin)
		}
	}
}

func TestGetPrices_Batch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("ids")
		response := make(map[string]map[string]float64)

		if ids == "bitcoin" {
			response = map[string]map[string]float64{
				"bitcoin": {"usd": 45250.50},
			}
		} else if ids == "ethereum" {
			response = map[string]map[string]float64{
				"ethereum": {"usd": 3250.00},
			}
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	testService := NewPriceServiceWithURL(server.URL)
	testService.lastCall = time.Now().Add(-1 * time.Second)

	ctx := context.Background()
	prices, err := testService.GetPrices(ctx, []string{"btc", "eth"})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(prices) != 2 {
		t.Errorf("Expected 2 prices, got %d", len(prices))
	}

	if prices["btc"] != 45250.50 {
		t.Errorf("Expected BTC price 45250.50, got %f", prices["btc"])
	}

	if prices["eth"] != 3250.00 {
		t.Errorf("Expected ETH price 3250.00, got %f", prices["eth"])
	}
}

func TestPriceService_ClearCache(t *testing.T) {
	testService := NewPriceService()
	testService.cache = map[string]PriceCache{
		"btc": {price: 45000, timestamp: time.Now()},
		"eth": {price: 3000, timestamp: time.Now()},
	}

	if testService.CacheSize() != 2 {
		t.Errorf("Expected cache size 2, got %d", testService.CacheSize())
	}

	testService.ClearCache()

	if testService.CacheSize() != 0 {
		t.Errorf("Expected cache size 0 after clear, got %d", testService.CacheSize())
	}
}

func TestIsStablecoin(t *testing.T) {
	testCases := []struct {
		coin     string
		expected bool
	}{
		{"usdt", true},
		{"USDT", true},
		{"usdc", true},
		{"dai", true},
		{"busd", true},
		{"btc", false},
		{"eth", false},
		{"xmr", false},
	}

	for _, tc := range testCases {
		result := IsStablecoin(tc.coin)
		if result != tc.expected {
			t.Errorf("IsStablecoin(%s): expected %v, got %v", tc.coin, tc.expected, result)
		}
	}
}

func TestGetPrice_ConcurrentAccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("ids")
		response := make(map[string]map[string]float64)

		if ids == "bitcoin" {
			response = map[string]map[string]float64{
				"bitcoin": {"usd": 45000.00},
			}
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	testService := NewPriceServiceWithURL(server.URL)
	testService.lastCall = time.Now().Add(-10 * time.Second)

	ctx := context.Background()
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := testService.GetPrice(ctx, "btc")
			if err != nil {
				errors <- err
			}
		}()
	}
	wg.Wait()
	close(errors)

	errorCount := 0
	for err := range errors {
		t.Logf("Error: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Errorf("Expected no errors in concurrent access, got %d", errorCount)
	}
}
