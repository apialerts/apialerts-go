# apialerts-go

Golang client for [apialerts.com](https://apialerts.com/)

## Installation

```bash
go get github.com/apialerts/apialerts-go
```

## Usage

### Event structure
```go
	event := model.APIAlertsEvent{
		Channel: "test_channel",           // optional, if not set events will be sent to the default channel
		Message: "Test message",           // required
		Tags   : []string{"tag1", "tag2"}, // optional
		Link   : "http://example.com",     // optional
	}
```


### Initialize client 
```go
func main() {
    // Custom config can be passed to the client
    customConfig := model.APIAlertsConfig{
        Logging: true,
        Timeout: 30 * time.Second,
        Debug:   false,
    }

	client := ApiAlertsClientWithConfig("test_api_key", customConfig)
	// or
	client := ApiAlertsClient("test_api_key")

	event := model.APIAlertsEvent{
		Channel: "test_channel",
		Message: "Test message",
		Tags:    []string{"tag1", "tag2"},
		Link:    "http://example.com",
	}

	err := client.Send(event)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
```

### Send message with custom api key
```go
func main() {
    //... client initialization
    client.SendWithApiKey("other project api key", event)
}
```


### Send message asynchronously
```go
func main() {
    //... client initialization
	client.SendAsync(event)
}
```


### Send message asynchronously with custom config
```go
func main() {
    //... client initialization
	client.SendAsyncWithApiKey("other project api key", event)
}
```





