package main

import (
	"flag"
	"fmt"
	"github.com/apialerts/apialerts-go"
	"os"
)

func main() {
	build, release, publish := parseFlags()

	apiKey := getApiKey()
	if apiKey == "" {
		fmt.Println("Error: APIALERTS_API_KEY environment variable is not set")
		return
	}

	apialerts.Configure(apiKey)

	if !*build && !*release && !*publish {
		fmt.Println("Usage: go run github.go --build|--release|--publish")
		return
	}

	event := createEvent(*build, *release, *publish)
	apialerts.Configure(apiKey)

	if err := apialerts.SendAsync(event); err != nil {
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

func createEvent(build, release, publish bool) apialerts.Event {
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

	return apialerts.Event{
		Channel: eventChannel,
		Message: eventMessage,
		Tags:    eventTags,
		Link:    eventLink,
	}
}
