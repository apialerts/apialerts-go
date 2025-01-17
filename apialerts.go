package apialerts

import (
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

// GetInstance returns the singleton instance of the Client.
func GetInstance() *Client {
	return instance
}

// SetApiKey sets a new API key for the client instance.
func SetApiKey(apiKey string) {
	instance.apiKey = apiKey
}

// Send sends an event asynchronously using the default API key.
func Send(event Event) {
	go instance.sendToUrlWithApiKey(ApiUrl, instance.apiKey, event)
}

// SendAsync sends an event asynchronously using the default API key, waits for the response, and returns an error if any.
func SendAsync(event Event) error {
	return instance.sendToUrlWithApiKeyAsync(ApiUrl, instance.apiKey, event)
}

// SendWithApiKey sends an event asynchronously using the provided API key.
func SendWithApiKey(apiKey string, event Event) {
	instance.sendToUrlWithApiKey(ApiUrl, apiKey, event)
}

// SendWithApiKeyAsync sends an event asynchronously using the provided API key, waits for the response, and returns an error if any.
func SendWithApiKeyAsync(apiKey string, event Event) error {
	return instance.sendToUrlWithApiKeyAsync(ApiUrl, apiKey, event)
}
