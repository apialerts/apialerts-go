package apialerts

import (
	"fmt"
	"log"
	"sync"
)

var (
	instance *Client
	once     sync.Once
)

func setupClient(apiKey string) {
	once.Do(func() {
		instance = initializeClient(apiKey)
	})
}

// Configure initializes the client with the provided API key.
// Subsequent calls are no-ops: the first call wins.
func Configure(apiKey string) {
	setupClient(apiKey)
}

// SetDebug enables or disables debug logging to stderr.
// When enabled, successful sends and errors are logged.
func SetDebug(debug bool) {
	if instance == nil {
		return
	}
	instance.debug = debug
}

// SetOverrides configures internal client settings used by first-party integrations
// such as the official apialerts CLI. This is not intended for general SDK use.
func SetOverrides(integration, version, baseURL string) {
	if instance == nil {
		return
	}
	instance.integration = integration
	instance.integrationVersion = version
	instance.baseURL = baseURL
}

func (client *Client) resolveURL() string {
	if client.baseURL != "" {
		return client.baseURL
	}
	return ApiUrl
}

// Send sends an event fire-and-forget using the default API key.
// It returns immediately and delivers in a background goroutine. In
// short-lived programs (CLI tools, CI scripts) that exit right after
// sending, use SendAsync instead so the process waits for delivery.
// Critical errors (missing key, not configured) are always logged.
// Other errors are only logged when debug is enabled.
func Send(event Event) {
	if instance == nil {
		log.Print("x (apialerts.com) Error: client not initialized, call Configure() first")
		return
	}
	go instance.sendToUrlWithApiKey(instance.resolveURL(), instance.apiKey, event)
}

// SendAsync sends an event using the default API key.
// Returns an error if the send fails. Check err before using result.
func SendAsync(event Event) (*Result, error) {
	if instance == nil {
		return nil, fmt.Errorf("client not initialized, call Configure() first")
	}
	return instance.sendToUrlWithApiKeyAsync(instance.resolveURL(), instance.apiKey, event)
}

// SendWithKey sends an event fire-and-forget using the provided API key.
// Like Send, delivery happens in a background goroutine; use
// SendWithKeyAsync in short-lived programs that exit right after sending.
func SendWithKey(apiKey string, event Event) {
	if instance == nil {
		log.Print("x (apialerts.com) Error: client not initialized, call Configure() first")
		return
	}
	go instance.sendToUrlWithApiKey(instance.resolveURL(), apiKey, event)
}

// SendWithKeyAsync sends an event using the provided API key.
// Returns an error if the send fails. Check err before using result.
func SendWithKeyAsync(apiKey string, event Event) (*Result, error) {
	if instance == nil {
		return nil, fmt.Errorf("client not initialized, call Configure() first")
	}
	return instance.sendToUrlWithApiKeyAsync(instance.resolveURL(), apiKey, event)
}
