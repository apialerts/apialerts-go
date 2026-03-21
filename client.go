package apialerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Client struct {
	apiKey             string
	integration        string
	integrationVersion string
	baseURL            string
	config             Config
	httpClient         *http.Client
}

func initializeClient(apiKey string, config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	return &Client{
		apiKey: apiKey,
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (client *Client) sendToUrlWithApiKey(url string, apiKey string, event Event) {
	// Critical checks — always log regardless of debug setting
	if apiKey == "" {
		log.Printf("x (apialerts.com) Error: api key is missing")
		return
	}
	if event.Message == "" {
		log.Printf("x (apialerts.com) Error: message is required")
		return
	}

	ch := make(chan *SendResult, 1)
	go func() {
		ch <- client.sendToUrlWithApiKeyAsync(url, apiKey, event)
	}()

	select {
	case result := <-ch:
		if !client.config.Debug {
			return
		}
		if !result.Success {
			log.Printf("x (apialerts.com) Error: %s", result.Error)
		} else {
			log.Printf("✓ (apialerts.com) Alert sent to %s (%s)", result.Workspace, result.Channel)
			for _, w := range result.Warnings {
				log.Printf("! (apialerts.com) Warning: %s", w)
			}
		}
	case <-time.After(client.config.Timeout):
		if client.config.Debug {
			log.Println("x (apialerts.com) Error: Send operation timed out")
		}
	}
}

func (client *Client) sendToUrlWithApiKeyAsync(url string, apiKey string, event Event) *SendResult {
	if apiKey == "" {
		return &SendResult{Success: false, Error: "api key is missing, use Configure() to set it"}
	}
	if event.Message == "" {
		return &SendResult{Success: false, Error: "message is required"}
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return &SendResult{Success: false, Error: "failed to serialize event"}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return &SendResult{Success: false, Error: "failed to create request"}
	}

	integration := IntegrationName
	if client.integration != "" {
		integration = client.integration
	}
	version := IntegrationVersion
	if client.integrationVersion != "" {
		version = client.integrationVersion
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Integration", integration)
	req.Header.Set("X-Version", version)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return &SendResult{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var data map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return &SendResult{Success: false, Error: "invalid response from server"}
		}
		result := &SendResult{
			Success:   true,
			Workspace: stringVal(data, "workspace"),
			Channel:   stringVal(data, "channel"),
		}
		if warnings, ok := data["warnings"].([]interface{}); ok {
			for _, w := range warnings {
				if s, ok := w.(string); ok {
					result.Warnings = append(result.Warnings, s)
				}
			}
		}
		return result
	case http.StatusBadRequest:
		return &SendResult{Success: false, Error: "bad request"}
	case http.StatusUnauthorized:
		return &SendResult{Success: false, Error: "unauthorized — check your API key"}
	case http.StatusForbidden:
		return &SendResult{Success: false, Error: "forbidden"}
	case http.StatusTooManyRequests:
		return &SendResult{Success: false, Error: "rate limit exceeded"}
	default:
		return &SendResult{Success: false, Error: fmt.Sprintf("unexpected status: %d", resp.StatusCode)}
	}
}

func stringVal(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
