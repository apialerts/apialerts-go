package main

import (
	"bytes"
	"encoding/json"
	"os"
	"errors"
	"net/http"
	"fmt"
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

func (client *Client) Send(message string, tags []string, link string) error {
	if client.apiKey == "" {
		return errors.New("Api Key is Missing")
	}
	if message == "" {
		return errors.New("Message is required")
	}
	
	payload := map[string]interface{}{
		"message": message,
		"tags": tags,
		"link": link,
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
	
	req.Header.Set("Authorization", "Bearer " + client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Integration", "golang")
	req.Header.Set("X-Version", "1.0.0")
	
	httpClient := &http.Client{}
	
	resp, err :=  httpClient.Do(req)
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
			fmt.Printf("✓ (apialerts.com) Alert sent to %v successfully.", data["project"])
			return nil
		case http.StatusBadRequest:
			return errors.New("Bad Request")
		case http.StatusUnauthorized:
			return errors.New("Unauthorized")
		case http.StatusForbidden:
			return errors.New("Forbidden")
		case http.StatusTooManyRequests:
			return errors.New("Rate Limit Exceeded")
		default:
			return errors.New("Unknown Error")
	}
	
}

// THIS IS SOME EXAMPLE CODE
/*
func main() {
	client := ApiAlertsClient()
	client.SetApiKey("API_KEY_GOES_HERE")
	err := client.Send("Golang Test Message", []string{"Golang is better than kotlin"}, "https://github.com/apialerts/")
	if err != nil {
		fmt.Println("Error:", err)
	}
}
*/