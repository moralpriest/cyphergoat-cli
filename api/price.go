package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	coinGeckoURL   = "https://api.coingecko.com/api/v3/simple/price"
	cacheDuration  = 5 * time.Minute
	rateLimitDelay = 100 * time.Millisecond
)

type PriceCache struct {
	price     float64
	timestamp time.Time
}

type PriceService struct {
	baseURL  string
	client   *http.Client
	cache    map[string]PriceCache
	mutex    sync.RWMutex
	lastCall time.Time
}

func NewPriceService() *PriceService {
	return NewPriceServiceWithURL(coinGeckoURL)
}

func NewPriceServiceWithURL(baseURL string) *PriceService {
	return &PriceService{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
		cache:   make(map[string]PriceCache),
	}
}

var coinIDMap = map[string]string{
	"btc":   "bitcoin",
	"eth":   "ethereum",
	"sol":   "solana",
	"bnb":   "binancecoin",
	"xmr":   "monero",
	"arrr":  "pirate-chain",
	"zec":   "zcash",
	"dero":  "dero",
	"wow":   "wownero",
	"firo":  "zcoin",
	"zano":  "zano",
	"dash":  "dash",
	"bdx":   "beldex",
	"ban":   "banano",
	"ltc":   "litecoin",
	"bch":   "bitcoin-cash",
	"doge":  "dogecoin",
	"dot":   "polkadot",
	"link":  "chainlink",
	"avax":  "avalanche-2",
	"matic": "matic-network",
	"uni":   "uniswap",
	"shib":  "shiba-inu",
	"etc":   "ethereum-classic",
	"hbar":  "hedera-hashgraph",
	"xtz":   "tezos",
	"ada":   "cardano",
	"xrp":   "binance-peg-xrp",
	"trx":   "tron",
	"atom":  "cosmos",
	"near":  "near",
	"apt":   "aptos",
	"sui":   "sui",
	"dcr":   "decred",
	"aave":  "aave",
	"bat":   "basic-attention-token",
	"paxg":  "pax-gold",
	"hive":  "hive",
	"zen":   "horizen",
	"scrt":  "secret",
	"leo":   "leo-token",
	"tusd":  "true-usd",
	"gusd":  "gemini-dollar",
	"nvdax": "nvidia-xstock",
}

var stablecoinMap = map[string]bool{
	"usdt": true,
	"usdc": true,
	"dai":  true,
	"busd": true,
	"tusd": true,
	"gusd": true,
	"fusd": true,
	"usdd": true,
	"crv":  false,
}

func (s *PriceService) GetPrice(ctx context.Context, coin string) (float64, error) {
	coinLower := strings.ToLower(coin)

	if stablecoinMap[coinLower] {
		return 1.0, nil
	}

	if price, found := s.getCachedPrice(coinLower); found {
		return price, nil
	}

	s.rateLimit()

	coinID := getCoinGeckoID(coinLower)
	if coinID == "" {
		return 0, fmt.Errorf("unknown coin: %s", coin)
	}

	price, err := s.fetchFromCoinGecko(ctx, coinID)
	if err != nil {
		return 0, err
	}

	s.cachePrice(coinLower, price)

	return price, nil
}

func (s *PriceService) getCachedPrice(coin string) (float64, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	cache, found := s.cache[coin]
	if !found {
		return 0, false
	}

	if time.Since(cache.timestamp) < cacheDuration {
		return cache.price, true
	}

	return 0, false
}

func (s *PriceService) cachePrice(coin string, price float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache[coin] = PriceCache{
		price:     price,
		timestamp: time.Now(),
	}
}

func (s *PriceService) rateLimit() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	elapsed := time.Since(s.lastCall)
	if elapsed < rateLimitDelay {
		time.Sleep(rateLimitDelay - elapsed)
	}
	s.lastCall = time.Now()
}

func (s *PriceService) fetchFromCoinGecko(ctx context.Context, coinID string) (float64, error) {
	url := fmt.Sprintf("%s?ids=%s&vs_currencies=usd", s.baseURL, coinID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("API request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 429 {
			return 0, fmt.Errorf("rate limit exceeded")
		}
		return 0, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if data, ok := result[coinID]; ok {
		if price, ok := data["usd"]; ok {
			return price, nil
		}
	}

	return 0, fmt.Errorf("price not found for %s", coinID)
}

func GetPrice(ctx context.Context, coin string) (float64, error) {
	service := NewPriceService()
	return service.GetPrice(ctx, coin)
}

func getCoinGeckoID(coin string) string {
	return coinIDMap[strings.ToLower(coin)]
}

func IsStablecoin(coin string) bool {
	return stablecoinMap[strings.ToLower(coin)]
}

func (s *PriceService) GetPrices(ctx context.Context, coins []string) (map[string]float64, error) {
	prices := make(map[string]float64)
	var errors []error

	for _, coin := range coins {
		price, err := s.GetPrice(ctx, coin)
		if err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", coin, err))
			continue
		}
		prices[coin] = price
	}

	if len(errors) > 0 && len(prices) == 0 {
		return nil, fmt.Errorf("all price lookups failed: %v", errors)
	}

	return prices, nil
}

func (s *PriceService) ClearCache() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cache = make(map[string]PriceCache)
}

func (s *PriceService) CacheSize() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.cache)
}
