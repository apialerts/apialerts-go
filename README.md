# apialerts-go

Golang client for the [apialerts.com](https://apialerts.com/) platform

[Docs](https://apialerts.com/docs/go) • [GitHub](https://github.com/apialerts/apialerts-go)

## Installation

Add the following dependency to your GO application

```bash
go get github.com/apialerts/apialerts-go
```

### Initialize the client

The client is implemented as a singleton, ensuring that only one instance is created and used throughout the application.

```go
// Basic initialization with default config
apialerts.Configure("your-api-key")

// or initialization with custom config
customConfig := apialerts.Config {
    Timeout: 45 * time.Second,  // default is 30 seconds
    Debug:   True,              // default is false
}
apialerts.ConfigureWithConfig("test_api_key", customConfig)
}
```

### Send Events

You can send alerts by constructing the Event struct and passing it to the Send() function.

```go
event := apialerts.Event {
    Channel: "test_channel",           // optional, uses the default channel if not provided
    Message: "Test message",           // required
    Tags:    []string{"tag1", "tag2"}, // optional
    Link:    "http://example.com",     // optional
}

apialerts.Send(event)
```

The apialerts.sendAsync() methods are also available if you need to wait for a successful execution. However, the send() functions are generally always preferred.

### Send with API Key functions

You may have the need to talk to different API Alerts workspaces in your application. You can use the SendWithAPIKey() functions to send alerts to override the default apikey for that single send call

```go
apialerts.SendWithApiKey("other_api_key", event)
```

### Feedback & Support

If you have any questions or feedback, please create an issue on our GitHub repository. We are always looking to improve our service and would love to hear from you. Thanks for using API Alerts!






