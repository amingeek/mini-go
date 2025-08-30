package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type responseData struct {
	Status string `json:"status"`
	Stats  map[string]struct {
		Latest string `json:"latest"`
	} `json:"stats"`
}

func GetExchangeRate(source, destination string) (string, error) {
	client := &http.Client{}

	srcLower := strings.ToLower(source)
	dstLower := strings.ToLower(destination)

	if dstLower == "" {
		dstLower = "rls"
	}

	baseURL := "http://localhost:4001/rates"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing base URL: %w", err)
	}

	query := reqURL.Query()
	query.Set("srcCurrency", srcLower)
	query.Set("dstCurrency", dstLower)
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var data responseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response JSON: %w", err)
	}

	key := fmt.Sprintf("%s-%s", srcLower, dstLower)

	if val, ok := data.Stats[key]; ok && data.Status == "OK" {
		return val.Latest, nil
	}

	return "", fmt.Errorf("no rate found or status not OK")
}
