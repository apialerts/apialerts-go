package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/apialerts/apialerts-go/model"
)

func TestApiAlertsClient(t *testing.T) {
	config := model.APIAlertsConfig{
		Logging: true,
		Timeout: 45 * time.Second,
		Debug:   false,
	}

	client := ApiAlertsClientWithConfig("test_api_key", config)

	if client.ApiKey != "test_api_key" {
		t.Errorf("Expected API key to be 'test_api_key', got '%s'", client.ApiKey)
	}

	if client.Config.Timeout != 45*time.Second {
		t.Errorf("Expected timeout to be 45 seconds, got %v", client.Config.Timeout)
	}
}

func TestSetApiKey(t *testing.T) {
	client := &APIAlertsClient{}
	client.SetApiKey("new_api_key")

	if client.ApiKey != "new_api_key" {
		t.Errorf("Expected API key to be 'new_api_key', got '%s'", client.ApiKey)
	}
}

func TestSendAsync(t *testing.T) {
	config := model.APIAlertsConfig{
		Logging: true,
		Timeout: 45 * time.Second,
		Debug:   true,
	}

	client := ApiAlertsClientWithConfig("test_api_key", config)

	event := model.APIAlertsEvent{
		Channel: "test_channel",
		Message: "Test message",
		Tags:    []string{"tag1", "tag2"},
		Link:    "http://example.com",
	}

	client.SendAsync(event)
}

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test_api_key" {
			t.Errorf("Expected Authorization header to be 'Bearer test_api_key'")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be 'application/json'")
		}
		if r.Header.Get("X-Integration") != "golang" {
			t.Errorf("Expected X-Integration header to be 'golang'")
		}
		if r.Header.Get("X-Version") != "2.0.0" {
			t.Errorf("Expected X-Version header to be '2.0.0'")
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
		err = json.NewEncoder(w).Encode(map[string]string{"project": "test_project"})

		if err != nil {
			t.Errorf("Error encoding response: %v", err)
		}
	}))
	defer server.Close()

	client := ApiAlertsClientWithConfig("test_api_key", defaultConfig)

	err := client.sendToUrlWithApiKey(
		server.URL,
		"test_api_key",
		model.APIAlertsEvent{
			Channel: "test_channel",
			Message: "Test message",
		})

	if err != nil {
		t.Errorf("Error sending message: %v", err)
	}

	err = client.sendToUrlWithApiKey(
		server.URL,
		"test_api_key",
		model.APIAlertsEvent{
			Channel: "test_channel",
		})

	if err == nil || err.Error() != "message is required" {
		t.Errorf("Expected 'message is required' error, got %v", err)
	}
}
