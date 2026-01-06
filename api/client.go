package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func GetHTTPClient() *http.Client {
	return httpClient
}

func SendRequestWithContext(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if API_KEY != "" {
		req.Header.Add("Authorization", "Bearer "+API_KEY)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseMap map[string]any
	if err := json.Unmarshal(data, &responseMap); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if errStr, ok := responseMap["error"].(string); ok {
		return nil, fmt.Errorf("API error: %s", errStr)
	}

	return data, nil
}
