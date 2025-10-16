package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MagaluCloud/mgc-sdk-go/blockstorage"
	"github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/MagaluCloud/mgc-sdk-go/helpers"
)

func main() {
	ExampleListVolumeTypes()
	ExampleListVolumes()
	id := ExampleCreateVolume()
	ExampleGetVolume(id)
	ExampleManageVolume(id)
	ExampleVolumeAttachments(id)
	ExampleSchedulers(id)
	ExampleDeleteVolume(id)
}

func ExampleGetVolume(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)

	volume, err := blockClient.Volumes().Get(context.Background(), id, []string{blockstorage.VolumeTypeExpand, blockstorage.VolumeAttachExpand})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Volume: %s (ID: %s)\n", volume.Name, volume.ID)
}

func ExampleListVolumes() {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)

	// List volumes with pagination and expansion (paginated)
	resp, err := blockClient.Volumes().List(context.Background(), blockstorage.ListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Expand: []string{blockstorage.VolumeTypeExpand, blockstorage.VolumeAttachExpand},
	})

	if err != nil {
		log.Fatal(err)
	}

	// Print pagination metadata
	fmt.Printf("Volumes (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	// Print volume details
	for _, vol := range resp.Volumes {
		fmt.Printf("Volume: %s (ID: %s)\n", vol.Name, vol.ID)
		fmt.Printf("  Size: %d GiB\n", vol.Size)
		if vol.Type.Name != nil {
			fmt.Printf("  Type: %s\n", *vol.Type.Name)
		}
		fmt.Printf("  Status: %s\n", vol.Status)
		fmt.Printf("  State: %s\n", vol.State)
		fmt.Printf("  Encrypted: %v\n", vol.Encrypted)
		fmt.Printf("  Created At: %s\n", vol.CreatedAt)

		if vol.Attachment != nil {
			fmt.Printf("  Attached to: %s\n", *vol.Attachment.Instance.ID)
		}
	}
}

func ExampleListAllVolumes() {
	// Create a new client
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)

	// List all volumes (fetches all pages automatically)
	volumes, err := blockClient.Volumes().ListAll(context.Background(), []blockstorage.VolumeExpand{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total Volumes: %d\n", len(volumes))

	// Print volume details
	for _, vol := range volumes {
		fmt.Printf("Volume: %s (ID: %s) - Size: %d GiB, Status: %s\n",
			vol.Name, vol.ID, vol.Size, vol.Status)
	}
}

func ExampleCreateVolume() string {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)

	// Create a new volume
	createReq := blockstorage.CreateVolumeRequest{
		Name: "my-test-volume",
		Size: 10,
		Type: blockstorage.IDOrName{
			Name: helpers.StrPtr("cloud_nvme1k"),
		},
		Encrypted: helpers.BoolPtr(true),
	}

	id, err := blockClient.Volumes().Create(context.Background(), createReq)
	if err != nil {
		var httpError *client.HTTPError
		if errors.As(err, &httpError) {
			fmt.Printf("Failed to create volume, status %s\n", httpError.Status)
			fmt.Printf("Error body: %s\n", string(httpError.Body))
		}
	}

	fmt.Printf("Created volume with ID: %s\n", id)
	return id
}

func ExampleManageVolume(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)
	ctx := context.Background()

	// Get volume details
	volume, err := blockClient.Volumes().Get(ctx, id, []string{blockstorage.VolumeTypeExpand})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Volume: %s (ID: %s)\n", volume.Name, volume.ID)

	// Rename the volume
	if err := blockClient.Volumes().Rename(ctx, volume.ID, "new-volume-name"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Volume renamed successfully")

	// Extend volume size
	extendReq := blockstorage.ExtendVolumeRequest{
		Size: 200,
	}
	if err := blockClient.Volumes().Extend(ctx, volume.ID, extendReq); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Volume size extended successfully")

	// Change volume type
	retypeReq := blockstorage.RetypeVolumeRequest{
		NewType: blockstorage.IDOrName{
			Name: helpers.StrPtr("cloud_nvme1k"),
		},
	}
	if err := blockClient.Volumes().Retype(ctx, volume.ID, retypeReq); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Volume type changed successfully")
}

func ExampleVolumeAttachments(volumeID string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)
	ctx := context.Background()

	// Replace with an actual instance ID
	instanceID := "your-instance-id"

	// Attach volume to instance
	if err := blockClient.Volumes().Attach(ctx, volumeID, instanceID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Volume %s attached to instance %s\n", volumeID, instanceID)

	// Get volume details with attachment info
	volume, err := blockClient.Volumes().Get(ctx, volumeID, []string{blockstorage.VolumeAttachExpand})
	if err != nil {
		log.Fatal(err)
	}

	if volume.Attachment != nil {
		fmt.Printf("Attachment details:\n")
		fmt.Printf("  Instance: %s\n", *volume.Attachment.Instance.ID)
		fmt.Printf("  Device: %s\n", *volume.Attachment.Device)
		fmt.Printf("  Attached At: %s\n", volume.Attachment.AttachedAt)
	}

	// Detach volume
	if err := blockClient.Volumes().Detach(ctx, volumeID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Volume %s detached successfully\n", volumeID)
}

func ExampleDeleteVolume(id string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)

	if err := blockClient.Volumes().Delete(context.Background(), id); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Volume deleted successfully")
}

