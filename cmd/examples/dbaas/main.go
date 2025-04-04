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
	ExampleListEngines()
	ExampleListInstanceTypes()
	ExampleListInstances()
	ExampleCreateInstance()
	ExampleListClusters()
	ExampleCreateCluster()
	ExampleGetCluster()
	ExampleUpdateCluster()
	ExampleCreateParameterGroup()
	ExampleGetParameterGroup()
	ExampleUpdateParameterGroup()
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
		Name:           "example-db-instance",
		EngineID:       helpers.StrPtr("your-engine-id"),        // Replace with actual engine ID
		InstanceTypeID: helpers.StrPtr("your-instance-type-id"), // Replace with actual instance type ID
		User:           "dbadmin",
		Password:       "YourStrongPassword123!",
		Volume: dbaas.InstanceVolumeRequest{
			Size: 20, // Size in GB
			Type: "NVME",
		},
		BackupStartAt: helpers.StrPtr("02:00"), // Start backup at 2 AM
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created database instance with ID: %s\n", instance.ID)
}

func ExampleListClusters() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	clusters, err := dbaasClient.Clusters().List(context.Background(), dbaas.ListClustersOptions{
		Limit:  helpers.IntPtr(10),
		Status: (*dbaas.ClusterStatus)(helpers.StrPtr(string(dbaas.ClusterStatusActive))),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d database clusters:\n", len(clusters))
	for _, cluster := range clusters {
		fmt.Printf("Cluster: %s (ID: %s)\n", cluster.Name, cluster.ID)
		fmt.Printf("  Engine ID: %s\n", cluster.EngineID)
		fmt.Printf("  Status: %s\n", cluster.Status)
		fmt.Printf("  Volume Size: %d GB\n", cluster.Volume.Size)
		fmt.Printf("  Volume Type: %s\n", cluster.Volume.Type)
		fmt.Printf("  Parameter Group ID: %s\n", cluster.ParameterGroupID)
		fmt.Printf("  Backup Retention: %d days\n", cluster.BackupRetentionDays)
		fmt.Printf("  Backup Start At: %s\n", cluster.BackupStartAt)

		if len(cluster.Addresses) > 0 {
			fmt.Println("  Addresses:")
			for _, addr := range cluster.Addresses {
				fmt.Printf("    %s: %s:%s\n", addr.Access, addr.Address, addr.Port)
			}
		}
	}
}

func ExampleCreateCluster() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	// Create a new database cluster
	volumeType := "NVME"
	paramGroupID := "your-parameter-group-id" // Replace with actual parameter group ID
	backupRetention := 7

	cluster, err := dbaasClient.Clusters().Create(context.Background(), dbaas.ClusterCreateRequest{
		Name:           "example-db-cluster",
		EngineID:       "your-engine-id",        // Replace with actual engine ID
		InstanceTypeID: "your-instance-type-id", // Replace with actual instance type ID
		User:           "dbadmin",
		Password:       "YourStrongPassword123!",
		Volume: dbaas.ClusterVolumeRequest{
			Size: 50, // Size in GB
			Type: &volumeType,
		},
		ParameterGroupID:    &paramGroupID,
		BackupRetentionDays: &backupRetention,
		BackupStartAt:       helpers.StrPtr("03:00"), // Start backup at 3 AM
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created database cluster with ID: %s\n", cluster.ID)
}

func ExampleGetCluster() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	clusterID := "your-cluster-id" // Replace with actual cluster ID
	cluster, err := dbaasClient.Clusters().Get(context.Background(), clusterID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cluster Details for %s (ID: %s):\n", cluster.Name, cluster.ID)
	fmt.Printf("  Status: %s\n", cluster.Status)
	fmt.Printf("  Engine ID: %s\n", cluster.EngineID)
	fmt.Printf("  Instance Type ID: %s\n", cluster.InstanceTypeID)
	fmt.Printf("  Parameter Group ID: %s\n", cluster.ParameterGroupID)
	fmt.Printf("  Volume Size: %d GB\n", cluster.Volume.Size)
	fmt.Printf("  Created At: %s\n", cluster.CreatedAt)
	fmt.Printf("  Apply Parameters Pending: %v\n", cluster.ApplyParametersPending)
}

func ExampleUpdateCluster() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	clusterID := "your-cluster-id" // Replace with actual cluster ID
	newParamGroupID := "new-parameter-group-id"
	newBackupRetention := 14

	updatedCluster, err := dbaasClient.Clusters().Update(context.Background(), clusterID, dbaas.ClusterUpdateRequest{
		ParameterGroupID:    &newParamGroupID,
		BackupRetentionDays: &newBackupRetention,
		BackupStartAt:       helpers.StrPtr("04:30"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully updated cluster %s\n", updatedCluster.ID)
	fmt.Printf("  New Parameter Group ID: %s\n", updatedCluster.ParameterGroupID)
	fmt.Printf("  New Backup Retention: %d days\n", updatedCluster.BackupRetentionDays)
	fmt.Printf("  New Backup Start Time: %s\n", updatedCluster.BackupStartAt)
}

// Example for parameter groups
func ExampleListParameters() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	paramGroups, err := dbaasClient.Parameters().List(context.Background(), dbaas.ListParameterGroupsOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d parameter groups:\n", len(paramGroups))
	for _, pg := range paramGroups {
		fmt.Printf("Parameter Group: %s (ID: %s)\n", pg.Name, pg.ID)
		fmt.Printf("  Description: %s\n", pg.Description)
		fmt.Printf("  Type: %s\n", pg.Type)
		fmt.Printf("  Engine ID: %s\n", pg.EngineID)
	}
}

func ExampleCreateParameterGroup() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	description := "Custom parameter group for MySQL production databases"

	paramGroup, err := dbaasClient.Parameters().Create(context.Background(), dbaas.ParameterGroupCreateRequest{
		Name:        "mysql-production-params",
		EngineID:    "your-engine-id", // Replace with actual engine ID
		Description: &description,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created parameter group with ID: %s\n", paramGroup.ID)
}

func ExampleGetParameterGroup() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	paramGroupID := "your-parameter-group-id" // Replace with actual parameter group ID
	paramGroup, err := dbaasClient.Parameters().Get(context.Background(), paramGroupID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parameter Group Details for %s (ID: %s):\n", paramGroup.Name, paramGroup.ID)
	fmt.Printf("  Description: %s\n", paramGroup.Description)
	fmt.Printf("  Type: %s\n", paramGroup.Type)
	fmt.Printf("  Engine ID: %s\n", paramGroup.EngineID)
}

func ExampleUpdateParameterGroup() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(apiToken)
	dbaasClient := dbaas.New(c)

	paramGroupID := "your-parameter-group-id" // Replace with actual parameter group ID
	newName := "mysql-optimized-params"
	newDescription := "Optimized parameter group for MySQL high-traffic workloads"

	updatedParamGroup, err := dbaasClient.Parameters().Update(context.Background(), paramGroupID, dbaas.ParameterGroupUpdateRequest{
		Name:        &newName,
		Description: &newDescription,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully updated parameter group %s\n", updatedParamGroup.ID)
	fmt.Printf("  New Name: %s\n", updatedParamGroup.Name)
	fmt.Printf("  New Description: %s\n", updatedParamGroup.Description)
}
