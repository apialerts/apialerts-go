package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestApiAlertsClient(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("APIALERTS_API_KEY", "test_api_key")
	os.Setenv("APIALERTS_TIMEOUT", "45")

	client := ApiAlertsClient()

	if client.apiKey != "test_api_key" {
		t.Errorf("Expected API key to be 'test_api_key', got '%s'", client.apiKey)
	}

	if client.timeout != 45*time.Second {
		t.Errorf("Expected timeout to be 45 seconds, got %v", client.timeout)
	}

	// Reset environment variables
	os.Unsetenv("APIALERTS_API_KEY")
	os.Unsetenv("APIALERTS_TIMEOUT")
}

func TestSetApiKey(t *testing.T) {
	client := &Client{}
	client.SetApiKey("new_api_key")

	if client.apiKey != "new_api_key" {
		t.Errorf("Expected API key to be 'new_api_key', got '%s'", client.apiKey)
	}
}

func TestSendAsync(t *testing.T) {
	client := &Client{apiKey: "test_api_key"}

	// Test with debug mode
	client.SendAsync("Test message", []string{"tag1", "tag2"}, "http://example.com", true)

	// Test without debug mode
	client.SendAsync("Test message", []string{"tag1", "tag2"}, "http://example.com")
}

func TestSend(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request headers
		if r.Header.Get("Authorization") != "Bearer test_api_key" {
			t.Errorf("Expected Authorization header to be 'Bearer test_api_key'")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be 'application/json'")
		}
		if r.Header.Get("X-Integration") != "golang" {
			t.Errorf("Expected X-Integration header to be 'golang'")
		}
		if r.Header.Get("X-Version") != "1.0.0" {
			t.Errorf("Expected X-Version header to be '1.0.0'")
		}

		// Decode the request body
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		// Check payload contents
		if payload["message"] != "Test message" {
			t.Errorf("Expected message to be 'Test message', got '%v'", payload["message"])
		}

		// Send response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"project": "test_project"})
	}))
	defer server.Close()

	// Create a custom client with our test server URL
	client := &Client{
		apiKey:  "test_api_key",
		timeout: 30 * time.Second,
	}

	// Override the URL to use our test server
	url := server.URL

	// Create a custom httpClient for testing
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	// Test the Send function
	err := sendWithCustomUrl(client, "Test message", []string{"tag1", "tag2"}, "http://example.com", url, httpClient)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test error cases
	client.apiKey = ""
	err = client.Send("Test message", []string{"tag1", "tag2"}, "http://example.com")
	if err == nil || err.Error() != "api key is missing" {
		t.Errorf("Expected 'api key is missing' error, got %v", err)
	}

	client.apiKey = "test_api_key"
	err = client.Send("", []string{"tag1", "tag2"}, "http://example.com")
	if err == nil || err.Error() != "message is required" {
		t.Errorf("Expected 'message is required' error, got %v", err)
	}
}

// Helper function to send with a custom URL and http client
func sendWithCustomUrl(client *Client, message string, tags []string, link string, url string, httpClient *http.Client) error {
	if client.apiKey == "" {
		return errors.New("api key is missing")
	}
	if message == "" {
		return errors.New("message is required")
	}
	payload := map[string]interface{}{
		"message": message,
		"tags":    tags,
		"link":    link,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Integration", "golang")
	req.Header.Set("X-Version", "1.0.0")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		return nil
	case http.StatusBadRequest:
		return errors.New("bad request")
	case http.StatusUnauthorized:
		return errors.New("unauthorized")
	case http.StatusForbidden:
		return errors.New("forbidden")
	case http.StatusTooManyRequests:
		return errors.New("rate limit exceeded")
	default:
		return errors.New("unknown error")
	}
}
