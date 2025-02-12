package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func main() {
	ExampleListMachineTypes()
	ExampleListImages()
	ExampleListInstances()
	id := ExampleCreateInstance()
	ExampleGetInstance(id)
	ExampleRenameAndRetypeInstance(id)
	ExampleDeleteInstance(id)
}

func ExampleRenameAndRetypeInstance(id string) {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	computeClient := compute.New(c)
	ctx := context.Background()
	// Rename the instance
	if err := computeClient.Instances().Rename(ctx, id, "new-name"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Instance renamed successfully")

	// Change machine type
	retypeReq := compute.RetypeRequest{
		MachineType: compute.IDOrName{
			Name: helpers.StrPtr("BV2-2-20"),
		},
	}
	if err := computeClient.Instances().Retype(ctx, id, retypeReq); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Instance machine type changed successfully")
}

func ExampleListInstances() {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	computeClient := compute.New(c)

	// List instances with pagination and sorting
	instances, err := computeClient.Instances().List(context.Background(), compute.ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Expand: []string{compute.InstanceMachineTypeExpand, compute.InstanceImageExpand, compute.InstanceNetworkExpand},
	})

	if err != nil {
		log.Fatal(err)
	}

	// Print instance details
	for _, instance := range instances {
		fmt.Printf("Instance: %s (ID: %s)\n", instance.Name, instance.ID)
		fmt.Printf("  Machine Type: %s\n", instance.MachineType.Name)
		fmt.Printf("  Image: %s\n", instance.Image.Name)
		fmt.Printf("  Status: %s\n", instance.Status)
		fmt.Printf("  State: %s\n", instance.State)
		fmt.Printf("  Created At: %s\n", instance.CreatedAt)
		fmt.Printf("  Updated At: %s\n", instance.UpdatedAt)
		fmt.Printf("  VPC ID: %s\n", *instance.Network.Vpc.ID)
		fmt.Printf("  VPC Name: %s\n", *instance.Network.Vpc.Name)
		for _, ni := range instance.Network.Interfaces {
			fmt.Println("  Interface ID: ", ni.ID)
			fmt.Println("  Interface Name: ", ni.Name)
			fmt.Println("  Interface IPv4: ", ni.AssociatedPublicIpv4)
			fmt.Println("  Interface IPv6: ", ni.IpAddresses.PublicIpv6)
			fmt.Println("  Interface Local IPv4: ", ni.IpAddresses.PrivateIpv4)
			fmt.Println("Is Primary: ", ni.Primary)
			for _, sg := range ni.SecurityGroups {
				fmt.Println("  Security Group ID: ", sg)
			}
			fmt.Println("--------")
		}
	}
}

func ExampleCreateInstance() string {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	computeClient := compute.New(c)

	// Create a new instance
	createReq := compute.CreateRequest{
		Name: "my-test-vm",
		MachineType: compute.IDOrName{
			Name: helpers.StrPtr("BV1-1-40"),
		},
		Image: compute.IDOrName{
			Name: helpers.StrPtr("cloud-ubuntu-24.04 LTS"),
		},
		Network: &compute.CreateParametersNetwork{
			AssociatePublicIp: helpers.BoolPtr(false),
		},
		SshKeyName: helpers.StrPtr("publio"),
	}

	id, err := computeClient.Instances().Create(context.Background(), createReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created instance with ID: %s\n", id)

	return id
}

func ExampleGetInstance(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	computeClient := compute.New(c)
	ctx := context.Background()

	// Get instance details
	instance, err := computeClient.Instances().Get(ctx, id, []string{compute.InstanceNetworkExpand})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Instance: %s (ID: %s)\n", instance.Name, instance.ID)
	fmt.Printf("  Machine Type: %s\n", instance.MachineType.Name)
	fmt.Printf("  Image: %s\n", instance.Image.Name)
	fmt.Printf("  Status: %s\n", instance.Status)
	fmt.Printf("  State: %s\n", instance.State)
	fmt.Printf("  Created At: %s\n", instance.CreatedAt)
	fmt.Printf("  Updated At: %s\n", instance.UpdatedAt)
	fmt.Printf("  VPC ID: %s\n", *instance.Network.Vpc.ID)
	fmt.Printf("  VPC Name: %s\n", *instance.Network.Vpc.Name)
	for _, ni := range instance.Network.Interfaces {
		fmt.Println("  Interface ID: ", ni.ID)
		fmt.Println("  Interface Name: ", ni.Name)
		fmt.Println("  Interface IPv4: ", ni.AssociatedPublicIpv4)
		fmt.Println("  Interface IPv6: ", ni.IpAddresses.PublicIpv6)
		fmt.Println("  Interface Local IPv4: ", ni.IpAddresses.PrivateIpv4)
		fmt.Println("Is Primary: ", ni.Primary)
		for _, sg := range ni.SecurityGroups {
			fmt.Println("  Security Group ID: ", sg)
		}
	}
}

func ExampleDeleteInstance(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	computeClient := compute.New(c)

	// Delete instance and its public IP
	if err := computeClient.Instances().Delete(context.Background(), id, true); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Instance deleted successfully")
}

func ExampleListMachineTypes() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}

	c := client.NewMgcClient(apiToken)
	computeClient := compute.New(c)

	// List machine types
	machineTypes, err := computeClient.MachineTypes().List(context.Background(), compute.MachineTypeListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Print machine type details
	for _, mt := range machineTypes {
		fmt.Printf("Machine Type: %s (ID: %s)\n", mt.Name, mt.ID)
		fmt.Printf("  VCPUs: %d\n", mt.VCPUs)
		fmt.Printf("  RAM: %d MB\n", mt.RAM)
		fmt.Printf("  Disk: %d GB\n", mt.Disk)
		fmt.Printf("  GPU: %d\n", mt.GPU)
		fmt.Printf("  Status: %s\n", mt.Status)
	}
}

func ExampleListImages() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}

	c := client.NewMgcClient(apiToken)
	computeClient := compute.New(c)

	// List images
	images, err := computeClient.Images().List(context.Background(), compute.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Print image details
	for _, img := range images {
		fmt.Printf("Image: %s (ID: %s)\n", img.Name, img.ID)
		fmt.Printf("  Status: %s\n", img.Status)
		fmt.Printf("  Version: %s\n", *img.Version)
		fmt.Printf("  Platform: %s\n", *img.Platform)
		fmt.Printf("  Release At: %s\n", *img.ReleaseAt)
		fmt.Printf("  End Standard Support At: %s\n", *img.EndStandardSupportAt)
		fmt.Printf("  End Life At: %s\n", *img.EndLifeAt)
		fmt.Printf("  Minimum Requirements: %d VCPUs, %d RAM, %d Disk\n", img.MinimumRequirements.VCPU, img.MinimumRequirements.RAM, img.MinimumRequirements.Disk)
	}
}
