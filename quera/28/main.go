package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetExchangeRate(source, destination string) (string, error) {
	client := &http.Client{}

	baseURL := "http://localhost:4001/rates"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing base URL: %w", err)
	}
	if destination == "" {
		destination = "rls"
	}

	query := reqURL.Query()
	query.Set("srcCurrency", source)
	query.Set("dstCurrency", destination)
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

	return string(body), nil
}
