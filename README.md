# API Alerts • Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/apialerts/apialerts-go.svg)](https://pkg.go.dev/github.com/apialerts/apialerts-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

[GitHub](https://github.com/apialerts/apialerts-go) • [API Alerts](https://apialerts.com)

Effortless project notifications. Send once, deliver everywhere.

## Installation

```bash
go get github.com/apialerts/apialerts-go@v1.2.0
```

## Quick Start

```go
import apialerts "github.com/apialerts/apialerts-go"

apialerts.Configure("your-api-key")
apialerts.Send(apialerts.Event{Message: "Deploy complete"})
```

## Setup

The client is a singleton. Configure it once at startup; later calls to `Configure` are ignored.

```go
apialerts.Configure("your-api-key")

// Enable debug logging to stderr
apialerts.SetDebug(true)
```

## Send Events

### Fire and forget

```go
event := apialerts.Event{
    Message: "Deploy complete",          // required
    Channel: "deployments",              // optional, uses default channel if not set
    Event:   "deploy.success",           // optional, event name for routing
    Title:   "Production Deploy",        // optional
    Tags:    []string{"deploy", "prod"}, // optional
    Link:    "https://example.com",      // optional
    Data:    map[string]any{     // optional, arbitrary JSON data
        "version": "1.4.2",
        "region":  "us-east-1",
    },
}

apialerts.Send(event)
```

`Send` returns immediately and delivers in a background goroutine. In short-lived programs (CLI tools, CI scripts) that exit right after sending, use `SendAsync` so the process waits for delivery.

### Wait for response

`SendAsync` blocks until the request completes and returns the result or an error.

```go
result, err := apialerts.SendAsync(event)
if err != nil {
    log.Println("Failed to send:", err)
    return
}

fmt.Printf("Sent to %s (%s)\n", result.Workspace, result.Channel)

for _, warning := range result.Warnings {
    fmt.Println("Warning:", warning)
}
```

The `Result` contains:

| Field | Description |
|-------|-------------|
| `Workspace` | Name of the workspace the event was delivered to |
| `Channel` | Name of the channel the event was delivered to |
| `Warnings` | Any non-fatal warnings returned by the API (e.g. unknown fields) |

## Event Fields

| Field     | Type             | Required | Description                      |
|-----------|------------------|----------|----------------------------------|
| `Message` | `string`         | Yes      | Main notification message        |
| `Channel` | `string`         | No       | Target channel name              |
| `Event`   | `string`         | No       | Event key (e.g. `ci.deploy`)     |
| `Title`   | `string`         | No       | Short title                      |
| `Tags`    | `[]string`       | No       | Categorisation tags              |
| `Link`    | `string`         | No       | URL attached to the notification |
| `Data`    | `map[string]any` | No       | Arbitrary key-value metadata     |

Null/empty fields are omitted from the JSON payload automatically.

## Send to Multiple Workspaces

Use `SendWithKey` or `SendWithKeyAsync` to override the API key for a single call.

```go
apialerts.SendWithKey("other-workspace-api-key", event)

result, err := apialerts.SendWithKeyAsync("other-workspace-api-key", event)
```

## Links

- [Documentation](https://apialerts.com/docs)
- [Sign up](https://apialerts.com)
- [GitHub Issues](https://github.com/apialerts/apialerts-go/issues)
