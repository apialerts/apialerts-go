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
	Timeout time.Duration
	Debug   bool
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

func Configure(apiKey string) {
	setupClient(apiKey, defaultConfig)
}

func ConfigureWithConfig(apiKey string, config Config) {
	setupClient(apiKey, config)
}

func SetApiKey(apiKey string) {
	instance.apiKey = apiKey
}

func Send(event Event) {
	go sendToUrlWithApiKey(ApiUrl, instance.apiKey, event)
}

func SendAsync(event Event) error {
	return sendToUrlWithApiKeyAsync(ApiUrl, instance.apiKey, event)
}

func SendWithApiKey(apiKey string, event Event) {
	sendToUrlWithApiKey(ApiUrl, apiKey, event)
}

func SendWithApiKeyAsync(apiKey string, event Event) error {
	return sendToUrlWithApiKeyAsync(ApiUrl, apiKey, event)
}

func sendToUrlWithApiKey(url string, apiKey string, event Event) {
	instance.sendToUrlWithApiKey(url, apiKey, event)
}

func sendToUrlWithApiKeyAsync(url string, apiKey string, event Event) error {
	return instance.sendToUrlWithApiKeyAsync(url, apiKey, event)
}
