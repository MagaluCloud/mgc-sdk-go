package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/dbaas"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func main() {
	ExampleListEngines()
	ExampleListInstanceTypes()
	ExampleListInstances()
	ExampleCreateInstance()
}

func ExampleListEngines() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	engines, err := dbaasClient.Engines().List(context.Background(), dbaas.ListEngineOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d database engines:\n", len(engines))
	for _, engine := range engines {
		fmt.Printf("Engine: %s (ID: %s)\n", engine.Name, engine.ID)
		fmt.Printf("  Version: %s\n", engine.Version)
		fmt.Printf("  Status: %s\n", engine.Status)
	}
}

func ExampleListInstanceTypes() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	instanceTypes, err := dbaasClient.InstanceTypes().List(context.Background(), dbaas.ListInstanceTypeOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d instance types:\n", len(instanceTypes))
	for _, instanceType := range instanceTypes {
		fmt.Printf("Instance Type: %s (ID: %s)\n", instanceType.Name, instanceType.ID)
		fmt.Printf("  Label: %s\n", instanceType.Label)
		fmt.Printf("  VCPU: %s\n", instanceType.VCPU)
		fmt.Printf("  RAM: %s\n", instanceType.RAM)
		fmt.Printf("  Family: %s (%s)\n", instanceType.FamilyDescription, instanceType.FamilySlug)
	}
}

func ExampleListInstances() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	instances, err := dbaasClient.Instances().List(context.Background(), dbaas.ListInstanceOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d database instances:\n", len(instances))
	for _, instance := range instances {
		fmt.Printf("Instance: %s (ID: %s)\n", instance.Name, instance.ID)
		fmt.Printf("  Engine ID: %s\n", instance.EngineID)
		fmt.Printf("  Status: %s\n", instance.Status)
		fmt.Printf("  Volume Size: %d GB\n", instance.Volume.Size)
		fmt.Printf("  Volume Type: %s\n", instance.Volume.Type)
		if len(instance.Addresses) > 0 {
			fmt.Println("  Addresses:")
			for _, addr := range instance.Addresses {
				if addr.Address != nil {
					fmt.Printf("    %s (%s): %s\n", addr.Access, *addr.Type, *addr.Address)
				}
			}
		}
	}
}

func ExampleCreateInstance() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	// Create a new database instance
	instance, err := dbaasClient.Instances().Create(context.Background(), dbaas.InstanceCreateRequest{
		Name:          "example-db-instance",
		EngineID:      "your-engine-id",     // Replace with actual engine ID
		InstanceTypeID: "your-instance-type-id", // Replace with actual instance type ID
		User:          "dbadmin",
		Password:      "YourStrongPassword123!",
		Volume: dbaas.InstanceVolumeRequest{
			Size: 20, // Size in GB
			Type: dbaas.VolumeTypeCloudNVME,
		},
		BackupRetentionDays: 7,
		BackupStartAt:       "02:00", // Start backup at 2 AM
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created database instance with ID: %s\n", instance.ID)
}
