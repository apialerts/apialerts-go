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

	type outcome struct {
		result *Result
		err    error
	}
	ch := make(chan outcome, 1)
	go func() {
		r, e := client.sendToUrlWithApiKeyAsync(url, apiKey, event)
		ch <- outcome{r, e}
	}()

	select {
	case o := <-ch:
		if !client.config.Debug {
			return
		}
		if o.err != nil {
			log.Printf("x (apialerts.com) Error: %s", o.err)
		} else {
			log.Printf("✓ (apialerts.com) Alert sent to %s (%s)", o.result.Workspace, o.result.Channel)
			for _, w := range o.result.Warnings {
				log.Printf("! (apialerts.com) Warning: %s", w)
			}
		}
	case <-time.After(client.config.Timeout):
		if client.config.Debug {
			log.Println("x (apialerts.com) Error: Send operation timed out")
		}
	}
}

func (client *Client) sendToUrlWithApiKeyAsync(url string, apiKey string, event Event) (*Result, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is missing, use Configure() to set it")
	}
	if event.Message == "" {
		return nil, fmt.Errorf("message is required")
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize event")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request")
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
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var data map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("invalid response from server")
		}
		result := &Result{
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
		return result, nil
	case http.StatusBadRequest:
		return nil, fmt.Errorf("bad request")
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized — check your API key")
	case http.StatusForbidden:
		return nil, fmt.Errorf("forbidden")
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("rate limit exceeded")
	default:
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
}

func stringVal(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
