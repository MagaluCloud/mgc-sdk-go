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
	ExampleListAllEngines()
	ExampleListInstanceTypes()
	ExampleListAllInstanceTypes()
	ExampleListInstances()
	ExampleListAllInstances()
	ExampleCreateInstance()
	ExampleUpdateInstance()
	ExampleGetInstance()
	ExampleListClusters()
	ExampleListAllClusters()
	ExampleCreateCluster()
	ExampleGetCluster()
	ExampleUpdateCluster()
	ExampleCreateParameterGroup()
	ExampleGetParameterGroup()
	ExampleUpdateParameterGroup()
	ExampleListParametersGroup()
	ExampleListAllParametersGroup()
	ExampleListParameters()
	ExampleListAllParameters()
	ExampleCreateParameter()
	ExampleUpdateParameter()
	ExampleDeleteParameter()
}

func ExampleListEngines() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	resp, err := dbaasClient.Engines().List(context.Background(), dbaas.ListEngineOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print pagination metadata
	fmt.Printf("Database Engines (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	for _, engine := range resp.Results {
		fmt.Printf("Engine: %s (ID: %s)\n", engine.Name, engine.ID)
		fmt.Printf("  Version: %s\n", engine.Version)
		fmt.Printf("  Status: %s\n", engine.Status)
	}
}

func ExampleListAllEngines() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	// List all engines across all pages
	engines, err := dbaasClient.Engines().ListAll(context.Background(), dbaas.EngineFilterOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d database engines:\n", len(engines))
	for _, engine := range engines {
		fmt.Printf("Engine: %s (ID: %s) - Version: %s, Status: %s\n",
			engine.Name, engine.ID, engine.Version, engine.Status)
	}
}

func ExampleListInstanceTypes() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	resp, err := dbaasClient.InstanceTypes().List(context.Background(), dbaas.ListInstanceTypeOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print pagination metadata
	fmt.Printf("Instance Types (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	for _, instanceType := range resp.Results {
		fmt.Printf("Instance Type: %s (ID: %s)\n", instanceType.Name, instanceType.ID)
		fmt.Printf("  Label: %s\n", instanceType.Label)
		fmt.Printf("  VCPU: %s\n", instanceType.VCPU)
		fmt.Printf("  RAM: %s\n", instanceType.RAM)
		fmt.Printf("  Family: %s (%s)\n", instanceType.FamilyDescription, instanceType.FamilySlug)
		fmt.Printf("  Compatible Product: %s\n", instanceType.CompatibleProduct)
	}
}

func ExampleListAllInstanceTypes() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	// List all instance types across all pages
	instanceTypes, err := dbaasClient.InstanceTypes().ListAll(context.Background(), dbaas.InstanceTypeFilterOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d instance types:\n", len(instanceTypes))
	for _, instanceType := range instanceTypes {
		fmt.Printf("Instance Type: %s (ID: %s) - VCPU: %s, RAM: %s\n",
			instanceType.Name, instanceType.ID, instanceType.VCPU, instanceType.RAM)
	}
}

