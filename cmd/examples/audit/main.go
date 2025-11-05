package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MagaluCloud/mgc-sdk-go/audit"
	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func main() {
	ExampleListEvents()
	ExampleListEventTypes()
}

func ExampleListEvents() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	eventsClient := audit.New(c)

	eventsList, err := eventsClient.Events().List(context.Background(), &audit.ListEventsParams{
		Limit: helpers.IntPtr(1),
		EventFilterParams: audit.EventFilterParams{
			TypeLike: helpers.StrPtr("cloud.magalu.block-storage.snapshot.create"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d events:\n", len(eventsList.Results))
	for _, event := range eventsList.Results {
		fmt.Printf("Event: %s\n", event.ID)
		fmt.Printf("  Type: %s\n", event.Type)
		fmt.Printf("  Source: %s\n", event.Source)
		fmt.Printf("  Time: %s\n", event.Time)
		fmt.Printf("  Subject: %s\n", event.Subject)
		fmt.Printf("  AuthID: %s\n", event.AuthID)
		fmt.Printf("  Product: %s\n", event.Product)
		if event.Region != nil {
			fmt.Printf("  Region: %s\n", *event.Region)
		}
	}
}

func ExampleListEventTypes() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	eventsClient := audit.New(c)

	types, err := eventsClient.EventTypes().List(context.Background(), &audit.ListEventTypesParams{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nFound %d event types:\n", len(types.Results))
	for _, eventType := range types.Results {
		fmt.Printf("Event Type: %s\n", eventType.Type)
	}
}
