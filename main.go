package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/apialerts/apialerts-go/model"
)

const X_VERSION = "1.0.0"
const X_INTEGRATION = "golang"
const DEFAULT_API_URL = "https://api.apialerts.com/event"

var defaultConfig = model.APIAlertsConfig{
	Logging: true,
	Timeout: 30 * time.Second,
	Debug:   false,
}

type APIAlertsClient struct {
	ApiKey string
	Config model.APIAlertsConfig
}

func ApiAlertsClientWithConfig(apiKey string, config model.APIAlertsConfig) *APIAlertsClient {
	return &APIAlertsClient{
		ApiKey: apiKey,
		Config: config,
	}
}

func ApiAlertsClient(apiKey string) *APIAlertsClient {
	return ApiAlertsClientWithConfig(apiKey, defaultConfig)
}

func (client *APIAlertsClient) SetApiKey(apiKey string) {
	client.ApiKey = apiKey
}

func (client *APIAlertsClient) sendToUrlWithApiKey(
	url string,
	apiKey string,
	event model.APIAlertsEvent,
) error {
	if apiKey == "" {
		return errors.New("api key is missing")
	}

	if event.Message == "" {
		return errors.New("message is required")
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
	req.Header.Set("X-Integration", X_INTEGRATION)
	req.Header.Set("X-Version", X_VERSION)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		if client.Config.Logging {
			log.Printf("✓ (apialerts.com) Alert sent to %v successfully.", data["project"])
		}
		return nil
	case http.StatusBadRequest:
		return errors.New("bad request")
	case http.StatusUnauthorized:
		return errors.New("unauthorized")
	case http.StatusForbidden:
		return errors.New("forbidden")
	case http.StatusTooManyRequests:
		return errors.New("rate limit exceeded")
	default:
		return errors.New("unknown error")
	}
}

func (client *APIAlertsClient) SendWithApiKey(apiKey string, event model.APIAlertsEvent) error {
	return client.sendToUrlWithApiKey(DEFAULT_API_URL, apiKey, event)
}

func (client *APIAlertsClient) SendAsyncWithApiKey(apiKey string, event model.APIAlertsEvent) {
	if client.Config.Debug {
		errChan := make(chan error, 1)
		go func() {
			errChan <- client.SendWithApiKey(apiKey, event)
		}()

		select {
		case err := <-errChan:
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		case <-time.After(30 * time.Second):
			log.Println("Send operation timed out")
		}
	} else {
		go func() {
			_ = client.SendWithApiKey(apiKey, event)
		}()
	}
}

func (client *APIAlertsClient) SendAsync(event model.APIAlertsEvent) {
	if client.Config.Debug {
		errChan := make(chan error, 1)
		go func() {
			errChan <- client.Send(event)
		}()

		select {
		case err := <-errChan:
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		case <-time.After(client.Config.Timeout):
			log.Println("Send operation timed out")
		}
	} else {
		go func() {
			_ = client.Send(event)
		}()
	}
}

func (client *APIAlertsClient) Send(event model.APIAlertsEvent) error {
	return client.SendWithApiKey(client.ApiKey, event)
}