func ExampleListVolumeTypes() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)

	// List volume types (paginated)
	resp, err := blockClient.VolumeTypes().List(context.Background(), blockstorage.ListVolumeTypesOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Print pagination metadata
	fmt.Printf("Volume Types (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	// Print volume type details
	for _, vt := range resp.Types {
		fmt.Printf("Volume Type: %s (ID: %s)\n", vt.Name, vt.ID)
		fmt.Printf("  Disk Type: %s\n", vt.DiskType)
		fmt.Printf("  Status: %s\n", vt.Status)
		fmt.Printf("  IOPS: Read=%d, Write=%d, Total=%d\n", vt.IOPS.Read, vt.IOPS.Write, vt.IOPS.Total)
		fmt.Printf("  Availability Zones: %v\n", vt.AvailabilityZones)
		fmt.Printf("  Allows Encryption: %v\n", vt.AllowsEncryption)
	}
}

func ExampleListAllVolumeTypes() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)

	// List all volume types
	volumeTypes, err := blockClient.VolumeTypes().ListAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total Volume Types: %d\n", len(volumeTypes))

	// Print volume type details
	for _, vt := range volumeTypes {
		fmt.Printf("Volume Type: %s (ID: %s)\n", vt.Name, vt.ID)
		fmt.Printf("  Disk Type: %s\n", vt.DiskType)
		fmt.Printf("  Status: %s\n", vt.Status)
		fmt.Printf("  IOPS: Read=%d, Write=%d, Total=%d\n", vt.IOPS.Read, vt.IOPS.Write, vt.IOPS.Total)
		fmt.Printf("  Availability Zones: %v\n", vt.AvailabilityZones)
		fmt.Printf("  Allows Encryption: %v\n", vt.AllowsEncryption)
	}
}

func ExampleSchedulers(volumeID string) {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	blockClient := blockstorage.New(c)
	ctx := context.Background()

	fmt.Println("\n=== Scheduler Examples ===")

	// Create a new scheduler
	schedulerReq := blockstorage.SchedulerPayload{
		Name:        "daily-backup-scheduler",
		Description: helpers.StrPtr("Daily backup scheduler for important volumes"),
		Snapshot: blockstorage.SnapshotConfig{
			Type: "instant",
		},
		Policy: blockstorage.Policy{
			RetentionInDays: 7,
			Frequency: blockstorage.Frequency{
				Daily: blockstorage.DailyFrequency{
					StartTime: "02:00:00",
				},
			},
		},
	}

	schedulerID, err := blockClient.Schedulers().Create(ctx, schedulerReq)
	if err != nil {
		log.Printf("Failed to create scheduler: %v", err)
		return
	}
	fmt.Printf("Created scheduler with ID: %s\n", schedulerID)

	// Get scheduler details
	scheduler, err := blockClient.Schedulers().Get(ctx, schedulerID, []blockstorage.ExpandSchedulers{blockstorage.ExpandSchedulersVolume})
	if err != nil {
		log.Printf("Failed to get scheduler: %v", err)
		return
	}
	fmt.Printf("Scheduler: %s (ID: %s)\n", scheduler.Name, scheduler.ID)
	fmt.Printf("  Description: %v\n", scheduler.Description)
	fmt.Printf("  State: %s\n", scheduler.State)
	fmt.Printf("  Retention: %d days\n", scheduler.Policy.RetentionInDays)
	fmt.Printf("  Start Time: %s\n", scheduler.Policy.Frequency.Daily.StartTime)
	fmt.Printf("  Created At: %s\n", scheduler.CreatedAt)

	// Attach volume to scheduler
	attachReq := blockstorage.SchedulerVolumeIdentifierPayload{
		Volume: blockstorage.IDOrName{
			ID: &volumeID,
		},
	}
	if err := blockClient.Schedulers().AttachVolume(ctx, schedulerID, attachReq); err != nil {
		log.Printf("Failed to attach volume to scheduler: %v", err)
	} else {
		fmt.Printf("Volume %s attached to scheduler %s\n", volumeID, schedulerID)
	}

	// List schedulers with expansion
	schedulerList, err := blockClient.Schedulers().List(ctx, blockstorage.SchedulerListOptions{
		Limit:  helpers.IntPtr(10),
		Offset: helpers.IntPtr(0),
		Expand: []blockstorage.ExpandSchedulers{blockstorage.ExpandSchedulersVolume},
	})
	if err != nil {
		log.Printf("Failed to list schedulers: %v", err)
		return
	}

	fmt.Printf("\nFound %d schedulers:\n", len(schedulerList.Schedulers))
	for _, sched := range schedulerList.Schedulers {
		fmt.Printf("  Scheduler: %s (ID: %s)\n", sched.Name, sched.ID)
		fmt.Printf("    State: %s\n", sched.State)
		fmt.Printf("    Attached Volumes: %d\n", len(sched.Volumes))
		for _, volID := range sched.Volumes {
			fmt.Printf("      - %s\n", volID)
		}
	}

	// Wait a moment to demonstrate functionality
	time.Sleep(2 * time.Second)

	// Detach volume from scheduler
	detachReq := blockstorage.SchedulerVolumeIdentifierPayload{
		Volume: blockstorage.IDOrName{
			ID: &volumeID,
		},
	}
	if err := blockClient.Schedulers().DetachVolume(ctx, schedulerID, detachReq); err != nil {
		log.Printf("Failed to detach volume from scheduler: %v", err)
	} else {
		fmt.Printf("Volume %s detached from scheduler %s\n", volumeID, schedulerID)
	}

	// Delete the scheduler
	if err := blockClient.Schedulers().Delete(ctx, schedulerID); err != nil {
		log.Printf("Failed to delete scheduler: %v", err)
	} else {
		fmt.Printf("Scheduler %s deleted successfully\n", schedulerID)
	}
}
