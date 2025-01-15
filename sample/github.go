package main

import (
	"flag"
	"fmt"
	"github.com/apialerts/apialerts-go"
	"github.com/apialerts/apialerts-go/model"
	"os"
)

func main() {
	build, release, publish := parseFlags()

	if !*build && !*release && !*publish {
		fmt.Println("Usage: go run github.go --build|--release|--publish")
		return
	}

	apiKey := getApiKey()
	if apiKey == "" {
		fmt.Println("Error: APIALERTS_API_KEY environment variable is not set")
		return
	}

	event := createEvent(*build, *release, *publish)
	client := apialerts.ApiAlertsClient(apiKey)

	if err := client.Send(event); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Alert sent successfully.")
	}
}

func parseFlags() (build, release, publish *bool) {
	build = flag.Bool("build", false, "Build the project")
	release = flag.Bool("release", false, "Release the project")
	publish = flag.Bool("publish", false, "Publish the project")
	flag.Parse()
	return
}

func getApiKey() string {
	return os.Getenv("APIALERTS_API_KEY")
}

func createEvent(build, release, publish bool) model.APIAlertsEvent {
	eventChannel := "developer"
	eventMessage := "apialerts-go"
	var eventTags []string
	eventLink := "https://github.com/apialerts/apialerts-go/actions"

	if build {
		eventMessage = "Go - PR build success"
		eventTags = []string{"CI/CD", "Go", "Build"}
	} else if release {
		eventMessage = "Go -Build for publish success"
		eventTags = []string{"CI/CD", "Go", "Build"}
	} else if publish {
		eventMessage = "Go - GitHub publish success"
		eventTags = []string{"CI/CD", "Go", "Deploy"}
	}

	return model.APIAlertsEvent{
		Channel: eventChannel,
		Message: eventMessage,
		Tags:    eventTags,
		Link:    eventLink,
	}
}
