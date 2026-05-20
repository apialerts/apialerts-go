package main

import (
	"flag"
	"fmt"
	"github.com/apialerts/apialerts-go"
	"os"
)

func main() {
	build            := flag.Bool("build",            false, "Send build notification")
	release          := flag.Bool("release",          false, "Send release notification")
	publish          := flag.Bool("publish",          false, "Send publish notification")
	integrationTests := flag.Bool("integration-tests", false, "Run integration tests")
	channel          := flag.String("channel",        "testing", "Channel for integration test sends")
	flag.Parse()

	apiKey := os.Getenv("APIALERTS_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: APIALERTS_API_KEY environment variable is not set")
		os.Exit(1)
	}

	apialerts.Configure(apiKey)

	link := "https://github.com/apialerts/apialerts-go/actions"

	switch {
	// SDK CI notifications — called from build-release.yml / publish.yml
	case *build:
		result, err := apialerts.SendAsync(apialerts.Event{
			Channel: "developer",
			Event:   "ci.build",
			Title:   "Build Passed",
			Message: "Go - PR build success",
			Tags:    []string{"CI/CD", "Go", "Build"},
			Link:    link,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Sent to %s (%s)\n", result.Workspace, result.Channel)

	case *release:
		result, err := apialerts.SendAsync(apialerts.Event{
			Channel: "developer",
			Event:   "ci.release",
			Title:   "Release Build Passed",
			Message: "Go - Build for publish success",
			Tags:    []string{"CI/CD", "Go", "Build"},
			Link:    link,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Sent to %s (%s)\n", result.Workspace, result.Channel)

	case *publish:
		result, err := apialerts.SendAsync(apialerts.Event{
			Channel: "releases",
			Event:   "ci.publish",
			Title:   "Published",
			Message: "Go - GitHub publish success",
			Tags:    []string{"CI/CD", "Go", "Deploy"},
			Link:    link,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Sent to %s (%s)\n", result.Workspace, result.Channel)

	case *integrationTests:
		r1, err := apialerts.SendAsync(apialerts.Event{Message: "Go SDK - minimal", Channel: *channel})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error (minimal):", err)
			os.Exit(1)
		}
		fmt.Printf("✓ sent to %s (%s)\n", r1.Workspace, r1.Channel)

		r2, err := apialerts.SendAsync(apialerts.Event{
			Message: "Go SDK - full",
			Channel: *channel,
			Event:   "sdk.test",
			Title:   "Integration Test",
			Tags:    []string{"CI/CD", "Go"},
			Link:    link,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error (full):", err)
			os.Exit(1)
		}
		fmt.Printf("✓ sent to %s (%s)\n", r2.Workspace, r2.Channel)

	default:
		fmt.Fprintln(os.Stderr, "Error: pass --build, --release, --publish, or --integration-tests")
		os.Exit(1)
	}
}
