package apialerts

import (
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
