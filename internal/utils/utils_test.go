package utils

import (
	"io"
	"net/http"
	"testing"
)

func TestNewHttpClientWithProxy(t *testing.T) {
	proxyURL := "PROXY:PORT"
	client, err := NewHttpClient(proxyURL)
	if err != nil {
		t.Fatalf("Failed to create HTTP client: %v", err)
	}
	req, err := http.NewRequest("GET", "http://icanhazip.com", nil)
	standardHeaders := StandardHeaders()
	if err != nil {
		t.Fatalf("Failed to construct request: %v", err)
	}
	for key, value := range standardHeaders {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	ip := string(body)
	ip = ip[:len(ip)-1]
	expectedIP := "PROXY"
	if ip != expectedIP {
		t.Errorf("Expected IP %v, got %v", expectedIP, ip)
	}
}
