package apialerts

import (
	"sync"
	"testing"
	"time"
)

func resetInstance() {
	instance = nil
	once = sync.Once{}
}

func TestApiAlertsClient(t *testing.T) {
	resetInstance()
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
		t.Errorf("Expected Debug to be true, got %v", instance.config.Debug)
	}
}

func TestSetApiKey(t *testing.T) {
	resetInstance()
	Configure("test_api_key")
	SetApiKey("new_api_key")

	if instance.apiKey != "new_api_key" {
		t.Errorf("Expected API key to be 'new_api_key', got '%s'", instance.apiKey)
	}
}

func TestDefaultTimeout(t *testing.T) {
	resetInstance()
	Configure("test_api_key")

	if instance.config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30 seconds, got %v", instance.config.Timeout)
	}
}

func TestNotInitialized(t *testing.T) {
	resetInstance()

	_, err := SendAsync(Event{Message: "test"})
	if err == nil {
		t.Error("Expected error when not initialized, got nil")
	}
}
