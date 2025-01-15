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
	apiKey     string
	config     Config
	httpClient *http.Client
}

func initializeClient(apiKey string, config Config) *Client {
	return &Client{
		apiKey: apiKey,
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (client *Client) sendToUrlWithApiKey(url string, apiKey string, event Event) {
	errChan := make(chan error, 1)
	go func() {
		errChan <- client.sendToUrlWithApiKeyAsync(url, apiKey, event)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("x (apialerts.com) Error: %v", err)
		}
	case <-time.After(client.config.Timeout):
		log.Println("x (apialerts.com) Error: Send operation timed out")
	}
}

func (client *Client) sendToUrlWithApiKeyAsync(url string, apiKey string, event Event) error {
	log.Println("x (apialerts.com) Sending alert...")
	log.Println(apiKey)
	log.Println(event)

	if apiKey == "" {
		return errors.New("x (apialerts.com) Error: api key is missing, use SetApiKey() to set it")
	}

	if event.Message == "" {
		return errors.New("x (apialerts.com) Error: message is required")
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Integration", IntegrationName)
	req.Header.Set("X-Version", IntegrationVersion)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("x (apialerts.com) Error closing response body: %v", err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		if client.config.Debug {
			log.Printf("✓ (apialerts.com) Alert sent to %v (%v) successfully.", data["workspace"], data["channel"])
		}
		return nil
	case http.StatusBadRequest:
		return errors.New("x (apialerts.com) Error: bad request")
	case http.StatusUnauthorized:
		return errors.New("x (apialerts.com) Error: unauthorized")
	case http.StatusForbidden:
		return errors.New("x (apialerts.com) Error: forbidden")
	case http.StatusTooManyRequests:
		return errors.New("x (apialerts.com) Error: rate limit exceeded")
	default:
		return errors.New("x (apialerts.com) Error: unknown error")
	}
}
