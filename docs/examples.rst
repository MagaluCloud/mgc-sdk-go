Examples
========

This section provides comprehensive examples of how to use the MGC Go SDK for common tasks. Each example is complete and ready to run.

Basic Setup
----------

All examples start with this basic setup:

.. code-block:: go

   package main

   import (
       "context"
       "log"
       "os"
       "time"

       "github.com/MagaluCloud/mgc-sdk-go/client"
       "github.com/MagaluCloud/mgc-sdk-go/helpers"
   )

   func main() {
       // Initialize the client
       apiToken := os.Getenv("MGC_API_TOKEN")
       if apiToken == "" {
           log.Fatal("MGC_API_TOKEN environment variable is required")
       }

       c := client.NewMgcClient(
           apiToken,
           client.WithTimeout(30*time.Second),
           client.WithBaseURL(client.BrSe1),
       )

       // Your code here...
   }

Compute Examples
----------------

List Instances
~~~~~~~~~~~~~

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/compute"

   func listInstances(ctx context.Context, c *client.CoreClient) error {
       computeClient := compute.New(c)
       
       instances, err := computeClient.Instances().List(ctx, compute.ListOptions{
           Limit:  helpers.IntPtr(10),
           Offset: helpers.IntPtr(0),
           Expand: []string{compute.InstanceMachineTypeExpand, compute.InstanceImageExpand},
       })
       if err != nil {
           return err
       }

       for _, instance := range instances {
           log.Printf("Instance: %s (ID: %s, Status: %s)", 
               *instance.Name, *instance.ID, instance.Status)
       }
       
       return nil
   }

Create an Instance
~~~~~~~~~~~~~~~~~

.. code-block:: go

   func createInstance(ctx context.Context, c *client.CoreClient, name string) (string, error) {
       computeClient := compute.New(c)
       
       createReq := compute.CreateRequest{
           Name: name,
           MachineType: compute.IDOrName{
               Name: helpers.StrPtr("BV1-1-40"),
           },
           Image: compute.IDOrName{
               Name: helpers.StrPtr("cloud-ubuntu-24.04 LTS"),
           },
           Network: &compute.CreateParametersNetwork{
               AssociatePublicIp: helpers.BoolPtr(false),
           },
           SshKeyName: helpers.StrPtr("my-ssh-key"),
       }

       id, err := computeClient.Instances().Create(ctx, createReq)
       if err != nil {
           return "", err
       }

       log.Printf("Created instance with ID: %s", id)
       return id, nil
   }

Delete an Instance
~~~~~~~~~~~~~~~~~

.. code-block:: go

   func deleteInstance(ctx context.Context, c *client.CoreClient, id string) error {
       computeClient := compute.New(c)
       
       err := computeClient.Instances().Delete(ctx, id)
       if err != nil {
           return err
       }

       log.Printf("Deleted instance: %s", id)
       return nil
   }

Network Examples
----------------

Create a VPC
~~~~~~~~~~~

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/network"

   func createVPC(ctx context.Context, c *client.CoreClient, name string) (string, error) {
       networkClient := network.New(c)
       
       req := network.CreateVPCRequest{
           Name:        name,
           Description: helpers.StrPtr("My VPC for testing"),
       }

       id, err := networkClient.VPCs().Create(ctx, req)
       if err != nil {
           return "", err
       }

       log.Printf("Created VPC with ID: %s", id)
       return id, nil
   }

Create a Subnet
~~~~~~~~~~~~~~

.. code-block:: go

   func createSubnet(ctx context.Context, c *client.CoreClient, vpcID, name string) (string, error) {
       networkClient := network.New(c)
       
       req := network.SubnetCreateRequest{
           Name:        name,
           Description: helpers.StrPtr("My subnet"),
           CIDRBlock:   "10.0.1.0/24",
           IPVersion:   4,
       }

       opts := network.SubnetCreateOptions{
           Zone: helpers.StrPtr("br-se1-1"),
       }

       id, err := networkClient.VPCs().CreateSubnet(ctx, vpcID, req, opts)
       if err != nil {
           return "", err
       }

       log.Printf("Created subnet with ID: %s", id)
       return id, nil
   }

Create a Security Group
~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   func createSecurityGroup(ctx context.Context, c *client.CoreClient, name string) (string, error) {
       networkClient := network.New(c)
       
       req := network.SecurityGroupCreateRequest{
           Name:        name,
           Description: helpers.StrPtr("My security group"),
       }

       id, err := networkClient.SecurityGroups().Create(ctx, req)
       if err != nil {
           return "", err
       }

       log.Printf("Created security group with ID: %s", id)
       return id, nil
   }

