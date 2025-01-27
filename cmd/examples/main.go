package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/compute"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func main() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	sdkClient := initClient(apiToken)
	computeClient := initComputeClient(sdkClient)

	ExampleListMachineTypes(computeClient)
	ExampleListImages(computeClient)
	ExampleListInstances(computeClient)

	id := ExampleCreateInstance(computeClient)
	ExampleManageInstance(computeClient, id)
	ExampleDeleteInstance(computeClient, id)
}

func initClient(apiToken string) *client.CoreClient {
	return client.NewMgcClient(apiToken)
}

func initComputeClient(sdkClient *client.CoreClient) *compute.VirtualMachineClient {
	return compute.New(sdkClient)
}

func ExampleListInstances(computeClient *compute.VirtualMachineClient) {
	// List instances with pagination and sorting
	instances, err := computeClient.Instances().List(context.Background(), compute.ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Expand: []string{compute.InstanceMachineTypeExpand, compute.InstanceImageExpand},
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

func ExampleCreateInstance(computeClient *compute.VirtualMachineClient) string {
	// Create a new instance
	createReq := compute.CreateRequest{
		Name: "my-test-vm-" + strconv.Itoa(int(time.Now().Unix())),
		MachineType: compute.IDOrName{
			Name: helpers.StrPtr("BV1-1-20"),
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

	for {
		instance, err := computeClient.Instances().Get(context.Background(), id, []string{compute.InstanceNetworkExpand})
		if err != nil {
			log.Fatal(err)
		}
		if instance.State == "running" {
			break
		}
		time.Sleep(5 * time.Second)
	}

	fmt.Printf("Created instance with ID: %s\n", id)

	return id
}

func ExampleManageInstance(computeClient *compute.VirtualMachineClient, id string) {
	ctx := context.Background()

	// Get instance details
	instance, err := computeClient.Instances().Get(ctx, id, []string{compute.InstanceNetworkExpand})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Instance: %s (ID: %s)\n", instance.Name, instance.ID)

	// Rename the instance
	if err := computeClient.Instances().Rename(ctx, instance.ID, "new-name"+strconv.Itoa(int(time.Now().Unix()))); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Instance renamed successfully")

	// Change machine type
	retypeReq := compute.RetypeRequest{
		MachineType: compute.IDOrName{
			Name: helpers.StrPtr("BV1-1-40"),
		},
	}
	if err := computeClient.Instances().Retype(ctx, instance.ID, retypeReq); err != nil {
		log.Fatal(err)
	}
	for {
		instance, err := computeClient.Instances().Get(context.Background(), id, []string{compute.InstanceNetworkExpand})
		if err != nil {
			log.Fatal(err)
		}
		if instance.Status == "completed" {
			break
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Instance machine type changed successfully")
}

func ExampleDeleteInstance(computeClient *compute.VirtualMachineClient, id string) {
	// Delete instance and its public IP
	if err := computeClient.Instances().Delete(context.Background(), id, true); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Instance deleted successfully")
}

func ExampleListMachineTypes(computeClient *compute.VirtualMachineClient) {
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

func ExampleListImages(computeClient *compute.VirtualMachineClient) {
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
