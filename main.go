package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type Client struct {
	apiKey string
}

func ApiAlertsClient() *Client {
	return &Client{
		apiKey: os.Getenv("APIALERTS_API_KEY"),
	}
}

func (client *Client) SetApiKey(apiKey string) {
	client.apiKey = apiKey
}

func (client *Client) Send(message string, tags []string, link string, debugMode ...bool) {
	isDebug := false
	if len(debugMode) > 0 {
		isDebug = debugMode[0]
	}

	if isDebug {
		errChan := make(chan error, 1)
		go func() {
			errChan <- client.send(message, tags, link)
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
			_ = client.send(message, tags, link)
		}()
	}
}

func (client *Client) send(message string, tags []string, link string) error {
	if client.apiKey == "" {
		return errors.New("api key is missing")
	}
	if message == "" {
		return errors.New("message is required")
	}

	payload := map[string]interface{}{
		"message": message,
		"tags":    tags,
		"link":    link,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := "https://api.apialerts.com/event"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Integration", "golang")
	req.Header.Set("X-Version", "1.0.0")

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
		log.Printf("✓ (apialerts.com) Alert sent to %v successfully.", data["project"])
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
