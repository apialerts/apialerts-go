package apialerts

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	instance *Client
	once     sync.Once
)

type Config struct {
	Timeout time.Duration // Timeout specifies the duration to wait before timing out a request.
	Debug   bool          // Debug enables or disables debug logging.
}

var defaultConfig = Config{
	Timeout: 30 * time.Second,
	Debug:   false,
}

func setupClient(apiKey string, config Config) {
	once.Do(func() {
		instance = initializeClient(apiKey, config)
	})
}

// Configure initializes the client with the default configuration using the provided API key.
func Configure(apiKey string) {
	setupClient(apiKey, defaultConfig)
}

// ConfigureWithConfig initializes the client with a custom configuration using the provided API key.
func ConfigureWithConfig(apiKey string, config Config) {
	setupClient(apiKey, config)
}

// SetApiKey sets a new API key for the client instance.
func SetApiKey(apiKey string) {
	if instance == nil {
		return
	}
	instance.apiKey = apiKey
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
// Critical errors (missing key, not configured) are always logged.
// Other errors are only logged when debug is enabled.
func Send(event Event) {
	if instance == nil {
		log.Printf("x (apialerts.com) Error: client not initialized — call Configure() first")
		return
	}
	go instance.sendToUrlWithApiKey(instance.resolveURL(), instance.apiKey, event)
}

// SendAsync sends an event using the default API key.
// Returns an error if the send fails — check err before using result.
func SendAsync(event Event) (*Result, error) {
	if instance == nil {
		return nil, fmt.Errorf("client not initialized — call Configure() first")
	}
	return instance.sendToUrlWithApiKeyAsync(instance.resolveURL(), instance.apiKey, event)
}

// SendWithApiKey sends an event fire-and-forget using the provided API key.
func SendWithApiKey(apiKey string, event Event) {
	if instance == nil {
		log.Printf("x (apialerts.com) Error: client not initialized — call Configure() first")
		return
	}
	go instance.sendToUrlWithApiKey(instance.resolveURL(), apiKey, event)
}

// SendWithApiKeyAsync sends an event using the provided API key.
// Returns an error if the send fails — check err before using result.
func SendWithApiKeyAsync(apiKey string, event Event) (*Result, error) {
	if instance == nil {
		return nil, fmt.Errorf("client not initialized — call Configure() first")
	}
	return instance.sendToUrlWithApiKeyAsync(instance.resolveURL(), apiKey, event)
}