func ExampleListInstances() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	resp, err := dbaasClient.Instances().List(context.Background(), dbaas.ListInstanceOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print pagination metadata
	fmt.Printf("Database Instances (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	for _, instance := range resp.Results {
		fmt.Printf("Instance: %s (ID: %s)\n", instance.Name, instance.ID)
		fmt.Printf("  Engine ID: %s\n", instance.EngineID)
		fmt.Printf("  Status: %s\n", instance.Status)
		fmt.Printf("  Volume Size: %d GiB\n", instance.Volume.Size)
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

func ExampleListAllInstances() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	// List all instances across all pages
	instances, err := dbaasClient.Instances().ListAll(context.Background(), dbaas.InstanceFilterOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d database instances:\n", len(instances))
	for _, instance := range instances {
		fmt.Printf("Instance: %s (ID: %s) - Status: %s, Volume: %d GiB\n",
			instance.Name, instance.ID, instance.Status, instance.Volume.Size)
	}
}

func ExampleCreateInstance() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	// Create a new database instance
	instance, err := dbaasClient.Instances().Create(context.Background(), dbaas.InstanceCreateRequest{
		Name:           "example-db-instance",
		EngineID:       helpers.StrPtr("063f3994-b6c2-4c37-96c9-bab8d82d36f7"), // Replace with actual engine ID
		InstanceTypeID: helpers.StrPtr("6111d89a-3bc0-41e6-98c2-fb23bfa5a56a"), // Replace with actual instance type ID
		User:           "dbadmin",
		Password:       "YourStrongPassword123!",
		Volume: dbaas.InstanceVolumeRequest{
			Size: 20, // Size in GiB
			Type: "CLOUD_NVME15K",
		},
		BackupStartAt:     helpers.StrPtr("02:00"), // Start backup at 2 AM
		DeletionProtected: helpers.BoolPtr(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created database instance with ID: %s\n", instance.ID)
}

func ExampleUpdateInstance() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	instanceId := "your-instance-id"

	// Create a new database instance
	instance, err := dbaasClient.Instances().Update(
		context.Background(),
		instanceId,
		dbaas.DatabaseInstanceUpdateRequest{
			DeletionProtected: helpers.BoolPtr(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully updating database instance with ID: %s\n", instance.ID)
}

func ExampleGetInstance() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	instanceId := "your-instance-id"

	// Create a new database instance
	instance, err := dbaasClient.Instances().Get(
		context.Background(),
		instanceId,
		dbaas.GetInstanceOptions{},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s\n", instance.ID)
	fmt.Printf("Status: %s\n", instance.Status)
	fmt.Printf("Volume: %d GiB\n", instance.Volume.Size)
	fmt.Printf("Name: %s\n", instance.Name)
	fmt.Printf("Engine ID: %s\n", instance.EngineID)
	fmt.Printf("Instance Type ID: %s\n", instance.InstanceTypeID)
	fmt.Printf("Deletion Protected: %v\n", instance.DeletionProtected)
	fmt.Printf("Generation: %s\n", instance.Generation)
	fmt.Printf("Parameter Group ID: %s\n", instance.ParameterGroupID)
	fmt.Printf("Availability Zone: %s\n", instance.AvailabilityZone)
	fmt.Printf("Backup Retention Days: %d\n", instance.BackupRetentionDays)
	fmt.Printf("Backup Start at: %s\n", instance.BackupStartAt)
}

func ExampleListClusters() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	resp, err := dbaasClient.Clusters().List(context.Background(), dbaas.ListClustersOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print pagination metadata
	fmt.Printf("Database Clusters (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	for _, cluster := range resp.Results {
		fmt.Printf("Cluster: %s (ID: %s)\n", cluster.Name, cluster.ID)
		fmt.Printf("  Engine ID: %s\n", cluster.EngineID)
		fmt.Printf("  Status: %s\n", cluster.Status)
		fmt.Printf("  Volume Size: %d GiB\n", cluster.Volume.Size)
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

func ExampleListAllClusters() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	// List all clusters across all pages
	clusters, err := dbaasClient.Clusters().ListAll(context.Background(), dbaas.ClusterFilterOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d database clusters:\n", len(clusters))
	for _, cluster := range clusters {
		fmt.Printf("Cluster: %s (ID: %s) - Status: %s, Retention: %d days\n",
			cluster.Name, cluster.ID, cluster.Status, cluster.BackupRetentionDays)
	}
}

func ExampleCreateCluster() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	// Create a new database cluster
	volumeType := "CLOUD_NVME15K"
	// paramGroupID := "your-parameter-group-id" // Replace with actual parameter group ID
	backupRetention := 7

	cluster, err := dbaasClient.Clusters().Create(context.Background(), dbaas.ClusterCreateRequest{
		Name:           "example-db-cluster",
		EngineID:       "063f3994-b6c2-4c37-96c9-bab8d82d36f7", // Replace with actual engine ID
		InstanceTypeID: "8bbe8e01-40c8-4d2b-80e8-189debc44b1c", // Replace with actual instance type ID
		User:           "dbadmin",
		Password:       "YourStrongPassword123!",
		Volume: dbaas.ClusterVolumeRequest{
			Size: 50, // Size in GiB
			Type: &volumeType,
		},
		// ParameterGroupID:    &paramGroupID,
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
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
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
	fmt.Printf("  Volume Size: %d GiB\n", cluster.Volume.Size)
	fmt.Printf("  Created At: %s\n", cluster.CreatedAt)
	fmt.Printf("  Apply Parameters Pending: %v\n", cluster.ApplyParametersPending)
}

func ExampleUpdateCluster() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
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
func ExampleListParametersGroup() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	resp, err := dbaasClient.ParametersGroup().List(context.Background(), dbaas.ListParameterGroupsOptions{
		Limit: helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print pagination metadata
	fmt.Printf("Parameter Groups (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	for _, pg := range resp.Results {
		fmt.Printf("Parameter Group: %s (ID: %s)\n", pg.Name, pg.ID)
		fmt.Printf("  Description: %s\n", *pg.Description)
		fmt.Printf("  Type: %s\n", pg.Type)
		fmt.Printf("  Engine ID: %s\n", pg.EngineID)
	}
}

func ExampleListAllParametersGroup() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	// List all parameter groups across all pages
	paramGroups, err := dbaasClient.ParametersGroup().ListAll(context.Background(), dbaas.ParameterGroupFilterOptions{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d parameter groups:\n", len(paramGroups))
	for _, pg := range paramGroups {
		fmt.Printf("Parameter Group: %s (ID: %s) - Type: %s, Engine: %s\n",
			pg.Name, pg.ID, pg.Type, pg.EngineID)
	}
}

func ExampleCreateParameterGroup() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	description := "Custom parameter group for MySQL production databases"

	paramGroup, err := dbaasClient.ParametersGroup().Create(context.Background(), dbaas.ParameterGroupCreateRequest{
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
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	paramGroupID := "your-parameter-group-id" // Replace with actual parameter group ID
	paramGroup, err := dbaasClient.ParametersGroup().Get(context.Background(), paramGroupID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Parameter Group Details for %s (ID: %s):\n", paramGroup.Name, paramGroup.ID)
	fmt.Printf("  Description: %s\n", *paramGroup.Description)
	fmt.Printf("  Type: %s\n", paramGroup.Type)
	fmt.Printf("  Engine ID: %s\n", paramGroup.EngineID)
}

func ExampleUpdateParameterGroup() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	paramGroupID := "your-parameter-group-id" // Replace with actual parameter group ID
	newName := "mysql-optimized-params"
	newDescription := "Optimized parameter group for MySQL high-traffic workloads"

	updatedParamGroup, err := dbaasClient.ParametersGroup().Update(context.Background(), paramGroupID, dbaas.ParameterGroupUpdateRequest{
		Name:        &newName,
		Description: &newDescription,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully updated parameter group %s\n", updatedParamGroup.ID)
	fmt.Printf("  New Name: %s\n", updatedParamGroup.Name)
	fmt.Printf("  New Description: %s\n", *updatedParamGroup.Description)
}

func ExampleListParameters() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	resp, err := dbaasClient.Parameters().List(context.Background(), dbaas.ListParametersOptions{
		ParameterGroupID: "88bd17e0-779c-43a5-9695-5cb9f6f918c0",
		Limit:            helpers.IntPtr(10),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print pagination metadata
	fmt.Printf("Parameters (Page %d-%d of %d total)\n",
		resp.Meta.Page.Offset,
		resp.Meta.Page.Offset+resp.Meta.Page.Count,
		resp.Meta.Page.Total)

	for _, p := range resp.Results {
		fmt.Printf("Parameter: %s (ID: %s) = %v\n", p.Name, p.ID, p.Value)
	}
}

func ExampleListAllParameters() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	// List all parameters across all pages for a specific parameter group
	params, err := dbaasClient.Parameters().ListAll(context.Background(), dbaas.ParameterFilterOptions{
		ParameterGroupID: "88bd17e0-779c-43a5-9695-5cb9f6f918c0",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d parameters:\n", len(params))
	for _, p := range params {
		fmt.Printf("Parameter: %s (ID: %s) = %v\n", p.Name, p.ID, p.Value)
	}
}

func ExampleCreateParameter() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	created, err := dbaasClient.Parameters().Create(context.Background(),
		"88bd17e0-779c-43a5-9695-5cb9f6f918c0",
		dbaas.ParameterCreateRequest{
			Name:  "LOWER_CASE_TABLE_NAMES",
			Value: 1,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created parameter with ID: %s\n", created.ID)
}

func ExampleUpdateParameter() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	updated, err := dbaasClient.Parameters().Update(context.Background(),
		"88bd17e0-779c-43a5-9695-5cb9f6f918c0",
		"68378760-c4e0-484a-b71a-b900942e7758",
		dbaas.ParameterUpdateRequest{
			Value: 0,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated parameter %s (ID: %s) to %v\n", updated.Name, updated.ID, updated.Value)
}

func ExampleDeleteParameter() {
	apiToken := os.Getenv("MGC_API_TOKEN")
	if apiToken == "" {
		log.Fatal("MGC_API_TOKEN environment variable is not set")
	}
	c := client.NewMgcClient(client.WithAPIKey(apiToken))
	dbaasClient := dbaas.New(c)

	err := dbaasClient.Parameters().Delete(context.Background(),
		"88bd17e0-779c-43a5-9695-5cb9f6f918c0",
		"68378760-c4e0-484a-b71a-b900942e7758",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parameter deleted successfully")
}
