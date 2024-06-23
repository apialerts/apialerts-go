# apialerts-go

Golang client for [apialerts.com](https://apialerts.com/)

## Installation

```bash
go get github.com/apialerts/apialerts-go
```

## Usage

```go
package main

import (
	"log"

	"github.com/apialerts/apialerts-go"
)

func main() {
	client := apialerts.ApiAlertsClient()
	client.SetApiKey("API_KEY_GOES_HERE")
	err := client.Send("Golang Test Message", []string{"Golang is better than kotlin"}, "https://github.com/apialerts/")
	if err != nil {
		log.Println("Error:", err)
	}
}
```


