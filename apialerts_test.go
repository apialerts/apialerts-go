package apialerts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApiAlertsClient(t *testing.T) {
	config := Config{
		Timeout: 45 * time.Second,
		Debug:   true,
	}
	ConfigureWithConfig("test_api_key", config)

	if instance.apiKey != "test_api_key" {
		t.Errorf("Expected API key to be 'test_api_key', got '%s'", instance.apiKey)
	}

	if instance.config.Timeout != 45*time.Second {
		t.Errorf("Expected timeout to be 45 seconds, got %v", instance.config.Timeout)
	}

	if !instance.config.Debug {
		t.Errorf("Expected Debug to be false, got %v", instance.config.Debug)
	}
}

func TestSetApiKey(t *testing.T) {
	Configure("test_api_key")
	SetApiKey("new_api_key")

	if instance.apiKey != "new_api_key" {
		t.Errorf("Expected API key to be 'new_api_key', got '%s'", instance.apiKey)
	}
}

func TestSendAsync(t *testing.T) {
	server := createTestServer(t)
	defer server.Close()

	Configure("test_api_key")

	event := Event{
		Channel: "test_channel",
		Message: "Test message",
		Tags:    []string{"tag1", "tag2"},
		Link:    "https://example.com",
	}

	err := sendToUrlWithApiKeyAsync(
		server.URL,
		"test_api_key",
		event)

	if err != nil {
		t.Errorf("Error sending message: %v", err)
	}
}

func TestSend(t *testing.T) {
	server := createTestServer(t)
	defer server.Close()

	Configure("test_api_key")

	event := Event{
		Channel: "test_channel",
		Message: "Test message",
		Tags:    []string{"tag1", "tag2"},
		Link:    "https://example.com",
	}

	err := sendToUrlWithApiKeyAsync(
		server.URL,
		"test_api_key",
		event)

	if err != nil {
		t.Errorf("Error sending message: %v", err)
	}

	err = sendToUrlWithApiKeyAsync(
		server.URL,
		"test_api_key",
		Event{
			Channel: "test_channel",
		})

	if err == nil || err.Error() != "x (apialerts.com) Error: message is required" {
		t.Errorf("Expected 'message is required' error, got %v", err)
	}
}

func createTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test_api_key" {
			t.Errorf("Expected Authorization header to be 'Bearer test_api_key'")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be 'application/json'")
		}
		if r.Header.Get("X-Integration") != "golang" {
			t.Errorf("Expected X-Integration header to be 'golang'")
		}
		if r.Header.Get("X-Version") != IntegrationVersion {
			t.Errorf("Expected X-Version header to be '%s'", IntegrationVersion)
		}

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Errorf("Error decoding payload: %v", err)
		}

		if payload["message"] != "Test message" {
			t.Errorf("Expected message to be 'Test message', got '%v'", payload["message"])
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"workspace": "test_workspace", "channel": "test_channel"})
		if err != nil {
			t.Errorf("Error encoding response: %v", err)
		}
	}))
}
