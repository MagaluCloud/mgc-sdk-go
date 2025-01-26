package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
	"github.com/MagaluCloud/mgc-sdk-go/sdk/virtualmachine"
)

func main() {
	fmt.Println("Running all examples...")

	fmt.Println("\n1. List Instances Example:")
	ExampleListInstances()

	fmt.Println("\n2. Create Instance Example:")
	id := ExampleCreateInstance()

	fmt.Println("\n3. Manage Instance Example:")
	// Wait a bit to ensure the instance is created
	time.Sleep(10 * time.Second)
	ExampleManageInstance(id)

	fmt.Println("\n4. Delete Instance Example:")
	// Wait a bit to ensure all management operations are complete
	time.Sleep(5 * time.Second)
	ExampleDeleteInstance(id)

	fmt.Println("\nAll examples completed successfully!")
}

func ExampleListInstances() {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.New(apiToken)
	vmClient := virtualmachine.New(c)

	// List instances with pagination and sorting
	instances, err := vmClient.Instances().List(context.Background(), virtualmachine.ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Expand: []string{"machine-type", "image"},
	})

	if err != nil {
		log.Fatal(err)
	}

	// Print instance details
	for _, instance := range instances {
		fmt.Printf("Instance: %s (ID: %s)\n", instance.Name, instance.ID)
		fmt.Printf("  Machine Type: %s\n", *instance.MachineType.Name)
		fmt.Printf("  Image: %s\n", *instance.Image.Name)
		fmt.Printf("  Status: %s\n", instance.Status)
		fmt.Printf("  State: %s\n", instance.State)
		fmt.Printf("  Created At: %s\n", instance.CreatedAt)
	}
}

func ExampleCreateInstance() string {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.New(apiToken)
	vmClient := virtualmachine.New(c)

	// Create a new instance
	createReq := virtualmachine.CreateRequest{
		Name: "my-test-vm",
		MachineType: virtualmachine.IDOrName{
			Name: helpers.StrPtr("BV1-1-40"),
		},
		Image: virtualmachine.IDOrName{
			Name: helpers.StrPtr("cloud-ubuntu-24.04 LTS"),
		},
		Network: &virtualmachine.CreateParametersNetwork{
			AssociatePublicIp: helpers.BoolPtr(false),
		},
		SshKeyName: helpers.StrPtr("publio"),
	}

	id, err := vmClient.Instances().Create(context.Background(), createReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created instance with ID: %s\n", id)

	return id
}

func ExampleManageInstance(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.New(apiToken)
	vmClient := virtualmachine.New(c)
	ctx := context.Background()

	// Get instance details
	instance, err := vmClient.Instances().Get(ctx, id, []string{"network"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Instance: %s (ID: %s)\n", instance.Name, instance.ID)

	// Rename the instance
	if err := vmClient.Instances().Rename(ctx, instance.ID, "new-name"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Instance renamed successfully")

	// Change machine type
	retypeReq := virtualmachine.RetypeRequest{
		MachineType: virtualmachine.IDOrName{
			Name: helpers.StrPtr("BV2-2-20"),
		},
	}
	if err := vmClient.Instances().Retype(ctx, instance.ID, retypeReq); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Instance machine type changed successfully")
}

func ExampleDeleteInstance(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.New(apiToken)
	vmClient := virtualmachine.New(c)

	// Delete instance and its public IP
	if err := vmClient.Instances().Delete(context.Background(), id, true); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Instance deleted successfully")
}
