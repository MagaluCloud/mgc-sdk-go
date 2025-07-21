# Availabilityzones

Example usage of the `availabilityzones` module.

**File:** `cmd/examples/availabilityzones/main.go`

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MagaluCloud/mgc-sdk-go/availabilityzones"
	"github.com/MagaluCloud/mgc-sdk-go/client"
)

func main() {
	ExampleListAvailabilityZones()
}

func ExampleListAvailabilityZones() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	azClient := availabilityzones.New(c)

	response, err := azClient.AvailabilityZones().List(context.Background(), availabilityzones.ListOptions{
		ShowBlocked: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d regions:\n", len(response))
	for _, region := range response {
		fmt.Printf("Region: %s\n", region.ID)
		fmt.Printf("  Availability Zones:\n")
		for _, az := range region.AvailabilityZones {
			fmt.Printf("    - ID: %s (Block Type: %s)\n", az.ID, az.BlockType)
		}
	}
}
```
