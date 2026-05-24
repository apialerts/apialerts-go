package apialerts

import (
	"sync"
	"testing"
)

func resetInstance() {
	instance = nil
	once = sync.Once{}
}

func TestConfigure(t *testing.T) {
	resetInstance()
	Configure("test_api_key")

	if instance.apiKey != "test_api_key" {
		t.Errorf("Expected API key to be 'test_api_key', got '%s'", instance.apiKey)
	}
	if instance.debug != false {
		t.Errorf("Expected debug to be false, got %v", instance.debug)
	}
}

func TestSetDebug(t *testing.T) {
	resetInstance()
	Configure("test_api_key")
	SetDebug(true)

	if !instance.debug {
		t.Errorf("Expected debug to be true, got %v", instance.debug)
	}
}

func TestDefaultTimeout(t *testing.T) {
	resetInstance()
	Configure("test_api_key")

	if instance.httpClient.Timeout != DefaultTimeout {
		t.Errorf("Expected default timeout to be %v, got %v", DefaultTimeout, instance.httpClient.Timeout)
	}
}

func TestNotInitialized(t *testing.T) {
	resetInstance()

	_, err := SendAsync(Event{Message: "test"})
	if err == nil {
		t.Error("Expected error when not initialized")
	}
}