Block Storage Examples
---------------------

Create a Volume
~~~~~~~~~~~~~~

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/blockstorage"

   func createVolume(ctx context.Context, c *client.CoreClient, name string, size int) (string, error) {
       storageClient := blockstorage.New(c)
       
       req := blockstorage.CreateVolumeRequest{
           Name:        name,
           Description: helpers.StrPtr("My volume"),
           Size:        size,
           VolumeType:  helpers.StrPtr("standard"),
       }

       id, err := storageClient.Volumes().Create(ctx, req)
       if err != nil {
           return "", err
       }

       log.Printf("Created volume with ID: %s", id)
       return id, nil
   }

Attach Volume to Instance
~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   func attachVolume(ctx context.Context, c *client.CoreClient, volumeID, instanceID string) error {
       storageClient := blockstorage.New(c)
       
       req := blockstorage.AttachVolumeRequest{
           InstanceID: instanceID,
       }

       err := storageClient.Volumes().Attach(ctx, volumeID, req)
       if err != nil {
           return err
       }

       log.Printf("Attached volume %s to instance %s", volumeID, instanceID)
       return nil
   }

Database Examples
-----------------

Create a Database Instance
~~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/dbaas"

   func createDatabase(ctx context.Context, c *client.CoreClient, name string) (string, error) {
       dbaasClient := dbaas.New(c)
       
       req := dbaas.CreateInstanceRequest{
           Name:        name,
           Description: helpers.StrPtr("My database"),
           Engine:      "mysql",
           Version:     "8.0",
           InstanceType: dbaas.IDOrName{
               Name: helpers.StrPtr("db.t3.micro"),
           },
           Storage: &dbaas.Storage{
               Size: 20,
               Type: "gp2",
           },
       }

       id, err := dbaasClient.Instances().Create(ctx, req)
       if err != nil {
           return "", err
       }

       log.Printf("Created database with ID: %s", id)
       return id, nil
   }

Container Registry Examples
--------------------------

Create a Repository
~~~~~~~~~~~~~~~~~~

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/containerregistry"

   func createRepository(ctx context.Context, c *client.CoreClient, name string) (string, error) {
       registryClient := containerregistry.New(c)
       
       req := containerregistry.CreateRepositoryRequest{
           Name:        name,
           Description: helpers.StrPtr("My repository"),
           Visibility:  helpers.StrPtr("private"),
       }

       id, err := registryClient.Repositories().Create(ctx, req)
       if err != nil {
           return "", err
       }

       log.Printf("Created repository with ID: %s", id)
       return id, nil
   }

Kubernetes Examples
-------------------

Create a Cluster
~~~~~~~~~~~~~~~

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/kubernetes"

   func createCluster(ctx context.Context, c *client.CoreClient, name string) (string, error) {
       k8sClient := kubernetes.New(c)
       
       req := kubernetes.CreateClusterRequest{
           Name:        name,
           Description: helpers.StrPtr("My Kubernetes cluster"),
           Version:     "1.28",
           Flavor: kubernetes.IDOrName{
               Name: helpers.StrPtr("k8s-small"),
           },
       }

       id, err := k8sClient.Clusters().Create(ctx, req)
       if err != nil {
           return "", err
       }

       log.Printf("Created cluster with ID: %s", id)
       return id, nil
   }

SSH Keys Examples
-----------------

Create an SSH Key
~~~~~~~~~~~~~~~~

.. code-block:: go

   import "github.com/MagaluCloud/mgc-sdk-go/sshkeys"

   func createSSHKey(ctx context.Context, c *client.CoreClient, name, publicKey string) (*sshkeys.SSHKey, error) {
       sshClient := sshkeys.New(c)
       
       req := sshkeys.CreateSSHKeyRequest{
           Name: name,
           Key:  publicKey,
       }

       key, err := sshClient.Keys().Create(ctx, req)
       if err != nil {
           return nil, err
       }

       log.Printf("Created SSH key with ID: %s", key.ID)
       return key, nil
   }

Complete Example: Infrastructure Setup
------------------------------------

Here's a complete example that sets up a basic infrastructure:

