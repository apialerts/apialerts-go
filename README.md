# apialerts-go

Golang client for the [apialerts.com](https://apialerts.com/) platform

[Docs](https://apialerts.com/docs/go) • [GitHub](https://github.com/apialerts/apialerts-go)

## Installation

Add the following dependency to your GO application

```bash
go get github.com/apialerts/apialerts-go
```

## Usage

### Event structure
```go
event := model.APIAlertsEvent {
    Channel: "test_channel",           // optional, uses the default channel if not provided
    Message: "Test message",           // required
    Tags   : []string{"tag1", "tag2"}, // optional
    Link   : "http://example.com",     // optional
}
```


### Initialize the client 
```go
// Custom config can be passed to the client
customConfig := model.APIAlertsConfig {
    Logging: true,
    Timeout: 30 * time.Second,
    Debug:   false,
}

client := ApiAlertsClientWithConfig("test_api_key", customConfig)
// or
client := ApiAlertsClient("test_api_key")
}
```

### Send an event

```go
event := model.APIAlertsEvent {
    Channel: "test_channel",           // optional, uses the default channel if not provided
    Message: "Test message",
    Tags:    []string{"tag1", "tag2"}, // optional
    Link:    "http://example.com",     // optional
}

err := client.Send(event)
if err != nil {
    log.Printf("Error sending message: %v", err)
}
```

### Send an event with an alternate api key

```go
client.SendWithApiKey("other project api key", event)
```


### Send message asynchronously

Async methods of the Send() function are available in situations (like AWS Lambda) where you need to wait for the event execution. However, using the Send() functions are generally always preferred.

```go
client.SendAsync(event)
```


### Send message asynchronously with custom config

Custom config can also be applied with a Send event

```go
func main() {
    //... client initialization
    client.SendAsyncWithApiKey("other project api key", event)
}
```





