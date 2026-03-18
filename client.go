package apialerts

import (
	"bytes"
	"encoding/json"
	"errors"
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

type sendOutcome struct {
	result *Result
	err    error
}

func (client *Client) sendToUrlWithApiKey(url string, apiKey string, event Event) {
	ch := make(chan sendOutcome, 1)
	go func() {
		result, err := client.sendToUrlWithApiKeyAsync(url, apiKey, event)
		ch <- sendOutcome{result, err}
	}()

	select {
	case outcome := <-ch:
		if !client.config.Debug {
			return
		}
		if outcome.err != nil {
			log.Printf("x (apialerts.com) Error: %v", outcome.err)
		} else if outcome.result != nil {
			log.Printf("✓ (apialerts.com) Alert sent to %v (%v) successfully.", outcome.result.Workspace, outcome.result.Channel)
			for _, w := range outcome.result.Warnings {
				log.Printf("! (apialerts.com) Warning: %v", w)
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
		return nil, errors.New("api key is missing, use Configure() to set it")
	}

	if event.Message == "" {
		return nil, errors.New("message is required")
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	integration := IntegrationName
	if client.integration != "" {
		integration = client.integration
	}
	version := IntegrationVersion
	if client.integrationVersion != "" {
		version = client.integrationVersion
	}
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
			return nil, err
		}
		result := &Result{
			Workspace: stringVal(data, "workspace"),
			Channel:   stringVal(data, "channel"),
		}
		if warnings, ok := data["errors"].([]interface{}); ok {
			for _, w := range warnings {
				if s, ok := w.(string); ok {
					result.Warnings = append(result.Warnings, s)
				}
			}
		}
		return result, nil
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusUnauthorized:
		return nil, errors.New("unauthorized — check your API key")
	case http.StatusForbidden:
		return nil, errors.New("forbidden")
	case http.StatusTooManyRequests:
		return nil, errors.New("rate limit exceeded")
	default:
		return nil, errors.New("unexpected error")
	}
}

func stringVal(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