.. code-block:: go

   package main

   import (
       "context"
       "log"
       "os"
       "time"

       "github.com/MagaluCloud/mgc-sdk-go/client"
       "github.com/MagaluCloud/mgc-sdk-go/compute"
       "github.com/MagaluCloud/mgc-sdk-go/network"
       "github.com/MagaluCloud/mgc-sdk-go/helpers"
   )

   func main() {
       // Initialize client
       apiToken := os.Getenv("MGC_API_TOKEN")
       if apiToken == "" {
           log.Fatal("MGC_API_TOKEN environment variable is required")
       }

       c := client.NewMgcClient(
           apiToken,
           client.WithTimeout(30*time.Second),
           client.WithBaseURL(client.BrSe1),
       )

       ctx := context.Background()

       // Create infrastructure
       if err := createInfrastructure(ctx, c); err != nil {
           log.Fatalf("Failed to create infrastructure: %v", err)
       }

       log.Println("Infrastructure created successfully!")
   }

   func createInfrastructure(ctx context.Context, c *client.CoreClient) error {
       networkClient := network.New(c)
       computeClient := compute.New(c)

       // 1. Create VPC
       log.Println("Creating VPC...")
       vpcReq := network.CreateVPCRequest{
           Name:        "my-vpc",
           Description: helpers.StrPtr("VPC for my application"),
       }
       vpcID, err := networkClient.VPCs().Create(ctx, vpcReq)
       if err != nil {
           return err
       }

       // 2. Create subnet
       log.Println("Creating subnet...")
       subnetReq := network.SubnetCreateRequest{
           Name:        "my-subnet",
           Description: helpers.StrPtr("Subnet for my application"),
           CIDRBlock:   "10.0.1.0/24",
           IPVersion:   4,
       }
       subnetOpts := network.SubnetCreateOptions{
           Zone: helpers.StrPtr("br-se1-1"),
       }
       _, err = networkClient.VPCs().CreateSubnet(ctx, vpcID, subnetReq, subnetOpts)
       if err != nil {
           return err
       }

       // 3. Create security group
       log.Println("Creating security group...")
       sgReq := network.SecurityGroupCreateRequest{
           Name:        "my-security-group",
           Description: helpers.StrPtr("Security group for my application"),
       }
       sgID, err := networkClient.SecurityGroups().Create(ctx, sgReq)
       if err != nil {
           return err
       }

       // 4. Create instance
       log.Println("Creating instance...")
       instanceReq := compute.CreateRequest{
           Name: "my-instance",
           MachineType: compute.IDOrName{
               Name: helpers.StrPtr("BV1-1-40"),
           },
           Image: compute.IDOrName{
               Name: helpers.StrPtr("cloud-ubuntu-24.04 LTS"),
           },
           Network: &compute.CreateParametersNetwork{
               AssociatePublicIp: helpers.BoolPtr(false),
           },
           SshKeyName: helpers.StrPtr("my-ssh-key"),
       }
       instanceID, err := computeClient.Instances().Create(ctx, instanceReq)
       if err != nil {
           return err
       }

       log.Printf("Created instance with ID: %s", instanceID)
       return nil
   }

Error Handling Example
---------------------

Here's an example with comprehensive error handling:

.. code-block:: go

   func createInstanceWithErrorHandling(ctx context.Context, c *client.CoreClient, name string) (string, error) {
       computeClient := compute.New(c)
       
       req := compute.CreateRequest{
           Name: name,
           MachineType: compute.IDOrName{
               Name: helpers.StrPtr("BV1-1-40"),
           },
           Image: compute.IDOrName{
               Name: helpers.StrPtr("cloud-ubuntu-24.04 LTS"),
           },
       }

       id, err := computeClient.Instances().Create(ctx, req)
       if err != nil {
           return "", handleCreateError(err, name)
       }

       return id, nil
   }

   func handleCreateError(err error, name string) error {
       if httpErr, ok := err.(*client.HTTPError); ok {
           switch httpErr.StatusCode {
           case 400:
               return fmt.Errorf("invalid request for instance %s: %s", name, string(httpErr.Body))
           case 403:
               return fmt.Errorf("insufficient permissions to create instance %s", name)
           case 409:
               return fmt.Errorf("instance %s already exists", name)
           case 429:
               return fmt.Errorf("rate limit exceeded while creating instance %s", name)
           default:
               return fmt.Errorf("HTTP %d error creating instance %s: %s", httpErr.StatusCode, name, string(httpErr.Body))
           }
       }
       
       if validErr, ok := err.(*client.ValidationError); ok {
           return fmt.Errorf("validation error for instance %s: field %s - %s", name, validErr.Field, validErr.Message)
       }
       
       return fmt.Errorf("unexpected error creating instance %s: %v", name, err)
   }

Running Examples
---------------

To run any of these examples:

1. Set your API token:
   .. code-block:: bash
      export MGC_API_TOKEN="your-api-token-here"

2. Create a Go file with the example code

3. Run the example:
   .. code-block:: bash
      go run main.go

For more examples, check the `cmd/examples <https://github.com/MagaluCloud/mgc-sdk-go/tree/main/cmd/examples>`_ directory in the repository. 