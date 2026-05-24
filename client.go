package apialerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Client struct {
	apiKey             string
	integration        string
	integrationVersion string
	baseURL            string
	debug              bool
	httpClient         *http.Client
}

func initializeClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		debug:  false,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

func (client *Client) sendToUrlWithApiKey(url string, apiKey string, event Event) {
	// Critical checks: always log regardless of debug setting
	if apiKey == "" {
		log.Print("x (apialerts.com) Error: api key is missing")
		return
	}
	if event.Message == "" {
		log.Print("x (apialerts.com) Error: message is required")
		return
	}

	// This already runs in its own goroutine; httpClient.Timeout bounds the request.
	result, err := client.sendToUrlWithApiKeyAsync(url, apiKey, event)
	if !client.debug {
		return
	}
	if err != nil {
		log.Printf("x (apialerts.com) Error: %s", err)
		return
	}
	log.Printf("✓ (apialerts.com) Alert sent to %s (%s)", result.Workspace, result.Channel)
	for _, w := range result.Warnings {
		log.Printf("! (apialerts.com) Warning: %s", w)
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
		return nil, fmt.Errorf("unauthorized, check your API key")
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
