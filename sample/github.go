package main

import (
	"flag"
	"fmt"
	"github.com/apialerts/apialerts-go"
	"os"
)

func main() {
	build := flag.Bool("build", false, "Send build notification")
	release := flag.Bool("release", false, "Send release notification")
	publish := flag.Bool("publish", false, "Send publish notification")
	flag.Parse()

	apiKey := os.Getenv("APIALERTS_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: APIALERTS_API_KEY environment variable is not set")
		return
	}

	apialerts.Configure(apiKey)

	link := "https://github.com/apialerts/apialerts-go/actions"

	switch {
	// SDK CI notifications — called from build-release.yml / publish.yml
	case *build:
		result := apialerts.SendAsync(apialerts.Event{
			Channel: "developer",
			Event:   "ci.build",
			Title:   "Build Passed",
			Message: "Go - PR build success",
			Tags:    []string{"CI/CD", "Go", "Build"},
			Link:    link,
		})
		if !result.Success {
			fmt.Println("Error:", result.Error)
			return
		}
		fmt.Printf("✓ Sent to %s (%s)\n", result.Workspace, result.Channel)

	case *release:
		result := apialerts.SendAsync(apialerts.Event{
			Channel: "developer",
			Event:   "ci.release",
			Title:   "Release Build Passed",
			Message: "Go - Build for publish success",
			Tags:    []string{"CI/CD", "Go", "Build"},
			Link:    link,
		})
		if !result.Success {
			fmt.Println("Error:", result.Error)
			return
		}
		fmt.Printf("✓ Sent to %s (%s)\n", result.Workspace, result.Channel)

	case *publish:
		result := apialerts.SendAsync(apialerts.Event{
			Channel: "releases",
			Event:   "ci.publish",
			Title:   "Published",
			Message: "Go - GitHub publish success",
			Tags:    []string{"CI/CD", "Go", "Deploy"},
			Link:    link,
		})
		if !result.Success {
			fmt.Println("Error:", result.Error)
			return
		}
		fmt.Printf("✓ Sent to %s (%s)\n", result.Workspace, result.Channel)

	// Integration test — called from apialerts-integration-tests with no args
	default:
		r1 := apialerts.SendAsync(apialerts.Event{Message: "Go SDK - minimal"})
		if !r1.Success {
			fmt.Println("Error (minimal):", r1.Error)
			return
		}
		fmt.Printf("✓ sent to %s (%s)\n", r1.Workspace, r1.Channel)

		r2 := apialerts.SendAsync(apialerts.Event{
			Message: "Go SDK - full",
			Channel: "developer",
			Event:   "sdk.test",
			Title:   "Integration Test",
			Tags:    []string{"CI/CD", "Go"},
			Link:    link,
		})
		if !r2.Success {
			fmt.Println("Error (full):", r2.Error)
			return
		}
		fmt.Printf("✓ sent to %s (%s)\n", r2.Workspace, r2.Channel)
	}
}
